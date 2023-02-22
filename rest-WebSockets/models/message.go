package models

type WebsocketMessag struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
