package main

import (
	"github.com/heesooh/go-dappley-commit-report/email"
	"github.com/heesooh/go-dappley-commit-report/helper"
	"flag"
	"fmt"
	"log"
)

func main() {
	var senderEmail, senderPasswd string
	flag.StringVar(&senderEmail, "senderEmail", "default_email", "Email address of the sender")
	flag.StringVar(&senderPasswd, "senderPasswd", "default_password", "Password of the sender's email address.")
	flag.Parse()

	err := helper.CheckFlags(senderEmail, senderPasswd)
	if err != nil {
		log.Fatal(err)
		return
	}

	dev_committer, dev_email, dev_fail := email.ComposeEmail("develop")
	master_committer, master_email, master_fail := email.ComposeEmail("master")

	fmt.Println("Branch: develop\nCommitter:", dev_committer,    "\n", dev_email)
	fmt.Println("Branch: master \nCommitter:", master_committer, "\n", master_email)

	//send develop branch info
	if dev_fail {
		email.SendEmail(dev_email, "develop", dev_committer, senderEmail, senderPasswd)
		fmt.Println("Email sent to develop branch committer:", dev_committer)
	} else {
		fmt.Println("No fail case on develop branch!")
	}

	//send master branch info
	if master_fail {
		email.SendEmail(master_email, "master", master_committer, senderEmail, senderPasswd)
		fmt.Println("Email sent to master branch committer:", master_committer)
	} else {
		fmt.Println("No fail case on master branch!")
	}
}