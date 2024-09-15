package handlers

import (
	"auto/pkg/communication"
	"auto/pkg/models"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) UploadData(c *gin.Context) {
	carNumber := c.PostForm("car_number")
	timestamp := c.PostForm("timestamp")

	// Получение файлов из запроса
	photo1, err := c.FormFile("photo1")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to open photo1 file"})
		return
	}

	photo2, err := c.FormFile("photo2")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to open photo2 file"})
		return
	}

	photo3, err := c.FormFile("photo3")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to open photo3 file"})
		return
	}

	video, err := c.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to open video file"})
		return
	}

	// Создание директории uploads, если она не существует
	uploadDir := "./uploads/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Сохранение фото на диск
	photoPath1 := uploadDir + photo1.Filename
	outPhoto1, err := os.Create(photoPath1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save photo1 file"})
		return
	}
	defer outPhoto1.Close()

	srcPhoto1, err := photo1.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open photo1 file"})
		return
	}
	defer srcPhoto1.Close()

	_, err = io.Copy(outPhoto1, srcPhoto1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save photo1 file"})
		return
	}

	photoPath2 := uploadDir + photo2.Filename
	outPhoto2, err := os.Create(photoPath2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save photo2 file"})
		return
	}
	defer outPhoto2.Close()

	srcPhoto2, err := photo2.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open photo2 file"})
		return
	}
	defer srcPhoto2.Close()

	_, err = io.Copy(outPhoto2, srcPhoto2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save photo2 file"})
		return
	}

	photoPath3 := uploadDir + photo3.Filename
	outPhoto3, err := os.Create(photoPath3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save photo3 file"})
		return
	}
	defer outPhoto3.Close()

	srcPhoto3, err := photo3.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open photo3 file"})
		return
	}
	defer srcPhoto3.Close()

	_, err = io.Copy(outPhoto3, srcPhoto3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save photo3 file"})
		return
	}

	// Сохранение видео на диск
	videoPath := uploadDir + video.Filename
	outVideo, err := os.Create(videoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save video file"})
		return
	}
	defer outVideo.Close()

	srcVideo, err := video.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open video file"})
		return
	}
	defer srcVideo.Close()

	_, err = io.Copy(outVideo, srcVideo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save video file"})
		return
	}

	// Сохранение метаданных и путей к файлам в базе данных
	mediaData := models.MediaData{
		CarNumber:  carNumber,
		Timestamp:  timestamp,
		PhotoPath1: photoPath1,
		PhotoPath2: photoPath2,
		PhotoPath3: photoPath3,
		VideoPath:  videoPath,
	}

	if err := h.DB.Create(&mediaData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data"})
		return
	}

	// Отправка данных на глобальный сервер
	err = communication.SendDataToGlobalServer(mediaData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to global server"})
		return
	}

	// Обновление статусов после успешной отправки данных
	if !mediaData.FirstRequest {
		mediaData.FirstRequest = true
		if err := h.DB.Save(&mediaData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update first request status"})
			return
		}
	}

	if mediaData.FirstRequest && !mediaData.SecondRequest {
		err = communication.SendDataToGlobalServer(mediaData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send second request to global server"})
			return
		}

		// Обновляем статус второго запроса после успешной отправки данных
		mediaData.SecondRequest = true
		if err := h.DB.Save(&mediaData).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update second request status"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data uploaded and sent successfully"})
}
