package main

import (
    "bytes"
    "net/http"
    "net/url"
    "golang.org/x/net/html"
    "io/ioutil"
)


// visit leetcode home page to get cookie
func Leetcode() (cookie *http.Cookie) {
    leetcode_url := "https://leetcode.com/"

    req, _ := http.NewRequest("GET", leetcode_url, nil)

    resp, _ := client.Do(req)
    cookie = resp.Cookies()[0]

    return
}


func LeetcodeLogin() (cookies []*http.Cookie) {

    // set post form data
    data := url.Values{}
    data.Set("login", login)
    data.Set("password", password)
    origin_cookie := Leetcode()
    data.Set("csrfmiddlewaretoken", origin_cookie.Value)

    leetcode_login := "https://leetcode.com/accounts/login/"
    req, _ := http.NewRequest("POST", leetcode_login, bytes.NewBufferString(data.Encode()))

    req.Header.Set("Referer", "https://leetcode.com/accounts/login/")
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Cookie", origin_cookie.Name+"="+origin_cookie.Value)

    resp, _ := transport.RoundTrip(req)
    cookies = resp.Cookies()

    if resp.Status=="302 FOUND" {
        logger.Println("login into leetcode successfully.")
    } else {
        logger.Fatalln("failed login into leetcode.")
    }

    return
}


func AlgorithmPage() (bytes []byte) {
    leetcode_algorithm := "https://leetcode.com/problemset/algorithms/"
    req, _ := http.NewRequest("GET", leetcode_algorithm, nil)

    cookies := LeetcodeLogin()
    for _, c := range cookies {
        req.AddCookie(c)
    }

    resp, _ := transport.RoundTrip(req)
    defer resp.Body.Close()
    bytes, _ = ioutil.ReadAll(resp.Body)

    return
}


func Progress() (solved, total int) {
    solved, total = 0, 0

    reader := bytes.NewReader(AlgorithmPage())

    // low level traverse
    z := html.NewTokenizer(reader)
    for {
        tt := z.Next()
        if tt==html.ErrorToken {
            logger.Printf("solved %d / %d problems.", solved, total)
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
