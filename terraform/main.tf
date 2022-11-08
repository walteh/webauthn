

locals {
  app                 = "auth"
  env                 = terraform.workspace
  app_stack           = "${local.env}-${local.app}"
  external_app_domain = "${local.app}.${local.env}.api.${local.rs_mesh_route53_zone}"
  internal_app_domain = "${local.app}.${local.rs_mesh_namespace}"
}

locals {
  new_relic_export_host = local.rs_newrelic_region == "EU" ? "otlp.eu01.nr-data.net" : "otlp.nr-data.net"
  otel_prefix           = "${local.env}:${local.app}"
}

locals {
  latest  = "latest"
  primary = "primary"
}




data "aws_caller_identity" "current" {}

data "aws_availability_zones" "available_zones" {
  state = "available"
}

data "aws_region" "current" {
}

data "aws_partition" "current" {}
