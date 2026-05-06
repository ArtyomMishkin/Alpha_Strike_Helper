# Описание выполненных работ по приложению Alpha Strike Helper

## 1. Цель приложения

Разработано веб-приложение для управления карточками юнитов BattleTech: Alpha Strike, формирования отрядов (Lance/Star), расчета суммарных параметров и подготовки данных к экспорту/печати.

Основная задача приложения: упростить подбор юнитов, фильтрацию по характеристикам и сборку игровых формаций с проверкой ограничений.

## 2. Реализованная архитектура

Приложение разделено на несколько уровней:

- `internal/domain` — доменные модели (`Card`, `User`, `Collection`, `Lance`, `Star`).
- `internal/repository` — доступ к PostgreSQL через GORM.
- `internal/service` — бизнес-логика (авторизация, карточки, коллекции, расчет статистики Lance/Star).
- `internal/handler` — HTTP-обработчики API на Gin.
- `internal/middleware` — middleware для CORS, JWT и логирования.
- `internal/sync` — модуль синхронизации внешних источников (в т.ч. Master Unit List: клиент, маппинг, импорт).
- `pkg` — инфраструктура (конфиг, подключение к БД, JWT-утилиты).
- `templates` / `static` — frontend часть на HTML/CSS/JavaScript.
- `cmd/*` — исполняемые команды:
  - `cmd/server` (web/API сервер),
  - `cmd/masterunitlist_sync` (ручной/разовый импорт),
  - `cmd/weekly_sync` (плановая периодическая синхронизация),
  - `cmd/sync_service` (микросервис sync с HTTP API),
  - `cmd/cards_service` (микросервис cards с read/admin API и chassis-sources).

### 2.1 Схема сервисов и портов

```mermaid
flowchart LR
    Browser["Browser UI (layout.html)"]
    AppServer["app-service :8080"]
    CardsService["cards-service :8082"]
    SyncService["sync-service :8081"]
    Postgres["PostgreSQL :5432"]
    MulApi["MasterUnitList API"]

    Browser -->|"GET /"| AppServer
    Browser -->|"GET cards + chassis sources"| CardsService
    Browser -->|"sync control (optional)"| SyncService

    AppServer -->|"read/write domain data"| Postgres
    CardsService -->|"read/admin cards"| Postgres
    SyncService -->|"import upsert cards"| Postgres
    SyncService -->|"QuickList/Faction API"| MulApi
```

## 3. Backend и API

### 3.1 Авторизация и пользователи

Реализована backend-подсистема аутентификации (JWT):

- регистрация пользователя;
- вход пользователя;
- генерация JWT-токена;
- защита части API через JWT middleware (`/api/v1/collections`, `/api/v1/admin/cards`).

Примечание:

- ключевые пользовательские сценарии текущего UI (каталог карточек, ростер, формации Lance/Star) работают по публичным маршрутам и не требуют токена.

### 3.2 Карточки юнитов

Реализованы:

- получение списка карточек с пагинацией;
- фильтрация по параметрам;
- поиск по имени/модели;
- получение карточки по ID;
- админские операции CRUD (создание/редактирование/удаление).

### 3.3 Формации Lance/Star

Реализованы:

- создание и редактирование Lance/Star;
- управление составом формаций;
- базовая валидация комплектности;
- расчет агрегированных характеристик;
- экспорт данных формаций.

### 3.4 Коллекции пользователя

Реализованы:

- создание пользовательских коллекций;
- добавление/удаление карточек в коллекциях;
- просмотр коллекций пользователя.

## 4. База данных и модели

Используется PostgreSQL + GORM AutoMigrate.

Основные таблицы:

- `users`
- `cards`
- `collections`
- `lances`
- `stars`
- таблицы участников формаций (`lance_members`, `star_members`)

Для карточек реализовано хранение как базовых характеристик Alpha Strike, так и расширенной информации по доступности:

- `available_factions` (JSONB массив),
- `available_eras` (JSONB массив),
- `faction_era_availability` (JSONB объект, соответствие фракций и эпох).

## 5. Импорт данных

### 5.1 Единоразовый импорт из Master Unit List

Добавлен отдельный CLI-импорт:

- `cmd/masterunitlist_sync`
- модуль: `internal/sync/masterunitlist`

Что делает:

- загружает карточки через `Unit/QuickList`;
- собирает доступность по эпохам (`AvailableEras`);
- собирает доступность по фракциям (`Factions`);
- опционально строит детализацию `фракция -> список эпох`;
- выполняет upsert/replace в таблицу `cards`.

