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

Changes to Outputs:
  + all_manufacturers = {
      + id            = "1651065682"
      + manufacturers = [
          + {
              + created             = "2022-03-08"
              + custom_fields       = {}
              + description         = ""
              + devicetype_count    = 4
              + display             = "Arista"
              + id                  = "832a9d3f-3d7e-40ec-b665-c4f1f056ccfd"
              + inventoryitem_count = 0
              + last_updated        = "2022-03-08T14:50:48.477597Z"
              + name                = "Arista"
              + platform_count      = 1
              + slug                = "arista"
              + url                 = "https://demo.nautobot.com/api/dcim/manufacturers/832a9d3f-3d7e-40ec-b665-c4f1f056ccfd/"
            },
          + {
              + created             = "2022-03-08"
              + custom_fields       = {}
              + description         = ""
              + devicetype_count    = 4
              + display             = "Cisco"
              + id                  = "1943f761-d4a3-45d2-814f-f3623d613789"
              + inventoryitem_count = 38
              + last_updated        = "2022-03-08T14:50:48.483444Z"
              + name                = "Cisco"
              + platform_count      = 5
              + slug                = "cisco"
              + url                 = "https://demo.nautobot.com/api/dcim/manufacturers/1943f761-d4a3-45d2-814f-f3623d613789/"
            },
          + {
              + created             = "2022-03-08"
              + custom_fields       = {}
              + description         = ""
              + devicetype_count    = 0
              + display             = "HP"
              + id                  = "bde6e2f0-72ef-4fe4-bbb3-fddaa2ab1451"
              + inventoryitem_count = 0
              + last_updated        = "2022-03-08T14:50:48.498452Z"
              + name                = "HP"
              + platform_count      = 0
              + slug                = "hp"
              + url                 = "https://demo.nautobot.com/api/dcim/manufacturers/bde6e2f0-72ef-4fe4-bbb3-fddaa2ab1451/"
            },
          + {
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
            },
          + {
              + created             = "2022-03-08"
              + custom_fields       = {}
              + description         = ""
              + devicetype_count    = 0
              + display             = "Mellanox"
              + id                  = "bbdf3da2-9f9e-486d-9367-0d3f5bde7fe3"
              + inventoryitem_count = 0
              + last_updated        = "2022-03-08T14:50:48.504294Z"
              + name                = "Mellanox"
              + platform_count      = 0
              + slug                = "mellanox"
              + url                 = "https://demo.nautobot.com/api/dcim/manufacturers/bbdf3da2-9f9e-486d-9367-0d3f5bde7fe3/"
            },
          + {
              + created             = "2022-03-08"
              + custom_fields       = {}
              + description         = ""
              + devicetype_count    = 0
              + display             = "Meraki"
              + id                  = "53f51c84-aad8-4d47-b1ed-acd3fa7274e4"
              + inventoryitem_count = 0
              + last_updated        = "2022-03-08T14:50:48.511030Z"
              + name                = "Meraki"
              + platform_count      = 0
              + slug                = "meraki"
              + url                 = "https://demo.nautobot.com/api/dcim/manufacturers/53f51c84-aad8-4d47-b1ed-acd3fa7274e4/"
            },
          + {
              + created             = "2022-03-08"
              + custom_fields       = {}
              + description         = ""
              + devicetype_count    = 0
              + display             = "NVIDIA"
              + id                  = "8d792462-8c55-4f39-8873-efba96c5de69"
              + inventoryitem_count = 0
              + last_updated        = "2022-03-08T14:50:48.517165Z"
              + name                = "NVIDIA"
              + platform_count      = 0
              + slug                = "nvidia"
              + url                 = "https://demo.nautobot.com/api/dcim/manufacturers/8d792462-8c55-4f39-8873-efba96c5de69/"
            },
          + {
              + created             = "2022-03-08"
              + custom_fields       = {}
              + description         = ""
              + devicetype_count    = 0
              + display             = "Opengear"
              + id                  = "864d37fd-3795-4e05-97af-bbba7a46230c"
              + inventoryitem_count = 0
              + last_updated        = "2022-03-08T14:50:48.523089Z"
              + name                = "Opengear"
              + platform_count      = 0
              + slug                = "opengear"
              + url                 = "https://demo.nautobot.com/api/dcim/manufacturers/864d37fd-3795-4e05-97af-bbba7a46230c/"
            },
          + {
              + created             = "2022-03-08"
              + custom_fields       = {}
              + description         = ""
              + devicetype_count    = 0
              + display             = "Palo Alto"
              + id                  = "4b7fa07f-66b4-434a-94b2-abaf7dae82ea"
              + inventoryitem_count = 0
              + last_updated        = "2022-03-08T14:50:48.528907Z"
              + name                = "Palo Alto"
              + platform_count      = 0
              + slug                = "pan"
              + url                 = "https://demo.nautobot.com/api/dcim/manufacturers/4b7fa07f-66b4-434a-94b2-abaf7dae82ea/"
            },
        ]
    }

You can apply this plan to save these new output values to the Terraform state, without changing any real infrastructure.

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```

### Explore

:-)

- https://networkop.co.uk/post/2019-04-tf-yang/
- https://github.com/networkop/terraform-yang