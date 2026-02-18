package services

import (
	"errors"
	"regexp"
	"unicode"

	"gorm.io/gorm"
)

// ValidatorService 验证服务
type ValidatorService struct {
	db *gorm.DB
}

func NewValidatorService(db *gorm.DB) *ValidatorService {
	return &ValidatorService{
		db: db,
	}
}

// ValidateAmount 验证金额（积分）
func (s *ValidatorService) ValidateAmount(amount int) error {
	if amount <= 0 {
		return errors.New("金额必须大于0")
	}
	if amount > 100000000 { // 100万积分
		return errors.New("金额超出允许范围")
	}
	return nil
}

// ValidatePhone 验证手机号
func (s *ValidatorService) ValidatePhone(phone string) error {
	if phone == "" {
		return errors.New("手机号不能为空")
	}
	matched, _ := regexp.MatchString(`^1[3-9]\d{9}$`, phone)
	if !matched {
		return errors.New("手机号格式错误")
	}
	return nil
}

// ValidateEmail 验证邮箱
func (s *ValidatorService) ValidateEmail(email string) error {
	if email == "" {
		return errors.New("邮箱不能为空")
	}
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	if !matched {
		return errors.New("邮箱格式错误")
	}
	return nil
}

// ValidatePassword 验证密码强度
func (s *ValidatorService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度至少8位")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("密码必须包含大写字母")
	}
	if !hasLower {
		return errors.New("密码必须包含小写字母")
	}
	if !hasNumber {
		return errors.New("密码必须包含数字")
	}
	if !hasSpecial {
		return errors.New("密码必须包含特殊字符")
	}

	return nil
}
