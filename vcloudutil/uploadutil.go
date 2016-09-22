package vcloudutil

/*
 * Description : 业务服务器调用网易视频云上传视频文件
 * Company : NetEase
 * Date : 2016-8-11
 */

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"golang-sdk/auth"
	"golang-sdk/config"
	"golang-sdk/util"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type UploadUtil struct {
	/* 封装视频上传初始化输入参数 */
	InitParams UploadInitParams

	/** 根据对象名查询视频或水印图片主ID输入参数的封装类*/
	//--------------------------------------

	/* 分片（上传）最大字节数 */
	UploadDataSize int64

	/* 存储上传文件的桶名 */
	Bucket string

	/* 上传加速节点地址 */
	UploadHost string

	/* 存储上传文件的对象名 */
	Object string

	/* 当前分片在整个对象中的起始偏移量 */
	Offset uint64

	/* 上传上下文 */
	Context string

	/* 上传文件剩余大小 */
	RemainderSize string

	/* 上传的token信息 */
	XNosToken string
}

/* 构造UploadUtil操作对象 */
func NewUploadUtil(params UploadInitParams) *UploadUtil {
	/* 检查视频上传初始化必填参数是否缺失 */
	if !checkInitParams(params) {
		panic("OriginFileName is null or invalid!")
	}

	uploadUtil := new(UploadUtil)
	uploadUtil.UploadDataSize = config.UploadDataSize
	uploadUtil.InitParams.CallbackUrl = params.CallbackUrl
	uploadUtil.InitParams.Description = params.Description
	uploadUtil.InitParams.OriginFileName = params.OriginFileName
	uploadUtil.InitParams.PresetId = params.PresetId
	uploadUtil.InitParams.TypeId = params.TypeId
	uploadUtil.InitParams.UploadCallbackUrl = params.UploadCallbackUrl
	uploadUtil.InitParams.UserDefInfo = params.UserDefInfo
	uploadUtil.InitParams.UserFileName = params.UserDefInfo
	uploadUtil.InitParams.WatermarkId = params.WatermarkId

	return uploadUtil
}

/* 视频上传初始化
 * @param  urlStr 上传初始化接口url
 * @param  key    访问凭证
 * @return UploadInitVideoParams 返回的初始化参数
 */
func (uutil *UploadUtil) InitUpload(urlStr string, key auth.Key) UploadInitVideoParams {
	/* 计算checksum */
	curTime := time.Now().Unix()
	curTimeStr := strconv.FormatInt(curTime, 10)
	nonce := "1"
	checkSum := key.GetCheckSum(key.GetSecretKey(), nonce, curTimeStr)

	/* 构造POST请求，传入请求参数 */
	client := &http.Client{}

	var postDataMap map[string]string
	postDataMap = make(map[string]string)

	if uutil.InitParams.OriginFileName != "" {
		postDataMap["originFileName"] = uutil.InitParams.OriginFileName
	}
	if uutil.InitParams.UserFileName != "" {
		postDataMap["userFileName"] = uutil.InitParams.UserFileName
	}
	if uutil.InitParams.TypeId != 0 {
		postDataMap["typeId"] = strconv.FormatInt(uutil.InitParams.TypeId, 10)
	}
	if uutil.InitParams.PresetId != 0 {
		postDataMap["presetId"] = strconv.FormatInt(uutil.InitParams.PresetId, 10)
	}
	if uutil.InitParams.CallbackUrl != "" {
		postDataMap["callbackUrl"] = uutil.InitParams.CallbackUrl
	}
	if uutil.InitParams.Description != "" {
		postDataMap["description"] = uutil.InitParams.Description
	}
	if uutil.InitParams.WatermarkId != 0 {
		postDataMap["watermarkId"] = strconv.FormatInt(uutil.InitParams.WatermarkId, 10)
	}
	if uutil.InitParams.UploadCallbackUrl != "" {
		postDataMap["uploadCallbackUrl"] = uutil.InitParams.UploadCallbackUrl
	}
	if uutil.InitParams.UserDefInfo != "" {
		postDataMap["userDefInfo"] = uutil.InitParams.UserDefInfo
	}

	postDataBytes, err := json.Marshal(postDataMap)
	if err != nil {
		panic("Failed to construct json object! error: " + err.Error())
	}

	postBytesReader := bytes.NewReader(postDataBytes)
	req, err := http.NewRequest("POST", urlStr, postBytesReader)
	if err != nil {
		panic("Failed to new http request! error: " + err.Error())
	}

	/* 设置请求头 */
	req.Header.Add("AppKey", key.GetAccessKey())
	req.Header.Add("Nonce", nonce)
	req.Header.Add("CurTime", curTimeStr)
	req.Header.Add("CheckSum", checkSum)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	/* 执行请求 */
	resp, err := client.Do(req)
	if err != nil {
		panic("Failed to excute http request! error: " + err.Error())
	}
	defer resp.Body.Close()

	/* 解析返回数据 */
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Failed to get response data! error: " + err.Error())
	}

	var initDataRet UploadInitVideoParams
	json.Unmarshal([]byte(respData), &initDataRet)

	return initDataRet
}

