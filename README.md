# Сервис назначения ревьюеров для Pull Request’ов
Простой, масштабируемый и типичный «чистый» Go-сервис с разделением на слои, валидацией, миграциями.
### Стек технологий
- **Go** 1.23+
- **Chi** — легковесный роутер
- **PostgreSQL** 15+
- **Docker + Docker Compose**
- **cleanenv** для загрузки конфигурации
- Миграции **goose**
- ORM **gorm*
- Валидация через **github.com/go-playground/validator**
- Логирование — **slog**
- **gofakeit** — сиды с реалистичными данными
####Готовая коллекция с примерами всех запросов:
/Pull Requests.postman_collection.json

## Быстрый старт (Docker Compose)
### 1. Запуск
```bash
make up
```
Сервис будет доступен по адресу:  
**http://localhost:8080/**
По желанию порт можно помеять в файле конфигурации`
**config/dev.yaml** ->http_server: address: "0.0.0.0:8080"

### 2. Остановка
```bash
make down                    # остановить контейнеры
make drop                    # + удалить volume с БД (полная очистка)
```

## Конфигурация
Все параметры загружаются из ``config/.yaml`` файла.

| Переменная окружения**         | Поле в YAML-конфиге                   | Описание                          | Значение в dev.yaml   | По умолчанию (если не задано) |
|--------------------------------|---------------------------------------|-----------------------------------|-----------------------|--------------------------------|
| `CONFIG_PATH`                  | —                                     | Путь к YAML-файлу конфигурации    | `/app/config/dev.yaml`| — в Docker.compose             |
| —                              | `env`                                 | Окружение (dev/prod/local)        | `dev`                 | —                              |
| —                              | `http_server.address`                 | Адрес и порт HTTP-сервера         | `0.0.0.0:8080`        | —                              |
| —                              | `http_server.timeout`                 | Общий таймаут сервера             | `4s`                  | —                              |
| —                              | `http_server.idle_timeout`            | Idle timeout                      | `30s`                 | —                              |
| —                              | `database.host`                       | Хост PostgreSQL                   | `pr_postgres`         | —                              |
| —                              | `database.port`                       | Порт PostgreSQL                   | `5432`                | —                              |
| —                              | `database.user`                       | Пользователь БД                   | `postgres`            | —                              |
| —                              | `database.password`                   | Пароль БД                         | `postgres`            | —                              |
| —                              | `database.dbname`                     | Имя базы данных                   | `pullrequest`         | —                              |
| —                              | `database.sslmode`                    | Режим SSL                         | `disable`             | —                              |

Миграции автоматически применяются при старте приложения.

## API Эндпоинты
| Метод  | Путь                            | Описание                                                                  |
|-------|----------------------------------|-------------------------------------------------------------------------- |
| `POST`  | `/pullRequest/create`            | Создать PR + автоматически назначить до 2 ревьюверов из команды автора  |
| `POST`  | `/pullRequest/merge`             | Пометить PR как MERGED (идемпотентно)                                   |
| `POST`  | `/pullRequest/reassign`          | Переназначить ревьювера на другого из его команды                       |
| `POST`  | `/team/add`                      | Создать команду (создаёт/обновляет пользователей)                       |
| `GET`   | `/team/get?team_name=Alpha`      | Получить команду с участниками                                          |
| `GET`   | `/users/getReview?user_id=xxx`   | Получить все PR, где пользователь назначен ревьювером                   |
| `POST`  | `/users/setIsActive`             | Установить флаг активности пользователя                                 |
## Gofakeit
После запуска приложение сидит базу данных одинаковым зерном. 
Таблица пользователей 
| user_id**                                | username**           | is_active | team_name |
|------------------------------------------|----------------------|-----------|-----------|
| 72b553cc-c00b-44e6-bb48-4710e784acb8     | Forest.Kihn          | t         | Alpha     |
| d97b3a8f-c1fd-440f-ac40-7375d5fd00c4     | weasel_54            | t         | Beta      |
| 1d98b332-69f8-41a6-8793-2dadbdc03295     | SierraLemonChiffon   | t         | Gamma     |
| 30b7531a-23fc-405f-895e-ea68af2a439d     | Bennieshiny          | t         | Alpha     |
| cd52f961-83a3-4b4f-be9a-b3a1854d7f2b     | LycheeUptight74      | t         | Beta      |
| 59306896-e84b-4b0f-9a9c-8e0c2a440002     | Unique473            | t         | Gamma     |
| 6b6bb0ea-89a3-4927-8343-33348a610f57     | Juliekuban           | t         | Alpha     |
| 29e13b17-6bbd-4cd6-9098-73b3dd5a2692     | Khaki_group          | t         | Beta      |
| 8175ff2c-eb81-4e1c-bdd6-06c652399928     | Kuphal936            | t         | Gamma     |
| 29bad72f-9665-482f-bf96-605dc7b88e02     | Breitenberg_Zkw      | t         | Alpha     |
| f3f2b353-42b0-4ecd-bedd-9e5204013257     | Ayana.43             | f         |           |
| 6d4d1f10-f822-4b89-af29-a9fab40c1db6     | drRomaguera          | f         |           |
| 61c8133d-b016-4e25-ac43-18a3fedb4b43     | SleepyUnderwear82    | f         |           |
| 7bb6ea93-563a-459d-9c22-23b835169a7c     | CortneyPurple        | f         |           |
| 729faab7-264c-4656-83df-a37319ed0541     | TiePainter           | f         |           |
| 7567dfc7-e49a-4b24-b666-bb943d9f830b     | SeaGreen_life        | t         |           |
| a4df52b4-4c2a-43e9-8d0a-da965155ffbd     | Janus2               | t         |           |
| 1e70ee5d-6e2f-4c37-84ec-b7b991111cfd     | Chloe124             | t         |           |

