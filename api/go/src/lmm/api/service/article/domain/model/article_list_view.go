package model

// ArticleListView is the model used to view article list
type ArticleListView struct {
	items       []*ArticleListViewItem
	page        int
	perPage     int
	tagFilter   string
	total       int
	hasNextPage bool
}

// NewArticleListView constructs a new ArticleListView
func NewArticleListView(items []*ArticleListViewItem, tagFilter string, page, perPage, total int, hasNextPage bool) *ArticleListView {
	return &ArticleListView{
		items:       items,
		page:        page,
		perPage:     perPage,
		tagFilter:   tagFilter,
		total:       total,
		hasNextPage: hasNextPage,
	}
}

// Items gets items of article list view
func (v *ArticleListView) Items() []*ArticleListViewItem {
	return v.items
}

func (v *ArticleListView) Page() int {
	return v.page
}

func (v *ArticleListView) PerPage() int {
	return v.perPage
}

func (v *ArticleListView) TagFilter() string {
	return v.tagFilter
}

// Total returns total articles number
func (v *ArticleListView) Total() int {
	return v.total
}

// HasNextPage returns true if there is next page
func (v *ArticleListView) HasNextPage() bool {
	return v.hasNextPage
}
