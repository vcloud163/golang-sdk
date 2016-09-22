package util

/*
 * Description : 进行文件操作的工具类
 * Company : NetEase
 * Date : 2016-8-16
 */

import (
	"os"
)

/* 判断文件是否存在 */
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil || os.IsExist(err)
}

/* 获取文件长度（单位：字节） */
func GetSize(path string) int64 {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	fileSize := fileInfo.Size()

	return fileSize
}

/* 获取文件名称 */
func GetFileName(path string) string {
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	fileName := fileInfo.Name()

	return fileName
}
