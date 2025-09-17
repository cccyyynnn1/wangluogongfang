package main

import (
	"encoding/json" // json数据的编码解码，读取配置文件config.json
	"fmt"           // 打印
	"net"           // I/O接口，TCP连接和端口扫描
	"net/http"      // HTTP客户端和服务端实现，用于发送HTTP请求检测子域名
	"os"
	"strings" // 字符串处理
	"sync"    // 并发控制，多线程端口扫描
	"time"    // 获取相应时间，设置超时

	"github.com/gocolly/colly/v2" // 网爬框架，抓取网站信息
)

// Config 结构体用于存储配置信息
type Config struct {
	ApiKey     string   `json:"api_key"` // 将json中的值映射到结构体的字段
	Proxy      string   `json:"proxy"`
	Subdomains []string `json:"subdomains"`
}

// 加载配置文件
func loadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// 验证配置文件是否完整
func validateConfig(config *Config) bool {
	hasError := false
	if config.ApiKey == "" {
		fmt.Println("[!] API Key 未设置！")
		hasError = true
	}
	if config.Proxy == "" {
		fmt.Println("[!] 代理未设置！")
		hasError = true
	}
	if len(config.Subdomains) == 0 {
		fmt.Println("[!] 子域名列表为空！")
		hasError = true
	}
	return !hasError
}

// 分析网站基础信息
func analyzeWebsite(url string) {
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("[+] 状态码: %d\n", r.StatusCode)
		fmt.Printf("[+] 响应头: %v\n", r.Headers)
		fmt.Printf("[+] 服务器类型: %s\n", r.Headers.Get("Server"))
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		fmt.Printf("[+] 网站标题: %s\n", e.Text)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("[!] 错误: %s\n", err.Error())
	})

	// 发起请求
	err := c.Visit(url)
	if err != nil {
		fmt.Printf("[!] 请求失败: %s\n", err.Error())
	}
}

// 子域名爆破 返回发现的子域名列表
func subdomainBrute(domain, wordlistPath string) []string {
	// 读取字典文件
	data, err := os.ReadFile(wordlistPath)
	if err != nil {
		fmt.Printf("[!] 无法读取字典文件: %s\n", err.Error())
		return nil
	}

	subdomains := strings.Split(string(data), "\n")
	var foundSubdomains []string // 用于存储发现的有效子域名

	for _, sub := range subdomains {
		sub = strings.TrimSpace(sub) // 去除空白字符
		if sub == "" {
			continue // 跳过空行
		}
		url := "http://" + sub + "." + domain
		resp, err := http.Head(url)
		if err == nil && resp.StatusCode == 200 {
			fmt.Printf("[+] 发现子域名: %s\n", url)
			foundSubdomains = append(foundSubdomains, url) // 将完整的URL添加到列表
		}
		// 注意：这里我们没有处理 HTTPS 重定向。生产环境可能需要更复杂的逻辑。
	}
	return foundSubdomains // 返回发现的子域名列表
}

// 端口扫描
func portScan(target string, ports []int) {
	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", target, port)
			conn, err := net.Dial("tcp", address)
			if err == nil {
				fmt.Printf("[+] 端口 %d: 开放\n", port)
				conn.Close()
			}
		}(port)
	}
	wg.Wait()
}

// getIP 根据域名获取IP地址
func getIP(host string) string {
	// 使用 net.LookupHost 函数解析域名
	ips, err := net.LookupHost(host)
	if err != nil {
		return fmt.Sprintf("解析失败: %v", err)
	}
	// 通常一个域名会解析出一个或多个IP地址
	// 这里我们简单返回第一个IP，实际应用中可能需要处理多个IP的情况
	if len(ips) > 0 {
		return ips[0] // 返回第一个IP地址
	}
	return "未知" // 如果没有返回任何IP
}

// getProtocolAndInfo 获取URL的协议、HTTP版本和响应时间
func getProtocolAndInfo(url string) (protocol, httpVersion string, responseTime float64) {
	// 创建一个 HTTP 客户端
	client := &http.Client{
		// 可以设置超时时间等
		Timeout: 10 * time.Second,
	}

	// 发起 GET 请求（HEAD 有时可能不包含完整信息，GET 更全面）
	start := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		return "Error", "", 0.0
	}
	defer resp.Body.Close()

	// 判断协议
	if strings.HasPrefix(url, "https://") {
		protocol = "HTTPS"
	} else {
		protocol = "HTTP"
	}

	// 获取 HTTP 版本
	httpVersion = fmt.Sprintf("HTTP/%d.%d", resp.ProtoMajor, resp.ProtoMinor)

	// 计算响应时间
	duration := time.Since(start)
	responseTime = duration.Seconds()

	return protocol, httpVersion, responseTime
}

