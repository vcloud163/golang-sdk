package vcloudutil

/*
 * Description : 分片上传视频输出参数
 * Company : NetEase
 * Date : 2016-8-17
 */

type UploadVideoFragmentParams struct {
	/* 服务器端生成的唯一UUID，用于记录日志排查问题使用 */
	RequestId string

	/* 输出参数中的下一个上传片在上传块中的偏移 */
	Offset int64

	/* 输出参数中的上传上下文 */
	Context string

	/* 输出参数中的上传回调信息 */
	CallbackRetMsg string

	/* 输出参数中的错误信息 */
	ErrMsg string
}
