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
	DEFAULT_SILENT_LIVE_DETECTION_REQUEST_ADDRESS = "https://api.xf-yun.com/v1/private/s67c9c78c"
	DEFAULT_SILENT_LIVE_DETECTION_REQUEST_LINE    = "POST /v1/private/s67c9c78c HTTP/1.1"
	DEFAULT_SILENT_LIVE_DETECTION_HOST            = "api.xf-yun.com"
)

// 静默活体检测请求体
type SilentLiveDetectionRequestBody struct {
	basic_client.XFYunAPICommonRequestBody
	Parameter struct {
		S67c9c78c struct {
			Service_kind      string `json:"service_kind"`
			Anti_spoof_result struct {
				Encoding string `json:"encoding"`
				Compress string `json:"compress"`
				Format   string `json:"format"`
			} `json:"anti_spoof_result"`
		} `json:"s67c9c78c"`
	} `json:"parameter"`
	Payload struct {
		Input1 struct {
			Encoding string `json:"encoding"`
			Status   int    `json:"status"`
			Image    string `json:"image"`
		} `json:"input1"`
	} `json:"payload"`
}

// NewSilentLiveDetectionRequestBody 创建并初始化一个新的 SilentLiveDetectionRequestBody 实例。
//
// 该函数主要用于设置默认的配置参数，为进行活体检测服务做准备。
//
// 返回值:
//
//	*SilentLiveDetectionRequestBody - 初始化后的 SilentLiveDetectionRequestBody 对象指针。
func NewSilentLiveDetectionRequestBody() *SilentLiveDetectionRequestBody {
	// 创建一个新的 SilentLiveDetectionRequestBody 实例。
	sldrb := &SilentLiveDetectionRequestBody{}

	// 设置请求头的状态码为3，表示请求正在处理中。
	sldrb.Header.Status = 3

	// 配置参数部分，指定服务类型为反欺骗服务。
	sldrb.Parameter.S67c9c78c.Service_kind = "anti_spoof"
	// 设置人脸识别结果的编码格式为 UTF-8。
	sldrb.Parameter.S67c9c78c.Anti_spoof_result.Encoding = "utf8"
	// 设置人脸识别结果的压缩格式为原始格式，即不压缩。
	sldrb.Parameter.S67c9c78c.Anti_spoof_result.Compress = "raw"
	// 设置人脸识别结果的数据格式为 JSON。
	sldrb.Parameter.S67c9c78c.Anti_spoof_result.Format = "json"

	// 设置负载部分的状态码为3，表示负载正在处理中。
	sldrb.Payload.Input1.Status = 3

	// 返回初始化完毕的 SilentLiveDetectionRequestBody 实例。
	return sldrb
}

// SetAppId 设置 SilentLiveDetectionRequestBody 实例的 AppId。
//
// 参数:
//
//	app_id - 要设置的 AppId 字符串。
//
// 返回值:
//
//	*SilentLiveDetectionRequestBody - 设置 AppId 后的 SilentLiveDetectionRequestBody 对象指针。
func (sldrb *SilentLiveDetectionRequestBody) SetAppId(app_id string) *SilentLiveDetectionRequestBody {
	sldrb.Header.AppId = app_id
	return sldrb
}

// SetInput 设置 SilentLiveDetectionRequestBody 实例的输入参数。
//
// 参数:
//
//	encoding - 输入图像的编码格式。
//	image - 输入的图像数据。
//
// 返回值:
//
//	*SilentLiveDetectionRequestBody - 设置输入参数后的 SilentLiveDetectionRequestBody 对象指针。
func (sldrb *SilentLiveDetectionRequestBody) SetInput(encoding string, image string) *SilentLiveDetectionRequestBody {
	sldrb.Payload.Input1.Encoding = encoding
	sldrb.Payload.Input1.Image = image
	return sldrb
}

