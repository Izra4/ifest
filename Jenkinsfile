pipeline{
    agent any
    stages{
        stage('Clone Repository'){
            steps{
                git 'https://github.com/Izra4/ifest.git'
            }
        }
        stage('Build image'){
            steps{
                script{
                    sh 'docker compose --build'
                }
            }

        }
        stage('Deploy'){
            steps{
                script{
                    sh 'docker compose up -d'
                }
            }

        }
    }
}