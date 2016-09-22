package auth

/*
 * Description : 访问凭证
 * Company : NetEase
 * Date : 2016-8-11
 */

import (
	"crypto/sha1"
	"fmt"
	"io"
)

type Key struct {
	AccessKey string
	SecretKey string
}

/* 构造访问凭证对象 */
func New(accessKey string, secretKey string) *Key {
	key := Key{accessKey, secretKey}
	return &key
}

/* 设置accessKey */
func (key *Key) SetAccessKey(accessKey string) {
	key.AccessKey = accessKey
}

/* 设置secretKey */
func (key *Key) SetSecretKey(secretKey string) {
	key.SecretKey = secretKey
}

/* 获取accessKey */
func (key *Key) GetAccessKey() string {
	return key.AccessKey
}

/* 获取secretKey */
func (key *Key) GetSecretKey() string {
	return key.SecretKey
}

/* 将appSecret、 nonce、 curTime 按照SHA-1加密算法生成 checkSum */
func (key *Key) GetCheckSum(accessKey string, nonce string, curTime string) string {
	value := accessKey + nonce + curTime
	ec := encode(value)
	return ec
}

/* 对字符串进行SHA1哈希 */
func encode(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
