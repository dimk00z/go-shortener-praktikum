# go-musthave-shortener-tpl
Шаблон репозитория для практического трека «Go в веб-разработке».

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` - адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона выполните следующую команду:

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

# Решение

## Запуск проекта

```bash
cp .env_example .env
docker-compose up
```
Пример запуска с БД `go run main.go -d "postgres://admin:password@localhost:5432/shortener_db"`
Пример запуска с файловым хранилищем `go run main.go -f "storage.json"`

`http://localhost:8080/swagger/` - описание API

### Использование HTTPS

Для запуска сервера с использованием tls необходимы ключ и сертификат шифрования.
Пример для генерации самоподписанных ключевых данных:
```bash
openssl genrsa -out cert/server.key 2048
openssl req -new -x509 -sha256 -key cert/server.key -out cert/server.crt -days 3650
```
флаг `-s` или `env:ENABLE_HTTPS` - для запуска сервера с шифрованием

### gRPC
Реализована gRPC сервер следующего контракта: 

```
service Shortener {
  rpc Ping(EmptyRequest) returns (PingResponse);
  rpc GetStats(EmptyRequest) returns (StatsResponse);
  rpc GetByShortURL(ShortURLRequest) returns (URLResponse);
  rpc SaveURLFromText(URLRequest) returns (URLFromTextResponse);
  rpc SaveBatch(BatchSaveRequest) returns (BatchSaveResponse);
  rpc GetUsersURLs(EmptyRequest) returns (UserURLsResponse);
  rpc DelBatch(DelBatchRequest) returns (DelBatchResponse);
}
```

bool флаг `enable_gprc` или `env:ENABLE_GPRC` - для запуска gRPC сервера
флаг `g` или `env:GRPC_PORT` - порт gRPC сервера формата: ":9000"
