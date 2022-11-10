data "terraform_remote_state" "mesh" {
  backend = "s3"
  config = {
    bucket = "nugg.xyz-terraform"
    key    = "mesh.tfstate"
    region = "us-east-1"
  }
  workspace = terraform.workspace
}

locals {
  rs_mesh_namespace                             = data.terraform_remote_state.mesh.outputs.namespace
  rs_mesh_namespace_id                          = data.terraform_remote_state.mesh.outputs.namespace_id
  rs_mesh_vpc                                   = data.terraform_remote_state.mesh.outputs.vpc
  rs_mesh_appmesh                               = data.terraform_remote_state.mesh.outputs.mesh
  rs_mesh_virtual_gateway                       = data.terraform_remote_state.mesh.outputs.virtual_gateway
  rs_mesh_virtual_gateway_endpoint              = data.terraform_remote_state.mesh.outputs.virtual_gateway_endpoint
  rs_mesh_private_subnet_a                      = data.terraform_remote_state.mesh.outputs.private_subnet_a
  rs_mesh_private_subnet_b                      = data.terraform_remote_state.mesh.outputs.private_subnet_b
  rs_mesh_public_subnet_a                       = data.terraform_remote_state.mesh.outputs.public_subnet_a
  rs_mesh_public_subnet_b                       = data.terraform_remote_state.mesh.outputs.public_subnet_b
  rs_mesh_ecs_cluster                           = data.terraform_remote_state.mesh.outputs.ecs_cluster
  rs_mesh_cloudwatch_log_group                  = data.terraform_remote_state.mesh.outputs.cloudwatch_log_group
  rs_mesh_http_security_group                   = data.terraform_remote_state.mesh.outputs.http_security_group
  rs_mesh_https_security_group                  = data.terraform_remote_state.mesh.outputs.https_security_group
  rs_mesh_egress_all_security_group             = data.terraform_remote_state.mesh.outputs.egress_all_security_group
  rs_mesh_route53_zone                          = data.terraform_remote_state.mesh.outputs.route53_zone
  rs_newrelic_otel_tracing_exporter_licence_key = data.terraform_remote_state.mesh.outputs.newrelic_otlp_tracing_exporter_licence_key
  rs_newrelic_region                            = data.terraform_remote_state.mesh.outputs.newrelic_region
  rs_newrelic_account_id                        = data.terraform_remote_state.mesh.outputs.newrelic_account_id
  rs_newrelic_api_key                           = data.terraform_remote_state.mesh.outputs.newrelic_api_key
  rs_mesh_apigw_host                            = data.terraform_remote_state.mesh.outputs.apigw_host
}
