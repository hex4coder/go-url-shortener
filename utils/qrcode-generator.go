package utils

import (
	"encoding/base64"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

func EncodeUrlToQrImage(url string) (error, []byte) {
	bytes, err := qrcode.Encode(url, qrcode.Medium, 256)
	return err, bytes
}

func EncodeURLToImageBase64(url string) string {
	var res string = ""

	_, png := EncodeUrlToQrImage(url)

	encodedImage := base64.StdEncoding.EncodeToString(png)

	res = fmt.Sprintf("%s%s", res, encodedImage)

	return res
}
