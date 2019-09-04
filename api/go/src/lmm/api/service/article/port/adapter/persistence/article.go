package persistence

import (
	"sort"
	"time"

	dsUtil "lmm/api/pkg/datastore"
	"lmm/api/pkg/transaction"
	"lmm/api/service/article/domain"
	"lmm/api/service/article/domain/model"

	"cloud.google.com/go/datastore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
)

type ArticleDataStore struct {
	dataStore *datastore.Client
	transaction.Manager
}

func NewArticleDataStore(dataStore *datastore.Client) *ArticleDataStore {
	return &ArticleDataStore{
		dataStore: dataStore,
		Manager:   dsUtil.NewTransactionManager(dataStore),
	}
}

func (s *ArticleDataStore) buildArticleKey(articleID, authorID int64) *datastore.Key {
	userKey := datastore.IDKey(dsUtil.UserKind, authorID, nil)

	return datastore.IDKey(dsUtil.ArticleKind, articleID, userKey)
}

func (s *ArticleDataStore) NextID(tx transaction.Transaction, authorID int64) (*model.ArticleID, error) {
	key := datastore.IncompleteKey(dsUtil.ArticleKind, datastore.IDKey(dsUtil.UserKind, authorID, nil))
	keys, err := s.dataStore.AllocateIDs(tx, []*datastore.Key{key})
	if err != nil || len(keys) == 0 {
		return nil, errors.Wrap(err, "failed to allocate new article key")
	}

	return model.NewArticleID(keys[0].Encode()), nil
}

type article struct {
	Title        string    `datastore:"Title"`
	Body         string    `datastore:"Body,noindex"`
	CreatedAt    time.Time `datastore:"CreatedAt"`
	PublishedAt  time.Time `datastore:"PublishedAt"`
	LastModified time.Time `datastore:"LastModified,noindex"`
}

type tag struct {
	ID    *datastore.Key `datastore:"__key__"`
	Name  string         `datastore:"Name"`
	Order int            `datastore:"Order"`
}

type tagV2 struct {
	tag
	CreatedAt time.Time `datastore:"CreatedAt"`
}

// Save saves article into datastore
func (s *ArticleDataStore) Save(tx transaction.Transaction, model *model.Article) error {
	articleKey := dsUtil.MustKey(model.ID().String())

	dstx := dsUtil.MustTransaction(tx)

	// save article
	if _, err := dstx.Mutate(datastore.NewUpsert(articleKey, &article{
		Title:        model.Content().Text().Title(),
		Body:         model.Content().Text().Body(),
		CreatedAt:    model.CreatedAt(),
		PublishedAt:  model.PublishedAt(),
		LastModified: model.LastModified(),
	})); err != nil {
		return errors.Wrap(err, "failed to put article into datastore")
	}

	// get all tag keys by article
	q := datastore.NewQuery(dsUtil.ArticleTagKind).Ancestor(articleKey).KeysOnly().Transaction(dstx)
	tagKeys, err := s.dataStore.GetAll(tx, q, nil)
	if err != nil {
		return errors.Wrap(err, "failed to get article's tags")
	}

	// delete all tags
	if err := dstx.DeleteMulti(tagKeys); err != nil {
		return errors.Wrap(err, "failed to clear article tags")
	}

	tagKeys = tagKeys[:0]
	tags := make([]*tag, len(model.Content().Tags()), len(model.Content().Tags()))
	for i, model := range model.Content().Tags() {
		tagKeys = append(tagKeys, datastore.IncompleteKey(dsUtil.ArticleTagKind, articleKey))
		tags[i] = &tag{Name: model.Name(), Order: int(model.Order())}
	}

	// save all tags
	if _, err := dstx.PutMulti(tagKeys, tags); err != nil {
		return errors.Wrap(err, "failed to put tags into datastore")
	}

	return nil
}

func (s *ArticleDataStore) FindByID(tx transaction.Transaction, id *model.ArticleID) (*model.Article, error) {
	articleKey, err := datastore.DecodeKey(id.String())
	if err != nil {
		return nil, errors.Wrapf(domain.ErrNoSuchArticle, "%s: %s", err.Error(), id.String())
	}

	dsTx := dsUtil.MustTransaction(tx)
	data := article{}
	if err := dsTx.Get(articleKey, &data); err != nil {
		return nil, errors.Wrap(domain.ErrNoSuchArticle, err.Error())
	}

	fetchTagsQuery := datastore.NewQuery(dsUtil.ArticleTagKind).Ancestor(articleKey).Transaction(dsTx)
	var tags []*tag
	if _, err := s.dataStore.GetAll(tx, fetchTagsQuery, &tags); err != nil {
		return nil, errors.Wrap(err, "failed to get article tags")
	}

	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Order < tags[j].Order
	})

	content, err := model.NewContent(data.Title, data.Body, func() []string {
		ss := make([]string, len(tags), len(tags))
		for i, t := range tags {
			ss[i] = t.Name
		}
		return ss
	}())
	if err != nil {
		return nil, errors.Wrap(err, "internal error")
	}

	author := model.NewAuthor(articleKey.Parent.ID)

	return model.NewArticle(id, author, content, data.CreatedAt, data.PublishedAt, data.LastModified), nil
}

