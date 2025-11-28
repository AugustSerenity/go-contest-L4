### Микросервис “Календарь событий”
Сервис управления событиями (календарём) c поддержкой уведомлений, обновления, удаления и выборки событий за день, неделю или месяц

### [Задание](./docs/task.md)

### Функционал

Сервис предоставляет следующие функции:
- Управление событиями
- Создание события
- Обновление события
- Удаление события
- Получение событий:
    - За конкретный день
    - За неделю
    - За месяц
- Напоминания:
    - Возможность указать время напоминания (`RemindIn: "2h"`)
    - Фоновый воркер отправляет напоминания через канал `ReminderWorker`
- Архивация:
    - Сервис архивации удаляет старые события с помощью `ArchiveWorker`

## Начало работы
### Установка
Клонирование репозитория
```sh
git clone https://github.com/AugustSerenity/go-contest-L4/tree/main/l4.3
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
#### Создать событие
Запрос
```sh
curl -X POST http://localhost:8081/create_event \
  -H "Content-Type: application/json" \
  -d '{
        "user_id": 1,
        "name": "meeting",
        "date": "2030-10-10",
        "remind_in": "2h"
      }'
```
ответ 
```sh
{"result":{"event_id":1,"message":"successfully created"}}
```
запрос
```sh
curl -X POST http://localhost:8081/create_event \
  -H "Content-Type: application/json" \
  -d '{
        "user_id": 2,
        "name": "watching TV",
        "date": "2030-11-10",
        "remind_in": "10h"
      }
```
ответ
```sh
{"result":{"event_id":2,"message":"successfully created"}}
```
запрос
```sh
curl -X POST http://localhost:8081/create_event \
  -H "Content-Type: application/json" \
  -d '{
        "user_id": 3,
        "name": "reading",    
        "date": "2030-11-11",
        "remind_in": "1h" 
      }'
```
ответ 
```sh
{"result":{"event_id":3,"message":"successfully created"}}
```      

### События на день, на неделю, на месяц
`фильтруем по дню`

запрос
```sh
curl "http://localhost:8081/events_for_day?user_id=1&date=2030-10-10"
```
ответ
```sh
{"result":[{"id":1,"user_id":1,"event_name":"","date":"2030-10-10T00:00:00+03:00","remind_at":"2030-10-09T22:00:00+03:00","created_at":"2025-11-28T10:46:56.77808+03:00","is_active":true}]}
```
`фильтруем по неделе`

запрос
```sh
curl "http://localhost:8081/events_for_week?user_id=1&date=2030-10-10"
```
ответ
```sh
{"result":[{"id":1,"user_id":1,"event_name":"","date":"2030-10-10T00:00:00+03:00","remind_at":"2030-10-09T22:00:00+03:00","created_at":"2025-11-28T10:46:56.77808+03:00","is_active":true}]}
```
`фильтруем по месяцу`

запрос
```sh
curl "http://localhost:8081/events_for_month?user_id=1&date=2030-10-10"
```
ответ
```sh
{"result":[{"id":1,"user_id":1,"event_name":"","date":"2030-10-10T00:00:00+03:00","remind_at":"2030-10-09T22:00:00+03:00","created_at":"2025-11-28T10:46:56.77808+03:00","is_active":true}]}
```

### Обновить событие
запрос
```sh
curl -X POST http://localhost:8081/update_event \
  -H "Content-Type: application/json" \
  -d '{
        "user_id": 1,
        "name": "meeting",
        "date": "2030-10-10",
        "new_name": "updated",
        "new_date": "2030-10-11",
        "remind_in": "1h"
      }'
```
ответ
```sh
{"result":{"message":"successfully updated"}}
```
### Удалить событие
запрос
```sh
curl -X POST http://localhost:8081/delete_event \
  -H "Content-Type: application/json" \
  -d '{
        "user_id": 1,
        "event_name": "updated",
        "date": "2030-10-11"
      }'
```
ответ
```sh
{"result":{"message":"successfully deleted"}}
```