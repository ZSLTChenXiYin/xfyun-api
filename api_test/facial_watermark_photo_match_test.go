package api_test

import (
	"io"
	"os"
	"testing"

	"github.com/ZSLTChenXiYin/xfyun-api/face_match"
)

func TestFacialWatermarkPhotoMatch(t *testing.T) {
	app_id := ""
	api_key := ""

	face_file_path := ""
	watermark_file_path := ""

	image_file, err := os.Open(face_file_path)
	if err != nil {
		t.Fatal(err)
	}

	face_image, err := io.ReadAll(image_file)
	if err != nil {
		t.Fatal(err)
	}

	image_file, err = os.Open(watermark_file_path)
	if err != nil {
		t.Fatal(err)
	}

	watermark_image, err := io.ReadAll(image_file)
	if err != nil {
		t.Fatal(err)
	}

	fwpmc := face_match.NewFacialWatermarkPhotoMatchClient(
		face_match.WithFacialWatermarkPhotoMatchClientBasicConfiguration(app_id, api_key),
		face_match.WithFacialWatermarkPhotoMatchClientRequestConfiguration(
			face_match.DEFAULT_FACIAL_WATERMARK_PHOTO_MATCH_REQUEST_ADDRESS,
			face_match.DEFAULT_FACIAL_WATERMARK_PHOTO_MATCH_HOST,
		),
	)

	err = fwpmc.AddFaceImage(face_image)
	if err != nil {
		t.Fatal(err)
	}

	err = fwpmc.AddWatermarkImage(watermark_image)
	if err != nil {
		t.Fatal(err)
	}

	err = fwpmc.Do(true)
	if err != nil {
		t.Fatal(err)
	}
	if fwpmc.FacialWatermarkPhotoMatchResult.Code != "0" {
		t.Fatalf("FacialWatermarkPhotoMatchResult == %v", fwpmc.FacialWatermarkPhotoMatchResult)
	}
	defer fwpmc.Flush()

	t.Logf("FacialWatermarkPhotoMatchResult.Data == %f", fwpmc.FacialWatermarkPhotoMatchResult.Data)
}
