package controllers

import (
	"my-go-webserver/global"
	"my-go-webserver/models"
	"my-go-webserver/services"
	"net/http"
	"strings"
)

// Обработчик POST /login
func LoginHandler(writer http.ResponseWriter, reader *http.Request) {
	users, err := models.LoadUsersFromJSON()
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	switch reader.Method {
	case http.MethodGet:
		// Показываем форму авторизации
		html := `<html>
            <head><title>Логин</title></head>
            <body>
                <h1>Авторизация</h1>
                <form action="/login" method="post">
                    Имя пользователя:<input type="text" name="username"><br />
                    Пароль:<input type="password" name="password"><br />
                    <button type="submit">Войти</button>
                </form>
            </body>
        </html>`
		writer.Write([]byte(html))

	case http.MethodPost:
		// Авторизуем пользователя
		reader.ParseForm()
		username := strings.TrimSpace(reader.PostFormValue("username"))
		password := strings.TrimSpace(reader.PostFormValue("password"))
		services.LogRequest(reader, "Попытка авторизации от: "+username+" "+password, true)
		for _, u := range users {
			// Проверка введенных данных
			if u.Username == username && u.Password == password {
				session, _ := global.Store.Get(reader, "session-name")
				session.Values["authenticated"] = true
				session.Values["username"] = username
				session.Save(reader, writer)
				http.Redirect(writer, reader, "/home", http.StatusFound)
				services.LogRequest(reader, "Пользователь авторизован - "+username, true)
				return
			}
		}
		http.Error(writer, "Неправильное имя пользователя или пароль", http.StatusUnauthorized)

	default:
		http.Error(writer, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

//// Домашняя страница выгрузка
//func HomeHandler(w http.ResponseWriter, r *http.Request) {
//	session, _ := store.Get(r, "session-name")
//	authenticated, ok := session.Values["authenticated"]
//	if !ok || authenticated != true {
//		http.Redirect(w, r, "/login", http.StatusFound)
//		return
//	}
//
//	filename := "data.txt" // Путь к файлу, который хотим отдать пользователю
//	services.SendFile(w, r, filename)
//}
