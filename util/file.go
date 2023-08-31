/*
文件控制器，包含判断是否存在、写入、
*/
package util

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	log "k8s.io/klog/v2"
)

func ExistsFile(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

type FdMananger struct {
	fdNu int //限制打开的Nu
}

var Fdm FdMananger
var FdNuMax = 6

func (f *FdMananger) OpenFile(filePath string) (*OBFile, error) {
	timeout := 0
	nu := f.fdNu
	for nu > FdNuMax {
		time.Sleep(1 * time.Second)
		timeout++
		nu = f.fdNu
	}
	f.fdNu++
	//fmt.Printf("open %s file \n", filePath)
	fileHanle, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {

		return nil, err
	}
	file := OBFile{
		*fileHanle,
	}
	return &file, nil

}

type OBFile struct {
	os.File
}

func (obf *OBFile) CloseFile() error {
	err := obf.File.Close()
	if err != nil {
		return err
	}
	Fdm.fdNu--
	if Fdm.fdNu < 0 {
		panic("painc:Fdm.fdNu<0")
	}
	return nil
}

type FileContext struct {
	context.Context
}

func GetAllFileNames(path string) ([]string, error) {
	var fileNames []string
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return fileNames, err
	}

	for _, fi := range rd {
		if !fi.IsDir() {

			fileNames = append(fileNames, fi.Name())
		}
	}
	return fileNames, nil

}
func GetAllDirNames(path string) ([]string, error) {
	var dirNames []string
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return dirNames, err
	}

	for _, fi := range rd {
		if fi.IsDir() {

			dirNames = append(dirNames, fi.Name())
		}
	}
	return dirNames, nil
}

func ReadFile(filePath string) (string, error) {
	//打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	//关闭文件
	defer file.Close()

	//读取文件内容
	buf := make([]byte, 1024*2)        // 2k大小
	_, err1 := file.Read(buf)          //n代表从文件读取内容的长度
	if err1 != nil && err1 != io.EOF { // 文件出错，同时没有到结尾
		return "", err1
	}

	return string(buf), nil
}

//basePath是固定目录路径,folderName是文件夹名
func CreateDateDir(basePath string, folderName string) error {
	folderPath := filepath.Join(basePath, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 若不存在文件，必须分成两步
		// 先创建文件夹
		if err = os.Mkdir(folderPath, 0777); err != nil {
			return err
		}

		// 再修改权限
		if err = os.Chmod(folderPath, 0777); err != nil {
			return err
		}
	}
	return nil
}
func CreateEmptyFileOnLocalPath(fileName string) error {
	if ExistsFile(fileName) {
		return nil
	}
	_, err := os.Create(fileName)
	return err

}

// RenameFile rename移动文件夹
func RenameFile(oldfile, newfile string) error {
	err := os.Rename(oldfile, newfile)
	if err != nil {
		return err
	}

	return nil
}
func CopyFile(oldfile, newfile string) error {
	return nil
}

func GetFileNames(path string, gettype string) ([]string, error) {
	var fileNames []string
	fileNames, err := GetAllFileNames(path)
	if err != nil {
		return fileNames, err
	}
	dirNames, err := GetAllDirNames(path)
	if err != nil {
		return dirNames, err
	}
	switch gettype {
	case "all":
		return append(fileNames, dirNames...), nil

	case "file":
		return fileNames, nil

	case "dir":
		return dirNames, err
	default:
		return nil, errors.New("gettype is err")
	}
	return nil, errors.New("意外错误")
}
func RMFile(filePath string) error {
	return os.Remove(filePath)
}
func WriteFileADD(filePath, data string) error {
	if data == "" {
		return nil
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Warning(err)
		return err
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(data + "\n")

	//Flush将缓存的文件真正写入到文件中
	err = write.Flush()
	return err
}

func SelectDataOnFileByString(path string, str string) (string, error) {
	redata := ""
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Warning(err)
		return "", err
	}
	//及时关闭file句柄
	defer file.Close()
	filename := path
	fileHanle, err := Fdm.OpenFile(filename)
	if err != nil {
		log.Warning(err)
		return "", err
	}
	defer func() {
		fileHanle.Close()
		Fdm.fdNu--
		runtime.GC()
	}()
	readBytes, err := ioutil.ReadAll(fileHanle)
	if err != nil {

		log.Warning(err)
		return "", err
	}
	//拆分行为切片
	results := strings.Split(string(readBytes), "\n")
	for i := 0; i < len(results); i++ {
		if strings.Contains(results[i], str) {
			redata += results[i] + "\n"
		}
	}

	return redata, nil

}
