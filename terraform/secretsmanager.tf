resource "aws_secretsmanager_secret" "apple_signinwithapple_privatekey" {
	name_prefix = "${local.app_stack}-apple-signinwithapple-privatekey"
}
