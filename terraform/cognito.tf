resource "aws_cognito_identity_pool" "main" {
  identity_pool_name               = "${local.app_stack}-apple-identity-pool"
  allow_unauthenticated_identities = false
  allow_classic_flow               = false

  developer_provider_name = "nuggid.nugg.xyz"
}
