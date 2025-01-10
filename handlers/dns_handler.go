package handlers

import (
	"database/sql"
	"encoding/json"
	"net"
	"net/http"

	"alifattahi.ir/go-dns-resolver/metrics"
	"alifattahi.ir/go-dns-resolver/models"
)

type Response struct {
	Domain      string `json:"domain"`
	IP          string `json:"ip"`
	DNSProvider string `json:"dns_provider"`
	Cached      bool   `json:"cached"`
}

func ResolveDomainHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := r.URL.Query().Get("domain")
		if domain == "" {
			http.Error(w, "Domain parameter is required", http.StatusBadRequest)
			return
		}
		// Check if domain exists in cache
		cachedData, err := models.GetCache(db, domain)
		if err == nil && cachedData != nil {
			metrics.DNSCacheHits.Inc()
			// If cache found, set 'Cached' flag to true and return cached response
			response := Response{
				Domain:      cachedData.Domain,
				IP:          cachedData.IP,
				DNSProvider: cachedData.DNSProvider,
				Cached:      true,
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		metrics.DNSCacheMisses.Inc()

		// If cache not found, resolve IP and DNS provider
		ipRecords, err := net.LookupIP(domain)
		if err != nil {
			http.Error(w, "Failed to resolve domain", http.StatusInternalServerError)
			return
		}

		dnsRecords, err := net.LookupNS(domain)
		if err != nil {
			metrics.FailedDNSResolver.Inc()
			http.Error(w, "Failed to resolve DNS provider", http.StatusInternalServerError)
			return
		}

		// Get the first IP and DNS provider from the lists
		ip := ipRecords[0].String()
		dnsProvider := dnsRecords[0].Host

		// Prepare the response
		response := Response{
			Domain:      domain,
			IP:          ip,
			DNSProvider: dnsProvider,
			Cached:      false,
		}

		// Save data to cache asynchronously
		cacheData := models.Cache{
			Domain:      domain,
			IP:          ip,
			DNSProvider: dnsProvider,
		}
		go models.SaveCache(db, cacheData)

		// Send response
		json.NewEncoder(w).Encode(response)
	}
}
