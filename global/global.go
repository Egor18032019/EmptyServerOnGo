package global

import "github.com/gorilla/sessions"

/*
Хранилище генерирует уникальный идентификатор сеанса и записывает его в cookie, отправляемый браузеру.
*/
var Store = sessions.NewCookieStore([]byte("secret-key")) // хранилище сесий
