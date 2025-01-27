package api_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/ZSLTChenXiYin/xfyun-api/face_match"
)

func TestFacialFeatureAnalysisTupuech(t *testing.T) {
	app_id := ""
	api_key := ""

	image_file_path := ""

	image_file, err := os.Open(image_file_path)
	if err != nil {
		t.Fatal(err)
	}

	image, err := io.ReadAll(image_file)
	if err != nil {
		t.Fatal(err)
	}

	ffatc := face_match.NewFacialFeatureAnalysisTupuechClient(
		face_match.WithFacialFeatureAnalysisTupuechClientBasicConfiguration(app_id, api_key),
		face_match.WithFacialFeatureAnalysisTupuechClientRequestConfiguration(
			face_match.DEFAULT_FACIAL_FEATURE_ANALYSIS_TUPUECH_FACE_SCORE_REQUEST_ADDRESS,
		),
	)

	if ffatc == nil {
		t.Fatal("ffatc is nil")
	}

	err = ffatc.Ready()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Standardize the use of AddFile", func(t *testing.T) {
		err = ffatc.AddFile(image)
		if err != nil {
			t.Fatal(err)
		}

		err = ffatc.Do(image_file.Name(), "")
		if err != nil {
			t.Fatal(err)
		}
		defer ffatc.Flush()

		ffatr_file_list := ffatc.FacialFeatureAnalysisTupuechResult.Data["fileList"].([]any)
		for _, ffatr_file := range ffatr_file_list {
			label := ffatr_file.(map[string]any)["label"].(float64)
			name := ffatr_file.(map[string]any)["name"].(string)
			var slabel string
			switch label {
			case 0:
				slabel = "beautiful"
			case 1:
				slabel = "nice"
			case 2:
				slabel = "ordinary"
			case 3:
				slabel = "ugly"
			case 4:
				slabel = "other"
			case 5:
				slabel = "half face"
			case 6:
				slabel = "multiple people"
			default:
				slabel = fmt.Sprintf("error label %v", label)
			}
			t.Logf("name == %s, label == %s", name, slabel)
		}
	})
}
