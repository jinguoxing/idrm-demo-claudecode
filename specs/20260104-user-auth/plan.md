# User Authentication Technical Plan

> **Branch**: `20260104-user-auth`
> **Spec Path**: `specs/20260104-user-auth/`
> **Created**: 2026-01-04
> **Status**: Draft

---

## Summary

实现用户注册和登录功能，采用 Go-Zero 微服务架构。使用 bcrypt 进行密码哈希，JWT 作为认证凭证。遵循 IDRM 分层架构（Handler → Logic → Model），使用 idrm-go-base 通用库处理错误、响应、校验和日志。

---

## Technical Context

| Item | Value |
|------|-------|
| **Language** | Go 1.24+ |
| **Framework** | Go-Zero v1.9+ |
| **Storage** | MySQL 8.0 |
| **Cache** | Redis 7.0 |
| **ORM** | GORM / SQLx |
| **Testing** | go test |
| **Common Lib** | idrm-go-base v0.1.0+ |
| **Password Hashing** | bcrypt (golang.org/x/crypto/bcrypt) |
| **JWT Library** | github.com/golang-jwt/jwt/v5 |

---

## 通用库 (idrm-go-base)

**安装**:
```bash
go get github.com/jinguoxing/idrm-go-base@latest
```

### 模块初始化

| 模块 | 初始化方式 |
|------|-----------|
| validator | `validator.Init()` 在 main.go |
| telemetry | `telemetry.Init(cfg)` 在 main.go |
| response | `httpx.SetErrorHandler(response.ErrorHandler)` |
| middleware | `rest.WithMiddlewares(...)` |

### 自定义错误码

| 功能 | 范围 | 位置 |
|------|------|------|
| 用户认证 | 30100-30199 | `internal/errorx/codes.go` |

**具体错误码定义**:
```go
const (
    ErrCodeUserPhoneExists     = 30101 // 手机号已被注册
    ErrCodeUserEmailExists     = 30102 // 邮箱已被注册
    ErrCodeUserInvalidPassword = 30103 // 密码强度不足
    ErrCodeUserNotFound        = 30104 // 用户不存在
    ErrCodeUserInvalidPassword = 30105 // 密码错误
    ErrCodeUserDisabled        = 30106 // 用户已禁用
)
```

### 第三方库确认

| 库 | 原因 | 确认状态 |
|----|------|----------|
| golang.org/x/crypto/bcrypt | 密码哈希（业界标准） | ✅ 已确认 |
| github.com/golang-jwt/jwt/v5 | JWT 生成和验证 | ✅ 已确认 |
| github.com/google/uuid | UUID v7 生成（符合宪法规范） | ✅ 已确认 |

---

## Go-Zero 开发流程

按以下顺序完成技术设计和代码生成：

| Step | 任务 | 方式 | 产出 |
|------|------|------|------|
| 1 | 定义 API 文件 | AI 实现 | `api/doc/auth/user.api` |
| 2 | 生成 Handler/Types | goctl | `api/internal/handler/`, `types/` |
| 3 | 定义 DDL 文件 | AI 手写 | `migrations/auth/users.sql` |
| 4 | 实现 Model 接口 | AI 手写 | `model/auth/user/` |
| 5 | 实现 Logic 层 | AI 实现 | `api/internal/logic/auth/` |

> ⚠️ **重要**：goctl 必须在 `api/doc/api.api` 入口文件上执行，不能针对单个功能文件！

**goctl 命令**:
```bash
# 步骤1：在 api/doc/api.api 中 import 新模块
# 步骤2：执行 goctl 生成代码（针对整个项目）
goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group
```

---

## File Structure

### 文件产出清单

| 序号 | 文件 | 生成方式 | 位置 |
|------|------|----------|------|
| 1 | API 文件 | AI 实现 | `api/doc/auth/user.api` |
| 2 | DDL 文件 | AI 实现 | `migrations/auth/users.sql` |
| 3 | Handler | goctl 生成 | `api/internal/handler/auth/` |
| 4 | Types | goctl 生成 | `api/internal/types/` |
| 5 | Logic | AI 实现 | `api/internal/logic/auth/` |
| 6 | Model | AI 实现 | `model/auth/user/` |

### 代码结构

```
api/internal/
├── handler/auth/
│   ├── register_handler.go       # goctl 生成
│   ├── login_handler.go          # goctl 生成
│   └── routes.go
├── logic/auth/
│   ├── register_logic.go         # AI 实现
│   └── login_logic.go
├── types/
│   └── types.go                  # goctl 生成
└── svc/
    └── servicecontext.go         # 手动维护

model/auth/user/
├── interface.go                  # 接口定义
├── types.go                      # 数据结构
├── vars.go                       # 常量/错误
├── factory.go                    # ORM 工厂
├── gorm_dao.go                   # GORM 实现
└── sqlx_model.go                 # SQLx 实现
```

---

## Architecture Overview

遵循 IDRM 分层架构：

```
HTTP Request → Handler → Logic → Model → Database
```

**注册流程**:
```
POST /api/v1/auth/register
  ↓
Handler: 参数解析、校验
  ↓
Logic: 密码哈希、唯一性检查、创建用户
  ↓
Model: 插入数据库
```

**登录流程**:
```
POST /api/v1/auth/login
  ↓
Handler: 参数解析、校验
  ↓
Logic: 查询用户、密码验证、生成 JWT
  ↓
Model: 查询数据库、更新登录时间
```

| 层级 | 职责 | 最大行数 |
|------|------|----------|
| Handler | 解析参数、格式化响应 | 30 |
| Logic | 业务逻辑实现 | 50 |
| Model | 数据访问 | 50 |

---

## Interface Definitions

