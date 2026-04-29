# Создать предмет
curl -s -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "barcode": "INV-2024-001",
    "name": "Dell OptiPlex 7090",
    "category_id": 1,
    "room_id": 1,
    "description": "Системный блок, SN: XYZ123"
  }' | jq .

# Создать второй предмет
curl -s -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "barcode": "INV-2024-002",
    "name": "Dell Monitor U2422H",
    "category_id": 1,
    "room_id": 1,
    "description": "Монитор 24 дюйма"
  }' | jq .

# Получить все предметы
curl -s http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $TOKEN" | jq .

# Получить по ID
curl -s http://localhost:8080/api/v1/items/1 \
  -H "Authorization: Bearer $TOKEN" | jq .

# Поиск по штрих-коду (сценарий сканирования)
curl -s http://localhost:8080/api/v1/items/barcode/INV-2024-001 \
  -H "Authorization: Bearer $TOKEN" | jq .

# Фильтр по комнате
curl -s "http://localhost:8080/api/v1/items?room_id=1" \
  -H "Authorization: Bearer $TOKEN" | jq .

# Фильтр по категории
curl -s "http://localhost:8080/api/v1/items?category_id=1" \
  -H "Authorization: Bearer $TOKEN" | jq .

# Фильтр по статусу
curl -s "http://localhost:8080/api/v1/items?status=active" \
  -H "Authorization: Bearer $TOKEN" | jq .

# Обновить описание (только description и photo_url доступны)
curl -s -X PATCH http://localhost:8080/api/v1/items/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"description":"Системный блок, SN: XYZ123, RAM: 16GB, обновлено"}' | jq .

# Попытка создать с дублирующим штрих-кодом (ожидаем 409 Conflict)
curl -s -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"barcode":"INV-2024-001","name":"Дубликат","category_id":1,"room_id":1}' | jq .

# Утилизировать предмет (ожидаем 204 No Content)
curl -s -o /dev/null -w "%{http_code}" -X DELETE http://localhost:8080/api/v1/items/1 \
  -H "Authorization: Bearer $TOKEN"

# Попытка изменить утилизированный предмет (ожидаем 422)
curl -s -X PATCH http://localhost:8080/api/v1/items/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"description":"попытка изменить"}' | jq .
