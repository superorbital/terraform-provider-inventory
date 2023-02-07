terraform {
  required_providers {
    inventory = {
      source = "superorbital/inventory"
    }
  }
}

# Configure the connection details for the inventory service
provider "inventory" {
  host = "127.0.0.1"
  port = "8080"
}

# Read in a existing inventory item
data "inventory_item" "example" {
 id = "1000"
}

# Create a new inventory item
resource "inventory_item" "example" {
 name = "car"
 tag  = "mustang"
}
