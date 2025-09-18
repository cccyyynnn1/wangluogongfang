# Yara 安全服务部署指南

## 目录
- [项目概述](#项目概述)
- [系统要求](#系统要求)
- [Windows 部署](#windows-部署)
- [Linux 部署](#linux-部署)
- [Docker 部署](#docker-部署)
- [配置说明](#配置说明)
- [故障排除](#故障排除)

---

## 项目概述

Yara 安全服务是一个基于 Go 语言开发的安全扫描和系统监控服务，提供文件扫描、进程管理、网络监控、注册表操作等功能。服务包含：

- **后端API服务**：运行在端口 8081，提供 RESTful API 接口
- **前端Web界面**：运行在端口 3001，提供用户友好的Web管理界面

本指南将详细介绍在不同环境下的部署步骤。

## 系统要求

### 基础要求
- **操作系统**: Windows 10/11, Linux (Ubuntu 18.04+, CentOS 7+)
- **内存**: 最少 2GB RAM，推荐 4GB+
- **存储**: 最少 1GB 可用空间
- **网络**: 支持 HTTP/HTTPS 访问

### 软件依赖
- **Go**: 1.19 或更高版本
- **Git**: 用于代码获取
- **Docker**: 仅 Docker 部署需要 (20.10+)
- **Web浏览器**: 用于访问Web管理界面

---

## Windows 部署

### 1. 环境准备

#### 安装 Go 语言环境
```powershell
# 下载并安装 Go (从 https://golang.org/dl/)
# 或使用 Chocolatey
choco install golang

# 验证安装
go version
```

#### 安装 Git
```powershell
# 使用 Chocolatey
choco install git

# 或从 https://git-scm.com/ 下载安装
```

### 2. 获取项目代码
```powershell
# 克隆项目
git clone <repository-url>
cd Yara

# 或直接下载源码包并解压
```

### 3. 编译项目
```powershell
# 进入项目目录
cd Yara

# 下载依赖
go mod download

# 编译项目
go build -o bin/yara-server.exe ./cmd/server

# 或使用 Makefile
make build
```

### 4. 配置文件设置
```powershell
# 复制配置文件模板
copy configs\config.yaml.example configs\config.yaml

# 编辑配置文件
notepad configs\config.yaml
```

配置文件示例：
```yaml
server:
  port: 8081
  host: "0.0.0.0"

security:
  yara_rules_path: "rules/"
  scan_timeout: 30s
  max_file_size: 100MB

logging:
  level: "info"
  file: "logs/yara-server.log"
```

### 5. 启动服务

#### 方法一：直接运行
```powershell
# 启动服务
.\bin\yara-server.exe

# 或指定配置文件
.\bin\yara-server.exe -config configs\config.yaml
```

#### 方法二：使用批处理脚本
```powershell
# 运行启动脚本
.\scripts\start_web.bat
```

#### 方法三：作为 Windows 服务运行
```powershell
# 使用 NSSM 安装为服务
nssm install YaraService "C:\path\to\yara-server.exe"
nssm set YaraService AppDirectory "C:\path\to\yara"
nssm start YaraService
```

### 6. 验证部署
```powershell
# 测试服务是否正常运行
curl http://localhost:8081/api/health

# 或使用浏览器访问后端API
start http://localhost:8081

# 访问Web管理界面
start http://localhost:3001
```

---

## Linux 部署

### 1. 环境准备

#### Ubuntu/Debian 系统
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装 Go
sudo apt install golang-go

# 或安装最新版本
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 Git
sudo apt install git -y
```

#### CentOS/RHEL 系统
```bash
# 安装 Go
sudo yum install golang

# 或安装最新版本
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 Git
sudo yum install git -y
```

### 2. 获取项目代码
```bash
# 克隆项目
git clone <repository-url>
cd Yara

# 或下载源码包
wget <download-url>
tar -xzf yara-project.tar.gz
cd Yara
```

### 3. 编译项目
```bash
# 下载依赖
go mod download

# 编译项目
go build -o bin/yara-server ./cmd/server

# 或使用 Makefile
make build
```

### 4. 配置文件设置
```bash
# 复制配置文件模板
cp configs/config.yaml.example configs/config.yaml

# 编辑配置文件
nano configs/config.yaml
```

### 5. 创建系统用户（可选）
```bash
# 创建专用用户
sudo useradd -r -s /bin/false yara

# 设置目录权限
sudo chown -R yara:yara /opt/yara
```

### 6. 启动服务

#### 方法一：直接运行
```bash
# 启动服务
./bin/yara-server

# 或指定配置文件
./bin/yara-server -config configs/config.yaml
```

#### 方法二：使用 Shell 脚本
```bash
# 运行启动脚本
./scripts/start.sh

# 或运行 Web 启动脚本
./scripts/start_web.sh
```

#### 方法三：使用 Systemd 服务
```bash
# 创建服务文件
sudo tee /etc/systemd/system/yara.service > /dev/null <<EOF
[Unit]
Description=Yara Security Service
After=network.target

[Service]
Type=simple
User=yara
WorkingDirectory=/opt/yara
ExecStart=/opt/yara/bin/yara-server
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable yara
sudo systemctl start yara

# 查看服务状态
sudo systemctl status yara
```

### 7. 验证部署
```bash
# 测试服务是否正常运行
curl http://localhost:8081/api/health

# 查看服务日志
sudo journalctl -u yara -f

# 访问Web管理界面
xdg-open http://localhost:3001
```

---

## Docker 部署

### 1. 环境准备

#### 安装 Docker
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install docker.io docker-compose -y
sudo systemctl start docker
sudo systemctl enable docker

# CentOS/RHEL
sudo yum install docker docker-compose -y
sudo systemctl start docker
sudo systemctl enable docker

# 验证安装
docker --version
docker-compose --version
```

### 2. 获取项目代码
```bash
# 克隆项目
git clone <repository-url>
cd Yara
```

### 3. 构建 Docker 镜像

#### 方法一：使用 Dockerfile
```bash
# 构建镜像
docker build -t yara-security-service .

# 查看构建的镜像
docker images
```

#### 方法二：使用 Docker Compose
```bash
# 构建并启动服务
docker-compose up -d --build

# 仅构建镜像
docker-compose build
```

### 4. 配置文件设置
```bash
# 复制配置文件
cp configs/config.yaml.example configs/config.yaml

# 编辑配置文件
nano configs/config.yaml
```

### 5. 启动容器

#### 方法一：直接使用 Docker
```bash
# 运行容器
docker run -d \
  --name yara-service \
  -p 8081:8081 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/rules:/app/rules \
  -v $(pwd)/logs:/app/logs \
  yara-security-service

# 查看容器状态
docker ps
```

#### 方法二：使用 Docker Compose
```bash
# 启动服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 6. 验证部署
```bash
# 测试服务
curl http://localhost:8081/api/health

# 或使用浏览器访问后端API
xdg-open http://localhost:8081

# 访问Web管理界面
xdg-open http://localhost:3001
```

---

## 配置说明

### 主要配置文件

#### configs/config.yaml
```yaml
server:
  port: 8081                    # 服务端口
  host: "0.0.0.0"              # 监听地址
  read_timeout: 30s             # 读取超时
  write_timeout: 30s            # 写入超时

security:
  yara_rules_path: "rules/"     # Yara 规则路径
  scan_timeout: 30s             # 扫描超时
  max_file_size: 100MB          # 最大文件大小
  cache_size: 1000              # 缓存大小
  max_concurrent_scans: 10      # 最大并发扫描数

logging:
  level: "info"                 # 日志级别
  file: "logs/yara-server.log"  # 日志文件
  max_size: 100MB               # 最大日志文件大小
  max_backups: 10               # 最大备份数量

database:
  type: "sqlite"                # 数据库类型
  path: "data/yara.db"          # 数据库路径

cache:
  enabled: true                 # 启用缓存
  ttl: 3600                     # 缓存生存时间（秒）
  max_size: 1000                # 最大缓存条目数
```

### 环境变量配置

可以通过环境变量覆盖配置文件中的设置：

```bash
# 服务配置
export YARA_SERVER_PORT=8081
export YARA_SERVER_HOST=0.0.0.0

# 安全配置
export YARA_SCAN_TIMEOUT=30s
export YARA_MAX_FILE_SIZE=100MB

# 日志配置
export YARA_LOG_LEVEL=info
export YARA_LOG_FILE=logs/yara-server.log
```

---

## 故障排除

### 常见问题

#### 1. 服务启动失败
```bash
# 检查后端端口占用
netstat -tulpn | grep 8081

# 检查前端端口占用
netstat -tulpn | grep 3001

# 检查配置文件
./bin/yara-server -config configs/config.yaml --validate

# 查看详细日志
tail -f logs/yara-server.log
```

#### 2. 权限问题
```bash
# Linux 系统权限
sudo chown -R $USER:$USER /opt/yara
chmod +x bin/yara-server

# Windows 权限
# 以管理员身份运行 PowerShell
```

#### 3. 网络连接问题
```bash
# 检查防火墙设置
sudo ufw allow 8081  # 后端API端口
sudo ufw allow 3001   # 前端Web端口

# 检查 SELinux (CentOS/RHEL)
sudo setsebool -P httpd_can_network_connect 1
```

#### 4. Docker 相关问题
```bash
# 清理 Docker 资源
docker system prune -a

# 重新构建镜像
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

#### 5. 性能问题
```bash
# 检查系统资源
top
htop
iostat

# 调整配置参数
# 增加缓存大小、并发数等
```

### 日志分析

#### 查看服务日志
```bash
# 实时查看日志
tail -f logs/yara-server.log

# 查看错误日志
grep ERROR logs/yara-server.log

# 查看访问日志
grep "HTTP" logs/yara-server.log
```

#### 系统日志
```bash
# Linux 系统日志
sudo journalctl -u yara -f

# Windows 事件查看器
eventvwr.msc
```

### 性能监控

#### 系统监控
```bash
# CPU 和内存使用
top -p $(pgrep yara-server)

# 后端网络连接
netstat -an | grep 8081

# 前端网络连接
netstat -an | grep 3001

# 磁盘 I/O
iotop
```

#### 应用监控
```bash
# 访问性能指标
curl http://localhost:8081/api/metrics

# 重置指标
curl -X POST http://localhost:8081/api/metrics/reset
```

### 备份和恢复

#### 数据备份
```bash
# 备份配置文件
tar -czf backup-config-$(date +%Y%m%d).tar.gz configs/

# 备份数据库
cp data/yara.db backup-yara-$(date +%Y%m%d).db

# 备份日志
tar -czf backup-logs-$(date +%Y%m%d).tar.gz logs/
```

#### 服务恢复
```bash
# 停止服务
sudo systemctl stop yara

# 恢复配置
tar -xzf backup-config-20231201.tar.gz

# 恢复数据库
cp backup-yara-20231201.db data/yara.db

# 启动服务
sudo systemctl start yara
```

---

## 安全建议

### 网络安全
- 使用防火墙限制访问端口
- 配置 HTTPS/TLS 加密
- 实施 IP 白名单访问控制

### 系统安全
- 定期更新系统和依赖包
- 使用专用用户运行服务
- 限制文件系统权限

### 应用安全
- 定期更新 Yara 规则
- 监控异常访问日志
- 实施访问频率限制

---

## 维护指南

### 日常维护
```bash
# 检查服务状态
sudo systemctl status yara

# 查看资源使用
ps aux | grep yara-server

# 清理日志文件
find logs/ -name "*.log" -mtime +30 -delete
```

### 定期维护
```bash
# 更新 Yara 规则
git pull origin main

# 重新编译服务
make build

# 重启服务
sudo systemctl restart yara
```

### 版本升级
```bash
# 备份当前版本
cp -r /opt/yara /opt/yara-backup

# 获取新版本
git pull origin main

# 重新编译
make build

# 测试新版本
./bin/yara-server --test

# 切换服务
sudo systemctl restart yara
```

*最后更新时间：2025年08月07日* 