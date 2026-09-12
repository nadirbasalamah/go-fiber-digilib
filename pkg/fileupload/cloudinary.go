package fileupload

import (
	"context"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryConfig struct {
	CloudinaryURL string
}

type CloudinaryUploader struct {
	Cld     *cloudinary.Cloudinary
	File    any
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

func (c *CloudinaryUploader) UploadFile(ctx context.Context) (string, error) {
	resp, err := c.Cld.Upload.Upload(ctx, c.File, c.Options)

	if err != nil {
		return "", err
	}

	return resp.SecureURL, nil
}
