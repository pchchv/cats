# Running the application
```
docker-compose up --build
```
## Running a database for development
```
bash run.sh
```
## Running tests (db must be running)
```
go test
```
### HTTP Methods
```
/ping — Checking the server connection
```
```
/cats — Getting a list of cats in JSON format
options: 
    attribute — Which key to sort by
    order — asc or desc
    offset — Skip a specified amount of records
    limit — Output the specified number of records
    
example: http://localhost:8080/cats?attribute=color&order=asc&offset=5&limit=2
```
```
/cat — Inserts a new cat into the database. Waits for cat data(name, color, tail length, whiskers length)in JSON format
```
### Params for ```.env``` file
```
HOST=localhost
PORT=5432
DBNAME=wg_forge_db
USERNAME=wg_forge
PASSWORD=42a
```
