package gormid

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/Schrodinger-Box/gormid/models"
)

var maxNonceAge = flag.Duration("max-nonce-age",
	60*time.Second,
	"Maximum accepted age for openid nonces. The bigger, the more"+
		"memory is needed to store used nonces.")

type nonceStore struct {
	db *gorm.DB
	mu sync.Mutex
}

func (d *nonceStore) Accept(endpoint, nonce string) error {
	if len(nonce) < 20 || len(nonce) > 256 {
		return errors.New("Invalid nonce")
	}

	ts, err := time.Parse(time.RFC3339, nonce[0:20])
	if err != nil {
		return err
	}

	now := time.Now()
	diff := now.Sub(ts)
	if diff > *maxNonceAge {
		return fmt.Errorf("Nonce too old: %ds", diff.Seconds())
	}

	s := nonce[20:]

	d.mu.Lock()
	defer d.mu.Unlock()
	nonces := []*models.Nonce{}
	if err := d.db.Where("endpoint = ?", endpoint).Find(&nonces).Error; err != nil {
		return err
	}

	newNonces := []*models.Nonce{{ts, s, endpoint}}
	if len(nonces) > 0 {
		d.db.Delete(new(models.Nonce))
		for _, n := range nonces {
			if n.T.UTC() == ts && n.S == s {
				return errors.New("Nonce already used")
			}
			if now.Sub(n.T) < *maxNonceAge {
				newNonces = append(newNonces, n)
			}
		}
		if ok := storeNonces(d.db, newNonces); !ok {
			return errors.New("Could not store nonces")
		}
	} else {
		if ok := storeNonces(d.db, newNonces); !ok {
			return errors.New("Could not store nonces")
		}
	}

	return nil
}

func storeNonces(db *gorm.DB, nonces []*models.Nonce) bool {

	var wg sync.WaitGroup
	for _, n := range nonces {
		wg.Add(1)
		go func(nonce *models.Nonce) {
			if err := db.Create(&nonce).Error; err != nil {
				log.Println(err)
			}
			wg.Done()
		}(n)
	}

	wg.Wait()

	return true
}
