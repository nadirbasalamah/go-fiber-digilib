package handlers

import (
	"context"
	"mime/multipart"

	"go-fiber-digilib/pkg/ai"
	"go-fiber-digilib/pkg/rajaongkir"
)

type stubUploader struct {
	url string
	err error
}

func (s stubUploader) UploadFile(_ context.Context, _ *multipart.FileHeader) (string, error) {
	return s.url, s.err
}

type stubRajaOngkir struct {
	fee float64
	err error
}

func (s stubRajaOngkir) GetDeliveryFee(_ rajaongkir.GetFeeRequest) (float64, error) {
	return s.fee, s.err
}

type stubRecommendation struct {
	response ai.PromptResponse
	err      error
}

func (s stubRecommendation) GetBookRecommendation(_ ai.BookRecommendationRequest) (ai.PromptResponse, error) {
	return s.response, s.err
}
