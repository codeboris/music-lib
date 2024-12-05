init:
	docker compose up --build -d

down:
	docker compose down --remove-orphans

logs:
	docker logs music-lib-app
	
db:
	docker compose exec postgres psql -U postgres

swag:
	swag init -g cmd/main.go

migrate-up:
	migrate -path=migrations -database 'postgres://postgres:12345@localhost:5432/postgres?sslmode=disable' up

migrate-down:
	migrate -path=migrations -database 'postgres://postgres:12345@localhost:5432/postgres?sslmode=disable' down

# Вы можете создавать миграции командой make migrate-add NAME=create_table...
migrate-add:
	migrate create -ext sql -dir ./migrations -seq $(NAME)