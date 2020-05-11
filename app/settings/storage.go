package settings

import (
	"github.com/Syncano/orion/pkg/storage"
)

const (
	BucketData    = storage.BucketKey("data")
	BucketHosting = storage.BucketKey("hosting")
)

var (
	Buckets = map[storage.BucketKey]string{
		BucketData:    "STORAGE_BUCKET",
		BucketHosting: "STORAGE_HOSTING_BUCKET",
	}
)
