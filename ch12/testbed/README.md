# Setting Up a Testing Environment

You can run the Go code example of the first part of the book in the [Go Playground][] or on your personal computer if you have already installed Go. Code examples in chapters six and later interact with a virtual network topology, which we run in containers with Containerlab. 

## What Is a Testing Environment?

A test environment is the hardware and software that meets the minimum requirements to execute test cases.

For the hardware, we rely on a cloud provider to have a common denominator. In this case, we leverage Amazon Web Services (AWS) to provision a virtual machine (VM), where we can install all the software we need to run the code examples of the book. We call this virtual machine an EC2 instance, or just an instance (we use these terms interchangeably). By default, the instance type we run in AWS is a `t2.micro`, which you can run for free as part of the [AWS Free Tier][], but we recommend you run at least a `t2.medium`, ideally a `t2.large` size instance that has more CPU and memory, to support all the virtual network devices in the lab topology.

An Ansible playbook describes and automates all the tasks required to create this VM in AWS, as well as the tasks that prescribe the software---such as Containerlab---that needs to be present in the VM. It also defines how to configure the software. You can find this playbook in the [book's GitHub repository]. Another playbook is also available to delete the resources you create when you no longer need them.

The playbooks in the book are:
- [create-EC2-testbed.yml](create-EC2-testbed.yml):  Creates the testbed. It takes around 10 minutes to run.
- [delete-EC2-testbed.yml](delete-EC2-testbed.yml): Deletes the resources you create when you no longer need them.

If you have Docker installed in your computer, you can build the AWS test setup with the command `make build-env`.

## What You Need to Create a Testing Environment

Before you run the playbook in your computer to create a Linux testing environment in AWS, you need to make sure you meet the following requirements:

1. An AWS account. [AWS Free Tier] is fine.
2. A computer with: 
    * Git
    * Docker
    * Make

With all this in place, you can go ahead and clone the [book's GitHub repository][] with the `git clone` command.

```bash
$ git clone https://github.com/PacktPublishing/Network-Automation-with-Go
```

## Creating the Testing Environment

The testing environment is a single Linux instance in AWS running Containerlab to create container-based network topologies. The next figure illustrates what the environment looks like. 

<p align="center">
  <img width=800 title="Target Environment" src="pictures/ch12-test_env.png"><br>
</p>

The first step is to make your AWS account credentials (`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`) available as environment variables with the `export` command. Check out [AWS Programmatic access][] to create an access key, and to save your secret access key.

```bash
export AWS_ACCESS_KEY_ID=’…’
export AWS_SECRET_ACCESS_KEY=’…’
```

Then, from the book repository directory, run the command `make env-build`.

```bash
Network-Automation-with-Go$ make env-build
AWS_ACCESS_KEY_ID is AKIAVFPUEFZCFVFGXXXX
AWS_SECRET_ACCESS_KEY is **************************
Using /etc/ansible/ansible.cfg as config file

PLAY [Create EC2 instance] *************************************************************************************************************************************************************

TASK [Gathering Facts] *****************************************************************************************************************************************************************
ok: [localhost]

### ... <omitted for brevity > ... ###

TASK [Print out instance information for the user] *************************************************************************************************************************************
ok: [testbed] => {}

MSG:

['SSH: ssh -i lab-state/id_rsa fedora@ec2-54-86-51-96.compute-1.amazonaws.com\n', 'To upload cEOS image: scp -i lab-state/id_rsa ~/Downloads/cEOS64-lab-4.28.0F.tar fedora@ec2-54-86-51-96.compute-1.amazonaws.com:./network-automation-with-go\n']

PLAY RECAP *****************************************************************************************************************************************************************************
localhost                  : ok=28   changed=9    unreachable=0    failed=0    skipped=3    rescued=0    ignored=0   
testbed                    : ok=36   changed=24   unreachable=0    failed=0    skipped=11   rescued=0    ignored=0  
```

You can find the VM access details in the logs, as the preceding output shows. Look for something like: `ssh -i lab-state/id_rsa fedora@ec2-54-86-51-96.compute-1.amazonaws.com` or `ssh -i lab-state/id_rsa ubuntu@ec2-54-86-51-96.compute-1.amazonaws.com` depending on the target Linux distribution.

## Virtual Machine Options

### Instance Size

AWS offers different instance types. You can select any type you prefer, based on your vCPU/Memory preference, and price constraints. By default, the playbook selects a `t2.micro` instance, which is the only free option you have as part of the [AWS Free Tier][]. You can check the hourly pricing for other instance types at [On-Demand Plans for Amazon EC2][]. The next table shows some examples.

Instance name | On-Demand hourly rate | vCPU | Memory
--- | --- | --- | ---
t3.medium | $0.0376 | 2 | 4 GiB
t3.large | $0.0832 | 2 | 8 GiB
t2.large | $0.0928 | 2 | 8 GiB
m5.large | $0.096 | 2 | 8 GiB
t3.xlarge | $0.1664 | 4 | 16 GiB
t2.xlarge | $0.1856| 4 | 16 GiB
m4.xlarge | $0.2 | 4 | 16 GiB
r5.xlarge | $0.252 | 4 | 32 GiB
m5.2xlarge | $0.384 | 8 | 32 GiB
r5.2xlarge | $0.504 | 8 | 64 GiB
c5.metal | $4.08 | 96 | 192 GiB

To run the testing environment on a `m5.large` instance, you need to export the environment variables `VM_SIZE` with the value `m5.large`, like in the next output.

```bash
$ export VM_SIZE=m5.large
$ make env-build
```

### AWS Region

We recommend you launch the instance in the AWS region ([EC2 Available Regions][]) that is closer to your current location. By default, the playbook selects `us-east-1` and you can go with it. The next table shows other regions you can choose from if you prefer so.

Code | Region Name
--- | ---
us-east-1| US East (N. Virginia)
us-east-2| US East (Ohio)
us-west-1| US West (N. California)
eu-west-2| EU West (London)
eu-central-1| Europe (Frankfurt)
sa-east-1| South America (São Paulo)
ca-central-1| Canada (Central)
ap-northeast-1| Asia Pacific (Tokyo)
ap-southeast-2| Asia Pacific (Sydney)
ap-south-1| Asia Pacific (Mumbai)

To run the testing environment in London, you need to export the environment variables `AWS_REGION` with the value `eu-west-2`, like in the next output

```bash
$ export AWS_REGION=eu-west-2
$ make env-build
```

### Linux Distribution

If you have a preference between Fedora vs Debian based Linux distributions, you have the option to run the testing environment on either Fedora (35) or Ubuntu (22.04). Export the environment variables `AWS_DISTRO` to select one or the other. The default option is `fedora`.

The next example shows how you can create a `t2.medium` instance running `ubuntu` in Ohio.

```bash
$ export AWS_DISTRO=ubuntu
$ export VM_SIZE=t2.medium
$ export AWS_REGION=us-east-2
$ make env-build
```

## Connecting to the Test VM

After you create the instance, you can connect to it using the info provided in the logs. The playbook generates an SSH private key (`lab-state/id_rsa`), which we use to authenticate to the test VM. Connect to the VM and verify that Go is present in the system with the `go version` command. 

```bash
fedora@testbed ~ ⇨  go version
go version go1.18.1 linux/amd64
```

## Uploading Container Images to the Test VM

Some networking vendors make it simpler than others to access their container-based network operating systems (NOS). If you can't pull the image directly from a container registry, like Docker Hub, you might need to download the image from their website and upload it to the test VM. Keep in mind most container images might require more resources that what a `t2.micro` instance can offer.

Let's pretend you downloaded a cEOS image (`cEOS64-lab-4.28.0F.tar`) to your Downloads folder. You can copy the image to the test VM with the `scp` command using the generated SSH private key. See an example next or check [Get Arista cEOS](get_arista_ceos.md)

```bash
Network-Automation-with-Go$ scp -i lab-state/id_rsa ~/Downloads/cEOS64-lab-4.28.0F.tar fedora@ec2-54-86-51-96.compute-1.amazonaws.com:./network-automation-with-go
cEOS64-lab-4.28.0F.tar                        100%  434MB  26.6MB/s   00:16
```

Then, SSH to the instance and import the image with the `docker` command.

```bash
Network-Automation-with-Go$ ssh -i lab-state/id_rsa fedora@ec2-54-86-51-96.compute-1.amazonaws.com
fedora@testbed ~ ⇨  cd network-automation-with-go && make clone
docker import cEOS64-lab-4.28.0F.tar ceos:4.28
sha256:dcdc721054804ed4ea92f970b5923d8501c28526ef175242cfab0d158ac0085c
```

You can now reference this image (`ceos:4.28`) in the `image` section of one or more routers in the topology file.

```bash
ubuntu@testbed topo ⇨  docker exec -it clab-netgo-ceos Cli
ceos>show ver | i Software
Software image version: 4.28.0F-26924507.4280F (engineering build)
ceos>
```

## Building a Virtual Network Topology

The Linux environment comes with [Containerlab] in it. We use this application to wire together different containerized NOS and create a virtual network topology we can interact with to run the book code examples. [Containerlab] offers a hassle-free and quick way to define and run network topologies. It also worth mentioning that Go is the programming language of choice for this tool.

You can find the topology definition file in the `topo-base` folder of [book's GitHub repository].

```yaml
name: netgo

topology:
  nodes:
    srl:
      kind: srl
      image: ghcr.io/nokia/srlinux:21.6.4
    ceos:
      kind: ceos
      image: ceos:4.28
    cvx:
      kind: cvx
      image: networkop/cx:5.0.0
      runtime: docker

  links:
    - endpoints: ["srl:e1-1", "ceos:eth1"]
    - endpoints: ["cvx:swp1", "ceos:eth2"]
```

This topology file defines a three node topology as the next figure shows. One node runs Nokia SR Linux, another NVIDIA Cumulus Linux and the last one Arista's cEOS. We chose the first two, because you can conveniently get their images from a public container registry. You can download cEOS from Arista's website after creating a free account.

<p align="center">
  <img width=600 title="Network Topology" src="pictures/ch6-topo.png"><br>
</p>


## Launching a Virtual Network Topology

To launch the virtual topology file, run the command `make make lab-up` from the $HOME/network-automation-with-go in the VM, as the next output shows.

```bash
fedora@testbed network-automation-with-go$ make lab-up
INFO[0000] Containerlab v0.25.1 started                 
INFO[0000] Parsing & checking topology file: topo.yml   
INFO[0000] Removing /home/fedora/network-automation-with-go/clab-netgo directory... 
INFO[0000] Could not read docker config: open /root/.docker/config.json: no such file or directory 
INFO[0000] Pulling docker.io/networkop/cx:5.0.0 Docker image 
INFO[0021] Done pulling docker.io/networkop/cx:5.0.0    
INFO[0021] Could not read docker config: open /root/.docker/config.json: no such file or directory 
INFO[0021] Pulling ghcr.io/nokia/srlinux:21.6.4 Docker image 
INFO[0039] Done pulling ghcr.io/nokia/srlinux:21.6.4    
WARN[0039] it appears that container host has low memory available: ~0Gi. This might lead to runtime errors. Consider freeing up more memory. 
INFO[0039] Creating lab directory: /home/fedora/network-automation-with-go/clab-netgo 
INFO[0040] Creating docker network: Name="clab", IPv4Subnet="172.20.20.0/24", IPv6Subnet="2001:172:20:20::/64", MTU="1500" 
INFO[0040] Creating container: "ceos"                   
INFO[0040] Creating container: "cvx"                    
INFO[0041] Creating container: "srl"                    
INFO[0043] Creating virtual wire: cvx:swp1 <--> ceos:eth2 
INFO[0046] Creating virtual wire: srl:e1-1 <--> ceos:eth1 
INFO[0046] Running postdeploy actions for Arista cEOS 'ceos' node 
INFO[0046] Running postdeploy actions for Nokia SR Linux 'srl' node 
INFO[0107] Adding containerlab host entries to /etc/hosts file 
INFO[0107] 🎉 New containerlab version 0.26.2 is available! Release notes: https://containerlab.dev/rn/0.26/#0262
Run 'containerlab version upgrade' to upgrade or go check other installation options at https://containerlab.dev/install/ 
+---+-----------------+--------------+------------------------------+------+---------+----------------+----------------------+
| # |      Name       | Container ID |            Image             | Kind |  State  |  IPv4 Address  |     IPv6 Address     |
+---+-----------------+--------------+------------------------------+------+---------+----------------+----------------------+
| 1 | clab-netgo-ceos | 53274759376c | ceos:4.28                    | ceos | running | 172.20.20.3/24 | 2001:172:20:20::3/64 |
| 2 | clab-netgo-cvx  | 404c58026c11 | networkop/cx:5.0.0           | cvx  | running | 172.20.20.2/24 | 2001:172:20:20::2/64 |
| 3 | clab-netgo-srl  | 16a5cf8add7a | ghcr.io/nokia/srlinux:21.6.4 | srl  | running | 172.20.20.4/24 | 2001:172:20:20::4/64 |
+---+-----------------+--------------+------------------------------+------+---------+----------------+----------------------+
```

You now have routers `clab-netgo-ceos`, `clab-netgo-cvx` and `clab-netgo-srl` ready to go.

### Connecting to the Routers

One of the first changes network engineers notice when they embark on their network automation journey, is they don't need to connect to individual devices too often, as you can perform most of the tasks via programming interfaces instead.

Some code examples in this book take advantage of these interfaces. Still, do not hesitate to login to the network elements via the CLI interface, to check the result of running a code example. This is a good way to build confidence in this novel approach to manage and operate networks.

[Containerlab] uses Docker to run the containers. This means we can leverage some of Docker capabilities to connect to the routers, for example the `docker exec` command with the name of the container, and corresponding command-line interface process.

```bash
fedora@testbed ~ ⇨  docker exec -it clab-netgo-srl sr_cli
Welcome to the srlinux CLI.                      
A:srl# show version | grep Software
Software Version  : v21.6.4
```

`sr_cli` in the preceding example is the command-line interface process for an SR Linux device. Other examples in the next table.

NOS | Command | 
--- | --- |
CVX | vtysh
SR Linux | sr_cli
EOS | Cli

You can also SSH to the same device. Use the `ssh` command with the credentials in the next table.

| Device | Username | Password |
| ------ | -------- | -------- |
| clab-netgo-srl | admin | admin |
| clab-netgo-ceos | admin | admin |
| clab-netgo-cvx | cumulus | cumulus |

For example: 

```bash
fedora@testbed ~ ⇨  ssh admin@clab-netgo-ceos
(admin@clab-netgo-ceos) Password: admin
ceos>en
ceos#
```

### Destroying the Network Topology

You can destroy the network topology using the `make lab-down` command.

```bash
fedora@testbed network-automation-with-go$ make lab-down
```

## Delete All Resources

As important or even more important that automating the VM setup process, is automating how you delete all cloud resources afterwards, so you don't pay for something you might no longer need. 

```bash
Network-Automation-with-Go$ make env-delete
AWS_ACCESS_KEY_ID is AKIAVFPUEFZCFVFGXXXX
AWS_SECRET_ACCESS_KEY is **************************
[WARNING]: No inventory was parsed, only implicit localhost is available
[WARNING]: provided hosts list is empty, only localhost is available. Note that the implicit localhost does not match 'all'

PLAY [Delete EC2 instance] *************************************************************************************************************************************************************

TASK [Gathering Facts] *****************************************************************************************************************************************************************
ok: [localhost]

### ... <omitted for brevity > ... ###

TASK [aws_delete_resources : Delete SSH Key Pair for instance] *************************************************************************************************************************
changed: [localhost]

TASK [Cleanup state files] *************************************************************************************************************************************************************
changed: [localhost] => (item=.region)
changed: [localhost] => (item=.vm)

PLAY RECAP *****************************************************************************************************************************************************************************
localhost                  : ok=21   changed=8    unreachable=0    failed=0    skipped=3    rescued=0    ignored=0   
```

## Other Testing Options

Not every networking offers public access to the container images of their network operating systems. We aim to make the examples as useful to you as possible. For this reason, we also take advantage of a couple of additional resources to get access to networking operating systems you might run in your organization.

The [DevNet Sandbox] is a **free** option to access always-on devices that we target in some examples in the book. They have Cisco Nexus, IOS XR and IOS shared devices. Keep in mind their hostname/fqdn, and credentials might change in the future. You can also [reserve a DevNet Sandbox].

Other great resources to run virtual network topologies are [GNS3], [EVE-NG], [netsim-tools], and [Vagrant]. You need to have a contract with a networking vendor company to get access to their virtual images to run on any of these though.

Last, but not least, [Cisco Modeling Labs] offer access to Cisco virtual images to create network simulations with this tool. The personal license is available for 199 dollars a year ([Cisco Modeling Labs - Personal]).

<!-- links -->
[AWS Free Tier]: https://aws.amazon.com/free/
[AWS Programmatic access]: https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys
[On-Demand Plans for Amazon EC2]: https://aws.amazon.com/ec2/pricing/on-demand/
[Containerlab]: https://github.com/srl-labs/containerlab
[Cisco Modeling Labs]: https://developer.cisco.com/modeling-labs/
[netsim-tools]: https://netsim-tools.readthedocs.io/en/latest/
[GNS3]: https://www.gns3.com/
[EVE-NG]: https://www.eve-ng.net/
[Vagrant]: https://www.vagrantup.com/
[Add cEOS]: https://github.com/nleiva/aws-testbed/blob/main/lab/get_arista_ceos.md#add-image-to-your-local-image-repository
[Arista cEOS in Containerlab]: https://containerlab.dev/manual/kinds/ceos/#arista-ceos
[DevNet Sandbox]: https://developer.cisco.com/site/sandbox/
[reserve a DevNet Sandbox]: https://developer.cisco.com/docs/sandbox/#!first-reservation-guide/reservation-hello-world
[EC2 Available Regions]: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html#concepts-available-regions
[Cisco Modeling Labs - Personal]: https://learningnetworkstore.cisco.com/cisco-modeling-labs-personal/cisco-cml-personal
[book's GitHub repository]: https://github.com/PacktPublishing/Network-Automation-with-Go
[AWS Free Tier]: https://aws.amazon.com/free/
[Beginners Guide: Downloading Python]: https://wiki.python.org/moin/BeginnersGuide/Download
[Installing Ansible with pip]: https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html#installing-ansible-with-pip
[Getting Started - Installing Git]: https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
[Installing pip - Supported Methods]: https://pip.pypa.io/en/stable/installation/#supported-methods
[AWS Programmatic access]: https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys
[Go Playground]: https://play.golang.org/