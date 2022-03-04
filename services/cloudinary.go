package services

import (
	"context"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

var (
	Cloudinary *cloudinary.Cloudinary
)

func NewCloudinaryInstance(cloud string, cloudKey string, secret string) error {

	cld, err := cloudinary.NewFromParams(cloud, cloudKey, secret)
	if err != nil {
		return err
	}

	Cloudinary = cld

	return nil
}

func UploadImage(file interface{}, filename string) (string, error) {

	ctx := context.Background()
	uploadResult, err := Cloudinary.Upload.Upload(
		ctx,
		file,
		uploader.UploadParams{
			PublicID:       filename,
			UniqueFilename: true,
			Folder:         "authentication-app",
			Overwrite:      true,
		})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}

func DeleteImage(publicID string) (string, error) {
	ctx := context.Background()
	resp, err := Cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image"})

	if err != nil {
		return "", err
	}

	return resp.Result, nil
}
