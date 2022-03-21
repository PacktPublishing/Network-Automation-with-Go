# cEOS image

## Create an account

Go to [User Registration](https://www.arista.com/en/user-registration) in their website.


If you use a nom-corporate email to register, you will get this message: "_This Username Available. (Note: You are trying to register with a non corporate email ID. In order to get customer or partner access, you must provide a corporate email ID.)_"

After registration, wait five minutes to get an email for account confirmation, and then another five to login.

## Download the container image

Once your account is confirmed, you can log in and access the [Software Download section](https://www.arista.com/en/support/software-download).

<p align="center">
  <img height="200" title="Software Download section" src="pictures/Arista_Download.png"><br>
  <b>Arista Software Download section</b><br>
</p>

From there, you can download cEOS.

<p align="center">
  <img height="600" title="Download cEOS" src="pictures/Download_cEOS.png"><br>
  <b>Download cEOS</b><br>
</p>

## Add image to your local image repository

Upload container image from your computer to EC2 instance.

```bash
scp -i lab-state/id_rsa cEOS64-lab-4.26.4M.tar fedora@ec2-3-86-163-31.compute-1.amazonaws.com:.
```

```bash
$ docker import cEOS64-lab-4.26.4M.tar ceos:4.26.4M
Getting image source signatures
Copying blob 474b2ea4514a [--------------------------------------] 0.0b / 0.0b
Copying config 0134ea0180 done  
Writing manifest to image destination
Storing signatures
sha256:0134ea0180b1fcf39e0a3d00abf6c698f034874320f42a2fc6f6bd1c059fe327
```

Note on [cEOS and cgroups v2](https://github.com/srl-labs/containerlab/issues/467)