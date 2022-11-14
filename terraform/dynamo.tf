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

resource "aws_dynamodb_table" "user" {
  name         = "${local.app_stack}-users"
  hash_key     = "id"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "id"
    type = "S"
  }

  attribute {
    name = "apple_id"
    type = "S"
  }

  global_secondary_index {
    name            = "apple_id"
    hash_key        = "apple_id"
    projection_type = "ALL"
  }
}

/* resource "aws_dynamodb_table" "ceremony" {
  name         = "${local.app_stack}-webauthn-ceremony"
  hash_key     = "ceremony_user_id"
  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "ceremony_user_id"
    type = "S"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }
} */

/* resource "aws_dynamodb_table" "session" {
  name      = "${local.app_stack}-sessions"
  hash_key  = "user_id"
  range_key = "auth_type"

  billing_mode = "PAY_PER_REQUEST"

  attribute {
    name = "user_id"
    type = "S"
  }

  attribute {
    name = "auth_type"
    type = "S"
  }

  attribute {
    name = "ttl"
    type = "S"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }
} */
