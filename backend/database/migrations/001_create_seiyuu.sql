DROP TABLE IF EXISTS seiyuu;
-- 声优基础信息表
CREATE TABLE IF NOT EXISTS seiyuu (
    id TEXT PRIMARY KEY,                    -- 声优唯一ID (UUID)
    name TEXT NOT NULL,                     -- 声优姓名
    avatar_url TEXT,                        -- 头像图片URL
    profile_markdown TEXT NOT NULL,         -- 完整Markdown格式资料
    raw_profile_data TEXT,                  -- 原始资料数据（用于AI关系生成）
    tags TEXT DEFAULT '[]',                 -- JSON数组格式的标签 ["萝莉音", "治愈系"]
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'inactive')),  -- 状态
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_seiyuu_status ON seiyuu(status);
CREATE INDEX IF NOT EXISTS idx_seiyuu_name ON seiyuu(name);

-- 触发器：自动更新 updated_at
CREATE TRIGGER IF NOT EXISTS update_seiyuu_timestamp
AFTER UPDATE ON seiyuu
FOR EACH ROW
BEGIN
    UPDATE seiyuu SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
