package vcloudutil

/*
 * Description : 查询上传进度的返回参数
 * Company : NetEase
 * Date : 2016-9-1
 */

type GetUploadProcessRet struct {
	RequestID string
	Offset    int64
}
