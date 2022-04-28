# Terraform

```
terraform init -upgrade
terraform plan
terraform apply -auto-approve
```

## Nautobot Provider

```bash
terraform-provider-nautobot ⇨  make install
```

### Example

```bash
  terraform init -upgrade

Initializing the backend...

Initializing provider plugins...
- Finding github.com/nleiva/nautobot versions matching "0.1.0"...
- Installing github.com/nleiva/nautobot v0.1.0...
- Installed github.com/nleiva/nautobot v0.1.0 (unauthenticated)

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
⇨  terraform apply --auto-approve

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # nautobot_manufacturer.new will be created
  + resource "nautobot_manufacturer" "new" {
      + created             = (known after apply)
      + description         = "Created with Terraform"
      + devicetype_count    = (known after apply)
      + id                  = (known after apply)
      + inventoryitem_count = (known after apply)
      + last_updated        = (known after apply)
      + name                = "Vendor I"
      + platform_count      = (known after apply)
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Changes to Outputs:
  + juniper = {
      + "4873d752-5dbe-4006-8345-8279a0dfbbda" = {
          + created             = "2022-03-08"
          + custom_fields       = {}
          + description         = ""
          + devicetype_count    = 1
          + display             = "Juniper"
          + id                  = "4873d752-5dbe-4006-8345-8279a0dfbbda"
          + inventoryitem_count = 0
          + last_updated        = "2022-03-08T14:50:48.492203Z"
          + name                = "Juniper"
          + platform_count      = 1
          + slug                = "juniper"
          + url                 = "https://demo.nautobot.com/api/dcim/manufacturers/4873d752-5dbe-4006-8345-8279a0dfbbda/"
        }
    }
nautobot_manufacturer.new: Creating...
nautobot_manufacturer.new: Creation complete after 1s [id=9fec5cd1-d23b-40e3-abe6-a220476631af]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

juniper = {
  "4873d752-5dbe-4006-8345-8279a0dfbbda" = {
    "created" = "2022-03-08"
    "custom_fields" = tomap({})
    "description" = ""
    "devicetype_count" = 1
    "display" = "Juniper"
    "id" = "4873d752-5dbe-4006-8345-8279a0dfbbda"
    "inventoryitem_count" = 0
    "last_updated" = "2022-03-08T14:50:48.492203Z"
    "name" = "Juniper"
    "platform_count" = 1
    "slug" = "juniper"
    "url" = "https://demo.nautobot.com/api/dcim/manufacturers/4873d752-5dbe-4006-8345-8279a0dfbbda/"
  }
}

```

### Explore

:-)

- https://networkop.co.uk/post/2019-04-tf-yang/
- https://github.com/networkop/terraform-yang