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
	DEFAULT_FACIAL_WATERMARK_PHOTO_MATCH_REQUEST_ADDRESS = "https://api.xfyun.cn/v1/service/v1/image_identify/watermark_verification"
)

type FacialWatermarkPhotoMatchResult struct {
	basic_client.XFYunAPICommonResult
	Data float64 `json:"data"`
}

type FacialWatermarkPhotoMatchClient struct {
	xfyun_api_basic_client basic_client.XFYunAPIHeaderVerificationClient

	app_id  string
	api_key string

	facial_watermark_photo_match_request_form url.Values

	FacialWatermarkPhotoMatchResult *FacialWatermarkPhotoMatchResult
}

type FacialWatermarkPhotoMatchClientOption func(*FacialWatermarkPhotoMatchClient)

func WithFacialWatermarkPhotoMatchClientRequestConfiguration(request_address string) FacialWatermarkPhotoMatchClientOption {
	return func(fwpmc *FacialWatermarkPhotoMatchClient) {
		fwpmc.xfyun_api_basic_client.RequestAddress = request_address
	}
}

func WithFacialWatermarkPhotoMatchClientBasicConfiguration(app_id string, api_key string) FacialWatermarkPhotoMatchClientOption {
	return func(fwpmc *FacialWatermarkPhotoMatchClient) {
		fwpmc.app_id = app_id
		fwpmc.api_key = api_key
	}
}

func NewFacialWatermarkPhotoMatchClient(options ...FacialWatermarkPhotoMatchClientOption) *FacialWatermarkPhotoMatchClient {
	fwpmc := &FacialWatermarkPhotoMatchClient{}

	for _, option := range options {
		option(fwpmc)
	}

	if fwpmc.xfyun_api_basic_client.RequestAddress == "" {
		fwpmc.xfyun_api_basic_client.RequestAddress = DEFAULT_FACIAL_WATERMARK_PHOTO_MATCH_REQUEST_ADDRESS
	}

	fwpmc.facial_watermark_photo_match_request_form = make(url.Values)

	return fwpmc
}

func (fwpmc *FacialWatermarkPhotoMatchClient) SetRequestConfiguration(request_address string, host string) *FacialWatermarkPhotoMatchClient {
	fwpmc.xfyun_api_basic_client.RequestAddress = request_address

	return fwpmc
}

func (fwpmc *FacialWatermarkPhotoMatchClient) SetBasicConfiguration(app_id string, api_key string) *FacialWatermarkPhotoMatchClient {
	fwpmc.app_id = app_id
	fwpmc.api_key = api_key

	return fwpmc
}

func (fwpmc *FacialWatermarkPhotoMatchClient) Ready() error {
	if fwpmc.app_id == "" || fwpmc.api_key == "" {
		return errors.New("facial_watermark_photo_match: The app_id, api_key is required")
	}

	return nil
}

func (fwpmc *FacialWatermarkPhotoMatchClient) AddFaceImage(face_image []byte) error {
	if len(face_image) > 5*1024*1024 {
		return errors.New("facial_watermark_photo_match: The image is too big")
	}

	dimg, img_type, err := image.Decode(bytes.NewReader(face_image))
	if err != nil {
		return err
	}

	if img_type != "jpg" && img_type != "jpeg" && img_type != "png" && img_type != "gif" && img_type != "bmp" && img_type != "tiff" {
		return errors.New("facial_watermark_photo_match: The image type is not supported")
	}

	bounds := dimg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width < 8 || height < 8 {
		return errors.New("facial_watermark_photo_match: The image is too small")
	}

	if width > 4000 || height > 4000 {
		return errors.New("facial_watermark_photo_match: The image is too large")
	}

	base64_image := base64.StdEncoding.EncodeToString(face_image)

	fwpmc.facial_watermark_photo_match_request_form.Set("face_image", base64_image)

	return nil
}

