DB_URL=postgres://root:root@localhost:5432/ecommerce?sslmode=disable

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root ecommerce

dropdb:
	docker exec -it postgres12 dropdb ecommerce

migrateup:
	migrate -path pkg/databases/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path pkg/databases/migrations -database "$(DB_URL)" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown