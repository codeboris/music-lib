init:
	docker compose up --build -d

down:
	docker compose down --remove-orphans

logs:
	docker compose logs music-lib-app
	
db:
	docker compose exec postgres psql -U postgres

swag:
	swag init -g cmd/main.go

migrate-up:
	docker compose run --rm migrate -path=/migrations -database 'postgresql://postgres:12345@postgres:5432/postgres?sslmode=disable' up

migrate-down:
	docker compose run --rm migrate -path=/migrations -database 'postgresql://postgres:12345@postgres:5432/postgres?sslmode=disable' down

migrate-create:
	docker compose run --rm migrate create -ext sql -dir ./migrations -seq $(NAME)