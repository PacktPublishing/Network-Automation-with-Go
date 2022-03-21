FROM docker.io/redhat/ubi8-minimal:8.5

ARG COLLECTIONS_DIR=/usr/share/ansible/collections

ENV PATH=/root/.local/bin:$PATH

RUN mkdir -p $COLLECTIONS_DIR

RUN microdnf --nodocs install git python38 \
    && alternatives --set python /usr/bin/python3 \
    && pip3.8 install --upgrade pip ansible-core==2.12.2 setuptools==60.9.1 botocore==1.24.0 boto3==1.21.0 cryptography==36.0.1 \
    && ansible-galaxy collection install amazon.aws:==3.0.0 -p $COLLECTIONS_DIR \
    && ansible-galaxy collection install community.crypto:==2.1.0 -p $COLLECTIONS_DIR \
    && ansible-galaxy collection install ansible.posix:==1.3.0 -p $COLLECTIONS_DIR \
    && ansible-galaxy collection install community.general:==4.6.1 -p $COLLECTIONS_DIR \
    && microdnf clean all \
    && rm -rf /tmp/*

COPY ch12/testbed/config /etc/ssh/ssh_config.d/testbed.conf