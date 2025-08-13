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
            image: golang:1.23-alpine
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
            vault-decrypt-v2.sh config-template/config-$BUILD_ENV.yaml $CHART_DIR/secrets/config.yaml

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
          sh 'apk add --no-cache openssh curl make git build-base'
          sh 'wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.61.0'

          sh 'mkdir -p -m 0600 /root/.ssh'
          sh 'touch /root/.ssh/known_hosts'
          sh 'ssh-keygen -f "/root/.ssh/known_hosts" -R "bitbucket.org"'
          sh 'curl https://bitbucket.org/site/ssh >> /root/.ssh/known_hosts'
          sh 'cp $REC_COMMON_LIB_KEY /root/.ssh/ai_rec_common'
          sh 'cp .gitconfig /root/.gitconfig'
          sh 'GIT_SSH_COMMAND="ssh -i /root/.ssh/ai_rec_common" go mod download'

          sh 'make test'
          // adding `buildvcs=false` to mitigate the "error obtaining VCS status: exit status 128" error
          // note: `buildvcs=false` flag is not a must when running golangci-lint
          sh 'GOFLAGS=-buildvcs=false ./bin/golangci-lint run'
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

    stage('QA system tests') {
      steps {
        container('golang') {
          sh 'curl --fail --max-time 30 --connect-timeout 10 https://rec-vendor-api-stg.arepa.appier.info/healthz'
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
