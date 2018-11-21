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

package libchange_name

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

type FileHandle struct {
}

func NewFileHandle() *FileHandle {
	return &FileHandle{}
}
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func (fh *FileHandle) FileReplace(files []string, rps map[string]string) {
	for _, file := range files {
		exist := checkFileIsExist(file)
		if !exist {
			fmt.Println(file)
			continue
		}

		readfile, rerr := os.OpenFile(file, os.O_RDWR, 0666)
		writefile, werr := os.OpenFile(file+"_temp", os.O_CREATE|os.O_RDWR, 0666)
		if rerr != nil || werr != nil {
			fmt.Println(file)
			fmt.Println(rerr.Error(), werr.Error())
			continue
		}
		rd := bufio.NewReader(readfile)
		wd := bufio.NewWriter(writefile)
		//CMakeList := path.Ext(file) == ".txt" && (strings.Contains(file,"CMakeList"))
		for {
			if path.Ext(file) == ".txt" && !(strings.Contains(file, "CMakeList")) {
				break
			}
			line, _, err := rd.ReadLine() //以'\n'为结束符读入一行

			if err != nil || io.EOF == err {
				break
			}
			for old, new := range rps {
				line = []byte(strings.Replace(string(line), old, new, -1))
			}
			wd.WriteString(string(line) + "\n")
			wd.Flush()
			//fmt.Println("otherline",string(line))
		}
		readfile.Close()
		err := os.Remove(file)
		if err != nil {
			fmt.Println(file)
			fmt.Println(err.Error())
			continue
		}
		writefile.Close()
		os.Rename(file+"_temp", file)
		if err != nil {
			fmt.Println(file)
			fmt.Println(err.Error())
			continue
		}
	}
}

func (fh *FileHandle) PathNameReplace(files []string, rps map[string]string) {
	for _, file := range files {
		exist := checkFileIsExist(file)
		if !exist {
			fmt.Printf("file (%s) not exists!\n", file)
			continue
		}

		index := strings.LastIndex(file, "\\")
		dirstr := file[0:index]
		//tmpolddir := dirstr

		for old, new := range rps {
			if strings.Contains(dirstr, old) {
				dirstr = strings.Replace(dirstr, old, new, -1)

				err := os.MkdirAll(dirstr, os.ModeDir)
				//err := os.Rename(tmpolddir, dirstr)

				if err != nil {
					fmt.Printf("mkdir  (%s) error!\n", dirstr)
					fmt.Println(err.Error())
					continue
				}

			}

		}

	}

}

func (fh *FileHandle) FileNameReplace(files []string, rps map[string]string) []string {

	res := files

	for _, file := range files {
		exist := checkFileIsExist(file)
		if !exist {
			fmt.Printf("file (%s) not exists!\n", file)
			continue
		}

		oldfile := file

		for old, new := range rps {
			if strings.Contains(file, old) {
				file = strings.Replace(file, old, new, -1)
			}

		}

		err := os.Rename(oldfile, file)
		if err != nil {
			fmt.Printf("rename file (%s) error!\n", file)
			fmt.Println(err.Error())
			continue
		}
	}

	return res
}

func (fh *FileHandle) PathClear(files []string, rps map[string]string) {

	for _, file := range files {

		index := strings.LastIndex(file, "\\")
		dirstr := file[0:index]

		for old, _ := range rps {
			if strings.Contains(dirstr, old) {
				index = strings.Index(file, old)
				dirstr = file[0 : index+len(old)]

				//fmt.Println(dirstr)
				err := os.RemoveAll(dirstr)

				if err != nil {
					fmt.Printf("rmdir  (%s) error!\n", dirstr)
					fmt.Println(err.Error())
					continue
				}

			}

		}

	}
}
