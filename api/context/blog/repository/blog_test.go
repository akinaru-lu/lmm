package repository

import (
	accountFactory "lmm/api/context/account/domain/factory"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/context/blog/domain/model"
	infra "lmm/api/db"
	"lmm/api/testing"
	"lmm/api/utils/uuid"
	"time"
)

func TestAddBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	title, text := uuid.New(), uuid.New()
	user, _ := accountFactory.NewUser(name, password)
	blog, _ := factory.NewBlog(user.ID(), title, text)
	err := repo.Add(blog)

	t.NoError(err)

	var (
		userID    uint64
		blogTitle string
		blogText  string
	)
	db := repo.DB()
	defer db.Close()

	stmt := db.MustPrepare(`SELECT user, title, text FROM blog WHERE id = ?`)
	defer stmt.Close()

	stmt.QueryRow(blog.ID()).Scan(&userID, &blogTitle, &blogText)

	t.Is(user.ID(), userID)
	t.Is(blog.Title(), blogTitle)
	t.Is(blog.Text(), blogText)
}

func TestAddBlog_DuplicateTitle(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	title, text := uuid.New(), uuid.New()
	user, _ := accountFactory.NewUser(name, password)
	blog, _ := factory.NewBlog(user.ID(), title, text)

	err := repo.Add(blog)
	t.NoError(err)

	blogWithSameTitle, _ := factory.NewBlog(user.ID(), title, text)
	err = repo.Add(blogWithSameTitle)
	t.Regexp(`Duplicate entry '[\w\d-]+' for key 'title'`, err.Error())
}

func TestFindAllBlog_FetchOneMore(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	user, _ := accountFactory.NewUser(name, password)

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	t.NoError(repo.Add(blog))

	title, text = uuid.New(), uuid.New()
	blog, _ = factory.NewBlog(user.ID(), title, text)
	t.NoError(repo.Add(blog))

	blogList, hasNextPage, err := repo.FindAll(1, 1)
	t.NoError(err)
	t.True(hasNextPage)
	t.Is(1, len(blogList))
}

func TestFindAllBlog_EmptyList(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	blogList, hasNextPage, err := repo.FindAll(100, 1)
	t.NoError(err)
	t.False(hasNextPage)
	t.NotNil(blogList)
	t.Is(0, len(blogList))
}

func TestFindAllBlogByCategory(tt *testing.T) {
	testing.Lock()
	defer testing.Unlock()

	testing.InitTable("blog")

	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	createBlogListWithCategory(repo, 10)

	c1, _ := model.NewCategory(1, "c1")
	c2, _ := model.NewCategory(1, "c2")

	if blogList, hasNextPage, err := repo.FindAllByCategory(c1, 5, 1); err == nil {
		t.False(hasNextPage)
		t.Is(5, len(blogList))

		blogList, hasNextPage, err = repo.FindAllByCategory(c1, 5, 2)
		t.NoError(err)
		t.False(hasNextPage)
		t.NotNil(blogList)
		t.Is(0, len(blogList))
	} else {
		t.Fatalf(err.Error())
	}

	if blogList, hasNextPage, err := repo.FindAllByCategory(c2, 2, 1); err == nil {
		t.True(hasNextPage)
		t.Is(2, len(blogList))

		blogList, hasNextPage, err = repo.FindAllByCategory(c2, 2, 2)
		t.NoError(err)
		t.True(hasNextPage)
		t.Is(2, len(blogList))

		blogList, hasNextPage, err = repo.FindAllByCategory(c2, 2, 3)
		t.NoError(err)
		t.False(hasNextPage)
		t.Is(1, len(blogList))
	} else {
		t.Fatalf(err.Error())
	}

	if blogList, hasNextPage, err := repo.FindAllByCategory(c2, 10, 1); err == nil {
		t.False(hasNextPage)
		t.Is(5, len(blogList))

		blogList, hasNextPage, err = repo.FindAllByCategory(c2, 10, 2)
		t.NoError(err)
		t.False(hasNextPage)
		t.Is(0, len(blogList))
	} else {
		t.Fatalf(err.Error())
	}
}

func TestFindBlogByID_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	title, text := uuid.New(), uuid.New()
	user, _ := accountFactory.NewUser(name, password)
	blog, _ := factory.NewBlog(user.ID(), title, text)
	err := repo.Add(blog)

	t.NoError(err)

	blogFound, err := repo.FindByID(blog.ID())
	t.NoError(err)
	t.Is(blog.ID(), blogFound.ID())
	t.Is(blog.Title(), blogFound.Title())
	t.Is(blog.Text(), blogFound.Text())
}

