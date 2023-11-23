# Configure provider

terraform {
  required_providers {
    glesys = {
      source  = "glesys/glesys"
      version = "~> 0.10.1"
    }
  }
}

# Configure provider credentials
provider "glesys" {
  token  = "ABC123"
  userid = "CL12345"
}

# Create a server resource
resource "glesys_server" "www" {
  # ...
}

