package refillrate

import "time"

type RefillRate struct {
	count int
	time  time.Duration
}

func New(count int, time time.Duration) RefillRate {
	return RefillRate{
		count: count,
		time:  time,
	}
}

func (r RefillRate) GetTime() time.Duration {
	return r.time
}

func (r RefillRate) GetCount() int {
	return r.count
}
