# go-dappley-commit-report
Sends email to the committer when there the committed update has error.

### Pipeline:
```
pipeline {
    agent any
    tools {
        go 'go-1.16.3'
    }
    environment {
        GO1163MODULE = 'on'
    }
    stages {
        stage('SCM Checkout') {
            steps {
                git 'https://github.com/heesooh/go-dappley-commit-report'
            }
        }
        stage('Compile') {
            steps {
                sh 'go mod init github.com/heesooh/go-dappley-commit-report'
                sh 'go mod tidy'
                sh 'go build EmailOnPush.go'
            }
        }
        stage('Deploy') {
            steps {
                sh "./EmailOnPush -senderEmail <Email address of the sender> -senderPasswd <Email password>"
            }
        }
        stage('Close') {
            steps {
                sh 'rm -r *'
            }
        }
    }
}
```