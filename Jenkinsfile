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
                        "JWT_SECRET=${JWT_SECRET}"
                    ]) {
                        sh 'docker compose up -d'
                    }
                }
            }
        }
    }
}