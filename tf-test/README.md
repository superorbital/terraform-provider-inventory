# Hands-on Testing

## Getting Started

* modify `~/.terraformrc`

```hcl
provider_installation {
  dev_overrides {
    "superorbital/inventory" = "/Users/me/dev/go/path/bin/"
  }
  direct {}
}
```

* Build the binary, by running `make` in the root of the git repo.
* Spin up a copy of the inventory service. `docker container run --rm --name inventory-service -p 8080:8080 superorbital/inventory-service`
* Run `terraform init`
* Run whatever other terraform commands you want.

## Cleaning up

* Stop the inventory service. `docker container stop inventory-service`
* Comment out the `dev_overrides` in `~/.terraformrc`.
