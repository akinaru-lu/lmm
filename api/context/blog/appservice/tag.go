package appservice

import (
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/domain"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	"lmm/api/utils/strings"
)

func (app *AppService) AddNewTagToBlog(user *account.User, blogIDStr, tagName string) error {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return domain.ErrNoSuchBlog
	}

	blog, err := app.blogService.GetBlogByID(blogID)
	if err != nil {
		return domain.ErrNoSuchBlog
	}

	tag, err := factory.NewTag(blog.ID(), tagName)
	if err != nil {
		return err
	}

	return app.tagRepository.Add(tag)
}

func (app *AppService) UpdateBlogTag(user *account.User, tagIDStr, tagName string) error {
	tagID, err := strings.StrToUint64(tagIDStr)
	if err != nil {
		return domain.ErrNoSuchTag
	}

	// TODO using transtraction
	tag, err := app.tagRepository.FindByID(tagID)
	if err != nil {
		return domain.ErrNoSuchTag
	}

	if err := tag.UpdateName(tagName); err != nil {
		return err
	}

	return app.tagRepository.Update(tag)
}

func (app *AppService) RemoveBlogTag(user *account.User, tagIDStr string) error {
	tagID, err := strings.StrToUint64(tagIDStr)
	if err != nil {
		return domain.ErrNoSuchTag
	}

	tag, err := app.tagRepository.FindByID(tagID)
	if err != nil {
		return domain.ErrNoSuchTag
	}

	return app.tagRepository.Remove(tag)
}

func (app *AppService) GetAllTags() ([]*model.Tag, error) {
	return app.tagRepository.FindAll()
}

func (app *AppService) GetAllTagsOfBlog(blogIDStr string) ([]*model.Tag, error) {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return nil, domain.ErrNoSuchBlog
	}

	blog, err := app.blogService.GetBlogByID(blogID)
	if err != nil {
		return nil, domain.ErrNoSuchBlog
	}

	return app.tagRepository.FindAllByBlog(blog)
}
