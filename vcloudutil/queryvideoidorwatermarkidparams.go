package vcloudutil

/*
 * Description : 根据对象名查询视频ID输出参数
 * Company : NetEase
 * Date : 2016-8-17
 */

type QueryVideoIDorWatermarkIDParams struct {
	/* 响应码 */
	Code int

	/* 输出参数的ret部分 */
	Ret QueryVideoIDorWatermarkIDRet

	/* 输出参数中的错误信息 */
	Msg string
}

type QueryVideoIDorWatermarkIDRet struct {
	/* 返回参数ret部分的vid个数 */
	Count int

	/* 输出参数中ret部分的切片 */
	List []QueryVideoIDorWatermarkIDList
}

type QueryVideoIDorWatermarkIDList struct {
	/* 上传文件对象名 */
	ObjectName string

	/* 视频ID */
	Vid uint64

	/* 水印图片ID */
	ImgId uint64
}
