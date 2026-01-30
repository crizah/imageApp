pipeline {
    agent any

    environment {
        DOCKERHUB_CREDENTIALS = 'docker-hub-creds-global'
        USERNAME = "shaizah"
        SERVER_IMAGE = 'shaizah/kube:imageApp-server'
        WEB_IMAGE    = 'shaizah/kube:imageApp-web'
        SERVER_PATH = './server'
        WEB_PATH = './web'
        
    }

    stages {

        stage('checkout') {
            steps {
                checkout scm
            }
        }

        stage('detect change') {
            steps {
                script {
                    // detect if changes made in server or web
                    def changedFiles = sh(
                        script: "git diff --name-only HEAD~1 HEAD",
                        returnStdout: true
                    ).trim()

                    env.SERVER_CHANGED = changedFiles.contains('server/') ? 'true' : 'false'
                    env.WEB_CHANGED    = changedFiles.contains('web/')    ? 'true' : 'false'
                }
            }
        }

        stage('build and push server images') {
            when {
                expression { env.SERVER_CHANGED == 'true' }
            }
            steps {
                script {

                    docker.withRegistry('', DOCKERHUB_CREDENTIALS) {
                        sh """
                        docker build -t ${SERVER_IMAGE} ${SERVER_PATH}
                        docker push ${SERVER_IMAGE}
                        """
                    }
                }


                
            }
        }

        stage('build and push web images') {
            when {
                expression { env.WEB_CHANGED == 'true' }
            }
            steps {
                script {
                   
                    docker.withRegistry('', DOCKERHUB_CREDENTIALS) {
                        sh """
                        docker build -t ${WEB_IMAGE} ${WEB_PATH}
                        docker push ${WEB_IMAGE}
                        """
                    }
                }
            }
        }
    }

    post {
        success {
            echo "yay yippee"
        }
        failure {
            echo ":("
        }
    }
}
