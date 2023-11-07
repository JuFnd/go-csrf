package csrf

import "time"

type Csrf struct {
	SID       string
	ExpiresAt time.Time
}
