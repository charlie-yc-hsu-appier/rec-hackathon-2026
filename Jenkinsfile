#!/usr/bin/env groovy

def pod_label = "worker-${JOB_NAME}".replace("@", "-").replace("/", "-").take(50) + "_job"

def genDeployMsg(status) {
  notifyAuthor = status == "success" ? "${env.COMMIT_AUTHOR}" : "@${env.COMMIT_AUTHOR}"
  attachments = [
    [
      title: "Deploy ${status} (${env.JOB_NAME})",
      fields: [
        [
            "title": "Project",
            "value" : "${env.JOB_NAME}",
            "short": true
        ],
        [
            "title": "Author",
            "value" : notifyAuthor,
            "short": true
        ],
        [
            "title": "Version",
            "value" : "${env.BUILD_VERSION}",
            "short": true
        ],
        [
            "title": "Build URL",
            "value" : "${env.BUILD_URL}",
            "short" : true
        ],
        [
            "title": "Build Duration",
            "value": "${currentBuild.durationString}",
            "short": true
        ]
      ]
    ]
  ]
  color = status == "success" ? "good" : "danger"
  attachments[0]['color'] = color
  return attachments
}

pipeline {
  agent {
    kubernetes {
      label pod_label
      defaultContainer 'deploy-kit'
      yaml """
        apiVersion: v1
        kind: Pod
        spec:
          containers:
          - name: dind
            image: docker:20.10-dind
            env:
              - name: DOCKER_TLS_CERTDIR
                value: ""
            securityContext:
              privileged: true
            resources:
              limits:
                memory: 2Gi
              requests:
                memory: 2Gi
                cpu: 2

          - name: golang
            image: golang:1.25-alpine
            tty: true
            resources:
              limits:
                memory: 2Gi
              requests:
                memory: 2Gi
                cpu: 2

          - name: gcr
            image: asia-docker.pkg.dev/appier-docker/docker-ai-rec-asia/gcr-alpine:1.3.9
            tty: true
            env:
              - name: DOCKER_HOST
                value: tcp://127.0.0.1:2375

          - name: deploy-kit
            image: asia-docker.pkg.dev/appier-docker/docker-ai-rec-asia/argo-deploy-kits:1.3.0
            tty: true

          - name: robot
            image: asia-docker.pkg.dev/appier-docker/docker-ai-rec-asia/qa/system_test_robot:v1.0.28
            tty: true
            command:
             - cat
            resources:
              limits:
                memory: 2Gi
                cpu: 2
              requests:
                memory: 2Gi
                cpu: 2

          volumes:
            - name: build-volume
              emptyDir: {}
      """
    }
  }

  options {
    disableConcurrentBuilds()
    buildDiscarder(logRotator(daysToKeepStr: '60'))
  }

  environment {
    // vault
    VAULT_ADDR = 'https://vault.appier.us/'
    VAULT_APPROLE = 'ai-jenkins-gcp'
    VAULT_APPROLE_SECRET = credentials('VAULT_SECRET_ID')
    VAULT_KEY_PATH = 'secret/project/recommendation'

    // SA credential
    CICD_GCLOUD_JSON_FILE = 'rec-jenkins-cicd.json'

    // SSH key
    REC_COMMON_LIB_KEY = 'rec_common_lib_key'

    // helm
    CHART_DIR = 'deploy/rec-vendor-api'

    // build version setting
    VERSION_MAJOR = 1
    VERSION_MINOR = 0
    VERSION_PATCH = "${BUILD_NUMBER}"

    // for QA system tests
    WORK_DIR = './tests/system_tests/API'
    TESTRAIL_URL = 'https://appier.testrail.io/index.php?/runs/view'
  }

  stages {
    stage('Prepare: prepare environment variable'){
      steps {
        script {
          sh 'apk update && apk add git'
          sh 'git config --global --add safe.directory $(pwd)'
          env.COMMIT_EMAIL = sh(script: 'git --no-pager show -s --pretty=%ae', returnStdout: true).replace('\n', '')
          env.COMMIT_AUTHOR = env.COMMIT_EMAIL.split('@')[0]

          env.BUILD_ENV = 'dev'
          env.BUILD_VERSION = "${env.BRANCH_NAME}-${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}"
          env.GKE_CLUSTER = 'nelson'

          if (env.BRANCH_NAME == 'staging') {
            env.BUILD_ENV = 'stg'
          }
          if (env.BRANCH_NAME == 'master') {
            env.BUILD_ENV = 'prd'
            env.BUILD_VERSION = "${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}"
            env.GKE_CLUSTER = 'echinata'
          }
        }
      }
    }

    stage("Credentials: get credentials from vault") {
      steps {
        container('deploy-kit') {
          sh '''
            vault-login.sh

            # APP credential
            vault-decrypt-v2.sh config-template/vendors.yaml $CHART_DIR/secrets/vendors.yaml
            vault-decrypt-v2.sh config-template/config-$BUILD_ENV.yaml $CHART_DIR/secrets/config.yaml
            cp ./config-template/nginx.conf $CHART_DIR/secrets/

            # SA credential
            vault kv get --field=data secret/project/_gcp/iam/appier-ai-recommendation/rec-jenkins-cicd > $CICD_GCLOUD_JSON_FILE

            vault kv get --field=private_key secret/project/recommendation/ssh_key/ai-rec-common > $REC_COMMON_LIB_KEY
            chmod 600 $REC_COMMON_LIB_KEY
          '''
        }
      }
    }

    stage('Pre-commit check') {
      steps {
        container('golang') {
          sh 'apk add --no-cache openssh make git build-base openssl'
          sh 'wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v2.8.0'
          sh 'go install golang.org/x/tools/cmd/goimports@v0.41.0'

          sh 'mkdir -p -m 0600 /root/.ssh'
          sh 'touch /root/.ssh/known_hosts'
          sh 'ssh-keyscan github.com >> ~/.ssh/known_hosts'
          sh 'cp $REC_COMMON_LIB_KEY /root/.ssh/ai_rec_common'
          sh 'cp .gitconfig /root/.gitconfig'
          sh 'GIT_SSH_COMMAND="ssh -i /root/.ssh/ai_rec_common" go mod download'

          sh 'make fmt-check'          
          // adding `buildvcs=false` to mitigate the "error obtaining VCS status: exit status 128" error
          // note: `buildvcs=false` flag is not a must when running golangci-lint
          sh 'GOFLAGS=-buildvcs=false ./bin/golangci-lint run'

          // Validate vendor configuration with real secrets from vault
          sh 'go run ./cmd/validate-config $CHART_DIR/secrets/vendors.yaml'
          sh 'make test'
        }
      }
    }

    stage('Publish image') {
      when {
        anyOf {
          branch 'master';
          branch 'staging';
        }
      }
      environment {
        GCR_JSON_FILE = "$CICD_GCLOUD_JSON_FILE"
        GCR_REGISTRY  = 'asia-docker.pkg.dev/appier-docker/docker-ai-rec-asia'
        GCR_REPO      = 'rec-vendor-api'

        DOCKER_CACHE      = 'asia-docker.pkg.dev/appier-docker/docker-ai-rec-asia/rec-vendor-api:cache'
        DOCKER_CONTEXT    = './Dockerfile'
        DOCKER_TAG        = "$BUILD_VERSION"
        DOCKER_EXTRA_ARGV = "--ssh ai-rec-common=$REC_COMMON_LIB_KEY"
        DOCKER_BUILDKIT   = 1

        WORKDIR = './'
      }
      steps {
        container('gcr') {
          sh 'gcloud auth configure-docker --quiet asia-docker.pkg.dev'
          sh 'docker-publish'
        }
      }
    }

    stage('Deploy rec-vendor-api') {
      when {
        anyOf {
          branch 'master';
          branch 'staging';
        }
      }
      environment {
        NAMESPACE = 'rec'
        CHART = "$CHART_DIR"
        RELEASE= "rec-vendor-api-$BUILD_ENV"

        GCP_PROJECT = 'appier-k8s-ai-rec'
        GKE_REGION = 'asia-east1'
        GCLOUD_JSON_FILE = "$CICD_GCLOUD_JSON_FILE"

        WAIT = 'true'
        FORCE = 'false'
        HISTORY_MAX = 5

        VALUE_FILES="${CHART_DIR}/values-${BUILD_ENV}.yaml"
        VALUES = "image.tag=$BUILD_VERSION"
      }
      steps {
        retry(3) {
          container('deploy-kit') {
            sh 'helm-deploy.sh'
          }
        }
      }
    }

    stage('QA prepare environment') {
      when {
        branch 'staging'
      }
      steps {
        container('deploy-kit') {
          sh '''
            vault-login.sh

            # APP credential
            vault-decrypt-v2.sh config-template/vendors.yaml $CHART_DIR/secrets/vendors.yaml
            vault-decrypt-v2.sh config-template/config-$BUILD_ENV.yaml $CHART_DIR/secrets/config.yaml

            # SA credential
            vault kv get --field=data secret/project/_gcp/iam/appier-ai-recommendation/rec-jenkins-cicd > $CICD_GCLOUD_JSON_FILE

            vault kv get --field=private_key secret/project/recommendation/ssh_key/ai-rec-common > $REC_COMMON_LIB_KEY
            chmod 600 $REC_COMMON_LIB_KEY
          '''
        }
      }
    }

    stage('QA system tests') {
      when {
        branch 'staging'
      }
      options {
        timeout(time: 10, unit: 'MINUTES')   // timeout on this stage
      }
      steps {
            container('robot'){
              withCredentials([[$class: 'UsernamePasswordMultiBinding',
              credentialsId: 'testrail',
              usernameVariable: 'TESTRAIL_USER',
              passwordVariable: 'TESTRAIL_API_KEY']]){

                  // Check the path
                  echo "Checking the permission ..."
                  sh "pwd"
                  sh "ls -al ${WORK_DIR}"
                  sh "ls -al ${WORK_DIR}/run_robot.sh"
                  sh "ls -al ${WORK_DIR}/res"

                  // The content of pass_rate.txt will be given by the run_robot.sh below after the automation run is completed
                  sh "touch ${WORK_DIR}/testsuite/pass_rate.txt"
                  sh "chmod 666 ${WORK_DIR}/testsuite/pass_rate.txt"

                  // run the robot script
                  script{
                    env.RUN_ID = sh(script: '${WORK_DIR}/run_robot.sh -u $TESTRAIL_USER  -k $TESTRAIL_API_KEY  -c rec_vendor -t rat -d no | grep  \'Task ID is :\' | cut -c\'22-\'',returnStdout: true).trim()
                    if (env.RUN_ID){
                        sh "echo 'rec_vendor run ID:' ${RUN_ID}"
                        env.EXECUTION_RATE = sh(script: "cat ${WORK_DIR}/testsuite/pass_rate.txt",returnStdout: true).trim()
                        sh "echo ${EXECUTION_RATE}"

                        // merge the report
                        sh "rebot -d ${WORK_DIR}/report/ --RemoveKeywords passed -o output.xml  ${WORK_DIR}/report/rec_vendor_rat.xml  &> /dev/null"
                        sh "sleep 60"
                        sh "ls -al ${WORK_DIR}/report/"

                        // Since multiple result will be merged into the same case_ID at testrail, we try to use the number of fail from the final XML report
                        env.FAIL_COUNT = sh(script: "cat ${WORK_DIR}/report/rec_vendor_rat.xml | grep -oP \"<stat (.+?)>All Tests</stat>\" | grep -P -o -e '(?<=fail=\").*?(?=\")'",returnStdout: true).trim()
                        sh "echo 'Failed Count:'  ${FAIL_COUNT}"
                    }
                    else {
                        // stop the automation process if we can not get the run_id from testrail API
                        error("[FAIL][QA-TESTRAIL] Can't get the test run_id. Stop executing the E2E automation...")
                    }
                    if (env.FAIL_COUNT != '0') {
                        error '!!! There are some failed cases in the E2E testing !!!'
                    }
                  }
              }
            }
      }
      post {
        always {
          script {
            if (env.BRANCH_NAME == 'staging' && env.RUN_ID ) {
              step(
                [
                  $class              : 'RobotPublisher',
                  outputPath          : "${WORK_DIR}/report",
                  outputFileName      : '**/output.xml',
                  reportFileName      : '**/report.html',
                  logFileName         : '**/log.html',
                  disableArchiveOutput: false,
                  passThreshold       : 100,
                  unstableThreshold   : 100,
                  otherFiles          : "**/*.png,**/*.jpg",
                ]
              )
            }
          }
        }
      }
    }

    stage('Remove unused docker images') {
      when {
        branch 'master'
      }

      environment {
        GCR_JSON_FILE = "$CICD_GCLOUD_JSON_FILE"
        GCR_REGISTRY = 'asia-docker.pkg.dev/appier-docker/docker-ai-rec-asia'
        KEEP_COUNT = 50
      }

      parallel {
        stage('rec-vendor-api') {
          environment {
            GCR_REPO = 'rec-vendor-api'
          }

          steps {
            container('gcr') {
              sh 'remove-unused-images'
            }
          }
        }

        stage('rec-vendor-api-dev') {
          environment {
            GCR_REPO = 'rec-vendor-api-dev'
          }

          steps {
            container('gcr') {
              sh 'remove-unused-images'
            }
          }
        }
      }
    }
  }  // end of stages

  post {
    success {
      script {
        if (env.BRANCH_NAME ==~ /(master|staging)/) {
            slackSend (channel: "#ec-rec-deployment-alerts",
                       attachments: genDeployMsg("success"))
        }
      }
    }

    failure {
      script {
        if (env.BRANCH_NAME ==~ /(master|staging)/) {
            slackSend (channel: "#ec-rec-deployment-alerts",
                       attachments: genDeployMsg("failure"))
        }
      }
    }
  }
}  // end of pipeline
