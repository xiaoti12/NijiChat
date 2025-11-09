DROP TABLE IF EXISTS groups;
-- 声优群组表
CREATE TABLE IF NOT EXISTS seiyuu_groups (
    id TEXT PRIMARY KEY,                    -- 群组唯一ID
    name TEXT NOT NULL,                     -- 群组名称
    description TEXT,                       -- 群组描述和规则
    member_ids TEXT NOT NULL DEFAULT '[]', -- JSON数组格式的成员声优ID列表
    is_discussion_mode BOOLEAN DEFAULT 0,   -- 是否群讨论模式 (SQLite中BOOLEAN用0/1表示)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_groups_name ON seiyuu_groups(name);

-- 触发器：自动更新 updated_at
CREATE TRIGGER IF NOT EXISTS update_groups_timestamp
AFTER UPDATE ON seiyuu_groups
FOR EACH ROW
BEGIN
    UPDATE seiyuu_groups SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
