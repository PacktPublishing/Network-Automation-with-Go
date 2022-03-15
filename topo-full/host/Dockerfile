FROM debian:stretch-slim

RUN apt update && apt install -y tcpdump procps iproute2 ssh net-tools wget unzip gcc git curl

RUN wget https://github.com/microsoft/ethr/releases/download/v1.0.0/ethr_linux.zip && unzip ethr_linux.zip && rm ethr_linux.zip

RUN wget https://github.com/cloudprober/cloudprober/releases/download/v0.11.5/cloudprober-v0.11.5-linux-x86_64.zip && \
         unzip cloudprober-v0.11.5-linux-x86_64.zip && \
         mv cloudprober-v0.11.5-linux-x86_64/cloudprober . && \
         rm -rf cloudprober-v0.11.5-linux-* 

RUN wget https://go.dev/dl/go1.17.7.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.7.linux-amd64.tar.gz && \
    rm go1.17.7.linux-amd64.tar.gz && \
    /usr/local/go/bin/go install github.com/go-delve/delve/cmd/dlv@v1.8.1 && \
    echo "PATH=$PATH:~/go/bin:/usr/local/go/bin" >>  ~/.bashrc


RUN bash -c "$(curl -sL https://get-gnmic.kmrd.dev)" && \
    /usr/local/go/bin/go install github.com/openconfig/gnmi/cmd/gnmi_cli@latest
