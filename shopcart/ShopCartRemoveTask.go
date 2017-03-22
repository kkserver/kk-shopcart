package shopcart

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ShopCartRemoveTaskResult struct {
	app.Result
	Items []ShopCart `json:"items,omitempty"`
}

type ShopCartRemoveTask struct {
	app.Task
	Uid      int64  `json:"uid"`      //用户ID
	Type     string `json:"type"`     //类型
	ShopId   int64  `json:"shopId"`   //店铺ID
	ItemId   int64  `json:"itemId"`   //商品ID
	OptionId int64  `json:"optionId"` //商品规格ID
	Count    int    `json:"count"`    //商品数量
	Result   ShopCartRemoveTaskResult
}

func (task *ShopCartRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *ShopCartRemoveTask) GetInhertType() string {
	return "shopcart"
}

func (task *ShopCartRemoveTask) GetClientName() string {
	return "ShopCart.Remove"
}
