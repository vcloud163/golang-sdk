package config

/*
 * Description : 公共参数配置
 * Company : NetEase
 * Date : 2016-8-11
 */

/* 访问凭证AppKey */
const AppKey string = "f4b51648c54e47777744653885cdd29c"

/* 访问凭证AppSecret */
const AppSecret string = "a7224d72e2e4403777770004da882fe6"

/* 分片上传视频：分片字节数（默认4M不变） */
const UploadDataSize int64 = 1024 * 1024 * 4

/* 视频上传初始化URL */
const InitUploadVideoURL string = "https://vcloud.163.com/app/vod/upload/init"

/* 获取上传加速节点URL */
const GetUploadHostURL string = "http://wanproxy.127.net/lbs?version=1.0"

/* 查询视频主ID的URL */
const QueryVideoIDURL string = "https://vcloud.163.com/app/vod/video/query"

/* 设置上传回调地址的URL */
const SetUploadCallback string = "https://vcloud.163.com/app/vod/upload/setcallback"
