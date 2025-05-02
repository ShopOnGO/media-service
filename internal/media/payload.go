package media

// Событие, которое приходит из product-service
type ProductCreatedEvent struct {
    Action    string   `json:"action"`
    ProductID uint     `json:"product_id"`
    ImageKeys  []string `json:"image_keys"`
	VideoKeys  []string `json:"video_keys"`
}

// Событие, которое уйдёт обратно в product-service
type MediaStoredEvent struct {
    Action    string   `json:"action"`
    ProductID uint     `json:"product_id"`
    ImageURLs []string `json:"image_urls"`
	VideoURLs []string `json:"video_urls"`
}