func TestFindBlogByID_NotFound(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	blog, err := repo.FindByID(uint64(777))
	t.Error(err, "sql: no rows in result set")
	t.Nil(blog)
}

func TestEditBlog_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	user, _ := accountFactory.NewUser(name, password)

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	t.NoError(repo.Add(blog))

	// no change
	t.Error(infra.ErrNoChange, repo.Update(blog))

	blog.UpdateTitle("new title")
	blog.UpdateText("new text")

	t.NoError(repo.Update(blog))

	blog, _ = repo.FindByID(blog.ID())

	t.Is("new title", blog.Title())
	t.Is("new text", blog.Text())
}

func TestEditBlog_NoSuchBlog(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	user, _ := accountFactory.NewUser(name, password)

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)

	// notice that it was not be saved
	t.Is(infra.ErrNoChange, repo.Update(blog))
}

func TestSetBlogCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	name, password := uuid.New()[:31], uuid.New()
	user, _ := accountFactory.NewUser(name, password)

	title, text := uuid.New(), uuid.New()
	blog, _ := factory.NewBlog(user.ID(), title, text)
	repo.Add(blog)

	category := newCategory()
	t.NoError(repo.SetBlogCategory(blog, category))

	blogList, hasNextPage, err := repo.FindAllByCategory(category, 10, 1)
	t.NoError(err)
	t.False(hasNextPage)
	t.Is(1, len(blogList))
	t.Is(blog.ID(), blogList[0].ID())

	otherCategory := newCategory()
	t.NoError(repo.SetBlogCategory(blog, otherCategory), "on duplicate")
	blogList, hasNextPage, err = repo.FindAllByCategory(otherCategory, 10, 1)
	t.NoError(err)
	t.False(hasNextPage)
	t.Is(1, len(blogList))
	t.Is(blog.ID(), blogList[0].ID())
}

func TestRemoveBlogCategory_Success(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	db := repo.DB()
	defer db.Close()

	user, _ := accountFactory.NewUser("username", "userpassword")
	blog, _ := factory.NewBlog(user.ID(), "blogtitle", "blogtext")
	category := newCategory()

	insertBlogCategory := db.MustPrepare(`INSERT INTO blog_category (blog, category) VALUES(?, ?)`)
	defer insertBlogCategory.Close()

	insertBlogCategory.Exec(blog.ID(), category.ID())

	selectBlogCategory := db.MustPrepare(`SELECT blog, category FROM blog_category WHERE blog = ? AND category = ?`)
	defer selectBlogCategory.Close()

	var (
		blogID     uint64
		categoryID uint64
	)

	t.NoError(selectBlogCategory.QueryRow(blog.ID(), category.ID()).Scan(&blogID, &categoryID))

	t.Is(blog.ID(), blogID)
	t.Is(category.ID(), categoryID)

	t.NoError(repo.RemoveBlogCategory(blog))

	err := selectBlogCategory.QueryRow(blog.ID(), category.ID()).Scan(&blogID, &categoryID)
	t.Is(infra.ErrNoRows, err)
}

func TestRemoveBlogCategory_NoRows(tt *testing.T) {
	t := testing.NewTester(tt)
	repo := NewBlogRepository()

	user, _ := accountFactory.NewUser("username", "userpassword")
	blog, _ := factory.NewBlog(user.ID(), "blogtitle", "blogtext")

	err := repo.RemoveBlogCategory(blog)
	t.Is(infra.ErrNoRows, err)
}

// blog.category.id = 1 if blog.id is even else 2
func createBlogListWithCategory(repo BlogRepository, amount int) {
	db := repo.DB()
	defer db.Close()

	now := time.Now()
	insertBlog := db.MustPrepare(`
		INSERT INTO blog (user, title, text, created_at, updated_at) VALUES(1, ?, ?, ?, ?)
	`)
	defer insertBlog.Close()

	setBlogCategory := db.MustPrepare(`
		INSERT INTO blog_category (blog, category) VALUES(?, ?)
	`)
	defer setBlogCategory.Close()

	for i := 0; i < amount; i++ {
		now = now.Add(1 * time.Second)

		res, _ := insertBlog.Exec(uuid.New(), uuid.New(), now, now)

		blogID, _ := res.LastInsertId()

		setBlogCategory.Exec(blogID, (i%2)+1)
	}
}
