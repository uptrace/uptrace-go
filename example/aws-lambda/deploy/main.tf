provider "aws" {
  region = "us-east-1"

  assume_role {
    role_arn = "arn:aws:iam::754388400681:role/devs-prod"
  }
}

module "hello-lambda-function" {
  source  = "terraform-aws-modules/lambda/aws"
  version = ">= 2.24.0"

  architectures = compact([var.architecture])
  function_name = var.name
  handler       = "bootstrap"
  runtime       = "provided.al2"

  create_package         = false
  local_existing_package = "${path.module}/../function/build/bootstrap.zip"

  memory_size = 384
  timeout     = 20

#  layers = compact([
#    var.collector_layer_arn
#  ])

  tracing_mode = var.tracing_mode

  attach_policy_statements = true
  attach_tracing_policy = true
  policy_statements = {
    s3 = {
      effect = "Allow"
      actions = [
        "s3:ListAllMyBuckets"
      ]
      resources = [
        "*"
      ]
    }
  }
}

module "api-gateway" {
  source = "./api-gateway-proxy"

  name                = var.name
  function_name       = module.hello-lambda-function.lambda_function_name
  function_invoke_arn = module.hello-lambda-function.lambda_function_invoke_arn
  enable_xray_tracing = var.tracing_mode == "Active"
}
