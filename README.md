# Сервис бронирования столиков в ресторанах

## Содержание

- [Задание](#Задание)
- [Подготовка](#Подготовка)
- [Запуск](#Запуск)
- [Зависимости](#Зависимости)

## Тестовое задание:

Представь знакомую всем ситуацию. Ты собираешься с друзьями в ресторан и хочешь забронировать столик. В вашем городе
открыты 3 ресторана: «Каравелла», «Молодость», «Мясо и Салат» (любые совпадения случайны).

**У ресторанов свои особенности:**

1. «Каравелла»:

   — 10 столиков: 6 столиков вмещает до 4 человек. 2 столика вмещают до 3 человек. 2 столика вмещает до 2 человек;

   — Среднее время ожидания блюда: 30 минут;

   — Средний чек: 2000 рублей.

2. «Молодость»

   — 3 столика, каждый из которых вмещает 3 человека;

   — Среднее время ожидания блюда: 15 минут;

   — Средний чек: 1000 рублей.

3. «Мясо и Салат»

   — 6 столиков: 2 столика вмещает до 8 человек. 4 столика вмещают до 3 человек;

   — Среднее время ожидания блюда: 60 минут;

   — Средний чек: 1500 рублей.

**Необходимо создать небольшую систему бронирования столиков.** Предусматриваются следующие шаги:

1. Пользователь указывает количество человек и желаемое время посещения ресторана.
2. Время брони - 2 часа. Рестораны работают с 9:00 до 23:00 (Последнюю бронь можно создать на 21:00).
3. Сдвигать столики - можно.
4. Система предлагает доступные варианты (рестораны).
5. Необходимо указывать актуальное количество свободных мест.
6. Необходимо отсортировать подходящие варианты по возрастанию среднего времени ожидания и среднего чека.
7. Необходимо скрывать недоступные варианты.
8. Пользователь указывает имя и номер телефона и завершает процесс бронирования.

**В качестве результата ожидается** ссылка на github-репозиторий, в котором находятся:

- [x] http-сервер, который общается с миром на языке REST API.
- [x] readme.md-файл, в котором подробно описана инструкция по установке системы

  — Все действия по установке должны быть автоматизированы (= вызов команд)

- [x] sql-дамп базы данных.

**Нужно не забыть:**

- [x] Проверить входные параметры. Сделать защиту от дурака.
- [x] Написать много комментариев к коду.
- [x] Проверить полученный результат дважды (а лучше трижды).

**Будет круто, если** (но совсем не критично, если не получится):

- [x] В качестве БД ты выберешь postgres.
- [x] Сделаешь визуальный интерфейс. Можно сверстать самому или использовать готовые решения.
- [x] Напишешь код на php или go (вообще идеально).
- [x] Используешь docker-compose.
- [x] Вынесешь все настройки подключения к базе в переменные окружения (env-параметры).

## Подготовка

Для запуска нужно установить 2 компоннта пакет языка на котором
написан сервер API и базу данных PosgreSql:

- [Go v1.18.4](https://go.dev/dl/)
- [Docker](https://www.docker.com/)
- [PostgreSql](https://hub.docker.com/_/postgres)

Далее для запуска необходимо скачать исходный код
сервиса бронирования:

- [Booking_restaurant]

## Запуск

Запуск осуществлятся срдствами:
- [Makefile]
- [docker-compose.yml]

Для запуска docker-compose можно воспользоваться командой:

```shell
# запуск компонентов в отдельных Docker-контейнерах
make compose-up
```
## Зависимости

* В качестве Роутра я взял [Gin](https://github.com/gin-gonic/gin):
Причины выбора:
1. Роутер GIN иапользует свой контекст.
2. У него много вспомагательных функций в контексте(различные парсинги например с JSON).
3. Он очень популярен, а значит будет поддерживаться еще долгое время.

* Драйвер общния с БД я взял из стандартного пакета Go [database/sql]
Причины выбора:
1. Большой функционал прдставлен в одном пакет
2. Хорошо протестирован

* Логгирование [logrus](https://github.com/sirupsen/logrus):
Причины выбора: 
1. Поддеерживается большим сообществом

* Для тестирования я воспользовался фреймворком [Check](gopkg.in/check.v1)

* Архитктуру выбрал хексоганальную:
Причины выбора:
1. Слои общаются между собой с помощью адаптеров в виде интерфейсов,
это приемущество дает нам легко мнять компоннты(масштабировать)
2. О базе данных не знает не один слой это дает нам безопасность
