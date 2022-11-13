data "archive_file" "core" {
  type        = "zip"
  source_dir  = "../core"
  excludes    = ["../core/bin/**"]
  output_path = "bin/core.zip"
}


resource "aws_ecr_repository" "core" {
  name = "${local.app_stack}-core"
}



