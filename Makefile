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
	docker exec -it -u root kong-custom-plugin-kong-gateway-1 /bin/bash
redis-master:
	docker exec -it -u root kong-custom-plugin-kong-redis-1 /bin/bash
redis-replicas:
	docker exec -it -u root kong-custom-plugin-kong-redis-replicas-1 /bin/bash
stop-redis-master:
	docker compose stop kong-redis
stop-redis-replicas:
	docker compose stop kong-redis-replicas
start-redis-master:
	docker compose start kong-redis
start-redis-replicas:
	docker compose start kong-redis-replicas
