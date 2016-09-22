package vcloudutil

/*
 * Description : 进行上传初始化需要的参数信息
 * Company : NetEase
 * Date : 2016-8-11
 */

type UploadInitParams struct {
	/** 上传文件的原始名称（包含后缀名） 此参数必填*/
	OriginFileName string

	/** 用户命名的上传文件名称  此参数非必填*/
	UserFileName string

	/** 视频所属的类别ID（不填写为默认分类）此参数非必填*/
	TypeId int64

	/** 视频所需转码模板ID（不填写为默认模板） 此参数非必填*/
	PresetId int64

	/** 转码成功后回调客户端的URL地址（需标准http格式）  此参数非必填*/
	CallbackUrl string

	/** 上传视频的描述信息  此参数非必填*/
	Description string

	/** 上传视频的视频水印Id  此参数非必填*/
	WatermarkId int64

	/** 上传成功后回调客户端的URL地址（需标准http格式） 此参数非必填 */
	UploadCallbackUrl string

	/** 用户自定义信息，会在上传成功或转码成功后通过回调返回给用户  此参数非必填*/
	UserDefInfo string
}
