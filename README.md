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
                sh 'git show > change_master.txt'
                sh 'make testall > log_master.txt'
                sh 'git log --pretty=fuller HEAD^..HEAD > commitInfo_master.txt'
            }
        }
        stage('Move Master Files') {
            steps {
                sh 'mv change_master.txt ../go-dappley-commit-report'
                sh 'mv log_master.txt ../go-dappley-commit-report'
                sh 'mv commitInfo_master.txt ../go-dappley-commit-report'
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
                sh 'git show > change_develop.txt'
                sh 'make testall > log_develop.txt'
                sh 'git log --pretty=fuller HEAD^..HEAD > commitInfo_develop.txt'
            }
        }
        stage('Move Develop Files') {
            steps {
                sh 'mv change_develop.txt ../go-dappley-commit-report'
                sh 'mv log_develop.txt ../go-dappley-commit-report'
                sh 'mv commitInfo_develop.txt ../go-dappley-commit-report'
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
        GO1157MODULE = 'on'
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
        // stage('Deploy') {
        //     steps {
        //         sh "./EmailOnPush -change change.txt -testResult log.txt -commitInfo commitInfo.txt -sender username@example.com -senderPasswd password"
        //     }
        // }
        stage('Close') {
            steps {
                // sh 'rm EmailOnPush'
                // sh 'rm log_master.txt'
                // sh 'rm log_develop.txt'
                // sh 'rm change_master.txt'
                // sh 'rm change_develop.txt'
                // sh 'rm commitInfo_master.txt'
                // sh 'rm commitInfo_develop.txt'
                sh 'rm -r *'
            }
        }
    }
}
```