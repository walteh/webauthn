output "init_lambda_function_name" {
	value = aws_lambda_function.init.function_name
}


output "init_lambda_function_arn" {
	value = aws_lambda_function.init.arn
}

output "apigw_lambda_authorizer_function_name" {
	value = aws_lambda_function.apigw_authorizer.function_name
}

output "apigw_lambda_authorizer_function_arn" {
	value = aws_lambda_function.apigw_authorizer.arn
}

output "appsync_lambda_authorizer_function_name" {
	value = aws_lambda_function.appsync_authorizer.function_name
}

output "appsync_lambda_authorizer_function_arn" {
	value = aws_lambda_function.appsync_authorizer.arn
}


output "ios_register_passkey_lambda_function_name" {
	value = aws_lambda_function.ios_register_passkey.function_name
}


output "ios_register_passkey_lambda_function_arn" {
	value = aws_lambda_function.ios_register_passkey.arn
}

output "passkey_assert_lambda_function_name" {
	value = aws_lambda_function.passkey_assert.function_name
}

output "passkey_assert_lambda_function_arn" {
	value = aws_lambda_function.passkey_assert.arn
}

output "ios_register_device_lambda_function_name" {
	value = aws_lambda_function.ios_register_device.function_name
}


output "ios_register_device_lambda_function_arn" {
	value = aws_lambda_function.ios_register_device.arn
}

output "devicecheck_assert_lambda_function_name" {
	value = aws_lambda_function.devicecheck_assert.function_name
}

output "devicecheck_assert_lambda_function_arn" {
	value = aws_lambda_function.devicecheck_assert.arn
}

output "apple_identity_pool_id" {
	value = aws_cognito_identity_pool.main.id
}

output "apple_identity_pool_arn" {
	value = aws_cognito_identity_pool.main.arn
}

output "apigw_host" {
	value = local.rs_mesh_apigw_host
}

output "auth_exp_api_invoke_url" {
	value = aws_apigatewayv2_stage.default.invoke_url
}
