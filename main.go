package main

import (
	"crypto/tls"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"os"
	"strconv"
)


func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/invite/send", SendMail)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8080", r))
}

func SendMail(w http.ResponseWriter, r *http.Request) {

	email, ok := r.URL.Query()["email"]

	if !ok || len(email[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	body := "<p>Email: " + string(email[0]) + "</p>"

	m := gomail.NewMessage()
	m.SetHeader("Subject", "Invite.")
	m.SetBody("text/html", body)

	go send(m)
	w.Write([]byte("Gorilla!\n"))
}

func send(message *gomail.Message) {

	from, ok := os.LookupEnv("FROM_EMAIL")

	if !ok {
		return
	}

	to, ok := os.LookupEnv("TO_EMAIL")

	if !ok {
		return
	}

	message.SetHeader("From", from)
	message.SetHeader("To", to)

	host, ok := os.LookupEnv("SMTP_HOST")

	if !ok {
		return
	}

	portStr, ok := os.LookupEnv("SMTP_PORT")

	if !ok {
		return
	}

	port, _ := strconv.Atoi(portStr)

	username, ok := os.LookupEnv("SMTP_USER")

	if !ok {
		return
	}

	password, ok := os.LookupEnv("SMTP_PASSWORD")

	if !ok {
		return
	}

	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(message); err != nil {
		log.Fatal(err.Error())
	}
}