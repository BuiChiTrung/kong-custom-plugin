up:
	docker compose up
down:
	docker compose down
restart:
	GOOS=linux GOARCH=amd64 go build -o ./bin  ./kong/plugins/proxy-cache-graphql/ && docker compose up
bootstrap:
	docker compose up && docker compose down
run fiber:
	cd fiber-app && go run .
kong:
	docker exec -it -u root kong-custom-plugin-kong-gateway-1 /bin/bash
redis:
	docker exec -it -u root kong-custom-plugin-kong-redis-1 /bin/bash
