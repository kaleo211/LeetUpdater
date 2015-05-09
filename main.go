package main

import (
    "bytes"
    "fmt"
    "net/http"
    "net/url"
    // "golang.org/x/net/html"
    // "io/ioutil"
    "io"
    "os"
    "bufio"
)

const user_agent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36"

// this would not REDIRECT!!!!! instead of http.Client
var transport = &http.Transport{}

// automatic redirect
var client = &http.Client{}

func leetcode() (cookie *http.Cookie) {
    leetcode_url := "https://leetcode.com/"

    req, _ := http.NewRequest("GET", leetcode_url, nil)
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
    // req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
    req.Header.Set("Accept-Language", "en-US,en;q=0.8")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("DNT", "1")
    req.Header.Set("Host", "leetcode.com")
    req.Header.Set("Referer", "https://leetcode.com/accounts/login/")
    req.Header.Set("User-Agent", user_agent)

    resp, _ := client.Do(req)
    cookie = resp.Cookies()[0]

    return
}


func login() (cookies []*http.Cookie) {

    reader := bufio.NewReader(os.Stdin)
    fmt.Print("password: ")
    password, _ := reader.ReadString('\n')

    data := url.Values{}
    data.Set("login", "kaleo211")
    data.Set("password", password)
    origin_cookie := leetcode()
    data.Set("csrfmiddlewaretoken", origin_cookie.Value)

    leetcode_login := "https://leetcode.com/accounts/login/"
    req, _ := http.NewRequest("POST", leetcode_login, bytes.NewBufferString(data.Encode()))
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
    req.Header.Set("Accept-Encoding", "gzip, deflate")
    req.Header.Set("Accept-Language", "en-US,en;q=0.8")
    req.Header.Set("Cache-Control", "max-age=0")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Content-Length", "92")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("DNT", "1")
    req.Header.Set("Host", "leetcode.com")
    req.Header.Set("Origin", "https://leetcode.com")
    req.Header.Set("Referer", "https://leetcode.com/accounts/login/")
    req.Header.Set("User-Agent", user_agent)
    req.Header.Set("Cookie", origin_cookie.Name+"="+origin_cookie.Value)

    resp, _ := transport.RoundTrip(req)
    cookies = resp.Cookies()

    return
}


func download_algorithm_html() {
    leetcode_algorithm := "https://leetcode.com/problemset/algorithms/"
    req, _ := http.NewRequest("GET", leetcode_algorithm, nil)

    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
    // Remove this header from req, cause GO COULD NOT decode one of them!!!!
    // req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
    req.Header.Set("Accept-Language", "en-US,en;q=0.8")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("DNT", "1")
    req.Header.Set("Host", "leetcode.com")
    req.Header.Set("Referer", "https://leetcode.com/accounts/login/")
    req.Header.Set("User-Agent", user_agent)

    cookies := login()
    for _, c := range cookies {
        req.AddCookie(c)
    }

    resp, _ := transport.RoundTrip(req)
    fmt.Println(resp.Header, resp.Status)
    defer resp.Body.Close()

    out, _ := os.Create("leetcode.html")
    defer out.Close()
    io.Copy(out, resp.Body)
}


func main() {

    download_algorithm_html()

}
