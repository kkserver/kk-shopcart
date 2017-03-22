package shopcart

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type ShopCartJoinTaskResult struct {
	app.Result
	Item *ShopCart `json:"item,omitempty"`
}

type ShopCartJoinTask struct {
	app.Task
	Uid      int64       `json:"uid"`      //用户ID
	Type     string      `json:"type"`     //类型
	ShopId   int64       `json:"shopId"`   //店铺ID
	ItemId   int64       `json:"itemId"`   //商品ID
	OptionId int64       `json:"optionId"` //商品规格ID
	Count    int         `json:"count"`    //商品数量
	Options  interface{} `json:"options"`  //其他选项
	Result   ShopCartJoinTaskResult
}

func (task *ShopCartJoinTask) GetResult() interface{} {
	return &task.Result
}

func (task *ShopCartJoinTask) GetInhertType() string {
	return "shopcart"
}

func (task *ShopCartJoinTask) GetClientName() string {
	return "ShopCart.Join"
}
