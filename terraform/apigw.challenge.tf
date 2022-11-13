resource "aws_apigatewayv2_route" "challenge" {
  api_id             = aws_apigatewayv2_api.auth.id
  route_key          = "POST /challenge"
  authorization_type = "NONE"
  target             = "integrations/${aws_apigatewayv2_integration.challenge_lambda.id}"
}

resource "aws_lambda_permission" "challenge" {
  statement_id  = "${local.app_stack}-challenge-AllowExecutionFromApiGatewayChallengeRoute"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.challenge.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.auth.execution_arn}/*/*/challenge"
}



resource "aws_apigatewayv2_integration" "challenge_lambda" {
  api_id             = aws_apigatewayv2_api.auth.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.challenge.invoke_arn
  /* credentials_arn        = aws_iam_role.challenge_integration_role.arn */
  payload_format_version = "2.0"
}
