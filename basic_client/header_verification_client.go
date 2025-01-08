package basic_client

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"
)

type XFYunAPICommonResult struct {
	Code string `json:"code"`
	Desc string `json:"desc"`
	Sid  string `json:"sid"`
}

// 讯飞云接口基础客户端
type XFYunAPIHeaderVerificationClient struct {
	RequestAddress string
	Host           string
}

func (xfyhvc *XFYunAPIHeaderVerificationClient) NowCurTime() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

func (xfyhvc *XFYunAPIHeaderVerificationClient) GetParam(param map[string]any) (string, error) {
	json_param, err := json.Marshal(param)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(json_param)), nil
}

func (xfyhvc *XFYunAPIHeaderVerificationClient) GetCheckSum(api_key string, cur_time string, param string) string {
	combined := api_key + cur_time + param

	hash := md5.Sum([]byte(combined))

	return hex.EncodeToString(hash[:])
}
