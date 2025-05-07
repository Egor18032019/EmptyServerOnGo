package services

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	logsDir        = "logs"    // Директория для хранения логов
	logFileName    = "app.log" // Имя основного файла лога
	maxLogFileSize = 5 << 20   // 5MB - максимальный размер файла перед ротацией
	logPrefix      = "HTTP: "  // Префикс для лог-сообщений
)

var (
	fileLogger *log.Logger
)

// InitLogging инициализирует систему логирования
func InitLogging() error {
	// Создаем директорию для логов, если ее нет
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return err
	}

	// Полный путь к файлу лога
	logPath := filepath.Join(logsDir, logFileName)

	// Открываем файл лога
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Настраиваем логгер
	fileLogger = log.New(logFile, logPrefix, log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

// LogRequest записывает информацию о запросе в лог
func LogRequest(r *http.Request, message string, success bool) {
	if fileLogger == nil {
		log.Println("Логгер не инициализирован")
		return
	}

	if err := rotateLogIfNeeded(); err != nil {
		fileLogger.Printf("Ошибка ротации лога: %v", err)
	}

	status := "FAIL"
	if success {
		status = "SUCCESS"
	}

	fileLogger.Printf("%s %s %s %s %s",
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		status,
		message)
}
func LogMessage(message string) {
	if fileLogger == nil {
		log.Println("Логгер не инициализирован")
		return
	}

	if err := rotateLogIfNeeded(); err != nil {
		fileLogger.Printf("Ошибка ротации лога: %v", err)
	}

	fileLogger.Printf(message)
}

// rotateLogIfNeeded проверяет размер файла лога и при необходимости создает новый
func rotateLogIfNeeded() error {
	logPath := filepath.Join(logsDir, logFileName)

	info, err := os.Stat(logPath)
	if err != nil {
		return err
	}

	if info.Size() < maxLogFileSize {
		return nil
	}

	// Формируем имя для архивного файла
	backupName := filepath.Join(logsDir, "app_"+time.Now().Format("20060102_150405")+".log")

	// Закрываем текущий файл перед ротацией
	if f, ok := fileLogger.Writer().(*os.File); ok {
		if err := f.Close(); err != nil {
			return err
		}
	}

	// Переименовываем текущий лог-файл
	if err := os.Rename(logPath, backupName); err != nil {
		return err
	}

	// Создаем новый лог-файл
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Обновляем writer для логгера
	fileLogger.SetOutput(logFile)
	return nil
}
