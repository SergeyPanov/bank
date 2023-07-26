create_container:
	docker run --name bank_db -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -p 55432:5432 -d postgres

createdb:
	docker exec -it bank_db createdb --username=root --owner=root bank

dropdb:
	docker exec -it bank_db dropdb bank

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:55432/bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:55432/bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

.PHONY: createdb