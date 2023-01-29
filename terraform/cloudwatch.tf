
resource "aws_cloudwatch_log_group" "stack" {
	name              = "/stack/${local.app_stack}"
	retention_in_days = 1
}


resource "aws_cloudwatch_log_group" "apigwv2_accesslogs" {
	name              = "/aws/apigwv2/access-logs/${local.app_stack}"
	retention_in_days = 1
}
