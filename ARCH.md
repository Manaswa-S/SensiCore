## Deep Dive into SensiCore Architecture

### Backend

##### Layers
`Router -> Handler -> Service -> MySQL (via sqlc)`

- **Handler**: Performs request data extraction and response formatting and then passes context and data to the service.
- **Service**: Implements core logic, data validation, parsing, filtering, database calls, etc.
- **MySQL (sqlc)**: Type safe database interaction layer, auto generated from raw .sql query files.


##### DI Model
The backend is constructed around the Dependency Injection (DI) Model, 
where dependencies like DB connections, service layer, etc are injected at runtime, maintaining testability and separation of concerns.


##### Ingest Data
- `POST/public/data`

  Accepts an array of sensor events.
  It processes the entire batch at once, in a SQL transaction.
  Partial failure tolerant i.e. if one fails, others proceed. Failures are only logged (silent drops).

  
- `POST/public/data/stream`

  Accepts a stream of sensor events via an open request body.
  An infinite for loop reads and processes one record at a time.
  Has a max read timeout of 30 seconds.


##### Retrieve Data
- `GET/public/data`

  All query parameters are optional and are used in an (OR)AND based filter implemented directly in SQL.
  Returns the latest `limit` values after skipping `offset` records, ordering the result by `read_at`.
  
```sql    
    FROM sensors_data
    WHERE 
        (sensors_data.id1 = ? OR ? = 0) AND
        (sensors_data.id2 = ? OR ? = "") AND
        (sensors_data.read_at >= ? OR ? = ?) AND
        (sensors_data.read_at <= ? OR ? = ?)
    ORDER BY read_at DESC
    LIMIT ?
    OFFSET ?;
```

---

### MySQL Database

##### Tables
-  `sensors_data`

  Stores the sensor events.
  
  ```sql
    data_id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    value FLOAT NOT NULL,
    unit VARCHAR(24),
    id1 INT NOT NULL,
    id2 VARCHAR(8) NOT NULL,
    read_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
  ```

-  `read_at`: Sensor's event timestamp (external, from sensor)
-  `created_at`: Server-generated insertion timestamp



##### Indexes
| Index                     | Column    |
|---------------------------|-----------|
|`idx_sensors_data_id1`     | `id1`     |
|`idx_sensors_data_id2`     | `id2`     |
|`idx_sensors_data_read_at` | `read_at` |
