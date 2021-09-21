## What You Need Create a Testing Environment

Before you run the playbook to create a Linux testing environment, you need: 

1. An AWS Free Tier account
2. A computer with Git, Python3, pip, and Ansible 2.9 or later installed.

## Creating a Testing Environment

The testing environment is a single Linux instance in AWS running Docker to create container-based network topologies. The playbook also creates all additional AWS logical resources necessary to provision this VM; a Subnet, VPC, Security group, SSH Key pair, and Internet gateway.

Before you run the playbook, you need to make your AWS account credentials (`AWS_ACCESS_KEY` and `AWS_SECRET_KEY`) available as environment variables with the `export` command.

```bash
export AWS_ACCESS_KEY=’…’
export AWS_SECRET_KEY=’…’
```

The next step is to execute the playbook with the `ansible-playbook` command.

```bash
$  ansible-playbook create-EC2-testbed.yml -v

<snip>

TASK [Print out SSH access details] ***********************************************************************************************************************************************************
ok: [testbed-fedora34] => {
    "msg": "ssh -i testbed-private.pem fedora@ec2-54-175-179-XXX.compute-1.amazonaws.com"
}

RUNNING HANDLER [configure_instance : Reboot machine] *****************************************************************************************************************************************
changed: [testbed-fedora34] => {"changed": true, "elapsed": 22, "rebooted": true}

PLAY RECAP ************************************************************************************************************************************************************************************
localhost                  : ok=25   changed=7    unreachable=0    failed=0    skipped=1    rescued=0    ignored=0   
testbed-fedora34           : ok=33   changed=22   unreachable=0    failed=0    skipped=0    rescued=0    ignored=0    
```

## Delete All Resources

As important or even more important that automating the VM setup process, is automating how you delete all cloud resources afterwards, so you don't pay for something you might no longer need. 

You can do this by running the delete playbook with the `ansible-playbook` command. You need to provide the AWS region in case you didn't use the default value.

```bash
ansible-playbook delete-EC2-testbed.yml -v --extra-vars "aws_region=us-east-2"

<snip>

TASK [aws_delete_resources : Delete SSH Key Pair for instance] ********************************************************************************************
changed: [localhost] => {
    "changed": true,
    "key": null
}

MSG:

key deleted

PLAY RECAP ************************************************************************************************************************************************
localhost                  : ok=19   changed=7    unreachable=0    failed=0    skipped=2    rescued=0    ignored=0   
```