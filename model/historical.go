package model

import (
	"time"
)

type Quarter [90]*Historic

type Historic struct {
	Symbol      string    `gorm:"primaryKey"`
	DayDate     time.Time `gorm:"primaryKey"`
	OrigDate    time.Time
	High        float32
	Low         float32
	Volume      int
	Close       float32
	Open        float32
	SymbolModel *Symbol `gorm:"foreignKey:Symbol"`
}
