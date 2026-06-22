loglint

loglint — линтер для Go, проверяющий лог-записи в коде.

Линтер реализован на базе golang.org/x/tools/go/analysis и совместим с golangci-lint.

Проверяемые правила

1. Лог-сообщение должно начинаться со строчной буквы.
2. Лог-сообщение должно быть написано на английском языке.
3. Лог-сообщение не должно содержать спецсимволы или эмодзи.
4. Лог-сообщение не должно содержать чувствительные данные.

Поддерживаемые логгеры

* log/slog
* go.uber.org/zap-style logger calls

Примеры:

slog.Info("message")
slog.Error("message")
logger.Info("message")
logger.Error("message")

Примеры

Неправильно:

slog.Info("Starting server")
slog.Info("запуск сервера")
slog.Info("server started!!!")
logger.Error("api_key=" + apiKey)

Правильно:

slog.Info("starting server")
slog.Info("server started")
logger.Error("api request failed")

Конфигурация

Файл .loglint.yml:

check_lowercase: true
check_english: true
check_special_chars: true
check_sensitive_data: true
sensitive_patterns:
  - password
  - passwd
  - token
  - secret
  - api_key
  - apikey
  - access_key
  - private_key
allowed_logger_names:
  - slog
  - logger

Запуск

Запуск линтера:

go run ./cmd/loglint ./...

Запуск с конфигом:

go run ./cmd/loglint -config .loglint.yml ./...

Сборка бинарника:

go build ./cmd/loglint
./loglint ./...

Тесты

go test ./...

Автоисправления

Реализованы SuggestedFixes для:

* исправления первой заглавной буквы на строчную;
* удаления спецсимволов из строкового литерала.

Интеграция с golangci-lint

Используется Module Plugin System.

Сборка кастомного golangci-lint:

golangci-lint custom

Запуск:

./custom-gcl run ./...

CI/CD

Настроен GitHub Actions workflow:

* запуск тестов
* сборка проекта

Ограничения

* анализируются строковые литералы и простые конкатенации;
* тип логгера определяется эвристически;
* сложные выражения анализируются частично.