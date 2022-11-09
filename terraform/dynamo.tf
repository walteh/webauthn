resource "aws_dynamodb_table" "challenge" {
  name         = "${local.app_stack}-challenges"
  hash_key     = "id"
  billing_mode = "ON_DEMAND"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "ttl"
    type = "N"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

}

