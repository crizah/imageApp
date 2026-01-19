# Serverless Chat Application (AWS)

Event-driven messaging system built on AWS serverless infrastructure. All resources provisioned using **Terraform**.

## Infrastructure Stack

| Service | Purpose |
|---------|---------|
| **Amazon Cognito** | User Pools for authentication, Identity Pools for temporary credentials |
| **AWS Lambda** | Event-driven compute for user provisioning and message handling |
| **Amazon API Gateway** | RESTful API endpoints with Lambda integration |
| **Amazon DynamoDB** | User metadata and message indices |
| **Amazon S3** | Encrypted message payload storage |
| **AWS KMS** | Customer master keys for envelope encryption |
| **Amazon SNS** | Message notification delivery |
 

## Infrastructure as Code

All AWS resources managed via Terraform:
- Cognito User Pool
- DynamoDB tables (Users, Messages)
- S3 buckets with encryption policies
- Lambda functions with IAM execution roles

  

## Architecture

![Sign-Up Flow](ss/signUp.png)

### Sign-Up Flow
- User creates an account through the frontend using **username, email, and password**.
- **Amazon Cognito User Pool** registers the user and associates their identity via an **Identity Pool**.
- User completes **email verification** to activate the account.
- Upon successful verification, a **Lambda function** is triggered to:
  - Create a **user-specific AWS KMS key**
  - Create an **SNS topic** for user notifications
- The Lambda function persists user metadata in the **Users DynamoDB table**.

---

![Message Flow](ss/message.png)

### Message Flow
- The frontend invokes an **API Gateway endpoint** to trigger a Lambda function that retrieves the list of registered users.
- User A uploads a file/message intended for User B.
- The backend retrieves **User B's KMS key reference** from the Users table.
- A **data encryption key (DEK)** is generated and encrypted using User B's KMS key.
- The message payload is encrypted using the DEK.
- The encrypted message is stored in **Amazon S3**.
- Message metadata, including the **S3 object key and encrypted DEK**, is stored in the **Messages DynamoDB table**.
- User B receives a **notification email** indicating a new message.

![SNS Message](ss/sns.png)

#### Message Retrieval
- User B logs in and requests unread message count.
- The backend queries the Messages table to determine unread messages.
- Upon selecting a message:
  - The backend retrieves the encrypted message and encrypted DEK using the message ID.
  - **AWS KMS** decrypts the DEK.
  - The message is decrypted using the DEK.
- The decrypted message is securely delivered to User B.

![Decrypted Message](ss/got.png)

---
## Getting Started

### Prerequisites

* Docker >= 24.0
* Docker Compose >= 2.20

### Run

Clone Repo 
```
git clone https://github.com/crizah/imageApp.git
cd fontend-example
```
Run with docker compose

```
docker compose up --build
```
The frontend should be running on http://localhost:3000

---

## Security Implementation

**Envelope encryption:**
- Messages encrypted with symmetric data encryption keys (AES-256)
- DEKs encrypted with recipient-specific KMS customer master keys
- Encrypted DEK stored in DynamoDB, encrypted payload in S3

**Access control:**
- KMS key policies restrict decrypt operations to key owner
- Cognito Identity Pool provides temporary AWS credentials
- Lambda execution roles follow least privilege principle

## Technical Notes

**Authentication:**
- Cognito issues JWT tokens (ID, access, refresh)
- API Gateway validates JWT signature and expiration
- Lambda receives authenticated user context





