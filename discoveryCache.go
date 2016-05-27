package gormid

import (
	"log"

	"github.com/Gacnt/gormid/models"
	"github.com/jinzhu/gorm"
	"github.com/yohcop/openid-go"
)

type discoveryCache struct {
	db *gorm.DB
}

type DiscoveredInfo struct {
	opEndpoint string
	opLocalID  string
	claimedID  string
}

func (s *DiscoveredInfo) OpEndpoint() string {
	return s.opEndpoint
}

func (s *DiscoveredInfo) OpLocalID() string {
	return s.opLocalID
}

func (s *DiscoveredInfo) ClaimedID() string {
	return s.claimedID
}

func (d *discoveryCache) Put(id string, info openid.DiscoveredInfo) {
	dC := &models.DiscoveryCache{id, info.OpEndpoint(), info.OpLocalID(), info.ClaimedID()}
	if err := d.db.Create(&dC).Error; err != nil {
		log.Println(err)
	}
}

func (d *discoveryCache) Get(id string) openid.DiscoveredInfo {
	dC := &models.DiscoveryCache{}
	if err := d.db.Where("cache_id = ?", id).First(&dC).Error; err != nil {
		return nil
	}

	sDI := &DiscoveredInfo{dC.Endpoint, dC.LocalID, dC.ClaimedID}

	return sDI
}
