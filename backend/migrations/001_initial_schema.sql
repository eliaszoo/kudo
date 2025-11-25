-- 奖励系统数据库迁移脚本
-- 基于 MySQL 8.0+

-- 创建数据库（如果不存在）
-- 已通过连接字符串选择数据库，无需在脚本中使用 USE

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

-- 仅创建表结构，不插入默认数据



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