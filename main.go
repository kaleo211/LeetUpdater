package main

import (
    "bytes"
    "fmt"
    "net/http"
    "net/url"
    "golang.org/x/net/html"
    // "io/ioutil"
    "io"
    "strings"
    "os"
    "bufio"
)

const user_agent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36"

// this would not REDIRECT!!!!! instead of http.Client
var transport = &http.Transport{}

// automatic redirect
var client = &http.Client{}

// visit leetcode home page to get cookie
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


func leetcode_login() (cookies []*http.Cookie) {

    // read password from command line
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("password: ")
    password, _ := reader.ReadString('\n')
    password = strings.Trim(password, "\r\n"+string(0))

    // set post form data
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

    cookies := leetcode_login()
    for _, c := range cookies {
        req.AddCookie(c)
    }

    resp, _ := transport.RoundTrip(req)
    // fmt.Println(resp.Header, resp.Status)
    defer resp.Body.Close()

    out, _ := os.Create("leetcode.html")
    defer out.Close()
    io.Copy(out, resp.Body)
}


func parse_html() (solved, total int) {
    solved, total = 0, 0

    file, err := os.Open("leetcode.html")
    if err!=nil {
        fmt.Println(err)
        return
    }
    z := html.NewTokenizer(file)

    for {
        tt := z.Next()
        if tt==html.ErrorToken {
            return
        }
        bytes, has_attr := z.TagName()

        if "span"==string(bytes) && has_attr {
            for {
                k, v, more_attr := z.TagAttr()

                if string(k)=="class" {
                    if string(v)=="ac" {
                        solved += 1
                        total += 1
                    }
                    if string(v)=="None" || string(v)=="notac" {
                        total += 1
                    }
                }
                if !more_attr {
                    break
                }
            }
        }
    }
}


// visit github page to get cookie
func github() (cookies []*http.Cookie, auth_token string) {
    github := "https://github.com/login"

    req, _ := http.NewRequest("GET", github, nil)
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
    // req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
    req.Header.Set("Accept-Language", "en-US,en;q=0.8")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("DNT", "1")
    req.Header.Set("Host", "github.com")
    req.Header.Set("Referer", "https://github.com/")
    req.Header.Set("User-Agent", user_agent)

    resp, err := client.Do(req)
    if err!=nil {
        fmt.Println(err)
    }
    cookies = resp.Cookies()

    z := html.NewTokenizer(resp.Body)
    for {
        tt := z.Next()
        if tt==html.ErrorToken {
            return
        }
        bytes, has_attr := z.TagName()

        if "meta"==string(bytes) && has_attr {
            for {
                k, v, more_attr := z.TagAttr()
                if string(k)=="content" && strings.HasSuffix(string(v), "==") {
                    auth_token = string(v)
                }
                if !more_attr {
                    break
                }
            }
        }
    }

    return
}


func github_login() (cookies []*http.Cookie, token string) {
    auth_cookies, auth_token := github()


    // read password from command line
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("password: ")
    password, _ := reader.ReadString('\n')
    password = strings.Trim(password, "\r\n"+string(0))


    // set post form data
    data := url.Values{}
    data.Set("utf8", "✓")
    data.Set("login", "kaleo211")
    data.Set("password", password)
    data.Set("authenticity_token", auth_token)

    github_login := "https://github.com/session"
    req, _ := http.NewRequest("POST", github_login, bytes.NewBufferString(data.Encode()))
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
    req.Header.Set("Accept-Encoding", "gzip, deflate")
    req.Header.Set("Accept-Language", "en-US,en;q=0.8")
    req.Header.Set("Cache-Control", "max-age=0")
    req.Header.Set("Connection", "keep-alive")
    // req.Header.Set("Content-Length", "174")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("DNT", "1")
    req.Header.Set("Host", "github.com")
    req.Header.Set("Origin", "https://github.com")
    req.Header.Set("Referer", "https://github.com/")
    req.Header.Set("User-Agent", user_agent)

    for _, c := range auth_cookies {
        req.AddCookie(c)
    }

    resp, _ := transport.RoundTrip(req)
    defer resp.Body.Close()
    cookies = resp.Cookies()

    z := html.NewTokenizer(resp.Body)
    for {
        tt := z.Next()
        if tt==html.ErrorToken {
            return
        }
        bytes, has_attr := z.TagName()

        if "meta"==string(bytes) && has_attr {
            for {
                k, v, more_attr := z.TagAttr()
                if string(k)=="content" && strings.HasSuffix(string(v), "==") {
                    token = string(v)
                }
                if !more_attr {
                    break
                }
            }
        }
    }
    fmt.Println(resp.Status, token)

    return
}

func update_description(description string, cookies []*http.Cookie, token string) {
    github := "https://github.com/kaleo211/Leetcode"
    req, _ := http.NewRequest("GET", github, nil)

    for _, c := range cookies {
        req.AddCookie(c)
    }

    resp, _ := client.Do(req)
    cookies = resp.Cookies()

    doc, _ := html.Parse(resp.Body)

    var f func(*html.Node)

    f = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "form" {
            for _, attr := range n.Attr {
                if attr.Key=="action" && attr.Val=="/kaleo211/Leetcode/settings/update_meta" {
                    child := n.FirstChild.FirstChild.NextSibling.NextSibling
                    for _, a := range child.Attr {
                        if a.Key=="value" {
                            token = a.Val
                            return
                        }
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
    }
    f(doc)


    fmt.Println(resp.Status, token)


    // set post form data
    data := url.Values{}
    data.Set("utf8", "✓")
    data.Set("_method", "put")
    data.Set("repo_description", description)
    data.Set("repo_homepage", "")
    data.Set("authenticity_token", token)

    github_project := "https://github.com/kaleo211/Leetcode/settings/update_meta"
    req, _ = http.NewRequest("POST", github_project, bytes.NewBufferString(data.Encode()))

    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
    req.Header.Set("Accept-Encoding", "gzip, deflate")
    req.Header.Set("Accept-Language", "en-US,en;q=0.8")
    req.Header.Set("Cache-Control", "max-age=0")
    req.Header.Set("Connection", "keep-alive")
    // req.Header.Set("Content-Length", "240")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("DNT", "1")
    req.Header.Set("Host", "github.com")
    req.Header.Set("Origin", "https://github.com")
    req.Header.Set("Referer", "https://github.com/kaleo211/Leetcode")
    req.Header.Set("User-Agent", user_agent)

    coo := "_gat=1; logged_in=yes; dotcom_user=kaleo211; _gh_sess="
    kie := "; tz=America/New_York; user_session="
    for _, c := range cookies {
        if c.Name=="user_session" {
            kie += c.Value
        } else if c.Name=="_gh_sess" {
            coo += c.Value
        }
    }
    req.Header.Set("Cookie", coo+kie)

    // fmt.Println(coo+kie)

    resp, _ = transport.RoundTrip(req)
    defer resp.Body.Close()

    fmt.Println(resp.Status)

    cookies = resp.Cookies()

    return
}



func main() {

    download_algorithm_html()

    solved, total := parse_html()
    description := fmt.Sprintf("Solutions to LeetCode %d/%d", solved, total)

    cookies, token := github_login()

    update_description(description, cookies, token)
}
