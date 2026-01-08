// 2 dynamo table

resource "aws_dynamodb_table" "table1" {
    // Messages
// messageID is pk
// recipient-index is index name for gsi with attr recipient


  name           = var.t1
  billing_mode = "PAY_PER_REQUEST"
  hash_key       = var.pk1

  attribute {
    name = var.pk1
    type = "S"
  }
}

resource "aws_dynamodb_table" "table2" {
    // Users
// username is pki


  name           = var.t2
  billing_mode = "PAY_PER_REQUEST"
  hash_key       = var.pk2

  attribute {
    name = var.pk1
    type = "S"
  }
}



resource "aws_s3_bucket" "bucket" {
    // 1 s3 bucket
// create bucket called encypted-files
    bucket = var.bucketName
  
}





// 2 lamda functions
module "lambda_function_existing_package_local" {
    // resource depends on users table
  source = "terraform-aws-modules/lambda/aws"

  function_name = "signUpTrigger"
  description   = "Create SNS topic, subscribe email, add to dynamodb"
  handler       = "trigger.lamda"
  runtime       = "python 3.10"

  create_package         = false
  local_existing_package = "./signUpTrigger.zip"
}


module "lambda_function_existing_package_local" {
    // resource depends on users table and api gateway
  source = "terraform-aws-modules/lambda/aws"

  function_name = "getUsers"
  description   = "gets list of all usernames in Users tables"
  handler       = "getUsers.lamda"
  runtime       = "python 3.10"

  create_package         = false
  local_existing_package = "./getUsers.zip"
}

// 1 api gateway