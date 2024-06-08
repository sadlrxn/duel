pipeline {
    agent any
    environment {
        JENKINS_HOME = "/var/lib/jenkins"
        AWS_ACCOUNT_ID ="714239636777"
        AWS_DEFAULT_REGION ="us-east-2" 
        IMAGE_REPO_NAME ="duelana_v1"
        REPOSITORY_URI = "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com/${IMAGE_REPO_NAME}"

        EC2_SSH_PRIV_URL = "${JENKINS_HOME}/.ssh/id_duelana_ec2_ssh"
        EC2_SSH_USER = "ubuntu"
        EC2_SSH_HOST = "3.143.88.5"
        EC2_SSH_COMMAND = "cd /home/ubuntu/workspace/duelana-v1-docker && sudo ./run.sh"

        EC2_SSH_PRIV_URL_TEST = "${JENKINS_HOME}/.ssh/id_duelana_ec2_test_ssh"
        EC2_SSH_USER_TEST = "ubuntu"
        EC2_SSH_HOST_TEST = "3.132.193.96"
        EC2_SSH_COMMAND_TEST = "cd /home/ubuntu/workspace/duelana-v1-test-docker && sudo ./run.sh"
    }
   
    stages {
        stage('Logging into AWS ECR') {
            steps {
                script {
                    sh "aws ecr get-login-password --region ${AWS_DEFAULT_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com"
                }
            }
        }
        
        stage('Cloning Git') {
            steps {
                checkout([$class: 'GitSCM', branches: [[name: '*/' + env.BRANCH_NAME]], extensions: [], userRemoteConfigs: [[credentialsId: 'GITHUB_CREDENTIAL', url: 'https://github.com/Duelana-Team/duelana-v1.git']]])
            }
        }
  
        // Building Docker images
        stage('Building image') {
            steps{
                script {
                    def IMAGE_TAG = env.BRANCH_NAME.replaceAll('/', '_')
                    def BUILD_ARGS = [
                        "prod" : "--build-arg GENERATE_SOURCEMAP=false --build-arg MASTER_WALLET_PUBLIC_KEY=DUELLrBB96snTu3Wn3Cjyj7s2pRRqiG5LpPCC1fmw2Wm --build-arg NETWORK=mainnet --build-arg SOLANA_ENDPOINT=https://duelana-mainbca-2158.mainnet.rpcpool.com/ --build-arg STAGE=beta --build-arg HAPPY_HOLIDAY=false",
                        "beta" : "--build-arg GENERATE_SOURCEMAP=false --build-arg MASTER_WALLET_PUBLIC_KEY=DUELjQHtLZDkf9rhrw96Sko7yxzxnWABbvV9F2AM7k3C --build-arg NETWORK=mainnet --build-arg SOLANA_ENDPOINT=https://duelana-mainbca-2158.mainnet.rpcpool.com/ --build-arg STAGE=beta --build-arg HAPPY_HOLIDAY=false",
                        "develop" : "--build-arg GENERATE_SOURCEMAP=false --build-arg MASTER_WALLET_PUBLIC_KEY=EEMxfcPwMK615YLbEhq8NVacdmxjkxkok6KXBJBHuZfB --build-arg SOLANA_ENDPOINT=https://duelana-dev2e70-9c52.devnet.rpcpool.com/ --build-arg STAGE=dev --build-arg HAPPY_HOLIDAY=false"
                    ]
                    def DEFAULT_BUILD_ARGS = "--build-arg GENERATE_SOURCEMAP=false --build-arg MASTER_WALLET_PUBLIC_KEY=EEMxfcPwMK615YLbEhq8NVacdmxjkxkok6KXBJBHuZfB"
                    sh """docker build -t ${REPOSITORY_URI}:${IMAGE_TAG} -f ./Dockerfile ${BUILD_ARGS[IMAGE_TAG] ?: DEFAULT_BUILD_ARGS} ."""
                }
            }
        }
   
        // Uploading Docker images into AWS ECR
        stage('Pushing to ECR') {
            steps {
                script {
                    def IMAGE_TAG = env.BRANCH_NAME.replaceAll('/', '_')
                    sh "docker push ${REPOSITORY_URI}:${IMAGE_TAG}"
                }
            }
        }

        // Publish the image to test link for beta branch
        // stage('Publishing beta branch to live website') {
        //     when {
        //         branch 'beta'
        //     }
        //     steps {
        //         script {
        //             sh "eval \$(ssh-agent) && ssh-add ${EC2_SSH_PRIV_URL} && ssh ${EC2_SSH_USER}@${EC2_SSH_HOST} '${EC2_SSH_COMMAND}'"
        //         }
        //     }
        // }

        // Publish the image to staging link for develop branch
        stage('Publishing develop branch to staging website') {
            when {
                branch 'develop'
            }
            steps {
                script {
                    sh "eval \$(ssh-agent) && ssh-add ${EC2_SSH_PRIV_URL_TEST} && ssh ${EC2_SSH_USER_TEST}@${EC2_SSH_HOST_TEST} '${EC2_SSH_COMMAND_TEST}'"
                }
            }
        }

        // Clearing docker caches while build
        // stage('Clearing docker build caches') {
        //     steps {
        //         script {
        //             sh "docker system prune -f"
        //         }
        //     }
        // }
    }
}