// 静默活体检测结果
type SilentLiveDetectionResult struct {
	H      int     `json:"h"`
	Passed bool    `json:"passed"`
	Ret    int     `json:"ret"`
	Score  float64 `json:"score"`
	W      int     `json:"w"`
	X      int     `json:"x"`
	Y      int     `json:"y"`
}

// 静默活体检测响应体
type SilentLiveDetectionResponseBody struct {
	basic_client.XFYunAPICommonResponseBody
	Payload struct {
		Anti_spoof_result struct {
			Compress string `json:"compress"`
			Encoding string `json:"encoding"`
			Format   string `json:"format"`
			Text     string `json:"text"`
		} `json:"anti_spoof_result"`
	} `json:"payload"`
}

// GetAntiSpoofResult 获取并解析防伪结果
//
// 该方法从SilentLiveDetectionResponseBody的Payload中提取防伪结果（anti-spoof result），
// 将其从Base64编码的字符串解码为JSON格式，并解析到SilentLiveDetectionResult结构体中。
//
// 返回值:
//
//	*SilentLiveDetectionResult - 解析后的SilentLiveDetectionResult对象。
//	error - 如果在Base64解码或JSON解析过程中发生错误，则返回错误信息。
func (sldrb *SilentLiveDetectionResponseBody) GetAntiSpoofResult() (*SilentLiveDetectionResult, error) {
	// 提取Base64编码的防伪结果文本
	base64_text := sldrb.Payload.Anti_spoof_result.Text

	// 将Base64编码的文本解码为JSON格式的字节切片
	json_text, err := base64.StdEncoding.DecodeString(base64_text)
	if err != nil {
		// 如果解码过程中出现错误，则返回nil和错误信息
		return nil, err
	}

	// 初始化SilentLiveDetectionResult对象
	sldr := &SilentLiveDetectionResult{}

	// 将解码后的JSON文本解析到SilentLiveDetectionResult对象中
	err = json.Unmarshal(json_text, sldr)
	if err != nil {
		// 如果JSON解析过程中出现错误，则返回nil和错误信息
		return nil, err
	}

	// 返回解析后的SilentLiveDetectionResult对象和nil错误
	return sldr, nil
}

// 静默活体检测客户端
type SilentLiveDetectionClient struct {
	// 讯飞云接口基础客户端
	xfyun_api_basic_client basic_client.XFYunAPIURLVerificationClient

	// 应用ID
	app_id string

	// API密钥
	api_secret string

	// API密钥
	api_key string

	// 静默活体检测请求体
	silent_live_detection_request_body *SilentLiveDetectionRequestBody

	// 静默活体检测响应体
	SilentLiveDetectionResponseBody *SilentLiveDetectionResponseBody
}

type SilentLiveDetectionClientOption func(*SilentLiveDetectionClient)

// WithSilentLiveDetectionClientRequestConfiguration 返回一个配置静默活体检测客户端请求的选项。
//
// 该选项允许设置请求地址、请求行和主机头。
//
// 参数:
//
//	request_address - 请求的地址。
//	request_line - 请求行，通常包括HTTP方法和路径。
//	host - 请求的主机名。
//
// 返回值:
//
//	SilentLiveDetectionClientOption - 用于配置SilentLiveDetectionClient。
func WithSilentLiveDetectionClientRequestConfiguration(request_address string, request_line string, host string) SilentLiveDetectionClientOption {
	return func(sldc *SilentLiveDetectionClient) {
		sldc.xfyun_api_basic_client.RequestAddress = request_address
		sldc.xfyun_api_basic_client.RequestLine = request_line
		sldc.xfyun_api_basic_client.Host = host
	}
}

// WithSilentLiveDetectionClientBasicConfiguration 返回一个配置静默活体检测客户端基本参数的选项。
//
// 该选项允许设置应用ID、API密钥和API密钥。
//
// 参数:
//
//	app_id - 应用的唯一标识符。
//	api_secret - API的密钥，用于服务器身份验证。
//	api_key - API的密钥，用于客户端身份验证。
//
// 返回值:
//
//	SilentLiveDetectionClientOption - 用于配置SilentLiveDetectionClient。
func WithSilentLiveDetectionClientBasicConfiguration(app_id string, api_secret string, api_key string) SilentLiveDetectionClientOption {
	return func(sldc *SilentLiveDetectionClient) {
		sldc.app_id = app_id
		sldc.api_secret = api_secret
		sldc.api_key = api_key
	}
}

