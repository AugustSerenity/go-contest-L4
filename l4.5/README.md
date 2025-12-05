### Оптимизация простого API-сервиса с профилировкой

### **[Задание](./docs/task.md)**

### **[API-сервис](./docs/README.md)**

Добавил в исходный проек `_ "net/http/pprof"`
Создал нагрузку 
```sh
wrk -t4 -c128 -d30s http://localhost:8080/events_for_day?user_id=1&date=2024-01-01
```
Вывод 
```sh
Running 30s test @ http://localhost:8080/events_for_day?user_id=1&date=2024-01-01
  4 threads and 128 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.58ms    1.61ms  62.70ms   92.39%
    Req/Sec    21.99k     3.71k   43.20k    85.58%
  2626633 requests in 30.04s, 475.94MB read
  Non-2xx or 3xx responses: 2626633
Requests/sec:  87437.73
Transfer/sec:     15.84MB
```
далее снимаю значение CPU профиля
```sh
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=20
(pprof) top
Showing nodes accounting for 30ms, 100% of 30ms total
Showing top 10 nodes out of 22
      flat  flat%   sum%        cum   cum%
      10ms 33.33% 33.33%       10ms 33.33%  runtime.(*unwinder).resolveInternal
      10ms 33.33% 66.67%       10ms 33.33%  runtime.kevent
      10ms 33.33%   100%       10ms 33.33%  runtime.pthread_cond_wait
         0     0%   100%       10ms 33.33%  runtime.(*unwinder).init (inline)
         0     0%   100%       10ms 33.33%  runtime.(*unwinder).initAt
         0     0%   100%       20ms 66.67%  runtime.findRunnable
         0     0%   100%       10ms 33.33%  runtime.gcBgMarkWorker
         0     0%   100%       10ms 33.33%  runtime.gcBgMarkWorker.func2
         0     0%   100%       10ms 33.33%  runtime.gcDrain
         0     0%   100%       10ms 33.33%  runtime.gcDrainMarkWorkerDedicated (inline)
```
значение Memory профиля
```sh
go tool pprof http://localhost:6060/debug/pprof/heap
(pprof) top
Showing nodes accounting for 1539.05kB, 100% of 1539.05kB total
Showing top 10 nodes out of 14
      flat  flat%   sum%        cum   cum%
     514kB 33.40% 33.40%      514kB 33.40%  bufio.NewWriterSize (inline)
     513kB 33.33% 66.73%      513kB 33.33%  runtime.allocm
  512.05kB 33.27%   100%   512.05kB 33.27%  runtime.(*scavengerState).init
         0     0%   100%      514kB 33.40%  net/http.(*conn).serve
         0     0%   100%      514kB 33.40%  net/http.newBufioWriterSize
         0     0%   100%   512.05kB 33.27%  runtime.bgscavenge
         0     0%   100%      513kB 33.33%  runtime.mstart
         0     0%   100%      513kB 33.33%  runtime.mstart0
         0     0%   100%      513kB 33.33%  runtime.mstart1
         0     0%   100%      513kB 33.33%  runtime.newm
```

значение Trace профиля
```sh
curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5
% Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 28163    0 28163    0     0   5626      0 --:--:--  0:00:05 --:--:--  6071
```
```sh
go tool trace trace.out
```
View trace by proc
<p align="center"> <img src="docs/image/noOptitraceByProc.png" alt="" width="80%" /> </p>

View trace by thread
<p align="center"> <img src="docs/image/noOptiThread.png" alt="" width="80%" /> </p>

Goroutine analysis
<p align="center"> <img src="docs/image/goroutine.png" alt="" width="80%" /> </p>

Garbage collection metrics
<p align="center"> <img src="docs/image/noOptGC.png" alt="" width="80%" /> </p>

Запускаю бенчмарк для `storage.go`, чтобы сравнить результаты до и после оптимизации. 
```sh
go test -bench=. -benchmem > bench_before.txt
```

Результат
<p align="center"> <img src="docs/image/bench_before.png" alt="" width="80%" /> </p>