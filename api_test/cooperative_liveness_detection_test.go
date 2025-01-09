package api_test

import (
	"io"
	"os"
	"testing"

	"github.com/ZSLTChenXiYin/xfyun-api/face_match"
)

func TestCooperativeLivenessDetection(t *testing.T) {
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

	cldc := face_match.NewCooperativeLivenessDetectionClient(
		face_match.WithCooperativeLivenessDetectionClientBasicConfiguration(app_id, api_secret, api_key),
		face_match.WithCooperativeLivenessDetectionClientRequestConfiguration(
			face_match.DEFAULT_COOPERATIVE_LIVENESS_DETECTION_REQUEST_ADDRESS,
			face_match.DEFAULT_COOPERATIVE_LIVENESS_DETECTION_REQUEST_LINE,
			face_match.DEFAULT_COOPERATIVE_LIVENESS_DETECTION_HOST,
		),
	)

	err = cldc.Ready()
	if err != nil {
		t.Fatal(err)
	}

	err = cldc.AddInput(image_file_type, image)
	if err != nil {
		t.Fatal(err)
	}

	err = cldc.Do()
	if err != nil {
		t.Fatal(err)
	}
	defer cldc.Flush()

	cooperative_liveness_detection_result, err := cldc.CooperativeLivenessDetectionResponseBody.GetCooperativeLivenessDetectionResult()
	if err != nil {
		t.Fatal(err)
	}
	if cooperative_liveness_detection_result.Ret != 0 {
		t.Fatalf("cooperative_liveness_detection_result.Ret == %d", cooperative_liveness_detection_result.Ret)
	}

	t.Logf("cooperative_liveness_detection_result.Face_1 == %v", cooperative_liveness_detection_result.Face_1)
}
