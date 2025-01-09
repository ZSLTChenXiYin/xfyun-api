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
	DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_ADDRESS      = "https://api.xf-yun.com/v1/private/s67c9c78c"
	DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_REQUEST_LINE = "POST /v1/private/s67c9c78c HTTP/1.1"
	DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_HOST         = "api.xf-yun.com"
)

type FacialFeatureAnalysisTupuechRequestBody struct {
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

func NewFacialFeatureAnalysisTupuechRequestBody() *FacialFeatureAnalysisTupuechRequestBody {
	ffatrb := &FacialFeatureAnalysisTupuechRequestBody{}

	ffatrb.Header.Status = 3

	ffatrb.Parameter.S67c9c78c.Service_kind = "face_detect"

	ffatrb.Parameter.S67c9c78c.Face_detect_result.Encoding = "utf8"
	ffatrb.Parameter.S67c9c78c.Face_detect_result.Compress = "raw"
	ffatrb.Parameter.S67c9c78c.Face_detect_result.Format = "json"

	ffatrb.Payload.Input1.Status = 3

	return ffatrb
}

func (ffatrb *FacialFeatureAnalysisTupuechRequestBody) SetAppId(app_id string) *FacialFeatureAnalysisTupuechRequestBody {
	ffatrb.Header.AppId = app_id
	return ffatrb
}

func (ffatrb *FacialFeatureAnalysisTupuechRequestBody) SetDetectPoints(detect_points bool) *FacialFeatureAnalysisTupuechRequestBody {
	if detect_points {
		ffatrb.Parameter.S67c9c78c.Detect_points = "1"
	} else {
		ffatrb.Parameter.S67c9c78c.Detect_points = "0"
	}
	return ffatrb
}

func (ffatrb *FacialFeatureAnalysisTupuechRequestBody) SetDetectProperty(detect_property bool) *FacialFeatureAnalysisTupuechRequestBody {
	if detect_property {
		ffatrb.Parameter.S67c9c78c.Detect_property = "1"
	} else {
		ffatrb.Parameter.S67c9c78c.Detect_property = "0"
	}
	return ffatrb
}

func (ffatrb *FacialFeatureAnalysisTupuechRequestBody) SetInput(encoding string, image string) *FacialFeatureAnalysisTupuechRequestBody {
	ffatrb.Payload.Input1.Encoding = encoding
	ffatrb.Payload.Input1.Image = image

	return ffatrb
}

type FacialFeatureAnalysisTupuechResult struct {
	Face_1   map[string]any `json:"face_1"`
	Face_num int            `json:"face_num"`
	Ret      int            `json:"ret"`
}

type FacialFeatureAnalysisTupuechResponseBody struct {
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

func (ffactc *FacialFeatureAnalysisTupuechResponseBody) GetFacialFeatureAnalysisTupuechResult() (*FacialFeatureAnalysisTupuechResult, error) {
	base64_text := ffactc.Payload.Face_detect_result.Text

	json_text, err := base64.StdEncoding.DecodeString(base64_text)
	if err != nil {
		return nil, err
	}

	ffacr := &FacialFeatureAnalysisTupuechResult{}

	err = json.Unmarshal(json_text, ffacr)
	if err != nil {
		return nil, err
	}

	return ffacr, nil
}

type FacialFeatureAnalysisTupuechClient struct {
	xfyun_api_basic_client basic_client.XFYunAPIURLVerificationClient

	app_id     string
	api_secret string
	api_key    string

	facial_feature_analysis_tupuech_request_body *FacialFeatureAnalysisTupuechRequestBody

	FacialFeatureAnalysisTupuechResponseBody *FacialFeatureAnalysisTupuechResponseBody
}

type FacialFeatureAnalysisTupuechClientOption func(*FacialFeatureAnalysisTupuechClient)

func WithFacialFeatureAnalysisTupuechClientRequestConfiguration(request_address string, request_line string, host string) FacialFeatureAnalysisTupuechClientOption {
	return func(ffactc *FacialFeatureAnalysisTupuechClient) {
		ffactc.xfyun_api_basic_client.RequestAddress = request_address
		ffactc.xfyun_api_basic_client.RequestLine = request_line
		ffactc.xfyun_api_basic_client.Host = host
	}
}

func WithFacialFeatureAnalysisTupuechClientBasicConfiguration(app_id string, api_secret string, api_key string) FacialFeatureAnalysisTupuechClientOption {
	return func(ffactc *FacialFeatureAnalysisTupuechClient) {
		ffactc.app_id = app_id
		ffactc.api_secret = api_secret
		ffactc.api_key = api_key
	}
}

func NewFacialFeatureAnalysisTupuechClient(options ...FacialFeatureAnalysisTupuechClientOption) *FacialFeatureAnalysisTupuechClient {
	ffactc := &FacialFeatureAnalysisTupuechClient{}

	for _, option := range options {
		option(ffactc)
	}

	if ffactc.xfyun_api_basic_client.RequestAddress == "" {
		ffactc.xfyun_api_basic_client.RequestAddress = DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_ADDRESS
	}

	if ffactc.xfyun_api_basic_client.RequestLine == "" {
		ffactc.xfyun_api_basic_client.RequestLine = DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_REQUEST_LINE
	}

	if ffactc.xfyun_api_basic_client.Host == "" {
		ffactc.xfyun_api_basic_client.Host = DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_HOST
	}

	ffactc.facial_feature_analysis_tupuech_request_body = NewFacialFeatureAnalysisTupuechRequestBody()

	ffactc.facial_feature_analysis_tupuech_request_body.SetAppId(ffactc.app_id).SetDetectPoints(false).SetDetectProperty(true)

	return ffactc
}

func (ffactc *FacialFeatureAnalysisTupuechClient) SetDetectConfiguration(detect_points bool, detect_property bool) *FacialFeatureAnalysisTupuechClient {
	ffactc.facial_feature_analysis_tupuech_request_body.SetDetectPoints(detect_points).SetDetectProperty(detect_property)

	return ffactc
}

func (ffactc *FacialFeatureAnalysisTupuechClient) SetRequestConfiguration(request_address string, request_line string, host string) *FacialFeatureAnalysisTupuechClient {
	ffactc.xfyun_api_basic_client.RequestAddress = request_address
	ffactc.xfyun_api_basic_client.RequestLine = request_line
	ffactc.xfyun_api_basic_client.Host = host

	return ffactc
}

func (ffactc *FacialFeatureAnalysisTupuechClient) SetBasicConfiguration(app_id string, api_secret string, api_key string) *FacialFeatureAnalysisTupuechClient {
	ffactc.app_id = app_id
	ffactc.api_secret = api_secret
	ffactc.api_key = api_key

	return ffactc
}

func (ffactc *FacialFeatureAnalysisTupuechClient) Ready() error {
	if ffactc.app_id == "" || ffactc.api_secret == "" || ffactc.api_key == "" {
		return errors.New("facial_feature_analysis_tupuech: The app_id, api_key is required")
	}

	return nil
}

func (ffactc *FacialFeatureAnalysisTupuechClient) AddInput(encoding string, image []byte) error {
	if encoding != "jpg" && encoding != "jpeg" && encoding != "png" && encoding != "bmp" {
		return errors.New("facial_feature_analysis_tupuech: The image type is not supported")
	}

	base64_image := base64.StdEncoding.EncodeToString(image)
	if len(base64_image) > 4*1024*1024 {
		return errors.New("facial_feature_analysis_tupuech: The base64 image is larger than 4M")
	}

	ffactc.facial_feature_analysis_tupuech_request_body.SetInput(encoding, base64_image)

	return nil
}

func (ffactc *FacialFeatureAnalysisTupuechClient) Do() error {
	format_date := ffactc.xfyun_api_basic_client.NowRFC1123()

	authorization, err := ffactc.xfyun_api_basic_client.GetAuthorization(format_date, ffactc.api_secret, ffactc.api_key)
	if err != nil {
		return err
	}

	request_url := ffactc.xfyun_api_basic_client.GetRequestURL(authorization, format_date)

	json_facial_feature_analysis_tupuech_request_body, err := json.Marshal(ffactc.facial_feature_analysis_tupuech_request_body)
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

	ffactc.FacialFeatureAnalysisTupuechResponseBody = &FacialFeatureAnalysisTupuechResponseBody{}
	err = json.Unmarshal(json_facial_feature_analysis_tupuech_response_body, ffactc.FacialFeatureAnalysisTupuechResponseBody)
	if err != nil {
		return err
	}

	return nil
}

func (ffactc *FacialFeatureAnalysisTupuechClient) Flush() {
	ffactc.FacialFeatureAnalysisTupuechResponseBody = nil
}
