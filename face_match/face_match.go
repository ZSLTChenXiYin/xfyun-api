package face_match

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ZSLTChenXiYin/xfyun-api/basic_client"
)

const (
	DEFAULT_FACE_MATCH_REQUEST_ADDRESS = "https://api.xf-yun.com/v1/private/s67c9c78c"
	DEFAULT_FACE_MATCH_REQUEST_LINE    = "POST /v1/private/s67c9c78c HTTP/1.1"
	DEFAULT_FACE_MATCH_HOST            = "api.xf-yun.com"
)

// 人脸比对请求体
type FaceMatchRequestBody struct {
	basic_client.XFYunAPICommonRequestBody
	Parameter struct {
		S67c9c78c struct {
			Service_kind        string `json:"service_kind"`
			Face_compare_result struct {
				Encoding string `json:"encoding"`
				Compress string `json:"compress"`
				Format   string `json:"format"`
			} `json:"face_compare_result"`
		} `json:"s67c9c78c"`
	} `json:"parameter"`
	Payload struct {
		Input1 struct {
			Encoding string `json:"encoding"`
			Status   int    `json:"status"`
			Image    string `json:"image"`
		} `json:"input1"`
		Input2 struct {
			Encoding string `json:"encoding"`
			Status   int    `json:"status"`
			Image    string `json:"image"`
		} `json:"input2"`
	} `json:"payload"`
}

// NewFaceMatchRequestBody 创建并初始化一个新的 FaceMatchRequestBody 对象。
//
// 此函数用于设置面部对比请求的初始状态和默认参数。
//
// 返回值:
//
//	*FaceMatchRequestBody - 初始化后的 FaceMatchRequestBody 对象指针。
func NewFaceMatchRequestBody() *FaceMatchRequestBody {
	// 创建一个新的 FaceMatchRequestBody 实例。
	fmrb := &FaceMatchRequestBody{}

	// 设置请求头的状态码为3，表示请求处于初始化状态。
	fmrb.Header.Status = 3

	// 配置面部对比服务的固定参数。
	fmrb.Parameter.S67c9c78c.Service_kind = "face_compare"
	fmrb.Parameter.S67c9c78c.Face_compare_result.Encoding = "utf8"
	fmrb.Parameter.S67c9c78c.Face_compare_result.Compress = "raw"
	fmrb.Parameter.S67c9c78c.Face_compare_result.Format = "json"

	// 设置请求负载中两个输入数据的状态码为3，表示输入数据有效。
	fmrb.Payload.Input1.Status = 3
	fmrb.Payload.Input2.Status = 3

	// 返回初始化后的 FaceMatchRequestBody 对象。
	return fmrb
}

// SetAppId 设置FaceMatchRequestBody对象的AppId。
//
// 参数:
//
//	app_id - 要设置的应用ID。
//
// 返回值:
//
//	*FaceMatchRequestBody - 返回FaceMatchRequestBody对象，用于链式调用。
func (fmrb *FaceMatchRequestBody) SetAppId(app_id string) *FaceMatchRequestBody {
	fmrb.Header.AppId = app_id
	return fmrb
}

// SetInput1 设置FaceMatchRequestBody的Input1字段。
//
// 该方法接收两个参数：encoding和image，并将它们分别赋值给Input1的Encoding和Image字段。
//
// 参数:
//
//	encoding - 表示图像编码格式的字符串，如 jpg或png 等。
//	image - 与面部特征关联的图像数据。
//
// 返回值:
//
//	*FaceMatchRequestBody - 以便进行链式调用。
func (fmrb *FaceMatchRequestBody) SetInput1(encoding string, image string) *FaceMatchRequestBody {
	fmrb.Payload.Input1.Encoding = encoding
	fmrb.Payload.Input1.Image = image
	return fmrb
}

// SetInput2 设置FaceMatchRequestBody的Input2字段的Encoding和Image。
//
// 这个方法允许调用者设置与人脸识别相关的编码信息和图片信息。
//
// 参数:
//
//	encoding - 表示图像编码格式的字符串，如 jpg或png 等。
//	image - 与面部特征关联的图像数据。
//
// 返回值:
//
//	*FaceMatchRequestBody - 以便进行链式调用。
func (fmrb *FaceMatchRequestBody) SetInput2(encoding string, image string) *FaceMatchRequestBody {
	fmrb.Payload.Input2.Encoding = encoding
	fmrb.Payload.Input2.Image = image
	return fmrb
}

