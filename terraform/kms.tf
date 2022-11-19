/* resource "aws_kms_key" "a" {
  description              = "KMS key 1"
  deletion_window_in_days  = 10
  key_usage                = "SIGN_VERIFY"
  customer_master_key_spec = "RSA_2048"
}

data "aws_kms_public_key" "a" {
  depends_on = [aws_kms_key.a]
  key_id     = aws_kms_key.a.arn
} */
