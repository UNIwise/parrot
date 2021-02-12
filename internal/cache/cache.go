package cache

type Cache interface {
	GetTranslation(projectID uint32, languageCode, format string)
}

type CacheImpl struct{}

func New() *CacheImpl {
	return &CacheImpl{}
}
