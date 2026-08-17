package errs

import "net/http"

// ErrorCode: httpStatus-code
// ex: 400-1
type ErrorCode string

// httpStatus reference to HTTP Status Code
type httpStatus int

const (
	statusBadRequest         httpStatus = http.StatusBadRequest
	statusUnauthorized       httpStatus = http.StatusUnauthorized
	statusForbidden          httpStatus = http.StatusForbidden
	statusNotFound           httpStatus = http.StatusNotFound
	statusServiceUnavailable httpStatus = http.StatusServiceUnavailable
)

// prefix code
type PrefixCode int

const (
	CodeCommon      PrefixCode = 0
	CodeDB          PrefixCode = 1
	CodeAuth        PrefixCode = 2
	CodeMember      PrefixCode = 3
	CodeWallet      PrefixCode = 5
	CodeOrder       PrefixCode = 6
	CodeDeposit     PrefixCode = 7
	CodeBalance     PrefixCode = 15
	CodePayCode     PrefixCode = 22
	CodeScanPayCode PrefixCode = 23
)

// common error code
var (
	DBNoRow           = new(statusBadRequest, CodeDB, 1, "资料不存在")
	DBUpdateFailed    = new(statusBadRequest, CodeDB, 2, "资料修改失败")
	DBUpdateDuplicate = new(statusBadRequest, CodeDB, 3, "资料修改后未异动")
	DBInsertFailed    = new(statusBadRequest, CodeDB, 4, "资料新增失败")
	DBInsertDuplicate = new(statusBadRequest, CodeDB, 5, "资料重复新增")
	DBFetchFailed     = new(statusBadRequest, CodeDB, 6, "资料取得失败")
	DBDeleteFailed    = new(statusBadRequest, CodeDB, 7, "资料删除失败")
	DBAccessDenied    = new(statusBadRequest, CodeDB, 8, "资料访问被拒绝")
	DBOperationFailed = new(statusBadRequest, CodeDB, 9, "资料操作失败")
	DBUpdateZeroRow   = new(statusBadRequest, CodeDB, 10, "资料更新比数为0")
)

var (
	CommonUnknownError                  = new(statusBadRequest, CodeCommon, 101, "未知错误")
	CommonNoData                        = new(statusBadRequest, CodeCommon, 102, "查无资料")
	CommonRawSQLNotFound                = new(statusBadRequest, CodeCommon, 103, "找不到执行档")
	CommonNoMemberID                    = new(statusBadRequest, CodeCommon, 104, "无法取得会员ID")
	CommonRequestParamInvalid           = new(statusBadRequest, CodeCommon, 105, "请求参数错误")
	CommonRequestParamParseFailed       = new(statusBadRequest, CodeCommon, 106, "请求参数解析失败")
	CommonRequestPageError              = new(statusBadRequest, CodeCommon, 107, "请求的页数错误")
	CommonParseError                    = new(statusBadRequest, CodeCommon, 108, "解析失败")
	CommonRecordOwnershipWrong          = new(statusBadRequest, CodeCommon, 110, "非该会员资料")
	CommonInterfaceImplementError       = new(statusBadRequest, CodeCommon, 111, "介面实作错误")
	CommonQRFormatError                 = new(statusBadRequest, CodeCommon, 112, "请使用保存至相册的收款码进行上传。") // 非可執行二維碼，請重新上傳。
	CommonParseTimeZoneError            = new(statusBadRequest, CodeCommon, 113, "時區解析錯誤")
	CommonFrequentOperationError        = new(statusBadRequest, CodeCommon, 114, "频繁操作，请稍后再尝试")
	CommonConfigureInvalid              = new(statusBadRequest, CodeCommon, 115, "设置参数错误")
	CommonImageFormatError              = new(statusBadRequest, CodeCommon, 116, "圖片格式必須為png或jpeg")
	CommonQRNotMatch                    = new(statusBadRequest, CodeCommon, 117, "我们目前不支持您的二维码，请重新生成")
	CommonQRNotMatchWX                  = new(statusBadRequest, CodeCommon, 118, "您的二维码不属于微信，请重新生成")
	CommonQRNotMatchZFB                 = new(statusBadRequest, CodeCommon, 119, "请使用保存至相册的收款码进行上传")
	CommonQRNotMatchYSF                 = new(statusBadRequest, CodeCommon, 120, "您的二维码不属于云闪付，请重新生成")
	CommonQRNotMatchQQ                  = new(statusBadRequest, CodeCommon, 121, "您的二维码不属于QQ，请重新生成")
	CommonCardNumberInvalid             = new(statusBadRequest, CodeCommon, 122, "卡号格式不正确，请重新输入")
	CommonAmountIntervalOverlapping     = new(statusBadRequest, CodeCommon, 126, "金额区间重叠，请修正后重新上传")
	CommonAmountDuplicate               = new(statusBadRequest, CodeCommon, 127, "金额重复，请修正后重新上传")
	CommonNotFindMemberID               = new(statusBadRequest, CodeCommon, 128, "会员编号不存在，请重新输入")
	CommonAmountEmpty                   = new(statusBadRequest, CodeCommon, 129, "金额设定不可为空")
	CommonServiceUnavailable            = new(statusServiceUnavailable, CodeCommon, 130, "系统维护中")
	CommonTimeIntervalForExportTooLarge = new(statusBadRequest, CodeCommon, 131, "汇出天数超过%d天，请重新选取")
	CommonExportExcelWithNullData       = new(statusBadRequest, CodeCommon, 132, "无资料汇出")
	CommonTimeIntervalInvalid           = new(statusBadRequest, CodeCommon, 133, "结束时间大于开始时间，时间区间不合理")
	CommonTimeIntervalForSearchTooLarge = new(statusBadRequest, CodeCommon, 134, "搜寻天数超过%d天，请重新选取")
	CommonTimeIntervalIsNecessary       = new(statusBadRequest, CodeCommon, 135, "时间区间为必填条件")
	CommonImageHasBeenAltered           = new(statusBadRequest, CodeCommon, 136, "图片已被窜改")
	CommonCardNumberLenError            = new(statusBadRequest, CodeCommon, 137, "卡号必须是介于1〜20位之间的数字")
	CommonNotFindThisMemberID           = new(statusBadRequest, CodeCommon, 138, "会员编号%s不存在，请重新输入")
	CommonAmountDecimalPlacesError      = new(statusBadRequest, CodeCommon, 139, "金额小数位不正确")
	CommonDoubleClickError              = new(statusBadRequest, CodeCommon, 140, "操作过于频繁，请稍后再试")
	CommonMemberLevelIsBlackList        = new(statusBadRequest, CodeCommon, 142, "功能受限，请联系客服")
	CommonImageWrongPath                = new(statusBadRequest, CodeCommon, 144, "图片路径错误")
	CommonDataDelay                     = new(statusBadRequest, CodeCommon, 146, "资料已更新，请重新查找")
)

