package api

import (
	"github.com/google/uuid"
)

// GenerateShortUUID 生成短UUID用于ID
func GenerateShortUUID() string {
	return uuid.NewString()[:8]
}
