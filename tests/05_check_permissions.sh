# Залогиниться как editor
EDITOR_TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ivanov@university.ru","password":"securepass123"}' | jq -r .token)

# Editor может создавать предметы (ожидаем 201)
curl -s -X POST http://localhost:8080/api/v1/items \
  -H "Authorization: Bearer $EDITOR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "barcode": "INV-2024-003",
    "name": "HP LaserJet Pro",
    "category_id": 1,
    "room_id": 1
  }' | jq .

# Editor НЕ может утилизировать (ожидаем 403 Forbidden)
curl -s -o /dev/null -w "%{http_code}" -X DELETE http://localhost:8080/api/v1/items/2 \
  -H "Authorization: Bearer $EDITOR_TOKEN"

# Без токена (ожидаем 401 Unauthorized)
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/v1/items
