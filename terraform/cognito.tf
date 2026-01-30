resource "aws_cognito_user_pool" "user_pool" {
  name = "user-pool"

   alias_attributes = ["email", "preferred_username"] // allows foir login with username as well as email
  auto_verified_attributes = ["email"]

  password_policy {
    minimum_length = 6
  }

  verification_message_template {
    default_email_option = "CONFIRM_WITH_CODE" # confirm with link needs a user domain to be associated. stupid 
    email_subject = "Account Confirmation"
    email_message = "This is your verification code: {####}"
  }


  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true
    name                     = "email"
    required                 = true

    string_attribute_constraints {
      min_length = 1
      max_length = 256
    }


    
  }

  schema {
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true  
    name                     = "name"
    required                 = true

  }


  # trigger a lamda function upon sign up and after user is verified  
    lambda_config {
    post_confirmation = module.lambda_function_existing_package_local.lambda_function_arn
    }

    # recovery
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }
    
}

resource "aws_cognito_user_pool_client" "client" {
  name = "cognito-client"

  user_pool_id = aws_cognito_user_pool.user_pool.id
  generate_secret = false
  refresh_token_validity = 90
  prevent_user_existence_errors = "ENABLED"
  explicit_auth_flows = [
    "ALLOW_REFRESH_TOKEN_AUTH",
    "ALLOW_USER_PASSWORD_AUTH",
    "ALLOW_ADMIN_USER_PASSWORD_AUTH"
  ]


  # validity
  access_token_validity  = 60   
  id_token_validity      = 60  
  token_validity_units {
    access_token  = "minutes"
    id_token      = "minutes"
    refresh_token = "days"
  }


  # read/write attributes
  read_attributes = [
    "email",
    "email_verified",
    "name"
  ]
  
  write_attributes = [
    "email",
    "name"
  ]

}


output "user_pool_id" {
  value = aws_cognito_user_pool.user_pool.id
}

output "user_pool_client_id" {
  value = aws_cognito_user_pool_client.client.id
}

output "user_pool_endpoint" {
  value = aws_cognito_user_pool.user_pool.endpoint
}





