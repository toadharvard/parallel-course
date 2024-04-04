
# Задача 2
## Использованный проект
https://github.com/AdayevKP/ParallelDataStructures

## Ошибки

TSAN и Helgrind нашли гонки данных в LazySyncSet.Add: [./Helgrind-fix-build.txt](./Helgrind-fix-build.txt) и [./TSAN-fix-build.txt](./TSAN-fix-build.txt)

При этом Helgrind и TSAN выдают разное количество предупреждений, но все об одной гонке данных в LazySyncSet.Add.

## Исправления
1. Исправлена ошибка компиляции заменой имени переменной size.
2. Исправлена гонка данных в LazySyncSet.Add почти так же, как это сделано в FineGrainedSyncSet.

## Анти-исправления

Кроме тупый вариантов: удалить мьютекс или поменять их порядок так, чтобы получился дедлок не было придумано.

Потому что даже полностью убрав все while true в исправленном LazySyncSet он продолжил работать без гонок.

Вот что выводит TSAN ([TSAN-new-race.txt](./TSAN-fixed-with-race.txt)), если "случайно"
```c++
pthread_mutex_lock(&root->mutex);
pthread_mutex_lock(&root->next->mutex);
```
поменять местами
```c++
pthread_mutex_lock(&root->next->mutex);
pthread_mutex_lock(&root->mutex);
```

При этом Helgrind не находит эту гонку данных за приемлимое время, скорее всего он просто впадает в состояние дедлока и умирает.
