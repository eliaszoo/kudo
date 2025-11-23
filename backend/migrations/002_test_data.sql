-- 测试数据插入脚本
-- 用于开发和测试环境

USE reward_system;

-- 清空现有数据（谨慎使用）
-- SET FOREIGN_KEY_CHECKS = 0;
-- TRUNCATE TABLE audit_logs;
-- TRUNCATE TABLE transactions;
-- TRUNCATE TABLE accounts;
-- TRUNCATE TABLE reward_types;
-- TRUNCATE TABLE users;
-- TRUNCATE TABLE families;
-- SET FOREIGN_KEY_CHECKS = 1;

-- 创建测试家庭
INSERT INTO families (name) VALUES 
('张家庭'),
('李家庭');

-- 创建测试用户
INSERT INTO users (family_id, role, display_name, wechat_openid) VALUES 
-- 张家庭
(1, 'guardian', '张爸爸', 'zhang_father_openid'),
(1, 'guardian', '张妈妈', 'zhang_mother_openid'),
(1, 'child', '张小明', 'zhang_xiaoming_openid'),
(1, 'child', '张小美', 'zhang_xiaomei_openid'),
-- 李家庭
(2, 'guardian', '李爸爸', 'li_father_openid'),
(2, 'child', '李小强', 'li_xiaoqiang_openid');

-- 创建奖励类型
INSERT INTO reward_types (family_id, name, unit_kind, unit_label) VALUES 
-- 张家庭奖励类型
(1, '零花钱', 'money', '元'),
(1, '看电视时间', 'time', '分钟'),
(1, '游戏时间', 'time', '分钟'),
(1, '积分', 'points', '积分'),
(1, '星星奖励', 'custom', '星星'),
-- 李家庭奖励类型
(2, '零用钱', 'money', '元'),
(2, '娱乐时间', 'time', '分钟');

-- 创建账户（为每个孩子和奖励类型组合）
INSERT INTO accounts (family_id, child_id, reward_type_id, balance) VALUES 
-- 张小明（张家庭，孩子ID: 3）
(1, 3, 1, 5000),    -- 零花钱: 50元
(1, 3, 2, 120),     -- 看电视时间: 120分钟
(1, 3, 3, 60),      -- 游戏时间: 60分钟
(1, 3, 4, 1000),    -- 积分: 1000分
(1, 3, 5, 25),      -- 星星奖励: 25个星星
-- 张小美（张家庭，孩子ID: 4）
(1, 4, 1, 3000),    -- 零花钱: 30元
(1, 4, 2, 90),      -- 看电视时间: 90分钟
(1, 4, 3, 30),      -- 游戏时间: 30分钟
(1, 4, 4, 800),     -- 积分: 800分
(1, 4, 5, 15),      -- 星星奖励: 15个星星
-- 李小强（李家庭，孩子ID: 6）
(2, 6, 6, 2000),    -- 零用钱: 20元
(2, 6, 7, 45);      -- 娱乐时间: 45分钟

-- 创建交易记录（最近30天的数据）
INSERT INTO transactions (account_id, type, value, note, created_by, created_at) VALUES 
-- 张小明的交易记录
(1, 'credit', 10000, '期末考试第一名', 1, DATE_SUB(NOW(), INTERVAL 25 DAY)),
(1, 'debit', 5000, '买新书包', 1, DATE_SUB(NOW(), INTERVAL 20 DAY)),
(1, 'credit', 5000, '帮忙做家务', 2, DATE_SUB(NOW(), INTERVAL 15 DAY)),
(1, 'debit', 3000, '买文具', 1, DATE_SUB(NOW(), INTERVAL 10 DAY)),
(1, 'credit', 8000, '数学竞赛获奖', 1, DATE_SUB(NOW(), INTERVAL 5 DAY)),

-- 张小美的交易记录
(6, 'credit', 8000, '舞蹈比赛第一名', 2, DATE_SUB(NOW(), INTERVAL 22 DAY)),
(6, 'debit', 3000, '买新裙子', 2, DATE_SUB(NOW(), INTERVAL 18 DAY)),
(6, 'credit', 3000, '帮忙洗碗', 1, DATE_SUB(NOW(), INTERVAL 12 DAY)),
(6, 'debit', 2000, '买彩色笔', 2, DATE_SUB(NOW(), INTERVAL 8 DAY)),

