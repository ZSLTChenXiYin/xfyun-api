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
	DEFAULT_COOPERATIVE_LIVENESS_DETECTION_REQUEST_ADDRESS = "https://api.xf-yun.com/v1/private/s67c9c78c"
	DEFAULT_COOPERATIVE_LIVENESS_DETECTION_REQUEST_LINE    = "POST /v1/private/s67c9c78c HTTP/1.1"
	DEFAULT_COOPERATIVE_LIVENESS_DETECTION_HOST            = "api.xf-yun.com"
)

// 人脸比对请求体
type CooperativeLivenessDetectionRequestBody struct {
	basic_client.XFYunAPICommonRequestBody
	Parameter struct {
		S67c9c78c struct {
			Service_kind       string `json:"service_kind"`
			Face_status_result struct {
				Encoding string `json:"encoding"`
				Compress string `json:"compress"`
				Format   string `json:"format"`
			} `json:"face_status_result"`
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

func NewCooperativeLivenessDetectionRequestBody() *CooperativeLivenessDetectionRequestBody {
	cldrb := &CooperativeLivenessDetectionRequestBody{}

	cldrb.Header.Status = 3

	cldrb.Parameter.S67c9c78c.Service_kind = "face_status"
	cldrb.Parameter.S67c9c78c.Face_status_result.Encoding = "utf8"
	cldrb.Parameter.S67c9c78c.Face_status_result.Compress = "raw"
	cldrb.Parameter.S67c9c78c.Face_status_result.Format = "json"

	cldrb.Payload.Input1.Status = 3

	return cldrb
}

func (cldrb *CooperativeLivenessDetectionRequestBody) SetAppId(app_id string) *CooperativeLivenessDetectionRequestBody {
	cldrb.Header.AppId = app_id

	return cldrb
}

func (cldrb *CooperativeLivenessDetectionRequestBody) SetInput(encoding string, image string) *CooperativeLivenessDetectionRequestBody {
	cldrb.Payload.Input1.Encoding = encoding
	cldrb.Payload.Input1.Image = image

	return cldrb
}

type CooperativeLivenessDetectionResult struct {
	Ret      int `json:"ret"`
	Face_num int `json:"face_num"`
	Face_1   struct {
		Ret              int     `json:"ret"`
		X                int     `json:"x"`
		Y                int     `json:"y"`
		W                int     `json:"w"`
		H                int     `json:"h"`
		Eye_status       string  `json:"eye_status"`
		Eye_status_score float64 `json:"eye_status_score"`
		Eye_threshold    string  `json:"eye_threshold"`
	} `json:"face_1"`
}

type CooperativeLivenessDetectionResponseBody struct {
	basic_client.XFYunAPICommonResponseBody
	Payload struct {
		Face_status_result struct {
			Compress string `json:"compress"`
			Encoding string `json:"encoding"`
			Format   string `json:"format"`
			Text     string `json:"text"`
		} `json:"face_status_result"`
	}
}

func (cldrb *CooperativeLivenessDetectionResponseBody) GetCooperativeLivenessDetectionResult() (*CooperativeLivenessDetectionResult, error) {
	base64_text := cldrb.Payload.Face_status_result.Text

	json_text, err := base64.StdEncoding.DecodeString(base64_text)
	if err != nil {
		return nil, err
	}

	cldr := &CooperativeLivenessDetectionResult{}

	err = json.Unmarshal(json_text, cldr)
	if err != nil {
		return nil, err
	}

	return cldr, nil
}

type CooperativeLivenessDetectionClient struct {
	xfyun_api_basic_client basic_client.XFYunAPIURLVerificationClient

	app_id string

	api_secret string

	api_key string

	cooperative_liveness_detection_request_body *CooperativeLivenessDetectionRequestBody

	CooperativeLivenessDetectionResponseBody *CooperativeLivenessDetectionResponseBody
}

type CooperativeLivenessDetectionClientOption func(*CooperativeLivenessDetectionClient)

func WithCooperativeLivenessDetectionClientRequestConfiguration(request_address string, request_line string, host string) CooperativeLivenessDetectionClientOption {
	return func(cldc *CooperativeLivenessDetectionClient) {
		cldc.xfyun_api_basic_client.RequestAddress = request_address
		cldc.xfyun_api_basic_client.RequestLine = request_line
		cldc.xfyun_api_basic_client.Host = host
	}
}

func WithCooperativeLivenessDetectionClientBasicConfiguration(app_id string, api_secret string, api_key string) CooperativeLivenessDetectionClientOption {
	return func(cldc *CooperativeLivenessDetectionClient) {
		cldc.app_id = app_id
		cldc.api_secret = api_secret
		cldc.api_key = api_key
	}
}

func NewCooperativeLivenessDetectionClient(options ...CooperativeLivenessDetectionClientOption) *CooperativeLivenessDetectionClient {
	cldc := &CooperativeLivenessDetectionClient{}

	for _, option := range options {
		option(cldc)
	}

	if cldc.xfyun_api_basic_client.RequestAddress == "" {
		cldc.xfyun_api_basic_client.RequestAddress = DEFAULT_COOPERATIVE_LIVENESS_DETECTION_REQUEST_ADDRESS
	}

	if cldc.xfyun_api_basic_client.RequestLine == "" {
		cldc.xfyun_api_basic_client.RequestLine = DEFAULT_COOPERATIVE_LIVENESS_DETECTION_REQUEST_LINE
	}

	if cldc.xfyun_api_basic_client.Host == "" {
		cldc.xfyun_api_basic_client.Host = DEFAULT_COOPERATIVE_LIVENESS_DETECTION_HOST
	}

	cldc.cooperative_liveness_detection_request_body = NewCooperativeLivenessDetectionRequestBody()

	cldc.cooperative_liveness_detection_request_body.SetAppId(cldc.app_id)

	return cldc
}

func (cldc *CooperativeLivenessDetectionClient) SetRequestConfiguration(request_address string, request_line string, host string) *CooperativeLivenessDetectionClient {
	cldc.xfyun_api_basic_client.RequestAddress = request_address
	cldc.xfyun_api_basic_client.RequestLine = request_line
	cldc.xfyun_api_basic_client.Host = host

	return cldc
}

func (cldc *CooperativeLivenessDetectionClient) SetBasicConfiguration(app_id string, api_secret string, api_key string) *CooperativeLivenessDetectionClient {
	cldc.app_id = app_id
	cldc.api_secret = api_secret
	cldc.api_key = api_key

	cldc.cooperative_liveness_detection_request_body.SetAppId(cldc.app_id)

	return cldc
}

func (cldc *CooperativeLivenessDetectionClient) Ready() error {
	if cldc.app_id == "" || cldc.api_secret == "" || cldc.api_key == "" {
		return errors.New("cooperative_liveness_detection: The app_id, api_key is required")
	}

	return nil
}

func (cldc *CooperativeLivenessDetectionClient) AddInput(encoding string, image []byte) error {
	if encoding != "jpg" && encoding != "jpeg" && encoding != "png" && encoding != "bmp" {
		return errors.New("cooperative_liveness_detection: The image type is not supported")
	}

	base64_image := base64.StdEncoding.EncodeToString(image)
	if len(base64_image) > 4*1024*1024 {
		return errors.New("cooperative_liveness_detection: The base64 image is larger than 4M")
	}

	cldc.cooperative_liveness_detection_request_body.SetInput(encoding, base64_image)

	return nil
}

func (cldc *CooperativeLivenessDetectionClient) Do() error {
	format_date := cldc.xfyun_api_basic_client.NowRFC1123()

	authorization, err := cldc.xfyun_api_basic_client.GetAuthorization(format_date, cldc.api_secret, cldc.api_key)
	if err != nil {
		return err
	}

	request_url := cldc.xfyun_api_basic_client.GetRequestURL(authorization, format_date)

	json_cooperative_liveness_detection_request_body, err := json.Marshal(cldc.cooperative_liveness_detection_request_body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", request_url, bytes.NewBuffer(json_cooperative_liveness_detection_request_body))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("host", cldc.xfyun_api_basic_client.Host)
	request.Header.Set("app_id", cldc.app_id)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	json_cooperative_liveness_detection_response_body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("cooperative_liveness_detection: %s", json_cooperative_liveness_detection_response_body)
	}

	cldc.CooperativeLivenessDetectionResponseBody = &CooperativeLivenessDetectionResponseBody{}
	err = json.Unmarshal(json_cooperative_liveness_detection_response_body, cldc.CooperativeLivenessDetectionResponseBody)
	if err != nil {
		return err
	}

	return nil
}

func (cldc *CooperativeLivenessDetectionClient) Flush() {
	cldc.CooperativeLivenessDetectionResponseBody = nil
}
