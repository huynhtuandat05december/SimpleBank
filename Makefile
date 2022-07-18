createdb:
	docker exec -it postgres_bank createdb --username=root simple_bank
dropdb:
	docker exec -it postgres_bank dropdb simple_bank
migrateup:
	migrate -path db/migration -database postgres://root:password@localhost:5000/simple_bank?sslmode=disable up
migratedown:
	migrate -path db/migration -database postgres://root:password@localhost:5000/simple_bank?sslmode=disable down
generate:
	docker run --rm -v "D:/VScode/golang-basic/SimpleBank:/src" -w /src kjconroy/sqlc generate
test:
	go test -v -cover ./...