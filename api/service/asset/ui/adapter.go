package ui

import "lmm/api/service/asset/domain/model"

type imageListItem struct {
	Name string `json:"name"`
}
type imageListJSON struct {
	Images     []imageListItem `json:"images"`
	NextCursor *uint           `json:"nextCursor"`
}

func imageCollectionToJSON(collection *model.ImageCollection) *imageListJSON {
	images := make([]imageListItem, len(collection.List()), len(collection.List()))
	for i, image := range collection.List() {
		images[i].Name = image.Name()
	}
	return &imageListJSON{
		Images:     images,
		NextCursor: collection.NextID(),
	}
}
