package controllers

import (
	"io"
	"net/http"
	"os"
)

// Максимальный размер файла (ограничивает объем загружаемого файла)
const maxUploadSize = int64(1 << 20) // 1MB

// Обработчик домашней страницы с возможностью загрузки файла
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	authenticated, ok := session.Values["authenticated"]
	if !ok || authenticated != true {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Обрабатываем загрузку файла
	if r.Method == http.MethodPost {
		err := r.ParseMultipartForm(maxUploadSize)
		if err != nil {
			http.Error(w, "Ошибка анализа формы", http.StatusBadRequest)
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
		dst, err := os.Create(fname)
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
