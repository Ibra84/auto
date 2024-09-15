package communication

import (
	"auto/pkg/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const globalServerURL = "http://globalserver.example.com/upload"

func SendDataToGlobalServer(mediaData models.MediaData) error {
	// Создание тела запроса
	body := map[string]interface{}{
		"car_number":     mediaData.CarNumber,
		"timestamp":      mediaData.Timestamp,
		"photo1":         mediaData.PhotoPath1,
		"photo2":         mediaData.PhotoPath2,
		"photo3":         mediaData.PhotoPath3,
		"video":          mediaData.VideoPath,
		"first_request":  mediaData.FirstRequest,
		"second_request": mediaData.SecondRequest,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// Отправка POST запроса на глобальный сервер
	resp, err := http.Post(globalServerURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send data, status code: %d", resp.StatusCode)
	}

	// Установка полей FirstRequest и SecondRequest в true после успешной отправки
	mediaData.FirstRequest = true
	mediaData.SecondRequest = true

	return nil
}
