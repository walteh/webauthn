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


output "apple_passkey_attest_lambda_function_name" {
  value = aws_lambda_function.apple_passkey_attest.function_name
}


output "apple_passkey_attest_lambda_function_arn" {
  value = aws_lambda_function.apple_passkey_attest.arn
}

output "apple_passkey_assert_lambda_function_name" {
  value = aws_lambda_function.apple_passkey_assert.function_name
}

output "apple_passkey_assert_lambda_function_arn" {
  value = aws_lambda_function.apple_passkey_assert.arn
}

output "apple_devicecheck_attest_lambda_function_name" {
  value = aws_lambda_function.apple_devicecheck_attest.function_name
}


output "apple_devicecheck_attest_lambda_function_arn" {
  value = aws_lambda_function.apple_devicecheck_attest.arn
}

output "apple_devicecheck_assert_lambda_function_name" {
  value = aws_lambda_function.apple_devicecheck_assert.function_name
}

output "apple_devicecheck_assert_lambda_function_arn" {
  value = aws_lambda_function.apple_devicecheck_assert.arn
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
