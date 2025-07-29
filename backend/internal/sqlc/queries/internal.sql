
-- name: InsertSensorData :exec
INSERT INTO sensors_data (value, unit, id1, id2, read_at)
VALUES (?, ?, ?, ?, ?);


-- name: GetAllSensorData :many
SELECT
    sensors_data.value,
    sensors_data.unit,
    sensors_data.id1,
    sensors_data.id2,
    sensors_data.read_at
FROM sensors_data
WHERE 
    (sensors_data.id1 = ? OR ? = 0) AND
    (sensors_data.id2 = ? OR ? = "") AND
    (sensors_data.read_at >= ? OR ? = ?) AND
    (sensors_data.read_at <= ? OR ? = ?)
ORDER BY read_at DESC
LIMIT ?
OFFSET ?;



-- -- name: GetAllSensorData :many
-- SELECT
--     *
-- FROM (
--     SELECT
--         sensors_data.value,
--         sensors_data.unit,
--         sensors_data.id1,
--         sensors_data.id2,
--         sensors_data.read_at
--     FROM sensors_data
--     WHERE 
--         (sensors_data.id1 = ? OR ? = 0) AND
--         (sensors_data.id2 = ? OR ? = "") AND
--         (sensors_data.read_at >= ? OR ? = ?) AND
--         (sensors_data.read_at <= ? OR ? = ?)
--     LIMIT ?
--     OFFSET ?) temp
-- ORDER BY temp.read_at DESC;