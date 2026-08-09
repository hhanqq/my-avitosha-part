# Карта backend-файлов и моего вклада

Ниже перечислены все файлы из каталога `backend/` и конкретная роль моей работы
в каждом из них. Полные файлы сохранены, чтобы код можно было читать вместе с
типами, импортами и package-контекстом.

## Сборка и зависимости

| Файл | Что реализовано |
|---|---|
| `backend/Dockerfile` | Сборка исполняемых файлов API Gateway, Auth Service и Game Service для микросервисного запуска. |
| `backend/go.mod` | Подключены зависимости gRPC, protobuf и компоненты, необходимые игровому и realtime-слою. |
| `backend/go.sum` | Зафиксированы контрольные суммы новых Go-зависимостей. |
| `backend/buf.yaml` | Настроен protobuf-модуль и lint-конфигурация Buf. |
| `backend/buf.gen.yaml` | Настроена воспроизводимая генерация Go и gRPC stubs. |

## Точки запуска сервисов

| Файл | Что реализовано |
|---|---|
| `backend/cmd/api/main.go` | Точка запуска публичного API Gateway. |
| `backend/cmd/auth/main.go` | Отдельный процесс Auth Service с gRPC transport и health checking. |
| `backend/cmd/game/main.go` | Отдельный Game Service, владеющий игровым доменом, AI и realtime hub. |

## REST, OpenAPI и gRPC-контракты

| Файл | Что реализовано |
|---|---|
| `backend/api/openapi.yaml` | Контракты игровых endpoint'ов: питомец, задания, действия, комната, история, summary, leaderboard, achievements, rewards, advice и WebSocket. |
| `backend/api/openapi_test.go` | Автоматическая проверка валидности и наличия обязательных частей OpenAPI-контракта. |
| `backend/api/proto/avitosha/v1/services.proto` | Версионированные AuthService и GameService RPC, команды, ответы, health/realtime-взаимодействие. |
| `backend/internal/gen/avitosha/v1/services.pb.go` | Сгенерированные protobuf-типы; generated-code, а не ручная реализация. |
| `backend/internal/gen/avitosha/v1/services_grpc.pb.go` | Сгенерированные gRPC clients/servers; generated-code. |

## Конфигурация и сборка приложения

| Файл | Что реализовано |
|---|---|
| `backend/internal/config/config.go` | Конфигурация AI, gRPC-адресов, HTTP gateway, timeout и проверка обязательных параметров. |
| `backend/internal/config/config_test.go` | Тесты default values, validation и различных режимов запуска сервисов. |
| `backend/internal/app/app.go` | Интеграция игровой логики, realtime и AI в исходный application wiring. |
| `backend/internal/app/game_wiring_test.go` | Проверка корректного подключения game service, repository, hub и advice generator. |
| `backend/internal/app/microservices.go` | Раздельная сборка API Gateway, Auth Service и Game Service, gRPC health, graceful shutdown и event streaming. |

## gRPC-клиенты, серверы и ошибки

| Файл | Что реализовано |
|---|---|
| `backend/internal/client/grpc/authclient/client.go` | Клиент gateway к Auth Service: регистрация, login, refresh, logout и token validation. |
| `backend/internal/client/grpc/gameclient/client.go` | Клиент gateway к Game Service для игровых query/command и realtime stream. |
| `backend/internal/transport/grpc/authserver/server.go` | gRPC transport Auth Service и преобразование transport/domain данных. |
| `backend/internal/transport/grpc/gameserver/server.go` | gRPC transport Game Service, проверка UUID/времени и подписка на игровые события. |
| `backend/internal/rpc/errors.go` | Стабильное преобразование domain errors в gRPC status/reason и обратно. |

## Доменная модель

| Файл | Что реализовано |
|---|---|
| `backend/internal/model/game.go` | Основные агрегаты Avitosha v2: действия, задания, питомец, комната, история, summary, leaderboard, achievements, character и rewards. |
| `backend/internal/model/domain_event.go` | Типы доменных событий и сохраняемый event envelope. |
| `backend/internal/model/pet.go` | Актуализированная модель питомца для игрового прогресса v2. |

## Игровые use cases и правила

| Файл | Что реализовано |
|---|---|
| `backend/internal/usecase/game_service.go` | Центральная orchestration: транзакционная обработка действий, чтение игровых экранов, публикация событий после commit и работа с наградами. |
| `backend/internal/usecase/game_service_test.go` | Большой набор unit/integration-style тестов прогресса, идемпотентности, concurrency, первой комнаты, summary, leaderboard, rewards и AI advice. |
| `backend/internal/usecase/game_contracts.go` | Интерфейсы repository, транзакций, event publisher и AI generator. |
| `backend/internal/usecase/game_errors.go` | Доменные ошибки игры для стабильного API/error mapping. |
| `backend/internal/usecase/game_rules.go` | Чистые правила соответствия действий заданиям, начисления XP, открытия этапов и расчёта состояния. |
| `backend/internal/usecase/game_rules_test.go` | Табличные тесты игровых правил и граничных условий. |
| `backend/internal/usecase/progression_rules.go` | Пороговые значения уровней и вычисление производного прогресса питомца. |
| `backend/internal/usecase/pet_name.go` | Нормализация, Unicode-aware validation и безопасное переименование питомца. |
| `backend/internal/usecase/pet_name_test.go` | Тесты пустых, длинных, Unicode- и некорректных имён. |
| `backend/internal/usecase/game_advice.go` | Сбор минимального AI-контекста, валидация ответа и локальный fallback. |

