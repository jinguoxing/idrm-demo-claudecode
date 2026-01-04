# User Authentication Tasks

> **Branch**: `20260104-user-auth`
> **Spec Path**: `specs/20260104-user-auth/`
> **Created**: 2026-01-04
> **Input**: spec.md, plan.md

---

## Task Format

```
[T001] [P?] [Story] Description
```

| 标记 | 含义 |
|------|------|
| `T001` | 任务 ID |
| `[P]` | 可并行执行（不同文件，无依赖） |
| `[US1]` | 关联 User Story 1 |
| `[US2]` | 关联 User Story 2 |

---

## Task Overview

| ID | Task | Story | Status | Parallel | Est. Lines |
|----|------|-------|--------|----------|------------|
| T001 | 环境准备检查 | Setup | ⏸️ | - | - |
| T002 | 安装依赖包 | Setup | ⏸️ | - | - |
| T003 | 创建自定义错误码 | Foundation | ⏸️ | - | 20 |
| T004 | 配置 ServiceContext | Foundation | ⏸️ | - | 30 |
| T005 | 创建 DDL 文件 | US1 | ⏸️ | [P] | 25 |
| T006 | 定义 API 文件 | US1 | ⏸️ | [P] | 50 |
| T007 | 更新 API 入口文件 | US1 | ⏸️ | - | 5 |
| T008 | 生成 Handler/Types | US1 | ⏸️ | - | - |
| T009 | 实现 Model 接口 | US1 | ⏸️ | - | 15 |
| T010 | 实现 Model 类型 | US1 | ⏸️ | [P] | 30 |
| T011 | 实现 Model 常量 | US1 | ⏸️ | [P] | 15 |
| T012 | 实现 GORM DAO | US1 | ⏸️ | - | 80 |
| T013 | 实现 SQLx Model | US1 | ⏸️ | [P] | 60 |
| T014 | 实现 Model 工厂 | US1 | ⏸️ | - | 25 |
| T015 | 实现注册 Logic | US1 | ⏸️ | - | 50 |
| T016 | 实现登录 Logic | US1 | ⏸️ | - | 50 |
| T017 | 注册单元测试 | US1 | ⏸️ | - | 60 |
| T018 | 登录单元测试 | US1 | ⏸️ | [P] | 50 |
| T019 | API 集成测试 | US1 | ⏸️ | - | 80 |
| T020 | 代码检查和格式化 | Polish | ⏸️ | - | - |
| T021 | 测试覆盖率验证 | Polish | ⏸️ | - | - |

**Total**: 21 tasks | **US1**: 15 tasks | **US2**: 0 tasks | **Foundation**: 2 tasks | **Setup**: 2 tasks | **Polish**: 2 tasks

---

## Phase 1: Setup

**目的**: 项目初始化和基础配置

- [x] T001 确认 Go-Zero 项目结构已就绪
- [x] T002 安装第三方依赖包
  ```bash
  go get github.com/jinguoxing/idrm-go-base@latest
  go get golang.org/x/crypto/bcrypt
  go get github.com/golang-jwt/jwt/v5
  go get github.com/google/uuid
  ```

**Checkpoint**: ✅ 开发环境就绪

---

## Phase 2: Foundation

**目的**: 必须完成后才能开始 User Story 实现

- [x] T003 创建 `api/internal/errorx/codes.go` 定义用户认证错误码
- [x] T004 配置 `api/internal/svc/servicecontext.go` 添加数据库连接和配置

**Checkpoint**: ✅ 基础设施就绪，可开始 User Story 实现

---

## Phase 3: User Story 1 - 用户注册 (P1) 🎯 MVP

**目标**: 用户可以使用手机号码或邮箱注册账号

**独立测试**: 用户可以成功提交注册表单，收到注册成功响应

### Step 1: 定义数据模型 (可并行)

- [x] T005 [P] [US1] 创建 `migrations/auth/users.sql` DDL 文件

- [x] T006 [P] [US1] 创建 `api/doc/auth/user.api` API 文件

### Step 2: 更新 API 入口

- [x] T007 [US1] 在 `api/doc/api.api` 中 import auth/user 模块

### Step 3: 生成代码

- [x] T008 [US1] 运行 `goctl api go -api api/doc/api.api -dir api/ --style=go_zero --type-group` 生成 Handler/Types

### Step 4: 实现 Model 层

- [x] T009 [US1] 创建 `model/auth/user/interface.go` 定义 Model 接口

- [x] T010 [P] [US1] 创建 `model/auth/user/types.go` 定义 User 结构体

- [x] T011 [P] [US1] 创建 `model/auth/user/vars.go` 定义常量和错误

- [x] T012 [US1] 实现 `model/auth/user/gorm_dao.go` GORM 实现

- [x] T013 [P] [US1] 实现 `model/auth/user/sqlx_model.go` SQLx 实现

- [x] T014 [US1] 实现 `model/auth/user/factory.go` ORM 工厂函数

### Step 5: 实现 Logic 层

- [x] T015 [US1] 实现 `api/internal/logic/user/registerlogic.go` 注册业务逻辑

- [x] T016 [US1] 实现 `api/internal/logic/user/loginlogic.go` 登录业务逻辑

### Step 6: 测试

- [ ] T017 [US1] 创建 `model/auth/user/gorm_dao_test.go` 单元测试

- [ ] T018 [P] [US1] 创建 `api/internal/logic/auth/registerlogic_test.go` 单元测试

- [ ] T019 [US1] 创建 `api/internal/logic/auth/integration_test.go` 集成测试

**Checkpoint**: ✅ User Story 1 可独立测试和验证

---

## Phase 4: User Story 2 - 用户登录 (P1)

**说明**: User Story 2 (登录功能) 已在 Phase 3 中与注册功能一起实现，因为它们共享相同的数据模型和基础设施。

**独立测试**: 用户可以提交正确的账号和密码，收到登录成功响应

**Checkpoint**: ✅ User Story 2 可独立测试和验证

---

## Phase 5: Polish

**目的**: 收尾工作

- [ ] T020 运行 `golangci-lint run` 检查代码质量
- [ ] T021 确认测试覆盖率 > 80%

---

## Dependencies

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundation)
    ↓
Phase 3 (US1 + US2)  # 注册和登录一起实现
    ↓
Phase 5 (Polish)
```

### 并行执行说明

**Phase 1 内并行**:
- 无（T001 和 T002 顺序执行）

**Phase 2 内并行**:
- 无（T003 和 T004 顺序执行）

**Phase 3 并行机会**:
- **Step 1**: T005 (DDL) 和 T006 (API) 可并行
- **Step 4**: T010 (types) 和 T011 (vars) 和 T013 (SQLx) 可并行
- **Step 6**: T017 和 T018 测试可并行

**跨 Phase**: 无并行，必须按 Phase 顺序完成

---

## Implementation Strategy

### MVP 范围

**最小可行产品 (MVP)**: Phase 1 + Phase 2 + Phase 3

包含功能:
- ✅ 用户注册（手机号/邮箱）
- ✅ 用户登录
- ✅ 基础错误处理
- ✅ 单元测试

### 增量交付

| 迭代 | 内容 | 价值 |
|------|------|------|
| v1.0 | Phase 1-3 | 核心注册登录功能 |
| v1.1 | Phase 5 | 代码质量提升 |

---

## Notes

- 每个 Task 完成后建议提交代码
- 每个 Checkpoint 进行手动验证
- 遵循 Go-Zero 最佳实践
- 所有函数控制在 50 行以内
- 使用 idrm-go-base 通用库
- 密码必须使用 bcrypt 加密
- JWT Token 有效期 2 小时
