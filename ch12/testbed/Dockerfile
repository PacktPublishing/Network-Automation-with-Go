FROM docker.io/redhat/ubi8-minimal:8.6

ARG COLLECTIONS_DIR=/usr/share/ansible/collections
ARG ANSIBLE_DIR=/etc/ansible

ENV PATH=/root/.local/bin:$PATH

RUN mkdir -p $COLLECTIONS_DIR
RUN mkdir -p $ANSIBLE_DIR

RUN microdnf --nodocs install git python38 \
    && alternatives --set python /usr/bin/python3 \
    && pip3.8 install --upgrade pip ansible-core==2.12.7 setuptools==62.1.0 botocore==1.24.46 boto3==1.21.46 cryptography==36.0.1 \
    && ansible-galaxy collection install amazon.aws:==3.4.0 -p $COLLECTIONS_DIR \
    && ansible-galaxy collection install community.crypto:==2.1.0 -p $COLLECTIONS_DIR \
    && ansible-galaxy collection install ansible.posix:==1.4.0 -p $COLLECTIONS_DIR \
    && ansible-galaxy collection install community.general:==5.1.1 -p $COLLECTIONS_DIR \
    && microdnf clean all \
    && rm -rf /tmp/*

COPY ssh_config /etc/ssh/ssh_config.d/testbed.conf
COPY ansible.cfg $ANSIBLE_DIR