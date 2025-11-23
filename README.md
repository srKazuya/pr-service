Простой, масштабируемый и типичный «чистый» Go-сервис с разделением на слои, валидацией, миграциями.
База данных изначально засижена тестовыми данными с помощью gofakeit \internal\infrastructure\storage\postgres\seeds.go
Postman коллкция \Pull Requests.postman_collection.json
Стек технологий
Go 1.23+
Chi — легковесный роутер
PostgreSQL 15+
Docker + Docker Compose
cleanenv для загрузки конфигурации
Миграции goose
ORM *gorm
Валидация через github.com/go-playground/validator
Логирование — slog
Быстрый старт (Docker Compose)
1. Запуск
make up
Сервис будет доступен по адресу:
http://localhost:8080/ По желанию порт можно помеять в файле конфигурации` config/dev.yaml ->http_server: address: "0.0.0.0:8080"

2. Остановка
make down                    # остановить контейнеры
make drop                    # + удалить volume с БД (полная очистка)
Конфигурация
Все параметры загружаются из config/.yaml файла.

Переменная окружения**	Поле в YAML-конфиге	Описание	Значение в dev.yaml	По умолчанию (если не задано)
CONFIG_PATH	—	Путь к YAML-файлу конфигурации	/app/config/dev.yaml	— в Docker.compose
—	env	Окружение (dev/prod/local)	dev	—
—	http_server.address	Адрес и порт HTTP-сервера	0.0.0.0:8080	—
—	http_server.timeout	Общий таймаут сервера	4s	—
—	http_server.idle_timeout	Idle timeout	30s	—
—	database.host	Хост PostgreSQL	qa_postgres	—
—	database.port	Порт PostgreSQL	5432	—
—	database.user	Пользователь БД	postgres	—
—	database.password	Пароль БД	postgres	—
—	database.dbname	Имя базы данных	questions	—
—	database.sslmode	Режим SSL	disable	—
Миграции автоматически применяются при старте приложения.

API Эндпоинты
Вопросы (Questions)
Метод	Путь	Описание
GET	/questions	Все вопросы
POST	/questions	Создать вопрос
GET	/questions/{questionID}	Получить вопрос с ответами
DELETE	/questions/{questionID}	Удалить вопрос с ответами
Ответы (Answers)
Метод	Путь	Описание
POST	/questions/{questionID}/answers	Добавить ответ к вопросу
GET	/answers/{answerID}	Получить конкретный ответ
DELETE	/answers/{answerID}	Удалить ответ
