pipeline {
    agent any

    environment {
        DB_USER = credentials('db-user')
        DB_PASS = credentials('db-pass')
        DB_NAME = credentials('db-name')
        DB_HOST = credentials('db-host')
        DB_PORT = credentials('db-port')
        JWT_EXP = credentials('jwt-exp')
        JWT_SECRET = credentials('jwt-secret')
        CLIENT_ID = credentials("client-id")
        CLIENT_SECRET = credentials("client-secret")
        REDIRECT_URL = credentials("redirect-url")
        AES_KEY = credentials('aes-key')
        SUPA_URL = credentials('supa-url')
        SUPA_API = credentials('supa-api')
    }

    stages {
        stage('Clone Repository') {
            steps {
                git 'https://github.com/Izra4/ifest.git'
            }
        }

        stage('Build Image') {
            steps {
                script {
                    sh 'docker compose build'
                }
            }
        }

        stage('Deploy') {
            steps {
                script {
                    withEnv([
                        "DB_USER=${DB_USER}",
                        "DB_PASS=${DB_PASS}",
                        "DB_NAME=${DB_NAME}",
                        "DB_HOST=${DB_HOST}",
                        "DB_PORT=${DB_PORT}",
                        "JWT_EXP=${JWT_EXP}",
                        "JWT_SECRET=${JWT_SECRET}",
                        "CLIENT_ID=${CLIENT_ID}",
                        "CLIENT_SECRET=${CLIENT_SECRET}",
                        "REDIRECT_URL=${REDIRECT_URL}",
                        "AES_KEY=${AES_KEY}",
                        "SUPA_API=${SUPA_API}",
                        "SUPA_URL=${SUPA_URL}"
                    ]) {
                        sh 'docker compose up -d'
                    }
                }
            }
        }
    }
}