/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2017 . All rights reserved.
 *
 * This library is free software under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3 of the License,
 * or (at your option) any later version.
 *
 */

package eth_lib

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func CheckErr(errMasg error) {
	if errMasg != nil {
		fmt.Println(errMasg)
		panic(errMasg)
	}
}

func Connect_db(dbstr string) *sql.DB {
	db, err := sql.Open("mysql", dbstr)

	CheckErr(err)

	return db
}

func Get_result(rows *sql.Rows) map[string]string {

	res := make(map[string]string)
	var addr string
	var balance string

	for rows.Next() {

		if err := rows.Scan(&addr, &balance); err != nil {
			CheckErr(err)
		}

		res[addr] = balance

	}

	return res

}
