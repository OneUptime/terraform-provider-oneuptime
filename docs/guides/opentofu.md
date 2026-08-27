---
page_title: "Using the Oneuptime provider with OpenTofu"
subcategory: "Guides"
description: |-
  The Oneuptime provider is published to the OpenTofu Registry and tested against OpenTofu on every change.
---

# Using the Oneuptime provider with OpenTofu

This provider works with [OpenTofu](https://opentofu.org) as well as Terraform, and is published to both registries. OpenTofu is a tested path, not an assumption inherited from Terraform compatibility: the provider's end-to-end suite runs the same fixtures against both engines on every pull request, and a break under `tofu` fails the build.

## Usage

Nothing in your configuration changes — declare the provider as usual and drive it with `tofu`:

```terraform
terraform {
  required_providers {
    oneuptime = {
      source  = "oneuptime/oneuptime"
      version = "~> 12.0"
    }
  }
}

provider "oneuptime" {
  # api_key is read from ONEUPTIME_API_KEY.
}
```

```shell
export ONEUPTIME_API_KEY="your-project-api-key"
tofu init
tofu plan
tofu apply
```

## Keep the registry hostname out of the source address

`source = "oneuptime/oneuptime"` carries no hostname, so each engine resolves it against its own default registry — `registry.opentofu.org` under OpenTofu, `registry.terraform.io` under Terraform. The provider is published to both, so one address covers both engines.

Writing `source = "registry.terraform.io/oneuptime/oneuptime"` pins the configuration to the Terraform Registry and breaks it under OpenTofu in registry-restricted environments.

## Differences worth knowing

| Topic | Behaviour |
|-------|-----------|
| The `terraform` block | Stays `terraform`. It is a language keyword, not a reference to the Terraform CLI. |
| `.tf` vs `.tofu` files | `.tf` works under both. OpenTofu also reads `.tofu` files and ignores any `.tf` file with a `.tofu` sibling. |
| `required_version` | OpenTofu's version series starts at 1.6.0, so a `>= 1.5.0` constraint written for Terraform is always satisfied. |
| Lock files | Both write `.terraform.lock.hcl`, but it records the registry it resolved against — one engine's lock file does not satisfy the other. |
| State files | Identical format and filenames; state moves between engines without conversion. |
| Variables | OpenTofu reads `TF_VAR_*` as well as `TOFU_VAR_*`. |

## Reusable modules

This repository ships hand-written modules under [`modules/`](https://github.com/OneUptime/terraform-provider-oneuptime/tree/master/modules). They are engine-agnostic HCL and work under Terraform too.

```terraform
module "storefront" {
  source = "github.com/OneUptime/terraform-provider-oneuptime//modules/monitoring-and-incident-response?ref=v12.0.24"

  service_name          = "storefront"
  status_page_is_public = true

  monitors = {
    homepage = { url = "https://example.com" }
    api      = { url = "https://api.example.com/health", expected_status_code = "204" }
  }
}
```

Pin `ref` to a published tag: this repository is regenerated on every release, so an unpinned source tracks whatever was published last.

## More

- OpenTofu guide: [oneuptime.com/docs/terraform/opentofu](https://oneuptime.com/docs/terraform/opentofu)
- Runnable examples: [`Examples/opentofu/`](https://github.com/OneUptime/oneuptime/tree/master/Examples/opentofu) in the main OneUptime repository
- Issues, including documentation issues: [github.com/OneUptime/oneuptime/issues](https://github.com/OneUptime/oneuptime/issues)
