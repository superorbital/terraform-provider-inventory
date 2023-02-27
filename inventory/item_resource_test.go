package inventory

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccItemResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "inventory_item" "test" {
    name = "bottle"
    tag  = "rare"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inventory_item.test", "name", "bottle"),
					resource.TestCheckResourceAttr("inventory_item.test", "tag", "rare"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("inventory_item.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "inventory_item.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "inventory_item" "test" {
    name = "jet"
    tag  = "SR-71 Blackbird"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inventory_item.test", "name", "jet"),
					resource.TestCheckResourceAttr("inventory_item.test", "tag", "SR-71 Blackbird"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("inventory_item.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
