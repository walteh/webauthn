locals {
  aws_partition = data.aws_partition.current.partition
  aws_region    = data.aws_region.current.name
  aws_account   = data.aws_caller_identity.current.account_id
}

locals {
  app       = "webauthn"
  env       = terraform.workspace
  app_stack = "${local.env}-${local.app}"

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

  appsync_dir              = "lambda/appsync/authorizer"
  apigw_dir                = "lambda/apigw/authorizer"
  init_dir                 = "lambda/init"
  ios_register_passkey_dir = "lambda/passkey/attest"
  passkey_assert_dir       = "lambda/passkey/assert"
  ios_register_device_dir  = "lambda/ios/register/device"

  devicecheck_assert_dir = "lambda/devicecheck/assert"

  apigw_tag                = replace("${local.apigw_dir}/${local.latest}", "/", "_")
  appsync_tag              = replace("${local.appsync_dir}/${local.latest}", "/", "_")
  init_tag                 = replace("${local.init_dir}/${local.latest}", "/", "_")
  ios_register_passkey_tag = replace("${local.ios_register_passkey_dir}/${local.latest}", "/", "_")
  passkey_assert_tag       = replace("${local.passkey_assert_dir}/${local.latest}", "/", "_")
  ios_register_device_tag  = replace("${local.ios_register_device_dir}/${local.latest}", "/", "_")
  devicecheck_assert_tag   = replace("${local.devicecheck_assert_dir}/${local.latest}", "/", "_")

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
