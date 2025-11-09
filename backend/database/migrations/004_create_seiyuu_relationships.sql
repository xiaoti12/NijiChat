DROP TABLE IF EXISTS seiyuu_relationships;
-- 声优关系表
CREATE TABLE IF NOT EXISTS seiyuu_relationships (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))), -- UUID主键
    seiyuu_id_a TEXT NOT NULL,              -- 声优A的ID
    seiyuu_id_b TEXT NOT NULL,              -- 声优B的ID
    relationship_description TEXT NOT NULL, -- 统一的关系描述（从上帝视角阐述）
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    -- 外键约束
    FOREIGN KEY (seiyuu_id_a) REFERENCES seiyuu(id) ON DELETE CASCADE,
    FOREIGN KEY (seiyuu_id_b) REFERENCES seiyuu(id) ON DELETE CASCADE,

    -- 确保每对声优关系唯一（避免重复关系）
    UNIQUE(seiyuu_id_a, seiyuu_id_b)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_seiyuu_relationships_a ON seiyuu_relationships(seiyuu_id_a);
CREATE INDEX IF NOT EXISTS idx_seiyuu_relationships_b ON seiyuu_relationships(seiyuu_id_b);

-- 触发器：自动更新 updated_at
CREATE TRIGGER IF NOT EXISTS update_seiyuu_relationships_timestamp
AFTER UPDATE ON seiyuu_relationships
FOR EACH ROW
BEGIN
    UPDATE seiyuu_relationships SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;