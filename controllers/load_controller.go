package controllers

import (
	"io"
	"my-go-webserver/global"
	"my-go-webserver/services"
	"net/http"
	"os"
	"path/filepath"
)

// Максимальный размер файла (ограничивает объем загружаемого файла)
const maxUploadSize = int64(1 << 20)
const usersFilePath = "uploads"

// Обработчик домашней страницы с возможностью загрузки файла
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := global.Store.Get(r, "session-name")
	// Получаем username из сессии
	username, ok := session.Values["username"].(string)
	if !ok {
		// Значение отсутствует или имеет неверный тип
		println("Имя пользователя не найдено в сессии")
		return
	}
	authenticated, ok := session.Values["authenticated"]
	if !ok || authenticated != true {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Ограничиваем размер всего запроса
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	// Обрабатываем загрузку файла
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(maxUploadSize)
		if err != nil {
			http.Error(w, "Файл превышает допустимый размер.", http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Сохраняем файл на сервере
		fname := handler.Filename
		dstPath := filepath.Join(usersFilePath, fname)
		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, "Ошибка при создании файла", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
			return
		}
		services.LogRequest(r, "Пользователь "+username+" загрузил файл "+fname, true)
		// Отвечаем пользователю об успешной загрузке
		w.Write([]byte("Файл успешно загружен."))
		return
	}

	// Если метод GET, выводим простую форму для загрузки файла
	html := `
<!DOCTYPE html>
<html lang="ru">
<head><title>Загрузка файла</title></head>
<body>
<h1>Загрузите файл:</h1>
<form method="post" enctype="multipart/form-data">
<input type="file" name="file"><br/><br/>
<input type="submit" value="Загрузить">
</form>
</body>
</html>
`
	w.Write([]byte(html))
}
