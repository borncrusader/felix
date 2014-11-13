package cache

import (
	"os"
)

type fileCache struct {
	size  int64
	files map[string][]byte
}

func (c *fileCache) evictCache(size int64) {

}

var cacheLimit int64 = 67108864
var cache fileCache

func RetrieveFile(fi os.FileInfo) ([]byte, bool) {
	/* we cannot handle more than the cache limit */
	return nil, false

	if fi.Size() > cacheLimit {
		return nil, false
	}

	/* check if already in cache */
	_, ok := cache.files[fi.Name()]
	if !ok {
		/* add to cache */
	}

	return nil, true
}
