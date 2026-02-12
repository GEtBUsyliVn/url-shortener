package model

type CreateShortCodeRequest struct {
	URL         string `json:"original_url" binding:"required" validate:"url"`
	UserId      string `json:"user_id" binding:"required"`
	ExpiredDays int    `json:"expired_days" binding:"required"`
}
