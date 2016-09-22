package vcloudutil

/*
 * Description : 断点续传的上传进度记录
 * Company : NetEase
 * Date : 2016-8-20
 */

import (
	"bufio"
	"golang-sdk/util"
	"io"
	"os"
	"strings"
)

type Recoder struct {
	/* 标记文件是否是续传，true：续传，false：非续传 */
	UploadAgain bool

	/* 文件路径 */
	FilePath string
}

/* 构造Recoder对象 */
func NewRecoder(recoderFilePath string) *Recoder {
	recoder := new(Recoder)
	recoder.FilePath = recoderFilePath

	/* 如果断点记录文件不存在 */
	if !util.FileExist(recoderFilePath) {
		//		recoderFile, err := os.OpenFile(recoderFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm|os.ModeTemporary)
		//		if err != nil {
		//			panic(err)
		//		}
		//      recoder.RecoderFile = recoderFile
		recoder.UploadAgain = false
	} else {
		//		recoderFile, err := os.OpenFile(recoderFilePath, os.O_RDWR, os.ModePerm|os.ModeTemporary)
		//		if err != nil {
		//			panic(err)
		//		}
		//		recoder.RecoderFile = recoderFile
		recoder.UploadAgain = isUploadAgain(recoderFilePath)
	}

	return recoder
}

/* 设置文件是否续传 */
func (recoder *Recoder) SetUploadAgain(uploadAgain bool) {
	recoder.UploadAgain = uploadAgain
}

/* 保存上传进度信息 */
func (recoder *Recoder) SaveRecoderInfo(uploadRecoderMap map[string]string) {
	recoderFile, err := os.OpenFile(recoder.FilePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm|os.ModeTemporary)
	if err != nil {
		panic(err)
	}
	defer recoderFile.Close()

	recoderFile.WriteString("uploadAgain:true\n")
	recoderFile.WriteString(uploadRecoderMap["bucket"] + "\n")
	recoderFile.WriteString(uploadRecoderMap["uploadHost"] + "\n")
	recoderFile.WriteString(uploadRecoderMap["object"] + "\n")
	recoderFile.WriteString(uploadRecoderMap["offset"] + "\n")
	recoderFile.WriteString(uploadRecoderMap["context"] + "\n")
	recoderFile.WriteString(uploadRecoderMap["remainderSize"] + "\n")
	recoderFile.WriteString(uploadRecoderMap["xNosToken"] + "\n")
}

/* 获取上传进度信息 */
func (recoder *Recoder) GetRecoderInfo() map[string]string {
	recoderFile, err := os.OpenFile(recoder.FilePath, os.O_RDONLY, os.ModePerm|os.ModeTemporary)
	if err != nil {
		panic(err)
	}
	defer recoderFile.Close()

	var uploadRecoderMap map[string]string
	uploadRecoderMap = make(map[string]string)

	buff := bufio.NewReader(recoderFile)

	uploadAgain, _ := buff.ReadString('\n')
	bucket, _ := buff.ReadString('\n')
	uploadHost, _ := buff.ReadString('\n')
	object, _ := buff.ReadString('\n')
	offset, _ := buff.ReadString('\n')
	context, _ := buff.ReadString('\n')
	remainderSize, _ := buff.ReadString('\n')
	xNosToken, _ := buff.ReadString('\n')

	/* 记录上传信息的文件损坏 */
	if uploadAgain == "" || bucket == "" || uploadHost == "" || object == "" || offset == "" || context == "" || remainderSize == "" || xNosToken == "" {
		panic("recoder file has been broken!")
	}

	bucket = strings.Replace(bucket, "\n", "", 1)
	uploadHost = strings.Replace(uploadHost, "\n", "", 1)
	object = strings.Replace(object, "\n", "", 1)
	offset = strings.Replace(offset, "\n", "", 1)
	context = strings.Replace(context, "\n", "", 1)
	remainderSize = strings.Replace(remainderSize, "\n", "", 1)
	xNosToken = strings.Replace(xNosToken, "\n", "", 1)

	uploadRecoderMap["bucket"] = bucket
	uploadRecoderMap["uploadHost"] = uploadHost
	uploadRecoderMap["object"] = object
	uploadRecoderMap["offset"] = offset
	uploadRecoderMap["context"] = context
	uploadRecoderMap["remainderSize"] = remainderSize
	uploadRecoderMap["xNosToken"] = xNosToken

	return uploadRecoderMap
}

/* 上传完成后，删除记录上传进度信息的文件 */
func (recoder *Recoder) DeleteRecoder() {
	os.Remove(recoder.FilePath)
}

/* 判断已存在的记录文件是否包含上传进度相关信息 */
func isUploadAgain(recoderFilePath string) bool {
	recoderFile, err := os.OpenFile(recoderFilePath, os.O_RDONLY, os.ModePerm|os.ModeTemporary)
	if err != nil {
		panic(err)
	}
	defer recoderFile.Close()

	buff := bufio.NewReader(recoderFile) //读入缓存
	/* 保存上传进度信息的第一行均为"uploadAgain:true" */
	line, err := buff.ReadString('\n') //以'\n'为结束符读入一行
	if err != nil || io.EOF == err {
		panic("Failed to read recoder file. err: " + err.Error())
	}

	return strings.EqualFold(line, "uploadAgain:true\n")
}
