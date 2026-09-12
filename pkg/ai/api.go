package ai

import (
	"encoding/json"
	"fmt"
	"go-digilib/pkg/clients"
	"go-digilib/pkg/constant"
	"go-digilib/pkg/utils"
	"net/http"
)

const BASE_URL = "https://openrouter.ai/api/v1"

type Service struct {
	client clients.HTTPClient
}

func InitService() Service {
	return Service{
		client: clients.InitHTTPClient(BASE_URL, 40, utils.GetConfig(constant.AI_API_KEY)),
	}
}

func (r *Service) GetBookRecommendation(req BookRecommendationRequest) (PromptResponse, error) {
	var response PromptResponse

	var model string = utils.GetConfig(constant.AI_MODEL)
	var userPrompt string = fmt.Sprintf("Suggest top %v book recommendations about %v", req.Quantity, req.Topic)

	payload := map[string]any{
		"model": model,
		"messages": []Message{
			{
				Role:    "system",
				Content: SYSTEM_PROMPT,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
	}

	res, err := r.client.SendJSON(
		"/chat/completions",
		http.MethodPost,
		payload,
	)

	if err != nil {
		return PromptResponse{}, err
	}

	if err := json.Unmarshal([]byte(res), &response); err != nil {
		return PromptResponse{}, err
	}

	return response, nil
}
