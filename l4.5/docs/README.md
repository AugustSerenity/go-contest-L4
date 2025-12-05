## HTTP-сервер «Календарь»

Это простой HTTP-сервер для управления календарём событий, реализующий базовые CRUD-операции. Сервер поддерживает создание, обновление, удаление и получение событий по дням, неделям и месяцам. Все события хранятся в оперативной памяти (in-memory storage).

### **[Задание](oldTask.md)**

## Начало работы
### Установка
Клонирование репозитория
```sh
git clone https://github.com/AugustSerenity/go-contest-L2/tree/main/l2.18
```
### Запуск сервиса
Запускаем сервер с помощью Makefile
```sh
make run
```

### Запуск тестов
Проверка vet, lint, data race, unit-тестов
```
make check
```
отдельный запуск unit-тестов c покрытием
```
make test
```

###  Функциональность
 Метод | Endpoint              | Описание                             |
|-------|------------------------|--------------------------------------|
| POST  | `/create_event`        | Создание события                     |
| POST  | `/update_event`        | Обновление события                   |
| POST  | `/delete_event`        | Удаление события                     |
| GET   | `/events_for_day`      | Получить события за день             |
| GET   | `/events_for_week`     | Получить события за неделю           |
| GET   | `/events_for_month`    | Получить события за месяц            |



### Формат запросов
#### создание события
Запросы
```
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"date":"2025-09-23 10:00:00","event":"My event name"}'
```
```
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"date":"2025-09-23 10:30:00","event":"My event name"}'
```
```
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"date":"2025-09-25 17:40:40","event":"My event name"}'
```
```
curl -X POST http://localhost:8080/create_event \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"date":"2025-10-01 10:00:00","event":"My event name"}'
```
Ответ
```
{"result":"successfully created"}
```
#### Отфильтруем события по времени
фильтруем по дню 
```
curl -X GET "http://localhost:8080/events_for_day?user_id=1&date=2025-09-23"

```
ответ 
```
{"result":[{"event_name":"","date":"2025-09-23T10:00:00Z"},{"event_name":"","date":"2025-09-23T10:30:00Z"}]}
```
фильтруем по неделе
```
curl -X GET "http://localhost:8080/events_for_week?user_id=1&date=2025-09-23"

```
ответ
```
{"result":[{"event_name":"","date":"2025-09-23T10:00:00Z"},{"event_name":"","date":"2025-09-23T10:30:00Z"},{"event_name":"","date":"2025-09-25T17:40:40Z"}]}
```
фильтруем по месяцу
```
curl -X GET "http://localhost:8080/events_for_month?user_id=1&date=2025-09-23"
```
ответ
```
{"result":[{"event_name":"","date":"2025-09-23T10:00:00Z"},{"event_name":"","date":"2025-09-23T10:30:00Z"},{"event_name":"","date":"2025-09-25T17:40:40Z"}]}
```
#### Обновляем событие
запрос
```
curl -X POST http://localhost:8080/update_event \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"date":"2025-09-23 10:30:00","event":"Fanta"}'      
```
ответ
```
{"result":"successfully updated"}
```
#### Удаляем событие
```
curl -X POST http://localhost:8080/delete_event \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"date":"2025-09-25 17:40:40","event":"My event name"}'
```
ответ 
```
{"result":"successfully deleted"}
```