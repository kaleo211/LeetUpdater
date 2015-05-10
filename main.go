package main

import (
    "fmt"
    "net/http"
    "log"
    "os"
    "time"
    "golang.org/x/crypto/ssh/terminal"
)

const user_agent string = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36"

// this would not REDIRECT!!!!! instead of http.Client
var transport = &http.Transport{}

// automatic redirect
var client = &http.Client{}

var logger = log.New(os.Stdout, "", log.Ltime | log.Lshortfile)

var password string

const sleeping = time.Duration(6)*time.Hour


func Update() {
    solved, total := Progress()

    description := fmt.Sprintf("Solutions to LeetCode %d/%d  - update automatically", solved, total)
    logger.Println(description)

    cookies, token := GithubLogin()

    UpdateDescription(description, cookies, token)
    logger.Println("------------------E-N-D----------------------")
}


func handler(w http.ResponseWriter, r *http.Request) {
    Update()
    fmt.Fprintln(w, "update successfully.")
}

func main() {
    fmt.Print("password ")
    raw, _ := terminal.ReadPassword(1)
    password = string(raw)

    go func() {
        for {
            Update()
            time.Sleep(sleeping)
        }
    }()

    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
