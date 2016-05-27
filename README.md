# gormid
Use GORM (Go Object Relational Mapping) to store OpenID DiscoveryCache / Nonce in a database instead of in memory

# Installation
`go get github.com/Gacnt/gormid`

# Usage
```
import "github.com/Gacnt/gormid"

var gormStore = gormid.CreateNewStore(db) // Pass in your *gorm.DB and then access the DiscoveryCache and NonceStore fields of the struct returned e.g.

func AuthCallback(w http.ResponseWriter, r *http.Request) {
	fullURL := "http://localhost:3000" + r.URL.String()
	log.Println(fullURL)
	id, err := openid.Verify(
		fullURL,
		gormStore.DiscoveryCache, gormStore.NonceStore)
	if err == nil {
		log.Println(id)
	} else {
		log.Println(err)
	}
}
```

Docs:

https://godoc.org/github.com/Gacnt/gormid
