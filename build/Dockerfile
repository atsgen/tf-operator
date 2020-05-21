#FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
FROM alpine as runtime

ENV OPERATOR=/usr/local/bin/tf-operator \
    USER_UID=1001 \
    USER_NAME=tf-operator

# install operator binary
COPY build/_output/bin/tf-operator ${OPERATOR}

COPY bindata /bindata
COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
