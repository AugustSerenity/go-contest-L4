### Утилита анализа GC и памяти (runtime, профилирование)
### [Задание](./docs/task.md) 
Эта утилита на Go собирает информацию о памяти и сборщике мусора (GC) в реальном времени и предоставляет её через HTTP-endpoint в формате, совместимом с Prometheus. Также подключено профилирование через pprof.

## Функционал

Программа:

- Считывает статистику памяти с помощью `runtime.ReadMemStats`.
- Управляет сборщиком мусора через `debug.SetGCPercent`.
- Предоставляет метрики через `/metrics` в формате Prometheus.
- Подключает профилирование с помощью `net/http/pprof`.

## Метрики Prometheus

Программа публикует следующие метрики:

| Метрика          | Описание                                   |
|-----------------|--------------------------------------------|
| `memory_Alloc`   | Количество используемой кучи (Heap Alloc) |
| `memory_Malloc`  | Общее количество аллокаций (`Mallocs`)    |
| `gc_Num`         | Количество завершённых сборок мусора      |
| `gc_Last`        | Время последнего завершённого GC (наносекунды) |
| `gc_percent`     | Текущий процент сборщика мусора (GC Percent) |

---

## Используемые технологии

- `runtime` — для получения информации о памяти и GC.
- `runtime/debug` — для управления параметрами GC.
- `prometheus/client_golang` — для экспорта метрик.
- `net/http/pprof` — для профилирования приложения.

---

## Начало работы
### Установка
Клонирование репозитория
```sh
git clone https://github.com/AugustSerenity/go-contest-L4/tree/main/l4.4
```
### Запуск сервиса
```sh
go run main.go
```
Сервер по умолчанию слушает на порту 8080.

### HTTP Endpoints
- Prometheus метрики:
  - [`http://localhost:8080/metrics`](http://localhost:8080/metrics)

- pprof профилирование:
    - [`http://localhost:8080/debug/pprof/`](http://localhost:8080/debug/pprof/)
    - [`http://localhost:8080/debug/pprof/cmdline`](http://localhost:8080/debug/pprof/cmdline)
    - [`http://localhost:8080/debug/pprof/profile`](http://localhost:8080/debug/pprof/profile)
    - [`http://localhost:8080/debug/pprof/symbol`](http://localhost:8080/debug/pprof/symbol)
    - [`http://localhost:8080/debug/pprof/trace`](http://localhost:8080/debug/pprof/trace)

### Пример запроса метрик
```sh
curl http://localhost:8080/metrics
```

### Профилирование pprof
```sh
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=10
```
После этого команда go tool pprof скачает профиль и откроет интерактивный CLI для анализа. Внутри можно использовать команды, например:
- `top` — показать функции с наибольшим потреблением CPU
- `list <func>` — показать исходный код функции с аннотациями времени CPU
- `web` — открыть визуализацию профиля в браузере (требует Graphviz)
- `pdf` — сохранить граф в PDF
- `help` — полный список команд
```sh
(pprof) top
(pprof) list main.someFunction
(pprof) web
(pprof) pdf > profile.pdf
(pprof) top10
```

### Сбор профиля памяти (heap)
```sh
go tool pprof http://localhost:8080/debug/pprof/heap
```
После этого можно посмотреть подбробный отчет о памяти 
```sh
(pprof) web
(pprof) list <имя_функции>
(pprof) pdf > heap_profile.pdf
(pprof) peek <имя_функции> # Показывает, сколько памяти использует конкретная функция
```





























