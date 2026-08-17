package condition

type FileImageUploadReq struct {
	// 上傳路徑
	Path string `json:"path"`
	// 圖片名稱
	FileName string `json:"file_name"`
	// 圖片檔案
	Image []byte `json:"image,omitempty"`
}
