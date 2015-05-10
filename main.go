package main

import (
    "fmt"
    "net/http"
    "log"
    "os"
    "time"
    "golang.org/x/crypto/ssh/terminal"
)

const sleeping = time.Duration(6)*time.Hour
const login = "kaleo211"
const description = "Solutions to LeetCode %d/%d  - update automatically"

var password string

// this would not REDIRECT!!!!!
var transport = &http.Transport{}

// automatic redirect
var client = &http.Client{}

var logger = log.New(os.Stdout, "", log.Ltime | log.Lshortfile)


func Update() {
    solved, total := Progress()

    descrip := fmt.Sprintf(description, solved, total)
    logger.Println(descrip)

    UpdateDescription(descrip)
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

    // visit 127.0.0.1:8080 could invoke update
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}
