package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

var emailNotifies = make(map[string]string)

func main() {
	loadEnv()
	urls := replaceSplitString(os.Getenv("URLS"), "\n", ",")

	for range time.Tick(time.Second * 10) {
		fmt.Println("Check Websites")
		for _, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				sendEmail(url)
			} else if !IsSuccessStatusCode(resp.StatusCode) {
				sendEmail(url, resp.Status)
			} else {
				fmt.Println(resp.StatusCode, ":", url)
			}
		}
		fmt.Println()
	}
}

func IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

func isValidToSendEmail(url string) bool {
	dateToday := time.Now().UTC().Format("2006-01-02")
	if emailNotifies[url] == "" || emailNotifies[url] < dateToday {
		emailNotifies[url] = dateToday
		return true
	}
	return false
}

func sendEmail(url string, statusCodeOptional ...string) {
	if !isValidToSendEmail(url) {
		return
	}

	username := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	smtpHost := os.Getenv("HOST")
	port, errConversion := strconv.Atoi(os.Getenv("PORT"))
	if errConversion != nil {
		port = 2525
	}
	from := os.Getenv("FROM")
	to := replaceSplitString(os.Getenv("TO"), "\n", ";")

	statusCodeMessage := ""
	if len(statusCodeOptional) > 0 {
		statusCode := statusCodeOptional[0]
		statusCodeMessage = " with Status Code " + statusCode + " : "
	}
	message := "[ALERT][Site is down] " + statusCodeMessage + url

	newMessage := gomail.NewMessage()
	newMessage.SetHeaders(map[string][]string{
		"From":    {from},
		"To":      to,
		"Subject": {message},
	})
	newMessage.SetBody("text/html", message)
	newDialer := gomail.NewDialer(smtpHost, port, username, password)
	err := newDialer.DialAndSend(newMessage)
	if err != nil {
		panic(err)
	}
}

func replaceSplitString(text string, replace string, split string) []string {
	text = strings.ReplaceAll(text, replace, "")
	return strings.Split(text, split)
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error from .env. Error: %s", err)
	}
}
