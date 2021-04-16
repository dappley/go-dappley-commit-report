package main

import (
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"strings"
	"bufio"
	"flag"
	"fmt"
)

func main() {
	var change, testResult, commitInfo, sender, senderPasswd string
	flag.StringVar(&change,       "change",       "change.txt",        "Changes made in the new commit")
	flag.StringVar(&testResult,   "testResult",   "log.txt",           "Test result file")
	flag.StringVar(&commitInfo,   "commitInfo",   "commit info",       "default commit info")
	flag.StringVar(&sender,       "sender",       "sender_username@example.com", "Email address of the sender")
	flag.StringVar(&senderPasswd, "senderPasswd", "PASSWORD",          "Password of the sender's email address.")
	flag.Parse()

	committer, emailContent, fail_exists := compose(testResult, commitInfo)
	if fail_exists {
		send(emailContent, change, testResult, committer, sender, senderPasswd)
		fmt.Println("Email sent to committer:", committer)
	} else {
		fmt.Println("ALL TESTS PASSED!")
	}
}

func compose(testResult string, commitInfo string) (string, string, bool){
	var committer string
	sendEmail := false

	//read log file
	testMSG_byte, err := ioutil.ReadFile(testResult)
	if err != nil {
		fmt.Printf("Failed to read file \"%s\"", testResult)
		return "", "", sendEmail
	}

	//read commitInfo file
	commitMSG_byte, err := ioutil.ReadFile(commitInfo)
	if err != nil {
		fmt.Printf("Failed to read file \"%s\"", commitInfo)
		return "", "", sendEmail
	}

	//convert to string
	testMSG   := string(testMSG_byte)
	commitMSG := string(commitMSG_byte)

	emailContents_commit   := "<p>Committer Information:"
	emailContents_testInfo := "<p>Failing Tests Information:"

	//Compose the commit information section of the email
	commitMsgScanner := bufio.NewScanner(strings.NewReader(commitMSG))
	for i := 0; commitMsgScanner.Scan(); i++ {
		MSG := commitMsgScanner.Text()
		if i == 6 {
			MSG = "<br> Commit Summary: " + MSG
		} else if MSG == "" {
			continue
		} else {
			if strings.Contains(MSG, "<") {
				if strings.Contains(MSG, "Commit:") {
					committer = between(MSG, "<", ">")
				}
				MSG = strings.Replace(MSG, "<", "", -1)
				MSG = strings.Replace(MSG, ">", "", -1)
			}
			MSG = "<br>" + MSG
		}
		emailContents_commit += MSG
	}
	emailContents_commit += "</p>"

	//Compose the test result information section of the email
	testMsgScanner := bufio.NewScanner(strings.NewReader(testMSG))
	for testMsgScanner.Scan() {
		MSG := testMsgScanner.Text()
		if (strings.Contains(MSG, "FAIL")) {
			if (strings.TrimLeft(MSG, "FAIL") != "") {
				sendEmail = true
				MSG = "<br>" + MSG
				emailContents_testInfo += MSG
			}
		}
	}
	emailContents_testInfo += "</p>"

	//Merge both sections together
	emailContents := emailContents_commit + emailContents_testInfo

	return committer, emailContents, sendEmail
}

func between(value string, a string, b string) string {
    // Get substring between two strings.
    posFirst := strings.Index(value, a)
    if posFirst == -1 {
        return ""
    }
    posLast := strings.Index(value, b)
    if posLast == -1 {
        return ""
    }
    posFirstAdjusted := posFirst + len(a)
    if posFirstAdjusted >= posLast {
        return ""
    }
    return value[posFirstAdjusted:posLast]
}

func send(emailBody string, change string, testResult string, committer string, sender string, senderPasswd string) {
	//send the email
	mail := gomail.NewMessage()
	mail.SetHeader("From", sender)
	mail.SetHeader("To",   committer)
	//mail.SetAddressHeader("Cc", "dan@example.com", "Dan")
	mail.SetHeader("Subject", "Go-Dappley Commit Test Result")
	mail.SetBody("text/html", emailBody)
	mail.Attach(change)
	mail.Attach(testResult)

	deliver := gomail.NewDialer("smtp.gmail.com", 587, sender, senderPasswd)

	if err := deliver.DialAndSend(mail); err != nil {
		panic(err)
	}
}