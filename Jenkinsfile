podTemplate(label: 'mypod', serviceAccount: 'jenkins', cloud:'openshift',containers: [
    containerTemplate(name: 'golang', image: 'golang:1.9.2', ttyEnabled: true, command: 'cat'),
    containerTemplate(name: 'oc', image:'registry.access.redhat.com/openshift3/jenkins-slave-base-rhel7:latest',
                    ttyEnabled: true, command: 'cat')
  ],volumes: [
    emptyDirVolume(mountPath: '/opt/binary', memory: false),
  ]) {
    node('mypod') {
        stage('Get a Golang project') {
            checkout scm
            container('golang') {
                stage('Build a Go project') {
                    sh """
                    mkdir -p /go/src/github.com/jmcshane
                    ln -s `pwd` /go/src/github.com/jmcshane/hipchat-openshift
                    cd /go/src/github.com/jmcshane/hipchat-openshift
                    go test ./...
                    go build .
                    ls -al
                    cp /go/src/github.com/jmcshane/hipchat-openshift/hipchat-openshift /opt/binary/main
                    """
                }
            }
            stage('Send Binary to build') {
              container('oc') {
                  sh "sleep 20"
                  sh "oc start-build server --from-file=/opt/binary/main"
                }
            }
        }
    }
}
