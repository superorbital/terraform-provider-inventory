package provider

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
    name = "Jones Extreme Sour Cherry Warhead Soda"
    tag = "USD:2.99"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inventory_item.test", "name", "Jones Extreme Sour Cherry Warhead Soda"),
					resource.TestCheckResourceAttr("inventory_item.test", "tag", "USD:2.99"),
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
    name = "1928 de Havilland DH-60GM"
    tag  = "USD:110,781"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("inventory_item.test", "name", "1928 de Havilland DH-60GM"),
					resource.TestCheckResourceAttr("inventory_item.test", "tag", "USD:110,781"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("inventory_item.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
