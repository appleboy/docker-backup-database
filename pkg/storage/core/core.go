package core

import "time"

// SignedURLOptions download options
type SignedURLOptions struct {
	Expiry          time.Duration
	DefaultFilename string
}
