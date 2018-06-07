package ids

import (
	"encoding/base64"
	"strings"

	"github.com/pborman/uuid"
)

// NewID Risk it for a biscuit
func NewID() string {
	uuid := uuid.New()
	return strings.ToLower(base64.StdEncoding.EncodeToString([]byte(uuid))[:9])
}