// 人脸比对结果
type FaceMatchResult struct {
	Ret   int     `json:"ret"`
	Score float64 `json:"score"`
}

// 人脸比对响应体
type FaceMatchResponseBody struct {
	basic_client.XFYunAPICommonResponseBody
	Payload struct {
		Face_compare_result struct {
			Compress string `json:"compress"`
			Encoding string `json:"encoding"`
			Format   string `json:"format"`
			Text     string `json:"text"`
		} `json:"face_compare_result"`
	} `json:"payload"`
}

// GetFaceCompareResult 获取人脸对比结果。
//
// 该方法解析FaceMatchResponseBody中的Face_compare_result字段，
// 将其从base64编码的字符串形式解码为JSON格式的文本，
// 然后将JSON文本反序列化为FaceMatchResult对象。
//
// 返回值:
//
//	*FaceMatchResult - 人脸对比结果对象，包含对比的相关信息。
//	error - 在解码或反序列化过程中遇到的错误，如果没有错误则为nil。
func (fmrb *FaceMatchResponseBody) GetFaceCompareResult() (*FaceMatchResult, error) {
	// 将base64编码的字符串解码为JSON格式的文本。
	base64_text := fmrb.Payload.Face_compare_result.Text
	json_text, err := base64.StdEncoding.DecodeString(base64_text)
	if err != nil {
		return nil, err
	}

	// 创建一个FaceMatchResult对象，并尝试将解码后的JSON文本解析到该对象中。
	rmr := &FaceMatchResult{}
	err = json.Unmarshal(json_text, rmr)
	if err != nil {
		return nil, err
	}

	// 返回解析后的人脸对比结果对象。
	return rmr, nil
}

// 讯飞云人脸比对客户端
type FaceMatchClient struct {
	// 讯飞云接口通用客户端
	xfyun_api_basic_client basic_client.XFYunAPIURLVerificationClient

	// 人脸比对服务申请的应用ID
	app_id string

	// 人脸比对服务申请的应用密钥
	api_secret string

	// 人脸比对服务申请的应用密钥
	api_key string

	// 人脸比对服务预制请求体
	face_match_request_body *FaceMatchRequestBody

	// 人脸比对响应体，每当正确调用人脸比对服务（即Do方法）时，都会向其中写入服务响应体
	FaceMatchResponseBody *FaceMatchResponseBody
}

type FaceMatchClientOption func(*FaceMatchClient)

// WithFaceMatchClientRequestConfiguration 是一个 FaceMatchClientOption 类型的函数，用于配置请求的地址和线路信息。
//
// 这个函数的主要作用是定制 FaceMatchClient 的请求配置，使其能够连接到指定的主机和请求地址。
//
// 参数:
//
//	request_address - 请求的地址。
//	request_line - 请求的线路。
//	host - 目标主机的地址。
//
// 返回值:
//
//	FaceMatchClientOption - 将配置应用到 FaceMatchClient 实例上。
func WithFaceMatchClientRequestConfiguration(request_address string, request_line string, host string) FaceMatchClientOption {
	return func(fmc *FaceMatchClient) {
		fmc.xfyun_api_basic_client.RequestAddress = request_address
		fmc.xfyun_api_basic_client.RequestLine = request_line
		fmc.xfyun_api_basic_client.Host = host
	}
}

// WithFaceMatchClientBasicConfiguration 返回一个配置了基本参数的FaceMatchClientOption。
// 这个函数主要用于简化FaceMatchClient的配置过程，通过接收app_id、api_secret和api_key作为参数，
// 生成一个闭包函数，该闭包函数将这些基本配置参数应用到FaceMatchClient实例上。
// 参数:
//
//	app_id - 应用程序的ID，用于标识和认证。
//	api_secret - API的密钥，用于服务器与服务器之间的安全通信。
//	api_key - 访问API的密钥，通常与api_secret一起使用进行认证。
//
// 返回值:
//
//	FaceMatchClientOption - 将配置应用到 FaceMatchClient 实例上。
func WithFaceMatchClientBasicConfiguration(app_id string, api_secret string, api_key string) FaceMatchClientOption {
	// 闭包函数将传入的配置参数应用到FaceMatchClient实例上。
	return func(fmc *FaceMatchClient) {
		fmc.app_id = app_id
		fmc.api_secret = api_secret
		fmc.api_key = api_key
	}
}

