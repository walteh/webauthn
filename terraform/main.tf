

locals {
  app                 = "auth"
  env                 = terraform.workspace
  app_stack           = "${local.env}-${local.app}"
  external_app_domain = "${local.app}.${local.env}.api.${local.rs_mesh_route53_zone}"
  internal_app_domain = "${local.app}.${local.rs_mesh_namespace}"

  apple_team_id      = "4497QJSAD3"
  apple_key_id       = "7K626D4KLV"
  apple_service_name = "xyz.nugg.app"
}

locals {
  new_relic_export_host = local.rs_newrelic_region == "EU" ? "otlp.eu01.nr-data.net" : "otlp.nr-data.net"
  otel_prefix           = "${local.env}:${local.app}"
}

locals {
  latest = "latest"

  challenge_tag = "${local.latest}-challenge"
  apigw_tag     = "${local.latest}-apigw"
  appsync_tag   = "${local.latest}-appsync"

  primary = "primary"

  lambda_docker_deploy_command = <<EOF
	    	aws ecr get-login-password --region ${local.aws_region} | docker login --username AWS --password-stdin ${local.aws_account}.dkr.ecr.${local.aws_region}.amazonaws.com
			cd ../source
		    docker build --platform=linux/arm64 -t $tag -f Dockerfile.lambda --build-arg CMD=$cmd .
			docker push $tag
			docker image rm $tag
		EOF
}




data "aws_caller_identity" "current" {}
data "aws_availability_zones" "available_zones" { state = "available" }
data "aws_region" "current" {}
data "aws_partition" "current" {}


locals {
  aws_partition = data.aws_partition.current.partition
  aws_region    = data.aws_region.current.name
  aws_account   = data.aws_caller_identity.current.account_id
}
