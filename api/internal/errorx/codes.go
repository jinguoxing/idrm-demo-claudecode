package errorx

// 用户认证模块错误码 (30100-30199)
const (
	// ErrCodeUserPhoneExists 手机号已被注册
	ErrCodeUserPhoneExists = 30101
	// ErrCodeUserEmailExists 邮箱已被注册
	ErrCodeUserEmailExists = 30102
	// ErrCodeUserPasswordWeak 密码强度不足
	ErrCodeUserPasswordWeak = 30103
	// ErrCodeUserNotFound 用户不存在
	ErrCodeUserNotFound = 30104
	// ErrCodeUserPasswordWrong 密码错误
	ErrCodeUserPasswordWrong = 30105
	// ErrCodeUserDisabled 用户已禁用
	ErrCodeUserDisabled = 30106
)
