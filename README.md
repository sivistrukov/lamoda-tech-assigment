# Lamoda-tech assigment

## Установка и запуск

```shell
git clone https://github.com/sivistrukov/lamoda-tech-assigment.git
cd lamoda-tech-assigment
make
```

## Методы с примерами запросов и ответов

*(swagger документация для API расположена на http://localhost:8080/swagger/index.html)*

- [post]  {{baseUrl}}/v1/products/ - добавление нового товара

Запрос:

```shell
curl --location 'http://localhost:8080/api/v1/products/' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--data '{
  "code": "D",
  "name": "D Product",
  "size": "size"
}'
```

Тело ответа:

```json
{
  "ok": true
}
```

- [post]  {{baseUrl}}/v1/warehouses/ - добавление нового склада

Запрос:

```shell
curl --location 'http://localhost:8080/api/v1/warehouses/' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--data '{
  "name": "New warehouse"
}'
```

Тело ответа:

```json
{
  "id": 3,
  "name": "New warehouse",
  "isAvailable": true
}
```

- [get]   {{baseUrl}}/v1/warehouses/{id}/products/ - получить запас хранимых товаров на складе

Запрос:

```shell
curl --location 'http://localhost:8080/api/v1/warehouses/1/products/' \
--header 'Accept: application/json'
```

Тело ответа:

```json
[
  {
    "code": "A001",
    "name": "Product A",
    "size": "?",
    "quantity": 100
  },
  {
    "code": "B001",
    "name": "Product B",
    "size": "?",
    "quantity": 50
  }
]
```

- [post]  {{baseUrl}}/v1/warehouses/{id}/products/ - добавление товаров на склад

Запрос:

```shell
curl --location 'http://localhost:8080/api/v1/warehouses/1/products/' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--data '[
  {
    "code": "A001",
    "quantity": 10
  },
  {
    "code": "B001",
    "quantity": 10
  }
]'
```

Тело ответа:

```json
[
  {
    "code": "A001",
    "quantity": 110
  },
  {
    "code": "B001",
    "quantity": 60
  }
]
```

- [get]   {{baseUrl}}/v1/warehouses/{id}/products/{code}/quantity/ - получить кол-во доступного товара на складе

Запрос:

```shell
curl --location 'http://localhost:8080/api/v1/warehouses/1/products/A001/quantity' \
--header 'Accept: application/json'
```

Тело ответа:

```json
{
  "code": "A001",
  "quantity": 100
}
```

- [post]  {{baseUrl}}/v1/warehouses/{id}/reserve/ - резервирование товара на складе для доставки

Запрос:

```shell
curl --location 'http://localhost:8080/api/v1/warehouses/1/reserve/' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--data '[
  {
    "code": "A001",
    "quantity": 50
  },
  {
    "code": "B001",
    "quantity": 20
  }
]'
```

Тело ответа:

```json
{
  "ok": true
}
```

- [post]  {{baseUrl}}/v1/warehouses/{id}/cancel-reservation/ - освобождение резерва товаров

Запрос:

```shell
curl --location 'http://localhost:8080/api/v1/warehouses/1/cancel-reservation/' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--data '[
  {
    "code": "A001",
    "quantity": 50
  },
  {
    "code": "C001",
    "quantity": 20
  }
]'
```

Тело ответа:

```json
{
  "ok": true
}
```

## Данные

Склады:

| id | name   | is_available |
|----|--------|--------------|
| 1  | Main   | true         |
| 2  | Second | true         |

Товары:

| code | name      | size |
|------|-----------|------|
| A001 | Product A | ?    |
| B001 | Product B | ?    |
| C001 | Product C | ?    |

Товары на складе:

| id | warehouse_id | product_code | quantity | status      |
|----|--------------|--------------|----------|-------------|
| 1  | 1            | A001         | 100      | available   |
| 2  | 1            | B001         | 50       | available   |
| 3  | 1            | C001         | 100      | reservation |
| 4  | 1            | A001         | 100      | reservation |
| 5  | 1            | C001         | 50       | available   |
| 6  | 1            | A001         | 50       | available   |
| 7  | 1            | A001         | 50       | reservation |
