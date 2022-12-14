
resource "aws_apigatewayv2_api" "auth" {
	name                       = "${local.app_stack}-exp-api"
	protocol_type              = "HTTP"
	route_selection_expression = "$request.method $request.path"
	cors_configuration {
		allow_origins  = ["*"]
		allow_methods  = ["POST"]
		allow_headers  = ["*"]
		expose_headers = ["*"]
	}
}

resource "aws_apigatewayv2_authorizer" "apple" {
	api_id                            = aws_apigatewayv2_api.auth.id
	authorizer_type                   = "REQUEST"
	authorizer_uri                    = aws_lambda_function.apigw_authorizer.invoke_arn
	identity_sources                  = ["$request.header.Authorization"]
	name                              = "${local.app_stack}-apple-apigw-authorizer"
	authorizer_payload_format_version = "2.0"
}

resource "aws_apigatewayv2_stage" "main" {
	api_id      = aws_apigatewayv2_api.auth.id
	name        = "main"
	auto_deploy = true
	access_log_settings {
		destination_arn = aws_cloudwatch_log_group.stack.arn
		format          = "$context.identity.sourceIp - $context.identity.caller - $context.identity.user [$context.requestTime] \"$context.httpMethod $context.routeKey $context.protocol\" $context.status $context.responseLength $context.requestId $context.integrationErrorMessage $context.error.message "
	}
}



resource "aws_apigatewayv2_api_mapping" "auth" {
	api_id          = aws_apigatewayv2_api.auth.id
	domain_name     = local.rs_mesh_apigw_host
	stage           = aws_apigatewayv2_stage.main.id
	api_mapping_key = "exp/${local.app}"
}

resource "aws_apigatewayv2_deployment" "main" {
	depends_on = [
		aws_apigatewayv2_route.init,
		aws_apigatewayv2_route.ios_register_device,
		aws_apigatewayv2_route.passkey_assert,
		aws_apigatewayv2_route.ios_register_passkey,
		aws_apigatewayv2_route.devicecheck_assert,

	]
	api_id = aws_apigatewayv2_api.auth.id
	lifecycle {
		create_before_destroy = true
	}
}
