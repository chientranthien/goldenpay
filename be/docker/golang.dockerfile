FROM golang:1.22

RUN apt update -y
RUN apt install iputils-ping -y
RUN apt install telnet -y
RUN apt install dnsutils -y
