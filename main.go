package main

import (
	"flag"
	"fmt"
)

func main() {
	var senderEmail, senderPasswd string
	flag.StringVar(&senderEmail, "senderEmail", "sender_username@example.com", "Email address of the sender")
	flag.StringVar(&senderPasswd, "senderPasswd", "PASSWORD", "Password of the sender's email address.")
	flag.Parse()

	dev_committer, dev_email, dev_fail := composeEmail("develop")
	master_committer, master_email, master_fail := composeEmail("master")

	fmt.Println("Branch: develop\nCommitter:", dev_committer,    "\n", dev_email)
	fmt.Println("Branch: master \nCommitter:", master_committer, "\n", master_email)

	//send develop branch info
	if dev_fail {
		sendEmail(dev_email, "develop", dev_committer, senderEmail, senderPasswd)
		fmt.Println("Email sent to develop branch committer:", dev_committer)
	} else {
		fmt.Println("No fail case on develop branch!")
	}

	//send master branch info
	if master_fail {
		sendEmail(master_email, "master", master_committer, senderEmail, senderPasswd)
		fmt.Println("Email sent to master branch committer:", master_committer)
	} else {
		fmt.Println("No fail case on master branch!")
	}
}