// NewFaceMatchClient 创建一个新的FaceMatchClient实例，并根据提供的选项进行配置。
//
// 该函数接受一系列FaceMatchClientOption作为参数，这些选项用于配置FaceMatchClient的行为。
//
// 参数:
//
//	options - 可变参数，包含FaceMatchClient的配置选项。
//
// 返回值:
//
//	*FaceMatchClient - 配置好的FaceMatchClient实例的指针。
func NewFaceMatchClient(options ...FaceMatchClientOption) *FaceMatchClient {
	// 初始化FaceMatchClient实例。
	fmc := &FaceMatchClient{}

	// 应用配置选项到FaceMatchClient实例。
	for _, option := range options {
		option(fmc)
	}

	// 如果未设置请求地址，则使用默认的请求地址。
	if fmc.xfyun_api_basic_client.RequestAddress == "" {
		fmc.xfyun_api_basic_client.RequestAddress = DEFAULT_FACE_MATCH_REQUEST_ADDRESS
	}

	// 如果未设置请求行，则使用默认的请求行。
	if fmc.xfyun_api_basic_client.RequestLine == "" {
		fmc.xfyun_api_basic_client.RequestLine = DEFAULT_FACE_MATCH_REQUEST_LINE
	}

	// 如果未设置主机名，则使用默认的主机名。
	if fmc.xfyun_api_basic_client.Host == "" {
		fmc.xfyun_api_basic_client.Host = DEFAULT_FACE_MATCH_HOST
	}

	// 初始化面部匹配请求体。
	fmc.face_match_request_body = NewFaceMatchRequestBody()

	// 设置请求体中的AppId。
	fmc.face_match_request_body.SetAppId(fmc.app_id)

	// 返回配置好的FaceMatchClient实例。
	return fmc
}

// SetRequestConfiguration 配置请求的地址和线路信息
//
// 该方法允许设置 FaceMatchClient 实例的请求配置信息，包括请求地址、请求线路和主机名。
// 这些配置信息对于建立与服务器的连接至关重要。
//
// 参数:
//
//	request_address - 请求的地址
//	request_line - 请求的线路
//	host - 主机名
//
// 返回值:
//
//	*FaceMatchClient - 返回配置后的 FaceMatchClient 实例，支持链式调用
func (fmc *FaceMatchClient) SetRequestConfiguration(request_address string, request_line string, host string) *FaceMatchClient {
	fmc.xfyun_api_basic_client.RequestAddress = request_address
	fmc.xfyun_api_basic_client.RequestLine = request_line
	fmc.xfyun_api_basic_client.Host = host
	return fmc
}

// SetBasicConfiguration 设置基本的认证信息
//
// 该方法用于设置 FaceMatchClient 实例的基本认证信息，包括应用ID、API密钥和API密钥。
// 这些信息是进行API调用时身份验证所必需的。
//
// 参数:
//
//	app_id - 应用ID
//	api_secret - API密钥
//	api_key - API密钥
//
// 返回值:
//
//	*FaceMatchClient - 返回配置后的 FaceMatchClient 实例，支持链式调用
func (fmc *FaceMatchClient) SetBasicConfiguration(app_id string, api_secret string, api_key string) *FaceMatchClient {
	fmc.app_id = app_id
	fmc.api_secret = api_secret
	fmc.api_key = api_key
	fmc.face_match_request_body.SetAppId(fmc.app_id)
	return fmc
}

// Ready 检查实例是否已配置必要的认证信息
//
// 该方法检查 FaceMatchClient 实例是否已设置必要的认证信息，即应用ID、API密钥和API密钥。
// 如果任何一项认证信息未被设置，则返回错误。
//
// 返回值:
//
//	error - 如果认证信息不完整，则返回错误；否则返回nil
func (fmc *FaceMatchClient) Ready() error {
	if fmc.app_id == "" || fmc.api_secret == "" || fmc.api_key == "" {
		return errors.New("face_match: The app_id, api_secret, api_key is required")
	}
	return nil
}

