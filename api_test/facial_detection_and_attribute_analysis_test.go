package api_test

import (
	"io"
	"os"
	"testing"

	"github.com/ZSLTChenXiYin/xfyun-api/face_match"
)

func TestFacialDetectionAndAttributeAnalysis(t *testing.T) {
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

	fdaaac := face_match.NewFacialDetectionAndAttributeAnalysisClient(
		face_match.WithFacialDetectionAndAttributeAnalysisClientBasicConfiguration(app_id, api_secret, api_key),
		face_match.WithFacialDetectionAndAttributeAnalysisClientRequestConfiguration(
			face_match.DEFAULT_FACIAL_DETECTION_AND_ATTRIBUTE_ANALYSIS_REQUEST_ADDRESS,
			face_match.DEFAULT_FACIAL_DETECTION_AND_ATTRIBUTE_ANALYSIS_REQUEST_LINE,
			face_match.DEFAULT_FACIAL_DETECTION_AND_ATTRIBUTE_ANALYSIS_HOST,
		),
	)

	err = fdaaac.Ready()
	if err != nil {
		t.Fatal(err)
	}

	err = fdaaac.AddInput(image_file_type, image)
	if err != nil {
		t.Fatal(err)
	}

	err = fdaaac.Do()
	if err != nil {
		t.Fatal(err)
	}
	defer fdaaac.Flush()

	facial_detection_and_attribute_analysis_result, err := fdaaac.FacialDetectionAndAttributeAnalysisResponseBody.GetFacialDetectionAndAttributeAnalysisResult()
	if err != nil {
		t.Fatal(err)
	}
	if facial_detection_and_attribute_analysis_result.Ret != 0 {
		t.Fatalf("facial_detection_and_attribute_analysis_result.Ret == %d", facial_detection_and_attribute_analysis_result.Ret)
	}

	t.Logf("facial_detection_and_attribute_analysis_result == %v", facial_detection_and_attribute_analysis_result)
}
