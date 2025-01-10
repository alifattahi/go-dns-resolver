package models

import (
	"database/sql"
	"time"
	"errors"
)

type Cache struct {
	Domain      string
	IP          string
	DNSProvider string
	CreatedAt   time.Time
}



func GetCache(db *sql.DB, domain string) (*Cache, error) {
	row := db.QueryRow("SELECT domain, ip, dns_provider, created_at FROM dns_cache WHERE domain = $1", domain)

	var cache Cache
	err := row.Scan(&cache.Domain, &cache.IP, &cache.DNSProvider, &cache.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &cache, nil
}

func SaveCache(db *sql.DB, data interface{}) error {
	// Assert the type of data to Cache
	cacheData, ok := data.(Cache)
	if !ok {
		return errors.New("invalid data type: expected Cache")
	}

	// Now you can access cacheData.Domain, cacheData.IP, cacheData.DNSProvider, and cacheData.CreatedAt
	_, err := db.Exec("INSERT INTO dns_cache (domain, ip, dns_provider, created_at) VALUES ($1, $2, $3, $4)",
		cacheData.Domain, cacheData.IP, cacheData.DNSProvider, time.Now())
	return err
}
