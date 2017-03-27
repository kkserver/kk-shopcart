package shopcart

import (
	"bytes"
	"fmt"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/json"
	"strings"
	"time"
)

type ShopCartService struct {
	app.Service

	Join   *ShopCartJoinTask
	Remove *ShopCartRemoveTask
	Query  *ShopCartQueryTask
}

func (S *ShopCartService) Handle(a app.IApp, task app.ITask) error {
	return app.ServiceReflectHandle(a, task, S)
}

func (S *ShopCartService) HandleShopCartJoinTask(a IShopCartApp, task *ShopCartJoinTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SHOPCART
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.Uid == 0 {
		task.Result.Errno = ERROR_SHOPCART_NOT_FOUND_UID
		task.Result.Errmsg = "未找到用户ID"
		return nil
	}

	rs, err := kk.DBQuery(db, a.GetShopCartTable(), a.GetPrefix(), " WHERE uid=? AND shopid=? AND itemid=? AND optionid=? AND type=?", task.Uid, task.ShopId, task.ItemId, task.OptionId, task.Type)

	if err != nil {
		task.Result.Errno = ERROR_SHOPCART
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rs.Close()

	v := ShopCart{}

	if rs.Next() {

		scanner := kk.NewDBScaner(&v)

		err = scanner.Scan(rs)

		if err != nil {
			task.Result.Errno = ERROR_SHOPCART
			task.Result.Errmsg = err.Error()
			return nil
		}

	} else {
		v.Ctime = time.Now().Unix()
	}

	v.Uid = task.Uid
	v.ItemId = task.ItemId
	v.ShopId = task.ShopId
	v.OptionId = task.OptionId
	v.Count = v.Count + task.Count
	v.Type = task.Type

	if task.Options != nil {
		b, err := json.Encode(task.Options)
		if err != nil {
			task.Result.Errno = ERROR_SHOPCART
			task.Result.Errmsg = err.Error()
			return nil
		}
		v.Options = string(b)
	}

	if v.Id == 0 {
		_, err = kk.DBInsert(db, a.GetShopCartTable(), a.GetPrefix(), &v)

		if err != nil {
			task.Result.Errno = ERROR_SHOPCART
			task.Result.Errmsg = err.Error()
			return nil
		}
	} else {

		_, err = kk.DBUpdate(db, a.GetShopCartTable(), a.GetPrefix(), &v)

		if err != nil {
			task.Result.Errno = ERROR_SHOPCART
			task.Result.Errmsg = err.Error()
			return nil
		}
	}

	task.Result.Item = &v

	return nil
}

func (S *ShopCartService) HandleShopCartRemoveTask(a IShopCartApp, task *ShopCartRemoveTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SHOPCART
		task.Result.Errmsg = err.Error()
		return nil
	}

	if task.Uid == 0 {
		task.Result.Errno = ERROR_SHOPCART_NOT_FOUND_UID
		task.Result.Errmsg = "未找到用户ID"
		return nil
	}

	sql := bytes.NewBuffer(nil)

	args := []interface{}{}

	sql.WriteString(" WHERE uid=?")

	args = append(args, task.Uid)

	if task.Type != "" {
		sql.WriteString(" AND type=?")
		args = append(args, task.Type)
	}

	if task.ShopId != 0 {
		sql.WriteString(" AND shopid=?")
		args = append(args, task.ShopId)
	}

	if task.ItemId != 0 {
		sql.WriteString(" AND itemid=?")
		args = append(args, task.ItemId)
	}

	if task.OptionId != 0 {
		sql.WriteString(" AND optionid=?")
		args = append(args, task.OptionId)
	}

	if task.Ids != "" {

		sql.WriteString(" AND id IN (")

		for i, s := range strings.Split(task.Ids, ",") {
			if i != 0 {
				sql.WriteString(",")
			}
			sql.WriteString("?")
			args = append(args, s)
		}

		sql.WriteString(")")

	}

	v := ShopCart{}

	items := []ShopCart{}

	tx, err := db.Begin()

	err = func() error {

		rows, err := kk.DBQuery(tx, a.GetShopCartTable(), a.GetPrefix(), sql.String(), args...)

		if err != nil {
			return err
		}

		for rows.Next() {

			scanner := kk.NewDBScaner(&v)

			err = scanner.Scan(rows)

			if err != nil {
				rows.Close()
				return err
			}

			v.Count = v.Count - task.Count

			items = append(items, v)

		}

		rows.Close()

		if len(items) == 0 {
			return app.NewError(ERROR_SHOPCART_NOT_FOUND, "未找到商品")
		}

		for _, v = range items {

			if v.Count > 0 {

				_, err = kk.DBUpdateWithKeys(tx, a.GetShopCartTable(), a.GetPrefix(), &v, map[string]bool{"count": true})

				if err != nil {
					return err
				}
			} else {

				_, err = kk.DBDelete(tx, a.GetShopCartTable(), a.GetPrefix(), " WHERE id=?", v.Id)

				if err != nil {
					return err
				}
			}
		}

		return nil
	}()

	if err == nil {
		err = tx.Commit()
	}

	if err != nil {
		e, ok := err.(*app.Error)
		if ok {
			task.Result.Errno = e.Errno
			task.Result.Errmsg = e.Errmsg
			return nil
		}
		task.Result.Errno = ERROR_SHOPCART
		task.Result.Errmsg = err.Error()
		return nil
	}

	task.Result.Items = items

	return nil
}

