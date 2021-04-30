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
	var senderEmail, senderPasswd string
	flag.StringVar(&senderEmail, "senderEmail", "sender_username@example.com", "Email address of the sender")
	flag.StringVar(&senderPasswd, "senderPasswd", "PASSWORD", "Password of the sender's email address.")
	flag.Parse()

	dev_committer, dev_email, dev_fail := compose("develop")
	master_committer, master_email, master_fail := compose("master")

	fmt.Println("Branch: develop")
	fmt.Println("Committer:", dev_committer)
	fmt.Println()
	fmt.Println(dev_email)
	
	fmt.Println("Branch: master")
	fmt.Println("Committer:", master_committer)
	fmt.Println()
	fmt.Println(master_email)
	
	//send develop branch info
	if dev_fail {
		send(dev_email, "develop", dev_committer, senderEmail, senderPasswd)
		fmt.Println("Email sent to develop branch committer:", dev_committer)
	} else {
		fmt.Println("No fail case on develop branch!")
	}

	//send master branch info
	if master_fail {
		send(master_email, "master", master_committer, senderEmail, senderPasswd)
		fmt.Println("Email sent to master branch committer:", master_committer)
	} else {
		fmt.Println("No fail case on master branch!")
	}
}

func compose(branch string) (string, string, bool){
	var committer string
	sendEmail := false

	//read log file
	testMSG_byte, err := ioutil.ReadFile(branch + "/log.txt")
	if err != nil {
		fmt.Printf("Failed to read from origin/%s branch", branch)
		return "", "", sendEmail
	}

	//read commitInfo file
	commitMSG_byte, err := ioutil.ReadFile(branch + "/commitInfo.txt")
	if err != nil {
		fmt.Printf("Failed to read from origin/%s branch", branch)
		return "", "", sendEmail
	}

	//convert to string
	testMSG   := string(testMSG_byte)
	commitMSG := string(commitMSG_byte)

	emailContents_commit   := "<p>Committer Information:"
	emailContents_testInfo := "<p>Failing Tests Information:"

	//Compose the commit information section of the email
	commitMsgScanner := bufio.NewScanner(strings.NewReader(commitMSG))
	for i := 0; commitMsgScanner.Scan() && i < 7; i++ {
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

	branch_info := "<p>Origin/" + branch + "::</p>"

	//Merge both sections together
	emailContents := branch_info + emailContents_commit + emailContents_testInfo

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

func send(emailBody string, branch string, committer string, senderEmail string, senderPasswd string) {
	//send the email
	mail := gomail.NewMessage()
	mail.SetHeader("From", senderEmail)
	fmt.Println(committer)
	mail.SetHeader("To",   "blockchainwarning@omnisolu.com")
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