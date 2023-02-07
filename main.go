package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/superorbital/terraform-provider-inventory/inventory"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name inventory

func main() {
	providerserver.Serve(context.Background(), inventory.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/superorbital/inventory",
	})
}
