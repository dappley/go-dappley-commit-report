package main

import (
	"io/ioutil"
	"strings"
	"bufio"
	"fmt"
)


func composeEmail(branch string) (string, string, bool){
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