terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.43.0"
    }
  }

  backend "s3" {
    bucket = "infrastructure-terraform-bucket-ywnlj"
    key    = "otel-package-updater"
    region = "eu-west-1"
  }
}

provider "aws" {
  region = "eu-west-1"

  default_tags {
    tags = {
      Environment = "dev"
      Owner       = "Jakub Jantosik"
      Project     = local.service_name
    }
  }
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

data "aws_ssm_parameter" "deploy_key" {
  name = "/otel-grpc-healthcheck/private_deploy_key"
}

locals {
  service_name = "otel-package-updater"

  account_id = data.aws_caller_identity.current.account_id
  region     = data.aws_region.current.name

  log_retention_days = 7

  ecr_repository_url = "${local.account_id}.dkr.ecr.${local.region}.amazonaws.com"
  build_tag          = "latest"

  trigger_lambda_name = "${local.service_name}-trigger"
  dist_trigger_lambda = "dist/handler.zip"

  build_project_name = "${local.service_name}-build"
}

resource "aws_cloudwatch_log_group" "trigger_lambda" {
  name              = "/aws/lambda/${local.trigger_lambda_name}"
  retention_in_days = local.log_retention_days
}

data "aws_iam_policy_document" "lambda_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }
  }
}

data "aws_iam_policy_document" "trigger_policy" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = ["arn:aws:logs:*:*:*"]
  }
}

resource "aws_iam_role" "trigger_lambda" {
  name               = "${local.trigger_lambda_name}-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.lambda_assume_role.json

  inline_policy {
    name   = "default"
    policy = data.aws_iam_policy_document.trigger_policy.json
  }
}

resource "aws_lambda_function" "trigger" {
  function_name = local.trigger_lambda_name
  role          = aws_iam_role.trigger_lambda.arn
  handler       = "handler.default"

  timeout = 30

  source_code_hash = filebase64sha256(local.dist_trigger_lambda)

  filename = local.dist_trigger_lambda
  runtime  = "nodejs20.x"
}

resource "aws_ecr_repository" "build" {
  name                 = "${local.service_name}-build"
  image_tag_mutability = "MUTABLE"
  force_delete         = true

  image_scanning_configuration {
    scan_on_push = false
  }
}

resource "aws_ecr_lifecycle_policy" "build" {
  repository = aws_ecr_repository.build.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Expire untagged images older than 3 days"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 3
        }
        action = {
          type = "expire"
        }
      },
    ]
  })
}

resource "terraform_data" "build" {
  triggers_replace = [filesha256("${path.root}/Dockerfile")]

  provisioner "local-exec" {
    working_dir = path.root
    interpreter = ["/bin/bash", "-ce"]
    command     = <<-EOT
      aws ecr get-login-password --region ${local.region} | docker login --username AWS --password-stdin ${local.ecr_repository_url}
      docker build \
	--file Dockerfile \
	--platform linux/amd64 \
	--output type=image,name=${aws_ecr_repository.build.repository_url}:${local.build_tag},push=true .
    EOT
  }
}

data "aws_ecr_image" "build" {
  depends_on = [terraform_data.build]

  repository_name = aws_ecr_repository.build.name
  image_tag       = local.build_tag
}

resource "aws_cloudwatch_log_group" "build" {
  name              = "/codebuild/${local.build_project_name}"
  retention_in_days = local.log_retention_days
}

data "aws_iam_policy_document" "build_assume_role" {
  statement {
    principals {
      type        = "Service"
      identifiers = ["codebuild.amazonaws.com"]
    }
    actions = ["sts:AssumeRole"]
  }
}

data "aws_iam_policy_document" "build_policy" {
  statement {
    actions = [
      "logs:CreateLogStream",
      "logs:PutLogEvents",
    ]
    resources = [
      aws_cloudwatch_log_group.build.arn,
      "${aws_cloudwatch_log_group.build.arn}:log-stream:*",
    ]
  }

  statement {
    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetDownloadUrlForLayer",
      "ecr:BatchGetImage",
      "ecr:GetAuthorizationToken",
    ]
    resources = ["*"]
  }

  statement {
    actions   = ["ssm:GetParameter*"]
    resources = [data.aws_ssm_parameter.deploy_key.arn]
  }
}

resource "aws_iam_role" "build" {
  name               = "${local.build_project_name}-build-role"
  assume_role_policy = data.aws_iam_policy_document.build_assume_role.json

  inline_policy {
    name   = "default"
    policy = data.aws_iam_policy_document.build_policy.json
  }
}

resource "aws_codebuild_project" "build" {
  name          = local.build_project_name
  build_timeout = 30
  service_role  = aws_iam_role.build.arn

  artifacts {
    type = "NO_ARTIFACTS"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_SMALL"
    image                       = "${aws_ecr_repository.build.repository_url}@${data.aws_ecr_image.build.id}"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "SERVICE_ROLE"
  }

  logs_config {
    cloudwatch_logs {
      group_name = aws_cloudwatch_log_group.build.name
    }
  }

  source {
    type      = "NO_SOURCE"
    buildspec = file("${path.root}/buildspec.yaml")
  }
}
