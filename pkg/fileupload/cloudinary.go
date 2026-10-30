package fileupload

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryConfig struct {
	CloudinaryURL string
}

type Uploader interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (*FileUploadResult, error)
	DeleteFile(ctx context.Context, fileID string, resourceType string) (*uploader.DestroyResult, error)
}

type FileUploadResult struct {
	SecureURL string
	PublicID  string
}

type CloudinaryUploader struct {
	Cld     *cloudinary.Cloudinary
	Options uploader.UploadParams
}

func (c *CloudinaryConfig) InitCloudinary() *cloudinary.Cloudinary {
	cld, err := cloudinary.NewFromURL(c.CloudinaryURL)

	if err != nil {
		log.Fatalf("error when connecting to the cloudinary: %s\n", err)
	}

	log.Println("cloudinary initialization succeed")

	return cld
}

func (c *CloudinaryUploader) UploadFile(ctx context.Context, file *multipart.FileHeader) (*FileUploadResult, error) {
	resp, err := c.Cld.Upload.Upload(ctx, file, c.Options)

	if err != nil {
		return nil, err
	}

	return &FileUploadResult{
		SecureURL: resp.SecureURL,
		PublicID:  resp.PublicID,
	}, nil
}

func (c *CloudinaryUploader) DeleteFile(ctx context.Context, fileID string, resourceType string) (*uploader.DestroyResult, error) {
	resp, err := c.Cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     fileID,
		ResourceType: resourceType, // Options: "image", "video", "raw"
		Invalidate:   c.Options.UseFilename,
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
