package image

import (
	model "lmm/api/domain/model/image"
)

func Add(userID int64, t model.ImageType, url string) error {
	return nil
}

func Fetch(userID int64, t model.ImageType) ([]model.Image, error) {
	return nil, nil
}
