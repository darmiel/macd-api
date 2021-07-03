package pg

import "github.com/darmiel/macd-api/yahoo"

func (p *Postgres) FindAllSymbols() (res []string, err error) {
	tx := p.Model(&yahoo.Historical{}).
		Distinct("symbol").
		Order("symbol asc").
		Find(&res)

	err = tx.Error
	return
}

func (p *Postgres) FindAllSymbolsWithMinData(num int) (res []string, err error) {
	tx := p.Model(&yahoo.Historical{}).
		Select("symbol").
		Having("count(symbol) >= ?", num).
		Group("symbol").
		Order("symbol asc").
		Find(&res)

	err = tx.Error
	return
}

func (p *Postgres) FindHistoricalsWithMinData(num int) (res map[string][]*yahoo.Historical, err error) {
	h := make([]*yahoo.Historical, 0)
	tx := p.Model(&yahoo.Historical{}).
		Raw("SELECT * FROM historicals h WHERE symbol IN (SELECT symbol FROM historicals GROUP BY symbol HAVING COUNT(symbol) >= ?)", num).
		Find(&h)
	if err = tx.Error; err != nil {
		return
	}
	res = make(map[string][]*yahoo.Historical)
	for _, x := range h {
		if _, o := res[x.Symbol]; !o {
			res[x.Symbol] = make([]*yahoo.Historical, 0)
		}
		res[x.Symbol] = append(res[x.Symbol], x)
	}
	return
}

func (p *Postgres) GetHistoricalData(symbol string) (res []*yahoo.Historical, err error) {
	tx := p.Model(&yahoo.Historical{}).
		Select("*").
		Where("symbol = ?", symbol).
		Order("date ASC").
		Find(&res)
	err = tx.Error
	return
}
