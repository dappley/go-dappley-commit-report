package main

import (
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/mail"
	"strings"
	"bufio"
	"log"
	"fmt"
)

func sendEmail(emailBody string, branch string, committer string, senderEmail string, senderPasswd string) {
	var recipients []string

	file_byte, err := ioutil.ReadFile("recipients.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(file_byte)))
	for scanner.Scan() {
		line := scanner.Text()
		if !valid_email(line) {
			fmt.Println("Invalid email address: \"" + line + "\"")
			continue
		}
		recipients = append(recipients, line)
	}

	//send the email
	mail := gomail.NewMessage()
	mail.SetHeader("From", senderEmail)
	if !contains(recipients, committer) {
		recipients = append(recipients, committer)
	}
	addresses := make([]string, len(recipients))
	for i, recipient := range recipients {
		addresses[i] = mail.FormatAddress(recipient, "")
	}
	mail.SetHeader("To", addresses...)
	//mail.SetAddressHeader("Cc", "dan@example.com", "Dan")
	mail.SetHeader("Subject", "Go-Dappley Commit Test Result - " + branch)
	mail.SetBody("text/html", emailBody)
	mail.Attach(branch + "/change.txt")
	mail.Attach(branch + "/log.txt")

	deliver := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, senderPasswd)

	if err := deliver.DialAndSend(mail); err != nil {
		panic(err)
	}
}

//Checks if slice contains the given value
func contains(slice []string, val string) bool {
	for _, elem := range slice {
		if elem == val {
			return true
		}
	}
	return false
}

//Checks the validity of the email address
func valid_email(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}