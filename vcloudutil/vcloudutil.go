package vcloudutil

/*
 * Description : 业务服务器调用网易视频云功能
 * Company : NetEase
 * Date : 2016-8-11
 */

import (
	"bytes"
	"encoding/json"
	"golang-sdk/auth"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type VcloudUtil struct {
	key    auth.Key
	Upload *UploadUtil
}

/* 构造VcloudUtil操作对象 */
func NewVcloudUtil(key auth.Key) *VcloudUtil {
	vcloudUtil := new(VcloudUtil)
	vcloudUtil.key = key

	return vcloudUtil
}

/* 构造VcloudUtil操作对象（不传入key） */
func NewVcloudUtilNoKey() *VcloudUtil {
	vcloudUtil := new(VcloudUtil)
	return vcloudUtil
}

/* 上传完成后查询视频主ID和水印图片ID
 * @param  urlStr      接口url
 * @param  callback    回调url
 * @return CommonRet   请求返回参数
 */
func (vutil *VcloudUtil) SetUploadCallback(urlStr string, callback string) CommonRet {
	/* 计算checksum */
	curTime := time.Now().Unix()
	curTimeStr := strconv.FormatInt(curTime, 10)
	nonce := "1"
	checkSum := vutil.key.GetCheckSum(vutil.key.GetSecretKey(), nonce, curTimeStr)

	/* 构造POST请求，传入请求参数 */
	client := &http.Client{}

	var postDataMap map[string]string
	postDataMap = make(map[string]string)

	postDataMap["callbackUrl"] = callback

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
	req.Header.Add("AppKey", vutil.key.GetAccessKey())
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

	var ret CommonRet
	json.Unmarshal([]byte(respData), &ret)

	return ret
}

/* 上传完成后查询视频主ID和水印图片ID
 * @param  recordFilePath  保存断点信息的本地文件
 * @return int64           上传进度偏移量
 */
func (vutil *VcloudUtil) GetUploadProcess(recordFilePath string) int64 {
	recoder := NewRecoder(recordFilePath)

	/* 如果断点信息文件存在 */
	if recoder.UploadAgain == true {
		/* 读取上传信息 */
		uploadRecoderMap0 := recoder.GetRecoderInfo()

		bucket, _ := uploadRecoderMap0["bucket"]
		uploadHost, _ := uploadRecoderMap0["uploadHost"]
		object, _ := uploadRecoderMap0["object"]
		context, _ := uploadRecoderMap0["context"]
		xNosToken, _ := uploadRecoderMap0["xNosToken"]

		var urlStr string
		urlStr = uploadHost + "/" + bucket + "/" + object + "?uploadContext&context=" + context + "&version=1.0"

		client := &http.Client{}
		reqest, _ := http.NewRequest("GET", urlStr, nil)
		reqest.Header.Set("x-nos-token", xNosToken)

		/* 解析返回数据 */
		response, err := client.Do(reqest)
		if err != nil {
			panic("Failed to get upload process! error: " + err.Error())
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			respData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				panic("Failed to get response data! error: " + err.Error())
			}

			var ret GetUploadProcessRet
			json.Unmarshal([]byte(respData), &ret)

			return ret.Offset
		} else {
			panic("Failed to get upload process!")
		}

	} else {
		panic("Record file not exist!")
	}

}
