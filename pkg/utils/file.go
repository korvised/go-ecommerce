package utils

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

func RandFileName(ext string) string {
	filename := fmt.Sprintf(
		"%s%v",
		strings.ReplaceAll(uuid.NewString()[:6], "-", ""),
		time.Now().UnixMilli(),
	)

	if ext != "" {
		filename += fmt.Sprintf(".%s", ext)
	}

	return filename
}
