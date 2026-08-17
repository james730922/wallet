package tools

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/nfnt/resize"

	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/errs"
)

var (
	wxRegexp     = regexp.MustCompile("^wxp://(?:[A-Za-z0-9-_]{4})*(?:[A-Za-z0-9-_]{2}|[A-Za-z0-9-_]{3})?$")
	zfbRegexp    = regexp.MustCompile("^[hH][tT][tT][pP][sS]://[qQ][rR].[aA][lL][iI][pP][aA][yY].[cC][oO][mM]/[0-9_a-z_A-Z]+$")
	ysfRegexp    = regexp.MustCompile("^https://qr.95516.com/[0-9]+/[0-9]+$")
	qqRegexp     = regexp.MustCompile("^https://i.qianbao.qq.com/wallet/sqrcode.htm\\?m=tenpay&[0-9_a-z_A-Z]+=")
	paybeiRegexp = regexp.MustCompile("^https://sh.51youdian.com")
)

const (
	qrcodeTypeWX    = "WX_QR"
	qrcodeTypeWXZZ  = "WX_QR_ZZ"
	qrcodeTypeZFB   = "ZFB_QR"
	qrcodeTypeZFBZZ = "ZFB_QR_ZZ"
	qrcodeTypeYSF   = "YSF_QR"
	qrcodeTypeYSFZZ = "YSF_QR_ZZ"
	qrcodeTypeQQ    = "QQ_QR"
	qrcodeTypeQQZZ  = "QQ_QR_ZZ"
)

// ImageFormat return the format of imageData
func ImageFormat(imageData []byte) (string, error) {
	_, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		logger.ApLog().Errorf("err=%v, format=%s", err, format)
		return "", err
	}
	return format, nil
}

func IsQrCode(imageData []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		logger.ApLog().Error(err)
		return "", err
	}

	// prepare BinaryBitmap
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		logger.ApLog().Error(err)
		return "", err
	}

	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		logger.ApLog().Error(err)
		return "", err
	}

	content := result.String()

	return content, nil
}

func CheckQrCodeContent(content string, types string) error {
	switch types {
	case qrcodeTypeWX, qrcodeTypeWXZZ: // 微信

		if !wxRegexp.MatchString(content) && !paybeiRegexp.MatchString(content) {
			return errs.CommonQRNotMatchWX
		}
	case qrcodeTypeZFB, qrcodeTypeZFBZZ: // 支付宝
		if !zfbRegexp.MatchString(content) && !paybeiRegexp.MatchString(content) {
			return errs.CommonQRNotMatchZFB
		}
	case qrcodeTypeYSF, qrcodeTypeYSFZZ: // 云闪付
		if !ysfRegexp.MatchString(content) && !paybeiRegexp.MatchString(content) {
			return errs.CommonQRNotMatchYSF
		}
	case qrcodeTypeQQ, qrcodeTypeQQZZ: // QQ付
		if !qqRegexp.MatchString(content) {
			return errs.CommonQRNotMatchQQ
		}
	default:
		return errs.CommonQRNotMatch
	}
	return nil
}

func ImageCheck(imageData []byte) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		logger.ApLog().Error(err)
		return nil, err
	}

	bounds := img.Bounds()
	width := uint(bounds.Max.X)
	height := uint(bounds.Max.Y)
	img = resize.Resize(width, height, img, resize.Lanczos3)

	out := &bytes.Buffer{}
	switch format {
	case "jpeg":
		err = jpeg.Encode(out, img, nil)
		if err != nil {
			return out.Bytes(), err
		}
	case "png":
		err = png.Encode(out, img)
		if err != nil {
			return out.Bytes(), err
		}
	default:
		logger.ApLog().Error("image format is not jpeg or png")
		return nil, errs.CommonImageFormatError
	}

	return out.Bytes(), nil
}

func ImageHash(imageData []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(imageData))
}

func IsImageAltered(imageURL, hash string) (bool, error) {
	// 如果沒有哈希紀錄，就不驗證圖片是否被竄改
	if hash == "" {
		return false, nil
	}

	// 根據圖片路徑取得圖片
	imageData, err := GetImageDataFrom(imageURL)
	if err != nil {
		return false, err
	}

	// 比對圖片哈希值與當初儲存起來的哈希值
	hashFromDownload := ImageHash(imageData)
	if hashFromDownload != hash {
		logger.SysLog().Errorf("image_path: %s, hashFromDownload:%v, originalHash:%v", imageURL, hashFromDownload, hash)
		return true, nil
	}
	return false, nil
}

func GetImageDataFrom(imageURL string) ([]byte, error) {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Get(imageURL)
	if err != nil {
		return nil, errs.FileServerResponseNotOK
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.FileServerResponseNotOK
	}

	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.CommonParseError
	}

	return imageData, nil
}