// NewSilentLiveDetectionClient 创建并返回一个新的SilentLiveDetectionClient实例。
//
// 它会应用默认的请求地址、请求行和主机名，如果这些值没有被其他选项覆盖的话。
//
// 参数:
//
//	options - 可选的SilentLiveDetectionClientOption，用于配置SilentLiveDetectionClient。
//
// 返回值:
//
//	*SilentLiveDetectionClient - 配置好的SilentLiveDetectionClient实例。
func NewSilentLiveDetectionClient(options ...SilentLiveDetectionClientOption) *SilentLiveDetectionClient {
	// 初始化SilentLiveDetectionClient实例。
	sldc := &SilentLiveDetectionClient{}

	// 应用每个配置选项到新创建的客户端实例上。
	for _, option := range options {
		option(sldc)
	}

	// 确保请求地址已设置，若未设置则使用默认值。
	if sldc.xfyun_api_basic_client.RequestAddress == "" {
		sldc.xfyun_api_basic_client.RequestAddress = DEFAULT_SILENT_LIVE_DETECTION_REQUEST_ADDRESS
	}

	// 确保请求行已设置，若未设置则使用默认值。
	if sldc.xfyun_api_basic_client.RequestLine == "" {
		sldc.xfyun_api_basic_client.RequestLine = DEFAULT_SILENT_LIVE_DETECTION_REQUEST_LINE
	}

	// 确保主机名已设置，若未设置则使用默认值。
	if sldc.xfyun_api_basic_client.Host == "" {
		sldc.xfyun_api_basic_client.Host = DEFAULT_SILENT_LIVE_DETECTION_HOST
	}

	// 初始化用于静默活体检测请求的请求体。
	sldc.silent_live_detection_request_body = NewSilentLiveDetectionRequestBody()

	// 设置请求体中的AppId。
	sldc.silent_live_detection_request_body.SetAppId(sldc.app_id)

	// 返回配置完毕的SilentLiveDetectionClient实例。
	return sldc
}

// SetRequestConfiguration 设置请求的配置信息。
//
// 该方法通过对SilentLiveDetectionClient实例进行配置，以使其能够根据指定的请求地址、请求行和主机名进行请求。
//
// 参数:
//
//	request_address - 请求地址。
//	request_line - 请求行。
//	host - 主机名。
//
// 返回值:
//
//	*SilentLiveDetectionClient - 配置好的SilentLiveDetectionClient实例，支持链式调用。
func (sldc *SilentLiveDetectionClient) SetRequestConfiguration(request_address string, request_line string, host string) *SilentLiveDetectionClient {
	sldc.xfyun_api_basic_client.RequestAddress = request_address
	sldc.xfyun_api_basic_client.RequestLine = request_line
	sldc.xfyun_api_basic_client.Host = host
	return sldc
}

// SetBasicConfiguration 设置基础配置信息。
//
// 该方法用于设置SilentLiveDetectionClient实例的基础配置信息，包括app_id、api_secret和api_key，并设置请求体的AppId。
//
// 参数:
//
//	app_id - 应用ID。
//	api_secret - API密钥。
//	api_key - API密钥。
//
// 返回值:
//
//	*SilentLiveDetectionClient - 配置好的SilentLiveDetectionClient实例，支持链式调用。
func (sldc *SilentLiveDetectionClient) SetBasicConfiguration(app_id string, api_secret string, api_key string) *SilentLiveDetectionClient {
	sldc.app_id = app_id
	sldc.api_secret = api_secret
	sldc.api_key = api_key
	sldc.silent_live_detection_request_body.SetAppId(app_id)
	return sldc
}

