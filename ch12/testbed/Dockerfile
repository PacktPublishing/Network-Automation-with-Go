FROM docker.io/redhat/ubi8-minimal:8.5

ARG TESTBED_DIR=/Network-Automation-with-Go/ch12/testbed
ARG COLLECTIONS_DIR=$TESTBED_DIR/collections
ARG BOOK_REPO=https://github.com/PacktPublishing/Network-Automation-with-Go

ENV PATH=/root/.local/bin:$PATH

RUN microdnf --nodocs install git python38 \
    && alternatives --set python /usr/bin/python3 \
    && pip3.8 install --user --upgrade pip setuptools ansible>=2.11 botocore>=1.19.0 boto3>=1.16.0 cryptography \
    && git clone $BOOK_REPO \
    && ansible-galaxy collection install amazon.aws:==3.0.0 -p $COLLECTIONS_DIR \
    && ansible-galaxy collection install community.crypto:==2.1.0 -p $COLLECTIONS_DIR \
    && microdnf clean all \
    && rm -rf /tmp/*

WORKDIR $TESTBED_DIR