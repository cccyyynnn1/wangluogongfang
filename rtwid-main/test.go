package main

/*package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	_ "os"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

// Config 结构体用于存储配置信息
type Config struct {
	ApiKey     string   `json:"api_key"`
	Proxy      string   `json:"proxy"`
	Subdomains []string `json:"subdomains"`
}

// 加载配置文件
func loadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
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

// 子域名爆破
func subdomainBrute(domain, wordlistPath string) {
	// 读取字典文件
	data, err := ioutil.ReadFile(wordlistPath)
	if err != nil {
		fmt.Printf("[!] 无法读取字典文件: %s\n", err.Error())
		return
	}

	subdomains := strings.Split(string(data), "\n")

	for _, sub := range subdomains {
		url := "http://" + strings.TrimSpace(sub) + "." + domain
		resp, err := http.Head(url)
		if err == nil && resp.StatusCode == 200 {
			fmt.Printf("[+] 发现子域名: %s\n", url)
		}
	}
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

	// 示例目标URL
	targetUrl := "https://baidu.com"
	domain := strings.Replace(targetUrl, "https://", "", 1)

	// 分析网站信息
	analyzeWebsite(targetUrl)

	// 子域名爆破
	subdomainBrute(domain, "subdomains.txt")

	// 端口扫描
	portScan(domain, []int{21, 22, 80, 443, 8080})
}
*/
