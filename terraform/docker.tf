data "archive_file" "apple" {
  type        = "zip"
  source_dir  = "../apple"
  excludes    = ["../apple/bin/**"]
  output_path = "bin/apple.zip"
}

resource "null_resource" "docker" {
  triggers = { src_hash = "${data.archive_file.apple.output_sha}" }
  provisioner "local-exec" {
    command = <<EOF
           aws ecr get-login-password --region ${data.aws_region.current.name} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com
           cd ${path.module}/../apple

		   docker build --platform=linux/arm64 --target apple-apigw   -t ${aws_ecr_repository.apple_apigw.repository_url}:${local.latest} .
		   docker build --platform=linux/arm64 --target apple-appsync -t ${aws_ecr_repository.apple_appsync.repository_url}:${local.latest} .
		   docker build --platform=linux/arm64 --target challenge     -t ${aws_ecr_repository.challenge.repository_url}:${local.latest} .

           docker push ${aws_ecr_repository.apple_apigw.repository_url}:${local.latest}
           docker push ${aws_ecr_repository.apple_appsync.repository_url}:${local.latest}
		   docker push ${aws_ecr_repository.apple_appsync.repository_url}:${local.latest}
       EOF
  }
}
