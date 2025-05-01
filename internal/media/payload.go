package media

type ProductCreatedEvent struct {
    Action    string   `json:"action"`
    ProductID uint     `json:"product_id"`
    MediaKeys []string `json:"mediaKeys"`
}

type MediaStoredEvent struct {
    ProductID uint     `json:"product_id"`
    URLs      []string `json:"urls"`
}