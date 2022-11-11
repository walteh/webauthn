output "apple_apigw_lambda_authorizer_function_name" {
  value = aws_lambda_function.apple_apigw.function_name
}

output "apple_apigw_lambda_authorizer_function_arn" {
  value = aws_lambda_function.apple_apigw.arn
}

output "apple_appsync_lambda_authorizer_function_name" {
  value = aws_lambda_function.apple_appsync.function_name
}

output "apple_appsync_lambda_authorizer_function_arn" {
  value = aws_lambda_function.apple_appsync.arn
}

output "challenge_lambda_authorizer_function_name" {
  value = aws_lambda_function.challenge.function_name
}

output "challenge_lambda_authorizer_function_arn" {
  value = aws_lambda_function.challenge.arn
}

output "apple_identity_pool_id" {
  value = aws_cognito_identity_pool.main.id
}
