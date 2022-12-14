## file_parser

### Описание задания
Необходимо реализовать приложение, которое работает с файлом, в котором построчно указаны пары ключ-значение. Значения должны быть целыми числами.
Требуется раз в N времени считывать каждую последующую строку файла и записывать в БД полученные данные. Если значение кратно трем – его записывать не нужно.
Если ключ уже существует в БД – данные необходимо обновить.
Каждое событие обновления данных по ключу и отбрасывание значения кратного трем – необходимо логировать в консоль.
БД PostgreSQL или MongoDB.
Приложение должно принимать на вход параметры конфигурации:
1. Имя файла.
2. Частоту считывания строк.

### Способы запуска приложения:
1. Запуск с помощью go run:

```go
go run ./cmd/cli -f="test_data.txt" -t=3s
```
2. Запуск со сборкой бинарного файла:

```go
go build -o parser.exe ./cmd/cli
```

```go
./parser -f="test_data.txt" -t=2s
```


3. Запуск посредством докер образа. Значения параметров конфигурации
имени файла и частоты чтения задавать в строках:
```dockerfile
ENV f="./test_data.txt"
ENV t=2s
```
```dockerfile
docker build --tag parser .
```
```dockerfile
docker run parser
```

4. Запуск посредством docker-compose. Значения параметров конфигурации
имени файла и частоты чтения задавать в строках:
```dockerfile
parser:
...
FILE_NAME: ./test_data.txt
READ_FREQ: 2s
...
```
```dockerfile
docker compose up -d
```

 Для способов запуска 1-3 необходим запущенный экземпляр БД postgreSQL.
Необходимо изменить логин, пароль пользователя БД и номер порта на котором будет запущена БД.

Файл .env:
```markdown
DB_DSN="postgres://postgres:postgrespw@localhost:49159/user_db?sslmode=disable"
```
Также необходимо создать БД перед запуском приложения:
```postgresql
`CREATE DATABASE user_db`
```