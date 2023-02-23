up:
	docker compose up
down:
	docker compose down
restart:
	GOOS=linux GOARCH=amd64 go build -o ./bin  ./kong/plugins/proxy-cache-graphql/ && docker compose up
run fiber:
	cd fiber-app && go run .