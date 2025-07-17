# Commitlint

Инструмент для проверки сообщений коммитов в соответствии с [Conventional Commits](https://www.conventionalcommits.org/).

## 📋 Описание

Commitlint состоит из двух независимых приложений:

1. **CLI приложение** - инструмент командной строки для локальной проверки коммитов
2. **Webhook сервер** - веб-сервер для интеграции с Gitea, автоматически проверяющий коммиты в pull request'ах

## 🚀 Установка

### Из исходного кода

```bash
# Клонирование репозитория
git clone https://github.com/conventionalcommit/commitlint.git
cd commitlint

# Сборка CLI приложения
make build-cli

# Сборка webhook сервера
make build-webhook-server

# Или собрать оба приложения
make build
```

### Установка в систему

```bash
# Установка CLI
make install-cli

# Установка webhook сервера
make install-webhook-server
```

### Docker

```bash
# Сборка Docker образа для CLI
make docker-build-cli

# Сборка Docker образа для webhook сервера
make docker-build-webhook-server
```

## 📖 CLI Приложение

### Использование

```bash
# Проверка последнего коммита
commitlint

# Проверка конкретного коммита
commitlint --commit-hash <hash>

# Проверка диапазона коммитов
commitlint --from-hash <hash1> --to-hash <hash2>

# Использование конкретного конфигурационного файла
commitlint --config path/to/config.yml
```

### Конфигурация

CLI приложение использует файл `.commitlintrc.yml` для настройки правил проверки:

```yaml
# .commitlintrc.yml
min_version: "v0.13.0"

formatter: default

rules:
  header-min-length:
    severity: error
    options:
      min: 10

  header-max-length:
    severity: error
    options:
      max: 72

  type-enum:
    severity: error
    options:
      types:
        - feat     # Новая функциональность
        - fix      # Исправление ошибок
        - docs     # Изменения в документации
        - style    # Форматирование кода
        - refactor # Рефакторинг кода
        - perf     # Улучшение производительности
        - test     # Добавление тестов
        - build    # Изменения в системе сборки
        - ci       # Изменения в CI/CD
        - chore    # Рутинные задачи
        - revert   # Откат изменений

  body-max-line-length:
    severity: warning
    options:
      max: 100

  footer-max-line-length:
    severity: warning
    options:
      max: 100
```

### Git Hooks

Для автоматической проверки коммитов можно настроить git hook:

```bash
#!/bin/sh
# .git/hooks/commit-msg

commitlint --commit-msg-file $1
```

## 🌐 Webhook Сервер

### Описание

Webhook сервер принимает вебхуки от Gitea при создании или обновлении pull request'ов, проверяет все коммиты и отправляет статус проверки обратно в Gitea.

### Конфигурация

Создайте файл `webhook-server.yml`:

```yaml
# HTTP сервер
server:
  address: ":8080"              # Адрес сервера (host:port)
  read_timeout: 30              # Таймаут чтения в секундах
  write_timeout: 30             # Таймаут записи в секундах

# Настройки Gitea API
gitea:
  base_url: "https://gitea.example.com"  # URL вашего Gitea сервера
  token: "your-gitea-api-token"          # API токен с доступом к репозиториям
  username: "commitlint-bot"             # Имя пользователя бота (опционально)

# Настройки commitlint
commitlint:
  config_path: ".commitlintrc.yml"       # Путь к конфигурационному файлу

# Настройки webhook
webhook:
  secret: "your-webhook-secret"          # Секрет для проверки подписи webhook
  path: "/webhook"                       # Путь endpoint'а webhook
  events:                                # События Gitea для обработки
    - "pull_request"
    - "pull_request_sync"
```

### Переменные окружения

Вместо файла конфигурации можно использовать переменные окружения:

- `WEBHOOK_SERVER_ADDRESS` - адрес сервера
- `GITEA_BASE_URL` - URL Gitea сервера
- `GITEA_TOKEN` - API токен
- `GITEA_USERNAME` - имя пользователя
- `WEBHOOK_SECRET` - секрет webhook

### Запуск

#### Обычный запуск

```bash
# С файлом конфигурации
commitlint-webhook-server

# С переменными окружения
GITEA_BASE_URL=https://gitea.example.com \
GITEA_TOKEN=your-token \
WEBHOOK_SECRET=your-secret \
commitlint-webhook-server
```

#### Docker

```bash
# Создайте docker-compose.yml
version: '3.8'

services:
  commitlint-webhook:
    image: commitlint-webhook-server:latest
    ports:
      - "8080:8080"
    environment:
      - GITEA_BASE_URL=https://gitea.example.com
      - GITEA_TOKEN=${GITEA_TOKEN}
      - WEBHOOK_SECRET=${WEBHOOK_SECRET}
    volumes:
      - ./webhook-server.yml:/etc/commitlint/webhook-server.yml:ro
    restart: unless-stopped

# Запуск
docker-compose up -d
```

#### Systemd

```ini
# /etc/systemd/system/commitlint-webhook.service
[Unit]
Description=Commitlint Webhook Server
After=network.target

[Service]
Type=simple
User=webhook
Group=webhook
ExecStart=/usr/local/bin/commitlint-webhook-server
Restart=on-failure
RestartSec=5
Environment="GITEA_BASE_URL=https://gitea.example.com"
Environment="GITEA_TOKEN=your-token"
Environment="WEBHOOK_SECRET=your-secret"

[Install]
WantedBy=multi-user.target
```

### Настройка Gitea

1. Перейдите в настройки репозитория в Gitea
2. Откройте раздел "Webhooks"
3. Нажмите "Добавить webhook" → "Gitea"
4. Заполните поля:
   - **URL**: `http://your-server:8080/webhook`
   - **Секрет**: тот же секрет, что указан в конфигурации
   - **События**: выберите "Pull Request"
5. Сохраните webhook

### Как это работает

1. При создании или обновлении pull request Gitea отправляет webhook
2. Сервер проверяет подпись webhook (если настроен секрет)
3. Сервер устанавливает статус "pending" для последнего коммита
4. Сервер получает список всех коммитов в PR через Gitea API
5. Каждый коммит проверяется с помощью commitlint
6. Результаты проверки отправляются обратно в Gitea:
   - Статус коммита (success/failure/error)
   - Комментарий в PR с деталями ошибок (если есть)

## 🛠️ Разработка

### Структура проекта

```
commitlint/
├── cmd/
│   ├── cli/                    # CLI приложение
│   │   └── main.go
│   └── webhook-server/         # Webhook сервер
│       └── main.go
├── internal/
│   ├── cmd/                    # CLI команды
│   ├── webhook/                # Webhook сервер логика
│   │   ├── config.go          # Конфигурация
│   │   ├── server.go          # HTTP сервер
│   │   ├── gitea.go           # Gitea API клиент
│   │   ├── lint.go            # Интеграция с commitlint
│   │   └── payload.go         # Типы webhook payload
│   └── ...                    # Общая логика
├── config/                     # Парсеры конфигурации
├── lint/                       # Логика проверки коммитов
├── rule/                       # Правила проверки
├── Makefile                   # Команды сборки
├── Dockerfile.cli             # Docker образ для CLI
└── Dockerfile.webhook         # Docker образ для webhook сервера
```

### Тестирование

```bash
# Запуск всех тестов
make test

# Тестирование webhook сервера локально
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Gitea-Event: pull_request" \
  -d @test-payload.json
```

### Требования

- Go 1.23 или выше
- Git
- Make (опционально)

## 📝 Лицензия

MIT License

## 🤝 Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для вашей функциональности (`git checkout -b feature/amazing-feature`)
3. Закоммитьте изменения (`git commit -m 'feat: add amazing feature'`)
4. Запушьте ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## ❓ FAQ

### Как добавить новое правило проверки?

1. Создайте файл в директории `rule/`
2. Реализуйте интерфейс `Rule`
3. Зарегистрируйте правило в `internal/registry/`

### Как изменить формат вывода?

Создайте новый форматер в директории `formatter/` и укажите его в конфигурации.

### Поддерживаются ли другие Git-сервисы кроме Gitea?

На данный момент webhook сервер поддерживает только Gitea. Для поддержки других сервисов (GitHub, GitLab) необходимо реализовать соответствующие клиенты и обработчики webhook.

## 🚀 GitHub Actions

Проект использует GitHub Actions для автоматизации сборки и релизов:

### Workflows

1. **Build and Test** (`.github/workflows/build.yml`)
   - Запускается при push и pull request
   - Тестирует код
   - Собирает бинарные файлы для разных платформ
   - Проверяет конфигурацию GoReleaser

2. **Release** (`.github/workflows/release.yml`)
   - Запускается при создании тега версии (`v*`)
   - Использует GoReleaser для создания релиза
   - Создает бинарные файлы для всех платформ
   - Публикует релиз на GitHub

3. **Docker Build and Push** (`.github/workflows/docker.yml`)
   - Запускается при push в main и создании тегов
   - Собирает Docker образы для CLI и webhook сервера
   - Публикует образы в GitHub Container Registry

4. **Commit Lint** (`.github/workflows/commitlint.yml`)
   - Запускается для всех pull request
   - Проверяет коммиты на соответствие Conventional Commits

### Создание релиза

Для создания нового релиза:

```bash
# Создание тега версии
git tag v1.0.0

# Отправка тега на GitHub
git push origin v1.0.0
```

GitHub Actions автоматически:
- Соберет бинарные файлы для всех платформ
- Создаст архивы с бинарными файлами и документацией
- Опубликует Docker образы
- Создаст черновик релиза на GitHub

## 🔗 Полезные ссылки

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Gitea API Documentation](https://docs.gitea.io/en-us/api-usage/)
- [Оригинальный commitlint для Node.js](https://commitlint.js.org/)
- [GoReleaser Documentation](https://goreleaser.com/)