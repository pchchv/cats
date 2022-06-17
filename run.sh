#!/bin/bash
docker pull yzh44yzh/wg_forge_backend_env:1.1
docker run -p 5432:5432 -d yzh44yzh/wg_forge_backend_env:1.1
echo PostgreSQL started!
go run .