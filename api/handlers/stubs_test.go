package handlers

import (
	"context"
	"mime/multipart"

	"go-fiber-digilib/pkg/ai"
	"go-fiber-digilib/pkg/fileupload"
	"go-fiber-digilib/pkg/rajaongkir"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type stubUploader struct {
	url       string
	publicID  string
	err       error
	deleteErr error
	deletedID string
}

func (s *stubUploader) UploadFile(_ context.Context, _ *multipart.FileHeader) (*fileupload.FileUploadResult, error) {
	if s.err != nil {
		return nil, s.err
	}

	return &fileupload.FileUploadResult{SecureURL: s.url, PublicID: s.publicID}, nil
}

func (s *stubUploader) DeleteFile(_ context.Context, fileID string, _ string) (*uploader.DestroyResult, error) {
	s.deletedID = fileID

	if s.deleteErr != nil {
		return nil, s.deleteErr
	}

	return &uploader.DestroyResult{Result: "ok"}, nil
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
