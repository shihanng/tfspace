terraform {
  required_version = "~> 1.1.4"
  required_providers {
    local = {}
  }
  backend "local" {}
}

resource "local_file" "foo" {
  content  = var.env
  filename = "${var.env}/foo.bar"
}

variable "env" {
  type = string
}
