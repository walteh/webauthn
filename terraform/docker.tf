data "archive_file" "apple" {
  type        = "zip"
  source_dir  = "../apple"
  excludes    = ["../apple/bin/**"]
  output_path = "bin/apple.zip"
}

/* resource "null_resource" "docker" {
  triggers = { src_hash = "${data.archive_file.apple.output_sha}" }
  provisioner "local-exec" {
    command = <<EOF
           aws ecr get-login-password --region ${local.aws_region} | docker login --username AWS --password-stdin ${local.aws_account}.dkr.ecr.${local.aws_region}.amazonaws.com
           cd ${path.module}/../apple
		   docker build --platform=linux/arm64 --target builder .
       EOF
  }
} */
