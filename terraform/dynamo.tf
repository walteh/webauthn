resource "aws_dynamodb_table" "credentials" {
	name         = "${local.app_stack}-credentials"
	hash_key     = "credential_id"
	billing_mode = "PAY_PER_REQUEST"

	attribute {
		name = "credential_id"
		type = "S"
	}

	attribute {
		name = "user_id"
		type = "S"
	}

	global_secondary_index {
		name            = "user_id-index"
		hash_key        = "user_id"
		projection_type = "ALL"
	}
}

resource "aws_dynamodb_table" "users" {
	name         = "${local.app_stack}-users"
	hash_key     = "user_id"
	billing_mode = "PAY_PER_REQUEST"

	attribute {
		name = "user_id"
		type = "S"
	}
}

resource "aws_dynamodb_table" "ceremonies" {
	name         = "${local.app_stack}-ceremonies"
	hash_key     = "challenge_id"
	billing_mode = "PAY_PER_REQUEST"

	attribute {
		name = "challenge_id"
		type = "S"
	}
}

