# Quotes API

АПИ на Go для хранения и управления цитатами.

---

##  Запуск

```bash
git clone git@github.com:NessibeliY/quotes-api.git
cd quotes-api
go mod tidy
go run .
```


⸻

## Тестирование

1. Добавить цитату:

```
curl -X POST http://localhost:8080/quotes \
  -H "Content-Type: application/json" \
  -d '{"author":"Confucius","quote":"Life is simple, but we insist on making it complicated."}'
```

2. Получить все цитаты:
```
curl http://localhost:8080/quotes
```

3. Получить случайную цитату:
```
curl http://localhost:8080/quotes/random
```

4. Фильтрация по автору:
```
curl "http://localhost:8080/quotes?author=Confucius"
```

5. Удалить цитату:
```
curl -X DELETE http://localhost:8080/quotes/1
```