/* 获取上传加速节点
 * @param  urlStr 获取上传加速节点接口url
 * @param  bucket 上传桶名
 * @return UploadHostAddress 返回加速节点
 */
func (uutil *UploadUtil) GetUploadHost(urlStr string, bucket string) UploadHostAddress {
	/* 检查桶名是否合法 */
	if !checkStrParams(bucket) {
		panic("bucket is null or invalid!")
	}

	/* 发起Get请求 */
	response, _ := http.Get(urlStr + "&bucketname=" + bucket)
	defer response.Body.Close()

	/* 解析返回数据 */
	respData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic("Failed to get response data! error: " + err.Error())
	}

	var uploadHost UploadHostAddress
	json.Unmarshal([]byte(respData), &uploadHost)

	return uploadHost
}

/* 简单上传视频
 * @param  uploadParams 上传初始化参数
 * @param  uploadHost   上传加速节点
 * @param  filePath     上传文件本地路径
 * @return
 */
func (uutil *UploadUtil) UploadVideo(uploadParams UploadInitVideoParams, uploadHost UploadHostAddress, filePath string) {
	/* 校验参数是否合法 */
	if !checkStrParams(uploadParams.Ret.Bucket) {
		panic("bucket is null or invalid!")
	}
	if !checkStrParams(uploadParams.Ret.Object) {
		panic("object is null or invalid!")
	}
	if !checkStrParams(uploadParams.Ret.XNosToken) {
		panic("xNosToken is null or invalid!")
	}
	if len(uploadHost.Upload) == 0 || !checkStrParams(uploadHost.Upload[0]) {
		panic("upload host is null or invalid!")
	}

	/* 分片偏移量 */
	var offset int64 = 0
	/* 上传文件上下文 */
	var context string
	/* 上传文件剩余大小 */
	fileSize := util.GetSize(filePath)
	var remainderSize int64 = fileSize

	fin, err := os.Open(filePath)
	defer fin.Close()
	if err != nil {
		panic("Failed to read file! error: " + err.Error())
	}

	for remainderSize > 0 {
		uploadVideoParams := uutil.uploadVideoFragment(uploadHost.Upload[0], uploadParams.Ret.Bucket, uploadParams.Ret.Object, uploadParams.Ret.XNosToken, fin, offset, context, remainderSize, false)
		context = uploadVideoParams.Context
		offset = uploadVideoParams.Offset
		remainderSize = fileSize - offset
	}
}

/* 分片上传视频
 * @param  uploadHost    上传加速节点IP
 * @param  uploadParams  上传初始化参数
 * @param  fin           文件输入流
 * @param  offset        偏移量
 * @param  context       上传文件上下文
 * @param  remainderSize 剩余大小（字节）
 * @param  flag          true表示断点续传的第一次上传
 * @return UploadVideoFragmentParams 分片上传视频输出参数
 */
