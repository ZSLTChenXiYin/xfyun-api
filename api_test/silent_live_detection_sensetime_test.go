package api_test

import (
	"io"
	"os"
	"testing"

	"github.com/ZSLTChenXiYin/xfyun-api/face_match"
)

func TestSilentLiveDetectionSensetime(t *testing.T) {
	app_id := ""
	api_key := ""

	file_path := ""

	file, err := os.Open(file_path)
	if err != nil {
		t.Fatal(err)
	}

	mv, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	sldsc := face_match.NewSilentLiveDetectionSensetimeClient(
		face_match.WithSilentLiveDetectionSensetimeClientBasicConfiguration(app_id, api_key),
		face_match.WithSilentLiveDetectionSensetimeClientRequestConfiguration(face_match.DEFAULT_SILENT_LIVE_DETECTION_SENSETIME_REQUEST_ADDRESS),
	)

	err = sldsc.Ready()
	if err != nil {
		t.Fatal(err)
	}

	err = sldsc.AddFile(mv)
	if err != nil {
		t.Fatal(err)
	}

	err = sldsc.Do(false)
	if err != nil {
		t.Fatal(err)
	}
	defer sldsc.Flush()

	t.Logf("SilentLiveDetectionSensetimeResult.Data.Liveness_score == %v", sldsc.SilentLiveDetectionSensetimeResult.Data.Liveness_score)
}
