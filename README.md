## topenergy-interview

# How to run?

## Docker compose

Make sure that there are not containers that use our ports.
Ports in usage are: 6379, 8080, 16686, 14268.
After that run the following commands in your terminal.

```
docker-compose build
```

```
docker-compose up --remove-orphans
```

## Routes

From here your can hit `{{host}}/routes` route with `GET` method to get all the availabe routes.
And check if the app is ready by hitting `{{host}}/health/ready` route with `GET` method.

## Usage

1. POST /tasks: Создает новую задачу. Тело запроса должно содержать заголовок и описание задачи. Возвращает идентификатор новой задачи.

2. GET /tasks: Возвращает список всех задач.

3. GET /tasks/{id}: Возвращает детали задачи по идентификатору.

4. PUT /tasks/{id}: Обновляет задачу по идентификатору. Тело запроса должно содержать новый заголовок и описание задачи.

5. DELETE /tasks/{id}: Удаляет задачу по идентификатору.
