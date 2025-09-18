# Yara安全服务Web管理界面

这是一个基于Vue 3的现代化Web管理界面，用于管理Yara安全服务。

## 功能特性

- 🎯 **现代化UI**: 基于Vue 3 + Tailwind CSS的响应式设计
- 📊 **实时监控**: 实时显示系统状态、进程、网络连接等信息
- 🔍 **文件扫描**: 支持单文件、目录和缓冲区扫描
- 🛡️ **安全管理**: 安全规则管理、隔离文件处理
- 👥 **用户管理**: 系统用户信息和会话管理
- 📈 **数据可视化**: 使用Chart.js展示统计图表
- 🔄 **自动刷新**: 支持数据自动刷新和手动刷新

## 技术栈

- **前端框架**: Vue 3 (Composition API)
- **路由管理**: Vue Router 4
- **状态管理**: Pinia
- **UI框架**: Tailwind CSS
- **图表库**: Chart.js
- **HTTP客户端**: Axios
- **构建工具**: Vite
- **图标库**: Font Awesome

## 快速开始

### 环境要求

- Node.js >= 16.0.0
- npm >= 8.0.0

### 安装依赖

```bash
npm install
```

### 开发模式

```bash
npm run dev
```

开发服务器将在 `http://localhost:3001` 启动。

### 构建生产版本

```bash
npm run build
```

构建文件将输出到 `dist` 目录。

### 预览生产版本

```bash
npm run preview
```

## 项目结构

```
web/
├── src/
│   ├── components/          # 可复用组件
│   ├── views/              # 页面组件
│   │   ├── Dashboard.vue   # 仪表板
│   │   ├── Scan.vue        # 扫描页面
│   │   ├── Processes.vue   # 进程管理
│   │   ├── Network.vue     # 网络管理
│   │   ├── Security.vue    # 安全管理
│   │   └── Users.vue       # 用户管理
│   ├── stores/             # Pinia状态管理
│   │   └── app.js         # 应用状态
│   ├── services/           # API服务
│   │   └── api.js         # API接口定义
│   ├── router/             # 路由配置
│   │   └── index.js       # 路由定义
│   ├── App.vue            # 根组件
│   ├── main.js            # 应用入口
│   └── style.css          # 全局样式
├── public/                # 静态资源
├── index.html             # HTML模板
├── package.json           # 项目配置
├── vite.config.js         # Vite配置
├── tailwind.config.js     # Tailwind配置
└── README.md             # 项目文档
```

## 页面功能

### 仪表板 (Dashboard)
- 系统状态概览
- 实时统计数据
- 文件扫描功能
- 快速操作入口
- 系统日志显示

### 扫描页面 (Scan)
- 多种扫描模式（文件、目录、缓冲区）
- 扫描历史记录
- 详细的扫描结果展示

### 进程管理 (Processes)
- 系统进程列表
- 进程状态监控
- 进程操作（暂停、恢复、终止）
- 进程详细信息

### 网络管理 (Network)
- 网络连接监控
- 协议过滤和搜索
- 网络接口信息
- 连接详情查看

### 安全管理 (Security)
- 安全规则管理
- 隔离文件处理
- 缓存统计
- 规则引擎状态

### 用户管理 (Users)
- 系统用户列表
- 用户状态管理
- 活跃会话监控
- 用户详细信息

## API接口

项目使用统一的API服务层，主要接口包括：

### 文件相关
- `POST /api/v1/file/scan` - 文件扫描
- `POST /api/v1/file/scan-directory` - 目录扫描
- `POST /api/v1/file/scan-buffer` - 缓冲区扫描

### 进程相关
- `GET /api/v1/process/list` - 获取进程列表
- `GET /api/v1/process/{pid}` - 获取进程详情
- `DELETE /api/v1/process/{pid}` - 终止进程

### 网络相关
- `GET /api/v1/network/connections` - 获取网络连接
- `GET /api/v1/network/interfaces` - 获取网络接口

### 安全相关
- `GET /api/v1/security/status` - 获取安全状态
- `GET /api/v1/security/rules` - 获取安全规则
- `POST /api/v1/security/reload-rules` - 重新加载规则

### 用户相关
- `GET /api/v1/user/current` - 获取当前用户
- `GET /api/v1/user/all` - 获取所有用户
- `GET /api/v1/user/sessions/{username}` - 获取用户会话

## 开发指南

### 添加新页面

1. 在 `src/views/` 目录下创建新的Vue组件
2. 在 `src/router/index.js` 中添加路由配置
3. 在 `src/services/api.js` 中添加相关API接口

### 添加新组件

1. 在 `src/components/` 目录下创建可复用组件
2. 使用Composition API编写组件逻辑
3. 使用Tailwind CSS进行样式设计

### 状态管理

使用Pinia进行状态管理，主要store包括：

- `useAppStore`: 应用全局状态
- 可根据需要添加更多store

### 样式指南

- 使用Tailwind CSS进行样式设计
- 遵循响应式设计原则
- 保持UI组件的一致性

## 部署

### 开发环境

```bash
npm run dev
```

### 生产环境

```bash
npm run build
npm run preview
```

### Docker部署

```bash
# 构建镜像
docker build -t yara-web .

# 运行容器
docker run -p 3001:3001 yara-web
```

## 配置说明

### Vite配置

项目使用Vite作为构建工具，主要配置包括：

- 开发服务器端口: 3001
- API代理配置: 自动代理到后端服务
- 构建优化: 代码分割和压缩

### Tailwind配置

- 自定义颜色主题
- 响应式断点配置
- 自定义动画效果

## 故障排除

### 常见问题

1. **API连接失败**
   - 检查后端服务是否运行在8081端口
   - 确认API代理配置正确

2. **样式不生效**
   - 确认Tailwind CSS已正确安装
   - 检查样式文件导入顺序

3. **图表不显示**
   - 确认Chart.js已正确安装
   - 检查Canvas元素是否正确渲染

### 调试技巧

- 使用Vue DevTools进行组件调试
- 使用浏览器开发者工具查看网络请求
- 查看控制台错误信息

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 推送到分支
5. 创建Pull Request

## 许可证

MIT License

## 联系方式

如有问题或建议，请提交Issue或联系lys17876365245@qq.com。 