package media

import (
    "github.com/ShopOnGO/ShopOnGO/pkg/kafkaService"

    "github.com/gin-gonic/gin"
)

type MediaHandlerDeps struct {
    Mediasvc     *MediaService
	Kafka   *kafkaService.KafkaService
}

type MediaHandler struct {
    Mediasvc     *MediaService
	Kafka   *kafkaService.KafkaService
}

func NewMediaHandler(router *gin.Engine, deps MediaHandlerDeps) *MediaHandler {
    handler := &MediaHandler{
        Mediasvc:   deps.Mediasvc,
        Kafka:      deps.Kafka,
    }
	mediaGroup := router.Group("/media")
	{
		mediaGroup.POST("/uploads", handler.HandleUploadHTTP)
	}

    return handler
}

// HandleUploadHTTP обрабатывает multipart/form-data
func (h *MediaHandler) HandleUploadHTTP(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": "no file"})
        return
    }
    url, err := h.Mediasvc.UploadFile(file)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"url": url})
}
