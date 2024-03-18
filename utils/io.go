package utils

import (
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

//go:embed redc-templates
var local embed.FS
var dirs []string

// File copies a single file from src to dst
func File(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// Dir copies a whole directory recursively
func Dir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(srcfp, dstfp); err != nil {
				//fmt.Println(err)
			}
		} else {
			if err = File(srcfp, dstfp); err != nil {
				//fmt.Println(err)
			}
		}
	}
	return nil
}

func GetFilesAndDirs(dirPth string) (files []string, dirs []string) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		os.Exit(3)
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			GetFilesAndDirs(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".go")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	return files, dirs
}

// ReleaseDir 释放文件夹
func ReleaseDir(path string) {
	dirs, _ := local.ReadDir(path)
	for _, entry := range dirs {
		if entry.IsDir() {
			//_ = utils.Dir(path+"/"+entry.Name(), path+"/"+entry.Name())
			err := os.MkdirAll(path+"/"+entry.Name(), os.ModePerm)
			if err != nil {
				fmt.Println(err)
				return
			}
			ReleaseDir(path + "/" + entry.Name())
		} else {
			//_ = utils.File(path+"/"+entry.Name(), path+"/"+entry.Name())
			out, _ := os.Create(path + "/" + entry.Name())
			in, _ := local.Open(path + "/" + entry.Name())
			_, err := io.Copy(out, in)
			if err != nil {
				fmt.Println(err)
				return
			}
			in.Close()
		}
	}
}

// ChechDirMain 递归
func ChechDirMain(dirPth string) []string {
	ChechDirSub(dirPth)
	return dirs
}

func ChechDirSub(dirPth string) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		os.Exit(3)
	}

	PthSep := string(os.PathSeparator)

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			ChechDirSub(dirPth + "/" + fi.Name())
		}
	}

}

func CheckFileName(path string, key string) bool {
	files, err := ioutil.ReadDir(path)
	if err == nil {
		for _, fileInfo := range files {
			if strings.Contains(fileInfo.Name(), key) {
				return true
			}
		}
	}
	return false
}
