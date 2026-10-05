# :bell: Notify Service

Микросервис для проекта **Date Wishlist Hub**.  

Ссылка на центральный репозиторий проекта: **[Date Wishlist Hub Deploy](https://github.com/alexgul25/date-wishlist-hub-deploy)**

Ссылка на канбан-доску проекта: **[Date Wishlist Hub - Development](https://github.com/users/alexgul25/projects/2)**

*Стек технологий сервиса:* `Go`  `gRPC`  `Kafka`  `PostgreSQL`

## :bulb: Описание сервиса

**Notify Service** - внутренний сервис, вычитывает события из брокера сообщений и организует логику отправки уведомлений пользователям.

- Для идемпотентной обработки прочитанных событий **реализован паттерн Inbox**.
- Для получения почтовых адресов пользователей посылаются запросы в **[User Service](https://github.com/alexgul25/user-svc)** через `gRPC-Client` (Protobuf-контракты определены публично в **[Protos](https://github.com/alexgul25/protos)**).
- В качестве брокера сообщений используется `Kafka`.
- В качестве БД используется `PostgreSQL`.
- Отправка уведомлений на электронную почту симулируется с помощью логирования.

<!-- markdownlint-disable MD033 -->
<details>
<summary>Примечание</summary>

На данном этапе сервис читает и обрабатывает только один тип событий: создание пользователем нового места для посещения в своём списке. Несмотря на это, система спроектирована с расчётом на лёгкое расширение при появлении новых типов событий в будущем.

</details>
<!-- markdownlint-enable MD033 -->

## :gear: Структура сервиса

:open_file_folder: **[cmd](./cmd/)** - команды запуска приложения.

:open_file_folder: **[migrations](./migrations/)** - файлы миграций.

:open_file_folder: **[internal/app](./internal/app/)** - сборка всех компонентов в единое приложение.

:open_file_folder: **[internal/config](./internal/config/)** - работа с файлами конфигурации.

:open_file_folder: **[internal/domain](./internal/domain/)** - определения доменных сущностей.

:open_file_folder: **[internal/inbox](./internal/inbox/)** - идемпотентный консьюмер брокера сообщений.

:open_file_folder: **[internal/infrastructure](./internal/infrastructure/)** - конкретные реализации абстрактных сущностей, используемых для работы приложения.

:open_file_folder: **[internal/lib](./internal/lib/)** - общие вспомогательные функции и утилиты.

:open_file_folder: **[internal/service](./internal/service/)** - сервисный слой (бизнес-логика).

:open_file_folder: **[internal/storage](./internal/storage/)** - слой хранения данных.

## :desktop_computer: Локальный запуск и работа через терминал

В данном разделе приведена инструкция по запуску одного **Notify Service**. Сервис можно запустить двумя способами: собрать и запустить бинарник локально или запустить в Docker-контейнере.

Инструкция по запуску всего проекта целиком доступна по **[ссылке](https://github.com/alexgul25/date-wishlist-hub-deploy#desktop_computer-локальный-запуск-и-работа-через-терминал)**.

### 1. Подготовка окружения

В вашем дистрибутиве должны быть установлены и готовы к работе:

- актуальная для проекта версия Go (см. [go.mod](./go.mod)) - для локального запуска;
- Docker Engine с плагином `buildx` - для запуска в контейнере;
- сервер PostgreSQL (версия 13+) и утилита `psql`;
- сервер Kafka (версия 4.0+);
- утилита `make`.

### 2. Клонирование репозитория

Клонируйте этот репозиторий c помощью HTTP или SSH.

```bash
git clone https://github.com/alexgul25/notify-svc.git
```

```bash
git clone git@github.com:alexgul25/notify-svc.git
```

### 3. Настройка инфраструктуры

#### 3.1. PostgreSQL

Запустите сервер PostgreSQL, затем создайте пользователя и базу данных для **Notify Service**.

```bash
sudo -u postgres psql -c "CREATE USER <имя пользователя> WITH PASSWORD '<пароль>';"
```

```bash
sudo -u postgres psql -c "CREATE DATABASE <имя БД> OWNER <имя пользователя>;"
```

Проверьте доступ.

```bash
psql -h localhost -U <имя пользователя> -d <имя БД> -c "SELECT 1;"
```

Если всё работает корректно, вы увидите следующий вывод:

```bash
 ?column? 
----------
        1
(1 row)
```

#### 3.2. Kafka

Запустите сервер Kafka. Если у вас отключено автоматическое создание топиков, создайте их самостоятельно (названия топиков см. в **[topics.go](./internal/inbox/topics.go)**).

#### 3.3. User Service

Для корректной работы необходимо запустить ещё один сервис проекта **Date Wishlist Hub**. Подробную инструкцию можно найти по **[этой ссылке](https://github.com/alexgul25/user-svc#desktop_computer-локальный-запуск-и-работа-через-терминал)**.

#### 3.4. Файл конфигурации

***ВАЖНО!*** Создайте в корневой папке репозитория файл `.env` для переменных окружения и заполните его (см [.env.example](.env.example)).

Для переменных `DB_USER`, `DB_PASSWORD` и `DB_NAME` используйте значения из шага [3.1.](#31-postgresql)

Для переменной `KAFKA_CONSUMER_BROKERS` используйте значения из шага [3.2.](#32-kafka)

Для переменной `USER_SERVICE_ADDR` используйте значение из шага [3.3.](#33-user-service)

<!-- markdownlint-disable MD033 -->
<details>
<summary>Особенности .env при запуске в Docker</summary>

- Значения указывайте без кавычек: Docker передаёт их в контейнер как есть, вместе с кавычками.
- `localhost` в `DB_HOST`, `KAFKA_CONSUMER_BROKERS` и `USER_SERVICE_ADDR` внутри контейнера означает сам контейнер, а не вашу машину (см. подсказки к варианту запуска в Docker в [следующем шаге](#4-запуск-и-работа)).

</details>
<!-- markdownlint-enable MD033 -->

### 4. Запуск и работа

#### Вариант 1. Локальный запуск

Для удобства локальной работы в корне репозитория определён Makefile.

1. `make help` - узнайте о доступных командах.
2. `make run` - примените миграции, соберите бинарник и запустите сервис.
3. `CTRL + C` - отправьте сервису сигнал завершения, когда закончите работу.

#### Вариант 2. Запуск в Docker

В корне репозитория определены **[Dockerfile](./Dockerfile)** и **[.dockerignore](./.dockerignore)**. Итоговый образ содержит бинарник сервиса, бинарник мигратора и папку с миграциями, переменные окружения передаются в контейнер при запуске.

1. `docker build -t notify-svc .` - соберите образ.
2. `docker run --rm --network host --env-file .env --entrypoint /app/migrator notify-svc` - примените миграции.
3. `docker run --rm --name notify-svc --network host --env-file .env notify-svc` - запустите контейнер с сервисом.
4. `CTRL + C` или `docker stop notify-svc` из другого терминала - отправьте сервису сигнал завершения, когда закончите работу.

<!-- markdownlint-disable MD033 -->
<details>
<summary>Подсказки</summary>

- Флаг `--network host` запускает контейнер в сети вашей машины: `localhost` в `DB_HOST`, `KAFKA_CONSUMER_BROKERS` и `USER_SERVICE_ADDR` указывает на локальные PostgreSQL, Kafka и User Service. Режим работает в Docker Engine на Linux (в том числе в WSL2).
- Сервис не принимает входящих соединений, поэтому публиковать порты (`-p`) не нужно. Если PostgreSQL, Kafka и User Service доступны контейнеру по сети (например, запущены в других контейнерах), вместо `--network host` укажите в `.env` их сетевые адреса.
- Если при сборке не удаётся скачать Go-модули (например, `proxy.golang.org` недоступен), передайте другой прокси через аргумент сборки: `docker build --build-arg GOPROXY=https://goproxy.io,direct -t notify-svc .`
- Чтобы запустить контейнер в фоне, замените `--rm` на `-d`. Логи сервиса можно посмотреть командой `docker logs -f notify-svc`, остановить и удалить контейнер - командами `docker stop notify-svc` и `docker rm notify-svc`.

</details>
<!-- markdownlint-enable MD033 -->

#### Демонстрация работы

Для демонстрации работы сервиса необходимо посылать события в брокер сообщений. Это можно сделать с помощью **Place Service** (**[инструкция по запуску](https://github.com/alexgul25/place-svc#desktop_computer-локальный-запуск-и-работа-через-терминал)**).