### 5.2 Плановая синхронизация (еженедельная)

Добавлен отдельный CLI-процесс периодической синхронизации:

- `cmd/weekly_sync`

Что делает:

- запускает импорт по расписанию (`time.Ticker`);
- может выполнить синхронизацию сразу при старте (`--run-now=true`);
- использует тот же импортный модуль `internal/sync/masterunitlist`.

### 5.3 Sync Service (отдельный HTTP-процесс)

Добавлен отдельный микросервис синхронизации:

- `cmd/sync_service`
- orchestration-слой: `internal/sync/service`

HTTP-эндпоинты:

- `GET /health` — проверка состояния процесса;
- `GET /sync/status` — статус последнего запуска;
- `POST /sync/run` — ручной запуск импорта (асинхронно).

### 5.4 Cards Service (отдельный HTTP-процесс)

Добавлен отдельный микросервис карточек:

- `cmd/cards_service`

HTTP-эндпоинты:

- `GET /health`
- `GET /api/v1/cards`
- `GET /api/v1/cards/:id`
- `GET /api/v1/cards/search`
- `POST /api/v1/admin/cards`
- `PUT /api/v1/admin/cards/:id`
- `DELETE /api/v1/admin/cards/:id`
- `GET /api/v1/chassis-sources`

Примечание:

- frontend (`templates/layout.html`) получает карточки и `chassis_sources` напрямую из cards-service (`:8082`).

## 6. Обновления фильтрации API карточек

Добавлена поддержка фильтрации по новой структуре доступности:

- `faction` — работает по legacy-полю `faction` и по `available_factions`;
- `era` — работает по legacy-полю `era` и по `available_eras`;
- `available_faction` — фильтр только по `available_factions`;
- `available_era` — фильтр только по `available_eras`.

Это обеспечивает обратную совместимость и поддержку новых данных одновременно.

## 7. Frontend

Реализован интерфейс на Vanilla JS:

- каталог карточек;
- формации Lance/Star;
- модальные окна для CRUD операций;
- базовая интеграция с backend API;
- подготовка данных к печати/экспорту.

## 8. Инфраструктура и запуск

Реализованы:

- конфигурация через переменные окружения;
- Dockerfile для сборки и запуска приложения;
- подключение к PostgreSQL;
- автоматическая миграция схемы при старте.

## 9. Текущее состояние проекта

Проект в состоянии рабочего прототипа/предрелиза:

- backend и frontend работают в связке, ключевые пользовательские сценарии доступны;
- каталог карточек загружается полностью (с пагинацией), фильтры по фракции/группе/эпохе/типу работают;
- ростер, формации, ангар, печать/экспорт и блок `Где взять` реализованы и используются в UI;
- импорт из MUL работает как в ручном режиме (`masterunitlist_sync`), так и в периодическом (`weekly_sync`);
- данные и UI уже покрывают BattleMech + другие типы юнитов, включая правила по совместимости и валидации формаций.

Ограничения текущего этапа:

- автотесты и регрессионное покрытие пока неполные;
- часть бизнес-правил и справочников (источники миниатюр/алиасы) продолжает пополняться вручную;
- перед финальным релизом нужна дополнительная вычитка документации и smoke-проверка сценариев после синхронизации.

## 10. Рекомендуемые следующие шаги

1. Добавить автотесты для `repository/service` слоев и импортеров.
2. Описать API (включая новые query-параметры) в отдельной документации.
3. Подготовить полноценный `docker-compose.yml` (app + postgres).
4. Добавить сидирование/команды smoke-проверки после импорта.
5. Провести оптимизацию тяжелых запросов и индексацию (по необходимости).

## 11. Практическая эксплуатация (текущее состояние)

### 11.1 Быстрый запуск (Docker Compose)

Требования:

- Docker
- Docker Compose

Команды:

- `docker-compose -f docker/docker-compose.yml up -d --build`
- `docker-compose -f docker/docker-compose.yml logs -f app`
- `docker-compose -f docker/docker-compose.yml down`

После запуска:

- API: `http://localhost:8080`
- Sync Service API: `http://localhost:8081`
- Cards Service API: `http://localhost:8082`
- PostgreSQL: `localhost:5432`

### 11.2 Запуск без Docker

Необходим PostgreSQL и переменные окружения:

- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `SERVER_PORT`
- `JWT_SECRET`

Команда запуска сервера:

- `go run ./cmd/server`

