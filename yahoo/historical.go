package yahoo

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Historical90 [90]*Historical

type Historical struct {
	gorm.Model
	Symbol string    `gorm:"primaryKey"`
	Date   time.Time `gorm:"primaryKey"`
	High   float32
	Low    float32
	Volume int
	Close  float32
	Open   float32
}

var (
	ErrInvalidResponse = errors.New("invalid response")
	ErrHighEmpty       = errors.New("highs empty")
	ErrLowEmpty        = errors.New("lows empty")
	ErrVolumeEmpty     = errors.New("volumes empty")
	ErrCloseEmpty      = errors.New("closes empty")
	ErrOpenEmpty       = errors.New("opens empty")
)
