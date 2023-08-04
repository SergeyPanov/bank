create_container:
	docker run --name bank_db -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -p 5432:5432 -d postgres

createdb:
	docker exec -it bank_db createdb --username=root --owner=root bank

dropdb:
	docker exec -it bank_db dropdb bank

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover -count=1 ./... 

server:
	go run main.go

mock:
	go generate ./...
