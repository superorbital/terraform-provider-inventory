package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccItemDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `
resource "inventory_item" "test" {
  name = "2022 Mustang Shelby GT500"
  tag = "USD:79,420"
}

data "inventory_item" "test" {
	id = inventory_item.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the item to ensure all attributes are set
					resource.TestCheckResourceAttr("data.inventory_item.test", "name", "2022 Mustang Shelby GT500"),
					resource.TestCheckResourceAttr("data.inventory_item.test", "tag", "USD:79,420"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttrSet("data.inventory_item.test", "id"),
				),
			},
		},
	})
}
