-- 管理员表
CREATE TABLE IF NOT EXISTS admins (
    id TEXT PRIMARY KEY,                    -- 管理员ID (UUID)
    username TEXT UNIQUE NOT NULL,          -- 用户名
    password_hash TEXT NOT NULL,            -- 密码哈希 (bcrypt)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_admin_username ON admins(username);

-- 插入默认管理员账号 (用户名: admin, 密码: admin123)
-- 密码哈希使用 bcrypt 生成
INSERT OR IGNORE INTO admins (id, username, password_hash)
VALUES ('00000000-0000-0000-0000-000000000001', 'admin', '$2a$10$XQq5YxH5Z3H5Z3H5Z3H5Zu');
