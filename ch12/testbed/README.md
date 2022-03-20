# Setting Up a Testing Environment

You can run the Go code snippets in the first part of this book in the [Go Playground][] or on your personal computer if you have already installed Go. Some code examples in chapters six and later interact with other systems—such as virtual network devices—which we run as containers with Docker. Now, instead of asking you to install Docker on your computer to run these examples, we think is more practical for you to recreate the Linux environment we actually use to write and test the code examples of these chapters. This testing environment comes with Docker and all other software dependencies installed, so you can run the examples without a hitch. We want to make sure you have a pleasant experience running the examples.

## What Is a Testing Environment?

A test environment is the hardware and software that meets the minimum requirements to execute test cases.

For the hardware, we rely on a cloud provider to have a common denominator. In this case, we leverage Amazon Web Services (AWS) to provision a virtual machine (VM), where we can install all the software we need to run the code examples of the book. We call this virtual machine an EC2 instance, or just an instance (we use these terms interchangeably). By default, the instance type we run in AWS is a `t2.micro`, which you can run for free as part of the [AWS Free Tier][], but we recommend you run at least a `t2.medium`, ideally a `t2.large` size instance that has more CPU and memory, to support all the virtual network devices we need to run.

An Ansible playbook describes and automates all the tasks required to create this VM in AWS, as well as the tasks that prescribe the software—such as Docker—that needs to be present in the VM and it defines how to configure the software as well. You can find this playbook in the [book's GitHub repository]. Another playbook is also available to delete the resources you create when you no longer need them.

The playbooks in the book are:
- [create-EC2-testbed.yml](create-EC2-testbed.yml):  Creates the testbed. It takes around 10 minutes to run.
- [delete-EC2-testbed.yml](delete-EC2-testbed.yml): Deletes the resources you create when you no longer need them.

If you have Docker installed, you can trigger the AWS test setup from a container with `make env-build`, without having to install the dependencies we cover next.

## What You Need Create a Testing Environment

Before you run the playbook in your computer to create a Linux testing environment in AWS, you need to make sure you meet the following requirements:

1. An AWS account. [AWS Free Tier] is fine.
2. A computer with: 
    * Git --- if you don't have installed already, check out their [Getting Started - Installing Git][] guide. Example: `sudo apt install git-all` (Ubuntu)
    * Python3 --- they provide installation instructions in the [Beginners Guide: Downloading Python][]. Example: `sudo apt install python3.8` (Ubuntu)
    * pip --- if your Python environment does not have pip installed, they offer instruction at [Installing pip - Supported Methods][]. Example: `python3 -m pip install --upgrade pip`
    * Ansible 2.9 --- you can install Ansible with pip as they describe it in [Installing Ansible with pip][]. Example: `python3 -m pip install --user ansible==2.9.26` 

With all this in place, you can go ahead and clone the [book's GitHub repository][] with the `git clone` command.

```bash
$ git clone https://github.com/PacktPublishing/Network-Automation-with-Go && cd Network-Automation-with-Go/ch012/testbed
```

Now we are one step closer to execute the playbook to create the test environment and run the code examples.

## Getting Ready to Run the Playbook

Ansible executes some tasks of the playbook, for example to communication with AWS, with content or code that the authors of that content distribute in the form of an Ansible Collection. To install it, you need to execute this command from the repository folder:

```bash
ch012/testbed$ ansible-galaxy collection install -r collections/requirements.yml -p ./collections
```

Ansible Collections might depend on some Python libraries. To install these Python libraries, from the repository folder run `pip install` as in the next output: 

```bash
ch012/testbed$ pip install --user -r requirements.txt
```

You can verify `pip` installed the libraries correctly with the command `pip list`. See the following output for an example of what to expect.

```bash
ch012/testbed$ pip list | grep 'boto\|crypto'
boto                              2.49.0
boto3                             1.17.93
botocore                          1.20.112
cryptography                      3.4.7
```

## Creating the Testing Environment

The testing environment is a single Linux instance in AWS running Docker and Containerlab to create container-based network topologies. The next figure illustrates what the environment looks like. 

<p align="center">
  <img width=800 title="Target Environment" src="pictures/ch12-test_env.png"><br>
</p>

Before you run the playbook, you need to make your AWS account credentials (`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`) available as environment variables with the `export` command. Check out [AWS Programmatic access][] to create an access key, and to save your secret access key.

```bash
export AWS_ACCESS_KEY_ID=’…’
export AWS_SECRET_ACCESS_KEY=’…’
```

The next step is to execute the playbook with the `ansible-playbook` command.

```bash
ch012/testbed$ ansible-playbook create-EC2-testbed.yml --extra-vars "instance_type=t2.large" -v 

<snip>

TASK [Print out SSH access details] **********************************************************************************
ok: [testbed] => {}

MSG:

SSH: ssh -i lab-state/id_rsa fedora@ec2-3-94-8-154.compute-1.amazonaws.com
To upload cEOS image: scp -i lab-state/id_rsa ~/Downloads/cEOS64-lab-4.26.4M.tar fedora@ec2-3-94-8-154.compute-1.amazonaws.com:.

RUNNING HANDLER [configure_instance : Reboot machine] ****************************************************************
changed: [testbed] => {
    "changed": true,
    "elapsed": 47,
    "rebooted": true
}

PLAY RECAP **********************************************************************************************************************
localhost                  : ok=26   changed=1    unreachable=0    failed=0    skipped=6    rescued=0    ignored=0   
testbed                    : ok=30   changed=20   unreachable=0    failed=0    skipped=6    rescued=0    ignored=0      
```

You can find the VM access details in the logs, as the preceding output shows. Look for something similar to: `ssh -i lab-state/id_rsa fedora@ec2-3-94-8-154.compute-1.amazonaws.com` or `ssh -i lab-state/id_rsa ubuntu@ec2-3-88-180-178.compute-1.amazonaws.com`.

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

To run the testing environment on a `m5.large` instance, you need to pass the variable `instance_type` to the playbook with the value `m5.large`, like in the command example below.

```bash
ansible-playbook create-EC2-testbed.yml -v --extra-vars "instance_type=m5.large"
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

To run the testing environment in London, you need to pass the variable `aws_region` to the playbook with the value `eu-west-2`, like in the command example below.

```bash
ansible-playbook create-EC2-testbed.yml -v --extra-vars "aws_region=eu-west-2"
```

### Linux Distribution

If you have a preference between Fedora vs Debian based Linux distributions, you have the option to run the testing environment on either Fedora (34) or Ubuntu (20.04). Pass the the variable `aws_distro` to the playbook to select one or the other. The default option is `fedora`.

To run an `ubuntu` machine instead, you need to pass the variable `aws_distro` with value `ubuntu` to the `ansible playbook` command. The next example shows how you can create a `t2.medium` instance running `ubuntu` in Ohio.

```bash
$ ansible-playbook create-EC2-testbed.yml -v --extra-vars "aws_distro=ubuntu instance_type=t2.medium aws_region=us-east-2"

<skip>

TASK [Print out SSH access details] *********************************************************************************************************************************************************
ok: [testbed] => {}

MSG:

SSH: ssh -i lab-state/id_rsa ubuntu@ec2-3-142-51-83.us-east-2.compute.amazonaws.com
To upload cEOS image: scp -i lab-state/id_rsa ~/Downloads/cEOS64-lab-4.26.4M.tar ubuntu@ec2-3-142-51-83.us-east-2.compute:.

PLAY RECAP **********************************************************************************************************************
localhost                  : ok=26   changed=1    unreachable=0    failed=0    skipped=6    rescued=0    ignored=0   
testbed                    : ok=30   changed=20   unreachable=0    failed=0    skipped=6    rescued=0    ignored=0
```


## Connecting to the Test VM

After you create the instance, you can connect to it using the info provided in the logs. The playbook generates an SSH private key (`my_key`), which we use to authenticate to the test VM. Connect to the VM and verify that Go is present in the system with the `go version` command. 

```bash
fedora@testbed ~ ⇨  go version
go version go1.18beta2 linux/amd64
```

## Uploading Container Images to the Test VM

Some networking vendors make it simpler than others to access their container-based network operating systems (NOS). If you can't pull the image directly from a container registry, like Docker Hub, you might need to download the image from their website and upload it to the test VM. Keep in mind most container images might require more resources that what a `t2.micro` instance can offer.

Let's pretend you downloaded a cEOS image (`cEOS-lab-4.26.1F.tar`) to your Downloads folder. You can copy the image to the test VM with the `scp` command using the generated SSH private key. See an example next or check [Get Arista cEOS](get_arista_ceos.md)

```bash
$ scp -i lab-state/id_rsa ~/Downloads/cEOS64-lab-4.26.4M.tar fedora@ec2-3-94-8-154.compute-1.amazonaws.com:.
cEOS64-lab-4.26.4M.tar                                                                 100%  405MB  22.2MB/s   00:18
```

or

```bash
scp -i lab-state/id_rsa ~/Downloads/cEOS64-lab-4.26.4M.tar ubuntu@ec2-3-88-180-178.compute-1.amazonaws.com:.
```

Then, SSH to the instance and import the image with the `docker` command.

```bash
fedora@testbed ~ ⇨  docker import cEOS64-lab-4.26.4M.tar ceos:4.26.4M
sha256:67283d3fe45bcba28c6757cc5f1fafe63140e6cfafa005f5aa7060965d5e1bb9
```

You can now reference this image (`ceos:4.26.4M`) in the `image` section of one or more routers in the topology file.

```bash
ubuntu@testbed topo ⇨  docker exec -it clab-netgo-ceos Cli
ceos>show ver | i Software
Software image version: 4.26.4M-25280011.4264M (engineering build)
ceos>
```

## Building a Virtual Network Topology

The Linux environment comes with [Containerlab] in it. We use this application to wire together different containerized NOS and create a virtual network topology we can interact with to run the book code examples. [Containerlab] offers a hassle-free and quick way to define and run network topologies. It also worth mentioning that Go is the programming language of choice for this tool.

You can find the topology definition file in the `topo` folder of [book's GitHub repository].

```yaml
name: netgo

topology:
  nodes:
    srl:
      kind: srl
      image: ghcr.io/nokia/srlinux:21.6.4
    ceos:
      kind: ceos
      image: ceos:4.26.4M
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

...

## Launching a Virtual Network Topology

To launch the virtual topology file, run the command `make` from the home folder to access the `network-automation-with-go/topo/` folder and run `clab deploy` with root privilege, as the next output shows.

```bash
fedora@testbed ~ ⇨  make
sudo containerlab deploy -t ~/network-automation-with-go/topo-base/topo.yml --reconfigure
INFO[0000] Containerlab v0.24.1 started                 
INFO[0000] Parsing & checking topology file: topo.yml   
INFO[0000] Removing /home/fedora/network-automation-with-go/clab-netgo directory... 
INFO[0000] Could not read docker config: open /root/.docker/config.json: no such file or directory 
INFO[0000] Pulling docker.io/networkop/cx:5.0.0 Docker image 
INFO[0027] Done pulling docker.io/networkop/cx:5.0.0    
INFO[0027] Could not read docker config: open /root/.docker/config.json: no such file or directory 
INFO[0027] Pulling ghcr.io/nokia/srlinux:21.6.4 Docker image 
INFO[0052] Done pulling ghcr.io/nokia/srlinux:21.6.4    
INFO[0052] Creating lab directory: /home/fedora/network-automation-with-go/clab-netgo 
INFO[0053] Creating docker network: Name="clab", IPv4Subnet="172.20.20.0/24", IPv6Subnet="2001:172:20:20::/64", MTU="1500" 
INFO[0053] Creating container: "ceos"                   
INFO[0053] Creating container: "cvx"                    
INFO[0053] Creating container: "srl"                    
INFO[0056] Creating virtual wire: cvx:swp1 <--> ceos:eth2 
INFO[0056] Creating virtual wire: srl:e1-1 <--> ceos:eth1 
INFO[0056] Running postdeploy actions for Nokia SR Linux 'srl' node 
INFO[0056] Running postdeploy actions for Arista cEOS 'ceos' node 
INFO[0112] Adding containerlab host entries to /etc/hosts file 
+---+-----------------+--------------+------------------------------+------+---------+----------------+----------------------+
| # |      Name       | Container ID |            Image             | Kind |  State  |  IPv4 Address  |     IPv6 Address     |
+---+-----------------+--------------+------------------------------+------+---------+----------------+----------------------+
| 1 | clab-netgo-ceos | 761ab932e3f1 | ceos:4.26.4M                 | ceos | running | 172.20.20.3/24 | 2001:172:20:20::3/64 |
| 2 | clab-netgo-cvx  | 1903b058747e | networkop/cx:5.0.0           | cvx  | running | 172.20.20.2/24 | 2001:172:20:20::2/64 |
| 3 | clab-netgo-srl  | ec8b4357f5dd | ghcr.io/nokia/srlinux:21.6.4 | srl  | running | 172.20.20.4/24 | 2001:172:20:20::4/64 |
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

You can destroy the network topology using the `clab destroy` command.

```bash
$  sudo clab destroy --topo topo-base/topo.yml
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

## Other Testing Options

Not every networking offers public access to the container images of their network operating systems. We aim to make the examples as useful to you as possible. For this reason, we also take advantage of a couple of additional resources to get access to networking operating systems you might run in your organization.

- [DevNet Sandbox]: DevNet offers **free** access to always-on devices that we can target in some examples. They have Cisco Nexus, IOS XR and IOS shared devices. Keep in mind their hostname/fqdn, and credentials might change in the future. You can also [reserve a DevNet Sandbox].
- [NRE Labs]: NRE Labs is an open source educational project sponsored by Juniper. It provides free access to lab scenarios with JunOS devices.

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
[NRE Labs]: https://nrelabs.io/
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