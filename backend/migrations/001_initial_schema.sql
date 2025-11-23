-- 奖励系统数据库迁移脚本
-- 基于 MySQL 8.0+

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS reward_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE reward_system;

-- 家庭表
CREATE TABLE IF NOT EXISTS families (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    family_id BIGINT NOT NULL,
    role ENUM('guardian', 'child') NOT NULL,
    display_name VARCHAR(64) NOT NULL,
    wechat_openid VARCHAR(128) UNIQUE,
    is_active TINYINT(1) DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE,
    INDEX idx_family_role (family_id, role),
    INDEX idx_wechat_openid (wechat_openid),
    INDEX idx_is_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 奖励类型表
CREATE TABLE IF NOT EXISTS reward_types (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    family_id BIGINT NOT NULL,
    name VARCHAR(64) NOT NULL,
    unit_kind ENUM('money', 'time', 'points', 'custom') NOT NULL,
    unit_label VARCHAR(32),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE,
    UNIQUE KEY uniq_family_name (family_id, name),
    INDEX idx_family_id (family_id),
    INDEX idx_unit_kind (unit_kind)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 账户表（孩子在每种奖励类型上的账户）
CREATE TABLE IF NOT EXISTS accounts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    family_id BIGINT NOT NULL,
    child_id BIGINT NOT NULL,
    reward_type_id BIGINT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE,
    FOREIGN KEY (child_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reward_type_id) REFERENCES reward_types(id) ON DELETE CASCADE,
    UNIQUE KEY uniq_child_reward_type (child_id, reward_type_id),
    INDEX idx_family_id (family_id),
    INDEX idx_child_id (child_id),
    INDEX idx_reward_type_id (reward_type_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 交易记录表
CREATE TABLE IF NOT EXISTS transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_id BIGINT NOT NULL,
    type ENUM('credit', 'debit') NOT NULL,
    value BIGINT NOT NULL,
    note VARCHAR(255),
    created_by BIGINT NOT NULL,
    idempotency_key VARCHAR(64),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_account_id (account_id),
    INDEX idx_created_by (created_by),
    INDEX idx_created_at (created_at),
    INDEX idx_idempotency_key (idempotency_key),
    INDEX idx_account_created (account_id, created_at DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 审计日志表
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    family_id BIGINT,
    user_id BIGINT,
    action VARCHAR(32) NOT NULL,
    payload JSON,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    INDEX idx_family_id (family_id),
    INDEX idx_user_id (user_id),
    INDEX idx_action (action),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 插入默认数据
INSERT INTO families (name) VALUES ('默认家庭');

INSERT INTO users (family_id, role, display_name, wechat_openid) VALUES 
(1, 'guardian', '爸爸', 'guardian_openid_001'),
(1, 'child', '小明', 'child_openid_001');

INSERT INTO reward_types (family_id, name, unit_kind, unit_label) VALUES 
(1, '零花钱', 'money', '元'),
(1, '看电视时间', 'time', '分钟'),
(1, '积分', 'points', '积分');

-- 创建账户（为每个孩子和奖励类型组合）
INSERT INTO accounts (family_id, child_id, reward_type_id, balance)
SELECT 1, u.id, rt.id, 0
FROM users u
CROSS JOIN reward_types rt
WHERE u.role = 'child' AND u.family_id = 1 AND rt.family_id = 1;

-- 创建触发器：在授予或消费奖励时自动更新账户余额
DELIMITER //

CREATE TRIGGER update_account_balance_after_transaction
AFTER INSERT ON transactions
FOR EACH ROW
BEGIN
    IF NEW.type = 'credit' THEN
        UPDATE accounts 
        SET balance = balance + NEW.value 
        WHERE id = NEW.account_id;
    ELSEIF NEW.type = 'debit' THEN
        UPDATE accounts 
        SET balance = balance - NEW.value 
        WHERE id = NEW.account_id;
    END IF;
END//

DELIMITER ;

-- 创建存储过程：安全地授予奖励（带余额检查和事务）
DELIMITER //

CREATE PROCEDURE grant_reward(
    IN p_family_id BIGINT,
    IN p_child_id BIGINT,
    IN p_reward_type_id BIGINT,
    IN p_value BIGINT,
    IN p_note VARCHAR(255),
    IN p_created_by BIGINT,
    IN p_idempotency_key VARCHAR(64)
)
BEGIN
    DECLARE v_account_id BIGINT;
    DECLARE v_new_balance BIGINT;
    
    -- 开始事务
    START TRANSACTION;
    
    -- 检查幂等性
    IF p_idempotency_key IS NOT NULL AND p_idempotency_key != '' THEN
        IF EXISTS (SELECT 1 FROM transactions WHERE idempotency_key = p_idempotency_key) THEN
            -- 返回已存在的交易信息
            SELECT id, (SELECT balance FROM accounts WHERE id = account_id) as balance
            FROM transactions 
            WHERE idempotency_key = p_idempotency_key;
            COMMIT;
            LEAVE proc_label;
        END IF;
    END IF;
    
    -- 获取账户ID（如果不存在则创建）
    SELECT id INTO v_account_id 
    FROM accounts 
    WHERE child_id = p_child_id AND reward_type_id = p_reward_type_id;
    
    IF v_account_id IS NULL THEN
        -- 创建新账户
        INSERT INTO accounts (family_id, child_id, reward_type_id, balance) 
        VALUES (p_family_id, p_child_id, p_reward_type_id, 0);
        SET v_account_id = LAST_INSERT_ID();
    END IF;
    
    -- 锁定账户行（防止并发问题）
    SELECT balance INTO v_new_balance 
    FROM accounts 
    WHERE id = v_account_id 
    FOR UPDATE;
    
    -- 创建交易记录
    INSERT INTO transactions (account_id, type, value, note, created_by, idempotency_key)
    VALUES (v_account_id, 'credit', p_value, p_note, p_created_by, p_idempotency_key);
    
    -- 更新余额（触发器会自动处理）
    SET v_new_balance = v_new_balance + p_value;
    
    -- 提交事务
    COMMIT;
    
    -- 返回结果
    SELECT LAST_INSERT_ID() as transaction_id, v_new_balance as new_balance;
    
END//

DELIMITER ;

-- 创建存储过程：安全地消费奖励（带余额检查）
DELIMITER //

CREATE PROCEDURE spend_reward(
    IN p_family_id BIGINT,
    IN p_child_id BIGINT,
    IN p_reward_type_id BIGINT,
    IN p_value BIGINT,
    IN p_note VARCHAR(255),
    IN p_created_by BIGINT,
    IN p_idempotency_key VARCHAR(64)
)
BEGIN
    DECLARE v_account_id BIGINT;
    DECLARE v_current_balance BIGINT;
    DECLARE v_new_balance BIGINT;
    
    -- 开始事务
    START TRANSACTION;
    
    -- 检查幂等性
    IF p_idempotency_key IS NOT NULL AND p_idempotency_key != '' THEN
        IF EXISTS (SELECT 1 FROM transactions WHERE idempotency_key = p_idempotency_key) THEN
            -- 返回已存在的交易信息
            SELECT id, (SELECT balance FROM accounts WHERE id = account_id) as balance
            FROM transactions 
            WHERE idempotency_key = p_idempotency_key;
            COMMIT;
            LEAVE proc_label;
        END IF;
    END IF;
    
    -- 获取账户ID和当前余额
    SELECT id, balance INTO v_account_id, v_current_balance
    FROM accounts 
    WHERE child_id = p_child_id AND reward_type_id = p_reward_type_id;
    
    IF v_account_id IS NULL THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Account not found';
        ROLLBACK;
        LEAVE proc_label;
    END IF;
    
    -- 检查余额是否足够
    IF v_current_balance < p_value THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient balance';
        ROLLBACK;
        LEAVE proc_label;
    END IF;
    
    -- 锁定账户行（防止并发问题）
    SELECT balance INTO v_current_balance 
    FROM accounts 
    WHERE id = v_account_id 
    FOR UPDATE;
    
    -- 创建交易记录
    INSERT INTO transactions (account_id, type, value, note, created_by, idempotency_key)
    VALUES (v_account_id, 'debit', p_value, p_note, p_created_by, p_idempotency_key);
    
    -- 计算新余额（触发器会自动处理）
    SET v_new_balance = v_current_balance - p_value;
    
    -- 提交事务
    COMMIT;
    
    -- 返回结果
    SELECT LAST_INSERT_ID() as transaction_id, v_new_balance as new_balance;
    
END//

DELIMITER ;

-- 创建视图：账户余额汇总
CREATE OR REPLACE VIEW account_balances AS
SELECT 
    a.id as account_id,
    a.family_id,
    a.child_id,
    u.display_name as child_name,
    a.reward_type_id,
    rt.name as reward_type_name,
    rt.unit_kind,
    rt.unit_label,
    a.balance,
    a.created_at,
    a.updated_at
FROM accounts a
JOIN users u ON a.child_id = u.id
JOIN reward_types rt ON a.reward_type_id = rt.id;

-- 创建视图：交易历史（带详细信息）
CREATE OR REPLACE VIEW transaction_history AS
SELECT 
    t.id as transaction_id,
    t.account_id,
    a.family_id,
    a.child_id,
    u.display_name as child_name,
    a.reward_type_id,
    rt.name as reward_type_name,
    rt.unit_kind,
    rt.unit_label,
    t.type,
    t.value,
    t.note,
    t.created_by,
    creator.display_name as creator_name,
    t.idempotency_key,
    t.created_at
FROM transactions t
JOIN accounts a ON t.account_id = a.id
JOIN users u ON a.child_id = u.id
JOIN reward_types rt ON a.reward_type_id = rt.id
JOIN users creator ON t.created_by = creator.id;

-- 添加注释
ALTER TABLE families COMMENT = '家庭信息表';
ALTER TABLE users COMMENT = '用户信息表（监护人和孩子）';
ALTER TABLE reward_types COMMENT = '奖励类型定义表';
ALTER TABLE accounts COMMENT = '孩子账户表（按奖励类型）';
ALTER TABLE transactions COMMENT = '交易记录表';
ALTER TABLE audit_logs COMMENT = '审计日志表';

-- 权限设置（根据实际需求调整）
-- 创建专用数据库用户
-- CREATE USER 'reward_user'@'localhost' IDENTIFIED BY 'secure_password';
-- GRANT SELECT, INSERT, UPDATE, DELETE ON reward_system.* TO 'reward_user'@'localhost';
-- FLUSH PRIVILEGES;

-- 显示创建结果
SELECT '数据库表创建完成' as message;
SHOW TABLES;