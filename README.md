# Postgres
```shell
# Install Postgres in docker
# Access psql command on docker
docker exec -it postgres12 psql -U root
# list database
\l
# Create Database
CREATE DATABASE ecommerce;
# Drop Database
DROP DATABASE ecommerce;
```

# Migration
```shell
# Create migrate file version [current directory must in pkg/databases/migrations]
migrate create -ext sql -seq initial
migrate create -ext sql -seq seed
# Migrate up
migrate -path db/migration -database "postgres://root:root@localhost:5432/ecommerce?sslmode=disable" -verbose up
# Migrate Down
migrate -path db/migration -database "postgres://root:root@localhost:5432/ecommerce?sslmode=disable" -verbose down

```
