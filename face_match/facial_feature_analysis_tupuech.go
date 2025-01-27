package face_match

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ZSLTChenXiYin/xfyun-api/basic_client"
)

const (
	DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_AGE_REQUEST_ADDRESS        = "http://tupapi.xfyun.cn/v1/age"
	DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_FACE_SCORE_REQUEST_ADDRESS = "http://tupapi.xfyun.cn/v1/face_score"
	DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_SEX_REQUEST_ADDRESS        = "http://tupapi.xfyun.cn/v1/sex"
	DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_EXPRESSION_REQUEST_ADDRESS = "http://tupapi.xfyun.cn/v1/expression"
)

type FacialFeatureAnalysisTupuechResult struct {
	basic_client.XFYunAPICommonResult
	Code int            `json:"code"`
	Data map[string]any `json:"data"`
}

type FacialFeatureAnalysisTupuechClient struct {
	xfyun_api_basic_client basic_client.XFYunAPIHeaderVerificationClient

	app_id  string
	api_key string

	facial_feature_analysis_tupuech_request_form []byte

	FacialFeatureAnalysisTupuechResult *FacialFeatureAnalysisTupuechResult
}

type FacialFeatureAnalysisTupuechClientOption func(*FacialFeatureAnalysisTupuechClient)

func WithFacialFeatureAnalysisTupuechClientRequestConfiguration(request_address string) FacialFeatureAnalysisTupuechClientOption {
	return func(fcatc *FacialFeatureAnalysisTupuechClient) {
		fcatc.xfyun_api_basic_client.RequestAddress = request_address
	}
}

func WithFacialFeatureAnalysisTupuechClientBasicConfiguration(app_id string, api_key string) FacialFeatureAnalysisTupuechClientOption {
	return func(fcatc *FacialFeatureAnalysisTupuechClient) {
		fcatc.app_id = app_id
		fcatc.api_key = api_key
	}
}

func NewFacialFeatureAnalysisTupuechClient(options ...FacialFeatureAnalysisTupuechClientOption) *FacialFeatureAnalysisTupuechClient {
	ffatc := &FacialFeatureAnalysisTupuechClient{}

	for _, option := range options {
		option(ffatc)
	}

	return ffatc
}

func (ffatc *FacialFeatureAnalysisTupuechClient) SetRequestConfiguration(request_address string) *FacialFeatureAnalysisTupuechClient {
	ffatc.xfyun_api_basic_client.RequestAddress = request_address
	return ffatc
}

func (ffatc *FacialFeatureAnalysisTupuechClient) SetBasicConfiguration(app_id string, api_key string) *FacialFeatureAnalysisTupuechClient {
	ffatc.app_id = app_id
	ffatc.api_key = api_key

	return ffatc
}

func (ffatc *FacialFeatureAnalysisTupuechClient) Ready() error {
	if ffatc.app_id == "" || ffatc.api_key == "" {
		return errors.New("facial_feature_analysis_tupuech: The app_id, api_key is required")
	}

	if ffatc.xfyun_api_basic_client.RequestAddress == "" {
		return errors.New("facial_feature_analysis_tupuech: The request_address is required")
	}

	return nil
}

func (ffatc *FacialFeatureAnalysisTupuechClient) AddFile(file []byte) error {
	if file == nil {
		return errors.New("facial_feature_analysis_tupuech: The file is required")
	}

	ffatc.facial_feature_analysis_tupuech_request_form = file

	return nil
}

func (ffatc *FacialFeatureAnalysisTupuechClient) Do(image_name string, image_url string) error {
	format_time := ffatc.xfyun_api_basic_client.NowCurTime()

	param_map := make(map[string]any)
	param_map["image_name"] = image_name
	if image_url != "" {
		param_map["image_url"] = image_url
	}

	param, err := ffatc.xfyun_api_basic_client.GetParam(param_map)
	if err != nil {
		return err
	}

	check_sum := ffatc.xfyun_api_basic_client.GetCheckSum(ffatc.api_key, format_time, param)

	client := &http.Client{}
	var request *http.Request
	if image_url == "" {
		request, err = http.NewRequest("POST", ffatc.xfyun_api_basic_client.RequestAddress, bytes.NewReader(ffatc.facial_feature_analysis_tupuech_request_form))
	} else {
		request, err = http.NewRequest("POST", ffatc.xfyun_api_basic_client.RequestAddress, nil)
	}
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/octet-stream")
	request.Header.Set("X-Appid", ffatc.app_id)
	request.Header.Set("X-CurTime", format_time)
	request.Header.Set("X-Param", param)
	request.Header.Set("X-CheckSum", check_sum)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("facial_feature_analysis_tupuech: " + response.Status)
	}

	json_facial_feature_analysis_tupuech_response_body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	response_map := make(map[string]any)
	err = json.Unmarshal(json_facial_feature_analysis_tupuech_response_body, &response_map)
	if err != nil {
		return err
	}

	if response_map["code"].(float64) != 0 {
		return fmt.Errorf("facial_feature_analysis_tupuech: error_response == %v", response_map)
	}

	ffatc.FacialFeatureAnalysisTupuechResult = &FacialFeatureAnalysisTupuechResult{}
	err = json.Unmarshal(json_facial_feature_analysis_tupuech_response_body, ffatc.FacialFeatureAnalysisTupuechResult)
	if err != nil {
		return err
	}

	return nil
}

func (ffatc *FacialFeatureAnalysisTupuechClient) Flush() {
	ffatc.FacialFeatureAnalysisTupuechResult = nil
}
