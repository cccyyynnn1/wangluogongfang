# Yara 安全服务功能 API 接口实现文档

## 📋 目录索引

### 🏥 系统监控模块

- [1.1 健康检查](#11-健康检查)
- [1.2 获取性能指标](#12-获取性能指标)
- [1.3 重置性能指标](#13-重置性能指标)

### 🛡️ 安全扫描模块

- [2.1 扫描目录](#21-扫描目录)
- [2.2 扫描单个文件](#22-扫描单个文件)
- [2.3 扫描内存缓冲区](#23-扫描内存缓冲区)
- [2.4 获取安全状态](#24-获取安全状态)
- [2.5 获取规则信息](#25-获取规则信息)
- [2.6 重新加载规则](#26-重新加载规则)
- [2.7 获取缓存统计](#27-获取缓存统计)
- [2.8 清空所有缓存](#28-清空所有缓存)
- [2.9 获取隔离列表](#29-获取隔离列表)
- [2.10 隔离文件](#210-隔离文件)
- [2.11 恢复文件](#211-恢复文件)
- [2.12 获取扫描历史](#212-获取扫描历史)

### 📁 文件管理模块

- [3.1 获取文件列表](#31-获取文件列表)
- [3.2 获取文件详细信息](#32-获取文件详细信息)
- [3.3 获取文件哈希值](#33-获取文件哈希值)
- [3.4 验证文件哈希](#34-验证文件哈希)
- [3.5 获取文件所有哈希值](#35-获取文件所有哈希值)
- [3.6 复制文件](#36-复制文件)
- [3.7 移动文件](#37-移动文件)
- [3.8 删除文件](#38-删除文件)

### 🔍 进程管理模块

- [4.1 获取系统进程列表](#41-获取系统进程列表)
- [4.2 获取进程详细信息](#42-获取进程详细信息)
- [4.3 启动进程](#43-启动进程)
- [4.4 挂起进程](#44-挂起进程)
- [4.5 恢复进程](#45-恢复进程)
- [4.6 结束进程](#46-结束进程)
- [4.7 获取进程网络连接](#47-获取进程网络连接)
- [4.8 获取进程内存信息](#48-获取进程内存信息)
- [4.9 检查进程运行状态](#49-检查进程运行状态)
- [4.10 获取进程统计信息](#410-获取进程统计信息)
- [4.11 获取进程子进程](#411-获取进程子进程)
- [4.12 获取进程模块列表](#412-获取进程模块列表)
- [4.13 获取系统模块列表](#413-获取系统模块列表)
- [4.14 获取模块详细信息](#414-获取模块详细信息)
- [4.15 挂起系统模块](#415-挂起系统模块)
- [4.16 恢复系统模块](#416-恢复系统模块)
- [4.17 结束系统模块](#417-结束系统模块)
- [4.18 启用进程监控](#418-启用进程监控)
- [4.19 禁用进程监控](#419-禁用进程监控)
- [4.20 获取监控进程列表](#420-获取监控进程列表)

### 🔧 注册表管理模块

- [5.1 获取注册表键信息](#51-获取注册表键信息)
- [5.2 列出注册表键](#52-列出注册表键)
- [5.3 列出注册表值](#53-列出注册表值)
- [5.4 创建注册表键](#54-创建注册表键)
- [5.5 删除注册表键](#55-删除注册表键)
- [5.6 获取注册表值](#56-获取注册表值)
- [5.7 设置注册表值](#57-设置注册表值)
- [5.8 删除注册表值](#58-删除注册表值)
- [5.9 搜索注册表](#59-搜索注册表)

### 🌐 网络管理模块

- [6.1 获取所有网络连接](#61-获取所有网络连接)
- [6.2 获取 TCP 连接](#62-获取TCP连接)
- [6.3 获取 UDP 连接](#63-获取UDP连接)
- [6.4 按进程 ID 获取连接](#64-按进程ID获取连接)
- [6.5 按端口获取连接](#65-按端口获取连接)
- [6.6 按 IP 地址获取连接](#66-按IP地址获取连接)
- [6.7 根据连接 ID 获取连接信息](#67-根据连接ID获取连接信息)
- [6.8 关闭网络连接](#68-关闭网络连接)
- [6.9 获取正在使用的端口](#69-获取正在使用的端口)
- [6.10 检查端口占用](#610-检查端口占用)
- [6.11 获取监听端口](#611-获取监听端口)
- [6.12 获取已建立连接](#612-获取已建立连接)
- [6.13 获取网络接口](#613-获取网络接口)
- [6.14 获取网络统计信息](#614-获取网络统计信息)
- [6.15 启用网络监控](#615-启用网络监控)
- [6.16 禁用网络监控](#616-禁用网络监控)
- [6.17 获取监控连接列表](#617-获取监控连接列表)
- [6.18 获取连接历史记录](#618-获取连接历史记录)

### 👤 用户管理模块

- [7.1 获取当前用户信息](#71-获取当前用户信息)
- [7.2 获取所有用户列表](#72-获取所有用户列表)
- [7.3 根据用户 ID 获取用户信息](#73-根据用户ID获取用户信息)
- [7.4 根据用户名获取用户信息](#74-根据用户名获取用户信息)
- [7.5 获取用户组](#75-获取用户组)
- [7.6 检查用户权限](#76-检查用户权限)
- [7.7 验证用户密码](#77-验证用户密码)
- [7.8 获取密码策略](#78-获取密码策略)
- [7.9 检查账户状态](#79-检查账户状态)
- [7.10 获取用户会话](#710-获取用户会话)
- [7.11 结束用户会话](#711-结束用户会话)
- [7.12 锁定用户账户](#712-锁定用户账户)
- [7.13 解锁用户账户](#713-解锁用户账户)
- [7.14 修改用户密码](#714-修改用户密码)
- [7.15 获取用户登录历史](#715-获取用户登录历史)

---

## 🏥 系统监控模块

### 1.1 健康检查

**接口地址**: `GET /api/health`

**功能描述**: 检查服务运行状态

**请求参数**: 无

**响应状态码**:

- `200` - 服务正常运行

**响应格式**: `application/json`

**标签**: 系统监控

**相关类型定义**:

- [`MetricsCollector`](#type-metricscollector) - 指标收集器结构

**相关方法签名**:

- [`MetricsCollector.HealthCheck()`](#func-mc-metricscollector-healthcheck) - 指标收集器健康检查方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "服务正常",
    "time": "2025-09-14T20:34:50+08:00",
    "data": {
        "status": "healthy",
        "version": "1.0.0"
    }
}
```

---

### 1.2 获取性能指标

**接口地址**: `GET /api/metrics`

**功能描述**: 获取系统性能监控指标

**请求参数**: 无

**响应状态码**:

- `200` - 性能指标数据

**响应格式**: `application/json`

**标签**: 系统监控

**相关类型定义**:

- [`MetricsCollector`](#type-metricscollector) - 指标收集器结构

**相关方法签名**:

- [`MetricsCollector.GetMetrics()`](#func-mc-metricscollector-getmetrics) - 指标收集器获取指标方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取指标成功",
    "time": "2025-09-14T20:35:04+08:00",
    "data": {
        "TotalRequests": 1,
        "SuccessfulRequests": 1,
        "FailedRequests": 0,
        "TimeoutRequests": 0,
        "TotalResponseTime": 0,
        "MinResponseTime": 0,
        "MaxResponseTime": 0,
        "AvgResponseTime": 0,
        "CurrentConcurrency": 1,
        "MaxConcurrency": 1,
        "ErrorCounts": {},
        "EndpointStats": {
            "GET /api/health": {
                "Count": 1,
                "TotalTime": 0,
                "MinTime": 0,
                "MaxTime": 0,
                "ErrorCount": 0,
                "LastAccess": "2025-09-14T20:34:50.931328+08:00"
            }
        }
    }
}
```

---

### 1.3 重置性能指标

**接口地址**: `POST /api/metrics/reset`

**功能描述**: 重置所有性能监控指标

**请求参数**: 无

**响应状态码**:

- `200` - 指标重置成功

**响应格式**: `application/json`

**标签**: 系统监控

**相关类型定义**:

- [`MetricsCollector`](#type-metricscollector) - 指标收集器结构

**相关方法签名**:

- [`MetricsCollector.ResetMetrics()`](#func-mc-metricscollector-resetmetrics) - 指标收集器重置指标方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "指标重置成功",
    "time": "2025-09-14T20:35:21+08:00",
    "data": {}
}
```

---

## 🛡️ 安全扫描模块

### 2.1 扫描目录

**接口地址**: `POST /api/v1/file/scan-directory`

**功能描述**: 扫描指定目录下的所有文件

**请求参数**:

```json
{
  "path": "D:\\Code\\Go\\workspace\\Yara-Security\\Yara\\test_files",
  "recursive": true,
  "max_depth": 1,
  "include_patterns": ["*.txt", "*.bat", "*.ps1"],
  "exclude_patterns": ["*.tmp", "*.log"]
}
```

**参数说明**:

- `path` (string, 必需): 目录路径
- `recursive` (boolean, 可选): 是否递归扫描，默认 true
- `max_depth` (integer, 可选): 最大扫描深度
- `include_patterns` (array, 可选): 包含的文件模式
- `exclude_patterns` (array, 可选): 排除的文件模式

**响应状态码**:

- `200` - 目录扫描完成

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`ScanFileRequest`](#type-scanfilerequest) - 文件扫描请求结构
- [`ScanResult`](#type-scanresult) - 扫描结果结构
- [`ThreatInfo`](#type-threatinfo) - 威胁信息结构

**相关方法签名**:

- [`FileHandler.ScanDirectory()`](#func-h-filehandler-scandirectory) - 文件处理器目录扫描方法
- [`SecurityService.ScanDirectory()`](#func-s-securityservice-scandirectory) - 安全服务目录扫描方法
- [`Scanner.ScanDirectory()`](#func-s-scanner-scandirectory) - 扫描器目录扫描方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "目录扫描完成",
    "data": {
        "count": 6,
        "directory": "D:\\Code\\Go\\workspace\\Yara-Security\\Yara\\test_files",
        "exclude_patterns": [
            "*.tmp",
            "*.log"
        ],
        "include_patterns": [
            "*.txt",
            "*.bat",
            "*.ps1"
        ],
        "max_depth": 1,
        "recursive": true,
        "results": [
             {
                "file_path": "D:\\Code\\Go\\workspace\\Yara-Security\\Yara\\test_files\\test_powershell.ps1",
                "is_infected": true,
                "threats": [
                    {
                        "rule_name": "Consolidated_Network_Threat_Detection",
                        "description": "网络威胁检测规则",
                        "severity": "High",
                        "category": "Network",
                        "tags": "network,url,ip,malicious,consolidated"
                    },
                    {
                        "rule_name": "Consolidated_Shellcode_Detection",
                        "description": "Shellcode检测规则2",
                        "severity": "Critical",
                        "category": "Shellcode",
                        "tags": "shellcode,assembly,hex,critical,consolidated"
                    },
                    {
                        "rule_name": "Malicious_Char_Pattern_21",
                        "description": "敏感信息",
                        "severity": "medium",
                        "category": "sensitive_info",
                        "tags": "character,malicious"
                    }
                ],
                "scan_time": "2025-09-14T20:35:30.7566621+08:00",
                "scan_duration": 14421100,
                "file_info": {
                    "path": "D:\\Code\\Go\\workspace\\Yara-Security\\Yara\\test_files\\test_powershell.ps1",
                    "name": "test_powershell.ps1",
                    "size": 1143,
                    "is_dir": false,
                    "mod_time": "2025-08-12T10:09:08+08:00",
                    "create_time": "2025-08-25T01:20:33.128063+08:00",
                    "access_time": "2025-09-14T20:35:30.7454888+08:00",
                    "permissions": "-rw-rw-rw-",
                    "owner": "DESKTOP-ASUS\\ASUS",
                    "group": "default_group",
                    "md5": "c586300a1b004762718fdb0fd297920c",
                    "sha256": "0e125530937c6cf0e80dd15f5b98d41013269c9a0a8bce04655515fffba0b266",
                    "is_suspicious": true,
                    "threat_level": "Medium"
                },
                "metadata": {
                    "file_size": "1143",
                    "file_type": "script",
                    "scan_engine": "yara",
                    "threat_count": "2"
                }
            }
        ],
        "scan_time": "2025-09-14T20:35:30+08:00",
        "total_scanned": 6
    },
    "time": "2025-09-14T20:35:30.7571747+08:00"
}
```

---

### 2.2 扫描单个文件

**接口地址**: `POST /api/v1/file/scan`

**功能描述**: 使用 Yara 规则扫描指定文件

**请求参数**:

```json
{
    "path": "test_files/test_powershell.ps1"
}
```

**参数说明**:

- `path` (string, 必需): 要扫描的文件路径

**响应状态码**:

- `200` - 扫描完成

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`ScanFileRequest`](#type-scanfilerequest) - 文件扫描请求结构
- [`ScanResult`](#type-scanresult) - 扫描结果结构
- [`ThreatInfo`](#type-threatinfo) - 威胁信息结构

**相关方法签名**:

- [`FileHandler.ScanFile()`](#func-h-filehandler-scanfile) - 文件处理器扫描方法
- [`SecurityService.ScanFile()`](#func-s-securityservice-scanfile) - 安全服务文件扫描方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "扫描完成",
    "time": "2025-09-14T20:38:49+08:00",
    "data": {
        "file_path": "test_files/test_powershell.ps1",
        "is_infected": true,
        "threats": [
            {
                "rule_name": "Consolidated_Network_Threat_Detection",
                "description": "网络威胁检测规则",
                "severity": "High",
                "category": "Network",
                "tags": "network,url,ip,malicious,consolidated"
            },
            {
                "rule_name": "Consolidated_Shellcode_Detection",
                "description": "Shellcode检测规则2",
                "severity": "Critical",
                "category": "Shellcode",
                "tags": "shellcode,assembly,hex,critical,consolidated"
            },
            {
                "rule_name": "Malicious_Char_Pattern_21",
                "description": "敏感信息",
                "severity": "medium",
                "category": "sensitive_info",
                "tags": "character,malicious"
            }
        ],
        "scan_time": "2025-09-14T20:38:49.571953+08:00",
        "scan_duration": 3348500,
        "file_info": {
            "path": "test_files/test_powershell.ps1",
            "name": "test_powershell.ps1",
            "size": 1143,
            "is_dir": false,
            "mod_time": "2025-08-12T10:09:08+08:00",
            "create_time": "2025-08-25T01:20:33.128063+08:00",
            "access_time": "2025-09-14T20:38:49.5698887+08:00",
            "permissions": "-rw-rw-rw-",
            "owner": "DESKTOP-ASUS\\ASUS",
            "group": "default_group",
            "md5": "c586300a1b004762718fdb0fd297920c",
            "sha256": "0e125530937c6cf0e80dd15f5b98d41013269c9a0a8bce04655515fffba0b266",
            "is_suspicious": true,
            "threat_level": "Medium"
        },
        "metadata": {
            "file_size": "1143",
            "file_type": "script",
            "scan_engine": "yara",
            "threat_count": "2"
        }
    }
}
```

---

### 2.3 扫描内存缓冲区

**接口地址**: `POST /api/v1/file/scan-buffer`

**功能描述**: 扫描内存中的缓冲区内容

**请求参数**:

```json
{
  "identifier": "network_payload",
  "data": "http://baidu.com"
}
```

**参数说明**:

- `identifier` (string, 可选): 缓冲区标识符
- `data` (string, 必需): 要扫描的数据内容

**响应状态码**:

- `200` - 缓冲区扫描完成

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`BufferScanRequest`](#type-bufferscanrequest) - 缓冲区扫描请求结构
- [`ScanResult`](#type-scanresult) - 扫描结果结构
- [`ThreatInfo`](#type-threatinfo) - 威胁信息结构

**相关方法签名**:

- [`FileHandler.ScanBuffer()`](#func-h-filehandler-scanbuffer) - 文件处理器缓冲区扫描方法
- [`SecurityService.ScanBuffer()`](#func-s-securityservice-scanbuffer) - 安全服务缓冲区扫描方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "扫描完成",
    "time": "2025-09-14T20:39:16+08:00",
    "data": {
        "file_path": "network_payload",
        "is_infected": false,
        "threats": null,
        "scan_time": "2025-09-14T20:39:16.1887265+08:00",
        "scan_duration": 0,
        "file_info": {
            "path": "network_payload",
            "name": "",
            "size": 16,
            "is_dir": false,
            "mod_time": "0001-01-01T00:00:00Z",
            "create_time": "0001-01-01T00:00:00Z",
            "access_time": "0001-01-01T00:00:00Z",
            "permissions": "",
            "owner": "",
            "group": ""
        },
        "metadata": {
            "buffer_size": "16",
            "scan_type": "memory",
            "threat_count": "0"
        }
    }
}
```

---

### 2.4 获取安全状态

**接口地址**: `GET /api/v1/security/status`

**功能描述**: 获取系统安全状态信息

**请求参数**: 无

**响应状态码**:

- `200` - 安全状态信息

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`SecurityStatus`](#type-securitystatus) - 安全状态结构

**相关方法签名**:

- [`SecurityHandler.GetSecurityStatus()`](#func-h-securityhandler-getsecuritystatus) - 安全处理器状态获取方法
- [`SecurityService.GetSecurityStatus()`](#func-s-securityservice-getsecuritystatus) - 安全服务状态获取方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取状态成功",
    "time": "2025-09-14T20:39:26+08:00",
    "data": {
        "is_protected": true,
        "last_scan_time": "2025-09-14T20:39:16.1892531+08:00",
        "threat_count": 5,
        "quarantine_count": 1,
        "real_time_protection": true,
        "definitions_version": "1.0.0",
        "last_update_time": "2025-09-14T20:39:26.4802051+08:00",
        "scan_statistics": {
            "total_scans": 8,
            "infected_files": 5,
            "clean_files": 3,
            "failed_scans": 0,
            "total_scan_time": 18296200
        },
        "rules_info": {
            "config_version": "4.3",
            "rule_manager": {
                "config_version": "4.3",
                "files": [
                    {
                        "description": "统一检测规则（已完全整合到consolidated_detection.yar，避免重复）",
                        "enabled": false,
                        "filename": "consolidated_detection.yar",
                        "hash": "fe95202881be1cf8270ef25c85aff9e8",
                        "load_status": "loaded",
                        "load_time": "2025-09-14T19:33:29+08:00",
                        "priority": 6,
                        "rule_count": 9,
                        "size": 31197
                    },
                    {
                        "description": "统一检测规则（已完全整合到consolidated_detection.yar，避免重复）",
                        "enabled": false,
                        "filename": "yara.rule",
                        "hash": "0cf05631633ffe2ac131d8be11226c65",
                        "load_status": "loaded",
                        "load_time": "2025-09-14T19:33:29+08:00",
                        "priority": 6,
                        "rule_count": 11,
                        "size": 11757
                    }
                ],
                "last_reload": "0001-01-01T00:00:00Z",
                "reload_count": 0,
                "rules_dir": "D:\\Code\\Go\\workspace\\Yara-Security\\Yara\\rules",
                "status": "loaded",
                "total_files": 2
            },
            "rules_dir": "D:\\Code\\Go\\workspace\\Yara-Security\\Yara\\rules",
            "status": "loaded",
            "yara_engine": {
                "cache_hit_rate": 1,
                "cache_size": 7,
                "enabled_rules": 9,
                "total_rules": 9
            }
        }
    }
}
```

---

### 2.5 获取规则信息

**接口地址**: `GET /api/v1/security/rules`

**功能描述**: 获取当前加载的安全规则信息

**请求参数**: 无

**响应状态码**:

- `200` - 规则信息

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`RulesInfo`](#type-rulesinfo) - 规则信息结构
- [`RuleInfo`](#type-ruleinfo) - 规则信息结构

**相关方法签名**:

- [`SecurityHandler.GetRulesInfo()`](#func-h-securityhandler-getrulesinfo) - 安全处理器规则获取方法
- [`SecurityService.GetRulesInfo()`](#func-s-securityservice-getrulesinfo) - 安全服务规则获取方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T20:39:35+08:00",
    "data": {
        "rules": [
            {
                "category": "Malware",
                "description": "恶意软件检测规则",
                "enabled": true,
                "metadata": {
                    "author": "LYS",
                    "date": "2025-08-10",
                    "version": "4.0"
                },
                "name": "Consolidated_Malware_Detection",
                "priority": 999,
                "severity": "High",
                "source_file": "consolidated_detection.yar",
                "strings": 59,
                "tags": [
                    "malware",
                    "generic",
                    "enhanced",
                    "consolidated"
                ]
            },
            {
                "category": "Script",
                "description": "恶意脚本检测规则",
                "enabled": true,
                "metadata": {
                    "author": "LYS",
                    "date": "2025-08-10",
                    "version": "4.0"
                },
                "name": "Consolidated_Script_Detection",
                "priority": 999,
                "severity": "High",
                "source_file": "consolidated_detection.yar",
                "strings": 44,
                "tags": [
                    "script",
                    "powershell",
                    "batch",
                    "vbs",
                    "malicious",
                    "consolidated"
                ]
            },
            {
                "category": "Network",
                "description": "网络威胁检测规则",
                "enabled": true,
                "metadata": {
                    "author": "LYS",
                    "date": "2025-08-10",
                    "version": "4.0"
                },
                "name": "Consolidated_Network_Threat_Detection",
                "priority": 999,
                "severity": "High",
                "source_file": "consolidated_detection.yar",
                "strings": 62,
                "tags": [
                    "network",
                    "url",
                    "ip",
                    "malicious",
                    "consolidated"
                ]
            },
            {
                "category": "Shellcode",
                "description": "Shellcode检测规则2",
                "enabled": true,
                "metadata": {
                    "author": "LYS",
                    "date": "2025-08-10",
                    "version": "4.0"
                },
                "name": "Consolidated_Shellcode_Detection",
                "priority": 999,
                "severity": "Critical",
                "source_file": "consolidated_detection.yar",
                "strings": 57,
                "tags": [
                    "shellcode",
                    "assembly",
                    "hex",
                    "critical",
                    "consolidated"
                ]
            },
            {
                "category": "Ransomware",
                "description": "勒索软件检测规则",
                "enabled": true,
                "metadata": {
                    "author": "LYS",
                    "date": "2025-08-10",
                    "version": "4.0"
                },
                "name": "Consolidated_Ransomware_Detection",
                "priority": 999,
                "severity": "Critical",
                "source_file": "consolidated_detection.yar",
                "strings": 33,
                "tags": [
                    "ransomware",
                    "encryption",
                    "crypto",
                    "critical",
                    "consolidated"
                ]
            },
            {
                "category": "Backdoor",
                "description": "后门程序检测规则",
                "enabled": true,
                "metadata": {
                    "author": "LYS",
                    "date": "2025-08-10",
                    "version": "4.0"
                },
                "name": "Consolidated_Backdoor_Detection",
                "priority": 999,
                "severity": "High",
                "source_file": "consolidated_detection.yar",
                "strings": 31,
                "tags": [
                    "backdoor",
                    "network",
                    "shell",
                    "high",
                    "consolidated"
                ]
            },
            {
                "category": "Keylogger",
                "description": "键盘记录器检测规则",
                "enabled": true,
                "metadata": {
                    "author": "LYS",
                    "date": "2025-08-10",
                    "version": "4.0"
                },
                "name": "Consolidated_Keylogger_Detection",
                "priority": 999,
                "severity": "High",
                "source_file": "consolidated_detection.yar",
                "strings": 31,
                "tags": [
                    "keylogger",
                    "keyboard",
                    "hook",
                    "high",
                    "consolidated"
                ]
            },
            {
                "category": "Worm",
                "description": "网络蠕虫检测规则",
                "enabled": true,
                "metadata": {
                    "author": "LYS",
                    "date": "2025-08-10",
                    "version": "4.0"
                },
                "name": "Consolidated_Worm_Detection",
                "priority": 999,
                "severity": "High",
                "source_file": "consolidated_detection.yar",
                "strings": 31,
                "tags": [
                    "worm",
                    "network",
                    "scan",
                    "propagate",
                    "high",
                    "consolidated"
                ]
            },
            {
                "category": "Sensitive_Info",
                "description": "敏感信息检测规则",
                "enabled": true,
                "metadata": {
                    "author": "LYS",
                    "date": "2025-08-10",
                    "version": "4.0"
                },
                "name": "Consolidated_Sensitive_Info_Detection",
                "priority": 999,
                "severity": "Medium",
                "source_file": "consolidated_detection.yar",
                "strings": 27,
                "tags": [
                    "sensitive",
                    "info",
                    "data",
                    "medium",
                    "consolidated"
                ]
            }
        ],
        "total_count": 9
    }
}
```

---

### 2.6 重新加载规则

**接口地址**: `POST /api/v1/security/reload-rules`

**功能描述**: 重新加载安全扫描规则

**请求参数**: 无

**响应状态码**:

- `200` - 规则重新加载成功

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`RulesInfo`](#type-rulesinfo) - 规则信息结构

**相关方法签名**:

- [`SecurityHandler.ReloadRules()`](#func-h-securityhandler-reloadrules) - 安全处理器重新加载规则方法
- [`SecurityService.ReloadRules()`](#func-s-securityservice-reloadrules) - 安全服务重新加载规则方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "重新加载规则任务已提交成功",
    "time": "2025-09-14T20:39:44+08:00",
    "data": {}
}
```

---

### 2.7 获取缓存统计

**接口地址**: `GET /api/v1/security/cache/stats`

**功能描述**: 获取缓存统计信息

**请求参数**: 无

**响应状态码**:

- `200` - 缓存统计信息

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`SecurityStatus`](#type-securitystatus) - 安全状态结构

**相关方法签名**:

- [`SecurityHandler.GetCacheStats()`](#func-h-securityhandler-getcachestats) - 安全处理器缓存统计获取方法
- [`SecurityService.GetCacheStats()`](#func-s-securityservice-getcachestats) - 安全服务缓存统计获取方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T20:39:53+08:00",
    "data": {
        "cache_size": 2,
        "cache_ttl": "10m0s",
        "eviction_count": 0,
        "hit_count": 0,
        "hit_rate": 0,
        "last_cleanup": "2025-09-14T20:38:29.396821+08:00",
        "max_cache_size": 1000,
        "memory_usage": "0.00MB",
        "miss_count": 1,
        "miss_rate": 1,
        "total_requests": 1
    }
}
```

---

### 2.8 清空所有缓存

**接口地址**: `POST /api/v1/security/cache/clear`

**功能描述**: 清空安全扫描缓存

**请求参数**: 无

**响应状态码**:

- `200` - 缓存清空成功

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`SecurityStatus`](#type-securitystatus) - 安全状态结构

**相关方法签名**:

- [`SecurityHandler.ClearCache()`](#func-h-securityhandler-clearcache) - 安全处理器清空缓存方法
- [`SecurityService.ClearCache()`](#func-s-securityservice-clearcache) - 安全服务清空缓存方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "清除缓存成功",
    "time": "2025-09-14T20:40:01+08:00",
    "data": {}
}
```

---

### 2.9 获取隔离列表

**接口地址**: `GET /api/v1/security/quarantine/list`

**功能描述**: 获取隔离区文件列表

**请求参数**: 无

**响应状态码**:

- `200` - 隔离文件列表

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`QuarantineFileRequest`](#type-quarantinefilerequest) - 隔离文件请求结构
- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`SecurityHandler.GetQuarantineList()`](#func-h-securityhandler-getquarantinelist) - 安全处理器隔离列表获取方法
- [`SecurityService.GetQuarantineList()`](#func-s-securityservice-getquarantinelist) - 安全服务隔离列表获取方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T20:40:08+08:00",
    "data": {
        "count": 1,
        "data": [
            {
                "file_path": "D:\\code\\Go\\Yara\\quarantine\\ceshi.exe_1754792755_quarantined",
                "original_path": "C:\\temp\\test\\ceshi.exe",
                "quarantine_time": "2025-08-10T10:25:55.3344294+08:00",
                "reason": "检测到恶意软件",
                "status": "quarantined"
            }
        ]
    }
}
```

---

### 2.10 隔离文件

**接口地址**: `POST /api/v1/security/quarantine`

**功能描述**: 将可疑或恶意文件移动到隔离区

**请求参数**:

```json
{
    "file_path": "D:\\Code\\Go\\workspace\\Yara-Security\\Yara\\test_files\\隔离.txt"
}
```

**参数说明**:

- `file_path` (string, 必需): 要隔离的文件路径
- `reason` (string, 可选): 隔离原因
- `threat_level` (string, 可选): 威胁等级

**响应状态码**:

- `200` - 文件隔离成功

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`QuarantineFileRequest`](#type-quarantinefilerequest) - 隔离文件请求结构
- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`SecurityHandler.QuarantineFile()`](#func-h-securityhandler-quarantinefile) - 安全处理器隔离文件方法
- [`SecurityService.QuarantineFile()`](#func-s-securityservice-quarantinefile) - 安全服务隔离文件方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "隔离文件成功",
    "time": "2025-09-14T20:40:50+08:00",
    "data": {}
}
```

---

### 2.11 恢复文件

**接口地址**: `POST /api/v1/security/restore`

**功能描述**: 从隔离区恢复文件

**请求参数**:

```json
{
    "file_path": "D:\\Code\\Go\\workspace\\Yara-Security\\Yara\\quarantine\\隔离.txt_1757853650_quarantined"
}
```

**参数说明**:

- `file_path` (string, 必需): 要恢复的文件路径

**响应状态码**:

- `200` - 文件恢复成功

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`RestoreFileRequest`](#type-restorefilerequest) - 恢复文件请求结构
- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`SecurityHandler.RestoreFile()`](#func-h-securityhandler-restorefile) - 安全处理器恢复文件方法
- [`SecurityService.RestoreFile()`](#func-s-securityservice-restorefile) - 安全服务恢复文件方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "恢复文件成功",
    "time": "2025-09-14T20:41:40+08:00",
    "data": {}
}
```

---

### 2.12 获取扫描历史

**接口地址**: `GET /api/v1/security/scan-history`

**功能描述**: 获取安全扫描历史记录

**请求参数**: 无

**响应状态码**:

- `200` - 扫描历史记录

**响应格式**: `application/json`

**标签**: 安全扫描

**相关类型定义**:

- [`ScanResult`](#type-scanresult) - 扫描结果结构
- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`SecurityHandler.GetScanHistory()`](#func-h-securityhandler-getscanhistory) - 安全处理器扫描历史获取方法
- [`SecurityService.GetScanHistory()`](#func-s-securityservice-getscanhistory) - 安全服务扫描历史获取方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T20:41:53+08:00",
    "data": {
        "count": 1,
        "data": [
            {
                "file_path": "C:\\Windows\\System32\\notepad.exe",
                "is_infected": false,
                "threats": null,
                "scan_time": "2025-09-14T20:41:53.7591466+08:00",
                "scan_duration": 100000000,
                "file_info": {
                    "path": "C:\\Windows\\System32\\notepad.exe",
                    "name": "notepad.exe",
                    "size": 1024,
                    "is_dir": false,
                    "mod_time": "0001-01-01T00:00:00Z",
                    "create_time": "0001-01-01T00:00:00Z",
                    "access_time": "0001-01-01T00:00:00Z",
                    "permissions": "",
                    "owner": "",
                    "group": "",
                    "md5": "md5_hash_0",
                    "sha256": "sha256_hash_0"
                }
            }
        ]
    }
}
```

---

## 📁 文件管理模块

### 3.1 获取文件列表

**接口地址**: `GET /api/v1/file/list`

**功能描述**: 获取指定目录下的文件列表

**查询参数**:

- `path` (string, 可选): 目录路径，例如：C:\temp\test
- `recursive` (boolean, 可选): 是否递归，例如：true
- `max_depth` (string, 可选): 最大深度，例如：1

**参数说明**:

- `path`: 目录路径，可以是相对路径或绝对路径
- `recursive`: 是否递归扫描子目录，true 表示递归，false 表示只扫描当前目录
- `max_depth`: 递归扫描的最大深度，当 recursive 为 true 时生效，1 表示只扫描当前目录，2 表示扫描当前目录和一级子目录

**响应状态码**:

- `200` - 文件列表

**响应格式**: `application/json`

**标签**: 文件管理

**相关类型定义**:

- [`FileListResponse`](#type-filelistresponse) - 文件列表响应结构
- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`FileHandler.GetFileList()`](#func-h-filehandler-getfilelist) - 文件处理器获取文件列表方法
- [`FileManager.GetFileList()`](#func-f-filemanager-getfilelist) - 文件管理器获取文件列表方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取文件列表成功",
    "data": {
        "count": 0,
        "directory": "C:\\temp\\test",
        "files": null,
        "max_depth": 1,
        "recursive": true
    },
    "time": "2025-09-14T20:50:56.0228095+08:00"
}
```

---

### 3.2 获取文件详细信息

**接口地址**: `GET /api/v1/file/info/{path}`

**功能描述**: 获取指定文件的详细属性信息

**路径参数**:

- `path` (string, 必需): 文件路径（支持未编码与已编码，如：test_files\test_batch.bat 或 C%3A%5Ctemp%5Ctest.txt）

**查询参数**:

- `fullPath` (string, 可选): 完整路径（当路径参数不包含斜杠时使用）

**参数说明**:

- `path`: 文件路径，可以是相对路径或绝对路径
- `fullPath`: 当路径包含斜杠时，使用此参数传递完整路径

**响应状态码**:

- `200` - 文件信息

**响应格式**: `application/json`

**标签**: 文件管理

**相关类型定义**:

- [`FileInfo`](#type-fileinfo) - 文件信息结构
- [`BasicFileInfo`](#type-basicfileinfo) - 基础文件信息结构

**相关方法签名**:

- [`FileHandler.GetFileInfo()`](#func-h-filehandler-getfileinfo) - 文件处理器获取文件信息方法
- [`Scanner.getFileInfo()`](#func-s-scanner-getfileinfo) - 扫描器获取文件信息方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T20:51:06+08:00",
    "data": {
        "path": "test_files\\test_batch.bat",
        "name": "test_batch.bat",
        "size": 904,
        "is_dir": false,
        "mod_time": "2025-08-12T10:09:21+08:00",
        "create_time": "2025-08-25T01:20:33.128063+08:00",
        "access_time": "2025-09-14T20:51:06.3147253+08:00",
        "permissions": "-rw-rw-rw-",
        "owner": "DESKTOP-ASUS\\ASUS",
        "group": "default_group",
        "md5": "38ef5e9826fb14478b3d8fee15f34e83",
        "sha256": "c7968d46e7ef5c9b0a1a6068566df3e95e28c652f7a5118912428422a69f0259",
        "is_suspicious": true,
        "threat_level": "Medium"
    }
}
```

---

### 3.3 获取文件哈希值

**接口地址**: `GET /api/v1/file/hash/{path}`

**功能描述**: 计算并返回文件的哈希值

**路径参数**:

- `path` (string, 必需): 文件路径（如果路径包含斜杠，请使用 fullPath 查询参数）

**查询参数**:

- `fullPath` (string, 可选): 完整文件路径（当路径参数不包含斜杠时使用），例如：C:\Windows\notepad.exe
- `algorithm` (string, 可选): 哈希算法，例如：sha256

**参数说明**:

- `path`: 文件路径，可以是相对路径或绝对路径
- `fullPath`: 当路径包含斜杠时，使用此参数传递完整路径
- `algorithm`: 哈希算法，支持 md5、sha1、sha256 等

**响应状态码**:

- `200` - 文件哈希值

**响应格式**: `application/json`

**标签**: 文件管理

**相关类型定义**:

- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`FileHandler.GetFileHash()`](#func-h-filehandler-getfilehash) - 文件处理器获取文件哈希方法
- [`FileManager.GetFileHash()`](#func-f-filemanager-getfilehash) - 文件管理器获取文件哈希方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取文件哈希成功",
    "time": "2025-09-14T20:51:12+08:00",
    "data": {
        "algorithm": "sha256",
        "file_path": "C:\\Windows\\notepad.exe",
        "hash": "1552f6a579b77b61460df56cb4b2ce0a34fe96b6176829d7916275b806edc2bb"
    }
}
```

---

### 3.4 验证文件哈希

**接口地址**: `POST /api/v1/file/verify-hash`

**功能描述**: 验证文件的哈希值是否匹配

**请求参数**:

```json
{
  "file_path": "C:\\Windows\\notepad.exe",
  "algorithm": "sha256",
  "expected_hash": "1552f6a579b77b61460df56cb4b2ce0a34fe96b6176829d7916275b806edc2bb"
}
```

**参数说明**:

- `file_path` (string, 必需): 文件路径
- `algorithm` (string, 必需): 哈希算法
- `expected_hash` (string, 必需): 期望的哈希值

**响应状态码**:

- `200` - 哈希验证结果

**响应格式**: `application/json`

**标签**: 文件管理

**相关类型定义**:

- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`FileHandler.VerifyHash()`](#func-h-filehandler-verifyhash) - 文件处理器验证哈希方法
- [`FileManager.VerifyHash()`](#func-f-filemanager-verifyhash) - 文件管理器验证哈希方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "哈希验证完成",
    "time": "2025-09-14T20:51:26+08:00",
    "data": {
        "actual_hash": "1552f6a579b77b61460df56cb4b2ce0a34fe96b6176829d7916275b806edc2bb",
        "algorithm": "sha256",
        "expected_hash": "1552f6a579b77b61460df56cb4b2ce0a34fe96b6176829d7916275b806edc2bb",
        "file_path": "C:\\Windows\\notepad.exe",
        "valid": true
    }
}
```

---

### 3.5 获取文件所有哈希值

**接口地址**: `GET /api/v1/file/hashes/{path}`

**功能描述**: 计算并返回文件的所有哈希值（MD5、SHA1、SHA256）

**路径参数**:

- `path` (string, 必需): 文件路径（如果路径包含斜杠，请使用 fullPath 查询参数）

**查询参数**:

- `fullPath` (string, 可选): 完整文件路径（当路径参数不包含斜杠时使用），例如：C:\windows\notepad.exe

**参数说明**:

- `path`: 文件路径，可以是相对路径或绝对路径
- `fullPath`: 当路径包含斜杠时，使用此参数传递完整路径

**响应状态码**:

- `200` - 文件哈希值列表

**响应格式**: `application/json`

**标签**: 文件管理

**相关类型定义**:

- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`FileHandler.GetFileHashes()`](#func-h-filehandler-getfilehashes) - 文件处理器获取文件所有哈希方法
- [`FileManager.GetFileHashes()`](#func-f-filemanager-getfilehashes) - 文件管理器获取文件所有哈希方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取文件所有哈希成功",
    "time": "2025-09-14T20:51:39+08:00",
    "data": {
        "file_path": "C:\\windows\\notepad.exe",
        "hashes": {
            "md5": "fbac126407b5735583dac5ea7cf519b3",
            "sha256": "1552f6a579b77b61460df56cb4b2ce0a34fe96b6176829d7916275b806edc2bb"
        }
    }
}
```

---

### 3.6 复制文件

**接口地址**: `POST /api/v1/file/copy`

**功能描述**: 复制文件到指定位置

**请求参数**:

```json
{
  "source": "test_files\\copy.txt",
  "dest": "D:\\Code\\Go\\workspace\\Yara-Security\\Yara\\test_files\\二级菜单\\copy.txt"
}
```

**参数说明**:

- `source` (string, 必需): 源文件路径
- `dest` (string, 必需): 目标文件路径

**响应状态码**:

- `200` - 文件复制成功

**响应格式**: `application/json`

**标签**: 文件管理

**相关类型定义**:

- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`FileHandler.CopyFile()`](#func-h-filehandler-copyfile) - 文件处理器复制文件方法
- [`FileManager.CopyFile()`](#func-f-filemanager-copyfile) - 文件管理器复制文件方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "文件复制成功",
    "time": "2025-09-14T20:53:51+08:00",
    "data": {}
}
```

---

### 3.7 移动文件

**接口地址**: `POST /api/v1/file/move`

**功能描述**: 移动文件到指定位置

**请求参数**:

```json
{
  "source": "test_files\\move.txt",
  "dest": "test_files\\二级菜单\\copy.txt"
}
```

**参数说明**:

- `source` (string, 必需): 源文件路径
- `dest` (string, 必需): 目标文件路径

**响应状态码**:

- `200` - 文件移动成功

**响应格式**: `application/json`

**标签**: 文件管理

**相关类型定义**:

- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`FileHandler.MoveFile()`](#func-h-filehandler-movefile) - 文件处理器移动文件方法
- [`FileManager.MoveFile()`](#func-f-filemanager-movefile) - 文件管理器移动文件方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "文件移动成功",
    "time": "2025-09-14T20:55:43+08:00",
    "data": {}
}
```

---

### 3.8 删除文件

**接口地址**: `DELETE /api/v1/file/{path}`

**功能描述**: 删除指定文件

**路径参数**:

- `path` (string, 必需): 文件路径（如果路径包含斜杠，请使用 fullPath 查询参数）

**查询参数**:

- `fullPath` (string, 可选): 完整文件路径（当路径参数不包含斜杠时使用），例如：test_files\\二级菜单\\二级文件.txt

**参数说明**:

- `path`: 文件路径，可以是相对路径或绝对路径
- `fullPath`: 当路径包含斜杠时，使用此参数传递完整路径

**响应状态码**:

- `200` - 文件删除成功

**响应格式**: `application/json`

**标签**: 文件管理

**相关类型定义**:

- [`FileInfo`](#type-fileinfo) - 文件信息结构

**相关方法签名**:

- [`FileHandler.DeleteFile()`](#func-h-filehandler-deletefile) - 文件处理器删除文件方法
- [`FileManager.DeleteFile()`](#func-f-filemanager-deletefile) - 文件管理器删除文件方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "文件删除成功",
    "time": "2025-09-14T20:57:21+08:00",
    "data": {}
}
```

---

## 🔍 进程管理模块

### 4.1 获取系统进程列表

**接口地址**: `GET /api/v1/process/list`

**功能描述**: 获取系统当前运行的进程列表

**请求参数**: 无

**响应状态码**:

- `200` - 进程列表

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ProcessListResponse`](#type-processlistresponse) - 进程列表响应结构
- [`ProcessInfo`](#type-processinfo) - 进程信息结构

**相关方法签名**:

- [`ProcessHandler.GetProcesses()`](#func-h-processhandler-getprocesses) - 进程处理器获取进程列表方法
- [`ProcessManager.GetProcessList()`](#func-p-processmanager-getprocesslist) - 进程管理器获取进程列表方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取进程列表成功",
    "time": "2025-09-14T22:10:45+08:00",
    "data": {
        "count": 376,
        "processes": [
            {
                "pid": 0,
                "name": "System Idle Process",
                "exe": "",
                "cmdline": "",
                "cwd": "",
                "status": "",
                "cpu_percent": 0,
                "memory_info": null,
                "create_time": "0001-01-01T00:00:00Z",
                "username": "",
                "ppid": 0,
                "num_threads": 0,
                "num_files": 0,
                "priority": 0,
                "modules": null,
                "children": null,
                "connections": null,
                "architecture": "",
                "memory_percent": 0,
                "working_set": 0,
                "private_bytes": 0,
                "peak_working_set": 0,
                "peak_private_bytes": 0,
                "page_faults": 0,
                "io_counters": null,
                "context_switches": 0,
                "handle_count": 0,
                "session_id": 0,
                "integrity_level": "",
                "elevated": false,
                "dep_enabled": false,
                "aslr_enabled": false,
                "last_update": "2025-09-14T22:10:45.4084415+08:00"
            }
        ]
    }
}
```

---

### 4.2 获取进程详细信息

**接口地址**: `GET /api/v1/process/{pid}`

**功能描述**: 根据 PID 获取进程的详细信息

**路径参数**:

- `pid` (integer, 必需): 进程 ID，例如：42192

**参数说明**:

- `pid`: 进程 ID，必须是系统中正在运行的进程 ID

**响应状态码**:

- `200` - 进程详细信息

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ProcessDetailResponse`](#type-processdetailresponse) - 进程详细信息响应结构
- [`ProcessInfo`](#type-processinfo) - 进程信息结构

**相关方法签名**:

- [`ProcessHandler.GetProcessByPID()`](#func-h-processhandler-getprocessbypid) - 进程处理器获取进程详细信息方法
- [`ProcessManager.GetProcessDetail()`](#func-p-processmanager-getprocessdetail) - 进程管理器获取进程详细信息方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:13:15+08:00",
    "data": {
        "code": 200,
        "message": "获取成功",
        "data": {
            "pid": 42192,
            "name": "Notepad.exe",
            "exe": "C:\\Program Files\\WindowsApps\\Microsoft.WindowsNotepad_11.2507.26.0_x64__8wekyb3d8bbwe\\Notepad\\Notepad.exe",
            "cmdline": "\"C:\\Program Files\\WindowsApps\\Microsoft.WindowsNotepad_11.2507.26.0_x64__8wekyb3d8bbwe\\Notepad\\Notepad.exe\" /SESSION:KNvfIYoqxE+GD9IiN5XlwwAAAAABFAABQIYAAIoqxE8IAwAAAAAAAA==",
            "cwd": "C:\\windows\\",
            "status": "",
            "cpu_percent": 2.0537521650043566,
            "memory_info": {
                "rss": 161570816,
                "vms": 139526144,
                "hwm": 0,
                "data": 0,
                "stack": 0,
                "locked": 0,
                "swap": 0
            },
            "create_time": "2025-09-14T22:12:39+08:00",
            "username": "DESKTOP-ASUS\\ASUS",
            "ppid": 34368,
            "num_threads": 106,
            "num_files": 0,
            "priority": 8,
            "modules": null,
            "children": null,
            "connections": null,
            "architecture": "x64",
            "memory_percent": 0.4798494577407837,
            "working_set": 161570816,
            "private_bytes": 66640,
            "peak_working_set": 175697920,
            "peak_private_bytes": 68792,
            "page_faults": 48229,
            "io_counters": {
                "read_count": 268,
                "write_count": 40,
                "read_bytes": 21679183,
                "write_bytes": 4011,
                "other_count": 4215,
                "other_bytes": 42006
            },
            "context_switches": 0,
            "handle_count": 0,
            "session_id": 0,
            "integrity_level": "Medium",
            "elevated": false,
            "dep_enabled": true,
            "aslr_enabled": true,
            "last_update": "2025-09-14T22:13:15.2348431+08:00"
        },
        "time": "0001-01-01T00:00:00Z"
    }
}
```

---

### 4.3 启动进程

**接口地址**: `POST /api/v1/process/start`

**功能描述**: 启动新的进程

**请求参数**:

```json
{
    "command": "notepad.exe",
    "args": [],
    "working_dir": "C:\\windows"
}
```

**参数说明**:

- `command` (string, 必需): 命令
- `args` (array, 可选): 命令行参数
- `working_dir` (string, 可选): 工作目录

**响应状态码**:

- `200` - 进程启动成功

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`StartProcessRequest`](#type-startprocessrequest) - 启动进程请求结构
- [`StartProcessResponse`](#type-startprocessresponse) - 启动进程响应结构

**相关方法签名**:

- [`ProcessHandler.StartProcess()`](#func-h-processhandler-startprocess) - 进程处理器启动进程方法
- [`ProcessManager.StartProcess()`](#func-p-processmanager-startprocess) - 进程管理器启动进程方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "启动进程成功",
    "time": "2025-09-14T22:12:39+08:00",
    "data": {
        "command": "notepad.exe",
        "message": "进程启动成功，正在运行中",
        "process_path": "C:\\windows\\notepad.exe",
        "start_time": "2025-09-14T22:12:39.5405091+08:00",
        "status": "started",
        "working_dir": "C:\\windows"
    }
}
```

---

### 4.4 挂起进程

**接口地址**: `PUT /api/v1/process/{pid}/suspend`

**功能描述**: 挂起指定进程

**路径参数**:

- `pid` (integer, 必需): 进程 ID，例如：42192

**参数说明**:

- `pid`: 进程 ID，必须是系统中正在运行的进程 ID

**响应状态码**:

- `200` - 进程已挂起

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`SuspendProcessResponse`](#type-suspendprocessresponse) - 挂起进程响应结构

**相关方法签名**:

- [`ProcessHandler.SuspendProcess()`](#func-h-processhandler-suspendprocess) - 进程处理器挂起进程方法
- [`ProcessManager.SuspendProcess()`](#func-p-processmanager-suspendprocess) - 进程管理器挂起进程方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "暂停进程成功",
    "time": "2025-09-14T22:13:46+08:00",
    "data": {}
}
```

---

### 4.5 恢复进程

**接口地址**: `PUT /api/v1/process/{pid}/resume`

**功能描述**: 恢复挂起的进程

**路径参数**:

- `pid` (integer, 必需): 进程 ID，例如：42192

**参数说明**:

- `pid`: 进程 ID，必须是系统中正在运行的进程 ID

**响应状态码**:

- `200` - 进程已恢复

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ProcessStatusResponse`](#type-processstatusresponse) - 进程状态响应结构

**相关方法签名**:

- [`ProcessHandler.ResumeProcess()`](#func-h-processhandler-resumeprocess) - 进程处理器恢复进程方法
- [`ProcessManager.ResumeProcess()`](#func-p-processmanager-resumeprocess) - 进程管理器恢复进程方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "恢复进程成功",
    "time": "2025-09-14T22:14:32+08:00",
    "data": {}
}
```

---

### 4.6 结束进程

**接口地址**: `DELETE /api/v1/process/{pid}`

**功能描述**: 根据 PID 结束指定进程

**路径参数**:

- `pid` (integer, 必需): 进程 ID，例如：42192

**参数说明**:

- `pid`: 进程 ID，必须是系统中正在运行的进程 ID

**响应状态码**:

- `200` - 进程已结束

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ProcessStatusResponse`](#type-processstatusresponse) - 进程状态响应结构

**相关方法签名**:

- [`ProcessHandler.KillProcess()`](#func-h-processhandler-killprocess) - 进程处理器结束进程方法
- [`ProcessManager.KillProcess()`](#func-p-processmanager-killprocess) - 进程管理器结束进程方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "结束进程成功",
    "time": "2025-09-14T22:14:53+08:00",
    "data": {}
}
```

---

### 4.7 获取进程网络连接

**接口地址**: `GET /api/v1/process/{pid}/connections`

**功能描述**: 获取指定进程的网络连接

**路径参数**:

- `pid` (integer, 必需): 进程 ID，例如：2176

**参数说明**:

- `pid`: 进程 ID，必须是系统中正在运行的进程 ID

**响应状态码**:

- `200` - 进程网络连接

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ProcessConnectionsResponse`](#type-processconnectionsresponse) - 进程网络连接响应结构
- [`ProcessConnectionInfo`](#type-processconnectioninfo) - 进程连接信息结构

**相关方法签名**:

- [`ProcessHandler.GetProcessConnections()`](#func-h-processhandler-getprocessconnections) - 进程处理器获取进程网络连接方法
- [`ProcessManager.GetProcessConnections()`](#func-p-processmanager-getprocessconnections) - 进程管理器获取进程网络连接方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取进程网络连接成功",
    "time": "2025-09-14T22:15:40+08:00",
    "data": {
        "connections": [
            {
                "fd": 0,
                "family": 2,
                "type": 1,
                "laddr": {
                    "ip": "127.0.0.1",
                    "port": 51681
                },
                "raddr": {
                    "ip": "127.0.0.1",
                    "port": 7890
                },
                "status": "ESTABLISHED"
            }
        ]
    }
}
```

---

### 4.8 获取进程内存信息

**接口地址**: `GET /api/v1/process/{pid}/memory`

**功能描述**: 获取指定进程的内存使用情况

**路径参数**:

- `pid` (integer, 必需): 进程 ID，例如：2176

**参数说明**:

- `pid`: 进程 ID，必须是系统中正在运行的进程 ID

**响应状态码**:

- `200` - 进程内存信息

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ProcessMemoryResponse`](#type-processmemoryresponse) - 进程内存信息响应结构
- [`MemoryInfo`](#type-memoryinfo) - 内存信息结构

**相关方法签名**:

- [`ProcessHandler.GetProcessMemory()`](#func-h-processhandler-getprocessmemory) - 进程处理器获取进程内存信息方法
- [`ProcessManager.GetProcessMemory()`](#func-p-processmanager-getprocessmemory) - 进程管理器获取进程内存信息方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:16:39+08:00",
    "data": {
        "code": 200,
        "message": "获取成功",
        "data": {
            "pid": 0,
            "rss": 32927744,
            "vms": 28147712,
            "percent": 0.09779216349124908,
            "available": 0,
            "used": 0,
            "free": 0,
            "total": 0
        },
        "time": "0001-01-01T00:00:00Z"
    }
}
```

---

### 4.9 检查进程运行状态

**接口地址**: `GET /api/v1/process/{pid}/running`

**功能描述**: 检查指定进程是否正在运行

**路径参数**:

- `pid` (integer, 必需): 进程 ID，例如：2176

**参数说明**:

- `pid`: 进程 ID，必须是系统中正在运行的进程 ID

**响应状态码**:

- `200` - 进程运行状态

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ProcessRunningResponse`](#type-processrunningresponse) - 进程运行状态响应结构

**相关方法签名**:

- [`ProcessHandler.CheckProcessRunning()`](#func-h-processhandler-checkprocessrunning) - 进程处理器检查进程运行状态方法
- [`ProcessManager.CheckProcessRunning()`](#func-p-processmanager-checkprocessrunning) - 进程管理器检查进程运行状态方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取状态成功",
    "time": "2025-09-14T22:17:00+08:00",
    "data": {
        "is_running": {
            "code": 200,
            "message": "获取成功",
            "data": {
                "running": true
            },
            "time": "0001-01-01T00:00:00Z"
        }
    }
}
```

---

### 4.10 获取进程统计信息

**接口地址**: `GET /api/v1/process/statistics`

**功能描述**: 获取进程相关的统计信息

**请求参数**: 无

**响应状态码**:

- `200` - 进程统计信息

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ProcessStatisticsResponse`](#type-processstatisticsresponse) - 进程统计信息响应结构

**相关方法签名**:

- [`ProcessHandler.GetProcessStatistics()`](#func-h-processhandler-getprocessstatistics) - 进程处理器获取进程统计信息方法
- [`ProcessManager.GetProcessStatistics()`](#func-p-processmanager-getprocessstatistics) - 进程管理器获取进程统计信息方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:17:11+08:00",
    "data": {
        "code": 200,
        "message": "获取成功",
        "data": {
            "running_processes": 0,
            "suspended_processes": 0,
            "total_memory_usage": 13261746176,
            "total_processes": 374
        },
        "time": "0001-01-01T00:00:00Z"
    }
}
```

---

### 4.11 获取进程子进程

**接口地址**: `GET /api/v1/process/{pid}/children`

**功能描述**: 获取指定进程的子进程列表

**路径参数**:

- `pid` (integer, 必需): 进程 ID，例如：40044

**参数说明**:

- `pid`: 进程 ID，必须是系统中正在运行的进程 ID

**响应状态码**:

- `200` - 子进程列表

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ProcessChildrenResponse`](#type-processchildrenresponse) - 进程子进程响应结构
- [`ChildProcess`](#type-childprocess) - 子进程结构

**相关方法签名**:

- [`ProcessHandler.GetProcessChildren()`](#func-h-processhandler-getprocesschildren) - 进程处理器获取进程子进程方法
- [`ProcessManager.GetProcessChildren()`](#func-p-processmanager-getprocesschildren) - 进程管理器获取进程子进程方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:17:38+08:00",
    "data": {
        "children": {
            "code": 200,
            "message": "获取成功",
            "data": {
                "children": [
                    2176,
                    41668,
                    18588,
                    24188,
                    34572,
                    40772,
                    23572,
                    18008,
                    10792,
                    46336
                ]
            },
            "time": "0001-01-01T00:00:00Z"
        },
        "pid": 40044
    }
}
```

---

### 4.12 获取进程模块列表

**接口地址**: `GET /api/v1/process/{pid}/modules`

**功能描述**: 获取指定进程加载的模块列表

**路径参数**:

- `pid` (integer, 必需): 进程 ID，例如：2176

**参数说明**:

- `pid`: 进程 ID，必须是系统中正在运行的进程 ID

**响应状态码**:

- `200` - 进程模块列表

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ProcessModulesResponse`](#type-processmodulesresponse) - 进程模块列表响应结构
- [`ProcessModule`](#type-processmodule) - 进程模块结构

**相关方法签名**:

- [`ProcessHandler.GetProcessModules()`](#func-h-processhandler-getprocessmodules) - 进程处理器获取进程模块方法
- [`ProcessManager.GetProcessModules()`](#func-p-processmanager-getprocessmodules) - 进程管理器获取进程模块方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取进程模块相关进程信息成功",
    "time": "2025-09-14T22:19:21+08:00",
    "data": {
        "description": "与进程PID=2176共享模块的所有相关进程信息",
        "process_count": 11,
        "related_processes": {
            "description": "与进程PID=2176共享模块的所有相关进程信息",
            "process_count": 11,
            "related_processes": [
                {
                    "pid": 2176,
                    "name": "WeChatAppEx.exe",
                    "exe": "C:\\Users\\ASUS\\AppData\\Roaming\\Tencent\\xwechat\\XPlugin\\Plugins\\RadiumWMPF\\16467\\extracted\\runtime\\WeChatAppEx.exe",
                    "cmdline": "\"C:\\Users\\ASUS\\AppData\\Roaming\\Tencent\\xwechat\\xplugin\\plugins\\RadiumWMPF\\16467\\extracted\\runtime\\WeChatAppEx.exe\" --type=utility --utility-sub-type=network.mojom.NetworkService --lang=zh-CN --service-sandbox-type=none --client_version=4065595426 --enable-crash-reporter --wmpf_root_dir=\"C:\\Users\\ASUS\\AppData\\Roaming\\Tencent\\xwechat\\radium\" --product-id=1002 --disable-mojo-broker --field-trial-handle=3068,i,14485445038403588694,6596298628817470105,262144 --enable-features=OverlayScrollbar,XWorker --disable-features=AudioServiceOutOfProcess,AutoupgradeMixedContent,BackForwardCache,DigitalGoodsApi,NotificationTriggers,PeriodicBackgroundSync,TFLiteLanguageDetectionEnabled,Vulkan,WebOTP --variations-seed-version --log-level=2 --mojo-platform-channel-handle=3048 /prefetch:11",
                    "cwd": "C:\\Users\\ASUS\\AppData\\Roaming\\Tencent\\xwechat\\xplugin\\plugins\\RadiumWMPF\\16467\\extracted\\runtime\\",
                    "status": "",
                    "cpu_percent": 0.00747268865499527,
                    "memory_info": {
                        "rss": 32927744,
                        "vms": 28147712,
                        "hwm": 0,
                        "data": 0,
                        "stack": 0,
                        "locked": 0,
                        "swap": 0
                    },
                    "create_time": "2025-09-14T10:11:00+08:00",
                    "username": "DESKTOP-ASUS\\ASUS",
                    "ppid": 40044,
                    "num_threads": 13,
                    "num_files": 0,
                    "priority": 8,
                    "modules": null,
                    "children": null,
                    "connections": null,
                    "architecture": "x64",
                    "memory_percent": 0.09779216349124908,
                    "working_set": 32927744,
                    "private_bytes": 29840,
                    "peak_working_set": 39460864,
                    "peak_private_bytes": 187536,
                    "page_faults": 29550,
                    "io_counters": {
                        "read_count": 2088,
                        "write_count": 3533,
                        "read_bytes": 1364180,
                        "write_bytes": 803447,
                        "other_count": 22797,
                        "other_bytes": 463400
                    },
                    "context_switches": 0,
                    "handle_count": 0,
                    "session_id": 0,
                    "integrity_level": "Medium",
                    "elevated": false,
                    "dep_enabled": true,
                    "aslr_enabled": true,
                    "last_update": "2025-09-14T22:19:21.2696169+08:00"
                }
            ],
            "target_pid": 2176,
            "timestamp": "2025-09-14T22:19:21.4463533+08:00"
        },
        "target_pid": 2176,
        "timestamp": "2025-09-14T22:19:21.4463533+08:00"
    }
}
```

---

### 4.13 获取系统模块列表

**接口地址**: `GET /api/v1/process/module/list`

**功能描述**: 获取系统加载的所有模块

**请求参数**: 无

**响应状态码**:

- `200` - 系统模块列表

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`SystemModulesResponse`](#type-systemmodulesresponse) - 系统模块列表响应结构
- [`SystemModule`](#type-systemmodule) - 系统模块结构

**相关方法签名**:

- [`ProcessHandler.GetSystemModules()`](#func-h-processhandler-getsystemmodules) - 进程处理器获取系统模块列表方法
- [`ProcessManager.GetSystemModules()`](#func-p-processmanager-getsystemmodules) - 进程管理器获取系统模块列表方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取系统模块成功",
    "time": "2025-09-14T22:21:18+08:00",
    "data": {
        "count": 78,
        "data": [
            {
                "name": "process_ctfmon.exe",
                "process_count": 1,
                "total_memory": 0,
                "status": "",
                "last_accessed": "0001-01-01T00:00:00Z",
                "affected_processes": [
                    22144
                ],
                "path": "process_path_ctfmon.exe"
            }
        ]
    }
}
```

---

### 4.14 获取模块详细信息

**接口地址**: `GET /api/v1/process/module/{module}`

**功能描述**: 获取指定模块的详细信息

**路径参数**:

- `module` (string, 必需): 模块名称（支持未编码与已编码，如：kernel32.dll 或 kernel32%2Edll），例如：process_Apifox.exe

**参数说明**:

- `module`: 模块名称，必须是系统中正在运行的模块名称

**响应状态码**:

- `200` - 模块信息

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ModuleDetailResponse`](#type-moduledetailresponse) - 模块详细信息响应结构
- [`ModuleDetail`](#type-moduledetail) - 模块详细信息结构

**相关方法签名**:

- [`ProcessHandler.GetModuleDetail()`](#func-h-processhandler-getmoduledetail) - 进程处理器获取模块详细信息方法
- [`ProcessManager.GetModuleDetail()`](#func-p-processmanager-getmoduledetail) - 进程管理器获取模块详细信息方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:22:26+08:00",
    "data": {
        "code": 200,
        "message": "获取信息成功",
        "data": {
            "name": "process_Apifox.exe",
            "process_count": 8,
            "total_memory": 0,
            "status": "",
            "last_accessed": "0001-01-01T00:00:00Z",
            "affected_processes": [
                11248,
                13516,
                18516,
                20580,
                31632,
                33672,
                39160,
                39632
            ],
            "path": "process_path_Apifox.exe"
        },
        "time": "2025-09-14T22:22:26.5885618+08:00"
    }
}
```

---

### 4.15 挂起系统模块

**接口地址**: `PUT /api/v1/process/module/{module}/suspend`

**功能描述**: 挂起指定模块

**路径参数**:

- `module` (string, 必需): 模块名称（支持未编码与已编码，如：kernel32.dll 或 kernel32%2Edll），例如：process_Notepad.exe

**参数说明**:

- `module`: 模块名称，必须是系统中正在运行的模块名称

**响应状态码**:

- `200` - 模块已挂起

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`SuspendModuleResponse`](#type-suspendmoduleresponse) - 挂起模块响应结构

**相关方法签名**:

- [`ProcessHandler.SuspendModule()`](#func-h-processhandler-suspendmodule) - 进程处理器挂起模块方法
- [`ProcessManager.SuspendModule()`](#func-p-processmanager-suspendmodule) - 进程管理器挂起模块方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "挂起系统模块成功",
    "time": "2025-09-14T22:23:18+08:00",
    "data": {}
}
```

---

### 4.16 恢复系统模块

**接口地址**: `PUT /api/v1/process/module/{module}/resume`

**功能描述**: 恢复挂起的模块

**路径参数**:

- `module` (string, 必需): 模块名称（支持未编码与已编码，如：kernel32.dll 或 kernel32%2Edll），例如：process_Notepad.exe

**参数说明**:

- `module`: 模块名称，必须是系统中正在运行的模块名称

**响应状态码**:

- `200` - 模块已恢复

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`ResumeModuleResponse`](#type-resumemoduleresponse) - 恢复模块响应结构

**相关方法签名**:

- [`ProcessHandler.ResumeModule()`](#func-h-processhandler-resumemodule) - 进程处理器恢复模块方法
- [`ProcessManager.ResumeModule()`](#func-p-processmanager-resumemodule) - 进程管理器恢复模块方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "恢复系统模块成功",
    "time": "2025-09-14T22:23:55+08:00",
    "data": {}
}
```

---

### 4.17 结束系统模块

**接口地址**: `DELETE /api/v1/process/module/{module}`

**功能描述**: 结束指定模块

**路径参数**:

- `module` (string, 必需): 模块名称（支持未编码与已编码，如：kernel32.dll 或 kernel32%2Edll），例如：process_Notepad.exe

**参数说明**:

- `module`: 模块名称，必须是系统中正在运行的模块名称

**响应状态码**:

- `200` - 模块已结束

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`KillModuleResponse`](#type-killmoduleresponse) - 结束模块响应结构

**相关方法签名**:

- [`ProcessHandler.KillModule()`](#func-h-processhandler-killmodule) - 进程处理器结束模块方法
- [`ProcessManager.KillModule()`](#func-p-processmanager-killmodule) - 进程管理器结束模块方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "结束系统模块成功",
    "time": "2025-09-14T22:24:11+08:00",
    "data": {}
}
```

---

### 4.18 启用进程监控

**接口地址**: `POST /api/v1/process/monitoring/enable`

**功能描述**: 启用进程监控功能

**请求参数**: 无

**响应状态码**:

- `200` - 进程监控已启用

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`EnableMonitoringResponse`](#type-enablemonitoringresponse) - 启用监控响应结构

**相关方法签名**:

- [`ProcessHandler.EnableMonitoring()`](#func-h-processhandler-enablemonitoring) - 进程处理器启用监控方法
- [`ProcessManager.EnableMonitoring()`](#func-p-processmanager-enablemonitoring) - 进程管理器启用监控方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "进程监控已启用成功",
    "time": "2025-09-14T22:24:22+08:00",
    "data": {}
}
```

---

### 4.19 禁用进程监控

**接口地址**: `POST /api/v1/process/monitoring/disable`

**功能描述**: 禁用进程监控功能

**请求参数**: 无

**响应状态码**:

- `200` - 进程监控已禁用

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`DisableMonitoringResponse`](#type-disablemonitoringresponse) - 禁用监控响应结构

**相关方法签名**:

- [`ProcessHandler.DisableMonitoring()`](#func-h-processhandler-disablemonitoring) - 进程处理器禁用监控方法
- [`ProcessManager.DisableMonitoring()`](#func-p-processmanager-disablemonitoring) - 进程管理器禁用监控方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "进程监控已禁用成功",
    "time": "2025-09-14T22:24:32+08:00",
    "data": {}
}
```

---

### 4.20 获取监控进程列表

**接口地址**: `GET /api/v1/process/monitoring/list`

**功能描述**: 获取当前监控的进程列表

**请求参数**: 无

**响应状态码**:

- `200` - 监控进程列表

**响应格式**: `application/json`

**标签**: 进程管理

**相关类型定义**:

- [`MonitoringListResponse`](#type-monitoringlistresponse) - 监控列表响应结构
- [`MonitoredProcess`](#type-monitoredprocess) - 监控进程结构

**相关方法签名**:

- [`ProcessHandler.GetMonitoringList()`](#func-h-processhandler-getmonitoringlist) - 进程处理器获取监控列表方法
- [`ProcessManager.GetMonitoringList()`](#func-p-processmanager-getmonitoringlist) - 进程管理器获取监控列表方法

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取监控进程列表成功",
    "time": "2025-09-14T22:24:39+08:00",
    "data": {
        "count": 373,
        "monitored_processes": [
            {
                "pid": 4,
                "name": "System",
                "exe": "",
                "cmdline": "",
                "cwd": "",
                "status": "",
                "cpu_percent": 0,
                "memory_info": null,
                "create_time": "0001-01-01T00:00:00Z",
                "username": "",
                "ppid": 0,
                "num_threads": 0,
                "num_files": 0,
                "priority": 0,
                "modules": null,
                "children": null,
                "connections": null,
                "architecture": "",
                "memory_percent": 0,
                "working_set": 0,
                "private_bytes": 0,
                "peak_working_set": 0,
                "peak_private_bytes": 0,
                "page_faults": 0,
                "io_counters": null,
                "context_switches": 0,
                "handle_count": 0,
                "session_id": 0,
                "integrity_level": "",
                "elevated": false,
                "dep_enabled": false,
                "aslr_enabled": false,
                "last_update": "2025-09-14T22:24:39.603699+08:00"
            }
        ],
        "monitoring_enabled": true
    }
}
```

---

## 🔧 注册表管理模块

### 5.1 获取注册表键信息

**接口地址**: `GET /api/v1/registry/key/{path}`

**功能描述**: 获取指定注册表路径的键信息

**路径参数**:

- `path` (string, 必需): 注册表路径（如果路径包含斜杠，请使用 fullPath 查询参数）

**查询参数**:

- `fullPath` (string, 可选): 完整注册表路径（当路径参数不包含斜杠时使用），例如：SOFTWARE\TestKey3

**响应状态码**:

- `200` - 注册表键信息

**响应格式**: `application/json`

**标签**: 注册表管理

**相关类型定义**: [`RegistryKeyResponse`](#type-registrykeyresponse), [`RegistryKeyInfo`](#type-registrykeyinfo)

**相关方法签名**: [`RegistryHandler.GetRegistryKey()`](#func-registryhandlergetregistrykey), [`RegistryManager.GetRegistryKey()`](#func-registrymanagergetregistrykey)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:34:55+08:00",
    "data": {
        "path": "SOFTWARE\\TestKey3",
        "name": "TestKey3",
        "type": "REG_KEY",
        "value": null,
        "values": {
            "TestString3": "TestValue3"
        },
        "last_modified": "2025-09-14T22:34:55.3951135+08:00"
    }
}
```

---

### 5.2 列出注册表键

**接口地址**: `GET /api/v1/registry/keys/{path}`

**功能描述**: 列出指定注册表路径下的子键

**路径参数**:

- `path` (string, 必需): 注册表路径（如果路径包含斜杠，请使用 fullPath 查询参数），例如：CurrentVersion

**查询参数**:

- `fullPath` (string, 可选): 完整注册表路径（当路径参数不包含斜杠时使用），例如：SOFTWARE

**响应状态码**:

- `200` - 注册表键列表

**响应格式**: `application/json`

**标签**: 注册表管理

**相关类型定义**: [`RegistryKeysResponse`](#type-registrykeysresponse), [`RegistrySubKey`](#type-registrysubkey)

**相关方法签名**: [`RegistryHandler.ListRegistryKeys()`](#func-registryhandlerlistregistrykeys), [`RegistryManager.ListRegistryKeys()`](#func-registrymanagerlistregistrykeys)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:35:29+08:00",
    "data": {
        "count": 63,
        "data": [
            "a47a88c0-27d8-5a45-af46-ded84899f5d9",
            "appdatalow",
            "Baidu",
            "Bandizip",
            "BoosterX",
            "Bytedance",
            "ByteRtcSDK",
            "ChangeTracker",
            "Chromium",
            "Clash Verge Rev",
            "Clients",
            "DeviceInfo",
            "e38e520f-3e39-5799-a112-30abd8037c63",
            "ej-technologies",
            "fastpdf",
            "Feishu",
            "finalshell",
            "Geek Uninstaller",
            "Google",
            "GoProgrammingLanguage",
            "huandian.com",
            "Huorong",
            "JavaSoft",
            "kdiskmgr_sogou",
            "Kingsoft",
            "Kingsoft DataRecovery Master",
            "KSP",
            "kzip_sogou",
            "Microsoft",
            "Mobatek",
            "MozillaPlugins",
            "MySQL",
            "NetSarang",
            "Netscape",
            "NVIDIA Corporation",
            "ODBC",
            "Policies",
            "PremiumSoft",
            "QLV",
            "QtProject",
            "Realtek",
            "redisdesktop",
            "RegisteredApplications",
            "roamingdevice",
            "SogouInput",
            "SogouInput.ppup",
            "SogouInput.store.user",
            "SogouInput.tc",
            "sogoupdf",
            "sogoupdfsdk",
            "SyncEngines",
            "SystemTask",
            "Tencent",
            "TestKey3",
            "Typora",
            "UTForPC",
            "VMware, Inc.",
            "WeChatAppEx",
            "WinRAR",
            "Wow6432Node",
            "Classes",
            "SogouInput.ppup.user",
            "SogouInput.user"
        ]
    }
}
```

---

### 5.3 列出注册表值

**接口地址**: `GET /api/v1/registry/values/{path}`

**功能描述**: 列出指定注册表路径下的值

**路径参数**:

- `path` (string, 必需): 注册表路径（如果路径包含斜杠，请使用 fullPath 查询参数），例如：CurrentVersion

**查询参数**:

- `fullPath` (string, 可选): 完整注册表路径（当路径参数不包含斜杠时使用），例如：SOFTWARE\\TestKey3

**响应状态码**:

- `200` - 注册表值列表

**响应格式**: `application/json`

**标签**: 注册表管理

**相关类型定义**: [`RegistryValuesResponse`](#type-registryvaluesresponse), [`RegistryValue`](#type-registryvalue)

**相关方法签名**: [`RegistryHandler.ListRegistryValues()`](#func-registryhandlerlistregistryvalues), [`RegistryManager.ListRegistryValues()`](#func-registrymanagerlistregistryvalues)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:36:04+08:00",
    "data": {
        "count": 1,
        "data": [
            "TestString3"
        ]
    }
}
```

---

### 5.4 创建注册表键

**接口地址**: `POST /api/v1/registry/key`

**功能描述**: 创建新的注册表键

**请求体参数**:

- `path` (string, 必需): 注册表路径，例如：SOFTWARE\\TestKey2

**响应状态码**:

- `200` - 注册表键创建成功

**响应格式**: `application/json`

**标签**: 注册表管理

**相关类型定义**: [`CreateRegistryKeyRequest`](#type-createregistrykeyrequest), [`CreateRegistryKeyResponse`](#type-createregistrykeyresponse)

**相关方法签名**: [`RegistryHandler.CreateRegistryKey()`](#func-registryhandlercreateregistrykey), [`RegistryManager.CreateRegistryKey()`](#func-registrymanagercreateregistrykey)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "创建注册表键成功",
    "time": "2025-09-14T22:36:57+08:00",
    "data": {}
}
```

---

### 5.5 删除注册表键

**接口地址**: `DELETE /api/v1/registry/key/{path}`

**功能描述**: 删除指定的注册表键

**路径参数**:

- `path` (string, 必需): 注册表路径（如果路径包含斜杠，请使用 fullPath 查询参数），例如：TestKey

**查询参数**:

- `fullPath` (string, 可选): 完整注册表路径（当路径参数不包含斜杠时使用），例如：SOFTWARE\\TestKey2

**响应状态码**:

- `200` - 注册表键删除成功

**响应格式**: `application/json`

**标签**: 注册表管理

**相关类型定义**: [`DeleteRegistryKeyResponse`](#type-deleteregistrykeyresponse)

**相关方法签名**: [`RegistryHandler.DeleteRegistryKey()`](#func-registryhandlerdeleteregistrykey), [`RegistryManager.DeleteRegistryKey()`](#func-registrymanagerdeleteregistrykey)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "删除注册表键成功",
    "time": "2025-09-14T22:40:18+08:00",
    "data": {}
}
```

---

### 5.6 获取注册表值

**接口地址**: `GET /api/v1/registry/value/{path}`

**功能描述**: 获取指定注册表路径和名称的值（名称通过查询参数 name 传入）

**路径参数**:

- `path` (string, 必需): 注册表路径（如果路径包含斜杠，请使用 fullPath 查询参数）

**查询参数**:

- `fullPath` (string, 可选): 完整注册表路径（当路径参数不包含斜杠时使用），例如：SOFTWARE\TestKey3
- `name` (string, 必需): 值名称，例如：TestString3

**响应状态码**:

- `200` - 注册表值

**响应格式**: `application/json`

**标签**: 注册表管理

**相关类型定义**: [`RegistryValueResponse`](#type-registryvalueresponse), [`RegistryValueDetail`](#type-registryvaluedetail)

**相关方法签名**: [`RegistryHandler.GetRegistryValue()`](#func-registryhandlergetregistryvalue), [`RegistryManager.GetRegistryValue()`](#func-registrymanagergetregistryvalue)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:37:41+08:00",
    "data": {
        "type": "REG_SZ",
        "value": "TestValue3"
    }
}
```

---

### 5.7 设置注册表值

**接口地址**: `PUT /api/v1/registry/value`

**功能描述**: 设置注册表键的值

**请求体参数**:

- `path` (string, 必需): 注册表路径，例如：SOFTWARE\\TestKey2
- `name` (string, 必需): 值名称，例如：TestString2
- `value` (string, 必需): 值内容，例如：TestValue2
- `type` (string, 必需): 值类型，例如：REG_SZ

**响应状态码**:

- `200` - 注册表值设置成功

**响应格式**: `application/json`

**标签**: 注册表管理

**相关类型定义**: [`SetRegistryValueRequest`](#type-setregistryvaluerequest), [`SetRegistryValueResponse`](#type-setregistryvalueresponse)

**相关方法签名**: [`RegistryHandler.SetRegistryValue()`](#func-registryhandlersetregistryvalue), [`RegistryManager.SetRegistryValue()`](#func-registrymanagersetregistryvalue)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "设置注册表值成功",
    "time": "2025-09-14T22:38:10+08:00",
    "data": {}
}
```

---

### 5.8 删除注册表值

**接口地址**: `DELETE /api/v1/registry/value/{path}`

**功能描述**: 删除指定的注册表值（名称通过查询参数 name 传入）

**路径参数**:

- `path` (string, 必需): 注册表路径（如果路径包含斜杠，请使用 fullPath 查询参数），例如：TestKey

**查询参数**:

- `fullPath` (string, 可选): 完整注册表路径（当路径参数不包含斜杠时使用），例如：SOFTWARE\\TestKey2
- `name` (string, 必需): 值名称，例如：TestString2

**响应状态码**:

- `200` - 注册表值删除成功

**响应格式**: `application/json`

**标签**: 注册表管理

**相关类型定义**: [`DeleteRegistryValueResponse`](#type-deleteregistryvalueresponse)

**相关方法签名**: [`RegistryHandler.DeleteRegistryValue()`](#func-registryhandlerdeleteregistryvalue), [`RegistryManager.DeleteRegistryValue()`](#func-registrymanagerdeleteregistryvalue)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "删除注册表值成功",
    "time": "2025-09-14T22:39:50+08:00",
    "data": {}
}
```

---

### 5.9 搜索注册表

**接口地址**: `GET /api/v1/registry/search`

**功能描述**: 在注册表中搜索指定的内容

**查询参数**:

- `root_path` (string, 可选): 根路径，例如：SOFTWARE
- `search_term` (string, 可选): 搜索关键词，例如：TestKey3

**响应状态码**:

- `200` - 搜索结果

**响应格式**: `application/json`

**标签**: 注册表管理

**相关类型定义**: [`RegistrySearchResponse`](#type-registrysearchresponse), [`RegistrySearchResult`](#type-registrysearchresult)

**相关方法签名**: [`RegistryHandler.SearchRegistry()`](#func-registryhandlersearchregistry), [`RegistryManager.SearchRegistry()`](#func-registrymanagersearchregistry)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:39:18+08:00",
    "data": {
        "count": 1,
        "data": [
            {
                "path": "SOFTWARE\\TestKey3",
                "name": "TestKey3",
                "type": "REG_KEY",
                "value": null,
                "values": {
                    "TestString3": "TestValue3"
                },
                "last_modified": "2025-09-14T22:39:16.7057368+08:00"
            }
        ]
    }
}
```

---

## 🌐 网络管理模块

### 6.1 获取所有网络连接

**接口地址**: `GET /api/v1/network/connections`

**功能描述**: 获取系统当前的所有网络连接

**请求参数**: 无

**响应状态码**:

- `200` - 网络连接列表

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`NetworkConnectionsResponse`](#type-networkconnectionsresponse), [`NetworkConnection`](#type-networkconnection)

**相关方法签名**: [`NetworkHandler.GetAllConnections()`](#func-networkhandlergetallconnections), [`NetworkManager.GetAllConnections()`](#func-networkmanagergetallconnections)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:40:42+08:00",
    "data": {
        "count": 209,
        "data": [
            {
                "id": "4-tcp-0",
                "local_addr": "0.0.0.0",
                "remote_addr": "0.0.0.0",
                "local_port": 445,
                "remote_port": 0,
                "protocol": "TCP",
                "status": "LISTEN",
                "pid": 4,
                "process_name": "System",
                "type": "tcp"
            }
        ]
    }
}
```

---

### 6.2 获取 TCP 连接

**接口地址**: `GET /api/v1/network/connections/tcp`

**功能描述**: 获取系统当前的 TCP 连接

**请求参数**: 无

**响应状态码**:

- `200` - TCP 连接列表

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`TCPConnectionsResponse`](#type-tcpconnectionsresponse), [`NetworkConnection`](#type-networkconnection)

**相关方法签名**: [`NetworkHandler.GetTCPConnections()`](#func-networkhandlergettcpconnections), [`NetworkManager.GetTCPConnections()`](#func-networkmanagergettcpconnections)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:41:12+08:00",
    "data": {
        "count": 127,
        "data": [
            {
                "id": "19648-tcp-0",
                "local_addr": "10.44.1.102",
                "remote_addr": "103.37.45.130",
                "local_port": 60841,
                "remote_port": 14784,
                "protocol": "TCP",
                "status": "ESTABLISHED",
                "pid": 19648,
                "process_name": "FlClashCore.exe",
                "type": "tcp"
            }
        ]
    }
}
```

---

### 6.3 获取 UDP 连接

**接口地址**: `GET /api/v1/network/connections/udp`

**功能描述**: 获取系统当前的 UDP 连接

**请求参数**: 无

**响应状态码**:

- `200` - UDP 连接列表

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`UDPConnectionsResponse`](#type-udpconnectionsresponse), [`NetworkConnection`](#type-networkconnection)

**相关方法签名**: [`NetworkHandler.GetUDPConnections()`](#func-networkhandlergetudpconnections), [`NetworkManager.GetUDPConnections()`](#func-networkmanagergetudpconnections)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:41:43+08:00",
    "data": {
        "count": 71,
        "data": [
            {
                "id": "4048-udp-0",
                "local_addr": "0.0.0.0",
                "remote_addr": "",
                "local_port": 4500,
                "remote_port": 0,
                "protocol": "UDP",
                "status": "",
                "pid": 4048,
                "process_name": "unknown",
                "type": "udp"
            }
        ]
    }
}
```

---

### 6.4 按进程 ID 获取连接

**接口地址**: `GET /api/v1/network/connections/pid/{pid}`

**功能描述**: 获取指定进程的网络连接

**路径参数**:

- `pid` (integer, 必需): 进程 ID，例如：2176

**参数说明**:

- `pid`: 进程 ID，必须是系统中正在运行的进程 ID

**响应状态码**:

- `200` - 进程网络连接

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`ProcessConnectionsResponse`](#type-processconnectionsresponse), [`NetworkConnection`](#type-networkconnection)

**相关方法签名**: [`NetworkHandler.GetConnectionsByPID()`](#func-networkhandlergetconnectionsbypid), [`NetworkManager.GetConnectionsByPID()`](#func-networkmanagergetconnectionsbypid)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:42:22+08:00",
    "data": {
        "count": 1,
        "data": [
            {
                "id": "2176-tcp-0",
                "local_addr": "127.0.0.1",
                "remote_addr": "127.0.0.1",
                "local_port": 51681,
                "remote_port": 7890,
                "protocol": "TCP",
                "status": "ESTABLISHED",
                "pid": 2176,
                "process_name": "WeChatAppEx.exe",
                "type": "tcp"
            }
        ]
    }
}
```

---

### 6.5 按端口获取连接

**接口地址**: `GET /api/v1/network/connections/port/{port}`

**功能描述**: 获取指定端口的网络连接

**路径参数**:

- `port` (integer, 必需): 端口号，例如：8081

**参数说明**:

- `port`: 端口号，必须是系统中正在使用的端口号

**响应状态码**:

- `200` - 端口网络连接

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`PortConnectionsResponse`](#type-portconnectionsresponse), [`NetworkConnection`](#type-networkconnection)

**相关方法签名**: [`NetworkHandler.GetConnectionsByPort()`](#func-networkhandlergetconnectionsbyport), [`NetworkManager.GetConnectionsByPort()`](#func-networkmanagergetconnectionsbyport)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:42:45+08:00",
    "data": {
        "count": 5,
        "data": [
            {
                "id": "25444-tcp-0",
                "local_addr": "127.0.0.1",
                "remote_addr": "0.0.0.0",
                "local_port": 8081,
                "remote_port": 0,
                "protocol": "TCP",
                "status": "LISTEN",
                "pid": 25444,
                "process_name": "main.exe",
                "type": "tcp"
            }
        ]
    }
}
```

---

### 6.6 按 IP 地址获取连接

**接口地址**: `GET /api/v1/network/connections/ip/{ip}`

**功能描述**: 获取指定 IP 地址的网络连接

**路径参数**:

- `ip` (string, 必需): IP 地址，例如：127.0.0.1

**参数说明**:

- `ip`: IP 地址，必须是系统中正在使用的 IP 地址

**响应状态码**:

- `200` - IP 网络连接

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`IPConnectionsResponse`](#type-ipconnectionsresponse), [`NetworkConnection`](#type-networkconnection)

**相关方法签名**: [`NetworkHandler.GetConnectionsByIP()`](#func-networkhandlergetconnectionsbyip), [`NetworkManager.GetConnectionsByIP()`](#func-networkmanagergetconnectionsbyip)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:43:17+08:00",
    "data": {
        "count": 58,
        "data": [
            {
                "id": "9548-tcp-0",
                "local_addr": "127.0.0.1",
                "remote_addr": "127.0.0.1",
                "local_port": 4369,
                "remote_port": 49670,
                "protocol": "TCP",
                "status": "ESTABLISHED",
                "pid": 9548,
                "process_name": "unknown",
                "type": "tcp"
            }
        ]
    }
}
```

---

### 6.7 根据连接 ID 获取连接信息

**接口地址**: `GET /api/v1/network/connection/{connection_id}`

**功能描述**: 根据连接 ID 获取指定网络连接的详细信息

**路径参数**:

- `connection_id` (string, 必需): 连接 ID，格式：PID-type-FD，如：1234-tcp-0

**参数说明**:

- `connection_id`: 连接 ID，必须是系统中正在使用的连接 ID

**响应状态码**:

- `200` - 连接信息获取成功
- `404` - 连接不存在

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`ConnectionDetailResponse`](#type-connectiondetailresponse), [`ConnectionDetail`](#type-connectiondetail)

**相关方法签名**: [`NetworkHandler.GetConnectionDetail()`](#func-networkhandlergetconnectiondetail), [`NetworkManager.GetConnectionDetail()`](#func-networkmanagergetconnectiondetail)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:43:52+08:00",
    "data": {
        "count": 1,
        "data": [
            {
                "id": "2176-tcp-0",
                "local_addr": "127.0.0.1",
                "remote_addr": "127.0.0.1",
                "local_port": 51681,
                "remote_port": 7890,
                "protocol": "TCP",
                "status": "ESTABLISHED",
                "pid": 2176,
                "process_name": "WeChatAppEx.exe",
                "type": "tcp"
            }
        ]
    }
}
```

---

### 6.8 关闭网络连接

**接口地址**: `POST /api/v1/network/connection/close`

**功能描述**: 根据连接 ID 或进程名称关闭网络连接

**请求体参数**:

- `connection_id` (string, 可选): 连接 ID，格式：PID-type-FD，如：17056-tcp-0
- `process_name` (string, 可选): 进程名称，用于关闭指定进程的所有连接

**响应状态码**:

- `200` - 网络连接已关闭

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`CloseConnectionRequest`](#type-closeconnectionrequest), [`CloseConnectionResponse`](#type-closeconnectionresponse)

**相关方法签名**: [`NetworkHandler.CloseConnection()`](#func-networkhandlercloseconnection), [`NetworkManager.CloseConnection()`](#func-networkmanagercloseconnection)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "关闭网络连接成功",
    "time": "2025-09-14T22:44:40+08:00",
    "data": {}
}
```

---

### 6.9 获取正在使用的端口

**接口地址**: `GET /api/v1/network/ports/in-use`

**功能描述**: 获取所有正在使用的端口信息，包括监听端口和已建立连接的端口

**查询参数**:

- `protocol` (string, 可选): 协议类型，可选值：tcp, udp, all（默认：all）
- `status` (string, 可选): 端口状态，可选值：listening, established, all（默认：all）
- `limit` (integer, 可选): 返回结果数量限制，默认：100，最大：1000

**响应状态码**:

- `200` - 端口信息获取成功

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`PortsInUseResponse`](#type-portsinuseresponse), [`PortInfo`](#type-portinfo)

**相关方法签名**: [`NetworkHandler.GetPortsInUse()`](#func-networkhandlergetportsinuse), [`NetworkManager.GetPortsInUse()`](#func-networkmanagergetportsinuse)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:46:08+08:00",
    "data": {
        "count": 123,
        "data": [
            3306,
            5672,
            8080,            
        ]
    }
}
```

---

### 6.10 检查端口占用

**接口地址**: `GET /api/v1/network/port/{port}/in-use`

**功能描述**: 检查指定端口是否被占用

**路径参数**:

- `port` (integer, 必需): 端口号，范围：1-65535

**参数说明**:

- `port`: 端口号，必须是系统中正在使用的端口号，如：8081

**响应状态码**:

- `200` - 端口占用状态

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`PortStatusResponse`](#type-portstatusresponse), [`PortStatus`](#type-portstatus)

**相关方法签名**: [`NetworkHandler.CheckPortInUse()`](#func-networkhandlercheckportinuse), [`NetworkManager.CheckPortInUse()`](#func-networkmanagercheckportinuse)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "端口占用状态检查完成",
    "time": "2025-09-14T22:46:58+08:00",
    "data": {
        "in_use": true,
        "port": 8081
    }
}
```

---

### 6.11 获取监听端口

**接口地址**: `GET /api/v1/network/listening-ports`

**功能描述**: 获取所有监听端口信息

**请求参数**: 无

**响应状态码**:

- `200` - 监听端口列表

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`ListeningPortsResponse`](#type-listeningportsresponse), [`ListeningPort`](#type-listeningport)

**相关方法签名**: [`NetworkHandler.GetListeningPorts()`](#func-networkhandlergetlisteningports), [`NetworkManager.GetListeningPorts()`](#func-networkmanagergetlisteningports)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:47:20+08:00",
    "data": {
        "count": 39,
        "data": [
            3306,
            5672,
            8080,        
        ]
    }
}
```

---

### 6.12 获取已建立连接

**接口地址**: `GET /api/v1/network/established-connections`

**功能描述**: 获取所有已建立的网络连接信息

**请求参数**: 无

**响应状态码**:

- `200` - 已建立连接列表

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`EstablishedConnectionsResponse`](#type-establishedconnectionsresponse), [`NetworkConnection`](#type-networkconnection)

**相关方法签名**: [`NetworkHandler.GetEstablishedConnections()`](#func-networkhandlergetestablishedconnections), [`NetworkManager.GetEstablishedConnections()`](#func-networkmanagergetestablishedconnections)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:47:54+08:00",
    "data": {
        "count": 38,
        "data": [
            {
                "id": "0",
                "local_addr": "10.44.1.102",
                "remote_addr": "104.18.18.125",
                "local_port": 52161,
                "remote_port": 443,
                "protocol": "TCP",
                "status": "ESTABLISHED",
                "pid": 39276,
                "process_name": "Cursor.exe",
                "type": "tcp"
            }
        ]
    }
}
```

---

### 6.13 获取网络接口

**接口地址**: `GET /api/v1/network/interfaces`

**功能描述**: 获取系统中所有网络接口信息

**请求参数**: 无

**响应状态码**:

- `200` - 网络接口列表

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`NetworkInterfacesResponse`](#type-networkinterfacesresponse), [`NetworkInterface`](#type-networkinterface)

**相关方法签名**: [`NetworkHandler.GetNetworkInterfaces()`](#func-networkhandlergetnetworkinterfaces), [`NetworkManager.GetNetworkInterfaces()`](#func-networkmanagergetnetworkinterfaces)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:48:19+08:00",
    "data": {
        "count": 10,
        "data": [
            {
                "name": "以太网",
                "index": 6,
                "addresses": [
                    "fe80::493f:f846:d3d:c551",
                    "169.254.51.133"
                ],
                "ipv4_addresses": [
                    "169.254.51.133"
                ],
                "ipv6_addresses": [
                    "fe80::493f:f846:d3d:c551"
                ],
                "mac_address": "bc:fc:e7:ef:91:a7",
                "is_up": false,
                "is_loopback": false,
                "is_multicast": true,
                "is_broadcast": true,
                "mtu": 1500,
                "speed": 0,
                "flags": "broadcast,multicast",
                "hardware_addr": "bc:fc:e7:ef:91:a7",
                "windows_info": {
                    "ip_config": {},
                    "statistics": {}
                }
            }
        ]
    }
}
```

---

### 6.14 获取网络统计信息

**接口地址**: `GET /api/v1/network/stats`

**功能描述**: 获取系统网络连接统计信息

**请求参数**: 无

**响应状态码**:

- `200` - 网络统计信息

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`NetworkStatsResponse`](#type-networkstatsresponse), [`NetworkStats`](#type-networkstats)

**相关方法签名**: [`NetworkHandler.GetNetworkStats()`](#func-networkhandlergetnetworkstats), [`NetworkManager.GetNetworkStats()`](#func-networkmanagergetnetworkstats)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取状态成功",
    "time": "2025-09-14T22:48:47+08:00",
    "data": {
        "timestamp": "2025-09-14T22:48:47.3717198+08:00",
        "monitored": true,
        "interfaces": [],
        "total_bytes": 308543969,
        "total_packets": 649519,
        "threat_count": 0
    }
}
```

---

### 6.15 启用网络监控

**接口地址**: `POST /api/v1/network/monitoring/enable`

**功能描述**: 启用网络连接监控功能，开始后台监控进程

**请求参数**: 无

**响应状态码**:

- `200` - 网络监控已启用

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`EnableMonitoringResponse`](#type-enablemonitoringresponse)

**相关方法签名**: [`NetworkHandler.EnableMonitoring()`](#func-networkhandlerenablemonitoring), [`NetworkManager.EnableMonitoring()`](#func-networkmanagerenablemonitoring)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "网络监控已启用成功",
    "time": "2025-09-14T22:48:57+08:00",
    "data": {}
}
```

---

### 6.16 禁用网络监控

**接口地址**: `POST /api/v1/network/monitoring/disable`

**功能描述**: 禁用网络连接监控功能，停止后台监控进程

**请求参数**: 无

**响应状态码**:

- `200` - 网络监控已禁用

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`DisableMonitoringResponse`](#type-disablemonitoringresponse)

**相关方法签名**: [`NetworkHandler.DisableMonitoring()`](#func-networkhandlerdisablemonitoring), [`NetworkManager.DisableMonitoring()`](#func-networkmanagerdisablemonitoring)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "网络监控已禁用成功",
    "time": "2025-09-14T22:49:04+08:00",
    "data": {}
}
```

---

### 6.17 获取监控连接列表

**接口地址**: `GET /api/v1/network/monitoring/connections`

**功能描述**: 获取当前监控的网络连接列表

**请求参数**: 无

**响应状态码**:

- `200` - 监控连接列表

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`MonitoringConnectionsResponse`](#type-monitoringconnectionsresponse), [`NetworkConnection`](#type-networkconnection)

**相关方法签名**: [`NetworkHandler.GetMonitoringConnections()`](#func-networkhandlergetmonitoringconnections), [`NetworkManager.GetMonitoringConnections()`](#func-networkmanagergetmonitoringconnections)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:49:12+08:00",
    "data": {
        "count": 0,
        "data": []
    }
}
```

---

### 6.18 获取连接历史记录

**接口地址**: `GET /api/v1/network/monitoring/history`

**功能描述**: 获取网络连接历史记录

**请求参数**: 无

**响应状态码**:

- `200` - 连接历史记录

**响应格式**: `application/json`

**标签**: 网络管理

**相关类型定义**: [`ConnectionHistoryResponse`](#type-connectionhistoryresponse), [`ConnectionHistory`](#type-connectionhistory)

**相关方法签名**: [`NetworkHandler.GetConnectionHistory()`](#func-networkhandlergetconnectionhistory), [`NetworkManager.GetConnectionHistory()`](#func-networkmanagergetconnectionhistory)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:49:21+08:00",
    "data": {
        "count": 0,
        "data": []
    }
}
```

---

## 👤 用户管理模块

### 7.1 获取当前用户信息

**接口地址**: `GET /api/v1/user/current`

**功能描述**: 获取当前登录用户信息

**请求参数**: 无

**响应状态码**:

- `200` - 当前用户信息

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`CurrentUserResponse`](#type-currentuserresponse), [`CurrentUser`](#type-currentuser)

**相关方法签名**: [`UserHandler.GetCurrentUser()`](#func-userhandlergetcurrentuser), [`UserManager.GetCurrentUser()`](#func-usermanagergetcurrentuser)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:49:32+08:00",
    "data": {
        "uid": "S-1-5-21-43760745-2046232404-1378798285-1001",
        "gid": "S-1-5-21-43760745-2046232404-1378798285-1001",
        "username": "DESKTOP-ASUS\\ASUS",
        "name": "yongsheng li",
        "home_dir": "C:\\Users\\ASUS",
        "shell": "",
        "last_login": "2025-09-14T22:49:32.6095587+08:00",
        "is_active": true,
        "is_admin": true,
        "password_age": 0,
        "groups": [
            "Administrators"
        ]
    }
}
```

---

### 7.2 获取所有用户列表

**接口地址**: `GET /api/v1/user/all`

**功能描述**: 获取系统所有用户列表

**请求参数**: 无

**响应状态码**:

- `200` - 用户列表

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`AllUsersResponse`](#type-allusersresponse), [`UserInfo`](#type-userinfo)

**相关方法签名**: [`UserHandler.GetAllUsers()`](#func-userhandlergetallusers), [`UserManager.GetAllUsers()`](#func-usermanagergetallusers)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:49:49+08:00",
    "data": {
        "count": 3,
        "data": [
            {
                "uid": "S-1-5-21-43760745-2046232404-1378798285-500",
                "gid": "S-1-5-21-43760745-2046232404-1378798285-513",
                "username": "DESKTOP-ASUS\\Administrator",
                "name": "Administrator",
                "home_dir": "C:\\Users\\Administrator",
                "shell": "",
                "last_login": "2025-09-14T22:49:48.7450438+08:00",
                "is_active": true,
                "is_admin": true,
                "password_age": 0,
                "groups": [
                    "Administrators",
                    "None"
                ]
            },
            {
                "uid": "S-1-5-21-43760745-2046232404-1378798285-501",
                "gid": "S-1-5-21-43760745-2046232404-1378798285-513",
                "username": "DESKTOP-ASUS\\Guest",
                "name": "Guest",
                "home_dir": "C:\\Users\\Guest",
                "shell": "",
                "last_login": "2025-09-14T22:49:49.6532077+08:00",
                "is_active": true,
                "is_admin": false,
                "password_age": 0,
                "groups": [
                    "Guests",
                    "None"
                ]
            },
            {
                "uid": "S-1-5-21-43760745-2046232404-1378798285-1001",
                "gid": "S-1-5-21-43760745-2046232404-1378798285-1001",
                "username": "DESKTOP-ASUS\\ASUS",
                "name": "yongsheng li",
                "home_dir": "C:\\Users\\ASUS",
                "shell": "",
                "last_login": "2025-09-14T22:49:49.6532077+08:00",
                "is_active": true,
                "is_admin": true,
                "password_age": 0,
                "groups": [
                    "Administrators"
                ]
            }
        ]
    }
}
```

---

### 7.3 根据用户 ID 获取用户信息

**接口地址**: `GET /api/v1/user/id/{uid}`

**功能描述**: 根据用户 ID 获取用户详细信息

**路径参数**:

- `uid` (string, 必需): 用户 ID，例如：S-1-5-21-43760745-2046232404-1378798285-1001

**请求参数**: 无

**响应状态码**:

- `200` - 用户信息

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`UserByIdResponse`](#type-userbyidresponse), [`UserDetail`](#type-userdetail)

**相关方法签名**: [`UserHandler.GetUserByID()`](#func-userhandlergetuserbyid), [`UserManager.GetUserByID()`](#func-usermanagergetuserbyid)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:50:20+08:00",
    "data": {
        "uid": "S-1-5-21-43760745-2046232404-1378798285-1001",
        "gid": "S-1-5-21-43760745-2046232404-1378798285-1001",
        "username": "DESKTOP-ASUS\\ASUS",
        "name": "yongsheng li",
        "home_dir": "C:\\Users\\ASUS",
        "shell": "",
        "last_login": "2025-09-14T22:50:20.6413149+08:00",
        "is_active": true,
        "is_admin": true,
        "password_age": 0,
        "groups": [
            "Administrators"
        ]
    }
}
```

---

### 7.4 根据用户名获取用户信息

**接口地址**: `GET /api/v1/user/name/{username}`

**功能描述**: 根据用户名获取用户详细信息

**路径参数**:

- `username` (string, 必需): 用户名，例如：Administrator

**请求参数**: 无

**响应状态码**:

- `200` - 用户信息

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`UserByNameResponse`](#type-userbynameresponse), [`UserDetail`](#type-userdetail)

**相关方法签名**: [`UserHandler.GetUserByName()`](#func-userhandlergetuserbyname), [`UserManager.GetUserByName()`](#func-usermanagergetuserbyname)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:50:50+08:00",
    "data": {
        "uid": "S-1-5-21-43760745-2046232404-1378798285-500",
        "gid": "S-1-5-21-43760745-2046232404-1378798285-513",
        "username": "DESKTOP-ASUS\\Administrator",
        "name": "Administrator",
        "home_dir": "C:\\Users\\Administrator",
        "shell": "",
        "last_login": "2025-09-14T22:50:50.5633413+08:00",
        "is_active": true,
        "is_admin": true,
        "password_age": 0,
        "groups": [
            "Administrators",
            "None"
        ]
    }
}
```

---

### 7.5 获取用户组

**接口地址**: `GET /api/v1/user/groups/{username}`

**功能描述**: 获取指定用户所属的组

**路径参数**:

- `username` (string, 必需): 用户名，例如：Administrator

**请求参数**: 无

**响应状态码**:

- `200` - 用户组信息

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`UserGroupsResponse`](#type-usergroupsresponse), [`UserGroup`](#type-usergroup)

**相关方法签名**: [`UserHandler.GetUserGroups()`](#func-userhandlergetusergroups), [`UserManager.GetUserGroups()`](#func-usermanagergetusergroups)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:52:28+08:00",
    "data": {
        "count": 2,
        "data": [
            "Administrators",
            "None"
        ]
    }
}
```

---

### 7.6 检查用户权限

**接口地址**: `POST /api/v1/user/permissions/check`

**功能描述**: 检查指定用户的权限

**请求体参数**:

- `username` (string, 必需): 用户名，例如：Guest
- `permission` (string, 必需): 权限名称，例如：admin

**响应状态码**:

- `200` - 权限检查结果

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`CheckPermissionRequest`](#type-checkpermissionrequest), [`CheckPermissionResponse`](#type-checkpermissionresponse)

**相关方法签名**: [`UserHandler.CheckPermission()`](#func-userhandlercheckpermission), [`UserManager.CheckPermission()`](#func-usermanagercheckpermission)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:52:58+08:00",
    "data": {
        "username": "Guest",
        "permissions": {
            "admin": false
        },
        "is_admin": false,
        "check_time": "2025-09-14T22:52:58.4487147+08:00"
    }
}
```

---

### 7.7 验证用户密码

**接口地址**: `POST /api/v1/user/password/validate`

**功能描述**: 验证用户密码是否符合安全策略

**请求体参数**:

- `username` (string, 必需): 用户名，例如：testuser
- `password` (string, 必需): 密码，例如：TestPassword123!

**响应状态码**:

- `200` - 密码验证结果

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`ValidatePasswordRequest`](#type-validatepasswordrequest), [`ValidatePasswordResponse`](#type-validatepasswordresponse)

**相关方法签名**: [`UserHandler.ValidatePassword()`](#func-userhandlervalidatepassword), [`UserManager.ValidatePassword()`](#func-usermanagervalidatepassword)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:54:06+08:00",
    "data": {
        "is_valid": {
            "is_valid": true
        }
    }
}
```

---

### 7.8 获取密码策略

**接口地址**: `GET /api/v1/user/password/policy`

**功能描述**: 获取系统密码策略信息

**响应状态码**:

- `200` - 密码策略信息

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`PasswordPolicyResponse`](#type-passwordpolicyresponse), [`PasswordPolicy`](#type-passwordpolicy)

**相关方法签名**: [`UserHandler.GetPasswordPolicy()`](#func-userhandlergetpasswordpolicy), [`UserManager.GetPasswordPolicy()`](#func-usermanagergetpasswordpolicy)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取信息成功",
    "time": "2025-09-14T22:54:15+08:00",
    "data": {
        "min_length": 8,
        "require_uppercase": true,
        "require_lowercase": true,
        "require_numbers": true,
        "require_special_chars": true,
        "max_age": 90,
        "history_count": 5,
        "lockout_threshold": 5,
        "lockout_duration": 30
    }
}
```

---

### 7.9 检查账户状态

**接口地址**: `GET /api/v1/user/status/{username}`

**功能描述**: 检查指定用户账户的状态

**路径参数**:

- `username` (string, 必需): 用户名，例如：Guest

- `200` - 账户状态

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`UserStatusResponse`](#type-userstatusresponse), [`UserStatus`](#type-userstatus)

**相关方法签名**: [`UserHandler.CheckUserStatus()`](#func-userhandlercheckuserstatus), [`UserManager.CheckUserStatus()`](#func-usermanagercheckuserstatus)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取状态成功",
    "time": "2025-09-14T22:54:33+08:00",
    "data": {
        "username": "DESKTOP-ASUS\\Guest",
        "is_active": true,
        "is_locked": false,
        "is_password_expired": false,
        "last_login": "2025-09-14T22:54:33.0660757+08:00",
        "check_time": "2025-09-14T22:54:33.0660757+08:00"
    }
}
```

---

### 7.10 获取用户会话

**接口地址**: `GET /api/v1/user/sessions/{username}`

**功能描述**: 获取指定用户的会话信息

**路径参数**:

- `username` (string, 必需): 用户名，例如：Guest

**响应状态码**:

- `200` - 用户会话信息

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`UserSessionsResponse`](#type-usersessionsresponse), [`UserSession`](#type-usersession)

**相关方法签名**: [`UserHandler.GetUserSessions()`](#func-userhandlergetusersessions), [`UserManager.GetUserSessions()`](#func-usermanagergetusersessions)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:54:54+08:00",
    "data": {
        "count": 1,
        "data": [
            {
                "session_id": "session_1757861694",
                "username": "Guest",
                "login_time": "2025-09-14T21:54:54.2989475+08:00",
                "last_active": "2025-09-14T22:54:54.2989475+08:00",
                "ip_address": "127.0.0.1",
                "user_agent": "Console Session",
                "is_active": true
            }
        ]
    }
}
```

---

### 7.11 结束用户会话

**接口地址**: `DELETE /api/v1/user/session/{sessionId}`

**功能描述**: 结束指定的用户会话

**路径参数**:

- `sessionId` (string, 必需): 会话 ID，例如：session_1757861694

**响应状态码**:

- `200` - 会话结束成功

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`EndSessionResponse`](#type-endsessionresponse)

**相关方法签名**: [`UserHandler.EndUserSession()`](#func-userhandlerendusersession), [`UserManager.EndUserSession()`](#func-usermanagerendusersession)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "结束用户会话成功",
    "time": "2025-09-14T22:55:18+08:00",
    "data": {}
}
```

---

### 7.12 锁定用户账户

**接口地址**: `POST /api/v1/user/lock/{username}`

**功能描述**: 锁定指定的用户账户

**路径参数**:

- `username` (string, 必需): 用户名，例如：Guest

**响应状态码**:

- `200` - 账户锁定成功

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`LockUserResponse`](#type-lockuserresponse)

**相关方法签名**: [`UserHandler.LockUserAccount()`](#func-userhandlerlockuseraccount), [`UserManager.LockUserAccount()`](#func-usermanagerlockuseraccount)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "锁定用户账户成功",
    "time": "2025-09-14T22:55:47+08:00",
    "data": {}
}
```

---

### 7.13 解锁用户账户

**接口地址**: `POST /api/v1/user/unlock/{username}`

**功能描述**: 解锁指定的用户账户

**路径参数**:

- `username` (string, 必需): 用户名，例如：Guest

**响应状态码**:

- `200` - 账户解锁成功

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`UnlockUserResponse`](#type-unlockuserresponse)

**相关方法签名**: [`UserHandler.UnlockUserAccount()`](#func-userhandlerunlockuseraccount), [`UserManager.UnlockUserAccount()`](#func-usermanagerunlockuseraccount)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "解锁用户账户成功",
    "time": "2025-09-14T22:56:13+08:00",
    "data": {}
}
```

---

### 7.14 修改用户密码

**接口地址**: `POST /api/v1/user/password/change`

**功能描述**: 修改指定用户的密码

**请求体参数**:

- `username` (string, 必需): 用户名，例如：Guest
- `new_password` (string, 必需): 新密码，例如：NewPassword123!

**响应状态码**:

- `200` - 密码修改成功

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`ChangePasswordRequest`](#type-changepasswordrequest), [`ChangePasswordResponse`](#type-changepasswordresponse)

**相关方法签名**: [`UserHandler.ChangePassword()`](#func-userhandlerchangepassword), [`UserManager.ChangePassword()`](#func-usermanagerchangepassword)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "修改用户密码成功",
    "time": "2025-09-14T22:56:37+08:00",
    "data": {}
}
```

---

### 7.15 获取用户登录历史

**接口地址**: `GET /api/v1/user/history/{username}`

**功能描述**: 获取指定用户的登录历史记录

**路径参数**:

- `username` (string, 必需): 用户名，例如：Guest

**响应状态码**:

- `200` - 登录历史记录

**响应格式**: `application/json`

**标签**: 用户管理

**相关类型定义**: [`UserHistoryResponse`](#type-userhistoryresponse), [`LoginHistory`](#type-loginhistory)

**相关方法签名**: [`UserHandler.GetUserHistory()`](#func-userhandlergetuserhistory), [`UserManager.GetUserHistory()`](#func-usermanagergetuserhistory)

**处理结果示例**:

```json
{
    "code": 200,
    "message": "获取列表成功",
    "time": "2025-09-14T22:57:04+08:00",
    "data": {
        "count": 10,
        "data": [
            {
                "username": "Guest",
                "login_time": "2025-09-14T22:57:04.1922737+08:00",
                "logout_time": "2025-09-14T23:57:04.1922737+08:00",
                "ip_address": "192.168.1.100",
                "user_agent": "Windows NT 10.0; Win64; x64",
                "status": "Success"
            }
        ]
    }
}
```

---

**文档版本** 5.3  
**最后更新** 2025-09-14  
**维护人员** LYS