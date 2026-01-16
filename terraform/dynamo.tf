
resource "aws_dynamodb_table" "messages" {
    // Messages
    // messageID is pk
    // recipient-index is index name for gsi with attr recipient


  name           = var.msgs
  billing_mode = "PAY_PER_REQUEST"
  hash_key       = var.msgspk

  attribute {
    name = var.msgspk
    type = "S"
  }
  attribute {
    name = "recipient"
    type = "S"
}

  global_secondary_index {
    name            = var.msgsgsi
    hash_key        = "recipient"
    projection_type = "ALL"
  }
}

resource "aws_dynamodb_table" "users" {
    // Users
    // username is pki


  name           = var.users
  billing_mode = "PAY_PER_REQUEST"
  hash_key       = var.userspk

  attribute {
    name = var.userspk
    type = "S"
  }
}
