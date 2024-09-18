pipeline {
    agent any

    environment {
        DB_USER = credentials('db-user')
        DB_PASS = credentials('db-pass')
        DB_NAME = credentials('db-name')
        DB_HOST = credentials('db-host')
        DB_PORT = credentials('db-port')
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
                    sh 'docker compose up -d'
                }
            }
        }
    }
}