terraform {
  required_providers {
    inventory = {
      source = "superorbital/inventory"
    }
  }
}

# Configure the connection details for the Inventory service
provider "inventory" {
  host = "127.0.0.1"
  port = "8080"
}

#Create new Inventory item
resource "inventory_item" "example" {
  name = "Jones Extreme Sour Cherry Warhead Soda"
  tag  = "USD:2.99"
}
