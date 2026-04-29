TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@university.ru","password":"password"}' | jq -r .token)

echo $TOKEN
