-- name: GetTerrainByKey :one
SELECT data FROM terrains WHERE key = ?;

-- name: CreateTerrain :one
INSERT OR REPLACE INTO terrains (key, data, created_at, last_accessed)
VALUES (?, ?, ?, ?);

-- name: SetLastAccessed :exec
UPDATE terrains SET last_accessed = ? WHERE key = ?;

-- name: CleanOldest :exec
DELETE FROM terrains WHERE key IN (
    SELECT key FROM terrains ORDER BY last_accessed ASC LIMIT ?
)