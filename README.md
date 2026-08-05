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

### 1. Подготовка окружения

В вашем дистрибутиве должны быть установлены и готовы к работе:

- актуальная для проекта версия Go (см. [go.mod](./go.mod));

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

Для переменной `USER_SERVICE_ADDR` используйте значение из шага [3.3.](#33-user-service).

### 4. Запуск и работа

Для удобства локальной работы в корне репозитория определён Makefile.

1. `make run-svc` - примените миграции и запустите сервис.

2. `CTRL + C` - отправьте сервису сигнал завершения, когда закончите работу.

Для демонстрации работы сервиса необходимо посылать события в брокер сообщений. Это можно сделать с помощью **Place Service** (**[инструкция по запуску](https://github.com/alexgul25/place-svc#desktop_computer-локальный-запуск-и-работа-через-терминал)**).
