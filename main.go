package main

import (
	"fmt"
	"net/http"
	"net/smtp"
	"log"
	"os"
	"github.com/gorilla/mux"
	"github.com/go-ini/ini"
	"strconv"
	"time"
)

func send(body string, from string, to string, key string, test bool) {
	testTxt := ""
	if test{
		testTxt = "test"
	} else {
		testTxt = "real"
	}

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + time.Now().String() + "Alert! Notified!\n\n" +
		fmt.Sprintf(body, testTxt)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, key, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Printf("success")
}

func NotifyClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	//TODO: API refactor
        cfg, err := ini.Load("app.ini")
        if err != nil {
            fmt.Printf("Fail to read file: %v", err)
            os.Exit(1)
        }

        // Classic read of values, default section can be represented as empty string
        test, err := strconv.ParseBool(cfg.Section("settings").Key("test").String())
	if err != nil {
		log.Fatal("Error determining whether test or not")
	}

        from := cfg.Section("acct").Key("from").String()
	to := cfg.Section("acct").Key("to").String()
	key := cfg.Section("acct").Key("key").String()
	fmt.Println(from)
	fmt.Println(to)
	fmt.Println(key)
	fmt.Println(test)

	send("This has been a %s notification message", from, to, key, test)
	w.Write([]byte("Success!\n"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/send-message", NotifyClient).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", router))
}