var (
	FrameworkContextErrorType      = new(statusBadRequest, CodeCommon, 201, "Context 型别错误")
	FrameworkNotProtoBuf           = new(statusBadRequest, CodeCommon, 202, "型别不是protoBuf")
	FrameworkRequestBodyParseError = new(statusBadRequest, CodeCommon, 203, "请求解析失败")
	FrameworkIllegalParameter      = new(statusBadRequest, CodeCommon, 204, "请求参数错误")
)

var (
	FileServerUploadFailed  = new(statusBadRequest, CodeCommon, 301, "图片上传失败")
	FileServerResponseNotOK = new(statusBadRequest, CodeCommon, 302, "图片库异常")
)

var (
	BufSpaceFilled = new(statusBadRequest, CodeCommon, 501, "缓冲空间已满")
)

var (
	RegistrationMobileFmtErr            = new(statusBadRequest, CodeCommon, 701, "请输入正确的手机号码")
	RegistrationQQFmtErr                = new(statusBadRequest, CodeCommon, 702, "请输入正确的QQ号")
	RegistrationPasswdFmtErr            = new(statusBadRequest, CodeCommon, 703, "密码格式不符合要求，请至少包含英文、数字以及长度介于8-16码")
	PasswdNotEqual                      = new(statusBadRequest, CodeCommon, 704, "两次密码输入不相等")
	QQEmpty                             = new(statusBadRequest, CodeCommon, 705, "请填写QQ号")
	RegistrationConfirmPasswdEmpty      = new(statusBadRequest, CodeCommon, 706, "请填写确认密码")
	MobileEmpty                         = new(statusBadRequest, CodeCommon, 707, "请填写手机号")
	PasswdEmpty                         = new(statusBadRequest, CodeCommon, 708, "请填写密码")
	NewPasswdEmpty                      = new(statusBadRequest, CodeCommon, 709, "请填写新密码")
	NewConfirmPasswdEmpty               = new(statusBadRequest, CodeCommon, 710, "请填写确认新密码")
	NameEmpty                           = new(statusBadRequest, CodeCommon, 711, "请填写中文姓名")
	RegistrationNameFmtErr              = new(statusBadRequest, CodeCommon, 712, "只接受中文姓名，并且长度在20字内")
	DataNotQualified                    = new(statusBadRequest, CodeCommon, 713, "资料有误")
	MobileDuplicate                     = new(statusBadRequest, CodeCommon, 716, "手机号码重复")
	QQDuplicate                         = new(statusBadRequest, CodeCommon, 717, "QQ号码重复")
	RegistrationEmailFmtErr             = new(statusBadRequest, CodeCommon, 718, "请输入正确的Email")
	RegistrationCheckOldPasswdError     = new(statusBadRequest, CodeCommon, 721, "原密码错误")
	RegistrationScanPayPasswdFmtErr     = new(statusBadRequest, CodeCommon, 722, "支付密码格式不符合要求，仅数字以及长度介于6码")
	RegistrationCheckScanPayPasswdError = new(statusBadRequest, CodeCommon, 723, "支付密码错误")
)

