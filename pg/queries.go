package pg

import (
	"github.com/darmiel/macd-api/model"
	"github.com/darmiel/macd-api/yahoo"
)

func (p *Postgres) FindAllSymbols() (res []string, err error) {
	tx := p.Model(&model.Historic{}).
		Distinct("symbol").
		Order("symbol asc").
		Find(&res)

	err = tx.Error
	return
}

func (p *Postgres) FindAllSymbolsWithMinData(num int) (res []string, err error) {
	tx := p.Model(&model.Historic{}).
		Select("symbol").
		Having("count(symbol) >= ?", num).
		Group("symbol").
		Order("symbol asc").
		Find(&res)

	err = tx.Error
	return
}

func (p *Postgres) FindHistoricalsWithMinData(num int) (res map[string][]*model.Historic, err error) {
	h := make([]*model.Historic, 0)
	tx := p.Model(&model.Historic{}).
		Raw("SELECT * FROM historicals h WHERE symbol IN (SELECT symbol FROM historicals GROUP BY symbol HAVING COUNT(symbol) >= ?)", num).
		Find(&h)
	if err = tx.Error; err != nil {
		return
	}
	res = make(map[string][]*model.Historic)
	for _, x := range h {
		if _, o := res[x.Symbol]; !o {
			res[x.Symbol] = make([]*model.Historic, 0)
		}
		res[x.Symbol] = append(res[x.Symbol], x)
	}
	return
}

func (p *Postgres) GetHistoricalData(symbol string, days int) (res []*model.Historic, err error) {
	tx := p.Model(&model.Historic{}).
		Select("*").
		Where("symbol = ?", symbol).
		Order("day_date DESC").
		Limit(days).
		Find(&res)
	err = tx.Error
	return
}

func (p *Postgres) GetHistorical90Data(symbol string) (res model.Quarter, err error) {
	var data []*model.Historic
	if data, err = p.GetHistoricalData(symbol, 90); err != nil {
		return
	}
	return yahoo.Historical90From(data)
}
