FROM registry.access.redhat.com/openshift3/jenkins-slave-base-rhel7
RUN mkdir gobin
EXPOSE 8080
ADD http-server /gobin
CMD [ "/gobin/http-server" ]
