
module "lambda_function_existing_package_local" {
  
  source = "terraform-aws-modules/lambda/aws"

  function_name = "signUpTrigger"
  description   = "Create SNS topic, subscribe email, add to dynamodb"
  handler       = "trigger.handler"
  runtime       = "python3.10"
  timeout = 30 // cuz sns might take time

  create_package         = false
  local_existing_package = "${path.module}/lamda/signUpTrigger.zip"

  # Lambda needs permission to be invoked by Cognito
  attach_policy_statements = true
  policy_statements = {
    dynamodb = {
      effect = "Allow"
      actions = [ // adds username and email to Users table (depends on dynamo.tf)
        "dynamodb:PutItem",
        "dynamodb:GetItem",
        "dynamodb:UpdateItem"
      ]
      resources = [aws_dynamodb_table.users.arn] 
    }
    sns = {  // create sns topic for user and subscribes email 
      effect = "Allow"
      actions = [
        "sns:Subscribe",
        "sns:CreateTopic",
        "sns:Publish",
        "sns:SetTopicAttributes"

      ]
      resources = ["*"]
    }
  }
  
  environment_variables = {
    DYNAMODB_TABLE = aws_dynamodb_table.users.name
    AWS_REGION_NAME = var.region
  }
}





# cognito function permission to invoke lamda upon sign up
resource "aws_lambda_permission" "allow_cognito" {
  statement_id  = "AllowExecutionFromCognito"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda_function_existing_package_local.lambda_function_name
  principal     = "cognito-idp.amazonaws.com"
  source_arn    = aws_cognito_user_pool.user_pool.arn
}


