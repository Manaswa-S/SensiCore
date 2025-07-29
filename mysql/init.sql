CREATE DATABASE IF NOT EXISTS sensicore;

USE sensicore;

CREATE TABLE IF NOT EXISTS sensicore.sensors_data (
	data_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	value FLOAT NOT NULL,
	unit VARCHAR(24),
	id1 INT NOT NULL,
	id2 VARCHAR(8) NOT NULL,
	read_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_sensors_data_id1 ON sensicore.sensors_data (id1);

CREATE INDEX IF NOT EXISTS idx_sensors_data_id2 ON sensicore.sensors_data (id2);

CREATE INDEX IF NOT EXISTS idx_sensors_data_read_at ON sensicore.sensors_data (read_at);
