package gormid

import (
	"sync"

	"github.com/Schrodinger-Box/gormid/models"
	"github.com/jinzhu/gorm"
)

type gormStore struct {
	DiscoveryCache *discoveryCache
	NonceStore     *nonceStore
}

func CreateNewStore(db *gorm.DB) *gormStore {
	db.AutoMigrate(&models.DiscoveryCache{}, &models.Nonce{})
	return &gormStore{
		DiscoveryCache: &discoveryCache{db},
		NonceStore:     &nonceStore{db, sync.Mutex{}},
	}
}
