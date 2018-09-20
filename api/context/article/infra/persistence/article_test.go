package persistence

import (
	"lmm/api/context/article/domain/model"
	"lmm/api/testing"
)

func TestSaveArticle(tt *testing.T) {
	t := testing.NewTester(tt)

	author, err := authorService.AuthorFromUserID(user.ID())
	t.NoError(err)

	article, err := articleService.NewArticleToPost(author, "title", "body", make([]string, 0))
	t.NoError(err)

	// save new article
	t.NoError(articleRepository.Save(article))

	text, err := model.NewText("new title", "new body")
	t.NoError(err)
	t.NoError(article.EditText(text))

	// save updated article
	t.NoError(articleRepository.Save(article))
}