DB_URL=postgres://root:password@localhost:5000/simple_bank?sslmode=disable

createdb:
	docker exec -it postgres_bank createdb --username=root simple_bank
dropdb:
	docker exec -it postgres_bank dropdb simple_bank
migrateup:
	migrate -path db/migration -database "$(DB_URL)" up
migrateup1:
	migrate -path db/migration -database "$(DB_URL)" up 1
migratedown:
	migrate -path db/migration -database "$(DB_URL)" down
migratedown1:
	migrate -path db/migration -database "$(DB_URL)" down 1
generate:
	docker run --rm -v "D:/VScode/SimpleBank:/src" -w /src kjconroy/sqlc generate
test:
	go test -v -cover ./...
start:
	go run main.go
server:
	nodemon --exec go run main.go