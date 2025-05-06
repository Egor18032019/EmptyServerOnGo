package services

import (
	"net/http"
	"os"
	"path/filepath"
)

// SendFile Функция отправки файла пользователю
func SendFile(w http.ResponseWriter, r *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Ошибка при чтении метаданных файла", http.StatusInternalServerError)
		return
	}

	//// Установка правильного MIME-типа (можно использовать mime.TypeByExtension)
	//contentType := http.DetectContentType([]byte{})
	http.ServeContent(w, r, filepath.Base(filename), stat.ModTime(), file)
}