## PostgreSQL repository

| Файл | Что реализовано |
|---|---|
| `backend/internal/repository/postgres/game_repository.go` | Транзакционная реализация игровых контрактов и единая точка работы с PostgreSQL. |
| `backend/internal/repository/postgres/game_actions.go` | Идемпотентная запись действия и возврат сохранённого результата при повторном `eventId`. |
| `backend/internal/repository/postgres/game_tasks.go` | Загрузка, назначение и блокировка пользовательских заданий для безопасного обновления. |
| `backend/internal/repository/postgres/game_progress.go` | Атомарное сохранение XP, daily/weekly progress, комнаты, истории, достижений, характера и событий. |
| `backend/internal/repository/postgres/game_rewards.go` | Начисление reward balances, lifetime totals и запись auditable transactions. |
| `backend/internal/repository/postgres/game_repository_test.go` | Repository-тесты SQL-сценариев, rollback, идемпотентности, блокировок и конкурентной обработки. |

## HTTP handlers и DTO

| Файл | Что реализовано |
|---|---|
| `backend/internal/handler/game.go` | REST handlers всех игровых query/command, advice и переименования питомца. |
| `backend/internal/handler/game_dto.go` | Преобразование domain models в стабильные JSON DTO и эффективная сериализация событий. |
| `backend/internal/handler/game_identity.go` | Извлечение и проверка пользовательской identity для игровых запросов. |
| `backend/internal/handler/game_test.go` | Handler-тесты статусов, validation, JSON-контрактов, ошибок и основных игровых сценариев. |
| `backend/internal/handler/game_dto_test.go` | Проверка DTO mapping, event envelope и повреждённых payload. |
| `backend/internal/handler/game_dto_benchmark_test.go` | Воспроизводимый benchmark горячего пути сериализации игровых событий. |
| `backend/internal/handler/game_ws.go` | WebSocket upgrade, проверка origin/identity и доставка пользовательских событий. |
| `backend/internal/handler/game_ws_test.go` | Тесты подключения, доставки событий и удаления отключённого клиента. |
| `backend/internal/handler/middleware.go` | Изменения middleware для игрового и микросервисного transport-контекста. |
| `backend/internal/handler/router.go` | Регистрация REST, advice и WebSocket маршрутов с dependency wiring. |

## Realtime

| Файл | Что реализовано |
|---|---|
| `backend/internal/realtime/hub.go` | Потокобезопасный in-memory hub, per-user subscriptions, ограниченные буферы и отключение медленных клиентов. |
| `backend/internal/realtime/hub_test.go` | Тесты fan-out, фильтрации по пользователю, unsubscribe и backpressure-сценария. |

## AI-интеграция

| Файл | Что реализовано |
|---|---|
| `backend/internal/ai/proxyapi.go` | HTTP-клиент ProxyAPI/OpenRouter, ограничение ответа, timeout, безопасный prompt и обработка ошибок. |
| `backend/internal/ai/proxyapi_test.go` | Тесты успешной генерации, HTTP-ошибок, пустого/некорректного ответа и отмены context. |

## PostgreSQL migrations

| Файл | Что реализовано |
|---|---|
| `backend/migrations/000004_rebuild_avitosha_game.up.sql` | Полная схема Avitosha v2, constraints, indexes и seed первой комнаты, заданий и достижений. |
| `backend/migrations/000004_rebuild_avitosha_game.down.sql` | Контролируемый откат новой игровой схемы. |
| `backend/migrations/game_migration_test.go` | Проверка применения миграции, структуры таблиц, ограничений и seed-данных. |
| `backend/migrations/000005_create_reward_balances.up.sql` | Таблицы reward balances и журнала начислений. |
| `backend/migrations/000005_create_reward_balances.down.sql` | Откат reward-схемы. |
| `backend/migrations/reward_balance_migration_test.go` | Проверка reward migration и обязательных ограничений данных. |

## Служебные сценарии

| Файл | Что реализовано |
|---|---|
| `backend/scripts/init-test-db.sh` | Инициализация отдельной безопасной PostgreSQL-базы для integration-тестов. |
| `backend/scripts/smoke-game.sh` | Сквозная проверка регистрации, действий, задания, XP, комнаты, истории, summary и награды через публичный API. |

## Самые значимые файлы

Если времени на review мало, рекомендую начать с этих файлов:

1. `backend/internal/usecase/game_service.go` — бизнес-сценарий целиком.
2. `backend/internal/repository/postgres/game_repository.go` и
   `game_progress.go` — транзакции и сохранение состояния.
3. `backend/migrations/000004_rebuild_avitosha_game.up.sql` — структура домена.
4. `backend/internal/handler/game.go` и `game_dto.go` — публичный API.
5. `backend/internal/realtime/hub.go` — realtime и backpressure.
6. `backend/internal/ai/proxyapi.go` — безопасная внешняя интеграция.
7. `backend/api/proto/avitosha/v1/services.proto` и
   `backend/internal/app/microservices.go` — микросервисная архитектура.
8. `backend/internal/usecase/game_service_test.go` и
   `backend/internal/repository/postgres/game_repository_test.go` — доказательство
   корректности сложных и конкурентных сценариев.
