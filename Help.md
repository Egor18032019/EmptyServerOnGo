Ставим зависимости:

```shell
go get github.com/gorilla/sessions 
```

Авторизация
```shell
curl -i -X POST http://127.0.0.1:8080/login \
 -H 'Content-Type: application/json' -d '{"username":"slava" , "password":"123123123"}'
```
