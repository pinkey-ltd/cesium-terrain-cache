CREATE TABLE terrains
(
    key           TEXT PRIMARY KEY, -- 瓦片键（如 "example:EPSG:4326:5:10:15"）
    data          BLOB,             -- 瓦片数据（PNG/JPEG 二进制）
    created_at    INTEGER,          -- 创建时间戳（用于清理）
    last_accessed INTEGER           -- 最后访问时间戳（用于 LRU）
);

CREATE INDEX idx_last_accessed ON terrains (last_accessed);