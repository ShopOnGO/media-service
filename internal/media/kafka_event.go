package media

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ShopOnGO/ShopOnGO/pkg/logger"
)

func (h *MediaHandler) HandleMediaEvent(ctx context.Context, msg []byte) error {
	logger.Infof("Получено сообщение для медиа: %s", string(msg))

	var event ProductCreatedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
        return fmt.Errorf("invalid product-created payload: %w", err)
    }

	eventHandlers := map[string]func(context.Context, *MediaService, ProductCreatedEvent) error{
		"create": h.HandleCreateMediaEvent,
		// "update": HandleUpdateMediaEvent,
		// "delete": HandleDeleteMediaEvent,
	}

	handler, exists := eventHandlers[event.Action]
	if !exists {
		return fmt.Errorf("неизвестное действие для медиа: %s", event.Action)
	}

	return handler(ctx, h.Mediasvc, event)
}

func (h *MediaHandler) HandleCreateMediaEvent(ctx context.Context, mediaSvc *MediaService, event ProductCreatedEvent) error {
	logger.Infof("Создание медиа для продукта %d", event.ProductID)
	
	var urls []string
	for _, key := range event.MediaKeys {
		url := mediaSvc.GenerateURL(key)
		urls = append(urls, url)
	}

	out := MediaStoredEvent{
		ProductID: event.ProductID,
		URLs:      urls,
	}

	payload, err := json.Marshal(out)
	if err != nil {
		return fmt.Errorf("ошибка сериализации события media-stored: %w", err)
	}

	return h.Kafka.Produce(ctx, []byte("media-stored"), payload)
}
