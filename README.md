# go-dappley-commit-report
Sends email to the committer when there the committed update has error.

Jenkins Pipeline Setup:
```
pipeline {
    agent any
    tools {
        go 'go-1.15.7'
    }
    parameters {
        gitParameter branchFilter: 'origin/(.*)', defaultValue: 'develop', name: 'BRANCH', type: 'PT_BRANCH'
    }
    environment {
        GO1157MODULE = 'on'
    }
    stages {
        stage('SCM Checkout') {
            steps {
                git branch: "${params.BRANCH}", url: 'https://github.com/heesooh/go-dappley.git'
            }
        }
        stage('Compile') {
            steps {
                sh 'make build'
                sh 'go build EmailOnPush.go'
            }
        }
        stage('Test') {
            steps {
                sh 'git show > change.txt'
                sh 'make testall > log.txt'
                sh 'git log --pretty=fuller HEAD^..HEAD > commitInfo.txt'
            }
        }
        stage('Deploy') {
            steps {
                sh "./EmailOnPush -change change.txt -testResult log.txt -commitInfo commitInfo.txt -sender username@example.com -senderPasswd password"
            }
        }
        stage('Close') {
            steps {
                sh 'rm log.txt'
                sh 'rm change.txt'
                sh 'rm commitInfo.txt'
            }
        }
    }
}
```
