package basic_client

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
)

// 讯飞云接口通用请求体
type XFYunAPICommonRequestBody struct {
	Header struct {
		AppId  string `json:"app_id"`
		Status int    `json:"status"`
	} `json:"header"`
}

// 讯飞云接口通用响应体
type XFYunAPICommonResponseBody struct {
	Header struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		SID     string `json:"sid"`
	} `json:"header"`
}

// 讯飞云接口基础客户端
type XFYunAPIURLVerificationClient struct {
	RequestAddress string
	RequestLine    string
	Host           string
}

// NowRFC1123 返回当前时间的RFC1123格式字符串。
//
// 这个方法主要用于获取当前时间，并将其格式化为RFC1123标准的字符串格式。
// RFC1123日期格式常用于HTTP协议的头部字段中，以表示时间。
//
// 返回值:
//
//	string - RFC1123格式的时间字符串。
func (xfyac *XFYunAPIURLVerificationClient) NowRFC1123() string {
	return time.Now().UTC().Format(time.RFC1123)
}

// GetAuthorization 生成访问API所需的Authorization头信息。
//
// 此方法根据提供的日期、API密钥和API签名密钥生成授权信息。
// 它首先创建一个签名字符串，然后使用HMAC-SHA256算法进行签名，并将结果编码为Base64。
// 最后，它将所有必要信息编码进一个Authorization字符串中，同样使用Base64编码。
//
// 参数:
//
//	date - 用于签名的日期字符串，通常为RFC1123格式。
//	api_secret - 用于HMAC-SHA256签名的API密钥。
//	api_key - 用于生成Authorization头部的API密钥。
//
// 返回值:
//
//	authorization - 生成的Authorization字符串。
//	err - 如果签名过程中发生错误，返回该错误。
func (xfyac *XFYunAPIURLVerificationClient) GetAuthorization(date string, api_secret string, api_key string) (authorization string, err error) {
	// 生成signature的原始字段(signature_origin)
	signature_origin := fmt.Sprintf("host: %s\ndate: %s\n%s", xfyac.Host, date, xfyac.RequestLine)

	// 使用hmac-sha256算法结合apiSecret对signature_origin签名，获得签名后的摘要signature_sha
	hmac_sha256 := hmac.New(sha256.New, []byte(api_secret))
	_, err = hmac_sha256.Write([]byte(signature_origin))
	if err != nil {
		return "", err
	}
	signature_sha := hmac_sha256.Sum(nil)

	// 使用base64编码对signature_sha进行编码获得最终的signature
	signature := base64.StdEncoding.EncodeToString(signature_sha)

	// 生成authorization base64编码前（authorization_origin）的字符串
	authorization_origin := fmt.Sprintf("api_key=\"%s\",algorithm=\"hmac-sha256\",headers=\"host date request-line\",signature=\"%s\"", api_key, signature)

	// 对authorization_origin进行base64编码获得最终的authorization参数
	authorization = base64.StdEncoding.EncodeToString([]byte(authorization_origin))

	return authorization, nil
}

// GetRequestURL 构建最终的API请求URL。
//
// 该方法将Authorization字符串、日期和其他必要参数组合到一个URL中，
// 以便可以向API发起请求。
//
// 参数:
//
//	authorization - Base64编码的Authorization字符串。
//	date - 用于签名的日期字符串，通常为RFC1123格式。
//
// 返回值:
//
//	string - 完整的请求URL。
func (xfyac *XFYunAPIURLVerificationClient) GetRequestURL(authorization string, date string) string {
	return fmt.Sprintf("%s?authorization=%s&host=%s&date=%s", xfyac.RequestAddress, authorization, xfyac.Host, url.QueryEscape(date))
}
