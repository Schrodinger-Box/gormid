package gormid

import (
	"sync"

	"gorm.io/gorm"

	"github.com/Schrodinger-Box/gormid/models"
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
