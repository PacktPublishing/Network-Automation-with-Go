FROM debian:stretch-slim

RUN apt update && apt install -y tcpdump procps iproute2 ssh net-tools wget unzip

RUN wget https://github.com/microsoft/ethr/releases/download/v1.0.0/ethr_linux.zip && unzip ethr_linux.zip && rm ethr_linux.zip

RUN wget https://github.com/cloudprober/cloudprober/releases/download/v0.11.4/cloudprober-v0.11.4-linux-x86_64.zip && \
         unzip cloudprober-v0.11.4-linux-x86_64.zip && \
         rm cloudprober-v0.11.4-linux-x86_64.zip