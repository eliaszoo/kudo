-- 测试数据插入脚本
-- 用于开发和测试环境

-- 已通过连接字符串选择数据库，无需在脚本中使用 USE

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
INSERT IGNORE INTO users (family_id, role, display_name, wechat_openid) VALUES 
-- 张家庭
((SELECT MIN(id) FROM families WHERE name = '张家庭'), 'guardian', '张爸爸', 'zhang_father_openid'),
((SELECT MIN(id) FROM families WHERE name = '张家庭'), 'guardian', '张妈妈', 'zhang_mother_openid'),
((SELECT MIN(id) FROM families WHERE name = '张家庭'), 'child', '张小明', 'zhang_xiaoming_openid'),
((SELECT MIN(id) FROM families WHERE name = '张家庭'), 'child', '张小美', 'zhang_xiaomei_openid'),
-- 李家庭
((SELECT MIN(id) FROM families WHERE name = '李家庭'), 'guardian', '李爸爸', 'li_father_openid'),
((SELECT MIN(id) FROM families WHERE name = '李家庭'), 'child', '李小强', 'li_xiaoqiang_openid');

-- 创建奖励类型
INSERT IGNORE INTO reward_types (family_id, name, unit_kind, unit_label) VALUES 
-- 张家庭奖励类型
((SELECT MIN(id) FROM families WHERE name = '张家庭'), '零花钱', 'money', '元'),
((SELECT MIN(id) FROM families WHERE name = '张家庭'), '看电视时间', 'time', '分钟'),
((SELECT MIN(id) FROM families WHERE name = '张家庭'), '游戏时间', 'time', '分钟'),
((SELECT MIN(id) FROM families WHERE name = '张家庭'), '积分', 'points', '积分'),
((SELECT MIN(id) FROM families WHERE name = '张家庭'), '星星奖励', 'custom', '星星'),
-- 李家庭奖励类型
((SELECT MIN(id) FROM families WHERE name = '李家庭'), '零用钱', 'money', '元'),
((SELECT MIN(id) FROM families WHERE name = '李家庭'), '娱乐时间', 'time', '分钟');

-- 创建账户（为每个孩子和奖励类型组合）
INSERT IGNORE INTO accounts (family_id, child_id, reward_type_id, balance)
SELECT u.family_id, u.id, rt.id, 0
FROM users u
JOIN reward_types rt ON rt.family_id = u.family_id
WHERE u.role = 'child';

-- 测试数据暂不插入交易记录，避免依赖触发器更新余额

-- 暂不插入审计日志

-- 创建测试查询