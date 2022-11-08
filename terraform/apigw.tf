/* resource "aws_apigatewayv2_authorizer" "example" {
  api_id           = aws_apigatewayv2_api.example.id
  authorizer_type  = "REQUEST"
  authorizer_uri   = aws_lambda_function.example.invoke_arn
  identity_sources = ["route.request.header.Auth"]
  name             = "example-authorizer"
} */

/* resource "aws_apigatewayv2_api" "example" {
  name                       = "example-websocket-api"
  protocol_type              = "WEBSOCKET"
  route_selection_expression = "$request.body.action"
}

resource "aws_apigatewayv2_api_mapping" "example" {
  api_id      = aws_apigatewayv2_api.example.id
  domain_name = aws_apigatewayv2_domain_name.example.id
  stage       = aws_apigatewayv2_stage.example.id
}



resource "aws_apigatewayv2_domain_name" "example" {
  domain_name = "http-api.example.com"

  domain_name_configuration {
    certificate_arn = aws_acm_certificate.example.arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_route53_record" "example" {
  name    = aws_apigatewayv2_domain_name.example.domain_name
  type    = "A"
  zone_id = aws_route53_zone.example.zone_id

  alias {
    name                   = aws_apigatewayv2_domain_name.example.domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.example.domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_apigatewayv2_integration" "example" {
  api_id           = aws_apigatewayv2_api.example.id
  integration_type = "AWS_PROXY"

  connection_type           = "INTERNET"
  content_handling_strategy = "CONVERT_TO_TEXT"
  description               = "Lambda example"
  integration_method        = "POST"
  integration_uri           = aws_lambda_function.example.invoke_arn
  passthrough_behavior      = "WHEN_NO_MATCH"
}

resource "aws_apigatewayv2_integration_response" "example" {
  api_id                   = aws_apigatewayv2_api.example.id
  integration_id           = aws_apigatewayv2_integration.example.id
  integration_response_key = "/200/"
} */
