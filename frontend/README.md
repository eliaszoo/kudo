# Reward System Frontend

基于 React + TypeScript + Vite 的奖励系统前端界面

## 功能特性

- 现代化的用户界面设计
- 响应式布局，支持移动端
- 家庭管理功能
- 奖励类型配置
- 交易记录查看
- 余额查询与统计
- 实时通知提醒

## 技术栈

- **框架**: React 18 + TypeScript
- **构建工具**: Vite
- **路由**: React Router DOM
- **状态管理**: Zustand
- **UI 组件**: Tailwind CSS
- **图标**: Lucide React
- **通知**: Sonner

## 快速开始

### 1. 安装依赖

```bash
cd frontend
pnpm install
```

### 2. 启动开发服务器

```bash
pnpm dev
```

应用将在 `http://localhost:3000` 启动

### 3. 构建生产版本

```bash
pnpm build
```

## 项目结构

```
frontend/
├── src/
│   ├── components/       # 可复用组件
│   ├── pages/           # 页面组件
│   ├── stores/          # 状态管理
│   ├── utils/           # 工具函数
│   ├── App.tsx          # 主应用组件
│   ├── main.tsx         # 应用入口
│   └── index.css        # 全局样式
├── public/              # 静态资源
├── index.html           # HTML 入口
├── package.json         # 依赖配置
├── vite.config.ts       # Vite 配置
├── tailwind.config.js   # Tailwind 配置
└── tsconfig.json        # TypeScript 配置
```

## 页面说明

### 仪表板 (Dashboard)
- 家庭概览统计
- 孩子余额展示
- 快捷操作入口

### 奖励类型 (RewardTypes)
- 奖励类型列表
- 创建新类型
- 编辑和删除功能

### 交易记录 (Transactions)
- 交易历史查看
- 筛选和搜索功能
- 详细的交易信息

## API 集成

前端通过 RESTful API 与后端通信：

- 基础 URL: `http://localhost:8080/api/v1`
- 认证方式: Bearer Token
- 错误处理: 统一错误码和消息

### 示例 API 调用

```typescript
import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('token')}`,
  },
})

// 获取余额
const response = await api.get('/balances', {
  params: { family_id: 1, child_id: 2, reward_type_id: 1 }
})
```

## 状态管理

使用 Zustand 进行状态管理：

- `authStore`: 用户认证状态
- `familyStore`: 家庭数据管理
- `rewardStore`: 奖励相关状态

## 样式系统

基于 Tailwind CSS：

- 响应式设计
- 自定义颜色主题
- 组件级样式
- 暗色模式支持（可选）

## 开发指南

### 添加新页面

1. 在 `src/pages/` 创建页面组件
2. 在 `src/App.tsx` 添加路由
3. 更新导航菜单（如需要）

### 创建组件

1. 在 `src/components/` 创建组件文件
2. 遵循组件命名规范（PascalCase）
3. 添加必要的 TypeScript 类型定义

### 状态管理

1. 在 `src/stores/` 创建新的 store
2. 使用 Zustand 的 create 函数
3. 添加必要的 actions 和 selectors

## 部署

### 构建生产版本

```bash
pnpm build
```

构建输出在 `dist/` 目录

### 部署选项

- **静态托管**: Netlify, Vercel, GitHub Pages
- **容器部署**: Docker + Nginx
- **CDN**: 阿里云 OSS, 腾讯云 COS

### Docker 部署示例

```dockerfile
FROM node:18-alpine as builder

WORKDIR /app
COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 环境变量

- `VITE_API_BASE_URL`: API 基础地址
- `VITE_APP_TITLE`: 应用标题

## 浏览器支持

- Chrome (推荐)
- Firefox
- Safari
- Edge

## 许可证

MIT License