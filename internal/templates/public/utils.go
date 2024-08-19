package templates

import (
	"time"
)

func currYear() string {
	return time.Now().Format("2006")
}
