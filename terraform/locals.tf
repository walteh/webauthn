locals {
  aws_partition = data.aws_partition.current.partition
  aws_region    = data.aws_region.current.name
  aws_account   = data.aws_caller_identity.current.account_id
}

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

  appsync_dir                = "lambda/appsync-authorizer"
  apigw_dir                  = "lambda/apigw-authorizer"
  complete_dir               = "lambda/complete"
  apple_passkey_register_dir = "lambda/apple/passkey/register"
  apple_passkey_init_dir     = "lambda/apple/passkey/init"
  apple_passkey_login_dir    = "lambda/apple/passkey/login"

  complete_tag               = replace("${local.complete_dir}/${local.latest}", "/", "_")
  apigw_tag                  = replace("${local.apigw_dir}/${local.latest}", "/", "_")
  appsync_tag                = replace("${local.appsync_dir}/${local.latest}", "/", "_")
  apple_passkey_register_tag = replace("${local.apple_passkey_register_dir}/${local.latest}", "/", "_")
  apple_passkey_init_tag     = replace("${local.apple_passkey_init_dir}/${local.latest}", "/", "_")
  apple_passkey_login_tag    = replace("${local.apple_passkey_login_dir}/${local.latest}", "/", "_")

  primary = "primary"

  lambda_docker_deploy_command = <<EOF
	    	aws ecr get-login-password --region ${local.aws_region} | docker login --username AWS --password-stdin ${local.aws_account}.dkr.ecr.${local.aws_region}.amazonaws.com
			cd ../core
		    docker build --platform=linux/arm64 -t $tag -f Dockerfile.lambda --build-arg CMD=$dir/main.go .
			docker push $tag
			docker image rm $tag
		EOF
}

data "aws_iam_policy_document" "lambda_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}
