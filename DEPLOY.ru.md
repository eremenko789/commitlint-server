# Руководство по развертыванию Commitlint Server

Это руководство описывает различные способы развертывания веб-сервера commitlint для обработки webhook'ов от Gitea.

## Содержание

- [Руководство по развертыванию Commitlint Server](#руководство-по-развертыванию-commitlint-server)
  - [Содержание](#содержание)
  - [Требования](#требования)
  - [Подготовка](#подготовка)
  - [Способы развертывания](#способы-развертывания)
    - [1. Нативное развертывание](#1-нативное-развертывание)
    - [2. Развертывание с systemd](#2-развертывание-с-systemd)
    - [3. Развертывание с Docker](#3-развертывание-с-docker)
    - [4. Развертывание с Docker Compose](#4-развертывание-с-docker-compose)
  - [Настройка Gitea](#настройка-gitea)
  - [Мониторинг и обслуживание](#мониторинг-и-обслуживание)
  - [Безопасность](#безопасность)
  - [Устранение неполадок](#устранение-неполадок)

## Требования

- Go 1.23+ (для сборки из исходного кода)
- Linux/macOS/Windows сервер
- Доступ к экземпляру Gitea
- Токен доступа к API Gitea
- Опционально: Docker и Docker Compose

## Подготовка

1. **Получите токен доступа Gitea:**
   - Войдите в Gitea
   - Перейдите в Настройки → Токены доступа
   - Создайте новый токен с правами:
     - `repo` (чтение репозиториев)
     - `repo:status` (запись статусов коммитов)

2. **Подготовьте конфигурацию:**
   ```bash
   cp commitlint-server.example.yml commitlint-server.yml
   ```

3. **Отредактируйте конфигурацию:**
   ```yaml
   server:
     port: 8080
     webhook_secret: "ваш-секретный-ключ"
     gitea_url: "https://ваш-gitea.com"
     gitea_token: "ваш-токен-доступа"
     debug: false
   ```

## Способы развертывания

### 1. Нативное развертывание

**Сборка:**
```bash
# Клонирование репозитория
git clone https://github.com/conventionalcommit/commitlint.git
cd commitlint

# Сборка сервера
make build-server
```

**Запуск:**
```bash
# Запуск с конфигурацией по умолчанию
./commitlint-server

# Запуск с пользовательской конфигурацией
COMMITLINT_CONFIG=./commitlint-server.yml ./commitlint-server

# Запуск в фоновом режиме
nohup ./commitlint-server > server.log 2>&1 &
```

### 2. Развертывание с systemd

**Установка:**
```bash
# Создание пользователя
sudo useradd -r -s /bin/false commitlint

# Создание директорий
sudo mkdir -p /etc/commitlint /var/log/commitlint
sudo chown commitlint:commitlint /var/log/commitlint

# Копирование файлов
sudo cp commitlint-server /usr/local/bin/
sudo cp commitlint-server.yml /etc/commitlint/
sudo cp deployments/systemd/commitlint-server.service /etc/systemd/system/

# Установка прав
sudo chmod +x /usr/local/bin/commitlint-server
sudo chmod 644 /etc/commitlint/commitlint-server.yml
```

**Управление службой:**
```bash
# Включение автозапуска
sudo systemctl enable commitlint-server

# Запуск службы
sudo systemctl start commitlint-server

# Проверка статуса
sudo systemctl status commitlint-server

# Просмотр логов
sudo journalctl -u commitlint-server -f

# Остановка службы
sudo systemctl stop commitlint-server

# Перезапуск службы
sudo systemctl restart commitlint-server
```

### 3. Развертывание с Docker

**Сборка образа:**
```bash
docker build -f Dockerfile.server -t commitlint-server .
```

**Запуск контейнера:**
```bash
# Создание директории для конфигурации
mkdir -p ./config

# Копирование конфигурации
cp commitlint-server.yml ./config/

# Запуск контейнера
docker run -d \
  --name commitlint-server \
  --restart unless-stopped \
  -p 8080:8080 \
  -v $(pwd)/config/commitlint-server.yml:/etc/commitlint/commitlint-server.yml:ro \
  -e COMMITLINT_CONFIG=/etc/commitlint/commitlint-server.yml \
  commitlint-server
```

**Управление контейнером:**
```bash
# Просмотр логов
docker logs -f commitlint-server

# Остановка контейнера
docker stop commitlint-server

# Запуск контейнера
docker start commitlint-server

# Удаление контейнера
docker rm -f commitlint-server
```

### 4. Развертывание с Docker Compose

**Запуск:**
```bash
# Создание конфигурации
cp commitlint-server.example.yml commitlint-server.yml
# Отредактируйте конфигурацию

# Запуск сервисов
docker-compose up -d

# Просмотр логов
docker-compose logs -f

# Остановка сервисов
docker-compose down

# Перезапуск сервисов
docker-compose restart
```

## Настройка Gitea

1. **Перейдите в настройки репозитория:**
   - Откройте ваш репозиторий в Gitea
   - Перейдите в Настройки → Webhook'и

2. **Добавьте новый webhook:**
   - URL: `http://ваш-сервер:8080/webhook`
   - Тип содержимого: `application/json`
   - Секрет: тот же, что в конфигурации сервера
   - Триггеры: отметьте `Pull requests`
   - Активен: ✓

3. **Протестируйте webhook:**
   - Создайте тестовый pull request
   - Проверьте логи сервера на наличие обработки webhook'а

## Мониторинг и обслуживание

**Проверка здоровья:**
```bash
curl http://localhost:8080/health
```

**Мониторинг логов:**
```bash
# Нативное развертывание
tail -f server.log

# systemd
sudo journalctl -u commitlint-server -f

# Docker
docker logs -f commitlint-server

# Docker Compose
docker-compose logs -f
```

**Ротация логов (для нативного развертывания):**
```bash
# Создание файла logrotate
sudo tee /etc/logrotate.d/commitlint-server << EOF
/var/log/commitlint/server.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    postrotate
        systemctl reload commitlint-server
    endscript
}
EOF
```

## Безопасность

**Рекомендации:**

1. **Используйте HTTPS:**
   ```yaml
   server:
     cert_file: "/path/to/server.crt"
     key_file: "/path/to/server.key"
   ```

2. **Настройте файрвол:**
   ```bash
   # ufw
   sudo ufw allow 8080/tcp
   
   # iptables
   sudo iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
   ```

3. **Ограничьте доступ по IP:**
   ```bash
   # Разрешить доступ только с Gitea сервера
   sudo iptables -A INPUT -p tcp --dport 8080 -s IP_GITEA_СЕРВЕРА -j ACCEPT
   sudo iptables -A INPUT -p tcp --dport 8080 -j DROP
   ```

4. **Используйте обратный прокси (nginx):**
   ```nginx
   server {
       listen 443 ssl;
       server_name commitlint.example.com;
       
       ssl_certificate /path/to/cert.pem;
       ssl_certificate_key /path/to/key.pem;
       
       location / {
           proxy_pass http://127.0.0.1:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   ```

## Устранение неполадок

**Частые проблемы:**

1. **Сервер не запускается:**
   - Проверьте правильность конфигурации
   - Убедитесь, что порт не занят другим процессом
   - Проверьте права доступа к конфигурационному файлу

2. **Webhook'и не обрабатываются:**
   - Проверьте URL webhook'а в настройках Gitea
   - Убедитесь, что секрет webhook'а совпадает
   - Проверьте доступность сервера с Gitea

3. **Ошибки API Gitea:**
   - Проверьте корректность токена доступа
   - Убедитесь, что токен имеет необходимые права
   - Проверьте доступность API Gitea

4. **Проблемы с SSL:**
   - Убедитесь в корректности путей к сертификатам
   - Проверьте права доступа к файлам сертификатов
   - Убедитесь, что сертификаты действительны

**Отладка:**

Включите отладочное логирование в конфигурации:
```yaml
server:
  debug: true
```

Это добавит подробную информацию о:
- Получаемых webhook'ах
- Обработке pull request'ов
- API вызовах к Gitea
- Результатах проверки коммитов