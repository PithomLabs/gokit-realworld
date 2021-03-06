package inmem

import (
	realworld "github.com/xesina/gokit-realworld"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

func NewMemArticleRepo() realworld.ArticleRepo {
	return &memArticleRepo{
		m: map[string]realworld.Article{},
	}
}

type memArticleRepo struct {
	rwlock  sync.RWMutex
	m       map[string]realworld.Article
	counter int64
}

func (store *memArticleRepo) Create(a realworld.Article) (*realworld.Article, error) {
	store.rwlock.Lock()
	defer store.rwlock.Unlock()

	if _, ok := store.m[a.Slug]; ok {
		return nil, realworld.ErrArticleAlreadyExists
	}

	a.ID = atomic.AddInt64(&store.counter, 1)
	a.Favorites = make(realworld.Favorites, 0)
	a.Comments = make(realworld.Comments, 0)
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	store.m[a.Slug] = a
	return &a, nil
}

func (store *memArticleRepo) Update(slug string, a realworld.Article) (*realworld.Article, error) {
	store.rwlock.Lock()
	defer store.rwlock.Unlock()

	old, ok := store.m[slug]
	if !ok {
		return nil, realworld.ErrArticleNotFound
	}

	a.ID = old.ID
	a.Comments = old.Comments
	a.Favorites = old.Favorites
	a.CreatedAt = old.CreatedAt
	a.UpdatedAt = time.Now()

	store.m[a.Slug] = a

	if old.Slug != a.Slug {
		delete(store.m, old.Slug)
	}

	return &a, nil
}

func (store *memArticleRepo) Delete(a realworld.Article) error {
	store.rwlock.Lock()
	defer store.rwlock.Unlock()
	delete(store.m, a.Slug)
	return nil
}

func (store *memArticleRepo) Get(slug string) (*realworld.Article, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()

	article, ok := store.m[slug]

	if !ok {
		return nil, realworld.ErrArticleNotFound
	}

	return &article, nil
}

func (store *memArticleRepo) List(offset, limit int) ([]*realworld.Article, int, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()

	if len(store.m) == 0 {
		return []*realworld.Article{}, 0, nil
	}

	count := len(store.m)
	orderedByIDs := make([]*realworld.Article, count)
	ids := make([]int, 0)

	for _, v := range store.m {
		ids = append(ids, int(v.ID))
	}

	sort.Ints(ids)
	for i, id := range ids {
		for k, a := range store.m {
			if a.ID != int64(id) {
				continue
			}
			a := store.m[k]
			orderedByIDs[i] = &a
			break
		}
	}

	var limited []*realworld.Article
	for i := offset; i < limit; i++ {
		if i >= len(orderedByIDs) {
			break
		}
		limited = append(limited, orderedByIDs[i])
	}

	return limited, len(limited), nil
}

func (store *memArticleRepo) ListByTag(tag string, offset, limit int) ([]*realworld.Article, int, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()

	if len(store.m) == 0 {
		return []*realworld.Article{}, 0, nil
	}

	count := len(store.m)
	orderedByIDs := make([]*realworld.Article, count)
	ids := make([]int, 0)

	for _, v := range store.m {
		ids = append(ids, int(v.ID))
	}

	sort.Ints(ids)
	for i, id := range ids {
		for k, a := range store.m {
			if a.ID != int64(id) {
				continue
			}
			a := store.m[k]
			orderedByIDs[i] = &a
			break
		}
	}

	qualified := store.filterByTag(tag, orderedByIDs)

	var limited []*realworld.Article
	for i := offset; i < limit; i++ {
		if i >= len(qualified) {
			break
		}
		limited = append(limited, qualified[i])
	}

	return limited, len(limited), nil
}

func (store *memArticleRepo) ListByAuthorID(id int64, offset, limit int) ([]*realworld.Article, int, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()

	if len(store.m) == 0 {
		return []*realworld.Article{}, 0, nil
	}

	count := len(store.m)
	orderedByIDs := make([]*realworld.Article, count)
	ids := make([]int, 0)

	for _, v := range store.m {
		ids = append(ids, int(v.ID))
	}

	sort.Ints(ids)
	for i, id := range ids {
		for k, a := range store.m {
			if a.ID != int64(id) {
				continue
			}
			a := store.m[k]
			orderedByIDs[i] = &a
			break
		}
	}

	qualified := store.filterByAuthorID(id, orderedByIDs)

	var limited []*realworld.Article
	for i := offset; i < limit; i++ {
		if i >= len(qualified) {
			break
		}
		limited = append(limited, qualified[i])
	}

	return limited, len(limited), nil
}

func (store *memArticleRepo) ListByFavoriterID(id int64, offset, limit int) ([]*realworld.Article, int, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()

	if len(store.m) == 0 {
		return []*realworld.Article{}, 0, nil
	}

	count := len(store.m)
	orderedByIDs := make([]*realworld.Article, count)
	ids := make([]int, 0)

	for _, v := range store.m {
		ids = append(ids, int(v.ID))
	}

	sort.Ints(ids)
	for i, id := range ids {
		for k, a := range store.m {
			if a.ID != int64(id) {
				continue
			}
			a := store.m[k]
			orderedByIDs[i] = &a
			break
		}
	}

	qualified := store.filterByAuthorID(id, orderedByIDs)

	var limited []*realworld.Article
	for i := offset; i < limit; i++ {
		if i >= len(qualified) {
			break
		}
		limited = append(limited, qualified[i])
	}

	return limited, len(limited), nil
}

func (store *memArticleRepo) Feed(req realworld.FeedRequest) ([]*realworld.Article, int, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()

	if len(store.m) == 0 {
		return []*realworld.Article{}, 0, nil
	}

	count := len(store.m)
	orderedByIDs := make([]*realworld.Article, count)
	ids := make([]int, 0)

	for _, v := range store.m {
		ids = append(ids, int(v.ID))
	}

	sort.Ints(ids)
	for i, id := range ids {
		for k, a := range store.m {
			if a.ID != int64(id) {
				continue
			}
			a := store.m[k]
			orderedByIDs[i] = &a
			break
		}
	}

	qualified := make([]*realworld.Article, 0)
	for _, authorID := range req.FollowingIDs {
		qualified = append(qualified, store.filterByAuthorID(authorID, orderedByIDs)...)
	}

	offset, limit := 0, 20

	if req.Offset > 0 {
		offset = req.Limit
	}
	if req.Limit > 0 {
		limit = req.Limit
	}

	var limited []*realworld.Article
	for i := offset; i < limit; i++ {
		if i >= len(qualified) {
			break
		}
		limited = append(limited, qualified[i])
	}

	return limited, len(limited), nil
}

func (store *memArticleRepo) filterByTag(tag string, articles []*realworld.Article) []*realworld.Article {
	if tag == "" {
		return articles
	}

	var qualified []*realworld.Article
	for _, article := range articles {
		if article.Tags.HasTag(tag) {
			qualified = append(qualified, article)
		}
	}

	return qualified
}

func (store *memArticleRepo) filterByAuthorID(id int64, articles []*realworld.Article) []*realworld.Article {
	if id == 0 {
		return articles
	}

	var qualified []*realworld.Article
	for _, article := range articles {
		if article.Author.ID == id {
			qualified = append(qualified, article)
		}
	}

	return qualified
}

func (store *memArticleRepo) filterByFavotiterID(id int64, articles []*realworld.Article) []*realworld.Article {
	if id == 0 {
		return articles
	}

	var qualified []*realworld.Article
	for _, article := range articles {
		if article.Favorites.FavoritedBy(id) {
			qualified = append(qualified, article)
		}
	}

	return qualified
}

func (store *memArticleRepo) AddFavorite(a realworld.Article, u realworld.User) (*realworld.Article, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()

	article, ok := store.m[a.Slug]
	if !ok {
		return nil, realworld.ErrArticleNotFound
	}

	article.Favorites[u.ID] = struct{}{}

	return &article, nil
}

func (store *memArticleRepo) RemoveFavorite(a realworld.Article, u realworld.User) (*realworld.Article, error) {
	store.rwlock.RLock()
	defer store.rwlock.RUnlock()

	article, ok := store.m[a.Slug]
	if !ok {
		return nil, realworld.ErrArticleNotFound
	}

	delete(article.Favorites, u.ID)

	return &article, nil
}

func (store *memArticleRepo) Tags() ([]*realworld.Tag, error) {
	tags := make(map[string]struct{})
	tt := make([]*realworld.Tag, 0)

	for _, a := range store.m {
		for _, tag := range a.Tags {
			if _, ok := tags[tag.Tag]; ok {
				continue
			}

			tags[tag.Tag] = struct{}{}
			t := tag
			tt = append(tt, &t)
		}
	}

	return tt, nil
}
