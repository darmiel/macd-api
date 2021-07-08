package models

import (
	"time"
)

type Historical90 [90]*Historical

type Historical struct {
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
