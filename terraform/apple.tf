
resource "aws_ecr_repository" "apple" {
  name = "${local.app_stack}-apple-image"
}

resource "aws_ecr_repository_policy" "apple" {
  repository = aws_ecr_repository.apple.name

  policy = jsonencode({
    "Version" : "2008-10-17",
    "Statement" : [
      {
        "Sid" : "ReadOnlyPermissions",
        "Effect" : "Allow",
        "Principal" : "*",
        "Action" : [
          "ecr:BatchCheckLayerAvailability",
          "ecr:BatchGetImage",
          "ecr:DescribeImageScanFindings",
          "ecr:DescribeImages",
          "ecr:DescribeRepositories",
          "ecr:GetAuthorizationToken",
          "ecr:GetDownloadUrlForLayer",
          "ecr:GetLifecyclePolicy",
          "ecr:GetLifecyclePolicyPreview",
          "ecr:GetRepositoryPolicy",
          "ecr:ListImages",
          "ecr:ListTagsForResource"
        ]
      }
    ]
  })
}

resource "aws_lambda_function" "apple" {
  function_name    = "${local.app_stack}-apple"
  image_uri        = "${aws_ecr_repository.apple.repository_url}:${local.latest}"
  role             = aws_iam_role.apple.arn
  memory_size      = 128
  timeout          = 120
  package_type     = "Image"
  publish          = true
  architectures    = ["arm64"]
  source_code_hash = trimprefix(data.aws_ecr_image.apple.image_digest, "sha256:")
  environment {
    variables = {
      CHALLENGE_TABLE_NAME          = "unused"
      APPLE_IDENTITY_POOL_ID        = aws_cognito_identity_pool.main.id
      APPLE_JWT_PUBLIC_KEY_ENDPOINT = "https://appleid.apple.com/auth/keys"
    }
  }
  tracing_config {
    mode = "Active"
  }
  depends_on = [
    aws_ecr_repository.apple,
    data.aws_ecr_image.apple
  ]
}

data "archive_file" "apple" {
  type        = "zip"
  source_dir  = "../apple"
  excludes    = ["../apple/bin/**"]
  output_path = "bin/apple.zip"
}

resource "null_resource" "apple" {
  triggers = {
    src_hash = "${data.archive_file.apple.output_sha}"
  }
  provisioner "local-exec" {
    command = <<EOF
           aws ecr get-login-password --region ${data.aws_region.current.name} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com
           cd ${path.module}/../apple
		   docker build --platform=linux/arm64 -t ${aws_ecr_repository.apple.repository_url}:${local.latest} .
           docker push ${aws_ecr_repository.apple.repository_url}:${local.latest}
       EOF
  }
}

data "aws_ecr_image" "apple" {
  depends_on = [
    null_resource.apple
  ]
  repository_name = aws_ecr_repository.apple.name
  image_tag       = local.latest
}

resource "aws_iam_role" "apple" {
  name               = "${local.app_stack}-apple-ExecutionRole"
  assume_role_policy = data.aws_iam_policy_document.apple_assume.json
  inline_policy {
    name   = "${local.app_stack}-apple-ExecutionRolePolicy"
    policy = data.aws_iam_policy_document.apple_inline.json
  }

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
    "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess",
  ]
}

data "aws_iam_policy_document" "apple_assume" {
  statement {
    effect = "Allow"
    actions = [
      "sts:AssumeRole"
    ]
    principals {
      type = "Service"
      identifiers = [
        "lambda.amazonaws.com"
      ]
    }
  }
}

data "aws_iam_policy_document" "apple_inline" {
  statement {
    effect = "Allow"
    actions = [
      "ecr:GetAuthorizationToken",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:GetRepositoryPolicy",
      "ecr:DescribeRepositories",
      "ecr:ListImages",
      "ecr:DescribeImages",
      "ecr:BatchGetImage"
    ]
    resources = [
      "*"
    ]
  }
}

