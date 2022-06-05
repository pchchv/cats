## Running the environment
Download docker image of PostgreSQL
```
docker pull yzh44yzh/wg_forge_backend_env:1.1
```
Run it
```
docker run -p 5432:5432 -d yzh44yzh/wg_forge_backend_env:1.1
```
Params for ```.env``` file
```
HOST=localhost
PORT=5432
DBNAME=wg_forge_db
USERNAME=wg_forge
PASSWORD=42a
```
