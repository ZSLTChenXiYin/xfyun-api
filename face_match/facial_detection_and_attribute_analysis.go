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
	DEFAULT_FACIAL_DETECTION_AND_ATTRIBUTE_ANALYSIS_REQUEST_ADDRESS = "https://api.xf-yun.com/v1/private/s67c9c78c"
	DEFAULT_FACIAL_DETECTION_AND_ATTRIBUTE_ANALYSIS_REQUEST_LINE    = "POST /v1/private/s67c9c78c HTTP/1.1"
	DEFAULT_FACIAL_DETECTION_AND_ATTRIBUTE_ANALYSIS_HOST            = "api.xf-yun.com"
)

type FacialDetectionAndAttributeAnalysisRequestBody struct {
	basic_client.XFYunAPICommonRequestBody
	Parameter struct {
		S67c9c78c struct {
			Service_kind       string `json:"service_kind"`
			Face_detect_result struct {
				Encoding string `json:"encoding"`
				Compress string `json:"compress"`
				Format   string `json:"format"`
			} `json:"face_detect_result"`
			Detect_points   string `json:"detect_points"`
			Detect_property string `json:"detect_property"`
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

func NewFacialDetectionAndAttributeAnalysisRequestBody() *FacialDetectionAndAttributeAnalysisRequestBody {
	ffatrb := &FacialDetectionAndAttributeAnalysisRequestBody{}

	ffatrb.Header.Status = 3

	ffatrb.Parameter.S67c9c78c.Service_kind = "face_detect"

	ffatrb.Parameter.S67c9c78c.Face_detect_result.Encoding = "utf8"
	ffatrb.Parameter.S67c9c78c.Face_detect_result.Compress = "raw"
	ffatrb.Parameter.S67c9c78c.Face_detect_result.Format = "json"

	ffatrb.Payload.Input1.Status = 3

	return ffatrb
}

func (fdaaarb *FacialDetectionAndAttributeAnalysisRequestBody) SetAppId(app_id string) *FacialDetectionAndAttributeAnalysisRequestBody {
	fdaaarb.Header.AppId = app_id
	return fdaaarb
}

func (fdaaarb *FacialDetectionAndAttributeAnalysisRequestBody) SetDetectPoints(detect_points bool) *FacialDetectionAndAttributeAnalysisRequestBody {
	if detect_points {
		fdaaarb.Parameter.S67c9c78c.Detect_points = "1"
	} else {
		fdaaarb.Parameter.S67c9c78c.Detect_points = "0"
	}
	return fdaaarb
}

func (fdaaarb *FacialDetectionAndAttributeAnalysisRequestBody) SetDetectProperty(detect_property bool) *FacialDetectionAndAttributeAnalysisRequestBody {
	if detect_property {
		fdaaarb.Parameter.S67c9c78c.Detect_property = "1"
	} else {
		fdaaarb.Parameter.S67c9c78c.Detect_property = "0"
	}
	return fdaaarb
}

func (fdaaarb *FacialDetectionAndAttributeAnalysisRequestBody) SetInput(encoding string, image string) *FacialDetectionAndAttributeAnalysisRequestBody {
	fdaaarb.Payload.Input1.Encoding = encoding
	fdaaarb.Payload.Input1.Image = image

	return fdaaarb
}

type FacialDetectionAndAttributeAnalysisResult struct {
	Face_1   map[string]any `json:"face_1"`
	Face_num int            `json:"face_num"`
	Ret      int            `json:"ret"`
}

type FacialDetectionAndAttributeAnalysisResponseBody struct {
	basic_client.XFYunAPICommonResponseBody
	Payload struct {
		Face_detect_result struct {
			Compress string `json:"compress"`
			Encoding string `json:"encoding"`
			Format   string `json:"format"`
			Text     string `json:"text"`
		} `json:"face_detect_result"`
	} `json:"payload"`
}

func (fdaaarb *FacialDetectionAndAttributeAnalysisResponseBody) GetFacialDetectionAndAttributeAnalysisResult() (*FacialDetectionAndAttributeAnalysisResult, error) {
	base64_text := fdaaarb.Payload.Face_detect_result.Text

	json_text, err := base64.StdEncoding.DecodeString(base64_text)
	if err != nil {
		return nil, err
	}

	ffacr := &FacialDetectionAndAttributeAnalysisResult{}

	err = json.Unmarshal(json_text, ffacr)
	if err != nil {
		return nil, err
	}

	return ffacr, nil
}

type FacialDetectionAndAttributeAnalysisClient struct {
	xfyun_api_basic_client basic_client.XFYunAPIURLVerificationClient

	app_id     string
	api_secret string
	api_key    string

	facial_detection_and_attribute_analysis_request_body *FacialDetectionAndAttributeAnalysisRequestBody

	FacialDetectionAndAttributeAnalysisResponseBody *FacialDetectionAndAttributeAnalysisResponseBody
}

type FacialDetectionAndAttributeAnalysisClientOption func(*FacialDetectionAndAttributeAnalysisClient)

func WithFacialDetectionAndAttributeAnalysisClientRequestConfiguration(request_address string, request_line string, host string) FacialDetectionAndAttributeAnalysisClientOption {
	return func(fdaaac *FacialDetectionAndAttributeAnalysisClient) {
		fdaaac.xfyun_api_basic_client.RequestAddress = request_address
		fdaaac.xfyun_api_basic_client.RequestLine = request_line
		fdaaac.xfyun_api_basic_client.Host = host
	}
}

func WithFacialDetectionAndAttributeAnalysisClientBasicConfiguration(app_id string, api_secret string, api_key string) FacialDetectionAndAttributeAnalysisClientOption {
	return func(fdaaac *FacialDetectionAndAttributeAnalysisClient) {
		fdaaac.app_id = app_id
		fdaaac.api_secret = api_secret
		fdaaac.api_key = api_key
	}
}

func NewFacialDetectionAndAttributeAnalysisClient(options ...FacialDetectionAndAttributeAnalysisClientOption) *FacialDetectionAndAttributeAnalysisClient {
	fdaaac := &FacialDetectionAndAttributeAnalysisClient{}

	for _, option := range options {
		option(fdaaac)
	}

	if fdaaac.xfyun_api_basic_client.RequestAddress == "" {
		fdaaac.xfyun_api_basic_client.RequestAddress = DEFAULT_FACIAL_DETECTION_AND_ATTRIBUTE_ANALYSIS_REQUEST_ADDRESS
	}

	if fdaaac.xfyun_api_basic_client.RequestLine == "" {
		fdaaac.xfyun_api_basic_client.RequestLine = DEFAULT_FACIAL_DETECTION_AND_ATTRIBUTE_ANALYSIS_REQUEST_LINE
	}

	if fdaaac.xfyun_api_basic_client.Host == "" {
		fdaaac.xfyun_api_basic_client.Host = DEFAULT_FACIAL_DETECTION_AND_ATTRIBUTE_ANALYSIS_HOST
	}

	fdaaac.facial_detection_and_attribute_analysis_request_body = NewFacialDetectionAndAttributeAnalysisRequestBody()

	fdaaac.facial_detection_and_attribute_analysis_request_body.SetAppId(fdaaac.app_id).SetDetectPoints(false).SetDetectProperty(true)

	return fdaaac
}

func (fdaaac *FacialDetectionAndAttributeAnalysisClient) SetDetectConfiguration(detect_points bool, detect_property bool) *FacialDetectionAndAttributeAnalysisClient {
	fdaaac.facial_detection_and_attribute_analysis_request_body.SetDetectPoints(detect_points).SetDetectProperty(detect_property)

	return fdaaac
}

func (fdaaac *FacialDetectionAndAttributeAnalysisClient) SetRequestConfiguration(request_address string, request_line string, host string) *FacialDetectionAndAttributeAnalysisClient {
	fdaaac.xfyun_api_basic_client.RequestAddress = request_address
	fdaaac.xfyun_api_basic_client.RequestLine = request_line
	fdaaac.xfyun_api_basic_client.Host = host

	return fdaaac
}

func (fdaaac *FacialDetectionAndAttributeAnalysisClient) SetBasicConfiguration(app_id string, api_secret string, api_key string) *FacialDetectionAndAttributeAnalysisClient {
	fdaaac.app_id = app_id
	fdaaac.api_secret = api_secret
	fdaaac.api_key = api_key

	return fdaaac
}

func (fdaaac *FacialDetectionAndAttributeAnalysisClient) Ready() error {
	if fdaaac.app_id == "" || fdaaac.api_secret == "" || fdaaac.api_key == "" {
		return errors.New("facial_feature_analysis_tupuech: The app_id, api_key is required")
	}

	return nil
}

func (fdaaac *FacialDetectionAndAttributeAnalysisClient) AddInput(encoding string, image []byte) error {
	if encoding != "jpg" && encoding != "jpeg" && encoding != "png" && encoding != "bmp" {
		return errors.New("facial_feature_analysis_tupuech: The image type is not supported")
	}

	base64_image := base64.StdEncoding.EncodeToString(image)
	if len(base64_image) > 4*1024*1024 {
		return errors.New("facial_feature_analysis_tupuech: The base64 image is larger than 4M")
	}

	fdaaac.facial_detection_and_attribute_analysis_request_body.SetInput(encoding, base64_image)

	return nil
}

func (fdaaac *FacialDetectionAndAttributeAnalysisClient) Do() error {
	format_date := fdaaac.xfyun_api_basic_client.NowRFC1123()

	authorization, err := fdaaac.xfyun_api_basic_client.GetAuthorization(format_date, fdaaac.api_secret, fdaaac.api_key)
	if err != nil {
		return err
	}

	request_url := fdaaac.xfyun_api_basic_client.GetRequestURL(authorization, format_date)

	json_facial_feature_analysis_tupuech_request_body, err := json.Marshal(fdaaac.facial_detection_and_attribute_analysis_request_body)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", request_url, bytes.NewBuffer(json_facial_feature_analysis_tupuech_request_body))
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

	json_facial_feature_analysis_tupuech_response_body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return fmt.Errorf("facial_feature_analysis_tupuech: %s", json_facial_feature_analysis_tupuech_response_body)
	}

	fdaaac.FacialDetectionAndAttributeAnalysisResponseBody = &FacialDetectionAndAttributeAnalysisResponseBody{}
	err = json.Unmarshal(json_facial_feature_analysis_tupuech_response_body, fdaaac.FacialDetectionAndAttributeAnalysisResponseBody)
	if err != nil {
		return err
	}

	return nil
}

func (fdaaac *FacialDetectionAndAttributeAnalysisClient) Flush() {
	fdaaac.FacialDetectionAndAttributeAnalysisResponseBody = nil
}
