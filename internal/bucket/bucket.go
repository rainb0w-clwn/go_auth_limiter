package bucket

import "time"

type IBucket interface {
	GetSize() int
	GetTokenCount() int
	GetLastRefill() time.Time
	GetToken(int)
	Refill()
	Reset()
	Full() bool
}
