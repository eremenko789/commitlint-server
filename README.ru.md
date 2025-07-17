# commitlint

[![PkgGoDev](https://pkg.go.dev/badge/github.com/conventionalcommit/commitlint)](https://pkg.go.dev/github.com/conventionalcommit/commitlint)

commitlint проверяет, соответствуют ли ваши сообщения коммитов [формату conventional commits](https://www.conventionalcommits.org/ru/v1.0.0/)

```
<тип>[необязательная область]: <описание>

[необязательное тело]

[необязательный нижний колонтитул]
```

- [Зачем использовать Conventional Commits?](https://www.conventionalcommits.org/ru/v1.0.0/#%D0%B7%D0%B0%D1%87%D0%B5%D0%BC-%D0%B8%D1%81%D0%BF%D0%BE%D0%BB%D1%8C%D0%B7%D0%BE%D0%B2%D0%B0%D1%82%D1%8C-conventional-commits)

## Содержание

- [commitlint](#commitlint)
  - [Содержание](#содержание)
  - [Что входит в комплект](#что-входит-в-комплект)
  - [Установка](#установка)
    - [Релизы](#релизы)
    - [Используя go](#используя-go)
    - [Сборка из исходного кода](#сборка-из-исходного-кода)
  - [CLI приложение](#cli-приложение)
    - [Быстрый тест](#быстрый-тест)
    - [Настройка](#настройка)
    - [Команды](#команды)
      - [config](#config)
      - [lint](#lint)
      - [hook](#hook)
      - [debug](#debug)
  - [Веб-сервер для webhook'ов](#веб-сервер-для-webhookов)
    - [Быстрый старт](#быстрый-старт)
    - [Конфигурация сервера](#конфигурация-сервера)
    - [Настройка Gitea](#настройка-gitea)
    - [Запуск сервера](#запуск-сервера)
    - [Мониторинг](#мониторинг)
  - [Конфигурация по умолчанию](#конфигурация-по-умолчанию)
    - [Типы коммитов](#типы-коммитов)
  - [Доступные правила](#доступные-правила)
  - [Доступные форматеры](#доступные-форматеры)
  - [Расширяемость](#расширяемость)
  - [FAQ](#faq)
  - [Лицензия](#лицензия)

## Что входит в комплект

Этот проект предоставляет два приложения:

1. **CLI приложение** (`commitlint`) - инструмент командной строки для проверки сообщений коммитов
2. **Веб-сервер** (`commitlint-server`) - webhook сервер для автоматической проверки коммитов в pull request'ах Gitea

Оба приложения могут быть собраны и использованы независимо друг от друга.

## Установка

### Релизы

Скачайте бинарные файлы из [релизов](https://github.com/conventionalcommit/commitlint/releases) и добавьте их в ваш `PATH`.

### Используя go

```bash
# Установка CLI приложения
go install github.com/conventionalcommit/commitlint@latest

# Установка веб-сервера
go install github.com/conventionalcommit/commitlint/cmd/server@latest
```

### Сборка из исходного кода

```bash
git clone https://github.com/conventionalcommit/commitlint.git
cd commitlint

# Собрать оба приложения
make build

# Или собрать по отдельности
make build-cli      # Собрать только CLI
make build-server   # Собрать только веб-сервер

# Показать все доступные команды make
make help
```

## CLI приложение

### Быстрый тест

```bash
echo "wrong commit message" | commitlint lint --message=-
echo "feat: add new feature" | commitlint lint --message=-
```

### Настройка

#### Ручная настройка

```bash
# Инициализация конфигурационного файла
commitlint init

# Установка git hook для автоматической проверки
commitlint hook install
```

### Команды

#### config

```bash
# Показать текущую конфигурацию
commitlint config

# Показать конфигурацию с комментариями
commitlint config --verbose
```

#### lint

```bash
# Проверить сообщение коммита из файла
commitlint lint --message=commit-message.txt

# Проверить сообщение коммита из stdin
echo "feat: добавить новую функцию" | commitlint lint --message=-

# Проверить с пользовательским конфигом
commitlint lint --config=custom-config.yml --message=commit-message.txt

# Проверить определенный диапазон коммитов
commitlint lint --from=HEAD~5 --to=HEAD
```

##### Приоритет

###### Конфигурация

1. Путь из переменной окружения (`COMMITLINT_CONFIG`)
2. Файл `commitlint.yaml` в текущей директории
3. Конфигурация по умолчанию

###### Сообщение

1. Флаг `--message` (файл или stdin через `-`)
2. Диапазон коммитов `--from` и `--to`
3. Последний коммит

#### hook

```bash
# Установить git hook
commitlint hook install

# Удалить git hook  
commitlint hook uninstall

# Запустить проверку как git hook
commitlint hook run
```

#### debug

```bash
# Показать отладочную информацию
commitlint debug

# Показать информацию для определенного конфига
commitlint debug --config=custom-config.yml
```

## Веб-сервер для webhook'ов

Веб-сервер `commitlint-server` предназначен для автоматической проверки сообщений коммитов в pull request'ах Gitea через webhook'и.

### Быстрый старт

1. Создайте конфигурационный файл на основе примера:
```bash
cp commitlint-server.example.yml commitlint-server.yml
```

2. Отредактируйте конфигурацию:
```yaml
server:
  port: 8080
  webhook_secret: "ваш-секретный-ключ"
  gitea_url: "https://ваш-gitea.com"
  gitea_token: "ваш-токен-доступа"
```

3. Запустите сервер:
```bash
./commitlint-server
```

### Конфигурация сервера

Все параметры конфигурации находятся в секции `server` конфигурационного файла:

| Параметр | Тип | По умолчанию | Описание |
|----------|-----|--------------|----------|
| `address` | string | `"0.0.0.0"` | IP адрес для прослушивания |
| `port` | int | `8080` | Порт для прослушивания |
| `webhook_secret` | string | - | Секретный ключ для проверки webhook'ов |
| `gitea_url` | string | - | URL вашего экземпляра Gitea |
| `gitea_token` | string | - | Токен доступа к API Gitea |
| `cert_file` | string | - | Путь к SSL сертификату (опционально) |
| `key_file` | string | - | Путь к приватному ключу SSL (опционально) |
| `debug` | bool | `false` | Включить отладочное логирование |

### Настройка Gitea

1. **Создайте токен доступа:**
   - Перейдите в настройки пользователя → Токены доступа
   - Создайте новый токен с правами:
     - `repo` (для чтения репозиториев)
     - `repo:status` (для записи статусов коммитов)

2. **Настройте webhook в репозитории:**
   - Перейдите в настройки репозитория → Webhook'и
   - Добавьте новый webhook:
     - URL: `http://ваш-сервер:8080/webhook`
     - Тип содержимого: `application/json`
     - Секрет: тот же, что в конфигурации сервера
     - События: `Pull requests`

### Запуск сервера

```bash
# Запуск с конфигурацией по умолчанию
./commitlint-server

# Запуск с пользовательской конфигурацией
COMMITLINT_CONFIG=my-config.yml ./commitlint-server

# Запуск в фоновом режиме
nohup ./commitlint-server > server.log 2>&1 &

# Запуск с systemd (создайте service файл)
sudo systemctl start commitlint-server
```

### Мониторинг

Сервер предоставляет endpoint для проверки состояния:

```bash
# Проверка здоровья сервера
curl http://localhost:8080/health
```

Логи сервера содержат информацию о:
- Получении webhook'ов
- Обработке pull request'ов
- Результатах проверки коммитов
- Ошибках API

Пример вывода лога:
```
2024/01/15 10:30:15 Сервер запущен на 0.0.0.0:8080
2024/01/15 10:30:45 Получен webhook: тип=pull_request
2024/01/15 10:30:45 Начинаем обработку PR #123 в owner/repo
2024/01/15 10:30:46 Проверяем коммит a1b2c3d4: feat: добавить новую функцию
2024/01/15 10:30:46 ✓ Коммит a1b2c3d4 прошел проверку
2024/01/15 10:30:47 ✓ Все 3 коммитов в PR #123 прошли проверку
```

## Конфигурация по умолчанию

```yaml
version: "v1.6.0"
formatter: "default"
rules:
  - "type-enum"
  - "type-case"
  - "type-empty"
  - "scope-case"
  - "subject-case"
  - "subject-empty"
  - "subject-full-stop"
  - "header-max-length"
  - "body-leading-blank"
  - "footer-leading-blank"
severity:
  default: "error"
settings:
  type-enum:
    argument:
      - "feat"
      - "fix"
      - "docs"
      - "style"
      - "refactor"
      - "perf"  
      - "test"
      - "build"
      - "ci"
      - "chore"
      - "revert"
  header-max-length:
    argument: 100
```

### Типы коммитов

По умолчанию поддерживаются следующие типы коммитов:

| Тип | Описание |
|-----|----------|
| `feat` | Новая функциональность |
| `fix` | Исправление ошибки |
| `docs` | Изменения в документации |
| `style` | Изменения форматирования кода |
| `refactor` | Рефакторинг кода |
| `perf` | Улучшения производительности |
| `test` | Добавление или изменение тестов |
| `build` | Изменения системы сборки |
| `ci` | Изменения CI/CD |
| `chore` | Прочие изменения |
| `revert` | Откат предыдущих изменений |

## Доступные правила

- `type-enum` - проверяет, что тип коммита из разрешенного списка
- `type-case` - проверяет регистр типа коммита
- `type-empty` - проверяет, что тип коммита не пустой
- `scope-case` - проверяет регистр области
- `subject-case` - проверяет регистр описания
- `subject-empty` - проверяет, что описание не пустое
- `subject-full-stop` - проверяет отсутствие точки в конце описания
- `header-max-length` - проверяет максимальную длину заголовка
- `body-leading-blank` - проверяет пустую строку перед телом
- `footer-leading-blank` - проверяет пустую строку перед нижним колонтитулом

## Доступные форматеры

- `default` - стандартный форматер
- `json` - JSON форматер
- `junit` - JUnit XML форматер

## Расширяемость

commitlint поддерживает создание пользовательских правил и форматеров. См. документацию разработчика для подробностей.

## FAQ

**В: Можно ли использовать сервер без CLI приложения?**
О: Да, приложения полностью независимы. Вы можете установить и использовать только веб-сервер.

**В: Поддерживает ли сервер другие Git платформы кроме Gitea?**
О: В данный момент поддерживается только Gitea. Поддержка GitHub/GitLab может быть добавлена в будущем.

**В: Как настроить SSL для сервера?**
О: Укажите пути к сертификату и приватному ключу в конфигурации:
```yaml
server:
  cert_file: "/path/to/server.crt"
  key_file: "/path/to/server.key"
```

**В: Можно ли запускать несколько экземпляров сервера?**
О: Да, но убедитесь, что они используют разные порты или адреса.

## Лицензия

Этот проект лицензирован под лицензией MIT. См. файл [LICENSE.md](LICENSE.md) для подробностей.