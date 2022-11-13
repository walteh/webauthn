resource "aws_apigatewayv2_route" "register" {
  api_id             = aws_apigatewayv2_api.auth.id
  route_key          = "POST /register"
  authorization_type = "NONE"
  target             = "integrations/${aws_apigatewayv2_integration.register_lambda.id}"
}

resource "aws_lambda_permission" "register" {
  statement_id  = "${local.app_stack}-register-AllowExecutionFromApiGatewayregisterRoute"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.register.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.auth.execution_arn}/*/*/register"
}


resource "aws_apigatewayv2_integration" "register_lambda" {
  api_id             = aws_apigatewayv2_api.auth.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.register.invoke_arn
  /* credentials_arn        = aws_iam_role.register_integration_role.arn */
  payload_format_version = "2.0"
}
