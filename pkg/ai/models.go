package ai

type BookRecommendationRequest struct {
	Quantity uint   `json:"quantity" validate:"required,min=1"`
	Topic    string `json:"topic" validate:"required"`
}

type BookRecommendationResponse struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Publisher string `json:"publisher"`
	Year      int    `json:"year"`
}

type PromptRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PromptResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Provider          string   `json:"provider"`
	SystemFingerprint string   `json:"system_fingerprint"`
	ServiceTier       any      `json:"service_tier"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	Index              int64           `json:"index"`
	Logprobs           any             `json:"logprobs"`
	FinishReason       string          `json:"finish_reason"`
	NativeFinishReason string          `json:"native_finish_reason"`
	Message            ResponseMessage `json:"message"`
}

type ResponseMessage struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Refusal   any    `json:"refusal"`
	Reasoning any    `json:"reasoning"`
}

type Usage struct {
	PromptTokens            int64                   `json:"prompt_tokens"`
	CompletionTokens        int64                   `json:"completion_tokens"`
	TotalTokens             int64                   `json:"total_tokens"`
	Cost                    int64                   `json:"cost"`
	IsByok                  bool                    `json:"is_byok"`
	PromptTokensDetails     PromptTokensDetails     `json:"prompt_tokens_details"`
	CostDetails             CostDetails             `json:"cost_details"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
}

type CompletionTokensDetails struct {
	ReasoningTokens int64 `json:"reasoning_tokens"`
	ImageTokens     int64 `json:"image_tokens"`
	AudioTokens     int64 `json:"audio_tokens"`
}

type CostDetails struct {
	UpstreamInferenceCost            int64 `json:"upstream_inference_cost"`
	UpstreamInferencePromptCost      int64 `json:"upstream_inference_prompt_cost"`
	UpstreamInferenceCompletionsCost int64 `json:"upstream_inference_completions_cost"`
}

type PromptTokensDetails struct {
	CachedTokens     int64 `json:"cached_tokens"`
	CacheWriteTokens int64 `json:"cache_write_tokens"`
	AudioTokens      int64 `json:"audio_tokens"`
	VideoTokens      int64 `json:"video_tokens"`
}
