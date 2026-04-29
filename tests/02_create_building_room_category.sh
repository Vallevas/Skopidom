# Создать здание
curl -s -X POST http://localhost:8080/api/v1/buildings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Главный корпус","address":"ул. Университетская, 1"}' | jq .

# Создать комнату (building_id=1)
curl -s -X POST http://localhost:8080/api/v1/rooms \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Аудитория 301","building_id":1}' | jq .

# Создать категорию
curl -s -X POST http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Компьютер"}' | jq .

# Получить все здания
curl -s http://localhost:8080/api/v1/buildings \
  -H "Authorization: Bearer $TOKEN" | jq .

# Получить комнаты по зданию
curl -s "http://localhost:8080/api/v1/rooms?building_id=1" \
  -H "Authorization: Bearer $TOKEN" | jq .

# Получить все категории
curl -s http://localhost:8080/api/v1/categories \
  -H "Authorization: Bearer $TOKEN" | jq .