// getServerFingerprint 尝试获取服务器指纹
func getServerFingerprint(url string) string {
	c := colly.NewCollector()

	var serverFingerprint string

	// 从响应头的 Server 字段获取
	c.OnResponse(func(r *colly.Response) {
		serverHeader := r.Headers.Get("Server")
		if serverHeader != "" {
			serverFingerprint = serverHeader
			return // 找到就返回
		}

		// 如果 Server 头为空，尝试从 X-Powered-By 或其他头推测
		poweredBy := r.Headers.Get("X-Powered-By")
		if poweredBy != "" {
			serverFingerprint = "可能基于: " + poweredBy
			return
		}

		// 更复杂的指纹识别（示例）：
		// 可以检查特定的响应头、页面内容等
		// 例如，检查是否有特定的 meta 标签、JS 文件路径等
		// 这里只是一个简单的占位符，实际应用需要更复杂的规则库
	})

	// 防止因错误中断整个流程
	c.OnError(func(r *colly.Response, err error) {
		// 不做处理，serverFingerprint 保持默认值
	})

	// 必须发起请求才能触发 OnResponse
	err := c.Visit(url)
	if err != nil {
		// 如果请求失败，尝试直接用 HTTP 客户端获取头
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Head(url)
		if err == nil {
			defer resp.Body.Close()
			if server := resp.Header.Get("Server"); server != "" {
				serverFingerprint = server
				return serverFingerprint
			}
		}
		serverFingerprint = "获取失败"
	}

	return serverFingerprint
}

// generateSearchQueries 生成针对目标域名的常见搜索引擎查询语法
func generateSearchQueries(domain string) []string {
	queries := []string{}

	// 基础查询
	queries = append(queries, fmt.Sprintf("site:%s", domain))
	queries = append(queries, fmt.Sprintf("site:*.%s", domain))

	// 查找特定文件类型
	fileTypes := []string{"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "zip", "rar"}
	for _, ft := range fileTypes {
		queries = append(queries, fmt.Sprintf("site:%s filetype:%s", domain, ft))
	}

	// 查找可能的敏感页面
	sensitivePages := []string{
		"login", "admin", "dashboard", "wp-login", "cpanel", "phpmyadmin",
		"config", "backup", "sql", "env", "git", "svn",
	}
	for _, page := range sensitivePages {
		queries = append(queries, fmt.Sprintf("site:%s inurl:%s", domain, page))
	}

	// 查找可能的错误信息
	errorKeywords := []string{"error", "exception", "warning", "debug", "stack trace"}
	for _, keyword := range errorKeywords {
		queries = append(queries, fmt.Sprintf("site:%s intext:\"%s\"", domain, keyword))
	}

	// 查找特定的技术栈
	techStack := []string{"wordpress", "joomla", "drupal", "magento", "shopify", "laravel", "django", "spring"}
	for _, tech := range techStack {
		queries = append(queries, fmt.Sprintf("site:%s %s", domain, tech))
	}

	return queries
}

func main() {
	// 加载配置
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Printf("[!] 配置文件错误: %s\n", err.Error())
		return
	}

	// 验证配置
	if !validateConfig(config) {
		return
	}

	// 示例目标域名
	targetDomain := "baidu.com" // 这是主域名

	// 分析主域名信息
	fmt.Println("\n=== 分析主域名信息 ===")
	analyzeWebsite("https://" + targetDomain)

	// 子域名爆破，并获取结果
	fmt.Println("\n=== 子域名爆破 ===")
	foundUrls := subdomainBrute(targetDomain, "subdomains.txt")
	if len(foundUrls) == 0 {
		fmt.Println("[!] 未发现任何子域名。")
		return
	}
	/*
		if len(foundUrls) > 0 {
			// 假设目标主域名是固定的（如配置中的 targetDomain）
			targetDomain := "baidu.com" // 或从配置/输入获取

			fmt.Printf("\n=== 生成主域名 %s 的搜索引擎查询语法 ===\n", targetDomain)
			queries := generateSearchQueries(targetDomain)
			for _, q := range queries {
				fmt.Printf("[+] %s\n", q)
			}
		}
	*/
	// 对发现的每个子域名进行深度分析
	fmt.Println("\n=== 深度分析发现的子域名 ===")
	for _, url := range foundUrls {
		host := strings.TrimPrefix(url, "http://")
		host = strings.TrimPrefix(host, "https://")
		host = strings.Split(host, "/")[0]

		fmt.Printf("\n--- 深度分析: %s ---\n", url)
		fmt.Printf("[+] 主机名: %s\n", host)

		ip := getIP(host)
		fmt.Printf("[+] IP地址: %s\n", ip)

		// 获取协议和基础协议信息
		protocol, httpVersion, respTime := getProtocolAndInfo(url)
		fmt.Printf("[+] 协议: %s\n", protocol)
		fmt.Printf("[+] HTTP版本: %s\n", httpVersion)
		fmt.Printf("[+] 响应时间: %.3f 秒\n", respTime)

		// 获取服务器指纹
		fingerprint := getServerFingerprint(url)
		fmt.Printf("[+] 服务器指纹: %s\n", fingerprint)

		// 可以选择性地调用原始的 analyzeWebsite 来获取标题等
		// analyzeWebsite(url)
	}

	// 对发现的子域名进行端口扫描 (提取主机名)
	fmt.Println("\n=== 对发现的子域名进行端口扫描 ===")
	for _, url := range foundUrls {
		// 从URL中提取主机名 (例如: http://www.example.com -> www.example.com)
		host := strings.TrimPrefix(url, "http://")
		host = strings.TrimPrefix(host, "https://")
		// 移除可能的路径 (虽然Head请求不关心路径，但为了清晰)
		host = strings.Split(host, "/")[0]
		fmt.Printf("\n--- 扫描主机: %s ---\n", host)
		portScan(host, []int{21, 22, 80, 443, 8080})
	}
}
