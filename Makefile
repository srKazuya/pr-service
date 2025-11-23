# =============================================================================
# Makefile — только docker-compose + openapi генерация
# =============================================================================

## Запустить проект
up:
	docker-compose up --build

## Остановить проект
down:
	docker-compose down
	
## + удалить volume с БД (полная очистка)
drop:
	docker-compose down -v

## Перезапустить
restart: down up

## Логи в реальном времени
logs:
	docker-compose logs -f

## Пересобрать образы (если поменял код или Dockerfile)
build:
	docker-compose build --no-cache

## Сгенерировать openapi-код (один раз)

.PHONY: gen-openapi

OPENAPI_FILE=internal/infrastructure/http/openapi/openapi.yml
OPENAPI_OUT=internal/infrastructure/http/openapi

gen-openapi:
	oapi-codegen -generate types -o $(OPENAPI_OUT)/types.gen.go $(OPENAPI_FILE) 
	oapi-codegen -generate chi-server -o $(OPENAPI_OUT)/server.gen.go $(OPENAPI_FILE)
