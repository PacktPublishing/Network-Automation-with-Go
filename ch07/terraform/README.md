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
⇨  terraform init -upgrade

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
⇨  terraform apply --auto-approve

Terraform used the selected providers to generate the following execution plan. Resource actions
are indicated with the following symbols:
  + create

Terraform will perform the following actions:

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
nautobot_manufacturer.new: Creation complete after 1s [id=b5c5ada7-7f98-482e-916d-4ef5e8621d68]

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