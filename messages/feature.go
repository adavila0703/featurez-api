package messages

type CreateFeatureRequest struct {
	Name string `json:"name"`
}

type CreateFeatureResponse struct {
	Message     string `json:"message"`
	FeatureFlag string `json:"feature_flag"`
}

type GetFeatureListRequest struct {
}

type Feature struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type GetFeatureListResponse struct {
	FeatureList []*Feature `json:"feature_list"`
}

type DeleteFeatureRequest struct {
	Name []string `json:"name"`
}

type DeleteFeatureResponse struct {
	Message string `json:"message"`
}

type UpdateFeatureRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type UpdateFeatureResponse struct {
	Message string `json:"message"`
}
