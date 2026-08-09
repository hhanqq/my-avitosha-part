# Сергей Дудник — backend-вклад в Avitosha v2

Этот каталог содержит только backend-часть моей работы, извлечённую из
Git-коммитов автора:

```text
sergey <152873255+hhanqq@users.noreply.github.com>
```

Исходная точка: ветка `master`, commit
`d92752556c321b63871de8d036b370c6f81fb180`.

Клиентская часть, полный командный проект, assets, общие корневые файлы и код,
которого нет в моих backend-коммитах, сюда не переносились.

## Что находится в каталоге

```text
my-part-avitosha/
├── README.md
├── source/app/backend/        # полные актуальные версии затронутых файлов
├── patches/                   # точные backend-diff моих коммитов
└── provenance/
    ├── all-author-commits.tsv
    ├── all-historical-backend-paths.txt
    ├── current-backend-files.txt
    ├── current-file-attribution.tsv
    ├── current-line-ranges.tsv
    └── SHA256SUMS
```

`source/` — удобная для чтения версия кода. `patches/` — главный источник для
проверки авторства: патчи содержат только изменения внутри `app/backend`.

## Итог моей работы

Я переработал первоначальную механику питомца в backend Avitosha v2 и построил
вокруг неё целостный игровой контур. Реальные действия пользователя стали
продвигать задания, начислять XP и награды, менять состояние питомца, открывать
предметы комнаты и этапы истории.

В рамках backend я реализовал:

- доменную модель Avitosha v2 и PostgreSQL-схему;
- транзакционный игровой репозиторий;
- атомарную обработку прогресса, XP, наград и доменных событий;
- идемпотентность по `eventId` и защиту конкурентных запросов;
- задания, первую комнату, сюжет, достижения, характер и лидерборд;
- безопасное переименование питомца;
- reward balances и журнал начислений;
- REST API, DTO и OpenAPI-контракт;
- realtime-доставку событий через WebSocket;
- AI-советы через ProxyAPI с безопасным локальным fallback;
- профилирование и оптимизацию сериализации событий;
- разделение backend на API Gateway, Auth Service и Game Service через gRPC;
- unit-, repository-, migration-, handler-, concurrency-, integration- и
  smoke-тесты этих сценариев.

## Проверяемые метрики

| Метрика | Значение |
|---|---:|
| Всего моих коммитов в репозитории | 21 |
| Содержательных не-merge-коммитов | 17 |
| Коммитов с backend-изменениями | 14 |
| Backend diff | +11 029 / −2 311 строк |
| Исторически затронутых backend-путей | 87 |
| Актуальных перенесённых backend-файлов | 64 |
| Текущих строк с моим авторством по `git blame` | 10 889 |
| Из них сгенерированных protobuf-строк | 2 644 |
| Текущих несгенерированных строк | 8 245 |

`+11 029 / −2 311` — сумма Git diff, а не размер итогового проекта. Метрика
учитывает рефакторинг и generated-код. Поэтому generated protobuf stubs отдельно
отмечены и не выдаются за вручную написанные строки.

## Где конкретно мой код

Проверять вклад можно на трёх уровнях.

### 1. Точный diff каждого коммита

Каталог `patches/` содержит 14 backend-only патчей. В них нет файлов из других
частей приложения. Это наиболее точное доказательство того, какие строки были
добавлены, изменены или удалены моими коммитами.

Пример:

```bash
less patches/0001-feat-game-implement-atomic-avitosha-progression.patch
```

### 2. Точные диапазоны в актуальных файлах

`provenance/current-line-ranges.tsv` имеет формат:

```text
file    start_line    end_line    commit
```

Пример поиска по ключевому use case:

```bash
awk -F '\t' '$1 == "app/backend/internal/usecase/game_service.go"' \
  provenance/current-line-ranges.tsv
```

### 3. Сводка по полным файлам

`provenance/current-file-attribution.tsv` показывает для каждого файла:

```text
file    total_lines    sergey_lines    sergey_percent    note
```

Это важно для файлов, которые существовали до моей работы или позже менялись
командой. Наличие файла в `source/` означает, что я менял его в backend-коммите,
но не означает автоматического присвоения мне каждой строки файла.

