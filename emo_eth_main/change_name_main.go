/**
 * Author: wengqiang (email: wens.wq@gmail.com  site: qiangweng.site)
 *
 * Copyright © 2015--2017 . All rights reserved.
 */

package emo_eth_main

import (
	"fmt"

	"github.com/wens07/libchange_name"
)

func main() {

	fmt.Println("renaming, please wait......")

	fd := libchange_name.NewPathFinder()
	//
	//files := fd.PathFile("E:\\goopal3.0\\kgt_project\\bitshares1-core")
	//files := fd.PathFile("E:\\goopal3.0\\emo_project\\bitshares1-core")
	files := fd.PathFile("E:\\goopal3.0\\ctc_project\\include")
	//files := fd.PathFile("E:\\goopal3.0\\ctc_project\\src")

	frp := libchange_name.NewFileHandle()
	rps := make(map[string]string)
	//rps["BTS"] = "EMO"
	//rps["bts"] = "emo"
	//rps["Bts"] = "Emo"
	//rps["ubcore"] = "ctccore"
	//rps["UBCORE"] = "CTCCORE"
	//rps["Ubcore"] = "Ctccore"
	//rps["UbCore"] = "CtcCore"
	rps["UB"] = "CTC"
	rps["Ub"] = "Ctc"
	//rps["ub_address"] = "ctc_address"
	frp.FileReplace(files, rps)
	frp.PathNameReplace(files, rps)
	files = frp.FileNameReplace(files, rps)
	frp.PathClear(files, rps)

	fmt.Println("rename finished!")

}
