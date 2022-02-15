FROM docker.io/redhat/ubi8-minimal:8.5

ARG TESTBED_DIR=/Network-Automation-with-Go/ch12/testbed
ARG COLLECTIONS_DIR=$TESTBED_DIR/collections
ARG BOOK_REPO=https://github.com/PacktPublishing/Network-Automation-with-Go

ENV PATH=/root/.local/bin:$PATH

RUN microdnf --nodocs install git python38 \
    && alternatives --set python /usr/bin/python3 \
    && pip3.8 install --upgrade pip ansible-core==2.12.2 setuptools==60.9.1 botocore==1.24.0 boto3==1.21.0 cryptography==36.0.1 \
    && git clone $BOOK_REPO \
    && ansible-galaxy collection install amazon.aws:==3.0.0 -p $COLLECTIONS_DIR \
    && ansible-galaxy collection install community.crypto:==2.1.0 -p $COLLECTIONS_DIR \
    && microdnf clean all \
    && rm -rf /tmp/*

RUN cp $TESTBED_DIR/config ~/.ssh/config

WORKDIR $TESTBED_DIR