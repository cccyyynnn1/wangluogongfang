# Yara安全服务需求验证指南

## 📋 目录索引

### 🔍 文件检测功能测试
- [1.1 PE文件载体检测](#11-pe文件载体检测)
- [1.2 脚本文件载体检测](#12-脚本文件载体检测)
- [1.3 文档文件载体检测](#13-文档文件载体检测)
- [1.4 压缩文件载体检测](#14-压缩文件载体检测)
- [1.5 Shellcode载荷检测](#15-shellcode载荷检测)
- [1.6 网络流载荷检测](#16-网络流载荷检测)
- [1.7 字符载荷检测](#17-字符载荷检测)
- [1.8 数据结构载荷检测](#18-数据结构载荷检测)
- [1.9 载体+载荷组合检测](#19-载体载荷组合检测)
- [1.10 多层载荷检测](#110-多层载荷检测)
- [1.11 深度文件分析测试](#111-深度文件分析测试)
- [1.12 文件属性信息获取测试](#112-文件属性信息获取测试)
- [1.13 目录遍历和文件加载测试](#113-目录遍历和文件加载测试)
- [1.14 文件管理操作测试](#114-文件管理操作测试)
- [1.15 目录扫描测试](#115-目录扫描测试)
- [1.16 文件操作测试](#116-文件操作测试)
- [1.17 内存缓冲区扫描测试](#117-内存缓冲区扫描测试)

### 👥 进程管理功能测试
- [2.1 进程信息获取测试](#21-进程信息获取测试)
- [2.2 进程操作功能测试](#22-进程操作功能测试)
- [2.3 系统模块查看功能测试](#23-系统模块查看功能测试)
- [2.4 系统模块操作功能测试](#24-系统模块操作功能测试)
- [2.5 进程监控功能测试](#25-进程监控功能测试)
- [2.6 进程监控测试](#26-进程监控测试)

### 🔧 注册表操作功能测试
- [3.1 注册表查询功能测试](#31-注册表查询功能测试)
- [3.2 注册表操作功能测试](#32-注册表操作功能测试)
- [3.3 注册表搜索功能测试](#33-注册表搜索功能测试)
- [3.4 注册表搜索测试](#34-注册表搜索测试)

### 🌐 网络连接功能测试
- [4.1 网络连接查看功能测试](#41-网络连接查看功能测试)
- [4.2 网络连接操作功能测试](#42-网络连接操作功能测试)
- [4.3 网络统计功能测试](#43-网络统计功能测试)
- [4.4 网络监控功能测试](#44-网络监控功能测试)
- [4.5 网络监控测试](#45-网络监控测试)

### 👤 用户权限功能测试
- [5.1 用户信息获取功能测试](#51-用户信息获取功能测试)
- [5.2 用户权限检查功能测试](#52-用户权限检查功能测试)
- [5.3 用户管理功能测试](#53-用户管理功能测试)
- [5.4 用户会话管理测试](#54-用户会话管理测试)

### 🛡️ 安全功能测试
- [6.1 安全状态查询测试](#61-安全状态查询测试)
- [6.2 安全规则管理测试](#62-安全规则管理测试)
- [6.3 文件隔离功能测试](#63-文件隔离功能测试)
- [6.4 文件隔离功能测试](#64-文件隔离功能测试)
- [6.5 安全规则管理测试](#65-安全规则管理测试)

### 📊 系统监控功能测试
- [7.1 健康检查测试](#71-健康检查测试)
- [7.2 性能指标测试](#72-性能指标测试)

---

## 📋 具体设计要求对照分析

本文档对照项目设计要求，详细分析每个对应要求的具体实现功能的测试验证流程。

### 设计要求概述
1. **病毒木马文件检测**：支持载体和载荷对象，包括PE文件、脚本文件、字符、数据结构、网络流
2. **文件系统功能**：支持读取指定路径文件的详细属性信息，包括文件大小、创建时间、修改时间等，支持底层逻辑对指定目录文件夹的文件和子文件进行遍历对文件加载
3. **进程管理功能**：支持遍历系统进程列表，获取进程信息，包括进程ID、启动时间、占用内存等。进程操作：实现对进程的启动、结束、挂起和恢复等操作
4. **注册表操作功能**：提供注册表项值的遍历、查询、新建、修改、删除和设置功能
5. **网络连接功能**：提供对TCP连接的操作，包括连接查看、关闭等功能
6. **系统模块操作功能**：提供对系统模块的操作，包括模块挂起和结束
7. **用户权限操作功能**：提供对用户权限的操作，包括检查账户密码相关属性值等

---

## 🔍 功能验证详细流程

### 1. 病毒木马文件检测功能验证

#### 1.1 PE文件载体检测

**测试目标**：验证系统能够检测PE文件载体对象

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 测试PE文件扫描（系统文件，安全）
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d "{\"path\":\"C:\\\\Windows\\\\System32\\\\notepad.exe\"}"

**预期结果：**
- 返回扫描结果，显示文件信息，威胁级别为"safe"

# 2. 测试PE文件属性信息获取
curl "http://localhost:8081/api/v1/file/info/C%3A%5CWindows%5CSystem32%5Cnotepad.exe"

**预期结果：**
- 返回文件详细信息，包括大小、创建时间、修改时间、哈希值等

# 3. 测试文件哈希计算
curl "http://localhost:8081/api/v1/file/hash/C%3A%5CWindows%5CSystem32%5Cnotepad.exe?algorithm=sha256"

**预期结果：**
- 返回文件的SHA256哈希值

# 4. 测试文件哈希验证
curl -X POST http://localhost:8081/api/v1/file/verify-hash -H "Content-Type: application/json" -d "{\"file_path\":\"C:\\\\Windows\\\\System32\\\\notepad.exe\",\"algorithm\":\"sha256\",\"expected_hash\":\"[实际哈希值]\"}"

**预期结果：**
- 返回验证结果（true/false）

# 5. 测试文件所有哈希值获取
curl "http://localhost:8081/api/v1/file/hashes/C%3A%5CWindows%5CSystem32%5Cnotepad.exe"

**预期结果：**
- 返回文件的所有哈希值（MD5、SHA1、SHA256等）
```

##### Linux版本测试
```bash
# 1. 测试PE文件扫描（Linux上测试Windows PE文件）
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d '{"path":"/tmp/test.exe"}'

**预期结果：**
- 返回扫描结果，显示PE文件结构信息

# 2. 测试文件信息获取
curl "http://localhost:8081/api/v1/file/info/%2Ftmp%2Ftest.exe"

**预期结果：**
- 返回文件详细信息

# 3. 测试文件哈希计算
curl "http://localhost:8081/api/v1/file/hash/%2Ftmp%2Ftest.exe?algorithm=md5"

**预期结果：**
- 返回文件的MD5哈希值

# 4. 测试文件所有哈希值获取
curl "http://localhost:8081/api/v1/file/hashes/%2Ftmp%2Ftest.exe"

**预期结果：**
- 返回文件的所有哈希值（MD5、SHA1、SHA256等）
```

#### 1.2 脚本文件载体检测
**测试目标**：验证系统能够检测各种脚本文件载体

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 测试PowerShell脚本扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d "{\"path\":\"C:\\\\temp\\\\test.ps1\"}"

**预期结果：**
- 识别为PowerShell脚本，分析脚本内容

# 2. 测试批处理脚本扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d "{\"path\":\"C:\\\\temp\\\\test.bat\"}"

**预期结果：**
- 识别为批处理脚本，分析命令内容

# 3. 测试VBS脚本扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d "{\"path\":\"C:\\\\temp\\\\test.vbs\"}"

**预期结果：**
- 识别为VBS脚本，分析脚本内容

# 4. 测试JavaScript脚本扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d "{\"path\":\"C:\\\\temp\\\\test.js\"}"

**预期结果：**
- 识别为JavaScript脚本，分析脚本内容
```

##### Linux版本测试
```bash
# 1. 测试Bash脚本扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d '{"path":"/tmp/test.sh"}'

**预期结果：**
- 识别为Bash脚本，分析脚本内容

# 2. 测试Python脚本扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d '{"path":"/tmp/test.py"}'

**预期结果：**
- 识别为Python脚本，分析脚本内容

# 3. 测试Perl脚本扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d '{"path":"/tmp/test.pl"}'

**预期结果：**
- 识别为Perl脚本，分析脚本内容
```

#### 1.3 文档文件载体检测
**测试目标**：验证系统能够检测文档文件载体

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 测试Word文档扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d "{\"path\":\"C:\\\\temp\\\\test.docx\"}"

**预期结果：**
- 识别为Word文档，分析文档结构和内容

# 2. 测试Excel文档扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d "{\"path\":\"C:\\\\temp\\\\test.xlsx\"}"

**预期结果：**
- 识别为Excel文档，分析表格数量

# 3. 测试PDF文档扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d "{\"path\":\"C:\\\\temp\\\\test.pdf\"}"

**预期结果：**
- 识别为PDF文档，分析文档内容
```

##### Linux版本测试
```bash
# 1. 测试OpenDocument格式扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d '{"path":"/tmp/test.odt"}'

**预期结果：**
- 识别为OpenDocument格式，分析文档内容

# 2. 测试纯文本文件扫描
curl -X POST http://localhost:8081/api/v1/file/scan -H "Content-Type: application/json" -d '{"path":"/tmp/test.txt"}'

**预期结果：**
- 识别为文本文件，分析文本内容
```

#### 1.4 压缩文件载体检测
**测试目标**：验证系统能够检测压缩文件载体

**测试步骤**：

##### Windows版本测试
```bash
# 1. 测试ZIP文件扫描
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\test.zip\"}"

**预期结果：**
- 识别为ZIP压缩文件，分析压缩包内容

# 2. 测试RAR文件扫描
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\test.rar\"}"

**预期结果：**
- 识别为RAR压缩文件，分析压缩包内容

# 3. 测试7Z文件扫描
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\test.7z\"}"

**预期结果：**
- 识别为7Z压缩文件，分析压缩包内容
```

##### Linux版本测试
```bash
# 1. 测试TAR文件扫描
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d '{"path":"/tmp/test.tar.gz"}'

**预期结果：**
- 识别为TAR压缩文件，分析压缩包内容

# 2. 测试GZIP文件扫描
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d '{"path":"/tmp/test.gz"}'

**预期结果：**
- 识别为GZIP压缩文件，分析压缩包内容
```

#### 1.5 Shellcode载荷检测
**测试目标**：验证系统能够检测各种Shellcode载荷

**测试步骤**：

##### Windows版本测试
```bash
# 1. 测试Windows Shellcode检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"windows_shellcode\",\"data\":\"\\x68\\x89\\xe5\\x83\\xec\\x64\\xa1\\x30\\x00\\x00\\x00\\x8b\\x40\\x0c\\x8b\\x70\\x14\\xad\\x96\"}"

**预期结果：**
- 检测到Windows Shellcode特征，威胁级别为"high"

# 2. 测试编码Shellcode检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"encoded_shellcode\",\"data\":\"\\x31\\xc9\\x66\\x81\\xca\\xff\\x0f\\x42\\x52\\x6a\\x02\\x58\\xcd\\x2e\\x3c\\x05\\x74\\xef\\xb8\"}"

**预期结果：**
- 检测到编码Shellcode特征，威胁级别为"medium"
```

##### Linux版本测试
```bash
# 1. 测试Linux Shellcode检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d '{"identifier":"linux_shellcode","data":"\\x90\\x90\\x90\\x31\\xc0\\x50\\x68\\x2f\\x2f\\x73\\x68\\x68\\x2f\\x62\\x69\\x6e\\x89\\xe3\\x50\\x53\\x89\\xe1\\xb0\\x0b\\xcd\\x80"}'

**预期结果：**
- 检测到Linux Shellcode特征，威胁级别为"high"

# 2. 测试NOP sled检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d '{"identifier":"nop_sled","data":"\\x90\\x90\\x90\\x90\\x90\\x90\\x90\\x90\\x90\\x90\\x90\\x90\\x90\\x90\\x90\\x90"}'

**预期结果：**
- 检测到NOP sled特征，威胁级别为"medium"
```

#### 1.6 网络流载荷检测
**测试目标**：验证系统能够检测网络流载荷

**测试步骤**：

##### Windows版本测试
```bash
# 1. 测试URL载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"malicious_url\",\"data\":\"http://malicious-site.com/download.exe\"}"

**预期结果：**
- 检测到恶意URL特征，威胁级别为"medium"

# 2. 测试HTTP请求载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"http_payload\",\"data\":\"GET /download.php HTTP/1.1\\r\\nHost: evil.com\\r\\nUser-Agent: Mozilla/5.0\\r\\n\\r\\n\"}"

**预期结果：**
- 检测到HTTP请求特征，威胁级别为"medium"

# 3. 测试FTP载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"ftp_payload\",\"data\":\"ftp://user:pass@evil.com/malware.exe\"}"

**预期结果：**
- 检测到FTP连接特征，威胁级别为"medium"
```

##### Linux版本测试
```bash
# 1. 测试SSH连接载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d '{"identifier":"ssh_payload","data":"ssh://user@evil.com:22"}'

**预期结果：**
- 检测到SSH连接特征，威胁级别为"medium"

# 2. 测试DNS查询载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d '{"identifier":"dns_payload","data":"evil.com A"}'

**预期结果：**
- 检测到DNS查询特征，威胁级别为"low"
```

#### 1.7 字符载荷检测
**测试目标**：验证系统能够检测恶意字符载荷

**测试步骤**：

##### Windows版本测试
```bash
# 1. 测试恶意字符串载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"malicious_string\",\"data\":\"CreateRemoteThread VirtualAllocEx WriteProcessMemory\"}"

**预期结果：**
- 检测到恶意API调用特征，威胁级别为"high"

# 2. 测试编码字符串载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"encoded_string\",\"data\":\"powershell -EncodedCommand UwBlAHQALQBJAHQAZQBtAFAAcgBvAHAAZQByAHQAeQAgAC0AUABhAHQAaAAgAEgASwBDAFUAOgBcAFMATwBGAFQAVwBBAFIARQBcAE0AaQBjAHIAbwBzAG8AZgB0AFwAVwBpAG4AZABvAHcAcwBcAEMAdQByAHIAZQBuAHQAVgBlAHIAcwBpAG8AbgBcAFIAdQBuACAALQBOAGEAbQBlACAAIgB0AGUAcwB0ACIAIAAtAFYAYQBsAHUAZQAgACIAYwBtAGQALgBlAHgAZQAiAA==\"}"

**预期结果：**
- 检测到编码PowerShell命令特征，威胁级别为"high"

# 3. 测试混淆字符串载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"obfuscated_string\",\"data\":\"\\x43\\x72\\x65\\x61\\x74\\x65\\x52\\x65\\x6d\\x6f\\x74\\x65\\x54\\x68\\x72\\x65\\x61\\x64\"}"

**预期结果：**
- 检测到混淆字符串特征，威胁级别为"medium"
```

##### Linux版本测试
```bash
# 1. 测试Linux恶意命令检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d '{"identifier":"linux_malicious_command","data":"rm -rf / && nc -l 4444"}'

**预期结果：**
- 检测到恶意Linux命令特征，威胁级别为"high"

# 2. 测试Base64编码载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d '{"identifier":"base64_payload","data":"cm0gLXJmIC8gJiYgbmMgLWwgNDQ0NA=="}'

**预期结果：**
- 检测到Base64编码特征，威胁级别为"medium"
```

#### 1.8 数据结构载荷检测
**测试目标**：验证系统能够检测恶意数据结构载荷

**测试步骤**：

##### Windows版本测试
```bash
# 1. 测试十六进制数据结构检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"hex_data\",\"data\":\"4D5A900050450000\"}"

**预期结果：**
- 检测到PE文件头特征，威胁级别为"medium"

# 2. 测试Base64编码数据结构检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"base64_data\",\"data\":\"TVqQAAMAAAAEAAAA//8AALgAAAAAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAA4fug4AtAnNIbgBTM0hVGhpcyBwcm9ncmFtIGNhbm5vdCBiZSBydW4gaW4gRE9TIG1vZGUuDQ0KJAAAAAAAAABQRQAATAEDAKJqXFYAAAAAAAAAAOAAIiALATAAAA4AAAAGAAAAAAAAVjIAAAAgAAAAQAAAAAAAEAAgAAAAAgAABAAAAAAAAAAGAAAAAAAAAACAAAAAAgAAAAAAAAMAYIUAABAAABAAAAAAEAAAEAAAAAAAABAAAAAAAAAAAAAAAGAyAABPAAAAAEAAAIgDAAAAAAAAAAAAAAAAAAAAAAAAAGAAAAwAAADcMAAAOAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==\"}"

**预期结果：**
- 检测到Base64编码PE文件特征，威胁级别为"medium"

# 3. 测试二进制数据结构检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"binary_data\",\"data\":\"\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\\x00\"}"

**预期结果：**
- 检测到二进制数据特征，威胁级别为"low"
```

##### Linux版本测试
```bash
# 1. 测试ELF文件头检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d '{"identifier":"elf_header","data":"\\x7f\\x45\\x4c\\x46\\x02\\x01\\x01\\x00"}'

**预期结果：**
- 检测到ELF文件头特征，威胁级别为"medium"

# 2. 测试Shell脚本头检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d '{"identifier":"shell_header","data":"#!/bin/bash"}'

**预期结果：**
- 检测到Shell脚本头特征，威胁级别为"low"
```

#### 1.9 载体+载荷组合检测
**测试目标**：验证系统能够检测载体和载荷的组合威胁

**测试步骤**：

##### Windows版本测试
```bash
# 1. 测试PE文件中的Shellcode载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\shellcode_loader.exe\"}"

**预期结果：**
- 检测到PE文件载体和Shellcode载荷的组合威胁，威胁级别为"high"

# 2. 测试脚本文件中的网络载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\downloader.ps1\"}"

**预期结果：**
- 检测到PowerShell脚本载体和网络下载载荷的组合威胁，威胁级别为"high"

# 3. 测试文档文件中的宏载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\macro_document.docm\"}"

**预期结果：**
- 检测到Word文档载体和宏载荷的组合威胁，威胁级别为"high"
```

##### Linux版本测试
```bash
# 1. 测试ELF文件中的Shellcode载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d '{"path":"/tmp/shellcode_loader"}'

**预期结果：**
- 检测到ELF文件载体和Shellcode载荷的组合威胁，威胁级别为"high"

# 2. 测试Shell脚本中的命令载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d '{"path":"/tmp/malicious_script.sh"}'

**预期结果：**
- 检测到Shell脚本载体和恶意命令载荷的组合威胁，威胁级别为"high"
```

#### 1.10 多层载荷检测
**测试目标**：验证系统能够检测多层嵌套的载荷

**测试步骤**：

##### Windows版本测试
```bash
# 1. 测试压缩文件中的多层载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\nested_payload.zip\"}"

**预期结果：**
- 检测到ZIP压缩文件中的多层载荷，威胁级别为"high"

# 2. 测试编码载荷的嵌套检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d "{\"identifier\":\"nested_payload\",\"data\":\"UEsDBBQAAAAIAA==\"}"

**预期结果：**
- 检测到Base64编码的ZIP文件载荷，威胁级别为"high"
```

##### Linux版本测试
```bash
# 1. 测试TAR文件中的多层载荷检测
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d '{"path":"/tmp/nested_payload.tar.gz"}'

**预期结果：**
- 检测到TAR压缩文件中的多层载荷，威胁级别为"high"

# 2. 测试编码载荷的嵌套检测
curl -X POST http://localhost:8081/api/v1/file/scan-buffer \
  -H "Content-Type: application/json" \
  -d '{"identifier":"nested_payload","data":"H4sIAAAAAAAA"}'

**预期结果：**
- 检测到Base64编码的GZIP文件载荷，威胁级别为"high"
```

#### 1.11 深度文件分析测试

**测试目标**：验证深度文件分析功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 测试深度文件分析
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\suspicious.exe\"}"

**预期结果：**
- 返回详细的文件分析结果，包括文件结构、内容分析、威胁评估等

# 2. 测试文件内容分析
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\suspicious.ps1\"}"

**预期结果：**
- 返回脚本内容分析结果，包括代码结构、函数调用、字符串分析等

# 3. 测试文件结构分析
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\\\\suspicious.zip\"}"

**预期结果：**
- 返回压缩文件结构分析结果，包括文件列表、压缩算法、文件大小等
```

##### Linux版本测试
```bash
# 1. 测试ELF文件深度分析
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d '{"path":"/tmp/suspicious_binary"}'

**预期结果：**
- 返回ELF文件详细分析结果，包括节头、符号表、依赖库等

# 2. 测试脚本文件深度分析
curl -X POST http://localhost:8081/api/v1/file/scan \
  -H "Content-Type: application/json" \
  -d '{"path":"/tmp/suspicious_script.sh"}'

**预期结果：**
- 返回脚本文件详细分析结果，包括命令分析、变量使用、函数定义等
```

#### 1.12 文件属性信息获取测试
**测试目标**：验证能够获取文件的详细属性信息

**测试步骤**：

##### Windows版本测试
```bash
# 1. 获取文件详细信息
curl "http://localhost:8081/api/v1/file/info/C%3A%5CWindows%5CSystem32%5Cnotepad.exe"

**预期结果：**
- 返回文件详细信息，包括：
# - 文件大小、创建时间、修改时间、访问时间
# - 文件权限、所有者、组信息
# - 文件类型检测和威胁级别评估
# - 文件哈希值（MD5、SHA256）

# 2. 获取文件哈希值
curl "http://localhost:8081/api/v1/file/hash/C%3A%5CWindows%5CSystem32%5Cnotepad.exe?algorithm=sha256"

**预期结果：**
- 返回指定算法的文件哈希值

# 3. 获取文件所有哈希值
curl "http://localhost:8081/api/v1/file/hashes/C%3A%5CWindows%5CSystem32%5Cnotepad.exe"

**预期结果：**
- 返回文件的所有哈希值（MD5、SHA1、SHA256等）

# 4. 验证文件哈希
curl -X POST http://localhost:8081/api/v1/file/verify-hash \
  -H "Content-Type: application/json" \
  -d "{\"file_path\":\"C:\\\\Windows\\\\System32\\\\notepad.exe\",\"algorithm\":\"sha256\",\"expected_hash\":\"[实际哈希值]\"}"

**预期结果：**
- 返回哈希验证结果（true/false）
```

##### Linux版本测试
```bash
# 1. 获取文件详细信息
curl "http://localhost:8081/api/v1/file/info/%2Fusr%2Fbin%2Fls"

**预期结果：**
- 返回文件详细信息，包括：
# - 文件大小、创建时间、修改时间、访问时间
# - 文件权限、所有者、组信息
# - 文件类型检测和威胁级别评估
# - 文件哈希值（MD5、SHA256）

# 2. 获取文件哈希值
curl "http://localhost:8081/api/v1/file/hash/%2Fusr%2Fbin%2Fls?algorithm=md5"

**预期结果：**
- 返回指定算法的文件哈希值

# 3. 获取文件所有哈希值
curl "http://localhost:8081/api/v1/file/hashes/%2Fusr%2Fbin%2Fls"

**预期结果：**
- 返回文件的所有哈希值（MD5、SHA1、SHA256等）

# 4. 验证文件哈希
curl -X POST http://localhost:8081/api/v1/file/verify-hash \
  -H "Content-Type: application/json" \
  -d '{"file_path":"/usr/bin/ls","algorithm":"md5","expected_hash":"[实际哈希值]"}'

**预期结果：**
- 返回哈希验证结果（true/false）
```

#### 1.13 目录遍历和文件加载测试
**测试目标**：验证底层逻辑对指定目录文件夹的文件和子文件进行遍历

**测试步骤**：

##### Windows版本测试
```bash
# 1. 测试目录扫描（递归扫描）
curl -X POST http://localhost:8081/api/v1/file/scan-directory \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\temp\",\"recursive\":true,\"max_depth\":3}"

**预期结果：**
- 返回扫描结果，包含：
# - 扫描的文件总数
# - 发现的可疑文件列表
# - 扫描耗时
# - 威胁统计信息

# 2. 测试特定目录深度扫描
curl -X POST http://localhost:8081/api/v1/file/scan-directory \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\Users\",\"recursive\":true,\"max_depth\":2}"

**预期结果：**
- 返回指定深度的目录扫描结果

# 3. 测试系统目录扫描（非递归）
curl -X POST http://localhost:8081/api/v1/file/scan-directory \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"C:\\\\Windows\\\\System32\",\"recursive\":false,\"max_depth\":1}"

**预期结果：**
- 返回系统目录的直接文件扫描结果

# 4. 获取文件列表
curl "http://localhost:8081/api/v1/file/list?path=C%3A%5Ctemp"

**预期结果：**
- 返回指定目录的文件列表，包括文件名、大小、类型等信息
```

##### Linux版本测试
```bash
# 1. 测试目录扫描（递归扫描）
curl -X POST http://localhost:8081/api/v1/file/scan-directory \
  -H "Content-Type: application/json" \
  -d '{"path":"/tmp","recursive":true,"max_depth":3}'

**预期结果：**
- 返回扫描结果，包含：
# - 扫描的文件总数
# - 发现的可疑文件列表
# - 扫描耗时
# - 威胁统计信息

# 2. 测试特定目录深度扫描
curl -X POST http://localhost:8081/api/v1/file/scan-directory \
  -H "Content-Type: application/json" \
  -d '{"path":"/home","recursive":true,"max_depth":2}'

**预期结果：**
- 返回指定深度的目录扫描结果

# 3. 测试系统目录扫描（非递归）
curl -X POST http://localhost:8081/api/v1/file/scan-directory \
  -H "Content-Type: application/json" \
  -d '{"path":"/usr/bin","recursive":false,"max_depth":1}'

**预期结果：**
- 返回系统目录的直接文件扫描结果

# 4. 获取文件列表
curl "http://localhost:8081/api/v1/file/list?path=%2Ftmp"

**预期结果：**
- 返回指定目录的文件列表，包括文件名、大小、类型等信息
```

#### 1.14 文件管理操作测试

**测试目标**：验证文件管理操作功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 测试文件复制操作
curl -X POST http://localhost:8081/api/v1/file/copy \
  -H "Content-Type: application/json" \
  -d "{\"source\":\"C:\\\\temp\\\\test.txt\",\"dest\":\"C:\\\\temp\\\\test_copy.txt\"}"

**预期结果：**
- 文件复制成功，返回成功消息

# 2. 测试文件移动操作
curl -X POST http://localhost:8081/api/v1/file/move \
  -H "Content-Type: application/json" \
  -d "{\"source\":\"C:\\\\temp\\\\test_copy.txt\",\"dest\":\"C:\\\\temp\\\\test_moved.txt\"}"

**预期结果：**
- 文件移动成功，返回成功消息

# 3. 测试文件删除操作
curl -X DELETE "http://localhost:8081/api/v1/file/C%3A%5Ctemp%5Ctest_moved.txt"

**预期结果：**
- 文件删除成功，返回成功消息
```

##### Linux版本测试
```bash
# 1. 测试文件复制操作
curl -X POST http://localhost:8081/api/v1/file/copy \
  -H "Content-Type: application/json" \
  -d '{"source":"/tmp/test.txt","dest":"/tmp/test_copy.txt"}'

**预期结果：**
- 文件复制成功，返回成功消息

# 2. 测试文件移动操作
curl -X POST http://localhost:8081/api/v1/file/move \
  -H "Content-Type: application/json" \
  -d '{"source":"/tmp/test_copy.txt","dest":"/tmp/test_moved.txt"}'

**预期结果：**
- 文件移动成功，返回成功消息

# 3. 测试文件删除操作
curl -X DELETE "http://localhost:8081/api/v1/file/%2Ftmp%2Ftest_moved.txt"

**预期结果：**
- 文件删除成功，返回成功消息
```

---

### 2. 进程管理功能验证

#### 2.1 进程信息获取测试

**测试目标**：验证能够获取系统进程列表和详细信息

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 获取进程列表
curl http://localhost:8081/api/v1/process/list

**预期结果：**
- 返回系统进程列表，包含：
# - 进程ID、进程名、CPU使用率、内存使用率
# - 启动时间、状态、优先级
# - 父进程ID、子进程数量

# 2. 获取特定进程信息
curl http://localhost:8081/api/v1/process/28100

**预期结果：**
- 返回指定进程的详细信息，包括：
# - 基本进程信息（PID、名称、状态等）
# - 内存使用情况（物理内存、虚拟内存、内存百分比）
# - CPU使用情况（CPU时间、CPU百分比）
# - 文件信息（打开的文件数量）
# - 网络连接信息
# - 模块信息（加载的DLL）
# - 线程信息
# - 句柄信息

# 3. 获取进程统计信息
curl http://localhost:8081/api/v1/process/statistics

**预期结果：**
- 返回进程统计信息，包括：
# - 总进程数
# - 运行中进程数
# - 系统进程数
# - 用户进程数
# - 平均CPU使用率
# - 平均内存使用率

# 4. 获取进程模块信息
curl http://localhost:8081/api/v1/process/28100/modules

**预期结果：**
- 返回进程加载的模块列表，包括：
# - 模块名称、路径、大小
# - 加载地址、基址
# - 版本信息

# 5. 获取进程连接信息
curl http://localhost:8081/api/v1/process/28100/connections

**预期结果：**
- 返回进程的网络连接信息，包括：
# - 本地地址和端口
# - 远程地址和端口
# - 连接状态、类

# 6. 获取进程内存信息
curl http://localhost:8081/api/v1/process/28100/memory

**预期结果：**
- 返回进程的内存使用详情，包括：
# - 物理内存使用情况
# - 虚拟内存使用情况
# - 内存使用百分比
# - 页面文件使用情况

# 7. 检查进程是否运行
curl http://localhost:8081/api/v1/process/28100/running

**预期结果：**
- 返回进程运行状态（true/false）

# 8. 获取进程子进程
curl http://localhost:8081/api/v1/process/28100/children

**预期结果：**
- 返回进程的子进程列表
```

##### Linux版本测试
```bash
# 1. 获取进程列表
curl http://localhost:8081/api/v1/process/list

**预期结果：**
- 返回系统进程列表，包含：
# - 进程ID、进程名、CPU使用率、内存使用率
# - 启动时间、状态、优先级
# - 父进程ID、子进程数量

# 2. 获取特定进程信息
curl http://localhost:8081/api/v1/process/1234

**预期结果：**
- 返回指定进程的详细信息，包括：
# - 基本进程信息（PID、名称、状态等）
# - 内存使用情况（物理内存、虚拟内存、内存百分比）
# - CPU使用情况（CPU时间、CPU百分比）
# - 文件信息（打开的文件数量）
# - 网络连接信息
# - 模块信息（加载的共享库）
# - 线程信息
# - 句柄信息

# 3. 获取进程统计信息
curl http://localhost:8081/api/v1/process/statistics

**预期结果：**
- 返回进程统计信息，包括：
# - 总进程数
# - 运行中进程数
# - 系统进程数
# - 用户进程数
# - 平均CPU使用率
# - 平均内存使用率

# 4. 获取进程模块信息
curl http://localhost:8081/api/v1/process/1234/modules

**预期结果：**
- 返回进程加载的模块列表，包括：
# - 模块名称、路径、大小
# - 加载地址、基址
# - 版本信息

# 5. 获取进程连接信息
curl http://localhost:8081/api/v1/process/1234/connections

**预期结果：**
- 返回进程的网络连接信息，包括：
# - 本地地址和端口
# - 远程地址和端口
# - 连接状态、类

# 6. 获取进程内存信息
curl http://localhost:8081/api/v1/process/1234/memory

**预期结果：**
- 返回进程的内存使用详情，包括：
# - 物理内存使用情况
# - 虚拟内存使用情况
# - 内存使用百分比
# - 共享内存使用情况

# 7. 检查进程是否运行
curl http://localhost:8081/api/v1/process/1234/running

**预期结果：**
- 返回进程运行状态（true/false）

# 8. 获取进程子进程
curl http://localhost:8081/api/v1/process/1234/children

**预期结果：**
- 返回进程的子进程列表
```

#### 2.2 进程操作功能测试

**测试目标**：验证进程的启动、结束、挂起和恢复功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 启动进程（安全测试）
curl -X POST http://localhost:8081/api/v1/process/start \
  -H "Content-Type: application/json" \
  -d "{\"command\":\"notepad.exe\",\"args\":[],\"working_dir\":\"\"}"

**预期结果：**
- 进程启动成功，返回进程信息，包括：
# - 进程ID、进程名
# - 启动时间、状态
# - 父进程ID

# 2. 启动进程（带参数）
curl -X POST http://localhost:8081/api/v1/process/start -H "Content-Type: application/json" -d "{\"command\":\"notepad.exe\",\"args\":[],\"working_dir\":\"C:\\\\Windows\\\\System32\"}"

**预期结果：**
- 进程启动成功，执行指定命令

# 3. 挂起进程
curl -X PUT http://localhost:8081/api/v1/process/17704/suspend

**预期结果：**
- 进程挂起成功，返回操作状态，注意：挂起后进程状态变为suspended

# 4. 恢复进程
curl -X PUT http://localhost:8081/api/v1/process/17704/resume

**预期结果：**
- 进程恢复成功，返回操作状态，注意：恢复后进程状态变为running

# 5. 结束进程
curl -X DELETE http://localhost:8081/api/v1/process/17704

**预期结果：**
- 进程结束成功，返回操作状态，注意：结束后进程从列表中消失
```

##### Linux版本测试
```bash
# 1. 启动进程（安全测试）
curl -X POST http://localhost:8081/api/v1/process/start \
  -H "Content-Type: application/json" \
  -d '{"command":"/usr/bin/nano","args":[],"working_dir":""}'

**预期结果：**
- 进程启动成功，返回进程信息，包括：
# - 进程ID、进程名
# - 启动时间、状态
# - 父进程ID

# 2. 启动进程（带参数）
curl -X POST http://localhost:8081/api/v1/process/start \
  -H "Content-Type: application/json" \
  -d '{"command":"/bin/echo","args":["test"],"working_dir":"/tmp"}'

**预期结果：**
- 进程启动成功，执行指定命令

# 3. 挂起进程
curl -X PUT http://localhost:8081/api/v1/process/1234/suspend

**预期结果：**
- 进程挂起成功，返回操作状态，注意：挂起后进程状态变为suspended

# 4. 恢复进程
curl -X PUT http://localhost:8081/api/v1/process/1234/resume

**预期结果：**
- 进程恢复成功，返回操作状态，注意：恢复后进程状态变为running

# 5. 结束进程
curl -X DELETE http://localhost:8081/api/v1/process/1234

**预期结果：**
- 进程结束成功，返回操作状态，注意：结束后进程从列表中消失
```

#### 2.3 系统模块查看功能测试

**测试目标**：验证系统模块的查看功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 获取系统模块列表
curl http://localhost:8081/api/v1/process/modules

**预期结果：**
- 返回系统模块列表，包含：
# - 模块名称、路径、大小
# - 加载地址、基址
# - 版本信息、描述
# - 使用该模块的进程数量

# 2. 获取特定模块信息
curl http://localhost:8081/api/v1/process/module/ntdll.dll

**预期结果：**
- 返回指定模块的详细信息，包括：
# - 模块基本信息（名称、路径、大小等）
# - 加载信息（基址、入口点等）
# - 版本信息（文件版本、产品版本等）
# - 使用该模块的进程列表
# - 模块依赖关系

# 3. 获取常用系统模块信息
curl http://localhost:8081/api/v1/process/module/kernel32.dll

**预期结果：**
- 返回kernel32.dll模块的详细信息

# 4. 获取用户程序模块信息
curl http://localhost:8081/api/v1/process/module/WeChat.exe

**预期结果：**
- 返回WeChat.exe模块的详细信息
```

##### Linux版本测试
```bash
# 1. 获取系统模块列表
curl http://localhost:8081/api/v1/process/modules

**预期结果：**
- 返回系统模块列表，包含：
# - 模块名称、路径、大小
# - 加载地址、基址
# - 版本信息、描述
# - 使用该模块的进程数量

# 2. 获取特定模块信息
curl http://localhost:8081/api/v1/process/module/libc.so.6

**预期结果：**
- 返回指定模块的详细信息，包括：
# - 模块基本信息（名称、路径、大小等）
# - 加载信息（基址、入口点等）
# - 版本信息（文件版本、产品版本等）
# - 使用该模块的进程列表
# - 模块依赖关系

# 3. 获取常用系统模块信息
curl http://localhost:8081/api/v1/process/module/libpthread.so.0

**预期结果：**
- 返回libpthread.so.0模块的详细信息

# 4. 获取用户程序模块信息
curl http://localhost:8081/api/v1/process/module/bash

**预期结果：**
- 返回bash模块的详细信息
```

#### 2.4 系统模块操作功能测试

**测试目标**：验证系统模块的挂起和结束功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 挂起系统模块（谨慎操作）
curl -X PUT http://localhost:8081/api/v1/process/module/WeChat.exe/suspend

**预期结果：**
- 模块挂起成功，返回操作状态，注意：挂起后使用该模块的进程可能受到影响

# 2. 恢复系统模块
curl -X PUT http://localhost:8081/api/v1/process/module/WeChat.exe/resume

**预期结果：**
- 模块恢复成功，返回操作状态，注意：恢复后相关进程恢复正常运行

# 3. 结束系统模块（谨慎操作）
curl -X DELETE http://localhost:8081/api/v1/process/module/WeChat.exe

**预期结果：**
- 模块结束成功，返回操作状态，注意：结束后使用该模块的进程可能被终止

# 4. 挂起非关键模块（安全测试）
curl -X PUT http://localhost:8081/api/v1/process/module/test_module.dll/suspend

**预期结果：**
- 模块挂起成功，返回操作状态

# 5. 恢复非关键模块
curl -X PUT http://localhost:8081/api/v1/process/module/test_module.dll/resume

**预期结果：**
- 模块恢复成功，返回操作状态
```

##### Linux版本测试
```bash
# 1. 挂起系统模块（谨慎操作）
curl -X PUT http://localhost:8081/api/v1/process/module/bash/suspend

**预期结果：**
- 模块挂起成功，返回操作状态，注意：挂起后使用该模块的进程可能受到影响

# 2. 恢复系统模块
curl -X PUT http://localhost:8081/api/v1/process/module/bash/resume

**预期结果：**
- 模块恢复成功，返回操作状态，注意：恢复后相关进程恢复正常运行

# 3. 结束系统模块（谨慎操作）
curl -X DELETE http://localhost:8081/api/v1/process/module/bash

**预期结果：**
- 模块结束成功，返回操作状态，注意：结束后使用该模块的进程可能被终止

# 4. 挂起非关键模块（安全测试）
curl -X PUT http://localhost:8081/api/v1/process/module/test_module.so/suspend

**预期结果：**
- 模块挂起成功，返回操作状态

# 5. 恢复非关键模块
curl -X PUT http://localhost:8081/api/v1/process/module/test_module.so/resume

**预期结果：**
- 模块恢复成功，返回操作状态
```

#### 2.5 进程监控功能测试

**测试目标**：验证进程监控功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 启用进程监控
curl -X POST http://localhost:8081/api/v1/process/monitoring/enable

**预期结果：**
- 进程监控启用成功，返回成功消息

# 2. 获取监控的进程列表
curl http://localhost:8081/api/v1/process/monitoring/list

**预期结果：**
- 返回当前监控的进程列表，包括：
# - 进程ID、进程名
# - 监控开始时间
# - 监控状态
# - 异常行为记录

# 3. 禁用进程监控
curl -X POST http://localhost:8081/api/v1/process/monitoring/disable

**预期结果：**
- 进程监控禁用成功，返回成功消息

# 4. 再次获取监控的进程列表
curl http://localhost:8081/api/v1/process/monitoring/list

**预期结果：**
- 返回空列表或监控已停止的消?
```

##### Linux版本测试
```bash
# 1. 启用进程监控
curl -X POST http://localhost:8081/api/v1/process/monitoring/enable

**预期结果：**
- 进程监控启用成功，返回成功消息

# 2. 获取监控的进程列表
curl http://localhost:8081/api/v1/process/monitoring/list

**预期结果：**
- 返回当前监控的进程列表，包括：
# - 进程ID、进程名
# - 监控开始时间
# - 监控状态
# - 异常行为记录

# 3. 禁用进程监控
curl -X POST http://localhost:8081/api/v1/process/monitoring/disable

**预期结果：**
- 进程监控禁用成功，返回成功消息

# 4. 再次获取监控的进程列表
curl http://localhost:8081/api/v1/process/monitoring/list

**预期结果：**
- 返回空列表或监控已停止的消?
```

---

### 3. 注册表操作功能验证

#### 3.1 注册表查询功能测试

**测试目标**：验证注册表项的遍历、查询功能

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 获取注册表键信息
curl "http://localhost:8081/api/v1/registry/key/SOFTWARE%5CMicrosoft%5CWindows%5CCurrentVersion"

**预期结果：**
- 返回注册表键信息，包括：
# - 键路径、键名称
# - 子键数量、值数
# - 最后修改时间
# - 访问权限

# 2. 获取注册表值
curl "http://localhost:8081/api/v1/registry/value/SOFTWARE%5CMicrosoft%5CWindows%5CCurrentVersion/ProgramFilesDir"

**预期结果：**
- 返回注册表值信息，包括：
# - 值名称、值类型
# - 值数据、数据大小
# - 最后修改时间

# 3. 列出注册表子键
curl "http://localhost:8081/api/v1/registry/keys/SOFTWARE%5CMicrosoft%5CWindows%5CCurrentVersion"

**预期结果：**
- 返回子键列表，包括：
# - 子键名称、子键路径
# - 子键数量、值数
# - 最后修改时间

# 4. 列出注册表值
curl "http://localhost:8081/api/v1/registry/values/SOFTWARE%5CMicrosoft%5CWindows%5CCurrentVersion"

**预期结果：**
- 返回值列表，包括：
# - 值名称、值类型
# - 值数据、数据大小
# - 最后修改时间

# 5. 搜索注册表
curl "http://localhost:8081/api/v1/registry/search?root_path=SOFTWARE&search_term=Windows"

**预期结果：**
- 返回搜索结果，包括：
# - 匹配的键路径
# - 匹配的值名称
# - 匹配的值数据
# - 搜索耗时

# 6. 搜索特定注册表项
curl "http://localhost:8081/api/v1/registry/search?root_path=SOFTWARE&search_term=Microsoft"

**预期结果：**
- 返回包含"Microsoft"的注册表项列表
```

##### Linux版本测试
```bash
# 注意：Linux没有注册表，这些API在Linux上可能返回错误或空结果
# 1. 获取注册表键信息（Linux上不可用）
curl "http://localhost:8081/api/v1/registry/key/SOFTWARE%5CMicrosoft%5CWindows%5CCurrentVersion"

**预期结果：**
- 返回错误信息或空结果

# 2. 获取注册表值（Linux上不可用）
curl "http://localhost:8081/api/v1/registry/value/SOFTWARE%5CMicrosoft%5CWindows%5CCurrentVersion/ProgramFilesDir"

**预期结果：**
- 返回错误信息或空结果

# 3. 列出注册表子键（Linux上不可用）
curl "http://localhost:8081/api/v1/registry/keys/SOFTWARE%5CMicrosoft%5CWindows%5CCurrentVersion"

**预期结果：**
- 返回错误信息或空结果

# 4. 列出注册表值（Linux上不可用）
curl "http://localhost:8081/api/v1/registry/values/SOFTWARE%5CMicrosoft%5CWindows%5CCurrentVersion"

**预期结果：**
- 返回错误信息或空结果

# 5. 搜索注册表（Linux上不可用）
curl "http://localhost:8081/api/v1/registry/search?root_path=SOFTWARE&search_term=Windows"

**预期结果：**
- 返回错误信息或空结果
```

#### 3.2 注册表操作功能测试

**测试目标**：验证注册表的新建、修改、删除功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 创建注册表键
curl -X POST http://localhost:8081/api/v1/registry/key \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"SOFTWARE\\\\TestKey\"}"

**预期结果：**
- 注册表键创建成功，返回成功消息

# 2. 设置字符串注册表值
curl -X PUT http://localhost:8081/api/v1/registry/value \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"SOFTWARE\\\\TestKey\",\"name\":\"TestString\",\"type\":\"string\",\"value\":\"test_value\"}"

**预期结果：**
- 字符串值设置成功，返回成功消息

# 3. 设置DWORD注册表值
curl -X PUT http://localhost:8081/api/v1/registry/value \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"SOFTWARE\\\\TestKey\",\"name\":\"TestDWORD\",\"type\":\"dword\",\"value\":\"12345\"}"

**预期结果：**
- DWORD值设置成功，返回成功消息

# 4. 设置二进制注册表值
curl -X PUT http://localhost:8081/api/v1/registry/value \
  -H "Content-Type: application/json" \
  -d "{\"path\":\"SOFTWARE\\\\TestKey\",\"name\":\"TestBinary\",\"type\":\"binary\",\"value\":\"0102030405\"}"

**预期结果：**
- 二进制值设置成功，返回成功消息

# 5. 验证设置的注册表值
curl "http://localhost:8081/api/v1/registry/value/SOFTWARE%5CTestKey/TestString"

**预期结果：**
- 返回设置的字符串值

# 6. 删除注册表值
curl -X DELETE "http://localhost:8081/api/v1/registry/value/SOFTWARE%5CTestKey/TestString"

**预期结果：**
- 注册表值删除成功，返回成功消息

# 7. 删除注册表键
curl -X DELETE "http://localhost:8081/api/v1/registry/key/SOFTWARE%5CTestKey"

**预期结果：**
- 注册表键删除成功，返回成功消息

# 8. 验证删除结果
curl "http://localhost:8081/api/v1/registry/key/SOFTWARE%5CTestKey"

**预期结果：**
- 返回键不存在或错误信息
```

##### Linux版本测试
```bash
# 注意：Linux没有注册表，这些API在Linux上可能返回错误或空结果
# 1. 创建注册表键（Linux上不可用）
curl -X POST http://localhost:8081/api/v1/registry/key \
  -H "Content-Type: application/json" \
  -d '{"path":"SOFTWARE\\TestKey"}'

**预期结果：**
- 返回错误信息或空结果

# 2. 设置注册表值（Linux上不可用）
curl -X PUT http://localhost:8081/api/v1/registry/value \
  -H "Content-Type: application/json" \
  -d '{"path":"SOFTWARE\\TestKey","name":"TestValue","type":"string","value":"test"}'

**预期结果：**
- 返回错误信息或空结果

# 3. 删除注册表值（Linux上不可用）
curl -X DELETE "http://localhost:8081/api/v1/registry/value/SOFTWARE%5CTestKey/TestValue"

**预期结果：**
- 返回错误信息或空结果

# 4. 删除注册表键（Linux上不可用）
curl -X DELETE "http://localhost:8081/api/v1/registry/key/SOFTWARE%5CTestKey"

**预期结果：**
- 返回错误信息或空结果
```

#### 3.3 注册表搜索功能测试

**测试目标**：验证注册表搜索功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 搜索包含特定字符串的注册表项
curl "http://localhost:8081/api/v1/registry/search?root_path=SOFTWARE&search_term=Windows"

**预期结果：**
- 返回包含"Windows"的注册表项列表

# 2. 搜索包含特定字符串的注册表值
curl "http://localhost:8081/api/v1/registry/search?root_path=SOFTWARE&search_term=Microsoft"

**预期结果：**
- 返回包含"Microsoft"的注册表值列表

# 3. 搜索特定路径下的注册表项
curl "http://localhost:8081/api/v1/registry/search?root_path=SOFTWARE\\Microsoft&search_term=Windows"

**预期结果：**
- 返回指定路径下包含"Windows"的注册表项

# 4. 搜索系统路径下的注册表项
curl "http://localhost:8081/api/v1/registry/search?root_path=SYSTEM&search_term=CurrentControlSet"

**预期结果：**
- 返回系统路径下包含"CurrentControlSet"的注册表项

# 5. 搜索用户路径下的注册表项
curl "http://localhost:8081/api/v1/registry/search?root_path=HKEY_CURRENT_USER&search_term=Software"

**预期结果：**
- 返回用户路径下包含"Software"的注册表项
```

##### Linux版本测试
```bash
# 注意：Linux没有注册表，这些API在Linux上可能返回错误或空结果
# 1. 搜索注册表项（Linux上不可用）
curl "http://localhost:8081/api/v1/registry/search?root_path=SOFTWARE&search_term=Windows"

**预期结果：**
- 返回错误信息或空结果

# 2. 搜索注册表值（Linux上不可用）
curl "http://localhost:8081/api/v1/registry/search?root_path=SOFTWARE&search_term=Microsoft"

**预期结果：**
- 返回错误信息或空结果

# 3. 搜索特定路径（Linux上不可用）
curl "http://localhost:8081/api/v1/registry/search?root_path=SOFTWARE\\Microsoft&search_term=Windows"

**预期结果：**
- 返回错误信息或空结果
```

---

### 4. 网络连接功能验证

#### 4.1 网络连接查看功能测试

**测试目标**：验证网络连接的查看和统计功能

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 获取所有网络连接
curl http://localhost:8081/api/v1/network/connections

**预期结果：**
- 返回系统网络连接列表，包含：
# - 本地地址和端口
# - 远程地址和端口
# - 连接状态、类
# - 进程ID、进程名

# 2. 获取TCP连接
curl http://localhost:8081/api/v1/network/connections/tcp

**预期结果：**
- 返回TCP连接列表，包含：
# - TCP连接状态（LISTEN、ESTABLISHED、TIME_WAIT等）
# - 本地和远程地址端口
# - 相关进程信息

# 3. 获取UDP连接
curl http://localhost:8081/api/v1/network/connections/udp

**预期结果：**
- 返回UDP连接列表，包含：
# - UDP连接状态
# - 本地和远程地址端口
# - 相关进程信息

# 4. 按进程ID获取连接
curl http://localhost:8081/api/v1/network/connections/pid/22792

**预期结果：**
- 返回指定进程的网络连接列表

# 5. 按端口获取连接
curl http://localhost:8081/api/v1/network/connections/port/8081

**预期结果：**
- 返回使用指定端口的连接列表

# 6. 按IP地址获取连接
curl http://localhost:8081/api/v1/network/connections/ip/127.0.0.1

**预期结果：**
- 返回与指定IP地址相关的连接列表

# 7. 获取网络统计信息
curl http://localhost:8081/api/v1/network/stats

**预期结果：**
- 返回网络统计信息，包含：
# - 总连接数
# - TCP连接数、UDP连接数
# - 监听端口
# - 已建立连接数

# 8. 获取监听端口列表
curl http://localhost:8081/api/v1/network/listening-ports

**预期结果：**
- 返回监听端口列表，包含：
# - 端口号、协议类型
# - 监听地址、状态
# - 相关进程信息

# 9. 获取已建立连接列表
curl http://localhost:8081/api/v1/network/established-connections

**预期结果：**
- 返回已建立的连接列表

# 10. 获取网络接口信息
curl http://localhost:8081/api/v1/network/interfaces

**预期结果：**
- 返回网络接口信息，包含：
# - 接口名称、MAC地址
# - IP地址、子网掩码
# - 接口状态、类
```

##### Linux版本测试
```bash
# 1. 获取所有网络连接
curl http://localhost:8081/api/v1/network/connections

**预期结果：**
- 返回系统网络连接列表，包含：
# - 本地地址和端口
# - 远程地址和端口
# - 连接状态、类
# - 进程ID、进程名

# 2. 获取TCP连接
curl http://localhost:8081/api/v1/network/connections/tcp

**预期结果：**
- 返回TCP连接列表，包含：
# - TCP连接状态（LISTEN、ESTABLISHED、TIME_WAIT等）
# - 本地和远程地址端口
# - 相关进程信息

# 3. 获取UDP连接
curl http://localhost:8081/api/v1/network/connections/udp

**预期结果：**
- 返回UDP连接列表，包含：
# - UDP连接状态
# - 本地和远程地址端口
# - 相关进程信息

# 4. 按进程ID获取连接
curl http://localhost:8081/api/v1/network/connections/pid/1234

**预期结果：**
- 返回指定进程的网络连接列表

# 5. 按端口获取连接
curl http://localhost:8081/api/v1/network/connections/port/22

**预期结果：**
- 返回使用指定端口的连接列表

# 6. 按IP地址获取连接
curl http://localhost:8081/api/v1/network/connections/ip/127.0.0.1

**预期结果：**
- 返回与指定IP地址相关的连接列表

# 7. 获取网络统计信息
curl http://localhost:8081/api/v1/network/stats

**预期结果：**
- 返回网络统计信息，包含：
# - 总连接数
# - TCP连接数、UDP连接数
# - 监听端口
# - 已建立连接数

# 8. 获取监听端口列表
curl http://localhost:8081/api/v1/network/listening-ports

**预期结果：**
- 返回监听端口列表，包含：
# - 端口号、协议类型
# - 监听地址、状态
# - 相关进程信息

# 9. 获取已建立连接列表
curl http://localhost:8081/api/v1/network/established-connections

**预期结果：**
- 返回已建立的连接列表

# 10. 获取网络接口信息
curl http://localhost:8081/api/v1/network/interfaces

**预期结果：**
- 返回网络接口信息，包含：
# - 接口名称、MAC地址
# - IP地址、子网掩码
# - 接口状态、类
```

#### 4.2 网络连接操作功能测试

**测试目标**：验证网络连接的关闭功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 检查端口是否被占用
curl http://localhost:8081/api/v1/network/port/8081/in-use

**预期结果：**
- 返回端口占用状态（true/false）

# 2. 关闭指定连接（谨慎操作）
curl -X DELETE http://localhost:8081/api/v1/network/connection/400-tcp-0

**预期结果：**
- 连接关闭成功，返回操作状态，注意：关闭连接可能影响相关进程

# 3. 关闭特定端口的连接（谨慎操作）
curl -X DELETE "http://localhost:8081/api/v1/network/connection/port/8081"

**预期结果：**
- 关闭指定端口的所有连接，返回操作状态

# 4. 关闭特定IP的连接（谨慎操作）
curl -X DELETE "http://localhost:8081/api/v1/network/connection/ip/192.168.1.100"

**预期结果：**
- 关闭指定IP的所有连接，返回操作状态
```

##### Linux版本测试
```bash
# 1. 检查端口是否被占用
curl http://localhost:8081/api/v1/network/port/22/in-use

**预期结果：**
- 返回端口占用状态（true/false）

# 2. 关闭指定连接（谨慎操作）
curl -X DELETE http://localhost:8081/api/v1/network/connection/1234-tcp-0

**预期结果：**
- 连接关闭成功，返回操作状态，注意：关闭连接可能影响相关进程

# 3. 关闭特定端口的连接（谨慎操作）
curl -X DELETE "http://localhost:8081/api/v1/network/connection/port/22"

**预期结果：**
- 关闭指定端口的所有连接，返回操作状态

# 4. 关闭特定IP的连接（谨慎操作）
curl -X DELETE "http://localhost:8081/api/v1/network/connection/ip/192.168.1.100"

**预期结果：**
- 关闭指定IP的所有连接，返回操作状态
```

#### 4.3 网络统计功能测试

**测试目标**：验证网络统计功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 获取网络统计信息
curl http://localhost:8081/api/v1/network/stats

**预期结果：**
- 返回网络统计信息，包含：
# - 总连接数
# - TCP连接数、UDP连接数
# - 监听端口
# - 已建立连接数
# - 连接状态分类

# 2. 获取网络接口统计
curl http://localhost:8081/api/v1/network/interfaces

**预期结果：**
- 返回网络接口统计信息，包含：
# - 接口名称、状态
# - 接收/发送字节数
# - 接收/发送包数
# - 错误包数

# 3. 获取连接历史记录
curl http://localhost:8081/api/v1/network/monitoring/history

**预期结果：**
- 返回连接历史记录，包含：
# - 连接建立时间
# - 连接持续时间
# - 数据传输
# - 连接状态变化
```

##### Linux版本测试
```bash
# 1. 获取网络统计信息
curl http://localhost:8081/api/v1/network/stats

**预期结果：**
- 返回网络统计信息，包含：
# - 总连接数
# - TCP连接数、UDP连接数
# - 监听端口
# - 已建立连接数
# - 连接状态分类

# 2. 获取网络接口统计
curl http://localhost:8081/api/v1/network/interfaces

**预期结果：**
- 返回网络接口统计信息，包含：
# - 接口名称、状态
# - 接收/发送字节数
# - 接收/发送包数
# - 错误包数

# 3. 获取连接历史记录
curl http://localhost:8081/api/v1/network/monitoring/history

**预期结果：**
- 返回连接历史记录，包含：
# - 连接建立时间
# - 连接持续时间
# - 数据传输
# - 连接状态变化
```

#### 4.4 网络监控功能测试

**测试目标**：验证网络监控功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 启用网络监控
curl -X POST http://localhost:8081/api/v1/network/monitoring/enable

**预期结果：**
- 网络监控启用成功，返回成功消息

# 2. 获取监控的连接列表
curl http://localhost:8081/api/v1/network/monitoring/connections

**预期结果：**
- 返回当前监控的连接列表，包括：
# - 连接ID、连接状态
# - 本地和远程地址端口
# - 监控开始时间
# - 数据传输统计

# 3. 获取连接历史记录
curl http://localhost:8081/api/v1/network/monitoring/history

**预期结果：**
- 返回连接历史记录，包含：
# - 连接建立和断开时间
# - 连接持续时间
# - 数据传输
# - 连接状态变化

# 4. 禁用网络监控
curl -X POST http://localhost:8081/api/v1/network/monitoring/disable

**预期结果：**
- 网络监控禁用成功，返回成功消息

# 5. 再次获取监控的连接列表
curl http://localhost:8081/api/v1/network/monitoring/connections

**预期结果：**
- 返回空列表或监控已停止的消?
```

##### Linux版本测试
```bash
# 1. 启用网络监控
curl -X POST http://localhost:8081/api/v1/network/monitoring/enable

**预期结果：**
- 网络监控启用成功，返回成功消息

# 2. 获取监控的连接列表
curl http://localhost:8081/api/v1/network/monitoring/connections

**预期结果：**
- 返回当前监控的连接列表，包括：
# - 连接ID、连接状态
# - 本地和远程地址端口
# - 监控开始时间
# - 数据传输统计

# 3. 获取连接历史记录
curl http://localhost:8081/api/v1/network/monitoring/history

**预期结果：**
- 返回连接历史记录，包含：
# - 连接建立和断开时间
# - 连接持续时间
# - 数据传输
# - 连接状态变化

# 4. 禁用网络监控
curl -X POST http://localhost:8081/api/v1/network/monitoring/disable

**预期结果：**
- 网络监控禁用成功，返回成功消息

# 5. 再次获取监控的连接列表
curl http://localhost:8081/api/v1/network/monitoring/connections

**预期结果：**
- 返回空列表或监控已停止的消?
```

---

### 5. 用户权限功能验证

#### 5.1 用户信息获取功能测试

**测试目标**：验证用户信息的获取功能

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 获取当前用户信息
curl http://localhost:8081/api/v1/user/current

**预期结果：**
- 返回当前用户信息，包含：
# - 用户ID、用户名
# - 用户组、权限级别
# - 账户状态、登录时间
# - 用户描述、主目录

# 2. 获取所有用户列表
curl http://localhost:8081/api/v1/user/all

**预期结果：**
- 返回系统用户列表，包含：
# - 用户ID、用户名
# - 用户组、账户状态
# - 最后登录时间、账户类型

# 3. 根据用户ID获取用户信息
curl http://localhost:8081/api/v1/user/id/S-1-5-21-2079322416-3376148235-2999533130-1003

**预期结果：**
- 返回指定用户ID的详细信息

# 4. 根据用户名获取用户信息
curl http://localhost:8081/api/v1/user/name/Administrator

**预期结果：**
- 返回指定用户名的详细信息

# 5. 检查账户状态
curl http://localhost:8081/api/v1/user/status/Administrator

**预期结果：**
- 返回账户状态信息，包含：
# - 账户是否启用
# - 账户是否锁定
# - 密码是否过期
# - 账户过期时间

# 6. 获取用户会话信息
curl http://localhost:8081/api/v1/user/sessions/Administrator

**预期结果：**
- 返回用户会话信息，包含：
# - 会话ID、登录时间
# - 登录地址、会话状态
# - 会话持续时间
```

##### Linux版本测试
```bash
# 1. 获取当前用户信息
curl http://localhost:8081/api/v1/user/current

**预期结果：**
- 返回当前用户信息，包含：
# - 用户ID、用户名
# - 用户组、权限级别
# - 账户状态、登录时间
# - 用户描述、主目录

# 2. 获取所有用户列表
curl http://localhost:8081/api/v1/user/all

**预期结果：**
- 返回系统用户列表，包含：
# - 用户ID、用户名
# - 用户组、账户状态
# - 最后登录时间、账户类型

# 3. 根据用户ID获取用户信息
curl http://localhost:8081/api/v1/user/id/1000

**预期结果：**
- 返回指定用户ID的详细信息

# 4. 根据用户名获取用户信息
curl http://localhost:8081/api/v1/user/name/root

**预期结果：**
- 返回指定用户名的详细信息

# 5. 检查账户状态
curl http://localhost:8081/api/v1/user/status/root

**预期结果：**
- 返回账户状态信息，包含：
# - 账户是否启用
# - 账户是否锁定
# - 密码是否过期
# - 账户过期时间

# 6. 获取用户会话信息
curl http://localhost:8081/api/v1/user/sessions/root

**预期结果：**
- 返回用户会话信息，包含：
# - 会话ID、登录时间
# - 登录地址、会话状态
# - 会话持续时间
```

#### 5.2 用户权限检查功能测试

**测试目标**：验证用户权限的检查功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 检查用户权限
curl -X POST http://localhost:8081/api/v1/user/permissions/check \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"Administrator\",\"permissions\":[\"admin\"]}"

**预期结果：**
- 返回权限检查结果，包含：
# - 权限检查结果（true/false）
# - 用户拥有的权限列表
# - 缺失的权限列表
# - 权限级别

# 2. 验证密码
curl -X POST http://localhost:8081/api/v1/user/password/validate \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"Administrator\",\"password\":\"test\"}"

**预期结果：**
- 返回密码验证结果，包含：
# - 密码是否有效
# - 密码复杂度评估
# - 密码策略符合性
# - 改进建议

# 3. 获取密码策略
curl http://localhost:8081/api/v1/user/password/policy

**预期结果：**
- 返回密码策略信息，包含：
# - 最小密码长度
# - 密码复杂度要求
# - 密码历史要求
# - 密码过期策略

# 4. 检查特定权限
curl -X POST http://localhost:8081/api/v1/user/permissions/check \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"Administrator\",\"permissions\":[\"file_scan\",\"process_control\",\"registry_access\"]}"

**预期结果：**
- 返回特定权限检查结果

# 5. 检查管理员权限
curl -X POST http://localhost:8081/api/v1/user/permissions/check \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"Administrator\",\"permissions\":[\"admin\",\"system_control\"]}"

**预期结果：**
- 返回管理员权限检查结果
```

##### Linux版本测试
```bash
# 1. 检查用户权限
curl -X POST http://localhost:8081/api/v1/user/permissions/check \
  -H "Content-Type: application/json" \
  -d '{"username":"root","permissions":["admin"]}'

**预期结果：**
- 返回权限检查结果，包含：
# - 权限检查结果（true/false）
# - 用户拥有的权限列表
# - 缺失的权限列表
# - 权限级别

# 2. 验证密码
curl -X POST http://localhost:8081/api/v1/user/password/validate \
  -H "Content-Type: application/json" \
  -d '{"username":"root","password":"test"}'

**预期结果：**
- 返回密码验证结果，包含：
# - 密码是否有效
# - 密码复杂度评估
# - 密码策略符合性
# - 改进建议

# 3. 获取密码策略
curl http://localhost:8081/api/v1/user/password/policy

**预期结果：**
- 返回密码策略信息，包含：
# - 最小密码长度
# - 密码复杂度要求
# - 密码历史要求
# - 密码过期策略

# 4. 检查特定权限
curl -X POST http://localhost:8081/api/v1/user/permissions/check \
  -H "Content-Type: application/json" \
  -d '{"username":"root","permissions":["file_scan","process_control","network_access"]}'

**预期结果：**
- 返回特定权限检查结果

# 5. 检查超级用户权限
curl -X POST http://localhost:8081/api/v1/user/permissions/check \
  -H "Content-Type: application/json" \
  -d '{"username":"root","permissions":["admin","system_control"]}'

**预期结果：**
- 返回超级用户权限检查结果
```

#### 5.3 用户管理功能测试

**测试目标**：验证用户账户管理功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 锁定用户账户
curl -X POST http://localhost:8081/api/v1/user/lock/Administrator

**预期结果：**
- 用户账户锁定成功，返回操作状态，注意：锁定后用户无法登录

# 2. 解锁用户账户
curl -X POST http://localhost:8081/api/v1/user/unlock/Administrator

**预期结果：**
- 用户账户解锁成功，返回操作状态，注意：解锁后用户可以正常登录

# 3. 更改用户密码
curl -X POST http://localhost:8081/api/v1/user/password/change \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"Administrator\",\"old_password\":\"old\",\"new_password\":\"new\"}"

**预期结果：**
- 密码更改成功，返回操作状态，注意：新密码需要符合密码策略

# 4. 获取用户组信息
curl http://localhost:8081/api/v1/user/groups/Administrator

**预期结果：**
- 返回用户所属的组列表，包含：
# - 组ID、组名称
# - 组描述、组类型
# - 组成员数

# 5. 获取用户登录历史
curl http://localhost:8081/api/v1/user/history/Administrator

**预期结果：**
- 返回用户登录历史，包含：
# - 登录时间、登出时间
# - 登录地址、登录方式
# - 会话持续时间
# - 登录状态

# 6. 终止用户会话
curl -X DELETE http://localhost:8081/api/v1/user/session/session_id

**预期结果：**
- 用户会话终止成功，返回操作状态，注意：session_id需要替换为实际的会话ID

# 7. 检查账户状态变化
curl http://localhost:8081/api/v1/user/status/Administrator

**预期结果：**
- 返回最新的账户状态信息
```

##### Linux版本测试
```bash
# 1. 锁定用户账户
curl -X POST http://localhost:8081/api/v1/user/lock/root

**预期结果：**
- 用户账户锁定成功，返回操作状态，注意：锁定后用户无法登录

# 2. 解锁用户账户
curl -X POST http://localhost:8081/api/v1/user/unlock/root

**预期结果：**
- 用户账户解锁成功，返回操作状态，注意：解锁后用户可以正常登录

# 3. 更改用户密码
curl -X POST http://localhost:8081/api/v1/user/password/change \
  -H "Content-Type: application/json" \
  -d '{"username":"root","old_password":"old","new_password":"new"}'

**预期结果：**
- 密码更改成功，返回操作状态，注意：新密码需要符合密码策略

# 4. 获取用户组信息
curl http://localhost:8081/api/v1/user/groups/root

**预期结果：**
- 返回用户所属的组列表，包含：
# - 组ID、组名称
# - 组描述、组类型
# - 组成员数

# 5. 获取用户登录历史
curl http://localhost:8081/api/v1/user/history/root

**预期结果：**
- 返回用户登录历史，包含：
# - 登录时间、登出时间
# - 登录地址、登录方式
# - 会话持续时间
# - 登录状态

# 6. 终止用户会话
curl -X DELETE http://localhost:8081/api/v1/user/session/session_id

**预期结果：**
- 用户会话终止成功，返回操作状态，注意：session_id需要替换为实际的会话ID

# 7. 检查账户状态变化
curl http://localhost:8081/api/v1/user/status/root

**预期结果：**
- 返回最新的账户状态信息
```

---

### 6. 安全功能验证

#### 6.1 安全状态查询测试

**测试目标**：验证安全状态查询功能

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 获取安全状态
curl http://localhost:8081/api/v1/security/status

**预期结果：**
- 返回安全状态信息，包含：
# - 扫描引擎状态
# - 规则加载状态
# - 缓存状态
# - 威胁检测系统
# - 系统安全级别

# 2. 获取规则信息
curl http://localhost:8081/api/v1/security/rules

**预期结果：**
- 返回规则信息，包含：
# - 已加载规则数量
# - 规则文件路径
# - 规则更新时间
# - 规则统计信息

# 3. 获取缓存统计
curl http://localhost:8081/api/v1/security/cache/stats

**预期结果：**
- 返回缓存统计信息，包含：
# - 缓存命中次数
# - 缓存大小
# - 缓存条目数
# - 缓存清理时间
```

##### Linux版本测试
```bash
# 1. 获取安全状态
curl http://localhost:8081/api/v1/security/status

**预期结果：**
- 返回安全状态信息，包含：
# - 扫描引擎状态
# - 规则加载状态
# - 缓存状态
# - 威胁检测系统
# - 系统安全级别

# 2. 获取规则信息
curl http://localhost:8081/api/v1/security/rules

**预期结果：**
- 返回规则信息，包含：
# - 已加载规则数量
# - 规则文件路径
# - 规则更新时间
# - 规则统计信息

# 3. 获取缓存统计
curl http://localhost:8081/api/v1/security/cache/stats

**预期结果：**
- 返回缓存统计信息，包含：
# - 缓存命中次数
# - 缓存大小
# - 缓存条目数
# - 缓存清理时间
```

#### 6.2 安全规则管理测试

**测试目标**：验证安全规则管理功能

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 重新加载规则
curl -X POST http://localhost:8081/api/v1/security/reload-rules

**预期结果：**
- 规则重新加载成功，返回操作状态，注意：重新加载后规则会更改

# 2. 清空缓存
curl -X POST http://localhost:8081/api/v1/security/cache/clear

**预期结果：**
- 缓存清空成功，返回操作状态，注意：清空后下次扫描会重新计算

# 3. 验证规则重新加载
curl http://localhost:8081/api/v1/security/rules

**预期结果：**
- 返回更新后的规则信息

# 4. 验证缓存清空
curl http://localhost:8081/api/v1/security/cache/stats

**预期结果：**
- 返回清空后的缓存统计信息
```

##### Linux版本测试
```bash
# 1. 重新加载规则
curl -X POST http://localhost:8081/api/v1/security/reload-rules

**预期结果：**
- 规则重新加载成功，返回操作状态，注意：重新加载后规则会更改

# 2. 清空缓存
curl -X POST http://localhost:8081/api/v1/security/cache/clear

**预期结果：**
- 缓存清空成功，返回操作状态，注意：清空后下次扫描会重新计算

# 3. 验证规则重新加载
curl http://localhost:8081/api/v1/security/rules

**预期结果：**
- 返回更新后的规则信息

# 4. 验证缓存清空
curl http://localhost:8081/api/v1/security/cache/stats

**预期结果：**
- 返回清空后的缓存统计信息
```

#### 6.3 文件隔离功能测试

**测试目标**：验证文件隔离功能

**测试步骤**：

##### Windows版本测试
```bash
# 1. 隔离可疑文件
curl -X POST http://localhost:8081/api/v1/security/quarantine \
  -H "Content-Type: application/json" \
  -d "{\"file_path\":\"C:\\\\temp\\\\suspicious.exe\"}"

**预期结果：**
- 文件隔离成功，返回操作状态，注意：隔离后文件会被移动到隔离区

# 2. 获取隔离文件列表
curl http://localhost:8081/api/v1/security/quarantine/list

**预期结果：**
- 返回隔离文件列表，包含：
# - 文件路径、隔离时间
# - 威胁级别、威胁类型
# - 文件大小、哈希值

# 3. 恢复隔离文件
curl -X POST http://localhost:8081/api/v1/security/restore \
  -H "Content-Type: application/json" \
  -d "{\"file_path\":\"C:\\\\temp\\\\suspicious.exe\"}"

**预期结果：**
- 文件恢复成功，返回操作状态，注意：恢复后文件会从隔离区移回原位置

# 4. 获取扫描历史
curl http://localhost:8081/api/v1/security/scan-history

**预期结果：**
- 返回扫描历史记录，包含：
# - 扫描时间、文件路径
# - 扫描结果、威胁信息
# - 扫描耗时、扫描状态
```

##### Linux版本测试
```bash
# 1. 隔离可疑文件
curl -X POST http://localhost:8081/api/v1/security/quarantine \
  -H "Content-Type: application/json" \
  -d '{"file_path":"/tmp/suspicious_file"}'

**预期结果：**
- 文件隔离成功，返回操作状态，注意：隔离后文件会被移动到隔离区

# 2. 获取隔离文件列表
curl http://localhost:8081/api/v1/security/quarantine/list

**预期结果：**
- 返回隔离文件列表，包含：
# - 文件路径、隔离时间
# - 威胁级别、威胁类型
# - 文件大小、哈希值

# 3. 恢复隔离文件
curl -X POST http://localhost:8081/api/v1/security/restore \
  -H "Content-Type: application/json" \
  -d '{"file_path":"/tmp/suspicious_file"}'

**预期结果：**
- 文件恢复成功，返回操作状态，注意：恢复后文件会从隔离区移回原位置

# 4. 获取扫描历史
curl http://localhost:8081/api/v1/security/scan-history

**预期结果：**
- 返回扫描历史记录，包含：
# - 扫描时间、文件路径
# - 扫描结果、威胁信息
# - 扫描耗时、扫描状态
```

---

### 7. 系统监控功能验证

#### 7.1 健康检查测试

**测试目标**：验证系统健康检查功能

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 健康检查
curl http://localhost:8081/api/health

**预期结果：**
- 返回健康状态信息，包含：
# - 服务状态（healthy）
# - 运行时间
# - 版本信息

# 2. 检查服务状态
curl http://localhost:8081/api/health

**预期结果：**
- 返回服务正常运行状态

# 3. 检查API响应时间
time curl http://localhost:8081/api/health

**预期结果：**
- API响应时间应该在合理范围内（1秒）
```

##### Linux版本测试
```bash
# 1. 健康检查
curl http://localhost:8081/api/health

**预期结果：**
- 返回健康状态信息，包含：
# - 服务状态（healthy）
# - 运行时间
# - 版本信息

# 2. 检查服务状态
curl http://localhost:8081/api/health

**预期结果：**
- 返回服务正常运行状态

# 3. 检查API响应时间
time curl http://localhost:8081/api/health

**预期结果：**
- API响应时间应该在合理范围内（1秒）
```

#### 7.2 性能指标测试

**测试目标**：验证性能指标收集功能

**测试步骤**：

##### Windows版本测试
```cmd
# 1. 获取性能指标
curl http://localhost:8081/api/metrics

**预期结果：**
- 返回性能指标信息，包含：
# - API调用次数
# - 平均响应时间
# - 错误次数
# - 并发连接数
# - 内存使用情况
# - CPU使用情况

# 2. 重置性能指标
curl -X POST http://localhost:8081/api/metrics/reset

**预期结果：**
- 性能指标重置成功，返回操作状态

# 3. 验证指标重置
curl http://localhost:8081/api/metrics

**预期结果：**
- 返回重置后的性能指标信息

# 4. 检查指标收集
# 执行一些API调用后再次检查
curl http://localhost:8081/api/v1/process/list
curl http://localhost:8081/api/v1/network/connections
curl http://localhost:8081/api/metrics

**预期结果：**
- 返回更新后的性能指标信息
```

##### Linux版本测试
```bash
# 1. 获取性能指标
curl http://localhost:8081/api/metrics

**预期结果：**
- 返回性能指标信息，包含：
# - API调用次数
# - 平均响应时间
# - 错误次数
# - 并发连接数
# - 内存使用情况
# - CPU使用情况

# 2. 重置性能指标
curl -X POST http://localhost:8081/api/metrics/reset

**预期结果：**
- 性能指标重置成功，返回操作状态

# 3. 验证指标重置
curl http://localhost:8081/api/metrics

**预期结果：**
- 返回重置后的性能指标信息

# 4. 检查指标收集
# 执行一些API调用后再次检查
curl http://localhost:8081/api/v1/process/list
curl http://localhost:8081/api/v1/network/connections
curl http://localhost:8081/api/metrics

**预期结果：**
- 返回更新后的性能指标信息
```

---

## 📋 新增API接口测试用例

### 🔍 文件检测模块新增测试

#### 1.15 目录扫描测试
**测试目标**：验证系统能够扫描整个目录并检测威胁

**测试步骤**：
```cmd
# 1. 基本目录扫描
curl -X POST http://localhost:8081/api/v1/file/scan-directory -H "Content-Type: application/json" -d "{\"path\":\"C:\\\\temp\"}"

**预期结果：**
- 返回目录扫描结果，包含扫描文件数量和威胁统计

# 2. 带过滤条件的目录扫描
curl -X POST http://localhost:8081/api/v1/file/scan-directory -H "Content-Type: application/json" -d "{\"path\":\"C:\\\\temp\",\"recursive\":true,\"max_depth\":2,\"include_patterns\":[\"*.exe\",\"*.dll\"],\"exclude_patterns\":[\"*.tmp\"]}"

**预期结果：**
- 只扫描指定类型的文件，排除临时文件

# 3. 获取文件列表
curl "http://localhost:8081/api/v1/file/list?path=C:\temp&recursive=true&max_depth=2"

**预期结果：**
- 返回指定目录的文件列表
```

#### 1.16 文件操作测试
**测试目标**：验证文件复制、移动、删除功能

**测试步骤**：
```cmd
# 1. 复制文件
curl -X POST http://localhost:8081/api/v1/file/copy -H "Content-Type: application/json" -d "{\"source\":\"C:\\\\temp\\\\source.txt\",\"dest\":\"C:\\\\temp\\\\backup\\\\source.txt\"}"

**预期结果：**
- 文件复制成功

# 2. 移动文件
curl -X POST http://localhost:8081/api/v1/file/move -H "Content-Type: application/json" -d "{\"source\":\"C:\\\\temp\\\\old_name.txt\",\"dest\":\"C:\\\\temp\\\\new_name.txt\"}"

**预期结果：**
- 文件移动成功

# 3. 删除文件
curl -X DELETE "http://localhost:8081/api/v1/file/test.txt"

**预期结果：**
- 文件删除成功
```

#### 1.17 内存缓冲区扫描测试
**测试目标**：验证系统能够扫描内存中的缓冲区数据

**测试步骤**：
```cmd
# 1. 扫描网络载荷
curl -X POST http://localhost:8081/api/v1/file/scan-buffer -H "Content-Type: application/json" -d "{\"identifier\":\"network_payload\",\"data\":\"http://malicious-site.com/download.exe\"}"

**预期结果：**
- 检测到恶意URL或载荷

# 2. 扫描文件内容
curl -X POST http://localhost:8081/api/v1/file/scan-buffer -H "Content-Type: application/json" -d "{\"identifier\":\"file_content\",\"data\":\"MZ...PE文件头内容...\"}"

**预期结果：**
- 识别PE文件头结构

# 3. 扫描URL
curl -X POST http://localhost:8081/api/v1/file/scan-buffer -H "Content-Type: application/json" -d "{\"identifier\":\"url_check\",\"data\":\"https://example.com/suspicious.php\"}"

**预期结果：**
- 分析URL安全性
```

### 👥 进程管理模块新增测试

#### 2.6 进程监控测试
**测试目标**：验证进程监控功能

**测试步骤**：
```cmd
# 1. 启用进程监控
curl -X POST http://localhost:8081/api/v1/process/monitoring/enable -H "Content-Type: application/json" -d "{\"enabled\":true}"

**预期结果：**
- 进程监控启用成功

# 2. 获取监控进程列表
curl "http://localhost:8081/api/v1/process/monitoring/list"

**预期结果：**
- 返回正在监控的进程列表

# 3. 获取进程统计信息
curl "http://localhost:8081/api/v1/process/statistics"

**预期结果：**
- 返回进程统计信息
```

### 🔧 注册表操作模块新增测试

#### 3.4 注册表搜索测试
**测试目标**：验证注册表搜索功能

**测试步骤**：
```cmd
# 1. 搜索注册表键
curl -X POST http://localhost:8081/api/v1/registry/search -H "Content-Type: application/json" -d "{\"pattern\":\"malware\",\"search_type\":\"key\"}"

**预期结果：**
- 返回匹配的注册表键

# 2. 搜索注册表值
curl -X POST http://localhost:8081/api/v1/registry/search -H "Content-Type: application/json" -d "{\"pattern\":\"suspicious\",\"search_type\":\"value\"}"

**预期结果：**
- 返回匹配的注册表值

# 3. 列出注册表键
curl "http://localhost:8081/api/v1/registry/keys/HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion"

**预期结果：**
- 返回指定路径下的注册表键列表

# 4. 列出注册表值
curl "http://localhost:8081/api/v1/registry/values/HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion"

**预期结果：**
- 返回指定路径下的注册表值列表
```

### 🌐 网络连接模块新增测试

#### 4.5 网络监控测试
**测试目标**：验证网络监控功能

**测试步骤**：
```cmd
# 1. 启用网络监控
curl -X POST http://localhost:8081/api/v1/network/monitoring/enable -H "Content-Type: application/json" -d "{\"enabled\":true}"

**预期结果：**
- 网络监控启用成功

# 2. 获取监控连接列表
curl "http://localhost:8081/api/v1/network/monitoring/list"

**预期结果：**
- 返回正在监控的网络连接

# 3. 获取连接历史记录
curl "http://localhost:8081/api/v1/network/monitoring/history?limit=50"

**预期结果：**
- 返回网络连接历史记录

# 4. 禁用网络监控
curl -X POST http://localhost:8081/api/v1/network/monitoring/disable

**预期结果：**
- 网络监控禁用成功
```

### 👤 用户权限模块新增测试

#### 5.4 用户会话管理测试
**测试目标**：验证用户会话管理功能

**测试步骤**：
```cmd
# 1. 获取用户会话
curl "http://localhost:8081/api/v1/user/sessions"

**预期结果：**
- 返回当前用户的所有会话

# 2. 结束用户会话
curl -X DELETE "http://localhost:8081/api/v1/user/session/1234567890abcdef1234567890abcdef1234567890"

**预期结果：**
- 指定会话结束成功

# 3. 锁定用户账户
curl -X PUT http://localhost:8081/api/v1/user/lock/Administrator -H "Content-Type: application/json" -d "{\"reason\":\"Suspicious activity detected\",\"duration\":30}"

**预期结果：**
- 用户账户锁定成功

# 4. 解锁用户账户
curl -X PUT "http://localhost:8081/api/v1/user/unlock/Administrator"

**预期结果：**
- 用户账户解锁成功

# 5. 修改用户密码
curl -X PUT http://localhost:8081/api/v1/user/password/Administrator -H "Content-Type: application/json" -d "{\"new_password\":\"NewPassword123!\",\"old_password\":\"OldPassword123!\"}"

**预期结果：**
- 用户密码修改成功
```

### 🛡️ 安全功能模块新增测试

#### 6.4 文件隔离功能测试
**测试目标**：验证文件隔离和恢复功能

**测试步骤**：
```cmd
# 1. 隔离可疑文件
curl -X POST http://localhost:8081/api/v1/security/quarantine -H "Content-Type: application/json" -d "{\"file_path\":\"C:\\\\temp\\\\suspicious.exe\",\"reason\":\"Malware detected\",\"threat_level\":\"high\"}"

**预期结果：**
- 文件隔离成功

# 2. 获取隔离列表
curl "http://localhost:8081/api/v1/security/quarantine/list"

**预期结果：**
- 返回隔离文件列表

# 3. 恢复隔离文件
curl -X POST http://localhost:8081/api/v1/security/restore -H "Content-Type: application/json" -d "{\"file_path\":\"C:\\\\temp\\\\suspicious.exe\",\"quarantine_id\":\"quarantine-12345\"}"

**预期结果：**
- 文件恢复成功

# 4. 获取扫描历史
curl "http://localhost:8081/api/v1/security/scan-history?limit=20"

**预期结果：**
- 返回扫描历史记录
```

#### 6.5 安全规则管理测试
**测试目标**：验证安全规则管理功能

**测试步骤**：
```cmd
# 1. 获取规则信息
curl "http://localhost:8081/api/v1/security/rules"

**预期结果：**
- 返回当前加载的YARA规则信息

# 2. 重新加载规则
curl -X POST "http://localhost:8081/api/v1/security/reload-rules"

**预期结果：**
- 规则重新加载成功

# 3. 清空缓存
curl -X POST "http://localhost:8081/api/v1/security/clear-cache"

**预期结果：**
- 缓存清空成功

# 4. 获取缓存统计
curl "http://localhost:8081/api/v1/security/cache/stats"

**预期结果：**
- 返回缓存统计信息
```

---

## 验证完成标准

### 功能完整性验证
- [x] 病毒木马文件检测功能完整
- [x] 文件系统功能完整
- [x] 进程管理功能完整
- [x] 注册表操作功能完整
- [x] 网络连接功能完整
- [x] 系统模块操作功能完整
- [x] 用户权限操作功能完整
- [x] 安全功能完整
- [x] 系统监控功能完整

### 性能指标验证
- [x] API响应时间符合要求（10秒）
- [x] 并发处理能力符合要求
- [x] 内存使用稳定
- [x] CPU使用率合理

### 稳定性验证
- [x] 长时间运行稳定
- [x] 错误处理完善
- [x] 日志记录完整
- [x] 异常恢复能力

### 跨平台兼容性验证
- [x] Windows系统功能正常
- [x] Linux系统功能正常
- [x] API接口统一
- [x] 数据格式一致

---

## 测试注意事项

### 安全测试注意事项
1. **谨慎操作**：某些操作（如进程结束、连接关闭、模块挂起）可能影响系统稳定
2. **权限要求**：部分操作需要管理员权限
3. **数据备份**：重要数据操作前请备份
4. **测试环境**：建议在测试环境中进行验证

### 性能测试注意事项
1. **负载测试**：避免在高负载环境下进行测试
2. **资源监控**：测试时监控系统资源使用情况
3. **超时设置**：某些操作可能需要较长时间，请耐心等待
4. **并发限制**：避免同时执行大量并发操作

### 兼容性测试注意事项
1. **系统差异**：Windows和Linux系统功能实现可能有差异
2. **路径格式**：注意不同系统的路径格式差异
3. **权限模型**：不同系统的权限模型不同
4. **API响应**：某些API在不同系统上可能返回不同结果

---

## 📊 接口路径和参数更新总结

### ✅ 已完成的更新

#### 🔍 文件检测模块接口更新
1. **文件扫描接口**
   - `POST /api/v1/file/scan` - 扫描单个文件
   - `POST /api/v1/file/scan-directory` - 扫描目录
   - `POST /api/v1/file/scan-buffer` - 扫描内存缓冲区

2. **文件信息接口**
   - `GET /api/v1/file/info/:path` - 获取文件信息（支持路径参数和fullPath查询参数）
   - `GET /api/v1/file/hash/:path` - 获取文件哈希（支持算法参数）
   - `GET /api/v1/file/hashes/:path` - 获取文件所有哈希值
   - `POST /api/v1/file/verify-hash` - 验证文件哈希

3. **文件操作接口**
   - `POST /api/v1/file/copy` - 复制文件
   - `POST /api/v1/file/move` - 移动文件
   - `DELETE /api/v1/file/:path` - 删除文件
   - `GET /api/v1/file/list` - 获取文件列表

#### 👥 进程管理模块接口更新
1. **进程信息接口**
   - `GET /api/v1/process/list` - 获取进程列表
   - `GET /api/v1/process/:pid` - 获取进程信息
   - `GET /api/v1/process/:pid/modules` - 获取进程模块
   - `GET /api/v1/process/:pid/connections` - 获取进程连接
   - `GET /api/v1/process/:pid/memory` - 获取进程内存信息
   - `GET /api/v1/process/:pid/running` - 检查进程运行状态
   - `GET /api/v1/process/:pid/children` - 获取进程子进程

2. **进程操作接口**
   - `POST /api/v1/process/start` - 启动进程
   - `DELETE /api/v1/process/:pid` - 结束进程
   - `PUT /api/v1/process/:pid/suspend` - 挂起进程
   - `PUT /api/v1/process/:pid/resume` - 恢复进程

3. **系统模块接口**
   - `GET /api/v1/process/modules` - 获取系统模块列表
   - `GET /api/v1/process/module/:module` - 获取模块信息
   - `PUT /api/v1/process/module/:module/suspend` - 挂起模块
   - `PUT /api/v1/process/module/:module/resume` - 恢复模块
   - `DELETE /api/v1/process/module/:module` - 结束模块

4. **进程监控接口**
   - `POST /api/v1/process/monitoring/enable` - 启用进程监控
   - `POST /api/v1/process/monitoring/disable` - 禁用进程监控
   - `GET /api/v1/process/monitoring/list` - 获取监控进程列表
   - `GET /api/v1/process/statistics` - 获取进程统计信息

#### 🔧 注册表操作模块接口更新
1. **注册表查询接口**
   - `GET /api/v1/registry/key/:path` - 获取注册表键
   - `GET /api/v1/registry/keys/:path` - 列出注册表键
   - `GET /api/v1/registry/values/:path` - 列出注册表值
   - `GET /api/v1/registry/search` - 搜索注册表

2. **注册表值接口**
   - `GET /api/v1/registry/value/:path` - 获取注册表值（查询参数方式）
   - `GET /api/v1/registry/value/:path/:name` - 获取注册表值（路径参数方式）
   - `PUT /api/v1/registry/value` - 设置注册表值
   - `DELETE /api/v1/registry/value/:path` - 删除注册表值（查询参数方式）
   - `DELETE /api/v1/registry/value/:path/:name` - 删除注册表值（路径参数方式）

3. **注册表操作接口**
   - `POST /api/v1/registry/key` - 创建注册表键
   - `DELETE /api/v1/registry/key/:path` - 删除注册表键

#### 🌐 网络连接模块接口更新
1. **网络连接接口**
   - `GET /api/v1/network/connections` - 获取所有网络连接
   - `GET /api/v1/network/connections/tcp` - 获取TCP连接
   - `GET /api/v1/network/connections/udp` - 获取UDP连接
   - `GET /api/v1/network/connections/pid/:pid` - 按进程ID获取连接
   - `GET /api/v1/network/connections/port/:port` - 按端口获取连接
   - `GET /api/v1/network/connections/ip/:ip` - 按IP地址获取连接

2. **网络操作接口**
   - `DELETE /api/v1/network/connection/:id` - 关闭网络连接
   - `GET /api/v1/network/port/:port/in-use` - 检查端口占用
   - `GET /api/v1/network/listening-ports` - 获取监听端口
   - `GET /api/v1/network/established-connections` - 获取已建立连接

3. **网络信息接口**
   - `GET /api/v1/network/interfaces` - 获取网络接口
   - `GET /api/v1/network/stats` - 获取网络统计信息

4. **网络监控接口**
   - `POST /api/v1/network/monitoring/enable` - 启用网络监控
   - `POST /api/v1/network/monitoring/disable` - 禁用网络监控
   - `GET /api/v1/network/monitoring/connections` - 获取监控连接列表
   - `GET /api/v1/network/monitoring/history` - 获取连接历史记录

#### 👤 用户权限模块接口更新
1. **用户信息接口**
   - `GET /api/v1/user/current` - 获取当前用户信息
   - `GET /api/v1/user/all` - 获取所有用户列表
   - `GET /api/v1/user/id/:uid` - 根据用户ID获取用户信息
   - `GET /api/v1/user/name/:username` - 根据用户名获取用户信息
   - `GET /api/v1/user/status/:username` - 检查账户状态
   - `GET /api/v1/user/sessions/:username` - 获取用户会话
   - `GET /api/v1/user/groups/:username` - 获取用户组
   - `GET /api/v1/user/history/:username` - 获取用户登录历史

2. **用户权限接口**
   - `POST /api/v1/user/permissions/check` - 检查用户权限
   - `POST /api/v1/user/password/validate` - 验证用户密码
   - `GET /api/v1/user/password/policy` - 获取密码策略

3. **用户管理接口**
   - `POST /api/v1/user/lock/:username` - 锁定用户账户
   - `POST /api/v1/user/unlock/:username` - 解锁用户账户
   - `POST /api/v1/user/password/change` - 修改用户密码
   - `DELETE /api/v1/user/session/:sessionId` - 结束用户会话

#### 🛡️ 安全功能模块接口更新
1. **安全状态接口**
   - `GET /api/v1/security/status` - 获取安全状态
   - `GET /api/v1/security/rules` - 获取规则信息
   - `GET /api/v1/security/cache/stats` - 获取缓存统计

2. **安全操作接口**
   - `POST /api/v1/security/reload-rules` - 重新加载规则
   - `POST /api/v1/security/cache/clear` - 清空缓存
   - `POST /api/v1/security/quarantine` - 隔离文件
   - `POST /api/v1/security/restore` - 恢复文件

3. **安全信息接口**
   - `GET /api/v1/security/quarantine/list` - 获取隔离列表
   - `GET /api/v1/security/scan-history` - 获取扫描历史

#### 📊 系统监控模块接口更新
1. **系统状态接口**
   - `GET /api/health` - 健康检查
   - `GET /api/metrics` - 获取性能指标
   - `POST /api/metrics/reset` - 重置性能指标

### 🔄 接口路径参数说明

#### 路径参数格式
- 所有路径参数都使用 `:param` 格式，而不是 `/*param` 格式
- 支持 `fullPath` 查询参数来传递完整路径
- 路径参数支持URL编码和未编码两种格式

#### 查询参数说明
- `fullPath`: 完整路径（当路径参数不包含斜杠时使用）
- `algorithm`: 哈希算法（md5, sha1, sha256等）
- `recursive`: 是否递归（true/false）
- `max_depth`: 最大深度（整数）
- `limit`: 限制数量（整数）

#### 请求体参数说明
- 所有POST请求都使用JSON格式
- 必需参数使用 `binding:"required"` 标签
- 可选参数使用指针类型或默认值

### 📝 响应格式统一

所有API响应都遵循统一的格式：
```json
{
  "code": 200,
  "message": "操作成功",
  "time": "2024-01-15T10:30:00Z",
  "data": {}
}
```

### 🎯 完成状态

✅ **所有接口路径和参数信息已更新完成！**

- ✅ 文件检测模块：11个接口
- ✅ 进程管理模块：20个接口
- ✅ 注册表操作模块：9个接口
- ✅ 网络连接模块：16个接口
- ✅ 用户权限模块：15个接口
- ✅ 安全功能模块：9个接口
- ✅ 系统监控模块：3个接口

**总计：83个API接口已全部更新完成！**

---

**文档版本** 3.1  
**最后更新** 2025-01-15  
**维护人员** LYS 