func (S *ShopCartService) HandleShopCartQueryTask(a IShopCartApp, task *ShopCartQueryTask) error {

	var db, err = a.GetDB()

	if err != nil {
		task.Result.Errno = ERROR_SHOPCART
		task.Result.Errmsg = err.Error()
		return nil
	}

	var items = []ShopCart{}

	var args = []interface{}{}

	var sql = bytes.NewBuffer(nil)

	sql.WriteString(" WHERE 1")

	if task.Type != "" {
		sql.WriteString(" AND type=?")
		args = append(args, task.Type)
	}

	if task.Uid != 0 {
		sql.WriteString(" AND uid=?")
		args = append(args, task.Uid)
	}

	if task.ItemId != 0 {
		sql.WriteString(" AND itemid=?")
		args = append(args, task.ItemId)
	}

	if task.ShopId != 0 {
		sql.WriteString(" AND shopid=?")
		args = append(args, task.ShopId)
	}

	if task.OptionId != 0 {
		sql.WriteString(" AND optionid=?")
		args = append(args, task.OptionId)
	}

	if task.OrderBy == "asc" {
		sql.WriteString(" ORDER BY id ASC")
	} else if task.OrderBy == "shopId" {
		sql.WriteString(" ORDER BY shopId ASC, id ASC")
	} else if task.OrderBy == "uid" {
		sql.WriteString(" ORDER BY uid ASC, id ASC")
	} else {
		sql.WriteString(" ORDER BY id DESC")
	}

	var pageIndex = task.PageIndex
	var pageSize = task.PageSize

	if pageIndex < 1 {
		pageIndex = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	if task.Counter {

		var counter = ShopCartQueryCounter{}
		counter.PageIndex = pageIndex
		counter.PageSize = pageSize
		counter.RowCount, err = kk.DBQueryCount(db, a.GetShopCartTable(), a.GetPrefix(), sql.String(), args...)

		if err != nil {
			task.Result.Errno = ERROR_SHOPCART
			task.Result.Errmsg = err.Error()
			return nil
		}

		if counter.RowCount%pageSize == 0 {
			counter.PageCount = counter.RowCount / pageSize
		} else {
			counter.PageCount = counter.RowCount/pageSize + 1
		}

		task.Result.Counter = &counter
	}

	sql.WriteString(fmt.Sprintf(" LIMIT %d,%d", (pageIndex-1)*pageSize, pageSize))

	var v = ShopCart{}
	var scanner = kk.NewDBScaner(&v)

	rows, err := kk.DBQuery(db, a.GetShopCartTable(), a.GetPrefix(), sql.String(), args...)

	if err != nil {
		task.Result.Errno = ERROR_SHOPCART
		task.Result.Errmsg = err.Error()
		return nil
	}

	defer rows.Close()

	for rows.Next() {

		err = scanner.Scan(rows)

		if err != nil {
			task.Result.Errno = ERROR_SHOPCART
			task.Result.Errmsg = err.Error()
			return nil
		}

		items = append(items, v)
	}

	task.Result.Items = items

	return nil
}
