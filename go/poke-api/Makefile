init-postgres-db:
	docker run --name pokeapi-postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=pokeapi_test \
		-v "${PWD}/migrations:/poke-api/migrations" \
		-p 5432:5432 \
		-d postgres:latest

init-postgres-migration:
	docker exec pokeapi-postgres \
	createdb -U postgres pokeapi 2>/dev/null || true
	migrate -database "postgres://postgres:postgres@localhost:5432/pokeapi?sslmode=disable" \
		-path "${PWD}/migrations" up

postgres: init-postgres-db init-postgres-migration

init-elastic-db:
	docker run --name pokeapi-elastic \
		-e "discovery.type=single-node" \
		-e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
		-e "xpack.security.enabled=false" \
		-e ELASTIC_PASSWORD="elastic" \
		-p 9200:9200 \
		-d docker.elastic.co/elasticsearch/elasticsearch:9.1.5
	@until curl -s http://localhost:9200 >/dev/null; do \
		echo "Waiting for Elasticsearch..."; \
		sleep 2; \
	done

init-elastic-index:
	curl -X PUT "http://localhost:9200/pokedex" \
		-H 'Content-Type: application/json' \
		-d @pokeapi-index.json

elasticsearch: init-elastic-db init-elastic-index

stop-elasticsearch:
	docker stop pokeapi-elastic
	docker rm pokeapi-elastic

stop-postgres:
	docker stop pokeapi-postgres
	docker rm pokeapi-postgres

test:
	TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/pokeapi?sslmode=disable \
	TEST_ELASTIC_URL=http://localhost:9200 \
	go test -v ./...
