package inventory

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccItemDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + fmt.Sprintf(`
resource "inventory_item" "test" {
	name = "car"
	tag = "mustang"
}

data "inventory_item" "test" {
	id = inventory_item.test.id
}
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the item to ensure all attributes are set
					resource.TestCheckResourceAttr("data.inventory_item.test", "name", "car"),
					resource.TestCheckResourceAttr("data.inventory_item.test", "tag", "mustang"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttrSet("data.inventory_item.test", "id"),
				),
			},
		},
	})
}
