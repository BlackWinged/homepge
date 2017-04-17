package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

func main() {
	routes()
	http.ListenAndServe(":5000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("html/index.html")
	t.Execute(w, nil)
}

func routes() {
	http.HandleFunc("/", handler)
	http.Handle("/content/", http.StripPrefix("/content/", http.FileServer(http.Dir("/content"))))
	http.Handle("/js/", http.FileServer(http.Dir("content")))
	http.Handle("/css/", http.FileServer(http.Dir("content")))
	http.Handle("/images/", http.FileServer(http.Dir("content")))
	http.Handle("/vendor/", http.FileServer(http.Dir("content")))

	http.HandleFunc("/contact/", mailer)
}

func processRequest(r *http.Request) (string, error) {
	var body string
	if r.Method == "POST" {
		r.ParseForm()
		body = "Reply-To: " + r.Form.Get("email")
		body += "<p>From: " + r.Form.Get("name") + "</p><br/>"
		body += "<p>Phone: " + r.Form.Get("phone") + "</p><br/>"
		body += "<p>Phone: " + r.Form.Get("message") + "</p><br/>"
	} else {
		return "", errors.New("invalid method")
	}
	return body, nil
}

func mailer(w http.ResponseWriter, r *http.Request) {
	password := os.Getenv("SMTP_PASSWORD")
	// domain := os.Getenv("DOMAIN")
	username := os.Getenv("SMTP_USERNAME")
	// smtpClient := os.Getenv("SMTP_ADDRESS")
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: Mail from your site\n"
	auth := smtp.PlainAuth(
		"",
		username,
		password,
		"smtp.gmail.com",
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	body, errProcessing := processRequest(r)

	if errProcessing == nil {
		log.Print(errProcessing)
		errSending := smtp.SendMail(
			"smtp.gmail.com:587",
			auth,
			"sender@example.org",
			[]string{"lovro.gamulin@gmail.com"},
			[]byte(subject+mime+body),
		)
		if errSending != nil {
			log.Print(errSending)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened! Pass:"))
		}
	}
	w.WriteHeader(http.StatusOK)
}
