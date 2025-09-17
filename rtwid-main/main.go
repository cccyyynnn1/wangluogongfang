package main

import (
	"encoding/json" // json数据的编码解码，读取配置文件config.json
	"fmt"           // 打印
	"io"
	"net"      // I/O接口，TCP连接和端口扫描
	"net/http" // HTTP客户端和服务端实现，用于发送HTTP请求检测子域名
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

// 通过 crt.sh 查询子域名
func crtShSubdomainEnum(domain string) []string {
	url := fmt.Sprintf("https://crt.sh/?q=%s&output=json", domain)

	// 创建客户端
	client := &http.Client{
		Timeout: 30 * time.Second, // 设置超时时间30秒
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[!] 创建请求失败: %s\n", err.Error())
		return nil
	}

	// 设置更真实的浏览器 User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[!] 请求 crt.sh 失败: %s\n", err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("[!] crt.sh 返回状态码: %d\n", resp.StatusCode)
		return nil
	}

	// 检查响应内容类型
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "text/json") {
		fmt.Printf("[!] crt.sh 返回非JSON内容: %s\n", contentType)
		return nil
	}

	// 尝试读取响应体内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[!] 读取响应体失败: %s\n", err.Error())
		return nil
	}

	// 解析 JSON 响应
	var results []struct {
		NameValue  string `json:"name_value"`
		CommonName string `json:"common_name"`
	}

	if err := json.Unmarshal(body, &results); err != nil {
		fmt.Printf("[!] 解析 crt.sh 响应失败: %s\n", err.Error())
		fmt.Printf("[!] 响应内容前200字符: %s\n", string(body[:min(200, len(body))]))
		return nil
	}

	// 使用 map 去重
	subdomainSet := make(map[string]bool)
	var subdomains []string

	// 处理 name_value 字段
	for _, result := range results {
		// name_value 可能包含多个域名，用换行符分隔
		names := strings.Split(result.NameValue, "\n")
		for _, name := range names {
			name = strings.TrimSpace(name)
			// 过滤掉通配符域名和确保是有效域名
			if name != "" && !strings.HasPrefix(name, "*") && (strings.HasSuffix(name, "."+domain) || name == domain) {
				if !subdomainSet[name] {
					subdomainSet[name] = true
					subdomains = append(subdomains, name)
				}
			}
		}
	}

	// 如果 name_value 没有结果，尝试 common_name
	if len(subdomains) == 0 {
		for _, result := range results {
			name := strings.TrimSpace(result.CommonName)
			if name != "" && !strings.HasPrefix(name, "*") && (strings.HasSuffix(name, "."+domain) || name == domain) {
				if !subdomainSet[name] {
					subdomainSet[name] = true
					subdomains = append(subdomains, name)
				}
			}
		}
	}

	fmt.Printf("[+] crt.sh 查询发现 %d 个子域名\n", len(subdomains))
	return subdomains
}

