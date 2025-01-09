package api_test

import (
	"io"
	"os"
	"testing"

	"github.com/ZSLTChenXiYin/xfyun-api/face_match"
)

func TestSilentLiveDetection(t *testing.T) {
	app_id := ""
	api_secret := ""
	api_key := ""

	image_file_path := ""
	image_file_type := ""

	image_file, err := os.Open(image_file_path)
	if err != nil {
		t.Fatal(err)
	}

	image, err := io.ReadAll(image_file)
	if err != nil {
		t.Fatal(err)
	}

	sldc_request_configuration := face_match.WithSilentLiveDetectionClientRequestConfiguration(
		face_match.DEFAULT_SILENT_LIVE_DETECTION_ADDRESS,
		face_match.DEFAULT_SILENT_LIVE_DETECTION_REQUEST_LINE,
		face_match.DEFAULT_SILENT_LIVE_DETECTION_HOST,
	)

	sldc_basic_configuration := face_match.WithSilentLiveDetectionClientBasicConfiguration(
		app_id,
		api_secret,
		api_key,
	)

	sldc := face_match.NewSilentLiveDetectionClient(sldc_request_configuration, sldc_basic_configuration)

	err = sldc.Ready()
	if err != nil {
		t.Fatal(err)
	}

	err = sldc.AddInput(image_file_type, image)
	if err != nil {
		t.Fatal(err)
	}

	err = sldc.Do()
	if err != nil {
		t.Fatal(err)
	}

	silent_live_detection_result, err := sldc.SilentLiveDetectionResponseBody.GetAntiSpoofResult()
	if err != nil {
		t.Fatal(err)
	}
	if silent_live_detection_result.Ret != 0 {
		t.Fatalf("silent_live_detection_result.Ret == %d", silent_live_detection_result.Ret)
	}

	t.Logf("silent_live_detection_result.Score == %f", silent_live_detection_result.Score)
}