## Почему некоторые файлы перенесены целиком

В Git-коммите фича представлена отдельными hunks, но один hunk часто невозможно
понять без типов, интерфейсов и окружающей реализации. Поэтому для каждого
актуального затронутого пути в `source/` сохранена полная версия файла на HEAD.

Полный файл нужен, когда:

- use case использует package-level интерфейсы и доменные типы;
- транзакция распределена между repository, rules и service;
- handler должен совпадать с DTO и OpenAPI;
- миграция проверяется совместно с down-файлом и migration-тестом;
- WebSocket transport зависит от hub и router wiring;
- gRPC server должен совпадать с proto-контрактом и клиентом;
- тест использует общие fixtures и helpers своего Go package.

Где в полном файле находится именно мой код, однозначно показывают
`current-line-ranges.tsv` и соответствующий patch. Это позволяет сохранить
контекст, не приписывая себе чужие строки.

## Хронология backend-коммитов

| Commit | Область | Реализованный результат |
|---|---|---|
| `e58b114` | Domain и DB | Новая модель игры, доменные события, миграция `000004` и её тест |
| `738d79b` | PostgreSQL | Транзакционный репозиторий действий, заданий и прогресса |
| `9feea2b` | Use case | Атомарный игровой прогресс и детерминированные правила |
| `bac31e8` | API | Игровые endpoint'ы v2, DTO, OpenAPI и handler-тесты |
| `409b902` | Realtime | WebSocket hub и публикация событий после commit |
| `6098a6a` | Integration tests | Сквозной сценарий прохождения первой комнаты |
| `03c85d4` | Game features | Daily summary, leaderboard и характер питомца |
| `f7289e5` | Runtime | Backend smoke-сценарий и воспроизводимая конфигурация запуска |
| `1b83af6` | Reliability | Усиленные тесты идемпотентности и конкурентного выполнения |
| `65d4a25` | Pet profile | Валидация и безопасное переименование питомца |
| `d23d281` | Rewards | Балансы, ledger начислений, миграция и repository-тесты |
| `f9648cc` | AI | ProxyAPI-клиент, ограниченный контекст, валидация и fallback |
| `66d47d2` | Performance | Benchmark и оптимизация сериализации событий |
| `187be5d` | Architecture | gRPC-контракт и разделение на gateway/auth/game services |

Полная история всех 21 коммита, включая merge-коммиты, находится в
`provenance/all-author-commits.tsv`.

Три содержательных коммита не породили backend-патч в этом пакете:

- `b995e44` — архитектурный документ вне `app/backend`;
- `56681ac` — работа вне backend;
- `d927525` — удаление личных рабочих заметок.

Четыре merge-коммита также не считаются самостоятельным backend diff, потому
что могут объединять работу нескольких авторов.

## Карта ключевого кода

### Домен и PostgreSQL

- `source/app/backend/internal/model/game.go` — питомец, задания, комната,
  история, достижения, характер и leaderboard read models.
- `source/app/backend/internal/model/domain_event.go` — сохраняемые игровые
  события.
- `source/app/backend/migrations/000004_rebuild_avitosha_game.up.sql` — схема
  Avitosha v2 и seed первой комнаты.
- `source/app/backend/internal/repository/postgres/game_repository.go` —
  транзакционная точка входа.
- `game_actions.go` — идемпотентная регистрация действия.
- `game_progress.go` — обновление XP, комнаты, истории и агрегатов.
- `game_tasks.go` — выбор и блокировка пользовательских заданий.
- `game_rewards.go` — начисление reward balances и запись ledger.

### Игровой use case

- `source/app/backend/internal/usecase/game_service.go` — orchestration полного
  игрового сценария.
- `game_rules.go` — правила прогресса задач и наград.
- `progression_rules.go` — уровни и производные состояния питомца.
- `game_contracts.go` — границы repository и транзакций.

`ProcessAction` проверяет `eventId`, открывает транзакцию, блокирует нужный
прогресс, рассчитывает результат и сохраняет действие, награды и события.
Повторная обработка того же события возвращает прежний результат, не начисляя
XP или награду повторно.

