resource "aws_dynamodb_table" "challenge" {
  name         = "${local.app_stack}-challenges"
  hash_key     = "id"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "id"
    type = "S"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

}

