package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(upload))
}

func upload(w http.ResponseWriter, r *http.Request) {
	img, ext, err := image.Decode(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fp, err := os.Create("upload." + ext)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer fp.Close()

	switch ext {
	case "png":
		err = png.Encode(fp, img)
	case "jpeg":
		err = jpeg.Encode(fp, img, &jpeg.Options{Quality: 90})
	default:
		http.Error(w, "format not support", 400)
		return
	}
}
