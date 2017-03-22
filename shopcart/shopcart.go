package shopcart

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
)

type ShopCart struct {
	Id       int64  `json:"id"`
	Uid      int64  `json:"uid"`      //用户ID
	Type     string `json:"type"`     //类型
	ShopId   int64  `json:"shopId"`   //店铺ID
	ItemId   int64  `json:"itemId"`   //商品ID
	OptionId int64  `json:"optionId"` //商品规格ID
	Count    int    `json:"count"`    //商品数量
	Options  string `json:"options"`  //其他选项
	Ctime    int64  `json:"ctime"`
}

type IShopCartApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetShopCartTable() *kk.DBTable
}

type ShopCartApp struct {
	app.App

	DB *app.DBConfig

	Remote *remote.Service

	ShopCart      *ShopCartService
	ShopCartTable kk.DBTable
}

func (C *ShopCartApp) GetDB() (*sql.DB, error) {
	return C.DB.Get(C)
}

func (C *ShopCartApp) GetPrefix() string {
	return C.DB.Prefix
}

func (C *ShopCartApp) GetShopCartTable() *kk.DBTable {
	return &C.ShopCartTable
}
