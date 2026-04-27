sudo docker exec -i skopidom-postgres-1 psql -U skopidom -d skopidom -c "
INSERT INTO users (full_name, email, password_hash, role)
VALUES (
  'Администратор',
  'admin@university.ru',
  '\$2a\$12\$iB.8j9wRbmfza6qfuGTBn.l2dCkBc9ojVcIWnZi80nDM4bwO0RhEy',
  'admin'
);"
