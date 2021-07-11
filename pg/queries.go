package pg

import (
	"github.com/darmiel/macd-api/model"
	"github.com/darmiel/macd-api/yahoo"
)

func (p *Postgres) FindAllSymbols() (res []string, err error) {
	tx := p.Model(&model.Historical{}).
		Distinct("symbol").
		Order("symbol asc").
		Find(&res)

	err = tx.Error
	return
}

func (p *Postgres) FindAllSymbolsWithMinData(num int) (res []string, err error) {
	tx := p.Model(&model.Historical{}).
		Select("symbol").
		Having("count(symbol) >= ?", num).
		Group("symbol").
		Order("symbol asc").
		Find(&res)

	err = tx.Error
	return
}

func (p *Postgres) FindHistoricalsWithMinData(num int) (res map[string][]*model.Historical, err error) {
	h := make([]*model.Historical, 0)
	tx := p.Model(&model.Historical{}).
		Raw("SELECT * FROM historicals h WHERE symbol IN (SELECT symbol FROM historicals GROUP BY symbol HAVING COUNT(symbol) >= ?)", num).
		Find(&h)
	if err = tx.Error; err != nil {
		return
	}
	res = make(map[string][]*model.Historical)
	for _, x := range h {
		if _, o := res[x.Symbol]; !o {
			res[x.Symbol] = make([]*model.Historical, 0)
		}
		res[x.Symbol] = append(res[x.Symbol], x)
	}
	return
}

func (p *Postgres) GetHistoricalData(symbol string, days int) (res []*model.Historical, err error) {
	tx := p.Model(&model.Historical{}).
		Select("*").
		Where("symbol = ?", symbol).
		Order("date DESC").
		Limit(days).
		Find(&res)
	err = tx.Error
	return
}

func (p *Postgres) GetHistorical90Data(symbol string) (res model.Historical90, err error) {
	var data []*model.Historical
	if data, err = p.GetHistoricalData(symbol, 90); err != nil {
		return
	}
	return yahoo.Historical90From(data)
}
