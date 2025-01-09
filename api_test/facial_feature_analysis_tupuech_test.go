package api_test

import (
	"io"
	"os"
	"testing"

	"github.com/ZSLTChenXiYin/xfyun-api/face_match"
)

func TestFacialFeatureAnalysisTupuech(t *testing.T) {
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

	ffatc := face_match.NewFacialFeatureAnalysisTupuechClient(
		face_match.WithFacialFeatureAnalysisTupuechClientBasicConfiguration(app_id, api_secret, api_key),
		face_match.WithFacialFeatureAnalysisTupuechClientRequestConfiguration(
			face_match.DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_ADDRESS,
			face_match.DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_REQUEST_LINE,
			face_match.DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_HOST,
		),
	)

	err = ffatc.Ready()
	if err != nil {
		t.Fatal(err)
	}

	err = ffatc.AddInput(image_file_type, image)
	if err != nil {
		t.Fatal(err)
	}

	err = ffatc.Do()
	if err != nil {
		t.Fatal(err)
	}
	defer ffatc.Flush()

	facial_feature_analysis_tupuech_result, err := ffatc.FacialFeatureAnalysisTupuechResponseBody.GetFacialFeatureAnalysisTupuechResult()
	if err != nil {
		t.Fatal(err)
	}
	if facial_feature_analysis_tupuech_result.Ret != 0 {
		t.Fatalf("facial_feature_analysis_tupuech_result.Ret == %d", facial_feature_analysis_tupuech_result.Ret)
	}

	t.Logf("facial_feature_analysis_tupuech_result == %v", facial_feature_analysis_tupuech_result)
}
