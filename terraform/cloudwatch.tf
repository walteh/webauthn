
resource "aws_cloudwatch_log_group" "stack" {
  name              = "/stack/${local.app_stack}"
  retention_in_days = 1
}
