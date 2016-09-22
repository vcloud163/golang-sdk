package vcloudutil

/*
 * Description : 视频上传初始化返回参数
 * Company : NetEase
 * Date : 2016-8-16
 */

/* 上传初始化返回参数ret部分 */
type UploadInitRet struct {
	/* 上传token信息 */
	XNosToken string

	/* 存储上传文件的桶名 */
	Bucket string

	/* 存储上传文件的对象名 */
	Object string
}

type UploadInitVideoParams struct {
	/* 返回参数的ret部分 */
	Ret UploadInitRet

	/* 返回参数中的响应码 */
	Code int

	/* 返回参数中的错误信息 */
	Msg string
}
