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

output "complete_lambda_authorizer_function_name" {
  value = aws_lambda_function.complete.function_name
}


output "complete_lambda_function_arn" {
  value = aws_lambda_function.complete.arn
}

output "apple_identity_pool_id" {
  value = aws_cognito_identity_pool.main.id
}

output "apigw_host" {
  value = local.rs_mesh_apigw_host
}

output "auth_exp_api_invoke_url" {
  value = aws_apigatewayv2_stage.default.invoke_url
}
