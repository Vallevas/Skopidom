# Skopidom

Приложение для автоматизированного учета, регистрации и журналирования инвентаря

---

Skopidom — это полнофункциональная система управления инвентарём с возможностями аудита на основе блокчейна. Система предоставляет:

- **Управление инвентарём**: Отслеживание предметов по штрих-кодам, инвентарным номерам, категориям и местоположениям
- **Управление пользователями**: Ролевая модель доступа (администратор/пользователь)
- **Аудит операций**: Неизменяемый журнал аудита всех операций с инвентарём на базе блокчейна
- **Импорт CSV**: Массовый импорт данных для миграции существующего инвентаря
- **Современный UI**: Фронтенд на React с TypeScript, TailwindCSS и поддержкой интернационализации

### Технологический стек

**Бэкенд:**

- Go 1.24.0
- Chi v5 (HTTP-маршрутизатор)
- pgx/v5 (драйвер PostgreSQL)
- go-ethereum (интеграция с блокчейном)
- golang-migrate (миграции базы данных)
- JWT-аутентификация

**Фронтенд:**

- React 18.3.1
- TypeScript 5.5.3
- Vite 5.4.6
- TailwindCSS 3.4.11
- Компоненты Radix UI
- React Query 5.56.2
- React Router 6.26.2
- i18next (интернационализация)

**Инфраструктура:**

- PostgreSQL 16-alpine
- Docker и Docker Compose
- Ganache (блокчейн Ethereum для аудита)

---

### Предварительные требования

- Docker и Docker Compose
- Go 1.24+ (для локальной разработки)
- Node.js 18+ (для разработки фронтенда)

---

## Быстрый старт

### 1. Клонирование репозитория

```bash
git clone https://github.com/Vallevas/Skopidom.git
cd skopidom
```

### 2. Настройка переменных окружения

Скопируйте примеры файлов окружения:

```bash
cp .env.example .env
cp backend/.env.example backend/.env
```

Отредактируйте `.env` и `backend/.env` при необходимости.

### 3. Запуск

```bash
make dbup
make back
<Ctrl-C>
make createadmin
make seed FILE=items_100.csv USER=admin@university.ru

# Start Ganache
make ganache

# Start backend
make run

# Start frontend
cd frontend
npm run dev -- --host 0.0.0.0
```

Учётные данные администратора по умолчанию:

- Email: `admin@university.ru`
- Пароль: `password`

Доступ к приложению: <https://localhost:5173>

## 🛠️ Разработка

### Разработка бэкенда

```bash
# Запуск бэкенда локально
make back

# Сборка бинарного файла
make build

# Генерация кода SQLC
make generate

# Запуск линтера
make lint

# Статические проверки
make check
```

### Разработка фронтенда

```bash
cd frontend

# Установка зависимостей
npm install

# Запуск сервера разработки
npm run dev

# Сборка для продакшена
npm run build

# Проверка типов
npm run type-check

# Линтинг
npm run lint
```

### Операции с базой данных

```bash
# Открыть оболочку psql
make psql

# Запуск только PostgreSQL
make dbup

# Остановка PostgreSQL
make dbdown

# Полный сброс (остановка + очистка томов)
make reset
```

## 📦 Импорт CSV

Skopidom поддерживает массовый импорт предметов инвентаря из CSV-файлов.

### Формат CSV

```csv
barcode,inventory_number,name,category,building,room,description
12345,INV-001,Стол письменный,Мебель,Главный корпус,Аудитория 101,Деревянный стол
```

### Использование

```bash
# Проверка CSV без импорта
make seed-dry-run FILE=items.csv

# Импорт предметов
make seed FILE=items.csv USER=admin@university.ru
```

Подробная документация в [CSV_IMPORT.md](backend/CSV_IMPORT.md).

## 🔗 Интеграция с блокчейном

Приложение использует смарт-контракты Ethereum для неизменяемого журнала аудита.

### Генерация привязок контракта

```bash
make generate-contract
```

Эта команда генерирует Go-привязки из ABI и байт-кода Solidity-контракта.

## 📁 Структура проекта

```
skopidom/
├── backend/
│   ├── cmd/
│   │   ├── server/      # Точка входа приложения
│   │   └── seed/        # Утилита импорта CSV
│   ├── internal/
│   │   ├── domain/      # Бизнес-логика и сущности
│   │   ├── service/     # Сервисы приложения
│   │   ├── handler/     # HTTP-обработчики
│   │   └── infrastructure/
│   │       ├── postgres/    # Слой работы с БД
│   │       └── blockchain/  # Интеграция со смарт-контрактами
│   ├── pkg/             # Публичные пакеты
│   ├── contracts/       # Миграции базы данных
│   ├── blockchain/      # Смарт-контракты (Solidity)
│   └── uploads/         # Хранилище файлов
├── frontend/
│   ├── src/
│   │   ├── components/  # React-компоненты
│   │   ├── pages/       # Компоненты страниц
│   │   ├── hooks/       # Кастомные хуки
│   │   ├── services/    # API-клиенты
│   │   └── locales/     # Переводы i18n
│   └── public/
├── tests/               # Скрипты интеграционных тестов
├── docker-compose.yml
├── Makefile
└── README.md
```

## 🔐 Аутентификация

Приложение использует JWT-аутентификацию:

- Токены действительны 24 часа (настраивается через `JWT_TTL`)
- Пароли хешируются с помощью bcrypt
- Ролевой контроль доступа (роли администратор/пользователь)

## ⚙️ Конфигурация

### Переменные окружения

**Корневой `.env`:**

| Переменная        | Описание        | По умолчанию |
| ----------------- | --------------- | ------------ |
| POSTGRES_DB       | Имя базы данных | skopidom     |
| POSTGRES_USER     | Пользователь БД | skopidom     |
| POSTGRES_PASSWORD | Пароль БД       | secret       |

**Бэкенд `.env`:**

| Переменная      | Описание                             | По умолчанию |
| --------------- | ------------------------------------ | ------------ |
| DATABASE_URL    | Строка подключения PostgreSQL        | -            |
| JWT_SECRET      | Секрет подписи JWT (мин. 32 символа) | -            |
| JWT_TTL         | Время жизни токена                   | 24h          |
| SERVER_PORT     | Порт сервера бэкенда                 | 8080         |
| DEBUG           | Режим отладки                        | True         |
| ALLOWED_ORIGINS | Разрешённые CORS origin              | \*           |
| STORAGE_DIR     | Директория для загружаемых файлов    | ./uploads    |

## 🐳 Docker-команды

```bash
# Запуск всех сервисов
make up

# Остановка всех сервисов
make down

# Просмотр логов
make logs

# Логи только бэкенда
make logs-backend

# Перезапуск конкретного сервиса
docker compose restart backend
```

## 📄 Лицензия

Подробности в файле [LICENSE](LICENSE).
