package httpapi

import (
	"image"
	"image/jpeg"
	"image/png"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestScreenshotResponseMatchesEncodedFormat(t *testing.T) {
	for _, format := range []string{"png", "jpeg"} {
		t.Run(format, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "0.png")
			file, err := os.Create(path)
			if err != nil {
				t.Fatal(err)
			}
			img := image.NewRGBA(image.Rect(0, 0, 2, 2))
			if format == "png" {
				err = png.Encode(file, img)
			} else {
				err = jpeg.Encode(file, img, nil)
			}
			file.Close()
			if err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			serveScreenshot(response, httptest.NewRequest("GET", "/evidence/run/0.png", nil), path)
			if response.Code != 200 || response.Header().Get("Content-Type") != "image/"+format {
				t.Fatalf("incorrect screenshot response: %d %s", response.Code, response.Header().Get("Content-Type"))
			}
			_, decoded, err := image.DecodeConfig(response.Body)
			if err != nil || decoded != format {
				t.Fatalf("image bytes differ from response header: %s %v", decoded, err)
			}
		})
	}
}
