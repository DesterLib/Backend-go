package cache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/anonyindian/logger"
)

var Cache *bigcache.BigCache

func Load(l *logger.Logger) {
	log := l.Create("CACHE")
	defer log.ChangeLevel(logger.LevelMain).Println("LOADED")
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         10 * 24 * time.Hour,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       500,
		HardMaxCacheSize:   512,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}
	var err error
	Cache, err = bigcache.NewBigCache(config)
	if err != nil {
		log.ChangeLevel(logger.LevelError).Println(err.Error())
	}
}