func (uutil *UploadUtil) uploadVideoFragment(uploadHost string, bucket string, object string, xNosToken string, fin *os.File, offset int64, context string, remainderSize int64, flag bool) UploadVideoFragmentParams {
	var urlStr string
	var buf []byte
	/* 判断是否为最后一片 */
	if remainderSize <= config.UploadDataSize {
		urlStr = uploadHost + "/" + bucket + "/" + object + "?offset=" + strconv.FormatInt(offset, 10) + "&complete=true&version=1.0"
		buf = make([]byte, remainderSize)
	} else {
		urlStr = uploadHost + "/" + bucket + "/" + object + "?offset=" + strconv.FormatInt(offset, 10) + "&complete=false&version=1.0"
		buf = make([]byte, config.UploadDataSize)
	}
	/* 如果不是第一次传输，需要加上context参数 */
	if context != "" {
		urlStr = urlStr + "&context=" + context
	}

	var len int
	/* 读取文件内容至buf */
	if !flag {
		len, _ = fin.Read(buf)
	} else {
		len, _ = fin.ReadAt(buf, offset)
	}

	if len <= 0 {
		panic("Failed to read file content!")
	}

	/* 构造POST请求，传入请求参数 */
	client := &http.Client{}
	postBytesReader := bytes.NewReader(buf)
	req, err := http.NewRequest("POST", urlStr, postBytesReader)
	if err != nil {
		panic("Failed to new http request! error: " + err.Error())
	}

	/* 设置请求头 */
	req.Header.Add("x-nos-token", xNosToken)

	/* 执行请求 */
	resp, err := client.Do(req)
	if err != nil {
		panic("Failed to excute http request! error: " + err.Error())
	}
	defer resp.Body.Close()

	/* 解析返回数据 */
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Failed to get response data! error: " + err.Error())
	}

	var uploadFragmentParams UploadVideoFragmentParams
	json.Unmarshal([]byte(respData), &uploadFragmentParams)

	return uploadFragmentParams
}

/* 上传完成后查询视频主ID和水印图片ID
 * @param  urlStr      查询接口url
 * @param  key         访问凭证
 * @param  objcetNames 对象名
 * @return QueryVideoIDorWatermarkIDParams 查询返回参数
 */
func (uutil *UploadUtil) QueryVideoIDorWatermarkID(urlStr string, key auth.Key, objcetNames []string) QueryVideoIDorWatermarkIDParams {
	if objcetNames == nil {
		panic("object names is null!")
	}
	/* 计算checksum */
	curTime := time.Now().Unix()
	curTimeStr := strconv.FormatInt(curTime, 10)
	nonce := "1"
	checkSum := key.GetCheckSum(key.GetSecretKey(), nonce, curTimeStr)

	/* 构造POST请求，传入请求参数 */
	client := &http.Client{}

	var postDataMap map[string][]string
	postDataMap = make(map[string][]string)
	postDataMap["objectNames"] = objcetNames

	postDataBytes, err := json.Marshal(postDataMap)
	if err != nil {
		panic("Failed to construct json object! error: " + err.Error())
	}

	postBytesReader := bytes.NewReader(postDataBytes)
	req, err := http.NewRequest("POST", urlStr, postBytesReader)
	if err != nil {
		panic("Failed to new http request! error: " + err.Error())
	}

	/* 设置请求头 */
	req.Header.Add("AppKey", key.GetAccessKey())
	req.Header.Add("Nonce", nonce)
	req.Header.Add("CurTime", curTimeStr)
	req.Header.Add("CheckSum", checkSum)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	/* 执行请求 */
	resp, err := client.Do(req)
	if err != nil {
		panic("Failed to excute http request! error: " + err.Error())
	}
	defer resp.Body.Close()

	/* 解析返回数据 */
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Failed to get response data! error: " + err.Error())
	}

	var queryRslt QueryVideoIDorWatermarkIDParams
	json.Unmarshal([]byte(respData), &queryRslt)

	return queryRslt
}

/* 断点续传
 * @param  filePath 上传文件本地路径
 * @param  recoder  断点记录
 * @return 上传成功返回object，否则返回空字符串""
 */
