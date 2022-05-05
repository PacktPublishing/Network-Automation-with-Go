# Terraform

```
terraform init -upgrade
terraform plan
terraform apply -auto-approve
```

## Nautobot Provider

```bash
$  make install
```

### Example

```bash
$  terraform init -upgrade

Initializing the backend...

Initializing provider plugins...
- Finding nleiva/nautobot versions matching "0.2.3"...
- Installing nleiva/nautobot v0.2.3...
- Installed nleiva/nautobot v0.2.3 (self-signed, key ID A33D26E300F155FF)

Partner and community providers are signed by their developers.
If you'd like to know more about provider signing, you can read about it here:
https://www.terraform.io/docs/cli/plugins/signing.html

Terraform has created a lock file .terraform.lock.hcl to record the provider
selections it made above. Include this file in your version control repository
so that Terraform can guarantee to make the same selections by default when
you run "terraform init" in the future.

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

```bash
$  terraform apply --auto-approve

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create
 <= read (data resources)

Terraform will perform the following actions:

  # data.nautobot_manufacturers.all will be read during apply
  # (config refers to values not yet known)
 <= data "nautobot_manufacturers" "all"  {
      + id            = (known after apply)
      + manufacturers = (known after apply)
    }

  # nautobot_manufacturer.new will be created
  + resource "nautobot_manufacturer" "new" {
      + created             = (known after apply)
      + description         = "Created with Terraform"
      + devicetype_count    = (known after apply)
      + display             = (known after apply)
      + id                  = (known after apply)
      + inventoryitem_count = (known after apply)
      + last_updated        = (known after apply)
      + name                = "New Vendor"
      + platform_count      = (known after apply)
      + slug                = (known after apply)
      + url                 = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + data_source_example = (known after apply)
nautobot_manufacturer.new: Creating...
nautobot_manufacturer.new: Creation complete after 1s [id=eccd8b38-f6d6-41a4-aebd-73b53731b099]
data.nautobot_manufacturers.all: Reading...
data.nautobot_manufacturers.all: Read complete after 0s [id=1651744472]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

data_source_example = {
  "eccd8b38-f6d6-41a4-aebd-73b53731b099" = {
    "created" = "2022-05-05"
    "custom_fields" = tomap({})
    "description" = "Created with Terraform"
    "devicetype_count" = 0
    "display" = "New Vendor"
    "id" = "eccd8b38-f6d6-41a4-aebd-73b53731b099"
    "inventoryitem_count" = 0
    "last_updated" = "2022-05-05T09:54:32.661009Z"
    "name" = "New Vendor"
    "platform_count" = 0
    "slug" = "new-vendor"
    "url" = "https://demo.nautobot.com/api/dcim/manufacturers/eccd8b38-f6d6-41a4-aebd-73b53731b099/"
  }
}
```

### Explore

:-)

- https://networkop.co.uk/post/2019-04-tf-yang/
- https://github.com/networkop/terraform-yang