func (s *ArticleDataStore) Remove(tx transaction.Transaction, id *model.ArticleID) error {
	panic("not implemented")
}

type articleItem struct {
	Title       string `datastore:"Title"`
	PublishedAt int64  `datastore:"PublishedAt"`
}

func (s *ArticleDataStore) ViewArticle(tx transaction.Transaction, id string) (*model.Article, error) {
	return s.FindByID(tx, model.NewArticleID(id))
}

func (s *ArticleDataStore) ViewArticles(tx transaction.Transaction, count, page int, filter *model.ArticlesFilter) (*model.ArticleListView, error) {
	if filter != nil && filter.Tag != "" {
		return s.viewArticleFilteredByTag(tx, count, page, filter.Tag)
	}

	return s.viewAllArticles(tx, count, page)
}

func (s *ArticleDataStore) viewAllArticles(tx transaction.Transaction, count, page int) (*model.ArticleListView, error) {
	counting := datastore.NewQuery(dsUtil.ArticleKind)
	paging := datastore.NewQuery(dsUtil.ArticleKind).Project("__key__", "Title", "CreatedAt").Order("-CreatedAt").Limit(count + 1).Offset((page - 1) * count)

	total, err := s.dataStore.Count(tx, counting)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get total number of articles")
	}

	var entities []*articleItem
	keys, err := s.dataStore.GetAll(tx, paging, &entities)
	if err != nil {
		return nil, errors.Wrap(err, "internal error")
	}

	hasNextPage := false
	if len(entities) > int(count) {
		hasNextPage = true
		entities = entities[:int(count)]
	}

	items := make([]*model.ArticleListViewItem, len(entities), len(entities))
	for i, entity := range entities {
		id := model.NewArticleID(keys[i].Encode())
		item, err := model.NewArticleListViewItem(id, entity.Title, time.Unix(entity.PublishedAt/dsUtil.UnixFactor, 0))
		if err != nil {
			return nil, errors.Wrap(err, "internal error")
		}
		items[i] = item
	}

	return model.NewArticleListView(items, page, count, total, hasNextPage), nil
}

func (s *ArticleDataStore) viewArticleFilteredByTag(tx transaction.Transaction, count, page int, tag string) (*model.ArticleListView, error) {
	dstx := dsUtil.MustTransaction(tx)

	counting := datastore.NewQuery(dsUtil.ArticleTagKind).Filter("Name =", tag)
	paging := datastore.NewQuery(dsUtil.ArticleTagKind).Filter("Name =", tag).KeysOnly().Order("-CreatedAt").Limit(count - 1).Offset((page - 1) * count)

	total, err := s.dataStore.Count(tx, counting)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get total number of articles")
	}

	view := model.NewArticleListView([]*model.ArticleListViewItem{}, page, count, total, false)
	if view.Total() == 0 {
		// fallback V1
		return view, nil
	}

	var tags []*tagV2
	keys, err := s.dataStore.GetAll(tx, paging, &tags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count tags")
	}

	var articles []*articleItem
	if err := dstx.GetMulti(keys, &articles); err != nil {
		return nil, errors.Wrap(err, "failed to get articles")
	}

	// select * from TagV2 where name = ? order by CreatedAt desc
	// select * from TagV1 if not found
	// foreach tag in tags select * from Article where id = tag.id.parent
	// foreach tag in tags insert tag into TagV2
	return nil, nil
}

func (s *ArticleDataStore) ViewAllTags(tx transaction.Transaction) ([]*model.TagView, error) {
	q := datastore.NewQuery(dsUtil.ArticleTagKind).Project("Name").Order("Name").Distinct()

	var t tag
	items := make([]*model.TagView, 0)

	iter := s.dataStore.Run(tx, q)
	for {
		_, err := iter.Next(&t)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "internal error: invalid tag")
		}
		items = append(items, model.NewTagView(t.Name))
	}

	return items, nil
}
