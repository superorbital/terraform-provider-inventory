config {
  format = "compact"
  plugin_dir = "~/.tflint.d/plugins"

  module = false
  force  = false
  disabled_by_default = false
}

plugin "aws" {
  enabled = true
  version = "0.17.0"
  source  = "github.com/terraform-linters/tflint-ruleset-aws"

  // see: https://github.com/terraform-linters/tflint-ruleset-aws/blob/master/docs/configuration.md
  deep_check = false
}

rule "terraform_deprecated_interpolation" {
  enabled = true
}

rule "terraform_deprecated_index" {
  enabled = true
}

rule "terraform_unused_declarations" {
  # I'd prefer to see this enabled, but we might need to do some cleanup for that.
  enabled = false
}

rule "terraform_comment_syntax" {
  enabled = true
}

rule "terraform_documented_outputs" {
  enabled = true
}

rule "terraform_documented_variables" {
  enabled = true
}

rule "terraform_typed_variables" {
  enabled = true
}

rule "terraform_module_pinned_source" {
  enabled = true
}

rule "terraform_required_version" {
  enabled = false
}

rule "terraform_required_providers" {
  enabled = false
}

rule "terraform_standard_module_structure" {
  enabled = false
}

rule "terraform_workspace_remote" {
  enabled = true
}

rule "terraform_naming_convention" {
  enabled = true

  variable {
    format = "snake_case"
  }

  locals {
    format = "snake_case"
  }

  output {
    format = "snake_case"
  }

  resource {
    format = "none"
  }

  module {
    format = "none"
  }

  data {
    format = "none"
  }
}
