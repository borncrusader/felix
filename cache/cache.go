package cache

import (
	"container/list"
	"log"
	"os"
)

/* total capacity of the cache */
const cacheCapacity int64 = 67108864

/* total size of the cache */
var cacheSize int64 = 0

/* map of all files in the cache */
var cacheFiles map[string][]byte = make(map[string][]byte)

/* list of fileinfo objs in MRU order */
var cacheList list.List

func cacheAdd(fi os.FileInfo, buf []byte) {
	cacheSize += fi.Size()
	cacheFiles[fi.Name()] = buf
	cacheList.PushFront(fi)

	log.Printf("Cache Add: '%s'; Size: %d; CacheSize: %d\n", fi.Name(),
		fi.Size(), cacheSize)
}

func cacheRetrieve(fi os.FileInfo) ([]byte, bool) {
	stream, ok := cacheFiles[fi.Name()]

	if !ok {
		/* cache miss */
		log.Printf("Cache Miss: '%s'; Size: %d; CacheSize: %d\n", fi.Name(),
			fi.Size(), cacheSize)
		return stream, ok
	}

	/* cache hit, but let's check if the file hasn't changed yet */
	var e *list.Element

	for e = cacheList.Front(); e != nil; e = e.Next() {
		if e.Value.(os.FileInfo).Name() == fi.Name() {
			break
		}
	}

	fi_cache := e.Value.(os.FileInfo)

	/* if all these match, perhaps it's the same file */
	if fi_cache.ModTime().Equal(fi.ModTime()) &&
		fi_cache.Size() == fi.Size() &&
		fi_cache.Mode() == fi.Mode() {
		log.Printf("Cache Hit: '%s'; Size: %d; cacheSize: %d\n", fi.Name(),
			fi.Size(), cacheSize)
		cacheList.MoveToFront(e)
		return stream, ok
	}

	/* the file has changed under us! not reliable to have it in cache */
	log.Printf("Cache Evict: '%s'; Size: %d; Reason: File changed; CacheSize: %d\n",
		fi.Name(), fi.Size(), cacheSize)

	delete(cacheFiles, fi.Name())
	cacheList.Remove(e)

	return nil, false
}

func cacheEvictLRU(fi os.FileInfo) {
	e := cacheList.Back()

	for e != nil {
		if (cacheCapacity - cacheSize) >= fi.Size() {
			break
		}

		prev := e.Prev()

		fi_cache := e.Value.(os.FileInfo)

		cacheSize -= fi_cache.Size()

		log.Printf("Cache Evict: '%s'; Size: %d; Reason: LRU; CacheSize: %d\n",
			fi_cache.Name(), fi_cache.Size(), cacheSize)

		delete(cacheFiles, fi_cache.Name())
		cacheList.Remove(e)

		e = prev
	}
}

func RetrieveFile(fi os.FileInfo) ([]byte, error) {
	/* we cannot handle more than the cache limit */
	if fi.Size() > cacheCapacity {
		log.Printf("Cannot add to cache: '%s'; Size: %d; CacheSize: %d\n",
			fi.Name(), fi.Size(), cacheSize)
		return nil, nil
	}

	stream, ok := cacheRetrieve(fi)
	if ok {
		/* cache hit */
		return stream, nil
	}

	file, err := os.Open(fi.Name())
	if err != nil {
		/* caller will log */
		return nil, err
	}

	buf := make([]byte, fi.Size())
	_, err = file.Read(buf[0:])

	file.Close()

	if (cacheSize + fi.Size()) > cacheCapacity {
		/* need to evict something */
		cacheEvictLRU(fi)
	}

	/* can add it to cache */
	cacheAdd(fi, buf)

	return buf, nil
}