func (uutil *UploadUtil) UploadVideoWithRecorder(filePath string, recoder *Recoder, key auth.Key) string {
	/* 第一次上传 */
	if !recoder.UploadAgain {

		/* 检查上传文件是否存在 */
		if !util.FileExist(filePath) {
			panic("Upload file not exist!")
		}
		/* 检查用户传入的参数是否合法 */
		if !checkInitParams(uutil.InitParams) {
			panic("OriginFileName is null or invalid!")
		}
		/* 上传初始化，获取xNosToken（上传token）、bucket（存储对象的桶名）、object（生成的唯一对象名） */
		responseData := uutil.InitUpload(config.InitUploadVideoURL, key)
		if responseData.Code != 200 {
			panic("Failed to init upload! Error message: " + responseData.Msg)
		}
		/* 获取上传加速节点地址 */
		uploadHost := uutil.GetUploadHost(config.GetUploadHostURL, responseData.Ret.Bucket)

		/* 分片偏移量 */
		var offset int64 = 0
		/* 上传文件上下文 */
		var context string
		/* 上传文件剩余大小 */
		fileSize := util.GetSize(filePath)
		var remainderSize int64 = fileSize

		fin, err := os.Open(filePath)
		defer fin.Close()
		if err != nil {
			panic("Failed to read file! error: " + err.Error())
		}

		//		count := 0
		for remainderSize > 0 {
			uploadVideoParams := uutil.uploadVideoFragment(uploadHost.Upload[0], responseData.Ret.Bucket, responseData.Ret.Object, responseData.Ret.XNosToken, fin, offset, context, remainderSize, false)
			context = uploadVideoParams.Context
			offset = uploadVideoParams.Offset
			remainderSize = fileSize - offset

			var uploadRecoderMap map[string]string
			uploadRecoderMap = make(map[string]string)

			uploadRecoderMap["bucket"] = responseData.Ret.Bucket
			uploadRecoderMap["uploadHost"] = uploadHost.Upload[0]
			uploadRecoderMap["object"] = responseData.Ret.Object
			uploadRecoderMap["offset"] = strconv.FormatInt(offset, 10)
			uploadRecoderMap["context"] = context
			uploadRecoderMap["remainderSize"] = strconv.FormatInt(remainderSize, 10)
			uploadRecoderMap["xNosToken"] = responseData.Ret.XNosToken

			recoder.SaveRecoderInfo(uploadRecoderMap)
			//			count++
			//			if count == 1 {
			//				break
			//			}
		}

		if remainderSize == 0 {
			/* 上传完成则删除记录 */
			recoder.DeleteRecoder()
			/* 上传完成返回object */
			return responseData.Ret.Object
		}
	} else {

		/* 读取上传信息 */
		uploadRecoderMap0 := recoder.GetRecoderInfo()

		bucket, _ := uploadRecoderMap0["bucket"]
		uploadHost, _ := uploadRecoderMap0["uploadHost"]
		object, _ := uploadRecoderMap0["object"]
		offset, _ := strconv.ParseInt(uploadRecoderMap0["offset"], 10, 64)
		context, _ := uploadRecoderMap0["context"]
		remainderSize, _ := strconv.ParseInt(uploadRecoderMap0["remainderSize"], 10, 64)
		xNosToken, _ := uploadRecoderMap0["xNosToken"]

		/* 跳过已经上传的部分继续分片上传 */
		fin, err := os.Open(filePath)
		defer fin.Close()
		if err != nil {
			panic("Failed to read file! error: " + err.Error())
		}

		for remainderSize > 0 {
			uploadVideoParams := uutil.uploadVideoFragment(uploadHost, bucket, object, xNosToken, fin, offset, context, remainderSize, true)
			context = uploadVideoParams.Context
			offset = uploadVideoParams.Offset
			remainderSize = util.GetSize(filePath) - offset

			var uploadRecoderMap1 map[string]string
			uploadRecoderMap1 = make(map[string]string)

			uploadRecoderMap1["bucket"] = bucket
			uploadRecoderMap1["uploadHost"] = uploadHost
			uploadRecoderMap1["object"] = object
			uploadRecoderMap1["offset"] = strconv.FormatInt(offset, 10)
			uploadRecoderMap1["context"] = context
			uploadRecoderMap1["remainderSize"] = strconv.FormatInt(remainderSize, 10)
			uploadRecoderMap1["xNosToken"] = xNosToken

			recoder.SaveRecoderInfo(uploadRecoderMap1)
		}

		if remainderSize == 0 {
			/* 上传完成则删除记录 */
			recoder.DeleteRecoder()
			/* 上传完成返回object */
			return object
		}
	}

	return ""
}

/* 检查视频上传初始化是否合法 */
func checkInitParams(params UploadInitParams) bool {
	fileName := strings.TrimSpace(params.OriginFileName)

	if fileName == "" || len(fileName) == 0 {
		return false
	}

	return true
}

/* 检查字符串参数是否合法 */
func checkStrParams(str string) bool {
	strTrim := strings.TrimSpace(str)

	if strTrim == "" || len(strTrim) == 0 {
		return false
	}

	return true
}
