#!/bin/bash
set -e

# убиваем все старые процессы mycut
pkill -f mycut || true
sleep 0.5

# создаем тестовые файлы
echo -e "a\tb\tc\nd\te\tf\ng\th\ti" > data1.txt
echo -e "1\t2\t3\n4\t5\t6\n7\t8\t9" > data2.txt
echo -e "x\ty\tz\np\tq\tr\nu\tv\tw" > data3.txt

# запускаем координатор в фоне (только stdout в actual.txt)
./mycut --mode=coord --addr="localhost:9100" --peers="localhost:9101,localhost:9102,localhost:9103" --f="1,3" > actual.txt &
COORD_PID=$!

# ждем немного, чтобы координатор запустился
sleep 1

# запускаем серверы (их stderr будет видно в консоли)
./mycut --mode=server --addr="localhost:9101" --peers="localhost:9100" --f="1,3" < data1.txt 2>&1 &
PID1=$!
./mycut --mode=server --addr="localhost:9102" --peers="localhost:9100" --f="1,3" < data2.txt 2>&1 &
PID2=$!
./mycut --mode=server --addr="localhost:9103" --peers="localhost:9100" --f="1,3" < data3.txt 2>&1 &
PID3=$!

# ждем завершения всех процессов
wait $PID1 $PID2 $PID3 $COORD_PID

# проверка результата
echo "=== Expected ==="
cut -f1,3 data1.txt data2.txt data3.txt | sort
echo "=== Actual ==="
cat actual.txt
echo "=== Diff ==="

cut -f1,3 data1.txt data2.txt data3.txt | sort > expected.txt
if diff -u expected.txt actual.txt; then
    echo "✅ Test passed! Results match GNU cut."
else
    echo "❌ Test failed! Results differ."
fi

# удаляем временные файлы
rm -f data1.txt data2.txt data3.txt expected.txt actual.txt