// member error code, begin of 1000.
var (
	AuthTokenCreateFailed                     = new(statusBadRequest, CodeAuth, 1, "权证取得失败")
	AuthTokenVerifyFailed                     = new(statusBadRequest, CodeAuth, 2, "权证验证失败")
	AuthTokenUnauthorized                     = new(statusUnauthorized, CodeAuth, 3, "请重新登录")
	AuthAccountOrPasswdError                  = new(statusUnauthorized, CodeAuth, 4, "帐号/密码错误")
	AuthOperationForbidden                    = new(statusForbidden, CodeAuth, 5, "没有操作权限")
	AuthLogoutFailed                          = new(statusForbidden, CodeAuth, 6, "登出失败")
	AuthAccountDisable                        = new(statusUnauthorized, CodeAuth, 7, "帐号未启用")
	AuthLoginModeError                        = new(statusForbidden, CodeAuth, 8, "请求与登入模式冲突")
	AuthTokenGetFailed                        = new(statusForbidden, CodeAuth, 9, "无法取得token")
	AuthTokenDelFailed                        = new(statusBadRequest, CodeAuth, 10, "权证删除失败")
	AuthMemberSecurityPasswdTokenUnauthorized = new(statusForbidden, CodeAuth, 11, "安全密码错误")
)

// Auth error code, begin of 2000.
var (
	LoginMemberSendVerificationCodeFetchTooFast  = new(statusBadRequest, CodeAuth, 101, "不可连续取得验证码")
	LoginMemberIdentifierFailed                  = new(statusBadRequest, CodeAuth, 102, "验证失败，请确认验证码是否正确")
	LoginMemberIdentifierVerificationCodeInvalid = new(statusBadRequest, CodeAuth, 103, "验证码过期")
	LoginMemberSendVerificationTooManyTime       = new(statusBadRequest, CodeAuth, 104, "取得验证码次数过多，请稍后再试")
	LoginMemberReverseVerificationFailed         = new(statusBadRequest, CodeAuth, 105, "反转验证失败")
	LoginMemberReverseVerificationCodeError      = new(statusBadRequest, CodeAuth, 106, "验证码不正确，或者验证码与手机号不相符")
	LoginMemberNoReceiveVerificationCode         = new(statusBadRequest, CodeAuth, 107, "未收到验证码")
	LoginMemberRegistFailed                      = new(statusBadRequest, CodeAuth, 109, "注册失败")
	LoginMemberRegistDuplicated                  = new(statusBadRequest, CodeAuth, 110, "账号已存在")
	LoginMemberLoginFailed                       = new(statusBadRequest, CodeAuth, 111, "帐号/密码有误，请重新输入")
	LoginMemberCaptchaValidateFailed             = new(statusBadRequest, CodeAuth, 114, "验证失败")
)

// member error code
var (
	MemberCreateFailed        = new(statusBadRequest, CodeMember, 1, "会员建立失败")
	MemberNotFound            = new(statusBadRequest, CodeMember, 2, "查无会员")
	MemberDisabled            = new(statusBadRequest, CodeMember, 7, "账号被冻结，请咨询客服")
	MemberSimplifyIDExhausted = new(statusServiceUnavailable, CodeMember, 8, "簡化會員編號已用完")
)

// wallet member
var (
	WalletMemberCreateFailed          = new(statusBadRequest, CodeWallet, 1, "会员钱包建立失败")
	WalletMemberNoFound               = new(statusBadRequest, CodeWallet, 2, "查无会员钱包")
	WalletMemberSignValidateFailed    = new(statusBadRequest, CodeWallet, 3, "会员钱包签名验证失败")
	WalletMemberUpdateFailed          = new(statusBadRequest, CodeWallet, 4, "会员钱包更新失败")
	WalletMemberUpdateBonusIsNegative = new(statusBadRequest, CodeWallet, 5, "会员钱包红利金额为负")
	WalletMemberJournalCreateFailed   = new(statusBadRequest, CodeWallet, 6, "会员钱包日志建立失败")
)

