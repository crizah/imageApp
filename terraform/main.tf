// 2 dynamo table
// Messages
// messageID is pk
// recipient-index is index name for gsi with attr recipient

// Users
// username is pki

// 1 s3 bucket
// create bucket called encypted-files


// 2 lamda functions
module "lambda_function_existing_package_local" {
    // resource depends on users table
  source = "terraform-aws-modules/lambda/aws"

  function_name = "signUpTrigger"
  description   = "Create SNS topic, subscribe email, add to dynamodb"
  handler       = "index.lambda_handler"
  runtime       = "go"

  create_package         = false
  local_existing_package = "../existing_package.zip"
}


module "lambda_function_existing_package_local" {
    // resource depends on users table and api gateway
  source = "terraform-aws-modules/lambda/aws"

  function_name = "getUsers"
  description   = "gets list of all usernames in Users tables"
  handler       = "index.lambda_handler"
  runtime       = "go"

  create_package         = false
  local_existing_package = "../existing_package.zip"
}

// 1 api gateway