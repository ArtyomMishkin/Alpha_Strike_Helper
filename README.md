# 🎮 Alpha Strike Helper — Полная переделка

## ✅ Что было обновлено:

### Доменные модели (clean, единые)
- ✓ `card.go` — чистая GORM модель
- ✓ `user.go` — без дублирования
- ✓ `collection.go` — правильная связь Many-to-Many через `CollectionCard`

### БД слой
- ✓ `postgres.go` — исправленная миграция со всеми таблицами в правильном порядке

### Репозитории (полная реализация)
- ✓ `card_repository.go` — CRUD + Search + Count + Filters
- ✓ `user_repository.go` — CRUD + GetByUsername + GetByEmail + UpdateLastLogin
- ✓ `collection_repository.go` — CRUD + AddCard + RemoveCard + GetCard с Preload

### Сервисы (бизнес-логика)
- ✓ `card_service.go` — ListCards теперь возвращает total для пагинации
- ✓ `auth_service.go` — Register + Login + Token validation с bcrypt
- ✓ `collection_service.go` — полная реализация операций с коллекциями

### JWT утилиты
- ✓ `jwt.go` — GenerateToken + ValidateToken + GetClaims с github.com/golang-jwt/jwt/v5

### Middleware
- ✓ `auth_middleware.go` — правильная проверка Bearer токена

### Хендлеры (полная реализация)
- ✓ `card_handler.go` — List с total_pages, Get, Create, Search, Update, Delete
- ✓ `auth_handler.go` — Register + Login
- ✓ `collection_handler.go` — List, Create, Get, Update, Delete, AddCard, RemoveCard

### Main
- ✓ `main.go` — правильная инициализация всех слоёв, логирование

---

## 🚀 Как использовать:

### 1. Загрузи все файлы в соответствующие директории:

```
internal/domain/       → card.go, user.go, collection.go
internal/repository/   → card_repository.go, user_repository.go, collection_repository.go
internal/service/      → card_service.go, auth_service.go, collection_service.go
internal/middleware/   → auth_middleware.go
internal/handler/      → card_handler.go, auth_handler.go, collection_handler.go
pkg/database/          → postgres.go
pkg/utils/             → jwt.go
cmd/server/            → main.go
```

### 2. Загрузи новый go.mod и обнови зависимости:

```powershell
cd C:\Go\Alpha_Strike_Helper
go mod tidy
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto
go get gorm.io/datatypes
```

### 3. Запусти проект:

```powershell
go run cmd/server/main.go
```

### 4. Проверь логи:

```
🚀 Starting Alpha Strike Helper...
📋 Config loaded: DB=localhost:5432, User=alpha_user
✓ Successfully connected to PostgreSQL database
🔄 Running migrations...
✓ User table migrated
✓ Card table migrated
✓ Collection table migrated
✓ CollectionCard table migrated
✓ All migrations completed successfully
✓ Repositories initialized
✓ Services initialized
✓ Handlers initialized
✓ Server starting on port 8080
📖 API available at http://localhost:8080/api/v1
```

---

## 📝 Тестирование API:

### Регистрация:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'
```

### Вход:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### Получить список карточек:
```bash
curl "http://localhost:8080/api/v1/cards?page=1&page_size=20&type=Medium"
```

### Поиск карточек:
```bash
curl "http://localhost:8080/api/v1/cards/search?q=Enforcer"
```

### Создать коллекцию (нужен токен):
```bash
curl -X POST http://localhost:8080/api/v1/collections \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name":"My Collection","description":"Test collection"}'
```

### Добавить карточку в коллекцию:
```bash
curl -X POST http://localhost:8080/api/v1/collections/1/cards/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"quantity":2}'
```

---

## 🔐 Особенности:

✅ **Безопасность:**
- Пароли хешируются с bcrypt
- JWT токены с истечением через 24 часа
- Auth middleware проверяет все защищённые маршруты

✅ **Функциональность:**
- Полная CRUD для карточек, пользователей, коллекций
- Поиск с ILIKE (case-insensitive)
- Пагинация с total_pages
- Фильтры по type, size, faction, role, tech_base
- Many-to-Many коллекции <-> карточки с количеством

✅ **Архитектура:**
- Clean Architecture (Domain → Repository → Service → Handler)
- Интерфейсы для всех репозиториев (легко тестировать)
- Структурированное логирование
- Обработка ошибок на всех уровнях

---

## 📦 Следующий шаг:

Когда бэкенд запустится и заработают все эндпоинты → переходим на фронтенд!

Нужны обновления файлов на фронте:
- `index.html` — основная страница с фильтрами
- `cards.js` — работа с API (регистрация, вход, список карточек, коллекции)
- `cards.css` — оформление

---

## ⚠️ Если что-то не скомпилируется:

1. Проверь, что все файлы в правильных директориях
2. Выполни `go mod tidy`
3. Проверь версию Go: `go version` (должна быть 1.21 или выше)
4. Убедись, что Postgres запущена и доступна

**Удачи!** 🚀🎉
