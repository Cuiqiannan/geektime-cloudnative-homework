package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
)

func index(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("<h1>Welcome to Cloud Native</h1>"))
	// 03. 设置version
	os.Setenv("VERSION", "v0.0.1")
	version := os.Getenv("VERSION")
	w.Header().Set("VERSION", version)
	fmt.Printf("os version: %s \n", version)
	// 02. 将request中的header设置到response中
	for k, v := range r.Header {
		for _, vv := range v {
			fmt.Printf("Header key: %s, Header value: %s \n", k, v)
			w.Header().Set(k, vv)
		}
	}
	// 04. 记录日志并输出
	clientip := getCurrentIP(r)
	// fmt.Println(r.RemoteAddr)
	log.Printf("Success! Response code: %d", 200)
	log.Printf("Success! clientip: %s", clientip)
}

// 05. 健康检查的路由
func healthz(w http.ResponseWriter, r *http.Request) {
	// Fprintf: 格式化并输出到io.Writes而不是os.Stdout
	fmt.Fprintf(w, "working")
}

func getCurrentIP(r *http.Request) string {
	// 这里也可以通过X-Forwarded-For请求头的第一个值作为用户的ip
	// 但是要注意的是这两个请求头代表的ip都有可能是伪造的
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		// 当请求头不存在即不存在代理时直接获取ip
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return ip
}

// ClientIP尽最大努力实现获取客户端IP的算法
// 解析X-Real-IP和X-Forwarded-For以便于反向代理（nginx或haproxy）可以正常工作
func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func main() {
	mux := http.NewServeMux()
	// 06. debug
	mux.HandleFunc("/debug/pprof", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.HandleFunc("/", index)
	mux.HandleFunc("/healthz", healthz)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("start http failed, error: %s\n", err.Error())
	}
}

// test: http:127.0.0.1:8080
/*
[Running] go run "c:\Users\qiann\go\src\github.com\Cuiqiannan\geektime-cloudnative-homework\module2\httpserver\main.go"
os version: v0.0.1
Header key: Accept, Header value: [text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*\/*;q=0.8,application/signed-exchange;v=b3;q=0.9]
Header key: Sec-Fetch-Mode, Header value: [navigate]
Header key: Sec-Fetch-Dest, Header value: [document]
Header key: Sec-Ch-Ua, Header value: [" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"]
Header key: Sec-Ch-Ua-Mobile, Header value: [?0]
Header key: Sec-Ch-Ua-Platform, Header value: ["Windows"]
Header key: Upgrade-Insecure-Requests, Header value: [1]
Header key: Sec-Fetch-User, Header value: [?1]
Header key: Connection, Header value: [keep-alive]
Header key: User-Agent, Header value: [Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36]
Header key: Accept-Encoding, Header value: [gzip, deflate, br]
Header key: Accept-Language, Header value: [zh-CN,zh;q=0.9,en;q=0.8]
Header key: Sec-Fetch-Site, Header value: [none]
2022/02/27 22:42:53 Success! Response code: 200
2022/02/27 22:42:53 Success! clientip: 127.0.0.1
os version: v0.0.1
Header key: Connection, Header value: [keep-alive]
Header key: Sec-Ch-Ua, Header value: [" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"]
Header key: Sec-Ch-Ua-Platform, Header value: ["Windows"]
Header key: Sec-Fetch-Site, Header value: [same-origin]
Header key: Accept-Language, Header value: [zh-CN,zh;q=0.9,en;q=0.8]
Header key: Sec-Fetch-Dest, Header value: [image]
Header key: Referer, Header value: [http://127.0.0.1:8080/]
Header key: Accept-Encoding, Header value: [gzip, deflate, br]
Header key: Sec-Ch-Ua-Mobile, Header value: [?0]
Header key: User-Agent, Header value: [Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36]
Header key: Accept, Header value: [image/avif,image/webp,image/apng,image/svg+xml,image/*,*\/*;q=0.8]
Header key: Sec-Fetch-Mode, Header value: [no-cors]
2022/02/27 22:42:53 Success! Response code: 200
2022/02/27 22:42:53 Success! clientip: 127.0.0.1
os version: v0.0.1
Header key: Connection, Header value: [keep-alive]
Header key: Accept, Header value: [image/avif,image/webp,image/apng,image/svg+xml,image/*,*\/*;q=0.8]
Header key: Referer, Header value: [http://127.0.0.1:8080/healthz]
Header key: Accept-Language, Header value: [zh-CN,zh;q=0.9,en;q=0.8]
Header key: Sec-Ch-Ua, Header value: [" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"]
2022/02/27 22:43:00 Success! Response code: 200
2022/02/27 22:43:00 Success! clientip: 127.0.0.1
Header key: Sec-Ch-Ua-Mobile, Header value: [?0]
Header key: User-Agent, Header value: [Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36]
Header key: Sec-Ch-Ua-Platform, Header value: ["Windows"]
Header key: Sec-Fetch-Site, Header value: [same-origin]
Header key: Sec-Fetch-Mode, Header value: [no-cors]
Header key: Sec-Fetch-Dest, Header value: [image]
Header key: Accept-Encoding, Header value: [gzip, deflate, br]
os version: v0.0.1
Header key: Sec-Ch-Ua, Header value: [" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"]
Header key: User-Agent, Header value: [Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36]
Header key: Accept, Header value: [image/avif,image/webp,image/apng,image/svg+xml,image/*,*\/*;q=0.8]
Header key: Sec-Fetch-Site, Header value: [same-origin]
Header key: Sec-Fetch-Dest, Header value: [image]
2022/02/27 22:43:57 Success! Response code: 200
2022/02/27 22:43:57 Success! clientip: 127.0.0.1
Header key: Connection, Header value: [keep-alive]
Header key: Sec-Ch-Ua-Mobile, Header value: [?0]
Header key: Sec-Ch-Ua-Platform, Header value: ["Windows"]
Header key: Sec-Fetch-Mode, Header value: [no-cors]
Header key: Referer, Header value: [http://127.0.0.1:8080/healthz]
Header key: Accept-Encoding, Header value: [gzip, deflate, br]
Header key: Accept-Language, Header value: [zh-CN,zh;q=0.9,en;q=0.8]
os version: v0.0.1
Header key: Sec-Fetch-Site, Header value: [same-origin]
Header key: Sec-Fetch-Mode, Header value: [no-cors]
Header key: Referer, Header value: [http://127.0.0.1:8080/healthz]
Header key: Accept-Language, Header value: [zh-CN,zh;q=0.9,en;q=0.8]
2022/02/27 22:43:59 Success! Response code: 200
2022/02/27 22:43:59 Success! clientip: 127.0.0.1
Header key: Connection, Header value: [keep-alive]
Header key: Sec-Ch-Ua-Mobile, Header value: [?0]
Header key: Sec-Ch-Ua-Platform, Header value: ["Windows"]
Header key: Accept, Header value: [image/avif,image/webp,image/apng,image/svg+xml,image/*,*\/*;q=0.8]
Header key: Sec-Fetch-Dest, Header value: [image]
Header key: Accept-Encoding, Header value: [gzip, deflate, br]
Header key: Sec-Ch-Ua, Header value: [" Not A;Brand";v="99", "Chromium";v="98", "Google Chrome";v="98"]
Header key: User-Agent, Header value: [Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36]
*/
