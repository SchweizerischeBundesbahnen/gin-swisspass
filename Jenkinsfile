#!groovy
@Library(['pipeline-helper@master', 'esta-cloud-pipeline@master']) _

// we can't use the esta-cloud-pipline (yet) because in this repo we don't have a docker-image to build
pipeline {
    agent {
        label 'java'
    }
    parameters {
        string(name: 'DO_RELEASE', defaultValue: '', description: 'Type DORELEASE to create a release-build from master')
    }
    stages {
        stage('On all branches, just build') {
            steps {
                script {
                    bin_goBuildAndPublishSnapshot(
                        targetRepo: "sjm.go",
                        useMage: true,
                        additionalGoParams: "cibuild",
                        deployArtifacts: false)
                }
            }
        }
        stage('When on master, do Sonar analysis, create image and deploy') {
            when {
                branch 'master'
                expression { 'DORELEASE' != params.DO_RELEASE }
            }
            steps {
                script {
                    sh(script: "git fetch --tags --force -p -P")
                    String lastTag = sh(returnStdout: true, script: "git tag --merged origin/master '*[0-9]*.[0-9]*.[0-9]*' --sort 'v:refname' | tail -n1")
                    lastTag = lastTag.trim()
                    if (lastTag == null || lastTag == 'null' || lastTag == '') {
                        lastTag = 'v0.0.1'
                    }
                    echo "last tag: ${lastTag}"

                    String[] versionParts = lastTag.split(/\./)
                    int minorPlusOne = versionParts[1].toInteger()+1
                    def nextVersion = "${versionParts[0]}.${minorPlusOne}.${versionParts[2]}"
					
                    String builtVersion = bin_goBuildAndPublishSnapshot(
                        releaseVersion: "${nextVersion}",
                        targetRepo: "sjm.go",
                        useMage: true,
                        additionalGoParams: "cibuild",
                        deployArtifacts: true)

                    bin_goSonarAnalysis()
                }
            }
        }
        stage('When on master and release-parameter is set, create release') {
            when {
                branch 'master'
                expression { 'DORELEASE' == params.DO_RELEASE }
            }
            steps {
                script {
                    String newVersion = cloud_gitRelease(moveLatestTag: false, useGoVersionPrefix: true)

                    bin_goBuildAndPublish(
                            releaseVersion: "${newVersion}",
                            targetRepo: "sjm.go",
                            useMage: true,
                            additionalGoParams: "cibuild")
                }
            }
        }
    }
}
