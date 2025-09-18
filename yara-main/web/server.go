package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// 获取当前可执行文件所在目录
	execDir, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDir = filepath.Dir(execDir)
	
	// 如果当前目录是web目录，使用当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	
	// 检查当前目录是否包含index.html
	if _, err := os.Stat(filepath.Join(currentDir, "index.html")); err == nil {
		// 当前目录包含index.html，使用当前目录
		log.Printf("使用当前目录: %s", currentDir)
	} else {
		// 尝试在web子目录中查找
		webDir := filepath.Join(currentDir, "web")
		if _, err := os.Stat(filepath.Join(webDir, "index.html")); err == nil {
			currentDir = webDir
			log.Printf("使用web子目录: %s", currentDir)
		} else {
			log.Fatal("找不到index.html文件")
		}
	}

	// 设置静态文件目录
	fs := http.FileServer(http.Dir(currentDir))

	// 设置路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 默认返回index.html
		if r.URL.Path == "/" {
			indexPath := filepath.Join(currentDir, "index.html")
			log.Printf("提供index.html: %s", indexPath)
			http.ServeFile(w, r, indexPath)
			return
		}
		fs.ServeHTTP(w, r)
	})

	// 启动服务器
	port := ":3001"
	log.Printf("Web界面服务器启动在 http://localhost%s", port)
	log.Printf("静态文件目录: %s", currentDir)
	log.Printf("请确保Yara安全服务运行在 http://localhost:8081")
	log.Fatal(http.ListenAndServe(port, nil))
}
