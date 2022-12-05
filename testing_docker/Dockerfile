FROM ubuntu:focal
ARG AUTHORIZED_KEY
ARG DEBIAN_FRONTEND=noninteractive
WORKDIR /root/


RUN apt update -y -qq
RUN apt upgrade -y -qq
RUN apt install -y -qq openssh-server 

RUN service ssh start

RUN mkdir -p ~/.ssh
RUN chmod  0700 ~/.ssh
RUN echo $AUTHORIZED_KEY >> ~/.ssh/authorized_keys
RUN chmod 0600 ~/.ssh/authorized_keys

EXPOSE 22

CMD ["bash"] 
