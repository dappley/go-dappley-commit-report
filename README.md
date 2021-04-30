# go-dappley-commit-report
Sends email to the committer when there the committed update has error.

### Jenkins Pipeline 1:
```
pipeline {
    agent any
    tools {
        go 'go-1.16.3'
    }
    environment {
        GO1163MODULE = 'on'
    }
    parameters {
        gitParameter branchFilter: 'origin/(.*)', defaultValue: 'master',  name: 'MASTER',  type: 'PT_BRANCH'
        gitParameter branchFilter: 'origin/(.*)', defaultValue: 'develop', name: 'DEVELOP', type: 'PT_BRANCH'
    }
    stages {
        stage('SCM Checkout Master Branch') {
            steps {
                git branch: "${params.MASTER}", url: 'https://github.com/dappley/go-dappley.git'
            }
        }
        // stage('Make Master Branch') {
        //     steps {
        //         sh 'make build'
        //     }
        // }
        stage('Test Master Branch') {
            steps {
                sh 'mkdir master'
                sh 'git show > master/change.txt'
                sh 'make testall > master/log.txt'
                sh 'git log --pretty=fuller HEAD^..HEAD > master/commitInfo.txt'
            }
        }
        stage('Move Master Files') {
            steps {
                sh 'mv master ../go-dappley-commit-report'
            }
        }
        stage('Clear Master Directory') {
            steps {
                sh 'rm -r *'
            }
        }
        stage('SCM Checkout Develop Branch') {
            steps {
                git branch: "${params.DEVELOP}", url: 'https://github.com/dappley/go-dappley.git'
            }
        }
        // stage('Compile Develop Branch') {
        //     steps {
        //         sh 'make build'
        //     }
        // }
        stage('Test Develop Branch') {
            steps {
                sh 'mkdir develop'
                sh 'git show > develop/change.txt'
                sh 'make testall > develop/log.txt'
                sh 'git log --pretty=fuller HEAD^..HEAD > develop/commitInfo.txt'
            }
        }
        stage('Move Develop Files') {
            steps {
                sh 'mv develop ../go-dappley-commit-report'
            }
        }
        stage('Clear Develop Directory') {
            steps {
                sh 'rm -r *'
            }
        }
    }
}
```


### Jenkins Pipeline 2:
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
                sh "./EmailOnPush -senderEmail <EMAIL ADDRESS> -senderPasswd <PASSWORD>"
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