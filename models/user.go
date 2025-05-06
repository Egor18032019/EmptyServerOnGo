package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password_hash"`
}

const usersFilePath = "users.json"

// Читаем данные пользователей из файла
func LoadUsersFromJSON() ([]User, error) {
	data, err := ioutil.ReadFile(usersFilePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл '%s': %v", usersFilePath, err)
	}

	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора JSON: %v", err)
	}

	return users, nil
}