### UserModel Interface

```go
package user

import (
    "context"
)

// Model 用户数据访问接口
type Model interface {
    // Insert 插入用户
    Insert(ctx context.Context, data *User) (*User, error)

    // FindByPhone 根据手机号查找用户
    FindByPhone(ctx context.Context, phone string) (*User, error)

    // FindByEmail 根据邮箱查找用户
    FindByEmail(ctx context.Context, email string) (*User, error)

    // FindOne 根据 ID 查找用户
    FindOne(ctx context.Context, id string) (*User, error)

    // Update 更新用户信息
    Update(ctx context.Context, data *User) error

    // UpdateLastLogin 更新最后登录时间
    UpdateLastLogin(ctx context.Context, id string) error

    // WithTx 设置事务
    WithTx(tx interface{}) Model

    // Trans 执行事务
    Trans(ctx context.Context, fn func(ctx context.Context, model Model) error) error
}
```

---

## Data Model

### DDL

**位置**: `migrations/auth/users.sql`

```sql
CREATE TABLE `users` (
    `id` CHAR(36) NOT NULL COMMENT 'ID (UUID v7)',
    `phone` VARCHAR(20) DEFAULT NULL COMMENT '手机号（唯一）',
    `email` VARCHAR(100) DEFAULT NULL COMMENT '邮箱（唯一，小写）',
    `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希（bcrypt）',
    `status` TINYINT NOT NULL DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    `last_login_at` DATETIME DEFAULT NULL COMMENT '最后登录时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间（软删除）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_phone` (`phone`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `idx_status` (`status`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';
```

### Go Struct

```go
package user

import (
    "time"
    "gorm.io/gorm"
)

// User 用户实体
type User struct {
    Id           string         `gorm:"primaryKey;size:36"`           // UUID v7
    Phone        string         `gorm:"size:20;uniqueIndex:uk_phone"` // 手机号
    Email        string         `gorm:"size:100;uniqueIndex:uk_email"`// 邮箱
    PasswordHash string         `gorm:"size:255;not null"`            // 密码哈希
    Status       int8           `gorm:"type:tinyint;not null;default:1"` // 状态
    LastLoginAt  *time.Time     `gorm:"type:datetime"`                // 最后登录时间
    CreatedAt    time.Time      `gorm:"type:datetime;not null"`
    UpdatedAt    time.Time      `gorm:"type:datetime;not null"`
    DeletedAt    gorm.DeletedAt `gorm:"index"`
}
```

---

## API Contract

### 注册接口

**位置**: `api/doc/auth/user.api`

```api
syntax = "v1"

import "../base.api"

type (
    // RegisterReq 注册请求
    RegisterReq {
        Phone    string `json:"phone" validate:"required_without=Email,phone"`
        Email    string `json:"email" validate:"required_without=Phone,email"`
        Password string `json:"password" validate:"required,min=8,max=32"`
    }

    // RegisterResp 注册响应
    RegisterResp {
        Id    string `json:"id"`
        Phone string `json:"phone"`
        Email string `json:"email"`
    }
)

@server(
    prefix: /api/v1/auth
    group: user
)
service idrm-demo-claudecode-api {
    @handler Register
    post /register (RegisterReq) returns (RegisterResp)
}
```

### 登录接口

```api
type (
    // LoginReq 登录请求
    LoginReq {
        Account  string `json:"account" validate:"required"`  // 手机号或邮箱
        Password string `json:"password" validate:"required"`
    }

    // LoginResp 登录响应
    LoginResp {
        AccessToken string `json:"access_token"`
        ExpiresIn   int64  `json:"expires_in"`  // 秒
        TokenType   string `json:"token_type"`  // Bearer
        UserInfo    UserInfo `json:"user_info"`
    }

    // UserInfo 用户信息
    UserInfo {
        Id    string `json:"id"`
        Phone string `json:"phone"`
        Email string `json:"email"`
    }
)

@server(
    prefix: /api/v1/auth
    group: user
)
service idrm-demo-claudecode-api {
    @handler Login
    post /login (LoginReq) returns (LoginResp)
}
```

---

## Testing Strategy

| 类型 | 方法 | 覆盖率 |
|------|------|--------|
| 单元测试 | 表驱动测试，Mock Model | > 80% |
| 集成测试 | 测试数据库 | 核心流程 |

### 测试用例

**注册功能**:
- 正常注册（手机号）
- 正常注册（邮箱）
- 手机号格式错误
- 邮箱格式错误
- 密码强度不足
- 手机号重复
- 邮箱重复
- 参数为空

**登录功能**:
- 手机号登录成功
- 邮箱登录成功
- 账号不存在
- 密码错误
- 用户已禁用
- 参数为空

---

## Security Considerations

1. **密码安全**:
   - 使用 bcrypt 哈希，cost factor = 10
   - 密码长度 8-32 位，强制包含字母和数字
   - 不在日志中记录密码

2. **JWT 安全**:
   - 使用 HS256 算法
   - Token 有效期 2 小时
   - 签名密钥从环境变量读取

3. **防暴力破解**:
   - 登录失败不提示具体原因（账号或密码错误）
   - 考虑添加速率限制（未来增强）

4. **数据隐私**:
   - 邮箱存储前转为小写
   - 手机号脱敏显示（未来增强）

---

## Performance Considerations

| 指标 | 目标 | 优化方案 |
|------|------|----------|
| 注册响应时间 | < 500ms (P99) | 数据库索引优化 |
| 登录响应时间 | < 300ms (P99) | Redis 缓存用户信息 |
| 并发注册 | 支持 1000 TPS | 唯一索引防重复 |

---

## Revision History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-01-04 | - | 初始版本 |
