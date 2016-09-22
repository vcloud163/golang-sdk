package main

/*
 * SDK demo
 */

import (
	"fmt"
	"golang-sdk/auth"
	"golang-sdk/config"
	"golang-sdk/util"
	"golang-sdk/vcloudutil"
)

/* 上传初始化demo
 * 获取上传加速节点地址demo
 * 简单文件上传demo
 * 上传完成后，根据objectName查询demo
 */
func UploadDemos(filePath string) {
	key := auth.Key{config.AppKey, config.AppSecret}
	vutil := vcloudutil.NewVcloudUtil(key)

	initParams := vcloudutil.UploadInitParams{}
	initParams.OriginFileName = util.GetFileName(filePath)
	initParams.UploadCallbackUrl = "http://106.2.44.248/xhr/vod/client/transcode/callback"
	vutil.Upload = vcloudutil.NewUploadUtil(initParams)

	/* 上传初始化，获取xNosToken（上传token）、bucket（存储对象的桶名）、object（生成的唯一对象名） */
	responseDate := vutil.Upload.InitUpload(config.InitUploadVideoURL, key)
	if responseDate.Code != 200 {
		fmt.Println("Failed to init upload! Error message: " + responseDate.Msg)
		return
	}

	fmt.Println(responseDate.Code)
	fmt.Println(responseDate.Msg)
	fmt.Println(responseDate.Ret.Bucket)
	fmt.Println(responseDate.Ret.Object)
	fmt.Println(responseDate.Ret.XNosToken)

	/* 获取上传加速节点地址 */
	uploadHost := vutil.Upload.GetUploadHost(config.GetUploadHostURL, responseDate.Ret.Bucket)

	fmt.Println(uploadHost.Lbs)
	fmt.Println(uploadHost.Upload)
	fmt.Println(uploadHost.Upload[0])

	/* 简单文件上传 */
	vutil.Upload.UploadVideo(responseDate, uploadHost, filePath)

	/* 上传完成后，根据objectName查询视频ID */
	var objectNames []string
	objectNames = append(objectNames, responseDate.Ret.Object)
	queryRslt := vutil.Upload.QueryVideoIDorWatermarkID(config.QueryVideoIDURL, key, objectNames)
	fmt.Println(queryRslt.Code)
	fmt.Println(queryRslt.Ret.Count)
	fmt.Println(queryRslt.Ret.List[0].ImgId)
	fmt.Println(queryRslt.Ret.List[0].Vid)
	fmt.Println(queryRslt.Ret.List[0].ObjectName)
}

/*
 * 断点续传demo
 * @param  filePath 上传文件路径
 * @param  recordFilePath 保存断点信息的本地文件
 */
func UploadWithRecoderDemo(filePath string, recordFilePath string) {
	key := auth.Key{config.AppKey, config.AppSecret}
	vutil := vcloudutil.NewVcloudUtil(key)

	initParams := vcloudutil.UploadInitParams{}
	initParams.OriginFileName = util.GetFileName(filePath)
	initParams.UploadCallbackUrl = "http://106.2.44.248/xhr/vod/client/transcode/callback"
	vutil.Upload = vcloudutil.NewUploadUtil(initParams)

	recoder := vcloudutil.NewRecoder(recordFilePath)
	object := vutil.Upload.UploadVideoWithRecorder(filePath, recoder, key)

	/* 上传完成后，根据objectName查询视频ID */
	if object != "" {
		var objectNames []string
		objectNames = append(objectNames, object)
		queryRslt := vutil.Upload.QueryVideoIDorWatermarkID(config.QueryVideoIDURL, key, objectNames)
		fmt.Println(queryRslt.Code)
		fmt.Println(queryRslt.Ret.Count)
		fmt.Println(queryRslt.Ret.List[0].ImgId)
		fmt.Println(queryRslt.Ret.List[0].Vid)
		fmt.Println(queryRslt.Ret.List[0].ObjectName)
	}
}

/*
 * 设置上传完成回调地址demo
 * @param  callback 回调地址
 */
func SetUploadCallbackDemo(callback string) {
	key := auth.Key{config.AppKey, config.AppSecret}
	vutil := vcloudutil.NewVcloudUtil(key)

	ret := vutil.SetUploadCallback(config.SetUploadCallback, callback)
	fmt.Println(ret.Code)
	if ret.Code != 200 {
		fmt.Println(ret.Msg)
	}
}

/*
 * 获取上传进度demo
 * @param  recordFilePath 保存断点信息的本地文件
 */
func GetUploadProcessDemo(recordFilePath string) {
	vutil := vcloudutil.NewVcloudUtilNoKey()
	offset := vutil.GetUploadProcess(recordFilePath)

	fmt.Println(offset)
}

func main() {
	/* 上传文件路径 */
	filePath := "D:\\abc.mp4"

	/* 普通上传 */
	UploadDemos(filePath)

	/* 断点续传 */
	//recordFilePath := "D:\\recoder.txt"
	//UploadWithRecoderDemo(filePath, recordFilePath)

	/* 上传完成回调地址 */
	//callback := "http://106.2.44.248/xhr/vod/client/transcode/callback"
	//SetUploadCallbackDemo(callback)

	/* 获取上传进度 */
	//GetUploadProcessDemo(recordFilePath)

}
