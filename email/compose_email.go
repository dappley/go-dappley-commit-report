package email

import (
	"github.com/heesooh/go-dappley-commit-report/helper"
	"io/ioutil"
	"strings"
	"bufio"
	"fmt"
)

//Compose the email message content for the go-dappley-commit-report.
func ComposeEmail(branch string) (committer string, emailContents string, sendEmail bool){
	//read log.txt file.
	testMSG_byte, err := ioutil.ReadFile(branch + "/log.txt")
	if err != nil {
		fmt.Printf("Failed to read from origin/%s branch", branch)
		return
	}
	//read commitInfo.txt file.
	commitMSG_byte, err := ioutil.ReadFile(branch + "/commitInfo.txt")
	if err != nil {
		fmt.Printf("Failed to read from origin/%s branch", branch)
		return
	}

	//convert to string.
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
					committer = helper.Between(MSG, "<", ">")
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
	emailContents = branch_info + emailContents_commit + emailContents_testInfo

	return
}