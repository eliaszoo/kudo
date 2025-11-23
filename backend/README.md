# Reward System Backend

基于 Go + MySQL 的奖励系统后端服务

## 功能特性

- RESTful API 设计
- 多用户家庭管理
- 多种奖励类型支持（零花钱、时间、积分、自定义）
- 交易记录与余额查询
- 微信消息接入
- MCP 大模型集成
- 幂等性保证
- 事务一致性

## 技术栈

- **语言**: Go 1.21+
- **框架**: Gin
- **数据库**: MySQL 8.0+
- **ORM**: GORM
- **依赖管理**: Go Modules

## 快速开始

### 1. 安装依赖

```bash
cd backend
go mod download
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

### 3. 创建数据库

```sql
CREATE DATABASE reward_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 运行服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动

## API 文档

### 认证

所有 API 请求需要携带 Bearer Token：

```
Authorization: Bearer your_api_token
```

### 核心接口

#### 创建奖励类型
```http
POST /api/v1/reward_types
Content-Type: application/json

{
  "family_id": 1,
  "name": "零花钱",
  "unit_kind": "money",
  "unit_label": "元"
}
```

#### 授予奖励
```http
POST /api/v1/rewards/grant
Content-Type: application/json

{
  "family_id": 1,
  "child_id": 2,
  "reward_type_id": 1,
  "value": 10000,
  "note": "完成作业",
  "idempotency_key": "unique-key-123"
}
```

#### 消费奖励
```http
POST /api/v1/rewards/spend
Content-Type: application/json

{
  "family_id": 1,
  "child_id": 2,
  "reward_type_id": 1,
  "value": 5000,
  "note": "买文具",
  "idempotency_key": "unique-key-456"
}
```

#### 查询余额
```http
GET /api/v1/balances?family_id=1&child_id=2&reward_type_id=1
```

#### 交易记录
```http
GET /api/v1/transactions?family_id=1&child_id=2&limit=20
```

## 数据库结构

主要表结构：

- `families`: 家庭信息
- `users`: 用户信息（监护人/孩子）
- `reward_types`: 奖励类型定义
- `accounts`: 孩子账户（按奖励类型）
- `transactions`: 交易记录
- `audit_logs`: 审计日志

## 错误处理

统一错误响应格式：

```json
{
  "code": 400,
  "message": "错误描述",
  "details": {}
}
```

错误码说明：
- `0`: 成功
- `400`: 参数错误
- `401`: 未认证
- `403`: 无权限
- `404`: 资源不存在
- `409`: 余额不足
- `422`: 意图解析失败
- `500`: 服务器错误

## 开发指南

### 项目结构

```
backend/
├── cmd/server/          # 应用入口
├── internal/
│   ├── api/            # API 处理器
│   ├── config/         # 配置管理
│   ├── db/             # 数据库模型
│   └── services/       # 业务逻辑
├── go.mod              # 依赖管理
└── .env.example        # 环境变量示例
```

### 添加新功能

1. 在 `internal/db/models.go` 中定义数据模型
2. 在 `internal/services/` 中实现业务逻辑
3. 在 `internal/api/` 中添加 API 处理器
4. 在 `internal/api/router.go` 中注册路由

## 测试

```bash
go test ./...
```

## 部署

### Docker 部署

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./server"]
```

### 环境变量

- `DB_DSN`: MySQL 连接字符串
- `WECHAT_TOKEN`: 微信校验 Token
- `API_TOKEN`: API 认证 Token
- `MCP_SERVER_URL`: MCP 服务器地址
- `PORT`: 服务端口（默认 8080）

## 许可证

MIT License