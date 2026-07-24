# :bell: Notify Service

Микросервис для проекта **Date Wishlist Hub**.  

Ссылка на центральный репозиторий проекта: **[Date Wishlist Hub Deploy](https://github.com/alexgul25/date-wishlist-hub-deploy)**

Ссылка на канбан-доску проекта: **[Date Wishlist Hub - Development](https://github.com/users/alexgul25/projects/2)**

*Стек технологий сервиса:* `Go`  `gRPC`  `Kafka`  `PostgreSQL`

## :bulb: Описание сервиса

**Notify Service** - вычитывает события из брокера сообщений и организует логику отправки уведомлений пользователям.

- В качестве брокера сообщений используется `Kafka`.

- Для идемпотентной обработки сообщений ведётся таблица с помощью `PostgreSQL`. Также используются миграции БД.

- Для получения почтовых адресов пользователей посылаются запросы в **[User Service](https://github.com/alexgul25/user-svc)** через `gRPC-Client` (Protobuf-контракты определены публично в **[Protos](https://github.com/alexgul25/protos)**).

- Отправка уведомлений на электронную почту симулируется с помощью логирования.

<!-- markdownlint-disable MD033 -->
<details>
<summary>Примечание</summary>

На данном этапе сервис читает и обрабатывает только один тип событий: создание пользователем нового места для посещения в своём списке.

</details>
<!-- markdownlint-enable MD033 -->

## :gear: Структура сервиса

**[:open_file_folder: cmd](./cmd/)** - команды запуска приложения.

**[:open_file_folder: migrations](./migrations/)** - файлы миграций.

**[:open_file_folder: internal/app](./internal/app/)** - сборка всех компонентов в единое приложение.

**[:open_file_folder: internal/config](./internal/config/)** - работа с файлами конфигурации.

**[:open_file_folder: internal/domain](./internal/domain/)** - определения доменных сущностей.

**[:open_file_folder: internal/inbox](./internal/inbox/)** - идемпотентный консьюмер брокера сообщений.

**[:open_file_folder: internal/infrastructure](./internal/infrastructure/)** - конкретные реализации абстрактных сущностей, используемых для работы приложения.

**[:open_file_folder: internal/lib](./internal/lib/)** - общие вспомогательные функции и утилиты.

**[:open_file_folder: internal/service](./internal/service/)** - сервисный слой (бизнес-логика).

**[:open_file_folder: internal/storage](./internal/storage/)** - слой хранения данных.
