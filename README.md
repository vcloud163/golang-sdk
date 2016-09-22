# golang-SDK 说明

## 1 简介

golang-SDK 是用于服务器端点播上传的软件开发工具包，提供简单、便捷的方法，方便用户开发上传视频或图片文件的功能。

## 2 功能特性

- 文件上传
- 获取进度
- 断点续传
- 查询视频
- 设置回调

## 3 开发准备

### 3.1 下载地址

[golang sdk 的源码地址](https://github.com/vcloud163/golang-sdk.git "golang sdk 的源码地址")

### 3.2 环境准备

- 适用于 go1.7及以上版本。
- 通过管理控制台->账户信息获取AppKey和AppSecret；
- 下载 golang sdk，如果安装了git命令行，执行 git clone https://github.com/vcloud/golang-sdk.git或者直接在github下载zip包。
- 参照 API 说明和 sdk 中提供的 demo，开发代码。


### 3.3 https支持

默认使用https协议，如需修改为http协议，请在sdk包中onfig/configuration.go中修改。

## 4 使用说明

### 4.1 初始化

接入视频云点播，需要拥有一对有效的 AppKey 和 AppSecret 进行签名认证，可通过如下步骤获得：

- 开通视频云点播服务；
- 登陆视频云开发者平台，通过管理控制台->账户信息获取 AppKey 和 AppSecret。

在获取到 AppKey 和 AppSecret 之后，可按照如下方式进行初始化，AppKey和AppSecret配置在config/configuration.go文件中：

    import (
	"golang-sdk/auth"
	"golang-sdk/config"
	"golang-sdk/vcloudutil"
    )
	key := auth.Key{config.AppKey, config.AppSecret}
	vutil := vcloudutil.NewVcloudUtil(key)

### 4.2 文件上传

视频云点播在全国各地覆盖大量上传节点，会选择适合用户的最优节点进行文件上传，并根据用户传入的参数做不同处理，具体详见点播服务端 API 文档。

以下是使用示例：

	key := auth.Key{"", ""}
	vutil := vcloudutil.NewVcloudUtil(key)

	initParams := vcloudutil.UploadInitParams{}
	initParams.OriginFileName = util.GetFileName(filePath)
	initParams.UploadCallbackUrl = "http://127.0.0.1/abc"
	vutil.Upload = vcloudutil.NewUploadUtil(initParams)
	
    /* 上传初始化，获取xNosToken（上传token）、bucket（存储对象的桶名）、object（生成的唯一对象名） */
	responseDate := vutil.Upload.InitUpload(config.InitUploadVideoURL, key)
	if responseDate.Code != 200 {
		fmt.Println("Failed to init upload! Error message: " + responseDate.Msg)
		return
	}
	
    /* 获取上传加速节点地址 */
	uploadHost := vutil.Upload.GetUploadHost(config.GetUploadHostURL, responseDate.Ret.Bucket)
	from vcloud import Client
	client = Client(appKey, secretKey)
	
    /* 简单文件上传 */
	vutil.Upload.UploadVideo(responseDate, uploadHost, filePath)


**注：具体使用示例详见 sdk 包中 demo.go 文件的 UploadDemos 函数。**

### 4.3 查询进度

视频云点播文件上传采用分片处理，上传进度保存在用户自定义的本地文件中，可通过调用函数GetUploadProcess查询上传进度。

以下是使用示例：
	
	vutil := vcloudutil.NewVcloudUtilNoKey()
	/* "D:\\recoder.txt"为记录上传进度的本地文件 */
	offset := vutil.GetUploadProcess("D:\\recoder.txt")


**注：具体使用示例详见 sdk 包中 demo.go 文件的 GetUploadProcessDemo 函数。**

### 4.4 断点续传

在上传文件中，视频云点播通过唯一标识 context 标识正在上传的文件，可通过此标识获取到已经上传视频云的文件字节数。通过此方法可实现文件的断点续传。

为防止服务中止造成文件上传信息丢失，可通过在本地存储文件信息来记录断点信息，当服务重启启动，可根据文件继续上传文件。临时文件会在上传完成后删除记录。

以下是使用示例：
    
    /* recordFilePath为本地保存断点信息的文件路径 */
	recoder := vcloudutil.NewRecoder(recordFilePath)
	/* filePath为上传文件路径 */
	object := vutil.Upload.UploadVideoWithRecorder(filePath, recoder, key)
	
**注：具体使用示例详见 sdk 包中 demo.go 文件的 UploadWithRecoderDemo 函数。**

### 4.5 查询视频

视频上传成功后，可根据视频上传后的返回结果，主动查询视频唯一标识，支持批量查询。

以下是使用示例：

	/* 上传完成后，根据objectName查询视频ID */
	var objectNames []string
	objectNames = append(objectNames, responseDate.Ret.Object)
	queryRslt := vutil.Upload.QueryVideoIDorWatermarkID(config.QueryVideoIDURL, key, objectNames)

**注：具体使用示例详见 sdk 包中 demo.go 文件的 UploadDemos函数或 UploadWithRecoderDemo 函数。**

### 4.6 设置回调

如果设置回调，视频上传成功后会发送相关视频信息给回调接口。

以下是使用示例：

	ret := vutil.SetUploadCallback(config.SetUploadCallback, callback)
	fmt.Println(ret.Code)
	if ret.Code != 200 {
		fmt.Println(ret.Msg)
	}

**注：具体使用示例详见 sdk 包中 demo.go 文件的 SetUploadCallbackDemo 函数。**

## 5 版本更新记录

**v1.0.0**

1. golang SDK 的初始版本，提供点播上传的基本功能。包括：文件上传、获取进度、断点续传、查询视频、设置回调。
