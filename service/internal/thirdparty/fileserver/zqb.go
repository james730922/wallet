package fileserver

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/james730922/wallet/service/internal/thirdparty/logger"
	"github.com/james730922/wallet/service/internal/utils/conf"
	"github.com/james730922/wallet/service/internal/utils/errs"
	"github.com/james730922/wallet/service/internal/utils/tools"
)

const (
	defaultImagePath = "/default.png"
)

type FilePath int32

const (
	FilePathBankAccount         FilePath = 1
	FilePathBankCode            FilePath = 2
	FilePathBankDepositCategory FilePath = 3
)

func filePathName(filepath FilePath) string {
	filePathName := map[int32]string{
		1: "bank/account",
		2: "bank/code",
		3: "bank/deposit/category",
	}

	return filePathName[int32(filepath)]
}

func GetFilePath(str string) FilePath {
	filePathValue := map[string]int32{
		"bank-account":          1,
		"bank-code":             2,
		"bank-deposit-category": 3,
	}

	if s, ok := filePathValue[str]; ok {
		return FilePath(s)
	}
	return -1
}

func NewZQBFileServer(fs FileServer) IZQBFileServer {
	return &zqbFileServer{
		fs: fs,
	}
}

type IZQBFileServer interface {
	Upload(path FilePath, fileName string, image []byte, opt ...FuncFolderOpt) (string, error)
	FilePathForDB(path FilePath, fileName string, image []byte, opt ...FuncFolderOpt) (string, error) // 與Upload()得到的路徑一致
	Delete(fileName string) (bool, error)
	URL(fileName string, opt ...FuncURLOpt) string
	InternalURL(fileName string, opt ...FuncURLOpt) string
}

type zqbFileServer struct {
	fs FileServer
}

func (z *zqbFileServer) isConnect() bool {
	err := z.fs.HealthCheck()

	// 如果失敗原因是resp的status不為200 -> 再試一次
	if err == errs.FileServerResponseNotOK {
		z.fs.Login()
		err = z.fs.HealthCheck()
	}

	// 如果仍失敗，則回傳無連線
	if err != nil {
		return false
	}

	return true
}

func (z *zqbFileServer) Upload(path FilePath, fileName string, image []byte, opt ...FuncFolderOpt) (string, error) {
	beginTime := time.Now()

	defer func() {
		execTime := time.Now().Sub(beginTime)
		if execTime > 5*time.Second {
			logger.ApLog().Warnf("Upload fileName: %s, size: %d, execTime: %s", fileName, len(image), execTime.String())
		}
	}()

	// 檢查與fileserver連線是否正常
	if !z.isConnect() {
		return "", errs.FileServerResponseNotOK
	}

	format, err := tools.ImageFormat(image)
	if err != nil {
		return "", err
	}

	o := &FolderOpt{}
	for _, f := range opt {
		f(o)
	}

	uploadPath := o.SetUploadPath(path)
	uploadFileName := fmt.Sprintf("%s.%s", fileName, format)

	// 上傳
	err = z.fs.Upload(uploadPath, uploadFileName, bytes.NewReader(image))

	// 如果上傳失敗原因是resp的status不為200 -> 再試一次
	if err == errs.FileServerResponseNotOK {
		z.fs.Login()
		err = z.fs.Upload(uploadPath, uploadFileName, bytes.NewReader(image))
	}

	// 如果仍失敗，則回傳錯誤
	if err != nil {
		return "", err
	}

	forDBPath := o.GetPath(path, uploadFileName)

	return forDBPath, nil
}

func (z *zqbFileServer) FilePathForDB(path FilePath, fileName string, image []byte, opt ...FuncFolderOpt) (string, error) {
	format, err := tools.ImageFormat(image)
	if err != nil {
		return "", err
	}

	uploadFileName := fmt.Sprintf("%s.%s", fileName, format)

	o := &FolderOpt{}
	for _, f := range opt {
		f(o)
	}

	forDBPath := o.GetPath(path, uploadFileName)

	return forDBPath, nil
}

// fileName應直接使用mysql內的儲存值，也就是省略/zqb/app之前的URL
func (z *zqbFileServer) Delete(fileName string) (bool, error) {
	return z.fs.Delete(conf.FileServer().GetInternalFolder() + fileName)
}

func (z *zqbFileServer) URL(fileName string, opt ...FuncURLOpt) string {
	o := &URLOpt{FileName: fileName}
	for _, f := range opt {
		f(o)
	}

	return conf.FileServer().GetPath() + o.URL()
}

// Resp 回傳圖片時禁止用這方法。
// 用於 server 對 fileserver 取圖片（走內網）。
func (z *zqbFileServer) InternalURL(fileName string, opt ...FuncURLOpt) string {
	o := &URLOpt{FileName: fileName}
	for _, f := range opt {
		f(o)
	}

	return conf.FileServer().GetInternalPath() + o.URL()
}

type URLOpt struct {
	FileName   string
	UpdateTime *time.Time
}

func (u *URLOpt) URL() string {
	fileName := strings.TrimSpace(u.FileName)
	if fileName == "" {
		fileName = defaultImagePath
	}

	params := url.Values{}
	if u.UpdateTime != nil {
		params.Add("u", strconv.Itoa(int(u.UpdateTime.Unix())))
	}

	qs := ""
	if len(params) > 0 {
		qs = "?" + params.Encode()
	}

	return fileName + qs
}

type FuncURLOpt func(o *URLOpt)

func WithURLUpdateTime(t time.Time) FuncURLOpt {
	return func(o *URLOpt) {
		o.UpdateTime = &t
	}
}

type FolderOpt struct {
	FolderName *string
}

func (u *FolderOpt) GetPath(filePath FilePath, uploadFileName string) string {
	var dbPath string
	dbPath = fmt.Sprintf("/%s/%s", filePathName(filePath), uploadFileName)
	if u.FolderName != nil {
		dbPath = fmt.Sprintf("/%s/%s/%s", filePathName(filePath), aws.StringValue(u.FolderName), uploadFileName)
	}
	return dbPath
}

func (u *FolderOpt) SetUploadPath(path FilePath) string {
	var uploadPath string

	uploadPath = fmt.Sprintf(conf.FileServer().GetInternalFolder()+"/%s", filePathName(path))
	if u.FolderName != nil {
		uploadPath = fmt.Sprintf(conf.FileServer().GetInternalFolder()+"/%s/%s", filePathName(path), aws.StringValue(u.FolderName))
	}
	return uploadPath
}

type FuncFolderOpt func(o *FolderOpt)

func SetFolderName(folderName string) FuncFolderOpt {
	return func(o *FolderOpt) {
		o.FolderName = &folderName
	}
}
