// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/jinguoxing/idrm-demo-claudecode/api/internal/svc"
	"github.com/jinguoxing/idrm-demo-claudecode/api/internal/types"
	"github.com/jinguoxing/idrm-demo-claudecode/model/auth/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// 1. 参数校验
	if req.Account == "" || req.Password == "" {
		return nil, fmt.Errorf("账号和密码不能为空")
	}

	// 2. 查找用户（手机号或邮箱）
	userEntity, err := l.findUser(req.Account)
	if err != nil {
		return nil, fmt.Errorf("账号或密码错误")
	}

	// 3. 验证密码
	if err := l.verifyPassword(req.Password, userEntity.PasswordHash); err != nil {
		return nil, fmt.Errorf("账号或密码错误")
	}

	// 4. 检查用户状态
	if userEntity.Status != user.UserStatusNormal {
		return nil, fmt.Errorf("用户已被禁用")
	}

	// 5. 生成 JWT Token
	accessToken, err := l.generateToken(userEntity.Id)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	// 6. 更新最后登录时间
	err = l.svcCtx.UserModel.UpdateLastLogin(l.ctx, userEntity.Id)
	if err != nil {
		logx.Errorf("更新最后登录时间失败: %v", err)
	}

	logx.Infof("用户登录成功: id=%s, account=%s", userEntity.Id, req.Account)

	// 7. 返回响应
	return &types.LoginResp{
		AccessToken: accessToken,
		ExpiresIn:   7200, // 2 小时
		TokenType:   "Bearer",
		UserInfo: types.UserInfo{
			Id:    userEntity.Id,
			Phone: userEntity.Phone,
			Email: userEntity.Email,
		},
	}, nil
}

// findUser 查找用户（手机号或邮箱）
func (l *LoginLogic) findUser(account string) (*user.User, error) {
	// 首先尝试手机号
	u, err := l.svcCtx.UserModel.FindByPhone(l.ctx, account)
	if err == nil {
		return u, nil
	}
	if err != user.ErrUserNotFound {
		return nil, err
	}

	// 然后尝试邮箱（小写）
	u, err = l.svcCtx.UserModel.FindByEmail(l.ctx, strings.ToLower(account))
	if err == nil {
		return u, nil
	}
	if err != user.ErrUserNotFound {
		return nil, err
	}

	return nil, user.ErrUserNotFound
}

// verifyPassword 验证密码
func (l *LoginLogic) verifyPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

// generateToken 生成 JWT Token
func (l *LoginLogic) generateToken(userId string) (string, error) {
	// TODO: 从配置中获取密钥
	secret := "94dddffea62801cb4d60d45b80a8b745c59d742716faa96ae03cbd12d436263c"

	now := time.Now()
	expiresAt := now.Add(2 * time.Hour)

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     expiresAt.Unix(),
		"iat":     now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