// AddInput 向FaceMatchClient添加输入数据。
//
// 该方法根据input_option选择设置人脸匹配请求的输入1或输入2。
//
// 参数:
//
//	input_option - 指定输入选项，1代表输入1，2代表输入2。
//	encoding - 图像的编码类型。
//	image - 图像的字节数据。
//
// 返回值:
//
//	error - 如果输入的图像数据的Base64编码长度超过4M，将返回错误。如果input_option不是有效的选项，将返回错误。
func (fmc *FaceMatchClient) AddInput(input_option uint, encoding string, image []byte) error {
	if encoding != "jpg" && encoding != "jpeg" && encoding != "png" && encoding != "bmp" {
		return errors.New("face_match: The image type is not supported")
	}

	// 将图像数据编码为Base64字符串。
	base64_image := base64.StdEncoding.EncodeToString(image)
	// 检查Base64编码后的图像数据大小是否超过4M。
	if len(base64_image) > 4*1024*1024 {
		return errors.New("face_match: The base64 image is larger than 4M")
	}

	// 根据input_option选择设置请求体的输入数据。
	switch input_option {
	case 1:
		fmc.face_match_request_body.SetInput1(encoding, base64_image)
	case 2:
		fmc.face_match_request_body.SetInput2(encoding, base64_image)
	default:
		// 如果input_option不是1或2，返回无效选项的错误。
		return errors.New("face_match: The option is invalid")
	}

	// 没有错误发生，返回nil。
	return nil
}

// AddAllInput 向FaceMatchClient添加所有输入数据。
//
// 该方法同时设置人脸匹配请求的输入1和输入2。
//
// 参数:
//
//	encoding1 - 第一幅图像的编码类型。
//	image1 - 第一幅图像的字节数据。
//	encoding2 - 第二幅图像的编码类型。
//	image2 - 第二幅图像的字节数据。
//
// 返回值:
//
//	error - 如果任一输入的图像数据的Base64编码长度超过4M，将返回错误。
func (fmc *FaceMatchClient) AddAllInput(encoding1 string, image1 []byte, encoding2 string, image2 []byte) error {
	if encoding1 != "jpg" && encoding1 != "jpeg" && encoding1 != "png" && encoding1 != "bmp" {
		return errors.New("face_match: The image type is not supported")
	}

	if encoding2 != "jpg" && encoding2 != "jpeg" && encoding2 != "png" && encoding2 != "bmp" {
		return errors.New("face_match: The image type is not supported")
	}

	base64_image1 := base64.StdEncoding.EncodeToString(image1)
	if len(base64_image1) > 4*1024*1024 {
		return errors.New("face_match: The base64 image is larger than 4M")
	}
	base64_image2 := base64.StdEncoding.EncodeToString(image2)
	if len(base64_image2) > 4*1024*1024 {
		return errors.New("face_match: The base64 image is larger than 4M")
	}

	fmc.face_match_request_body.SetInput1(encoding1, base64_image1).SetInput2(encoding2, base64_image2)

	return nil
}

// Do 执行人脸匹配请求。
//
// 该方法构造请求、发送请求并处理响应。
//
// 返回值:
//
//	error - 如果在构造请求、发送请求或处理响应过程中发生错误，将返回错误。
func (fmc *FaceMatchClient) Do() error {
	format_date := fmc.xfyun_api_basic_client.NowRFC1123()

	authorization, err := fmc.xfyun_api_basic_client.GetAuthorization(format_date, fmc.api_secret, fmc.api_key)
	if err != nil {
		return err
	}

	request_url := fmc.xfyun_api_basic_client.GetRequestURL(authorization, format_date)

	json_face_match_request_body, err := json.Marshal(fmc.face_match_request_body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", request_url, bytes.NewBuffer(json_face_match_request_body))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	json_face_match_response_body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("face_match: %s", json_face_match_response_body)
	}

	fmc.FaceMatchResponseBody = &FaceMatchResponseBody{}
	err = json.Unmarshal(json_face_match_response_body, fmc.FaceMatchResponseBody)
	if err != nil {
		return err
	}

	return nil
}

// Flush 清空FaceMatchClient的响应体。
//
// 该方法用于重置FaceMatchClient，以便进行新的请求。
func (fmc *FaceMatchClient) Flush() {
	fmc.FaceMatchResponseBody = nil
}
