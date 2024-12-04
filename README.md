# Music Library API

REST API сервис для управления музыкальной библиотекой.

## Функциональность

- Получение списка песен с фильтрацией и пагинацией
- Получение текста песни с пагинацией по куплетам
- Добавление новых песен с получением информации из внешнего API
- Обновление информации о песнях
- Удаление песен

## Технологии

- Go 1.21
- PostgreSQL
- Gorilla Mux (маршрутизация)
- Zap (логирование)
- Golang-migrate (миграции БД)
- Swagger (документация API)

## Установка и запуск

1. Клонируйте репозиторий:
```bash
git clone https://github.com/silberRus/testGOforEffectiveMobile.git
```

2. Создайте файл .env со следующими параметрами:
```env
# Database configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=music_library
DB_SSL_MODE=disable

# Server configuration
SERVER_PORT=8080
```

3. Установите зависимости:
```bash
go mod download
```

4. Запустите PostgreSQL

5. Запустите приложение:
```bash
go run cmd/main.go
```

## API Endpoints

### GET /api/v1/songs
Получение списка песен с фильтрацией и пагинацией.

**Query параметры:**
- `group_name` - название группы
- `song_name` - название песни
- `from_date` - начальная дата (формат: YYYY-MM-DD)
- `to_date` - конечная дата (формат: YYYY-MM-DD)
- `text` - поиск по тексту песни
- `link` - поиск по ссылке
- `page` - номер страницы
- `page_size` - размер страницы

### GET /api/v1/songs/{id}/lyrics
Получение текста песни с пагинацией по куплетам.

**Path параметры:**
- `id` - ID песни

**Query параметры:**
- `page` - номер страницы
- `page_size` - размер страницы

### POST /api/v1/songs
Добавление новой песни.

**Body:** JSON объект с информацией о песне
```json
{
    "group_name": "string",
    "song_name": "string",
    "text": "string",
    "link": "string"
}
```

### PUT /api/v1/songs/{id}
Обновление информации о песне.

**Path параметры:**
- `id` - ID песни

**Body:** JSON объект с обновленной информацией о песне (аналогичен POST)

### DELETE /api/v1/songs/{id}
Удаление песни.

**Path параметры:**
- `id` - ID песни

Swagger документация доступна по адресу: http://localhost:8080/swagger/
где localhost:8080 - адрес вашего сервера (нужно изменить в файле .env)

# Это не рабочий проект

Это тестовое задание для проверки знаний по Go. Специально для Effective Mobile (https://effective-mobile.ru/).