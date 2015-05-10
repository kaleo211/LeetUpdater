package main

import (
    "bytes"
    "fmt"
    "net/http"
    "net/url"
    "golang.org/x/net/html"
    "strings"
)

// visit github page to get cookie
func Github() (cookies []*http.Cookie, auth_token string) {
    github := "https://github.com/login"

    req, _ := http.NewRequest("GET", github, nil)

    resp, err := client.Do(req)
    if err!=nil {
        fmt.Println(err)
    }
    cookies = resp.Cookies()

    // low level traverse html
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


func GithubLogin() (cookies []*http.Cookie) {
    auth_cookies, auth_token := Github()

    // set post form data
    data := url.Values{}
    data.Set("utf8", "✓")
    data.Set("login", login)
    data.Set("password", password)
    data.Set("authenticity_token", auth_token)

    github_login := "https://github.com/session"
    req, _ := http.NewRequest("POST", github_login, bytes.NewBufferString(data.Encode()))

    for _, c := range auth_cookies {
        req.AddCookie(c)
    }

    resp, _ := transport.RoundTrip(req)
    defer resp.Body.Close()

    if resp.Status=="302 Found" {
        logger.Println("login into github successfully.")
    } else {
        logger.Fatalln("failed login into github.")
    }

    cookies = resp.Cookies()

    return
}


// every action has his own token in github
func GetFormToken(origin_cookies []*http.Cookie) (cookies []*http.Cookie, token string) {
    github := "https://github.com/kaleo211/Leetcode"
    req, _ := http.NewRequest("GET", github, nil)

    for _, c := range origin_cookies {
        req.AddCookie(c)
    }

    resp, _ := client.Do(req)
    cookies = resp.Cookies()
    doc, _ := html.Parse(resp.Body)

    var f func(*html.Node)
    // high level traverse html
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
    return
}


func UpdateDescription(description string) {

    origin_cookies := GithubLogin()

    cookies, token := GetFormToken(origin_cookies)

    // set post form data
    data := url.Values{}
    data.Set("utf8", "✓")
    data.Set("_method", "put")
    data.Set("repo_description", description)
    data.Set("repo_homepage", "")
    data.Set("authenticity_token", token)

    github_project := "https://github.com/kaleo211/Leetcode/settings/update_meta"
    req, _ := http.NewRequest("POST", github_project, bytes.NewBufferString(data.Encode()))

    cookie := "_gat=1; logged_in=yes; dotcom_user=kaleo211"
    for _, c := range cookies {
        cookie += "; "+c.Name+"="+c.Value
    }
    req.Header.Set("Cookie", cookie)

    resp, _ := transport.RoundTrip(req)
    defer resp.Body.Close()

    if resp.Status=="302 Found" {
        logger.Println("update description for Leetcode project successfully.")
    } else {
        logger.Fatalln("failed update description.")
    }

    return
}
