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



### Sensor Data Generator

##### Generator

-   A new go routine is spawned corresponding to each sensor.
-   Each of these sensor generates random data values (with-in a certain realistic range from initial read).
-   The `id1` and `id2` are somewhat deterministic as they oscillate between a pre-defined range.
-   The read_at time is the time it was generated.
-   The generation is periodic, the interval being randomly decided between a range, and used forever without changing.



##### Flusher

-   Two data channels are initialized on startup. One being the main channel, the other one being the alternate/substitute.

##### Stream Events Mode

-   A connection is opened and an `io.Pipe` is used to continuosly push incoming events from main channel to the connection.
-   This connection is held open for a certain time, and then closed, and swapped for a new connection.

##### Batch Events Mode
-   In this mode, a channel-based dual-buffer design is used.
    -  A producer fills buffer A while buffer B is flushed to backend.
    -  Once buffer A is full, buffers are swapped, allowing ingestion to continue without blocking flush.
    -  This improves throughput without needing locks.

-   A time.Ticker periodically, over a short interval, checks for main data channel crossing its `BUFFER_LIMIT`.
-   Once the limit is crossed, the channels are swapped, the main channel, after swapping becomes alternate channel, which is then flushed in  a batch to the backend. 