FROM golang:1.14.8-buster
WORKDIR /go/src
LABEL maintainer="demo"

USER root
ENV DEBIAN_FRONTEND "noninteractive"
RUN apt update && apt install -qy \
        curl \
        joe \
        docker.io && \ 
    curl -L "https://github.com/docker/compose/releases/download/1.26.0/docker-compose-$(uname -s)-$(uname -m)"\
         -o /usr/local/bin/docker-compose && \
    chmod +x /usr/local/bin/docker-compose && \
    apt clean && \
    apt autoremove

COPY misc/.bashrc /root

ENTRYPOINT ["/bin/bash"]
