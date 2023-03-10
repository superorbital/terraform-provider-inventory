---
page_title: "inventory_item Resource - inventory"
subcategory: ""
description: |-
  Manage an item.
---

# inventory_item (Resource)

Manage an item.

## Example Usage

```terraform
# Create a new inventory item
resource "inventory_item" "example" {
  name = "car"
  tag  = "mustang"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name for this inventory item.

### Optional

- `tag` (String) The tag for this inventory item.

### Read-Only

- `id` (Number) Identifier for this inventory item.

## Import

Import is supported using the following syntax:

```shell
terraform import inventory_item.example 1000
```