-- 李小强的交易记录
(11, 'credit', 5000, '足球比赛MVP', 5, DATE_SUB(NOW(), INTERVAL 16 DAY)),
(11, 'debit', 2000, '买足球鞋', 5, DATE_SUB(NOW(), INTERVAL 12 DAY)),
(11, 'credit', 3000, '科学实验获奖', 5, DATE_SUB(NOW(), INTERVAL 6 DAY));

-- 创建审计日志记录
INSERT INTO audit_logs (family_id, user_id, action, payload, created_at) VALUES 
(1, 1, 'create_reward_type', '{"name": "星星奖励", "unit_kind": "custom", "unit_label": "星星"}', DATE_SUB(NOW(), INTERVAL 30 DAY)),
(1, 1, 'grant_reward', '{"child_id": 3, "reward_type_id": 1, "value": 10000, "note": "期末考试第一名"}', DATE_SUB(NOW(), INTERVAL 25 DAY)),
(1, 1, 'spend_reward', '{"child_id": 3, "reward_type_id": 1, "value": 5000, "note": "买新书包"}', DATE_SUB(NOW(), INTERVAL 20 DAY)),
(2, 5, 'create_family', '{"name": "李家庭"}', DATE_SUB(NOW(), INTERVAL 35 DAY));

-- 创建测试查询
-- 查看所有家庭的概览
SELECT 
    f.name as family_name,
    COUNT(DISTINCT u.id) as total_members,
    COUNT(DISTINCT CASE WHEN u.role = 'child' THEN u.id END) as children_count,
    COUNT(DISTINCT rt.id) as reward_types_count
FROM families f
LEFT JOIN users u ON f.id = u.family_id
LEFT JOIN reward_types rt ON f.id = rt.family_id
GROUP BY f.id;

-- 查看孩子的余额汇总
SELECT 
    f.name as family_name,
    u.display_name as child_name,
    rt.name as reward_type,
    rt.unit_kind,
    rt.unit_label,
    a.balance,
    CASE 
        WHEN rt.unit_kind = 'money' THEN CONCAT('¥', FORMAT(a.balance / 100, 2))
        WHEN rt.unit_kind = 'time' THEN CONCAT(a.balance, ' ', rt.unit_label)
        WHEN rt.unit_kind = 'points' THEN CONCAT(a.balance, ' ', rt.unit_label)
        WHEN rt.unit_kind = 'custom' THEN CONCAT(a.balance, ' ', rt.unit_label)
        ELSE CAST(a.balance AS CHAR)
    END as formatted_balance
FROM accounts a
JOIN users u ON a.child_id = u.id
JOIN reward_types rt ON a.reward_type_id = rt.id
JOIN families f ON a.family_id = f.id
WHERE u.role = 'child'
ORDER BY f.id, u.id, rt.id;

-- 查看最近的交易记录
SELECT 
    th.*,
    CASE 
        WHEN th.unit_kind = 'money' THEN CONCAT('¥', FORMAT(th.value / 100, 2))
        WHEN th.unit_kind = 'time' THEN CONCAT(th.value, ' ', th.unit_label)
        WHEN th.unit_kind = 'points' THEN CONCAT(th.value, ' ', th.unit_label)
        WHEN th.unit_kind = 'custom' THEN CONCAT(th.value, ' ', th.unit_label)
        ELSE CAST(th.value AS CHAR)
    END as formatted_value
FROM transaction_history th
WHERE th.created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
ORDER BY th.created_at DESC
LIMIT 20;

-- 显示测试数据创建完成
SELECT '测试数据创建完成' as message;
SELECT '共创建了:' as info;
SELECT CONCAT(COUNT(*), ' 个家庭') FROM families;
SELECT CONCAT(COUNT(*), ' 个用户') FROM users;
SELECT CONCAT(COUNT(*), ' 种奖励类型') FROM reward_types;
SELECT CONCAT(COUNT(*), ' 个账户') FROM accounts;
SELECT CONCAT(COUNT(*), ' 条交易记录') FROM transactions;
SELECT CONCAT(COUNT(*), ' 条审计日志') FROM audit_logs;