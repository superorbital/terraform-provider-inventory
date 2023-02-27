# Inventory service provider for Terraform

A modern (protocol 6) Terraform provider for the demo inventory service (work-in-progress).

## `terraform-provider-inventory`

This Terraform provider code for the inventory service is a proof of concept used primarily for teaching.

### Build provider

Run the following command to build & install the provider

```shell
make
```

### Test provider

```shell
make test
```

```shell
make testacc
```

- If you want to test out the provider with the `terraform` CLI.
  - Edit `$HOME/.terraformrc` and point "superorbital/inventory" to your ${GOBIN} directory.

```hcl
provider_installation {
  dev_overrides {
    "superorbital/inventory" = "/home/me/go/path/bin/"
  }
  direct {}
}
```

And then create and test a few runs based on the files under examples. **NOTE:** You should not run `terraform init` when using _dev\_overrides_.

### Documentation

Documentation is generated with [tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs) and exists in the [docs](./docs/) directory.

```shell
make generate
```

### Release

- Make sure that you update the auto-generated documentation and then commit any changes!

```shell
make generate
```

Then go ahead and create a new release in Github. This will kick of the Github action to handle the rest. Release tags should use semantic versioning and look something like this `v1.0.2`.

## Pre-Commit Hooks

- See: [pre-commit](https://pre-commit.com/)
  - [pre-commit/pre-commit-hooks](https://github.com/pre-commit/pre-commit-hooks)
  - [antonbabenko/pre-commit-terraform](https://github.com/antonbabenko/pre-commit-terraform)

### Install

#### Local Install (macOS)

- **IMPORTANT**: All developers committing any code to this repo, should have these pre-commit hooks installed locally. Github actions may also run these at some point, but it is generally faster and easier to run them locally, in most cases.

```sh
brew install pre-commit jq shellcheck shfmt git-secrets go-critic golangci-lint
go install github.com/BurntSushi/toml/cmd/tomlv@master
go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
go install golang.org/x/tools/cmd/goimports@latest

mkdir -p ${HOME}/.git-template/hooks
git config --global init.templateDir ${HOME}/.git-template
```

- Close and reopen your terminal
- Make sure that you run these commands from the root of this git repo!

```sh
cd terraform-provider-inventory
pre-commit init-templatedir -t pre-commit ${HOME}/.git-template
pre-commit install
```

- Test it

```sh
pre-commit run -a
git diff
```

### Checks

See:

- [.pre-commit-config.yaml](./.pre-commit-config.yaml)

#### Configuring Hooks

- [.pre-commit-config.yaml](./.pre-commit-config.yaml)
