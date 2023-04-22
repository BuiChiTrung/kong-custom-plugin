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
kong-gateway:
	docker exec -it -u root kong-gateway /bin/bash
redis-1:
	docker exec -it -u root kong-redis-1 /bin/bash
redis-2:
	docker exec -it -u root kong-redis-2 /bin/bash
stop-redis-1:
	docker compose stop kong-redis-1
stop-redis-2:
	docker compose stop kong-redis-2
start-redis-1:
	docker compose start kong-redis-1
start-redis-2:
	docker compose start kong-redis-2
