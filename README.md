# Online music libraries
## Stack
- Dokcer compose
- Make
- Git

## Start project
```bash
make init
```
Open in browse: http://localhost:8000/swagger/index.html
## Stop project
```bash
make down
```
## Commands for migrations
```bash
make migrate-add NAME=<name_migrate>
```
## Logs music-lib-app
```bash
make logs
```

## Generate swagger docs
```bash
make swag
```

## Connect to DB
```bash
make db
```