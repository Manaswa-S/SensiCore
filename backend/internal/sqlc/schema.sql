CREATE TABLE sensors_data (
	data_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	value FLOAT NOT NULL,
	unit VARCHAR(24),
	id1 INT NOT NULL,
	id2 VARCHAR(8) NOT NULL,
	read_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);