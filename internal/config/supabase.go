package config

import (
	"os"
	"strconv"
)

type SupabaseConfig struct {
	URL        string
	ServiceKey string
	Bucket     string
	Expire     int
	MaxSizeMB  int
}

func GetSupabaseConfig() SupabaseConfig {
	expire := 600
	if v := os.Getenv("STORAGE_SIGNED_URL_EXPIRES"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			expire = parsed
		}
	}

	maxSize := 2
	if v := os.Getenv("MAX_UPLOAD_SIZE_MB"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			maxSize = parsed
		}
	}

	return SupabaseConfig{
		URL:        os.Getenv("SUPABASE_URL"),
		ServiceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
		Bucket:     os.Getenv("SUPABASE_BUCKET_NAME"),
		Expire:     expire,
		MaxSizeMB:  maxSize,
	}
}
