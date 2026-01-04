// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jinguoxing/idrm-demo-claudecode/api/internal/svc"
	"github.com/jinguoxing/idrm-demo-claudecode/api/internal/types"
	"github.com/jinguoxing/idrm-demo-claudecode/model/auth/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// 1. 参数校验
	if err := l.validateRegisterReq(req); err != nil {
		return nil, err
	}

	// 2. 检查唯一性
	if err := l.checkUniqueness(req); err != nil {
		return nil, err
	}

	// 3. 密码哈希
	passwordHash, err := l.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 4. 创建用户实体
	userEntity := &user.User{
		Id:           uuid.New().String(),
		Phone:        req.Phone,
		Email:        strings.ToLower(req.Email),
		PasswordHash: passwordHash,
		Status:       user.UserStatusNormal,
	}

	// 5. 保存到数据库
	_, err = l.svcCtx.UserModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	logx.Infof("用户注册成功: id=%s, phone=%s, email=%s", userEntity.Id, userEntity.Phone, userEntity.Email)

	// 6. 返回响应
	return &types.RegisterResp{
		Id:    userEntity.Id,
		Phone: userEntity.Phone,
		Email: userEntity.Email,
	}, nil
}

// validateRegisterReq 校验注册请求
func (l *RegisterLogic) validateRegisterReq(req *types.RegisterReq) error {
	// 至少提供手机号或邮箱
	if req.Phone == "" && req.Email == "" {
		return fmt.Errorf("手机号和邮箱至少填写一个")
	}

	// 校验手机号格式
	if req.Phone != "" {
		if !l.isValidPhone(req.Phone) {
			return fmt.Errorf("手机号格式错误")
		}
	}

	// 校验邮箱格式
	if req.Email != "" {
		if !l.isValidEmail(req.Email) {
			return fmt.Errorf("邮箱格式错误")
		}
	}

	// 校验密码强度
	if err := l.validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

// checkUniqueness 检查手机号和邮箱唯一性
func (l *RegisterLogic) checkUniqueness(req *types.RegisterReq) error {
	// 检查手机号
	if req.Phone != "" {
		_, err := l.svcCtx.UserModel.FindByPhone(l.ctx, req.Phone)
		if err == nil {
			return fmt.Errorf("手机号已被注册")
		}
		if err != user.ErrUserNotFound {
			return err
		}
	}

	// 检查邮箱
	if req.Email != "" {
		_, err := l.svcCtx.UserModel.FindByEmail(l.ctx, strings.ToLower(req.Email))
		if err == nil {
			return fmt.Errorf("邮箱已被注册")
		}
		if err != user.ErrUserNotFound {
			return err
		}
	}

	return nil
}

// hashPassword 密码哈希
func (l *RegisterLogic) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// isValidPhone 校验手机号格式
func (l *RegisterLogic) isValidPhone(phone string) bool {
	// 中国大陆 11 位手机号，1 开头
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// isValidEmail 校验邮箱格式
func (l *RegisterLogic) isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// validatePassword 校验密码强度
func (l *RegisterLogic) validatePassword(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return fmt.Errorf("密码长度必须为 8-32 位")
	}

	var (
		hasLetter bool
		hasDigit  bool
	)

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return fmt.Errorf("密码必须包含字母和数字")
	}

	return nil
}
