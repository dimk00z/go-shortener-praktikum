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