package config

import (
	storage_go "github.com/supabase-community/storage-go"
	"os"
)

func SupabaseClient() *storage_go.Client {
	storage := storage_go.NewClient(
		os.Getenv("SUPA_URL"),
		os.Getenv("SUPA_API"),
		nil,
	)

	return storage
}
