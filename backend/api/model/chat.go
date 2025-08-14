package model

type SendMessageRequest struct {
	To      uint   `json:"to" binding:"required"`
	Type    string `json:"type"`
	Content string `json:"content" binding:"required"`
}
