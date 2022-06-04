## Running the environment
Download docker image of PostgreSQL
```
docker pull yzh44yzh/wg_forge_backend_env:1.1
```
Run it
```
docker run -p 5432:5432 -d yzh44yzh/wg_forge_backend_env:1.1
```
Params
```
host: localhost
port: 5432
dbname: wg_forge_db
username: wg_forge
password: 42a
```
