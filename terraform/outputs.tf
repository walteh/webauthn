output "apple_lambda_authorizer_function_name" {
  value = aws_lambda_function.apple.function_name
}

output "apple_lambda_authorizer_function_arn" {
  value = aws_lambda_function.apple.arn
}

output "apple_identity_pool_id" {
  value = aws_cognito_identity_pool.main.id
}
