
resource "aws_ecr_repository" "apple_auth" {
  name = "${local.app_stack}-apple-auth-image"
}

resource "aws_ecr_repository_policy" "apple_auth" {
  repository = aws_ecr_repository.apple_auth.name

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



resource "aws_lambda_function" "apple_auth" {
  function_name    = "${local.app_stack}-apple-auth"
  image_uri        = "${aws_ecr_repository.apple_auth.repository_url}:${local.latest}"
  role             = aws_iam_role.apple_auth.arn
  memory_size      = 128
  timeout          = 120
  package_type     = "Image"
  publish          = true
  architectures    = ["arm64"]
  source_code_hash = trimprefix(data.aws_ecr_image.apple_auth.image_digest, "sha256:")
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
  vpc_config {
    subnet_ids         = [local.rs_mesh_private_subnet_a, local.rs_mesh_private_subnet_b]
    security_group_ids = [local.rs_mesh_egress_all_security_group]
  }
  depends_on = [
    aws_ecr_repository.apple_auth,
    data.aws_ecr_image.apple_auth
  ]
}

data "archive_file" "apple_auth" {
  type        = "zip"
  source_dir  = "../go"
  excludes    = ["../go/bin/**"]
  output_path = "bin/apple_auth.zip"
}

resource "null_resource" "apple_auth" {
  triggers = {
    src_hash = "${data.archive_file.apple_auth.output_sha}"
  }
  provisioner "local-exec" {
    command = <<EOF
           aws ecr get-login-password --region ${data.aws_region.current.name} | docker login --username AWS --password-stdin ${data.aws_caller_identity.current.account_id}.dkr.ecr.${data.aws_region.current.name}.amazonaws.com
           cd ${path.module}/../go
		   docker build --platform=linux/arm64 -t ${aws_ecr_repository.apple_auth.repository_url}:${local.latest} .
           docker push ${aws_ecr_repository.apple_auth.repository_url}:${local.latest}
       EOF
  }
}

data "aws_ecr_image" "apple_auth" {
  depends_on = [
    null_resource.apple_auth
  ]
  repository_name = aws_ecr_repository.apple_auth.name
  image_tag       = local.latest
}

resource "aws_iam_role" "apple_auth" {
  name               = "${local.app_stack}-apple-auth-ExecutionRole"
  assume_role_policy = data.aws_iam_policy_document.apple_auth_assume.json
  inline_policy {
    name   = "${local.app_stack}-apple-auth-ExecutionRolePolicy"
    policy = data.aws_iam_policy_document.apple_auth_inline.json
  }

  managed_policy_arns = [
    "arn:aws:iam::aws:policy/service-role/AWSLambdaVPCAccessExecutionRole",
    "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole",
    "arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess",
  ]
}

data "aws_iam_policy_document" "apple_auth_assume" {
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

data "aws_iam_policy_document" "apple_auth_inline" {
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

/* resource "aws_lambda_permission" "apple_auth" {
  statement_id  = "${local.app_stack}-apple-auth-AllowExecutionFromAppSync"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.apple_auth.function_name
  principal     = "appsync.amazonaws.com"
  source_arn    = aws_appsync_graphql_api.appsync.arn
} */

