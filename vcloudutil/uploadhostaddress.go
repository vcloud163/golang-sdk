package vcloudutil

/*
 * Description : 获取上传加速节点返回参数
 * Company : NetEase
 * Date : 2016-8-16
 */

type UploadHostAddress struct {
	/* 返回参数中的lbs */
	Lbs string

	/* 返回参数中上传节点地址集合 */
	Upload []string
}
