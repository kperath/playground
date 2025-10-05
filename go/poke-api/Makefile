postgres:
	docker run --name pokeapi-postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=pokeapi_test \
		-p 5432:5432 \
		-d postgres:latest

stop-postgres:
	docker stop pokeapi-postgres
	docker rm pokeapi-postgres

test:
	TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/pokeapi?sslmode=disable \
	go test -v ./...
