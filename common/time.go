package common

import (
	"time"
)

func NormalizeTime(t time.Time) time.Time {
	res := t.Unix()
	res = res - res%86400
	return time.Unix(res, 0)
}

func NormalizeTimeNoon(t time.Time) time.Time {
	nm := NormalizeTime(t)
	return nm.Add(13*time.Hour + 37*time.Minute)
}

func IsToday(t time.Time) bool {
	const f = "02012006"
	return time.Now().Format(f) == t.Format(f)
}
