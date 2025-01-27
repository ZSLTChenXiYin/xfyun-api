package face_match

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ZSLTChenXiYin/xfyun-api/basic_client"
)

const (
	DEFAULT_SILENT_LIVE_DETECTION_SENSETIME_REQUEST_ADDRESS = "https://api.xfyun.cn/v1/service/v1/image_identify/silent_detection"
)

type SilentLiveDetectionSensetimeResult struct {
	basic_client.XFYunAPICommonResult
	Data struct {
		Passed          bool    `json:"passed"`
		Liveness_score  float64 `json:"liveness_score"`
		Image_timestamp int     `json:"imagetimestamp"`
		Base64_image    string  `json:"base64_image"`
	}
}

type SilentLiveDetectionSensetimeClient struct {
	xfyun_api_basic_client basic_client.XFYunAPIHeaderVerificationClient

	app_id  string
	api_key string

	silent_live_detection_sensetime_request_form url.Values

	SilentLiveDetectionSensetimeResult *SilentLiveDetectionSensetimeResult
}

type SilentLiveDetectionSensetimeClientOption func(*SilentLiveDetectionSensetimeClient)

func WithSilentLiveDetectionSensetimeClientRequestConfiguration(request_address string) SilentLiveDetectionSensetimeClientOption {
	return func(sldsc *SilentLiveDetectionSensetimeClient) {
		sldsc.xfyun_api_basic_client.RequestAddress = request_address
	}
}

func WithSilentLiveDetectionSensetimeClientBasicConfiguration(app_id string, api_key string) SilentLiveDetectionSensetimeClientOption {
	return func(sldsc *SilentLiveDetectionSensetimeClient) {
		sldsc.app_id = app_id
		sldsc.api_key = api_key
	}
}

func NewSilentLiveDetectionSensetimeClient(options ...SilentLiveDetectionSensetimeClientOption) *SilentLiveDetectionSensetimeClient {
	sldsc := &SilentLiveDetectionSensetimeClient{}

	for _, option := range options {
		option(sldsc)
	}

	if sldsc.xfyun_api_basic_client.RequestAddress == "" {
		sldsc.xfyun_api_basic_client.RequestAddress = DEFAULT_SILENT_LIVE_DETECTION_SENSETIME_REQUEST_ADDRESS
	}

	sldsc.silent_live_detection_sensetime_request_form = make(url.Values)

	return sldsc
}

func (sldsc *SilentLiveDetectionSensetimeClient) SetRequestConfiguration(request_address string) *SilentLiveDetectionSensetimeClient {
	sldsc.xfyun_api_basic_client.RequestAddress = request_address

	return sldsc
}

func (sldsc *SilentLiveDetectionSensetimeClient) SetBasicConfiguration(app_id string, api_key string) *SilentLiveDetectionSensetimeClient {
	sldsc.app_id = app_id
	sldsc.api_key = api_key

	return sldsc
}

func (sldsc *SilentLiveDetectionSensetimeClient) Ready() error {
	if sldsc.app_id == "" || sldsc.api_key == "" {
		return errors.New("silent_live_detection_sensetime: The app_id, api_key is required")
	}

	return nil
}

func (sldsc *SilentLiveDetectionSensetimeClient) AddFile(file []byte) error {
	base64_file := base64.StdEncoding.EncodeToString(file)

	if len(base64_file) > 10*1024*1024 {
		return errors.New("face_match_sensetime: The base64 file is larger than 10M")
	}

	sldsc.silent_live_detection_sensetime_request_form.Set("file", base64_file)

	return nil
}

func (sldsc *SilentLiveDetectionSensetimeClient) Do(get_image bool) error {
	format_time := sldsc.xfyun_api_basic_client.NowCurTime()

	param_map := make(map[string]any)
	param_map["get_image"] = get_image

	param, err := sldsc.xfyun_api_basic_client.GetParam(param_map)
	if err != nil {
		return err
	}

	check_sum := sldsc.xfyun_api_basic_client.GetCheckSum(sldsc.api_key, format_time, param)

	client := &http.Client{}

	request, err := http.NewRequest("POST", sldsc.xfyun_api_basic_client.RequestAddress, strings.NewReader(sldsc.silent_live_detection_sensetime_request_form.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("X-Appid", sldsc.app_id)
	request.Header.Set("X-CurTime", format_time)
	request.Header.Set("X-Param", param)
	request.Header.Set("X-CheckSum", check_sum)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("silent_live_detection_sensetime: response.StatusCode == %d", response.StatusCode)
	}

	json_silent_live_detection_sensetime_response_body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	response_map := make(map[string]any)
	err = json.Unmarshal(json_silent_live_detection_sensetime_response_body, &response_map)
	if err != nil {
		return err
	}

	if response_map["code"] != "0" {
		return fmt.Errorf("silent_live_detection_sensetime: error_response == %v", response_map)
	}

	sldsc.SilentLiveDetectionSensetimeResult = &SilentLiveDetectionSensetimeResult{}
	err = json.Unmarshal(json_silent_live_detection_sensetime_response_body, sldsc.SilentLiveDetectionSensetimeResult)
	if err != nil {
		return err
	}

	return nil
}

func (sldsc *SilentLiveDetectionSensetimeClient) Flush() {
	sldsc.SilentLiveDetectionSensetimeResult = nil
}
