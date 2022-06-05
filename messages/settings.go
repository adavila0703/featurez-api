package messages

type UpdateSettingsRequest struct {
	RedisAddress string `json:"redis_address"`
}

type UpdateSettingsResponse struct {
	Message string `json:"message"`
}

type GetSettingsRequest struct {
}

type GetSettingsResponse struct {
	RedisAddress string `json:"redis_address"`
}