func (fwpmc *FacialWatermarkPhotoMatchClient) AddWatermarkImage(watermark_image []byte) error {
	if len(watermark_image) > 1*1024*1024 {
		return errors.New("facial_watermark_photo_match: The image is too big")
	}

	dimg, img_type, err := image.Decode(bytes.NewReader(watermark_image))
	if err != nil {
		return err
	}

	if img_type != "jpg" && img_type != "jpeg" && img_type != "png" && img_type != "gif" && img_type != "bmp" && img_type != "tiff" {
		return errors.New("facial_watermark_photo_match: The image type is not supported")
	}

	bounds := dimg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width < 8 || height < 8 {
		return errors.New("facial_watermark_photo_match: The image is too small")
	}

	if width > 480 || height > 480 {
		return errors.New("facial_watermark_photo_match: The image is too large")
	}

	base64_image := base64.StdEncoding.EncodeToString(watermark_image)

	fwpmc.facial_watermark_photo_match_request_form.Set("watermark_image", base64_image)

	return nil
}

func (fwpmc *FacialWatermarkPhotoMatchClient) AddAllImage(face_image []byte, watermark_image []byte) error {
	if len(face_image) > 5*1024*1024 {
		return errors.New("facial_watermark_photo_match: The face_image is too big")
	}

	if len(watermark_image) > 1*1024*1024 {
		return errors.New("facial_watermark_photo_match: The watermark_image is too big")
	}

	dimg, img_type, err := image.Decode(bytes.NewReader(face_image))
	if err != nil {
		return err
	}

	if img_type != "jpg" && img_type != "jpeg" && img_type != "png" && img_type != "gif" && img_type != "bmp" && img_type != "tiff" {
		return errors.New("facial_watermark_photo_match: The image type is not supported")
	}

	bounds := dimg.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width < 8 || height < 8 {
		return errors.New("facial_watermark_photo_match: The image is too small")
	}

	if width > 4000 || height > 4000 {
		return errors.New("facial_watermark_photo_match: The image is too large")
	}

	dimg, img_type, err = image.Decode(bytes.NewReader(watermark_image))
	if err != nil {
		return err
	}

	if img_type != "jpg" && img_type != "jpeg" && img_type != "png" && img_type != "gif" && img_type != "bmp" && img_type != "tiff" {
		return errors.New("facial_watermark_photo_match: The image type is not supported")
	}

	bounds = dimg.Bounds()
	width = bounds.Dx()
	height = bounds.Dy()

	if width < 8 || height < 8 {
		return errors.New("facial_watermark_photo_match: The image is too small")
	}

	if width > 480 || height > 480 {
		return errors.New("facial_watermark_photo_match: The image is too large")
	}

	fwpmc.facial_watermark_photo_match_request_form.Set("face_image", base64.StdEncoding.EncodeToString(face_image))
	fwpmc.facial_watermark_photo_match_request_form.Set("watermark_image", base64.StdEncoding.EncodeToString(watermark_image))

	return nil
}

func (fwpmc *FacialWatermarkPhotoMatchClient) Do(auto_rotate bool) error {
	format_time := fwpmc.xfyun_api_basic_client.NowCurTime()

	param_map := make(map[string]any)
	param_map["auto_rotate"] = auto_rotate

	param, err := fwpmc.xfyun_api_basic_client.GetParam(param_map)
	if err != nil {
		return err
	}

	check_sum := fwpmc.xfyun_api_basic_client.GetCheckSum(fwpmc.api_key, format_time, param)

	client := &http.Client{}

	request, err := http.NewRequest("POST", fwpmc.xfyun_api_basic_client.RequestAddress, strings.NewReader(fwpmc.facial_watermark_photo_match_request_form.Encode()))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	request.Header.Set("X-Appid", fwpmc.app_id)
	request.Header.Set("X-CurTime", format_time)
	request.Header.Set("X-Param", param)
	request.Header.Set("X-CheckSum", check_sum)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("facial_watermark_photo_match: response.StatusCode == %d", response.StatusCode)
	}

	json_facial_watermark_photo_match_response_body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	response_map := make(map[string]any)
	err = json.Unmarshal(json_facial_watermark_photo_match_response_body, &response_map)
	if err != nil {
		return err
	}

	if response_map["code"] != "0" {
		return fmt.Errorf("facial_watermark_photo_match: error_response == %v", json_facial_watermark_photo_match_response_body)
	}

	fwpmc.FacialWatermarkPhotoMatchResult = &FacialWatermarkPhotoMatchResult{}
	err = json.Unmarshal(json_facial_watermark_photo_match_response_body, fwpmc.FacialWatermarkPhotoMatchResult)
	if err != nil {
		return err
	}

	return nil
}

func (fwpmc *FacialWatermarkPhotoMatchClient) Flush() {
	fwpmc.FacialWatermarkPhotoMatchResult = nil
}
