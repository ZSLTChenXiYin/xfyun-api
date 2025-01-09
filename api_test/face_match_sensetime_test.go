package api_test

import (
	"io"
	"os"
	"testing"

	"github.com/ZSLTChenXiYin/xfyun-api/face_match"
)

func TestFaceMatchSensetimeClient(t *testing.T) {
	app_id := ""
	api_key := ""

	image1_file_path := ""
	image2_file_path := ""

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

	fmsc := face_match.NewFaceMatchSensetimeClient(
		face_match.WithFaceMatchSensetimeClientBasicConfiguration(app_id, api_key),
		face_match.WithFaceMatchSensetimeClientRequestConfiguration(
			face_match.DEFAULT_FACE_MATCH_SENSETIME_REQUEST_ADDRESS,
		),
	)

	err = fmsc.Ready()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Standardize the use of AddInput1 add AddInput2", func(t *testing.T) {
		err = fmsc.AddInput(1, image1)
		if err != nil {
			t.Fatal(err)
		}

		err = fmsc.AddInput(2, image2)
		if err != nil {
			t.Fatal(err)
		}

		err = fmsc.Do(true)
		if err != nil {
			t.Fatal(err)
		}
		defer fmsc.Flush()

		t.Logf("FaceMatchSensetimeResult.Data == %f", fmsc.FaceMatchSensetimeResult.Data)
	})

	t.Run("Standardize the use of AddAllInput", func(t *testing.T) {
		err = fmsc.AddAllInput(image1, image2)
		if err != nil {
			t.Fatal(err)
		}

		err = fmsc.Do(true)
		if err != nil {
			t.Fatal(err)
		}
		defer fmsc.Flush()

		t.Logf("FaceMatchSensetimeResult.Data == %f", fmsc.FaceMatchSensetimeResult.Data)
	})
}