// Ready 检查配置是否完整。
//
// 该方法检查SilentLiveDetectionClient实例的配置信息是否完整，确保app_id、api_secret和api_key均已设置。
//
// 返回值:
//
//	error - 如果配置信息不完整，则返回错误；否则返回nil。
func (sldc *SilentLiveDetectionClient) Ready() error {
	if sldc.app_id == "" || sldc.api_secret == "" || sldc.api_key == "" {
		return errors.New("silent_live_detection: The app_id, api_secret, api_key is required")
	}
	return nil
}

// AddInput 添加输入图像及其编码信息。
//
// 该方法将给定的图像以Base64编码后，连同其编码格式一起添加到静默活体检测请求的输入数据中。
// 它确保了图像数据的大小不超过4MB，以符合接口处理能力。
//
// 参数:
//
//	encoding - 图像的编码格式。
//	image - 图像的字节数据。
//
// 返回值:
//
//	error - 如果图像数据超过4MB，则返回错误；否则返回nil。
func (sldc *SilentLiveDetectionClient) AddInput(encoding string, image []byte) error {
	if encoding != "jpg" && encoding != "jpeg" && encoding != "png" && encoding != "bmp" {
		return errors.New("silent_live_detection: The image type is not supported")
	}

	base64_image := base64.StdEncoding.EncodeToString(image)
	if len(base64_image) > 4*1024*1024 {
		return errors.New("silent_live_detection: The base64 image is larger than 4M")
	}

	sldc.silent_live_detection_request_body.SetInput(encoding, base64_image)

	return nil
}

// Do 执行静默活体检测请求。
//
// 该方法负责组装请求、发送请求并处理响应。它首先生成请求的授权头，然后构造请求URL和请求体，
// 发送POST请求到服务器，并解析响应。如果响应状态码不是200，它会返回错误。
//
// 返回值:
//
//	error - 如果在执行过程中遇到错误，则返回错误；否则返回nil。
func (sldc *SilentLiveDetectionClient) Do() error {
	// 获取当前时间，并格式化为RFC1123标准，用于请求头中的Date字段。
	format_date := sldc.xfyun_api_basic_client.NowRFC1123()

	// 根据当前时间、API密钥和API私钥生成请求的授权信息。
	authorization, err := sldc.xfyun_api_basic_client.GetAuthorization(format_date, sldc.api_secret, sldc.api_key)
	if err != nil {
		return err
	}

	// 构造请求URL，包含授权信息和时间戳。
	request_url := sldc.xfyun_api_basic_client.GetRequestURL(authorization, format_date)

	// 将请求体序列化为JSON格式。
	json_silent_live_detection_request_body, err := json.Marshal(sldc.silent_live_detection_request_body)
	if err != nil {
		return err
	}

	// 创建HTTP POST请求。
	request, err := http.NewRequest("POST", request_url, bytes.NewBuffer(json_silent_live_detection_request_body))
	if err != nil {
		return err
	}

	// 设置请求头，指定内容类型和主机。
	request.Header.Set("Content-Type", "application/json")

	// 发送HTTP请求。
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// 读取HTTP响应体。
	json_silent_live_detection_response_body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// 检查HTTP状态码，非200则视为错误。
	if response.StatusCode != 200 {
		return fmt.Errorf("silent_live_detection: %s", json_silent_live_detection_response_body)
	}

	// 初始化静默活体检测响应体对象，并反序列化JSON响应体。
	sldc.SilentLiveDetectionResponseBody = &SilentLiveDetectionResponseBody{}
	err = json.Unmarshal(json_silent_live_detection_response_body, sldc.SilentLiveDetectionResponseBody)
	if err != nil {
		return err
	}

	// 所有操作成功完成，返回nil。
	return nil
}

// Flush 清空静默活体检测的响应体。
//
// 该方法将静默活体检测的响应体设置为nil，以便下一次检测可以重新使用该客户端。
func (sldc *SilentLiveDetectionClient) Flush() {
	sldc.SilentLiveDetectionResponseBody = nil
}