### 11.2.1 Запуск sync-сервиса без Docker

Команда запуска:

- `go run ./cmd/sync_service`

Доступные HTTP-эндпоинты:

- `GET /health`
- `GET /sync/status`
- `POST /sync/run`

Пример body для `POST /sync/run`:

```json
{
  "unit_type_ids": [18, 19, 17, 21],
  "replace_first": true,
  "include_faction_eras": true,
  "http_timeout_seconds": 180
}
```

### 11.2.2 Запуск cards-сервиса без Docker

Команда запуска:

- `go run ./cmd/cards_service`

### 11.3 Импорт карточек MUL (Master Unit List)

Основная команда:

- `go run ./cmd/masterunitlist_sync --unit-type-id=18 --replace=true --include-faction-eras=true --era-ids=13,247,14`

Ключевые параметры:

- `--unit-type-id=18` (BattleMech)
- `--replace=true|false` (пересоздавать набор карточек)
- `--include-faction-eras=true|false`
- `--era-ids=13,247,14` (ограничение эпох)
- `--faction-ids=24,29,27` (ограничение фракций)
- `--http-timeout=120s` (рекомендуется 180s+ при нестабильной сети)
- `--batch-size=300`

Примечание:

- если внешний сервис нестабилен на широких era-запросах, используйте режим эпоха -> фракция (в текущей реализации это основная стратегия заполнения доступности).

### 11.4 Новые поля карточки по доступности

В модели карточки используются:

- `available_factions` (`jsonb`, массив строк)
- `available_faction_groups` (`jsonb`, массив строк)
- `available_eras` (`jsonb`, массив строк)
- `faction_era_availability` (`jsonb` object)

### 11.5 Новые фильтры API `/api/v1/cards`

Поддерживаются query-параметры:

- `faction`
- `available_faction`
- `faction_group`
- `era`
- `available_era`
- `unit_type` (рекомендуемый)
- `type` (legacy-синоним)
- `role`
- `size`
- `techbase`
- `name`
- `pvmin`
- `pvmax`

Примеры:

- `/api/v1/cards?faction_group=HW Clan&available_era=Jihad`
- `/api/v1/cards?available_faction=Clan Wolf&era=Clan Invasion`

### 11.6 Защита от дублей

При импорте карточек:

- выполняется проверка дублей по `model_number` до и после импорта;
- при обнаружении дублей лишние строки удаляются;
- запись карточек идет через upsert по `model_number`.

## 12. Последние доработки (UI, логика ростера, импорт)

### 12.1 Каталог источников миниатюр (Sarna) и блок "Где взять"

Реализован полный цикл наполнения `static/data/chassis_sources.json` на основе списка Sarna (Miniatures - Catalyst Game Labs):

- собран и очищен каталог `шасси -> [коробки/паки]`;
- удалены технические артефакты парсинга (`[edit]`, заголовки секций, мусорные ключи);
- добавлены алиасы шасси (например, варианты `Timber Wolf/Mad Cat`);
- добавлена поддержка составных имен формата `Alias (Canonical)` (например, `Masakari (Warhawk)`, `Daishi (Dire Wolf)`);
- для `Elemental` поддержаны разные варианты написания (исключая `Water Elemental`);
- нормализованы названия источников;
- на карточках в разделе "Список мехов" отображается блок `Где взять: ...`.

### 12.2 Полная загрузка карточек на frontend

Исправлена загрузка карточек в UI:

- ранее использовалась только первая страница API (`page=1`);
- теперь `loadCards()` загружает **все страницы** (`page..total_pages`);
- это устранило ситуацию, когда в фильтрах отображались не все эпохи/данные.

### 12.3 Расширенная совместимость фракций через General

Добавлена логика совместимости `General`-фракций в пределах группы:

- `IS Clan General` считается доступной для всех фракций группы `IS Clan`;
- аналогично для `HW Clan General`, `Inner Sphere General`, `Periphery General` и т.д.;
- фильтрация, добавление в ростер и подсказки используют единую совместимость;
- выбор такого юнита не должен принудительно менять выбранную пользователем фракцию ростера.

### 12.4 Доступность юнита по фракциям и эпохам

В карточках списка добавлена кнопка `Доступность`:

- открывается модальное окно с матрицей `фракция × эпоха`;
- используется `faction_era_availability` и поля доступности карточки;
- в таблице применяется хронологическая сортировка эпох (не алфавит).

### 12.5 Эпохи и синхронизация фильтров "Список мехов" -> "Ростер"