### REST и WebSocket

- `source/app/backend/api/openapi.yaml` — публичный API-контракт.
- `source/app/backend/internal/handler/game.go` — HTTP handlers.
- `source/app/backend/internal/handler/game_dto.go` — transport DTO и
  сериализация событий.
- `source/app/backend/internal/realtime/hub.go` — неблокирующий realtime hub.
- `source/app/backend/internal/handler/game_ws.go` — WebSocket transport.

События публикуются только после успешного commit. Клиентские очереди
ограничены; медленный подписчик отключается и не блокирует основной use case.

### AI-советы

- `source/app/backend/internal/ai/proxyapi.go` — клиент ProxyAPI.
- `source/app/backend/internal/usecase/game_advice.go` — формирование
  разрешённого контекста, проверка ответа и fallback.
- `proxyapi_test.go` и `game_service_test.go` — успешный ответ, ошибки,
  некорректные данные и fallback-сценарии.

Во внешний AI-запрос не передаются email, access token и пользовательские
сообщения. Модель генерирует только текст совета и не участвует в расчёте XP,
прогресса или баланса.

### gRPC-микросервисы

- `source/app/backend/api/proto/avitosha/v1/services.proto` — versioned RPC API.
- `source/app/backend/internal/app/microservices.go` — сборка трёх процессов.
- `source/app/backend/internal/client/grpc/` — клиенты API Gateway.
- `source/app/backend/internal/transport/grpc/` — auth/game servers.
- `source/app/backend/internal/rpc/errors.go` — отображение доменных ошибок.
- `source/app/backend/cmd/api`, `cmd/auth`, `cmd/game` — точки запуска.

API Gateway сохраняет публичные HTTP/WebSocket-контракты и не обращается к SQL.
Auth Service отвечает за пользователей и сессии, Game Service — за игровой
домен. Realtime-события передаются в gateway через server-streaming RPC.

### Производительность

Горячий путь сериализации событий раньше выполнял лишнее преобразование
`json.RawMessage → map[string]any → JSON`. После замены на проверяемое
объединение payload benchmark показал:

| Метрика | До | После | Изменение |
|---|---:|---:|---:|
| Время | 26 428 ns/op | 18 480 ns/op | −30,1% |
| Память | 11 798 B/op | 6 149 B/op | −47,9% |
| Аллокации | 272 | 44 | −83,8% |

Benchmark и тесты находятся в
`source/app/backend/internal/handler/game_dto_benchmark_test.go` и
`game_dto_test.go`.

## Тестовое покрытие моей части

Перенесённые тесты проверяют:

- уровни, XP, настроение и игровые правила;
- полный путь первой комнаты;
- идемпотентную повторную обработку `eventId`;
- параллельные одинаковые и независимые действия;
- отсутствие потерянных обновлений;
- repository и PostgreSQL migrations;
- API status, DTO и OpenAPI contract;
- WebSocket fan-out и очистку отключённых клиентов;
- публикацию событий только после commit;
- AI success/error/timeout/invalid-response/fallback;
- gRPC validation, error mapping и service wiring.

## Как проверить происхождение

Контрольные суммы:

```bash
shasum -a 256 -c provenance/SHA256SUMS
```

Список backend-файлов:

```bash
cat provenance/current-backend-files.txt
```

Сводка авторства по файлам:

```bash
column -t -s $'\t' provenance/current-file-attribution.tsv | less -S
```

Независимая проверка в исходном Git-репозитории:

```bash
git show <commit> -- app/backend
git blame --line-porcelain HEAD -- app/backend/<path>
```

## Границы пакета

- Это extract моего backend-вклада, а не самостоятельная копия всего проекта.
- Невключённые общие зависимости могут быть нужны для отдельной сборки.
- Полные файлы сохранены ради читаемости; точное авторство подтверждают патчи.
- Сгенерированные protobuf-файлы включены для понимания контракта, но отдельно
  отмечены как generated-code.
- Merge-коммиты не используются для присвоения кода одному автору.

Такой формат даёт проверяющей системе одновременно цельный backend-контекст и
однозначные Git-доказательства каждой заявленной части работы.
