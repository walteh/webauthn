resource "aws_apigatewayv2_route" "complete" {
  api_id             = aws_apigatewayv2_api.auth.id
  route_key          = "POST /complete"
  authorization_type = "NONE"
  target             = "integrations/${aws_apigatewayv2_integration.complete_lambda.id}"
}

resource "aws_lambda_permission" "complete" {
  statement_id  = "${local.app_stack}-complete-AllowExecutionFromApiGatewayChallengeRoute"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.complete.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.auth.execution_arn}/*/*/complete"
}



resource "aws_apigatewayv2_integration" "complete_lambda" {
  api_id             = aws_apigatewayv2_api.auth.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.complete.invoke_arn
  /* credentials_arn        = aws_iam_role.complete_integration_role.arn */
  payload_format_version = "2.0"
}