Реализован перенос выбранной эпохи:

- эпоха из фильтра списка юнитов автоматически подставляется в `Ростер`;
- ручной выбор эпохи в `Ростере` имеет приоритет и не перезаписывается;
- при очистке ростера флаги ручного выбора сбрасываются.

Также эпохи в фильтрах и связанных селекторах отсортированы по хронологии:

- Star League
- Early Succession War
- Late Succession War - LosTech
- Late Succession War - Renaissance
- Clan Invasion
- Civil War
- Jihad
- Early Republic
- Late Republic
- Dark Age
- ilClan

### 12.6 Исключение юнитов с PV=0

На клиентской стороне добавлено исключение юнитов с `PV <= 0`:

- в каталоге карточек;
- в фильтрации списка;
- в поиске добавления в ростер.

### 12.7 Рефактор и доработка формаций Clan/IS

Обновлена логика формирования и выбора типов:

- для Clan добавлены Star-аналоги Lance-типов (`Battle Lance -> Battle Star` и т.д.);
- размер Star-формаций зафиксирован как 5;
- `Omni Star` удален из активных типов (и мигрируется в `Battle Star` для старых данных UI);
- исправлено определение стороны ростера с приоритетом выбранной фракции (например, `Clan Sea Fox`);
- ограничено переполнение формации при drag&drop (проверка капа по размеру).

### 12.8 Skill юнитов (формации + печать)

Добавлена система `Skill` для юнитов:

- базовый skill = `4`;
- диапазон изменения skill: `0..8`;
- доступно изменение skill в блоке участника формации;
- PV юнита пересчитывается на основе базового PV и таблиц повышения/понижения skill;
- пересчитанный PV используется в:
  - карточке участника формации,
  - статистике ростера,
  - экспорте/preview,
  - печати карточек (`print`), где также выводится текущий `Skill`.

### 12.9 Ангар и рекомендации для неполных формаций

Доработан раздел "Ангар" и интеграция с "Формациями":

- ангар хранит owned-количество и показывает:
  - `Имеется: xN`,
  - `Используется: xM` (сколько взято в ростер через рекомендации);
- owned-количество не уменьшается при добавлении рекомендации;
- в неполной формации отображается блок:
  - `Сюда может подойти из ангара`;
- рекомендации учитывают:
  - выбранную фракцию (строгая совместимость),
  - выбранную эпоху (если задана),
  - доступность в ангаре (`имеется - используется`),
  - валидность/полезность для условий текущей формации;
- расширена выдача рекомендаций: больше вариантов + приоритет разнообразия по шасси.

### 12.10 Локализация и UX-полировка

Внесены дополнительные UI-правки:

- удалена кнопка `Авто-формирование` из панели выбора формаций;
- переведены элементы в "Список мехов":
  - `Add to Roster` -> `В ростер`,
  - `Min PV` -> `Мин. PV`,
  - `Max PV` -> `Макс. PV`;
- исправлен баг кнопки удаления из ростера (сравнение ID как строк/чисел).

### 12.11 Полный импорт по всем эпохам

Проведен и отлажен импорт по всем доступным эпохам и типам юнитов:

- использован сценарий `masterunitlist_sync` с полным `era-ids`;
- исправлен сценарий, при котором поэтапный импорт мог затирать агрегированные `available_eras`;
- итоговая стратегия: импорт по каждому типу юнита сразу по всему списку эпох.

### 12.12 Docker: команды для работы с проектом

Базовые команды (выполняются из корня репозитория):

- первый запуск/пересборка всех сервисов:
  - `docker compose -f docker/docker-compose.yml up --build -d`
- обычный запуск без пересборки:
  - `docker compose -f docker/docker-compose.yml up -d`
- запуск только одного сервиса с пересборкой:
  - `docker compose -f docker/docker-compose.yml up --build -d sync-service`
  - `docker compose -f docker/docker-compose.yml up --build -d cards-service`
  - `docker compose -f docker/docker-compose.yml up --build -d app`

Проверка состояния:

- список контейнеров и статус:
  - `docker compose -f docker/docker-compose.yml ps`
- просмотр логов всех сервисов:
  - `docker compose -f docker/docker-compose.yml logs -f`
- просмотр логов конкретного сервиса:
  - `docker compose -f docker/docker-compose.yml logs -f sync-service`
  - `docker compose -f docker/docker-compose.yml logs -f cards-service`
  - `docker compose -f docker/docker-compose.yml logs -f app`

