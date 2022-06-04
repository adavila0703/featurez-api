package messages

type UpdateSettingsRequest struct {
	RedisAddress string `json:"redis_address"`
}

type UpdateSettingsResponse struct {
	Message string `json:"message"`
}