// 辅助函数：取两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// virusTotalDomainInfo 通过 VirusTotal API 查询域名的详细信息
func virusTotalDomainInfo(domain, apiKey string) {
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%s", domain)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[!] 创建VirusTotal请求失败: %s\n", err.Error())
		return
	}

	req.Header.Set("x-apikey", apiKey)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[!] 请求VirusTotal失败: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		fmt.Println("[!] VirusTotal API密钥无效或已过期")
		return
	}

	if resp.StatusCode == 429 {
		fmt.Println("[!] VirusTotal API请求频率超限，请稍后再试")
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("[!] VirusTotal返回状态码: %d\n", resp.StatusCode)
		return
	}

	var result struct {
		Data struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				LastAnalysisStats struct {
					Harmless   int `json:"harmless"`
					Malicious  int `json:"malicious"`
					Suspicious int `json:"suspicious"`
					Undetected int `json:"undetected"`
					Timeout    int `json:"timeout"`
				} `json:"last_analysis_stats"`
				LastAnalysisDate int64             `json:"last_analysis_date"`
				Reputation       int               `json:"reputation"`
				Categories       map[string]string `json:"categories"`
				LastDnsRecords   []struct {
					Type  string `json:"type"`
					Value string `json:"value"`
					TTL   int    `json:"ttl"`
				} `json:"last_dns_records"`
				LastHttpResponseHeaders    map[string]string `json:"last_http_response_headers"`
				LastHttpResponseStatusCode int               `json:"last_http_response_status_code"`
				CreationDate               int64             `json:"creation_date"`
			} `json:"attributes"`
		} `json:"data"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[!] 读取VirusTotal响应失败: %s\n", err.Error())
		return
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("[!] 解析VirusTotal响应失败: %s\n", err.Error())
		return
	}

	// 打印域名基本信息
	fmt.Printf("\n=== VirusTotal 域名详细信息 ===\n")
	fmt.Printf("域名: %s\n", result.Data.ID)
	fmt.Printf("类型: %s\n", result.Data.Type)

	// 打印分析统计信息
	fmt.Printf("\n--- 安全分析统计 ---\n")
	fmt.Printf("无害: %d\n", result.Data.Attributes.LastAnalysisStats.Harmless)
	fmt.Printf("恶意: %d\n", result.Data.Attributes.LastAnalysisStats.Malicious)
	fmt.Printf("可疑: %d\n", result.Data.Attributes.LastAnalysisStats.Suspicious)
	fmt.Printf("未检测: %d\n", result.Data.Attributes.LastAnalysisStats.Undetected)
	fmt.Printf("超时: %d\n", result.Data.Attributes.LastAnalysisStats.Timeout)

	// 打印声誉和分类信息
	fmt.Printf("\n--- 声誉和分类 ---\n")
	fmt.Printf("声誉分数: %d\n", result.Data.Attributes.Reputation)

	if len(result.Data.Attributes.Categories) > 0 {
		fmt.Printf("分类:\n")
		for provider, category := range result.Data.Attributes.Categories {
			fmt.Printf("  %s: %s\n", provider, category)
		}
	}

	// 打印DNS记录
	if len(result.Data.Attributes.LastDnsRecords) > 0 {
		fmt.Printf("\n--- 最新DNS记录 ---\n")
		for _, record := range result.Data.Attributes.LastDnsRecords {
			fmt.Printf("  %s: %s (TTL: %d)\n", record.Type, record.Value, record.TTL)
		}
	}

	// 打印HTTP响应信息
	if result.Data.Attributes.LastHttpResponseStatusCode != 0 {
		fmt.Printf("\n--- HTTP响应信息 ---\n")
		fmt.Printf("状态码: %d\n", result.Data.Attributes.LastHttpResponseStatusCode)

		if len(result.Data.Attributes.LastHttpResponseHeaders) > 0 {
			fmt.Printf("响应头:\n")
			for key, value := range result.Data.Attributes.LastHttpResponseHeaders {
				fmt.Printf("  %s: %s\n", key, value)
			}
		}
	}

	// 打印时间信息
	if result.Data.Attributes.LastAnalysisDate != 0 {
		analysisTime := time.Unix(result.Data.Attributes.LastAnalysisDate, 0)
		fmt.Printf("\n--- 时间信息 ---\n")
		fmt.Printf("最后分析时间: %s\n", analysisTime.Format("2006-01-02 15:04:05"))
	}

	if result.Data.Attributes.CreationDate != 0 {
		creationTime := time.Unix(result.Data.Attributes.CreationDate, 0)
		fmt.Printf("创建时间: %s\n", creationTime.Format("2006-01-02 15:04:05"))
	}
}

// virusTotalSubdomainEnum 通过 VirusTotal API 查询子域名
func virusTotalSubdomainEnum(domain, apiKey string) []string {
	url := fmt.Sprintf("https://www.virustotal.com/api/v3/domains/%s/subdomains", domain)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("[!] 创建VirusTotal请求失败: %s\n", err.Error())
		return nil
	}

	req.Header.Set("x-apikey", apiKey)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[!] 请求VirusTotal失败: %s\n", err.Error())
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		fmt.Println("[!] VirusTotal API密钥无效或已过期")
		return nil
	}

	if resp.StatusCode == 429 {
		fmt.Println("[!] VirusTotal API请求频率超限，请稍后再试")
		return nil
	}

	if resp.StatusCode != 200 {
		fmt.Printf("[!] VirusTotal返回状态码: %d\n", resp.StatusCode)
		return nil
	}

	var result struct {
		Data []struct {
			ID         string `json:"id"`
			Type       string `json:"type"`
			Attributes struct {
				LastAnalysisStats struct {
					Harmless   int `json:"harmless"`
					Malicious  int `json:"malicious"`
					Suspicious int `json:"suspicious"`
					Undetected int `json:"undetected"`
					Timeout    int `json:"timeout"`
				} `json:"last_analysis_stats"`
				Reputation int `json:"reputation"`
			} `json:"attributes"`
		} `json:"data"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[!] 读取VirusTotal响应失败: %s\n", err.Error())
		return nil
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("[!] 解析VirusTotal响应失败: %s\n", err.Error())
		return nil
	}

	// 打印子域名详细信息
	fmt.Printf("\n=== VirusTotal 子域名详细信息 ===\n")
	var subdomains []string
	subdomainSet := make(map[string]bool)

	for _, item := range result.Data {
		if item.Type == "domain" && item.ID != "" {
			// 验证是否为目标域名的子域名
			if strings.HasSuffix(item.ID, "."+domain) || item.ID == domain {
				if !subdomainSet[item.ID] {
					subdomainSet[item.ID] = true
					subdomains = append(subdomains, item.ID)

					// 打印子域名详细信息
					fmt.Printf("\n子域名: %s\n", item.ID)
					fmt.Printf("  类型: %s\n", item.Type)
					fmt.Printf("  声誉分数: %d\n", item.Attributes.Reputation)
					fmt.Printf("  安全扫描结果:\n")
					fmt.Printf("    无害: %d\n", item.Attributes.LastAnalysisStats.Harmless)
					fmt.Printf("    恶意: %d\n", item.Attributes.LastAnalysisStats.Malicious)
					fmt.Printf("    可疑: %d\n", item.Attributes.LastAnalysisStats.Suspicious)
					fmt.Printf("    未检测: %d\n", item.Attributes.LastAnalysisStats.Undetected)
					fmt.Printf("    超时: %d\n", item.Attributes.LastAnalysisStats.Timeout)
				}
			}
		}
	}

	fmt.Printf("\n[+] VirusTotal查询发现 %d 个子域名\n", len(subdomains))
	return subdomains
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

	/*
		// 在子域名爆破部分，替换或添加 crt.sh 查询
		fmt.Println("\n=== 通过 crt.sh 查询子域名 ===")
		crtSubdomains := crtShSubdomainEnum(targetDomain)
		if len(crtSubdomains) > 0 {
			for _, subdomain := range crtSubdomains {
				fmt.Printf("[+] 发现子域名: %s\n", subdomain)
			}
		} else {
			fmt.Println("[!] crt.sh 未发现任何子域名。")
		}
	*/

	// 在 main 函数中替换之前的查询方法
	fmt.Println("\n=== 通过 VirusTotal 查询域名详细信息 ===")
	// 需要从配置中获取API密钥
	config, err = loadConfig("config.json")
	if err != nil {
		fmt.Printf("[!] 配置文件错误: %s\n", err.Error())
		return
	}

	virusTotalApiKey := config.ApiKey // 假设API密钥存储在配置文件中

	// 查询域名详细信息
	virusTotalDomainInfo(targetDomain, virusTotalApiKey)

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
