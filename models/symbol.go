package models

type Symbol struct {
	Symbol   string `gorm:"primaryKey"`
	Name     string
	ETF      bool
	Exchange string
}

func ConvertToGenericArray(inp []*Symbol) (res []interface{}) {
	res = make([]interface{}, len(inp))
	for i, v := range inp {
		res[i] = v
	}
	return
}
