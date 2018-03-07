package category

import (
	"lmm/api/domain/model/blog"
	model "lmm/api/domain/model/category"
	repo "lmm/api/domain/repository/category"
	blogUsecase "lmm/api/usecase/blog"
	"strings"

	"github.com/akinaru-lu/errors"
)

func Register(userID int64, name string) (int64, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, errors.New("Empty name")
	}
	return repo.Add(userID, name)
}

func Update(userID, categoryID int64, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("Empty name")
	}
	return repo.Update(userID, categoryID, name)
}

func FetchByID(categoryID int64) (*model.Category, error) {
	return repo.ByID(categoryID)
}

func FetchByUser(userID int64) ([]model.Category, error) {
	return repo.ByUser(userID)
}

func FetchByBlog(blogID int64) (*model.Category, error) {
	return repo.ByBlog(blogID)
}

func SetBlogCategory(userID, blogID, categoryID int64) error {
	if err := blogUsecase.CheckOwnership(userID, blogID); err != nil {
		return errors.Wrap(err, "Not allowed to edit target blog")
	}
	if err := CheckOwnership(userID, categoryID); err != nil {
		return errors.Wrap(err, "Not allowd to add targer category")
	}
	return repo.SetBlogCategory(blogID, categoryID)
}

func FetchAllBlog(categoryID int64) ([]blog.ListItem, error) {
	return repo.AllBlogByID(categoryID)
}

func CheckOwnership(userID, categoryID int64) error {
	category, err := FetchByID(categoryID)
	if err != nil {
		return err
	}
	if category.ID != categoryID {
		return errors.New("User doesn't own the target category")
	}
	return nil
}

func Delete(userID, categoryID int64) error {
	blogList, err := FetchAllBlog(categoryID)
	if err != nil {
		return err
	}
	if blogList != nil && len(blogList) != 0 {
		return errors.New("There are still blog in this category")
	}
	return repo.Delete(userID, categoryID)
}
