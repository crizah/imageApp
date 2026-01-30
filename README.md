# ImageApp 
A containerised cloud based chatting application built using AWS architecture with encrypted data storage, automated infrastructure provisionings and pipelines. 

---
## Technical stack
The technoligies used: 

| Service                | Technology used                                                  |
| ---------------------- | ------------------------------------------------------------------ |
| **Frontend**     |Built in React using Tailwind CSS        |
| **Backend**         | Built in Golang|
| **AWS services** | For storage, authentication, serverless infrastructure|
| **Terraform**    | Used to provison all the AWS servcies                             |
| **Docker**          | Used to containerise each component of the app                             |
| **Docker compose**            | Used to run and manage the multicontainer application              |
| **Jenkins**         | Pipeline to automate building and pushing of Docker images                            |


### Repository Structure
```
.
├── server/                # Go server
|   ├── Dockerfile
├── web/                    # React frontend 
|   ├── Dockerfile                  
├── terraform/             # All .tf files for AWS provisioning
│   ├── cognito.tf
|   ├── lambda/                # Contains the zip folder to be run on lamda 
│   ├── dynamodb.tf
│   ├── ec2.tf
│   ├── lambda.tf
│   └── s3.tf
├── docker-compose.yml    
├── Jenkinsfile            # pipeline definition
└── README.md
```


## AWS Services used
All the AWS services used: 

| Service                | Purpose                                                            |
| ---------------------- | ------------------------------------------------------------------ |
| **Amazon Cognito**     | User authentication via user pools and email verification after signup    |
| **AWS Lambda**         | Triggered after successful cognito verification to add user information to db |
| **Amazon DynamoDB**    | User informtaion and Message metadata                         |
| **Amazon S3**          | Stores encrypted image based messages                                |
| **AWS KMS**            | Customer managed keys for encryption                     |
| **Amazon SNS**         | Notification delivery for new messages                             |
| **Amazon EC2**         | Deploy the application|

---

## Infrastructure as Code (Terraform)
All AWS resources are provisioned using **Terraform** :
* All the AWS resources are defined in .tf files
* Terraform makes all the required infrastructure (Dynamo tables S3 buckets, cognito userpools, lamda functions) and makes and attatches the required iam roles and policies.
* Terraform makes an EC2 instance, installs the dependencies, clones the repo inside it and runs it with docker compose.

---

## CI/CD Automation (Jenkins)

This repository includes a **Jenkinsfile** whose build when ran:
* Detects any new commit made to the remote githib repo.
* If there is a change in ./server or ./web folders, It is detected and the docker images for those respective folders is rebuilt.
* The images are tagged and pushed to docker hub.

---

## Docker compose
This is a multi container application and docker compose is used to run and manage the containers.

## Architecture

![Sign-Up Flow](ss/signUp.png)

### Sign-Up Flow

* User signs up via frontend with **username, email, and password**
* **Cognito User Pool** made by terraform, registers the user 
* Cognito sends email verification to the entered emailID
* User enters the code into the frontend and is verified
* A post-confirmation **Lambda function** is triggered that:

  * Creates a **KMS key** for the user
  * Creates an **SNS topic** for the user
* User data is stored in the **Users DynamoDB table**

---

![Message Flow](ss/message.png)

### Message Flow

* Frontend fetches all registered users and a receiver is selected
* User A sends a file to User B
* Backend retrieves **User B’s KMS key reference**
* A **data encryption key (DEK)** is generated
* DEK is encrypted using User B’s KMS key
* Message is encrypted using the DEK
* Encrypted image is stored in **S3**
* Metadata (S3 key + encrypted DEK) is stored in **Messages DynamoDB**
* **SNS** notifies User B of a new message


---

### Message Retrieval

* User B logs in and queries unread messages
* Backend fetches message metadata from DynamoDB
* Encrypted DEK is decrypted via **KMS**
* Message is decrypted and displayed


---

## Getting Started (Terraform + Docker)

### Prerequisites

* Terraform >= 1.6
* Docker >= 24.0
* Docker Compose >= 2.20
* AWS account with sufficient permissions
* An existing EC2 key pair

---

### Terraform Setup

edit the `terraform.tfvars` file at `./terraform/terraform.tfvars`:

```hcl
ssh_key    = "your-ec2-keypair-name"
ur_ip      = "YOUR_PUBLIC_IP"
```

**Notes:**

* `ssh_key` must match an EC2 key pair already created in AWS
* `ur_ip` should be **your public IP**, used to restrict SSH access (`/32`)

---

### Provision Infrastructure

```bash
terraform init
terraform apply -var-file terraform.tfvars
```

Terraform will:

* Provision all AWS resources
* Attach IAM roles to EC2 
* Install all dependencies onto the EC2 instances
* Allocate a public Elastic IP
* Clone this repository
* Run `docker compose up --build`

Once complete, Terraform outputs the application URL.
The app can be access via port `3000` of the public ip

---

## Local Development (Docker)

Pull the remote docker images

```bash
docker pull shaizah/kube:imageApp-server
docker pull shaizah/kube:imageApp-web
```
Run the images

```bash
docker run -d -p 8082:8082 shaizah/kube:imageApp-server
docker run -d -p 3000:80 -e BACKEND_URL=http://localhost:8082 shaizah/kube:imageApp-web
```

Frontend will be available at:

```
http://localhost:3000
```

---

## Screenshots
### Signup page
![Signup page](ss/signUp_web.png)

### Home page
![HomePage](ss/dashboard.png)

### Send Message
![Send Message](ss/sendmsg.png)

### Get Message
![Messages](ss/getmsgs.png)


##  What I Learned

- **AWS Services**: Hands-on experience with 8+ AWS services (Cognito, Lambda, KMS, DynamoDB, S3, SNS, EC2, API Gateway)
- **Infrastructure as Code**: Transitioned from AWS console to Terraform to allocate resources.
- **Encryption**: Implemented envelope encryption from scratch instead of turing on encryption in S3
- **DevOps Practices**: Set up CI/CD pipeline to automate redundant tasks

