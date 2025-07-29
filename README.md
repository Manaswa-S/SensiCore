## SensiCore - Sensor Data Streaming & Analytics
A complete system to ingest, store, and query real-time sensor data using Go, Echo, MySQL, and Docker.

This project simulates a sensor network where devices continuously emit data which is captured and processed by a backend server. 
The backend is capable of ingesting real-time data, storing it efficiently, and exposing powerful RESTful APIs to query and analyze this data.

### Features
- Continuous random yet realistic sensor data generation
- Configurable ingestion via batched inserts or streaming individual records
- Extremely flexible filtering on `ID1`, `ID2`, `Start_Timestamp` or `End_Timestamp`, accepts any combination of filters
- Pagination support for large datasets via `limit` and `offset` query parameters
- All services containerized via Docker Compose

> #### To know more about architecture and design choices, please see <a href="https://github.com/Manaswa-S/SensiCore/blob/main/ARCH.md">ARCH.md</a>

---

### Getting Started

#### Option 1: Via Docker (Recommended)

##### Pre-requisites
- Docker & Docker Compose installed (refer <a href="https://docs.docker.com/engine/install/">this</a>)
- Bash

##### Clone and Run
```bash
git clone https://github.com/Manaswa-S/SensiCore.git
cd SensiCore/
docker compose up --build
```

The first time you run this, Docker may take a few minutes to download required images (MySQL, Alpine, etc.).

To see proper running logs of the program, open a new terminal window and do
```bash
docker attach sensicore-backend
```
OR 
```bash
docker attach sensicore-sensor
```

#### Option 2: Run Direct Locally

##### Pre-requisites
- Go 1.23
- MySQL 8.0.43

Note: Ensure relevant environment variables are set correctly for local runs.

##### Clone and Run
```bash
git clone https://github.com/Manaswa-S/SensiCore.git
cd SensiCore/backend/
go mod download
go run cmd/main.go
```
Now In a different terminal window, navigate to the `SensiCore` folder again, and do
```bash
cd sensor
go mod download
go run main.go
```
You can now see the running logs in there directly.

---

### API Usage

At this point, your backend server as well as the data generator should be running.

Visit 
```bash
http://localhost:8686/public/data
```
to retrieve the data.

You can use any combination of the query parameters and fetch desired data.

The data is sorted by timestamp in descending order, returning the `{limit}` most recent records after skipping `{offset}` newer entries. The most recent data appears first. This follows unless `{start}` or `{end}` is given.


#### Query Params:

| Param       |   Type   | Optional? | Description                | Example    |    Default |
|-------------|----------|-----------|----------------------------|------------|------------|
| `id1`       | integer  | ✅        | Sensor ID                  | 423        | x          |
| `id2`       | string   | ✅        | Subsensor ID               | A          | x          |
| `start`     | integer  | ✅        | Start UNIX time in seconds | 1753707407 | x          |
| `end`       | integer  | ✅        | End UNIX time in seconds   | 1753707507 | x          |
| `limit`     | integer  | ✅        | Max results                | 100        | 25         |
| `offset`    | integer  | ✅        | Skip results               | 200        | 0          |


### Example Calls

| Query                                                      | Result                |
|------------------------------------------------------------|-----------------------|
|`localhost:8686/public/data?id1=1&id2=A`                    | <a href="https://github.com/Manaswa-S/SensiCore/blob/main/zresults/1.A.0.0.json">result</a> |
|`localhost:8686/public/data?start=1753800631&end=1753800731`| <a href="https://github.com/Manaswa-S/SensiCore/blob/main/zresults/0.0.1753800631.1753800731.json">result</a> |

Paste the query in any browser or replace {query} in following with it.
```bash
curl -X GET "{query}"
```

##### Remember
- All timestamps are in UTC.

---

### Configurations
These configurations can be set via enviornment variables or command line arguments.
#### Sensor Generator Configs
| Variable        | Description                         | Example     | Default |
|-----------------|-------------------------------------|-------------|---------|
| `SENSORS_COUNT` | Number of sensors to spawn          | `25`        | `4`     |
| `STREAM_DATA`   | Stream mode toggle (`true`/`false`) | `true`      | `true`  |
| `DATA_CHAN_SIZE`| Size of internal channel buffer     | `1000`      | `100`   |
| `BUFFER_LIMIT`  | Batch limit before flush            | `25`        | `70`    |


### Care to be Taken
- Make sure the ports `3306` and `8686` are free and not binded to any other process. Adjust in ``docker-compose.yml`` if needed. 
- Ensure `docker` and `docker compose` are working and able to pull images.
- The backend includes a waiting mechanism, but may still sometimes start before MySQL has fully initialized.
  In that case, please restart the backend by doing `docker restart sensicore-backend`.
- If containers with the same name exist, you may need to remove them (refer <a href="https://docs.docker.com/reference/cli/docker/container/rm/">this</a>).
- If you're on Linux and not in the Docker group, you may need to prefix Docker commands with `sudo`.


### Architecture
- **Sensor Generator** (Go): Emits randomized, realistic sensor readings either in buffered batches or continuous stream. 
- **Backend API Server** (Go + Echo): Accepts incoming data, parses, cleans and stores in MySQL, exposes REST APIs for retrieval.
- **MySQL**: Stores sensor data.
- **Docker**: Manages multi-container setup

### Tech Stack
- **Language**: Go 1.23
- **Web Framework**: Echo v4.13.4
- **Database**: MySQL 8.0.43
- **ORM/Query Tool**: sqlc v1.28.0
- **Containerization**: Docker 28.3.2
- **Scripting**: Bash

---

### License
This project is open-sourced only for assignment review. All rights reserved © 2025.