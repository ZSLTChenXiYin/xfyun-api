package face_match

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ZSLTChenXiYin/xfyun-api/basic_client"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

const (
	DEFAULT_FACE_MATCH_SENSETIME_REQUEST_ADDRESS = "https://api.xfyun.cn/v1/service/v1/image_identify/face_verification"
	DEFAULT_FACE_MATCH_SENSETIME_HOST            = "api.xf-yun.com"
)

type FaceMatchSensetimeResult struct {
	basic_client.XFYunAPICommonResult
	Data float64 `json:"data"`
}

type FaceMatchSensetimeClient struct {
	xfyun_api_basic_client basic_client.XFYunAPIHeaderVerificationClient

	app_id  string
	api_key string

	face_match_request_form url.Values

	FaceMatchSensetimeResult *FaceMatchSensetimeResult
}

type FaceMatchSensetimeClientOption func(*FaceMatchSensetimeClient)

func WithFaceMatchSensetimeClientRequestConfiguration(request_address string, host string) FaceMatchSensetimeClientOption {
	return func(fmc *FaceMatchSensetimeClient) {
		fmc.xfyun_api_basic_client.RequestAddress = request_address
		fmc.xfyun_api_basic_client.Host = host
	}
}

func WithFaceMatchSensetimeClientBasicConfiguration(app_id string, api_key string) FaceMatchSensetimeClientOption {
	return func(fmc *FaceMatchSensetimeClient) {
		fmc.app_id = app_id
		fmc.api_key = api_key
	}
}

func NewFaceMatchSensetimeClient(options ...FaceMatchSensetimeClientOption) *FaceMatchSensetimeClient {
	fmsc := &FaceMatchSensetimeClient{}

	for _, option := range options {
		option(fmsc)
	}

	if fmsc.xfyun_api_basic_client.RequestAddress == "" {
		fmsc.xfyun_api_basic_client.RequestAddress = DEFAULT_FACE_MATCH_SENSETIME_REQUEST_ADDRESS
	}

	if fmsc.xfyun_api_basic_client.Host == "" {
		fmsc.xfyun_api_basic_client.Host = DEFAULT_FACE_MATCH_SENSETIME_HOST
	}

	fmsc.face_match_request_form = make(url.Values)

	return fmsc
}

func (fmsc *FaceMatchSensetimeClient) SetRequestConfiguration(request_address string, host string) *FaceMatchSensetimeClient {
	fmsc.xfyun_api_basic_client.RequestAddress = request_address
	fmsc.xfyun_api_basic_client.Host = host
	return fmsc
}

func (fmsc *FaceMatchSensetimeClient) SetBasicConfiguration(app_id string, api_key string) *FaceMatchSensetimeClient {
	fmsc.app_id = app_id
	fmsc.api_key = api_key

	return fmsc
}

func (fmsc *FaceMatchSensetimeClient) Ready() error {
	if fmsc.app_id == "" || fmsc.api_key == "" {
		return errors.New("face_match_sensetime: The app_id, api_key is required")
	}

	return nil
}

func (fmsc *FaceMatchSensetimeClient) AddInput(input_option uint, img []byte) error {
	dimg, img_type, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return err
	}

	if img_type != "jpg" && img_type != "jpeg" && img_type != "png" && img_type != "gif" && img_type != "bmp" && img_type != "tiff" {
		return errors.New("face_match_sensetime: The image type is not supported")
	}

	bounds := dimg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width < 8 || height < 8 {
		return errors.New("face_match_sensetime: The image is too small")
	}

	if width > 4000 || height > 4000 {
		return errors.New("face_match_sensetime: The image is too large")
	}

	base64_image := base64.StdEncoding.EncodeToString(img)

	if len(base64_image) > 5*1024*1024 {
		return errors.New("face_match_sensetime: The base64 image is larger than 5M")
	}

	switch input_option {
	case 1:
		fmsc.face_match_request_form.Set("first_image", base64_image)
	case 2:
		fmsc.face_match_request_form.Set("second_image", base64_image)
	default:
		return errors.New("face_match_sensetime: The option is invalid")
	}

	return nil
}

func (fmsc *FaceMatchSensetimeClient) AddAllInput(img1 []byte, img2 []byte) error {
	dimg, img_type, err := image.Decode(bytes.NewReader(img1))
	if err != nil {
		return err
	}

	if img_type != "jpg" && img_type != "jpeg" && img_type != "png" && img_type != "gif" && img_type != "bmp" && img_type != "tiff" {
		return errors.New("face_match_sensetime: The image type is not supported")
	}

	bounds := dimg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width < 8 || height < 8 {
		return errors.New("face_match_sensetime: The image is too small")
	}

	if width > 4000 || height > 4000 {
		return errors.New("face_match_sensetime: The image is too large")
	}

	base64_image1 := base64.StdEncoding.EncodeToString(img1)

	if len(base64_image1) > 5*1024*1024 {
		return errors.New("face_match_sensetime: The base64 image is larger than 5M")
	}

	dimg, img_type, err = image.Decode(bytes.NewReader(img2))
	if err != nil {
		return err
	}

	if img_type != "jpg" && img_type != "jpeg" && img_type != "png" && img_type != "gif" && img_type != "bmp" && img_type != "tiff" {
		return errors.New("face_match_sensetime: The image type is not supported")
	}

	bounds = dimg.Bounds()
	width = bounds.Dx()
	height = bounds.Dy()

	if width < 8 || height < 8 {
		return errors.New("face_match_sensetime: The image is too small")
	}

	if width > 4000 || height > 4000 {
		return errors.New("face_match_sensetime: The image is too large")
	}

	base64_image2 := base64.StdEncoding.EncodeToString(img2)

	if len(base64_image2) > 5*1024*1024 {
		return errors.New("face_match_sensetime: The base64 image is larger than 5M")
	}

	fmsc.face_match_request_form.Set("first_image", base64_image1)
	fmsc.face_match_request_form.Set("second_image", base64_image2)

	return nil
}

func (fmsc *FaceMatchSensetimeClient) Do(auto_rotate bool) error {
	format_time := fmsc.xfyun_api_basic_client.NowCurTime()

	param_map := make(map[string]any)
	param_map["auto_rotate"] = auto_rotate

	param, err := fmsc.xfyun_api_basic_client.GetParam(param_map)
	if err != nil {
		return err
	}

	check_sum := fmsc.xfyun_api_basic_client.GetCheckSum(fmsc.api_key, format_time, param)

	client := &http.Client{}

	request, err := http.NewRequest("POST", fmsc.xfyun_api_basic_client.RequestAddress, strings.NewReader(fmsc.face_match_request_form.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	request.Header.Set("X-Appid", fmsc.app_id)
	request.Header.Set("X-CurTime", format_time)
	request.Header.Set("X-Param", param)
	request.Header.Set("X-CheckSum", check_sum)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("face_match_sensetime: response.StatusCode == %d", response.StatusCode)
	}

	json_face_match_sensetime_response_body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmsc.FaceMatchSensetimeResult = &FaceMatchSensetimeResult{}
	err = json.Unmarshal(json_face_match_sensetime_response_body, fmsc.FaceMatchSensetimeResult)
	if err != nil {
		return err
	}

	return nil
}

func (fmsc *FaceMatchSensetimeClient) Flush() {
	fmsc.FaceMatchSensetimeResult = nil
}
