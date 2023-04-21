resource "aws_dynamodb_table" "credential" {
	name         = "${local.app_stack}-credential"
	hash_key     = "pk"
	billing_mode = "PAY_PER_REQUEST"

	attribute {
		name = "pk"
		type = "S"
	}

	attribute {
		name = "sk"
		type = "S"
	}

	global_secondary_index {
		name            = "user_id-index"
		hash_key        = "user_id"
		projection_type = "ALL"
	}
}

/* resource "aws_dynamodb_table" "users" {
	name         = "${local.app_stack}-users"
	hash_key     = "user_id"
	billing_mode = "PAY_PER_REQUEST"

	attribute {
		name = "user_id"
		type = "S"
	}
} */

resource "aws_dynamodb_table" "challenge" {
	name         = "${local.app_stack}-challenge"
	hash_key     = "pk"
	billing_mode = "PAY_PER_REQUEST"

	attribute {
		name = "pk"
		type = "S"
	}
}

