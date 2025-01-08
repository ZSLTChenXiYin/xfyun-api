package api_test

import (
	"io"
	"os"
	"testing"

	"github.com/ZSLTChenXiYin/xfyun-api/face_match"
)

func TestFaceMatch(t *testing.T) {
	app_id := ""
	api_secret := ""
	api_key := ""

	image1_file_path := ""
	image1_file_type := ""
	image2_file_path := ""
	image2_file_type := ""

	image_file, err := os.Open(image1_file_path)
	if err != nil {
		t.Fatal(err)
	}

	image1, err := io.ReadAll(image_file)
	if err != nil {
		t.Fatal(err)
	}

	image_file, err = os.Open(image2_file_path)
	if err != nil {
		t.Fatal(err)
	}

	image2, err := io.ReadAll(image_file)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Standardize the use of face_match_client", func(t *testing.T) {
		fmc_request_configuration := face_match.WithFaceMatchClientRequestConfiguration(face_match.DEFAULT_FACE_MATCH_REQUEST_ADDRESS, face_match.DEFAULT_FACE_MATCH_REQUEST_LINE, face_match.DEFAULT_FACE_MATCH_HOST)
		fmc_basic_configuration := face_match.WithFaceMatchClientBasicConfiguration(app_id, api_secret, api_key)

		fmc := face_match.NewFaceMatchClient(fmc_request_configuration, fmc_basic_configuration)

		err := fmc.Ready()
		if err != nil {
			t.Fatal(err)
		}

		t.Run("Standardize the use of AddInput1 add AddInput2", func(t *testing.T) {
			err = fmc.AddInput(1, image1_file_type, image1)
			if err != nil {
				t.Fatal(err)
			}

			err = fmc.AddInput(2, image2_file_type, image2)
			if err != nil {
				t.Fatal(err)
			}

			err = fmc.Do()
			if err != nil {
				t.Fatal(err)
			}

			face_match_result, err := fmc.FaceMatchResponseBody.GetFaceCompareResult()
			if err != nil {
				t.Fatal(err)
			}
			if face_match_result.Ret != 0 {
				t.Fatalf("face_match_result.Ret == %d", face_match_result.Ret)
			}
			defer fmc.Flush()

			t.Logf("face_match_result.Score == %f", face_match_result.Score)
		})

		t.Run("Standardize the use of AddAllInput", func(t *testing.T) {
			err = fmc.AddAllInput(image1_file_type, image1, image2_file_type, image2)
			if err != nil {
				t.Fatal(err)
			}

			err = fmc.Do()
			if err != nil {
				t.Fatal(err)
			}

			face_match_result, err := fmc.FaceMatchResponseBody.GetFaceCompareResult()
			if err != nil {
				t.Fatal(err)
			}
			if face_match_result.Ret != 0 {
				t.Fatalf("face_match_result.Ret == %d", face_match_result.Ret)
			}
			defer fmc.Flush()

			t.Logf("face_match_result.Score == %f", face_match_result.Score)
		})
	})
}
