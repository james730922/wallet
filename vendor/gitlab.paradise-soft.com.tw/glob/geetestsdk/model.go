package geetestsdk

type CaptchaConfig struct {
	GeeTestID                 string // geetest 後台設定
	GeeTestKey                string // geetest 後台設定
	GeeTestUserIDSalt         string // 密碼鹽，專案自行定義
	GeeTestByPassCycleTimeSec int    // 檢查 geetest 服務的 health check 間隔時間
}
