podTemplate(label: 'mypod', cloud:'openshift',containers: [
    containerTemplate(name: 'golang', image: 'golang:1.9.2', ttyEnabled: true, command: 'cat')
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
                    """
                }
            }
        }       
    }
}
