FROM debian:stretch-slim

RUN apt update && apt install -y tcpdump procps iproute2 ssh net-tools wget unzip

RUN wget https://github.com/microsoft/ethr/releases/download/v1.0.0/ethr_linux.zip && unzip ethr_linux.zip && rm ethr_linux.zip