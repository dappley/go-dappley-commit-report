package email

import (
	"github.com/heesooh/go-dappley-commit-report/helper"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"strings"
	"bufio"
	"log"
	"fmt"
)

//Send commit-report email to the recipients specified in the "recipients.txt" file.
func SendEmail(emailBody string, branch string, committer string, senderEmail string, senderPasswd string) {
	var recipients []string

	file_byte, err := ioutil.ReadFile("recipients.txt")
	if err != nil { log.Fatal(err) }
	scanner := bufio.NewScanner(strings.NewReader(string(file_byte)))
	for scanner.Scan() {
		line := scanner.Text()
		if !helper.Valid_email(line) {
			fmt.Println("Invalid email address: \"" + line + "\"")
			continue
		}
		recipients = append(recipients, line)
	}

	//send the email
	mail := gomail.NewMessage()
	mail.SetHeader("From", senderEmail)
	if !helper.Contains(recipients, committer) {
		recipients = append(recipients, committer)
	}
	addresses := make([]string, len(recipients))
	for i, recipient := range recipients {
		addresses[i] = mail.FormatAddress(recipient, "")
	}
	mail.SetHeader("To", addresses...)
	mail.SetHeader("Subject", "Go-Dappley Commit Test Result - " + branch)
	mail.SetBody("text/html", emailBody)
	mail.Attach(branch + "/change.txt")
	mail.Attach(branch + "/log.txt")

	deliver := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPasswd)
	if err := deliver.DialAndSend(mail); err != nil {
		panic(err)
	}
}