// order deposit
var (
	OrderDepositCreateFailed = new(statusBadRequest, CodeOrder, 1, "入款申请失败")
)

// order transaction
var (
	OrderTransactionCreateFailed     = new(statusBadRequest, CodeOrder, 201, "交易纪录失败")
	TransactionQueryMultiTypeIDError = new(statusBadRequest, CodeOrder, 202, "无法同时查询多种订单编号，请重新操作")
)

// order bind
var (
	OrderBindCreateFailed = new(statusBadRequest, CodeOrder, 301, "交易纪录绑定失败")
	OrderBindQueryFailed  = new(statusBadRequest, CodeOrder, 302, "交易绑定查询失败")
)

// deposit
var (
	DepositErrorParamPayName = new(statusBadRequest, CodeDeposit, 1, "入款姓名不可为空")
	DepositErrorParamAmount  = new(statusBadRequest, CodeDeposit, 2, "入款金额错误")
	DepositErrorParamType    = new(statusBadRequest, CodeDeposit, 3, "入款方式错误")
	DepositNoMethod          = new(statusBadRequest, CodeDeposit, 4, "尚无设置入款方式，请通知客服处理")
	DepositTooOften          = new(statusBadRequest, CodeDeposit, 5, "禁止在短时间内，存入多笔同金额款项")
	DepositBankNotExist      = new(statusBadRequest, CodeDeposit, 6, "银行代码不存在")
	DepositAmountShouldBeInt = new(statusBadRequest, CodeDeposit, 12, "入款金额必须为正整数")
	DepositAmountExceedRange = new(statusBadRequest, CodeDeposit, 13, "入款金额超出限制范围")
)

// balance 为了前端可以依据error code导转到入款页面，因此与wallet 独立
var (
	WalletMemberUpdateBalanceIsNegative      = new(statusBadRequest, CodeBalance, 1, "账户可用余额不足")
	WalletMemberUpdateFrozenAmountIsNegative = new(statusBadRequest, CodeBalance, 2, "会员冻结余额不足")
	WalletMemberNotEnoughForMargin           = new(statusBadRequest, CodeBalance, 3, "余额不足扣保证金")
	WalletMemberAmountUnreasonable           = new(statusBadRequest, CodeBalance, 4, "账户内部额度不合理")
)

// PayCode
var (
	PayCodeScanPayEncodeError = new(statusBadRequest, CodePayCode, 1, "二维码产生发生失败")
	PayCodeScanPayDecodeError = new(statusBadRequest, CodePayCode, 2, "错误二维码")
	PayCodeScanPayExpired     = new(statusBadRequest, CodePayCode, 3, "二维码已过期")
)

// ScanPay
var (
	ScanPayOrderCancel               = new(statusBadRequest, CodeScanPayCode, 1, "订单已取消，请勿重复扫码，如果再次支付请重新获取")
	ScanPayOrderAlreadyDone          = new(statusBadRequest, CodeScanPayCode, 2, "订单已完成，请勿重复扫码，如果再次支付请重新获取")
	ScanPayAddRecordOrderIDDuplicate = new(statusBadRequest, CodeScanPayCode, 3, "新增单号失败，单号重复")
	ScanPayAddRecordFailed           = new(statusBadRequest, CodeScanPayCode, 4, "请勿重复扫码，如果再次支付请重新获取") // 新增單號失敗
	OrderScanPaySignValidateFailed   = new(statusBadRequest, CodeScanPayCode, 5, "订单校验失败")             // 签名验证失败
	ScanPayAmountValidateFailed      = new(statusBadRequest, CodeScanPayCode, 6, "金额错误")
	ScanPayRecordStatusError         = new(statusBadRequest, CodeScanPayCode, 7, "状态错误")
	ScanPayRecordIDNotFound          = new(statusBadRequest, CodeScanPayCode, 8, "查无扫码订单")
	ScanPayGetRecordFailed           = new(statusBadRequest, CodeScanPayCode, 9, "取得扫码支付订单失败")
	ScanQRCodeVerifyFailed           = new(statusBadRequest, CodeScanPayCode, 10, "二维码错误") // QR CODE 解析失敗
	ScanQRCodeContentFailed          = new(statusBadRequest, CodeScanPayCode, 11, "二维码错误") // QR CODE 內文有錯
	ScanPayRecordExpired             = new(statusBadRequest, CodeScanPayCode, 12, "订单已过期")
	ScanQRCodeImageTooLarge          = new(statusBadRequest, CodeScanPayCode, 13, "二维码图片不可超过4 MiB")
	ScanPayOrderFailure              = new(statusBadRequest, CodeScanPayCode, 14, "订单支付失败已作废，如需再次支付请重新获取")
)