Остановка и очистка:

- остановить сервисы:
  - `docker compose -f docker/docker-compose.yml down`
- остановить и удалить volume БД (полный сброс данных):
  - `docker compose -f docker/docker-compose.yml down -v`

Полезные проверки API после запуска:

- health main app: `http://localhost:8080/health`
- health sync-service: `http://localhost:8081/health`
- health cards-service: `http://localhost:8082/health`
- статус/прогресс импорта:
  - `http://localhost:8081/sync/status`
  - `http://localhost:8081/sync/progress`

Пример точечного импорта (пехота + aerospace, все эпохи, все фракции):

- PowerShell:
  - ```powershell
    $body = @{
      unit_type_ids = @(21,17)           # 21 = Infantry, 17 = Aerospace
      replace_first = $false             # не очищать текущую базу
      include_faction_eras = $true
      http_timeout_seconds = 180
      batch_size = 300
    } | ConvertTo-Json

    irm -Method Post -Uri http://localhost:8081/sync/run -ContentType "application/json" -Body $body
    ```
- Проверка прогресса:
  - `irm http://localhost:8081/sync/progress | ConvertTo-Json -Depth 10`

### 12.13 Admin CRUD карточек (cards-service)

Админские операции по карточкам выполняются через `cards-service`:

- `POST http://localhost:8082/api/v1/admin/cards` — создать карточку
- `PUT http://localhost:8082/api/v1/admin/cards/:id` — обновить карточку
- `DELETE http://localhost:8082/api/v1/admin/cards/:id` — удалить карточку

#### Как получить токен для admin-эндпоинтов

В текущей реализации `cards-service` использует JWT-проверку, а сам логин выполняется через основной сервис (`app`, порт `8080`):

- `POST http://localhost:8080/api/v1/auth/login`

PowerShell пример:

```powershell
$loginBody = @{
  username = "admin"
  password = "admin123"
} | ConvertTo-Json

$loginResp = irm -Method Post -Uri "http://localhost:8080/api/v1/auth/login" -ContentType "application/json" -Body $loginBody
$token = $loginResp.token
```

#### Пример создания карточки (Create)

```powershell
$body = @{
  name = "Viper A"
  model_number = "VPR-A-EXAMPLE"
  unit_type = "BattleMech"
  type = "Striker"
  size = 1
  move = "16\"j"
  tmm = 4
  point_value = 34
  armor = 4
  structure = 2
  damage_short = "3"
  damage_medium = "3"
  damage_long = "0"
  overheat = 0
  abilities = "JMPS1"
  tech_base = "Clan"
  role = "Striker"
  source = "Manual"
  faction = "Clan Wolf"
  era = "Clan Invasion"
  available_factions = @("Clan Wolf")
  available_faction_groups = @("HW Clan")
  available_eras = @("Clan Invasion")
  faction_era_availability = @{
    "Clan Wolf" = @("Clan Invasion")
  }
} | ConvertTo-Json -Depth 10

irm -Method Post `
  -Uri "http://localhost:8082/api/v1/admin/cards" `
  -Headers @{ Authorization = "Bearer $token" } `
  -ContentType "application/json" `
  -Body $body
```

#### Пример обновления карточки (Update)

```powershell
$cardId = 123
$updateBody = @{
  name = "Viper A (Updated)"
  point_value = 35
  overheat = 1
} | ConvertTo-Json -Depth 10

irm -Method Put `
  -Uri "http://localhost:8082/api/v1/admin/cards/$cardId" `
  -Headers @{ Authorization = "Bearer $token" } `
  -ContentType "application/json" `
  -Body $updateBody
```

Примечание: текущий `Update` использует `Save` на полной модели, поэтому на практике рекомендуется отправлять полный объект карточки (не только частичный патч), чтобы не перезаписать поля нулевыми значениями.

#### Пример удаления карточки (Delete)

```powershell
$cardId = 123
irm -Method Delete `
  -Uri "http://localhost:8082/api/v1/admin/cards/$cardId" `
  -Headers @{ Authorization = "Bearer $token" }
```

#### Как сейчас определяется "админ"

На текущем этапе отдельной роли администратора (`is_admin`, `role=admin`) в модели пользователя нет.

- `cards-service` проверяет только валидность JWT (`Authorization: Bearer <token>`);
- значит фактически "админ" сейчас = любой аутентифицированный пользователь с корректным токеном.

Для строгого разграничения прав нужен отдельный флаг/роль в `users` и проверка этой роли в middleware.

