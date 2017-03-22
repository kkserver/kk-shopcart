package shopcart

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ShopCartQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
	RowCount  int `json:"rowCount"`
}

type ShopCartQueryTaskResult struct {
	app.Result
	Counter *ShopCartQueryCounter `json:"counter,omitempty"`
	Items   []ShopCart            `json:"items,omitempty"`
}

type ShopCartQueryTask struct {
	app.Task
	Uid       int64  `json:"uid"`
	Type      string `json:"type"`     //类型
	ShopId    int64  `json:"shopId"`   //店铺ID
	ItemId    int64  `json:"itemId"`   //商品ID
	OptionId  int64  `json:"optionId"` //商品规格ID
	OrderBy   string `json:"orderBy"`  //desc asc shopId uid
	PageIndex int    `json:"p"`
	PageSize  int    `json:"size"`
	Counter   bool   `json:"counter"`
	Result    ShopCartQueryTaskResult
}

func (task *ShopCartQueryTask) GetResult() interface{} {
	return &task.Result
}

func (task *ShopCartQueryTask) GetInhertType() string {
	return "shopcart"
}

func (task *ShopCartQueryTask) GetClientName() string {
	return "ShopCart.Query"
}
