# Yara安全服务完整API验证指南

## 📋 概述

本文档提供了Yara安全服务所有83个API接口的完整验证指南，包括测试步骤、预期结果和验证要点。

## 📋 目录索引

### 🔍 文件管理模块验证 (11个接口)
- [1.1 扫描单个文件](#11-扫描单个文件)
- [1.2 获取文件信息](#12-获取文件信息)
- [1.3 扫描目录](#13-扫描目录)
- [1.4 扫描内存缓冲区](#14-扫描内存缓冲区)
- [1.5 获取文件列表](#15-获取文件列表)
- [1.6 复制文件](#16-复制文件)
- [1.7 移动文件](#17-移动文件)
- [1.8 删除文件](#18-删除文件)
- [1.9 获取文件哈希](#19-获取文件哈希)
- [1.10 获取文件所有哈希值](#110-获取文件所有哈希值)
- [1.11 验证文件哈希](#111-验证文件哈希)

### 👥 进程管理模块验证 (20个接口)
- [2.1 获取进程列表](#21-获取进程列表)
- [2.2 获取进程信息](#22-获取进程信息)
- [2.3 启动进程](#23-启动进程)
- [2.4 结束进程](#24-结束进程)
- [2.5 挂起进程](#25-挂起进程)
- [2.6 恢复进程](#26-恢复进程)
- [2.7 获取进程模块](#27-获取进程模块)
- [2.8 获取进程连接](#28-获取进程连接)
- [2.9 获取进程内存信息](#29-获取进程内存信息)
- [2.10 检查进程运行状态](#210-检查进程运行状态)
- [2.11 获取进程子进程](#211-获取进程子进程)
- [2.12 获取系统模块列表](#212-获取系统模块列表)
- [2.13 获取模块信息](#213-获取模块信息)
- [2.14 挂起模块](#214-挂起模块)
- [2.15 恢复模块](#215-恢复模块)
- [2.16 结束模块](#216-结束模块)
- [2.17 启用进程监控](#217-启用进程监控)
- [2.18 禁用进程监控](#218-禁用进程监控)
- [2.19 获取监控进程列表](#219-获取监控进程列表)
- [2.20 获取进程统计信息](#220-获取进程统计信息)

### 🔧 注册表操作模块验证 (9个接口)
- [3.1 获取注册表键信息](#31-获取注册表键信息)
- [3.2 创建注册表键](#32-创建注册表键)
- [3.3 删除注册表键](#33-删除注册表键)
- [3.4 设置注册表值](#34-设置注册表值)
- [3.5 获取注册表值](#35-获取注册表值)
- [3.6 删除注册表值](#36-删除注册表值)
- [3.7 列出注册表键](#37-列出注册表键)
- [3.8 列出注册表值](#38-列出注册表值)
- [3.9 搜索注册表](#39-搜索注册表)

### 🌐 网络管理模块验证 (16个接口)
- [4.1 获取网络连接](#41-获取网络连接)
- [4.2 获取TCP连接](#42-获取tcp连接)
- [4.3 获取UDP连接](#43-获取udp连接)
- [4.4 按进程ID获取连接](#44-按进程id获取连接)
- [4.5 按端口获取连接](#45-按端口获取连接)
- [4.6 按IP获取连接](#46-按ip获取连接)
- [4.7 关闭网络连接](#47-关闭网络连接)
- [4.8 检查端口占用](#48-检查端口占用)
- [4.9 获取监听端口](#49-获取监听端口)
- [4.10 获取已建立连接](#410-获取已建立连接)
- [4.11 获取网络接口](#411-获取网络接口)
- [4.12 获取网络统计](#412-获取网络统计)
- [4.13 启用网络监控](#413-启用网络监控)
- [4.14 禁用网络监控](#414-禁用网络监控)
- [4.15 获取监控连接列表](#415-获取监控连接列表)
- [4.16 获取连接历史记录](#416-获取连接历史记录)

### 👤 用户管理模块验证 (15个接口)
- [5.1 获取当前用户信息](#51-获取当前用户信息)
- [5.2 根据用户ID获取用户信息](#52-根据用户id获取用户信息)
- [5.3 根据用户名获取用户信息](#53-根据用户名获取用户信息)
- [5.4 获取所有用户列表](#54-获取所有用户列表)
- [5.5 检查用户权限](#55-检查用户权限)
- [5.6 验证用户密码](#56-验证用户密码)
- [5.7 获取密码策略](#57-获取密码策略)
- [5.8 检查账户状态](#58-检查账户状态)
- [5.9 获取用户会话](#59-获取用户会话)
- [5.10 结束用户会话](#510-结束用户会话)
- [5.11 锁定用户账户](#511-锁定用户账户)
- [5.12 解锁用户账户](#512-解锁用户账户)
- [5.13 修改用户密码](#513-修改用户密码)
- [5.14 获取用户组](#514-获取用户组)
- [5.15 获取用户历史记录](#515-获取用户历史记录)

### 🛡️ 安全管理模块验证 (9个接口)
- [6.1 获取安全状态](#61-获取安全状态)
- [6.2 获取规则信息](#62-获取规则信息)
- [6.3 重新加载规则](#63-重新加载规则)
- [6.4 清空缓存](#64-清空缓存)
- [6.5 获取缓存统计](#65-获取缓存统计)
- [6.6 隔离文件](#66-隔离文件)
- [6.7 恢复文件](#67-恢复文件)
- [6.8 获取隔离列表](#68-获取隔离列表)
- [6.9 获取扫描历史](#69-获取扫描历史)

### 📊 系统监控模块验证 (3个接口)
- [7.1 健康检查](#71-健康检查)
- [7.2 获取性能指标](#72-获取性能指标)
- [7.3 重置性能指标](#73-重置性能指标)

---

## 📊 API接口统计

**总计**: 83个API接口

### 🔍 文件管理模块 (11个接口)
1. `POST /api/v1/file/scan` - 扫描单个文件
2. `GET /api/v1/file/info/{path}` - 获取文件信息
3. `POST /api/v1/file/scan-directory` - 扫描目录
4. `POST /api/v1/file/scan-buffer` - 扫描内存缓冲区
5. `GET /api/v1/file/list` - 获取文件列表
6. `POST /api/v1/file/copy` - 复制文件
7. `POST /api/v1/file/move` - 移动文件
8. `DELETE /api/v1/file/{path}` - 删除文件
9. `GET /api/v1/file/hash/{path}` - 获取文件哈希
10. `GET /api/v1/file/hashes/{path}` - 获取文件所有哈希值
11. `POST /api/v1/file/verify-hash` - 验证文件哈希

### 👥 进程管理模块 (20个接口)
12. `GET /api/v1/process/list` - 获取进程列表
13. `GET /api/v1/process/{pid}` - 获取进程信息
14. `POST /api/v1/process/start` - 启动进程
15. `DELETE /api/v1/process/{pid}` - 结束进程
16. `PUT /api/v1/process/{pid}/suspend` - 挂起进程
17. `PUT /api/v1/process/{pid}/resume` - 恢复进程
18. `GET /api/v1/process/{pid}/modules` - 获取进程模块
19. `GET /api/v1/process/{pid}/connections` - 获取进程连接
20. `GET /api/v1/process/{pid}/memory` - 获取进程内存信息
21. `GET /api/v1/process/{pid}/running` - 检查进程运行状态
22. `GET /api/v1/process/{pid}/children` - 获取进程子进程
23. `GET /api/v1/process/modules` - 获取系统模块列表
24. `GET /api/v1/process/module/{module}` - 获取模块信息
25. `PUT /api/v1/process/module/{module}/suspend` - 挂起模块
26. `PUT /api/v1/process/module/{module}/resume` - 恢复模块
27. `DELETE /api/v1/process/module/{module}` - 结束模块
28. `POST /api/v1/process/monitoring/enable` - 启用进程监控
29. `POST /api/v1/process/monitoring/disable` - 禁用进程监控
30. `GET /api/v1/process/monitoring/list` - 获取监控进程列表
31. `GET /api/v1/process/statistics` - 获取进程统计信息

### 🔧 注册表操作模块 (9个接口)
32. `GET /api/v1/registry/key/{path}` - 获取注册表键信息
33. `POST /api/v1/registry/key` - 创建注册表键
34. `DELETE /api/v1/registry/key/{path}` - 删除注册表键
35. `PUT /api/v1/registry/value` - 设置注册表值
36. `GET /api/v1/registry/value/{path}/{name}` - 获取注册表值
37. `DELETE /api/v1/registry/value/{path}/{name}` - 删除注册表值
38. `GET /api/v1/registry/keys/{path}` - 列出注册表键
39. `GET /api/v1/registry/values/{path}` - 列出注册表值
40. `GET /api/v1/registry/search` - 搜索注册表

### 🌐 网络管理模块 (16个接口)
41. `GET /api/v1/network/connections` - 获取网络连接
42. `GET /api/v1/network/connections/tcp` - 获取TCP连接
43. `GET /api/v1/network/connections/udp` - 获取UDP连接
44. `GET /api/v1/network/connections/pid/{pid}` - 按进程ID获取连接
45. `GET /api/v1/network/connections/port/{port}` - 按端口获取连接
46. `GET /api/v1/network/connections/ip/{ip}` - 按IP获取连接
47. `DELETE /api/v1/network/connection/{id}` - 关闭网络连接
48. `GET /api/v1/network/port/{port}/in-use` - 检查端口占用
49. `GET /api/v1/network/listening-ports` - 获取监听端口
50. `GET /api/v1/network/established-connections` - 获取已建立连接
51. `GET /api/v1/network/interfaces` - 获取网络接口
52. `GET /api/v1/network/stats` - 获取网络统计
53. `POST /api/v1/network/monitoring/enable` - 启用网络监控
54. `POST /api/v1/network/monitoring/disable` - 禁用网络监控
55. `GET /api/v1/network/monitoring/connections` - 获取监控连接列表
56. `GET /api/v1/network/monitoring/history` - 获取连接历史记录

### 👤 用户管理模块 (15个接口)
57. `GET /api/v1/user/current` - 获取当前用户信息
58. `GET /api/v1/user/id/{uid}` - 根据用户ID获取用户信息
59. `GET /api/v1/user/name/{username}` - 根据用户名获取用户信息
60. `GET /api/v1/user/all` - 获取所有用户列表
61. `POST /api/v1/user/permissions/check` - 检查用户权限
62. `POST /api/v1/user/password/validate` - 验证用户密码
63. `GET /api/v1/user/password/policy` - 获取密码策略
64. `GET /api/v1/user/status/{username}` - 检查账户状态
65. `GET /api/v1/user/sessions/{username}` - 获取用户会话
66. `DELETE /api/v1/user/session/{sessionId}` - 结束用户会话
67. `POST /api/v1/user/lock/{username}` - 锁定用户账户
68. `POST /api/v1/user/unlock/{username}` - 解锁用户账户
69. `POST /api/v1/user/password/change` - 修改用户密码
70. `GET /api/v1/user/groups/{username}` - 获取用户组
71. `GET /api/v1/user/history/{username}` - 获取用户历史记录

### 🛡️ 安全管理模块 (9个接口)
72. `GET /api/v1/security/status` - 获取安全状态
73. `GET /api/v1/security/rules` - 获取规则信息
74. `POST /api/v1/security/reload-rules` - 重新加载规则
75. `POST /api/v1/security/cache/clear` - 清空缓存
76. `GET /api/v1/security/cache/stats` - 获取缓存统计
77. `POST /api/v1/security/quarantine` - 隔离文件
78. `POST /api/v1/security/restore` - 恢复文件
79. `GET /api/v1/security/quarantine/list` - 获取隔离列表
80. `GET /api/v1/security/scan-history` - 获取扫描历史

### 📊 系统监控模块 (3个接口)
81. `GET /api/health` - 健康检查
82. `GET /api/metrics` - 获取性能指标
83. `POST /api/metrics/reset` - 重置性能指标

---

## 🔍 文件管理模块验证 (11个接口)

### 1.1 扫描单个文件
**接口**: `POST /api/v1/file/scan`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\test.txt\"}"
```

**预期结果**:
- 返回扫描结果，包含文件信息和威胁检测结果
- 状态码: 200

### 1.2 获取文件信息
**接口**: `GET /api/v1/file/info/{path}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/file/info/C%3A%5Ctemp%5Ctest.txt"
```

**预期结果**:
- 返回文件详细信息（大小、时间、权限等）
- 状态码: 200

### 1.3 扫描目录
**接口**: `POST /api/v1/file/scan-directory`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/file/scan-directory \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\"}"
```

**预期结果**:
- 返回目录扫描结果
- 状态码: 200

### 1.4 扫描内存缓冲区
**接口**: `POST /api/v1/file/scan-buffer`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"buffer\":\"base64_encoded_data\"}"
```

**预期结果**:
- 返回缓冲区扫描结果
- 状态码: 200

### 1.5 获取文件列表
**接口**: `GET /api/v1/file/list`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/file/list?path=C%3A%5Ctemp"
```

**预期结果**:
- 返回指定目录的文件列表
- 状态码: 200

### 1.6 复制文件
**接口**: `POST /api/v1/file/copy`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/file/copy \
  -H "Content-Type: application/json" \
  -d "{\"source\":\"C:\\\\temp\\\\source.txt\",\"dest\":\"C:\\\\temp\\\\dest.txt\"}"
```

**预期结果**:
- 返回复制操作结果
- 状态码: 200

### 1.7 移动文件
**接口**: `POST /api/v1/file/move`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/file/move \
  -H "Content-Type: application/json" \
  -d "{\"source\":\"C:\\\\temp\\\\source.txt\",\"dest\":\"C:\\\\temp\\\\moved.txt\"}"
```

**预期结果**:
- 返回移动操作结果
- 状态码: 200

### 1.8 删除文件
**接口**: `DELETE /api/v1/file/{path}`

**测试步骤**:
```bash
curl -X DELETE "http://localhost:8081/api/v1/file/C%3A%5Ctemp%5Ctest.txt"
```

**预期结果**:
- 返回删除操作结果
- 状态码: 200

### 1.9 获取文件哈希
**接口**: `GET /api/v1/file/hash/{path}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/file/hash/C%3A%5Ctemp%5Ctest.txt?algorithm=sha256"
```

**预期结果**:
- 返回文件哈希值
- 状态码: 200

### 1.10 获取文件所有哈希值
**接口**: `GET /api/v1/file/hashes/{path}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/file/hashes/C%3A%5Ctemp%5Ctest.txt"
```

**预期结果**:
- 返回文件的所有哈希值（MD5、SHA1、SHA256等）
- 状态码: 200

### 1.11 验证文件哈希
**接口**: `POST /api/v1/file/verify-hash`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/file/verify-hash \
  -H "Content-Type: application/json" \
  -d "{\"file_path\":\"C:\\\\temp\\\\test.txt\",\"algorithm\":\"sha256\",\"expected_hash\":\"expected_hash_value\"}"
```

**预期结果**:
- 返回哈希验证结果
- 状态码: 200

---

## 👥 进程管理模块验证 (20个接口)

### 2.1 获取进程列表
**接口**: `GET /api/v1/process/list`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/list"
```

**预期结果**:
- 返回系统进程列表
- 状态码: 200

### 2.2 获取进程信息
**接口**: `GET /api/v1/process/{pid}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/1234"
```

**预期结果**:
- 返回指定进程的详细信息
- 状态码: 200

### 2.3 启动进程
**接口**: `POST /api/v1/process/start`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/process/start \
  -H "Content-Type: application/json" \
  -d "{\"command\":\"notepad.exe\"}"
```

**预期结果**:
- 返回进程启动结果
- 状态码: 200

### 2.4 结束进程
**接口**: `DELETE /api/v1/process/{pid}`

**测试步骤**:
```bash
curl -X DELETE "http://localhost:8081/api/v1/process/1234"
```

**预期结果**:
- 返回进程结束结果
- 状态码: 200

### 2.5 挂起进程
**接口**: `PUT /api/v1/process/{pid}/suspend`

**测试步骤**:
```bash
curl -X PUT "http://localhost:8081/api/v1/process/1234/suspend"
```

**预期结果**:
- 返回进程挂起结果
- 状态码: 200

### 2.6 恢复进程
**接口**: `PUT /api/v1/process/{pid}/resume`

**测试步骤**:
```bash
curl -X PUT "http://localhost:8081/api/v1/process/1234/resume"
```

**预期结果**:
- 返回进程恢复结果
- 状态码: 200

### 2.7 获取进程模块
**接口**: `GET /api/v1/process/{pid}/modules`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/1234/modules"
```

**预期结果**:
- 返回进程加载的模块列表
- 状态码: 200

### 2.8 获取进程连接
**接口**: `GET /api/v1/process/{pid}/connections`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/1234/connections"
```

**预期结果**:
- 返回进程的网络连接列表
- 状态码: 200

### 2.9 获取进程内存信息
**接口**: `GET /api/v1/process/{pid}/memory`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/1234/memory"
```

**预期结果**:
- 返回进程的内存使用信息
- 状态码: 200

### 2.10 检查进程运行状态
**接口**: `GET /api/v1/process/{pid}/running`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/1234/running"
```

**预期结果**:
- 返回进程运行状态
- 状态码: 200

### 2.11 获取进程子进程
**接口**: `GET /api/v1/process/{pid}/children`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/1234/children"
```

**预期结果**:
- 返回进程的子进程列表
- 状态码: 200

### 2.12 获取系统模块列表
**接口**: `GET /api/v1/process/modules`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/modules"
```

**预期结果**:
- 返回系统模块列表
- 状态码: 200

### 2.13 获取模块信息
**接口**: `GET /api/v1/process/module/{module}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/module/kernel32.dll"
```

**预期结果**:
- 返回指定模块的详细信息
- 状态码: 200

### 2.14 挂起模块
**接口**: `PUT /api/v1/process/module/{module}/suspend`

**测试步骤**:
```bash
curl -X PUT "http://localhost:8081/api/v1/process/module/suspicious.dll/suspend"
```

**预期结果**:
- 返回模块挂起结果
- 状态码: 200

### 2.15 恢复模块
**接口**: `PUT /api/v1/process/module/{module}/resume`

**测试步骤**:
```bash
curl -X PUT "http://localhost:8081/api/v1/process/module/suspicious.dll/resume"
```

**预期结果**:
- 返回模块恢复结果
- 状态码: 200

### 2.16 结束模块
**接口**: `DELETE /api/v1/process/module/{module}`

**测试步骤**:
```bash
curl -X DELETE "http://localhost:8081/api/v1/process/module/suspicious.dll"
```

**预期结果**:
- 返回模块结束结果
- 状态码: 200

### 2.17 启用进程监控
**接口**: `POST /api/v1/process/monitoring/enable`

**测试步骤**:
```bash
curl -X POST "http://localhost:8081/api/v1/process/monitoring/enable"
```

**预期结果**:
- 返回进程监控启用结果
- 状态码: 200

### 2.18 禁用进程监控
**接口**: `POST /api/v1/process/monitoring/disable`

**测试步骤**:
```bash
curl -X POST "http://localhost:8081/api/v1/process/monitoring/disable"
```

**预期结果**:
- 返回进程监控禁用结果
- 状态码: 200

### 2.19 获取监控进程列表
**接口**: `GET /api/v1/process/monitoring/list`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/monitoring/list"
```

**预期结果**:
- 返回被监控的进程列表
- 状态码: 200

### 2.20 获取进程统计信息
**接口**: `GET /api/v1/process/statistics`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/process/statistics"
```

**预期结果**:
- 返回进程统计信息
- 状态码: 200

---

## 🔧 注册表操作模块验证 (9个接口)

### 3.1 获取注册表键信息
**接口**: `GET /api/v1/registry/key/{path}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/registry/key/SOFTWARE%5CMicrosoft%5CWindows%5CCurrentVersion"
```

**预期结果**:
- 返回注册表键信息
- 状态码: 200

### 3.2 创建注册表键
**接口**: `POST /api/v1/registry/key`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/registry/key \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"SOFTWARE\\\\TestKey\"}"
```

**预期结果**:
- 返回创建结果
- 状态码: 200

### 3.3 删除注册表键
**接口**: `DELETE /api/v1/registry/key/{path}`

**测试步骤**:
```bash
curl -X DELETE "http://localhost:8081/api/v1/registry/key/SOFTWARE%5CTestKey"
```

**预期结果**:
- 返回删除结果
- 状态码: 200

### 3.4 设置注册表值
**接口**: `PUT /api/v1/registry/value`

**测试步骤**:
```bash
curl -X PUT http://localhost:8081/api/v1/registry/value \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"SOFTWARE\\\\TestKey\",\"name\":\"TestValue\",\"value\":\"test\",\"type\":\"REG_SZ\"}"
```

**预期结果**:
- 返回设置结果
- 状态码: 200

### 3.5 获取注册表值
**接口**: `GET /api/v1/registry/value/{path}/{name}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/registry/value/SOFTWARE%5CTestKey/TestValue"
```

**预期结果**:
- 返回注册表值信息
- 状态码: 200

### 3.6 删除注册表值
**接口**: `DELETE /api/v1/registry/value/{path}/{name}`

**测试步骤**:
```bash
curl -X DELETE "http://localhost:8081/api/v1/registry/value/SOFTWARE%5CTestKey/TestValue"
```

**预期结果**:
- 返回删除结果
- 状态码: 200

### 3.7 列出注册表键
**接口**: `GET /api/v1/registry/keys/{path}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/registry/keys/SOFTWARE%5CMicrosoft"
```

**预期结果**:
- 返回子键列表
- 状态码: 200

### 3.8 列出注册表值
**接口**: `GET /api/v1/registry/values/{path}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/registry/values/SOFTWARE%5CMicrosoft%5CWindows%5CCurrentVersion"
```

**预期结果**:
- 返回值列表
- 状态码: 200

### 3.9 搜索注册表
**接口**: `GET /api/v1/registry/search`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/registry/search?query=test&path=SOFTWARE"
```

**预期结果**:
- 返回搜索结果
- 状态码: 200

---

## 🌐 网络管理模块验证 (16个接口)

### 4.1 获取网络连接
**接口**: `GET /api/v1/network/connections`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/connections"
```

**预期结果**:
- 返回所有网络连接
- 状态码: 200

### 4.2 获取TCP连接
**接口**: `GET /api/v1/network/connections/tcp`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/connections/tcp"
```

**预期结果**:
- 返回TCP连接列表
- 状态码: 200

### 4.3 获取UDP连接
**接口**: `GET /api/v1/network/connections/udp`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/connections/udp"
```

**预期结果**:
- 返回UDP连接列表
- 状态码: 200

### 4.4 按进程ID获取连接
**接口**: `GET /api/v1/network/connections/pid/{pid}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/connections/pid/1234"
```

**预期结果**:
- 返回指定进程的连接
- 状态码: 200

### 4.5 按端口获取连接
**接口**: `GET /api/v1/network/connections/port/{port}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/connections/port/80"
```

**预期结果**:
- 返回指定端口的连接
- 状态码: 200

### 4.6 按IP获取连接
**接口**: `GET /api/v1/network/connections/ip/{ip}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/connections/ip/192.168.1.1"
```

**预期结果**:
- 返回指定IP的连接
- 状态码: 200

### 4.7 关闭网络连接
**接口**: `DELETE /api/v1/network/connection/{id}`

**测试步骤**:
```bash
curl -X DELETE "http://localhost:8081/api/v1/network/connection/1234-tcp-0"
```

**预期结果**:
- 返回连接关闭结果
- 状态码: 200

### 4.8 检查端口占用
**接口**: `GET /api/v1/network/port/{port}/in-use`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/port/80/in-use"
```

**预期结果**:
- 返回端口占用状态
- 状态码: 200

### 4.9 获取监听端口
**接口**: `GET /api/v1/network/listening-ports`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/listening-ports"
```

**预期结果**:
- 返回监听端口列表
- 状态码: 200

### 4.10 获取已建立连接
**接口**: `GET /api/v1/network/established-connections`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/established-connections"
```

**预期结果**:
- 返回已建立的连接
- 状态码: 200

### 4.11 获取网络接口
**接口**: `GET /api/v1/network/interfaces`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/interfaces"
```

**预期结果**:
- 返回网络接口信息
- 状态码: 200

### 4.12 获取网络统计
**接口**: `GET /api/v1/network/stats`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/stats"
```

**预期结果**:
- 返回网络统计信息
- 状态码: 200

### 4.13 启用网络监控
**接口**: `POST /api/v1/network/monitoring/enable`

**测试步骤**:
```bash
curl -X POST "http://localhost:8081/api/v1/network/monitoring/enable"
```

**预期结果**:
- 返回网络监控启用结果
- 状态码: 200

### 4.14 禁用网络监控
**接口**: `POST /api/v1/network/monitoring/disable`

**测试步骤**:
```bash
curl -X POST "http://localhost:8081/api/v1/network/monitoring/disable"
```

**预期结果**:
- 返回网络监控禁用结果
- 状态码: 200

### 4.15 获取监控连接列表
**接口**: `GET /api/v1/network/monitoring/connections`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/monitoring/connections"
```

**预期结果**:
- 返回被监控的连接列表
- 状态码: 200

### 4.16 获取连接历史记录
**接口**: `GET /api/v1/network/monitoring/history`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/network/monitoring/history"
```

**预期结果**:
- 返回连接历史记录
- 状态码: 200

---

## 👤 用户管理模块验证 (15个接口)

### 5.1 获取当前用户信息
**接口**: `GET /api/v1/user/current`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/user/current"
```

**预期结果**:
- 返回当前用户信息
- 状态码: 200

### 5.2 根据用户ID获取用户信息
**接口**: `GET /api/v1/user/id/{uid}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/user/id/S-1-5-21-1000"
```

**预期结果**:
- 返回指定用户信息
- 状态码: 200

### 5.3 根据用户名获取用户信息
**接口**: `GET /api/v1/user/name/{username}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/user/name/Administrator"
```

**预期结果**:
- 返回指定用户信息
- 状态码: 200

### 5.4 获取所有用户列表
**接口**: `GET /api/v1/user/all`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/user/all"
```

**预期结果**:
- 返回所有用户列表
- 状态码: 200

### 5.5 检查用户权限
**接口**: `POST /api/v1/user/permissions/check`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/user/permissions/check \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"Administrator\",\"permissions\":[\"read\",\"write\"]}"
```

**预期结果**:
- 返回权限检查结果
- 状态码: 200

### 5.6 验证用户密码
**接口**: `POST /api/v1/user/password/validate`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/user/password/validate \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"Administrator\",\"password\":\"test123\"}"
```

**预期结果**:
- 返回密码验证结果
- 状态码: 200

### 5.7 获取密码策略
**接口**: `GET /api/v1/user/password/policy`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/user/password/policy"
```

**预期结果**:
- 返回密码策略信息
- 状态码: 200

### 5.8 检查账户状态
**接口**: `GET /api/v1/user/status/{username}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/user/status/Administrator"
```

**预期结果**:
- 返回账户状态信息
- 状态码: 200

### 5.9 获取用户会话
**接口**: `GET /api/v1/user/sessions/{username}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/user/sessions/Administrator"
```

**预期结果**:
- 返回用户会话列表
- 状态码: 200

### 5.10 结束用户会话
**接口**: `DELETE /api/v1/user/session/{sessionId}`

**测试步骤**:
```bash
curl -X DELETE "http://localhost:8081/api/v1/user/session/S-1-5-21-1000"
```

**预期结果**:
- 返回会话结束结果
- 状态码: 200

### 5.11 锁定用户账户
**接口**: `POST /api/v1/user/lock/{username}`

**测试步骤**:
```bash
curl -X POST "http://localhost:8081/api/v1/user/lock/testuser"
```

**预期结果**:
- 返回账户锁定结果
- 状态码: 200

### 5.12 解锁用户账户
**接口**: `POST /api/v1/user/unlock/{username}`

**测试步骤**:
```bash
curl -X POST "http://localhost:8081/api/v1/user/unlock/testuser"
```

**预期结果**:
- 返回账户解锁结果
- 状态码: 200

### 5.13 修改用户密码
**接口**: `POST /api/v1/user/password/change`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/user/password/change \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"testuser\",\"old_password\":\"old123\",\"new_password\":\"new123\"}"
```

**预期结果**:
- 返回密码修改结果
- 状态码: 200

### 5.14 获取用户组
**接口**: `GET /api/v1/user/groups/{username}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/user/groups/Administrator"
```

**预期结果**:
- 返回用户组信息
- 状态码: 200

### 5.15 获取用户历史记录
**接口**: `GET /api/v1/user/history/{username}`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/user/history/Administrator"
```

**预期结果**:
- 返回用户历史记录
- 状态码: 200

---

## 🛡️ 安全管理模块验证 (9个接口)

### 6.1 获取安全状态
**接口**: `GET /api/v1/security/status`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/security/status"
```

**预期结果**:
- 返回系统安全状态
- 状态码: 200

### 6.2 获取规则信息
**接口**: `GET /api/v1/security/rules`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/security/rules"
```

**预期结果**:
- 返回安全规则信息
- 状态码: 200

### 6.3 重新加载规则
**接口**: `POST /api/v1/security/reload-rules`

**测试步骤**:
```bash
curl -X POST "http://localhost:8081/api/v1/security/reload-rules"
```

**预期结果**:
- 返回规则重载结果
- 状态码: 200

### 6.4 清空缓存
**接口**: `POST /api/v1/security/cache/clear`

**测试步骤**:
```bash
curl -X POST "http://localhost:8081/api/v1/security/cache/clear"
```

**预期结果**:
- 返回缓存清空结果
- 状态码: 200

### 6.5 获取缓存统计
**接口**: `GET /api/v1/security/cache/stats`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/security/cache/stats"
```

**预期结果**:
- 返回缓存统计信息
- 状态码: 200

### 6.6 隔离文件
**接口**: `POST /api/v1/security/quarantine`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/security/quarantine \
  -H "Content-Type: application/json" \
  -d "{\"file_path\":\"C:\\\\temp\\\\malware.exe\"}"
```

**预期结果**:
- 返回文件隔离结果
- 状态码: 200

### 6.7 恢复文件
**接口**: `POST /api/v1/security/restore`

**测试步骤**:
```bash
curl -X POST http://localhost:8081/api/v1/security/restore \
  -H "Content-Type: application/json" \
  -d "{\"file_id\":\"quarantine-123\"}"
```

**预期结果**:
- 返回文件恢复结果
- 状态码: 200

### 6.8 获取隔离列表
**接口**: `GET /api/v1/security/quarantine/list`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/security/quarantine/list"
```

**预期结果**:
- 返回隔离文件列表
- 状态码: 200

### 6.9 获取扫描历史
**接口**: `GET /api/v1/security/scan-history`

**测试步骤**:
```bash
curl "http://localhost:8081/api/v1/security/scan-history"
```

**预期结果**:
- 返回扫描历史记录
- 状态码: 200

---

## 📊 系统监控模块验证 (3个接口)

### 7.1 健康检查
**接口**: `GET /api/health`

**测试步骤**:
```bash
curl "http://localhost:8081/api/health"
```

**预期结果**:
- 返回系统健康状态
- 状态码: 200

### 7.2 获取性能指标
**接口**: `GET /api/metrics`

**测试步骤**:
```bash
curl "http://localhost:8081/api/metrics"
```

**预期结果**:
- 返回系统性能指标
- 状态码: 200

### 7.3 重置性能指标
**接口**: `POST /api/metrics/reset`

**测试步骤**:
```bash
curl -X POST "http://localhost:8081/api/metrics/reset"
```

**预期结果**:
- 返回指标重置结果
- 状态码: 200

---

## 📝 验证要点

### 1. 接口可用性验证
- 所有83个接口都应该能够正常响应
- 检查HTTP状态码是否正确
- 验证返回的JSON数据格式

### 2. 功能完整性验证
- 文件管理功能：扫描、信息获取、操作等
- 进程管理功能：列表、操作、监控等
- 注册表操作：查询、创建、修改、删除等
- 网络管理：连接查看、监控等
- 用户管理：信息获取、权限检查、账户操作等
- 安全管理：状态查询、规则管理、隔离等
- 系统监控：健康检查、性能指标等

### 3. 错误处理验证
- 无效参数应返回400状态码
- 资源不存在应返回404状态码
- 服务器错误应返回500状态码

### 4. 性能验证
- 响应时间应在合理范围内
- 并发请求处理能力
- 内存和CPU使用情况

---

**文档版本**: 3.0  
**最后更新**: 2025-08-07  
**维护人员**: LYS 