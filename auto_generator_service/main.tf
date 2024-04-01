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

locals {
  service_name = "otel-package-updater"

  trigger_lambda_name = "${local.service_name}-trigger"
  dist_trigger_lambda = "dist/handler.zip"
}

resource "aws_cloudwatch_log_group" "trigger_lambda" {
  name              = "/aws/lambda/${local.trigger_lambda_name}"
  retention_in_days = 14
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
