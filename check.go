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

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
	DDMMYYYY       = "2006-01-02"
)

var emailNotifies = make(map[string]string)

func main() {
	fmt.Println("Init Check Websites")
	fmt.Println("Load .env")
	loadEnv()
	urls := replaceSplitString(os.Getenv("URLS"), "\n", ",")
	seconds, errConversion := strconv.Atoi(os.Getenv("SECONDS_TO_CHECK"))
	if errConversion != nil {
		port = 300
	}

	fmt.Println("Start Process")
	for range time.Tick(time.Second * seconds) {
		fmt.Println("Check Websites - " + time.Now().UTC().Format(DDMMYYYYhhmmss))
		for _, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				sendEmail(url, err.Error())
				fmt.Println(err.Error(), ":", url)
				continue
			} else if !IsSuccessStatusCode(resp.StatusCode) {
				sendEmail(url, resp.Status)
			}
			fmt.Println(resp.Status, ":", url)
		}
		fmt.Println()
	}
}

func IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

func isValidToSendEmail(url string) bool {
	dateToday := time.Now().UTC().Format(DDMMYYYY)
	if emailNotifies[url] == "" || emailNotifies[url] < dateToday {
		emailNotifies[url] = dateToday
		return true
	}
	return false
}

func sendEmail(url string, error string) {
	if !isValidToSendEmail(url) {
		return
	}

	username := os.Getenv("USER_EMAIL")
	password := os.Getenv("PASSWORD")
	smtpHost := os.Getenv("HOST")
	port, errConversion := strconv.Atoi(os.Getenv("PORT"))
	if errConversion != nil {
		port = 2525
	}
	from := os.Getenv("FROM")
	to := replaceSplitString(os.Getenv("TO"), "\n", ";")
	message := "[ALERT WEBSITE] " + url

	newMessage := gomail.NewMessage()
	newMessage.SetHeaders(map[string][]string{
		"From":    {from},
		"To":      to,
		"Subject": {message},
	})
	newMessage.SetBody("text/html", "Error: "+error+" : "+url)
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
