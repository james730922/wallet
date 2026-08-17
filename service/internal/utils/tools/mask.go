package tools

import (
	"strings"
)

func MaskWalletAddress(number string) string {
	prefixMaskLen := 4
	maskLen := 6
	numberLen := len(number)
	return number[:prefixMaskLen] + "****" + number[numberLen-maskLen:]
}

func MaskPhoneNumber(number string) string {
	mask := "*****"

	numberLen := len(number)
	if numberLen > 5 {
		idx := numberLen / 4
		return number[:idx] + mask + number[idx+len(mask):]
	}

	return mask
}

func MaskQQ(qq string) string {
	mask := "*"

	qqLen := len(qq)
	switch qqLen {
	case 1, 2, 3, 4: // 全部隱藏
		return strings.Repeat(mask, 4)
	case 5: // 只顯示最後1碼
		return strings.Repeat(mask, 4) + qq[4:]
	case 6, 7: // 只顯示後4碼
		return strings.Repeat(mask, qqLen-4) + qq[len(qq)-4:]
	case 8: // 只顯示前1碼和後4碼
		return qq[:1] + strings.Repeat(mask, qqLen-5) + qq[len(qq)-4:]
	default: // 9碼以上，只顯示前2碼和後4碼
		return qq[:2] + strings.Repeat(mask, qqLen-6) + qq[len(qq)-4:]
	}
}

// 姓名隐码规则
// 长度1不码
// 长度2码第2码
// 长度3码以上显示头尾的字, 中间隐码最多显示3码
func MaskName(name string) string {
	mask := "*"
	nameLen := len([]rune(name))

	switch nameLen {
	case 1:
		return name
	case 2:
		return string([]rune(name)[0]) + mask
	default:
		maskLenLimit := 3
		maskString := strings.Repeat(mask, maskTimes(nameLen, maskLenLimit))
		return string([]rune(name)[0]) + maskString + string([]rune(name)[nameLen-1])
	}
}

// 隱碼產生個數
func maskTimes(nameLen, maskLenLimit int) int {
	// 扣掉 最后一码的明字
	times := nameLen - maskLenLimit + 1
	if times > maskLenLimit {
		return maskLenLimit
	}
	return times
}

// 卡號隐码规则
// 长度>6僅顯示後6碼
// 长度2~6隱前半部
// 长度1不隱碼
func MaskCardNumber(cardNumber string) string {
	mask := "*"
	cardNumberRune := []rune(cardNumber)
	cardNumberLen := len(cardNumberRune)

	switch cardNumberLen {
	case 0, 1:
		return cardNumber
	case 2, 3, 4, 5, 6:
		maskLen := cardNumberLen / 2
		maskString := strings.Repeat(mask, maskLen)
		return maskString + string(cardNumberRune[maskLen:])
	default:
		maskString := strings.Repeat(mask, cardNumberLen-6)
		return maskString + string(cardNumberRune[cardNumberLen-6:])
	}
}

// APP卡號隐码规则
// 长度>4僅顯示後4碼
// 长度2~4隱前半部
// 长度1不隱碼
// 隱碼固定5
func AppMaskCardNumber(cardNumber string) string {
	mask := "*"
	cardNumberRune := []rune(cardNumber)
	cardNumberLen := len(cardNumberRune)

	switch cardNumberLen {
	case 0, 1:
		return cardNumber
	case 2, 3, 4:
		maskLen := cardNumberLen / 2
		maskString := strings.Repeat(mask, maskLen)
		return maskString + string(cardNumberRune[maskLen:])
	default:
		maskLen := 5
		maskString := strings.Repeat(mask, maskLen)
		return maskString + string(cardNumberRune[cardNumberLen-4:])
	}
}

// 姓名隐码规则
// 长度1不引碼
// 长度2码以上显示最後的字
// 固定隱碼最大2位
func AppMaskName(name string) string {
	mask := "*"
	nameLen := len([]rune(name))

	switch nameLen {
	case 0, 1:
		return name
	case 2:
		return mask + string([]rune(name)[nameLen-1])
	default:
		maskLenLimit := 2
		maskString := strings.Repeat(mask, maskTimes(nameLen, maskLenLimit))
		return maskString + string([]rune(name)[nameLen-1])
	}
}
