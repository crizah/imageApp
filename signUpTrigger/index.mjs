import { DynamoDBClient, PutItemCommand } from "@aws-sdk/client-dynamodb";
import { KMSClient, CreateKeyCommand } from "@aws-sdk/client-kms";
import { SNSClient, CreateTopicCommand, SubscribeCommand } from "@aws-sdk/client-sns";

const dynamoClient = new DynamoDBClient({ region: "eu-north-1" });
const kmsClient = new KMSClient({ region: "eu-north-1" });
const snsClient = new SNSClient({ region: "eu-north-1" });

export const handler = async (event) => {
  console.log("=== POST CONFIRMATION TRIGGER ===");
  console.log("Full event:", JSON.stringify(event, null, 2));

  try {
    // Cognito trigger structure
    const username = event.userName;
    const email = event.request.userAttributes.email;
    
    console.log("Username:", username);
    console.log("Email:", email);

    if (!username) {
      throw new Error("Missing username");
    }
    
    if (!email) {
      console.warn("Email not found, using placeholder");
      // Some sign-up flows might not require email
    }

    // 1. Create KMS key
    console.log("Creating KMS key...");
    const keyResult = await kmsClient.send(new CreateKeyCommand({
      Description: `Key for user ${username}`,
      Tags: [
        {
          TagKey: 'Username',
          TagValue: username
        }
      ]
    }));

    const kmsKeyId = keyResult.KeyMetadata.KeyId;
    console.log("✅ Created KMS key:", kmsKeyId);

    // 2. Create SNS topic
    console.log("Creating SNS topic...");
    const topicResult = await snsClient.send(new CreateTopicCommand({
      Name: `user-${username}-notifications`
    }));

    const topicArn = topicResult.TopicArn;
    console.log("✅ Created SNS topic:", topicArn);

    // 3. Subscribe email if available
    if (email) {
      console.log("Subscribing email to SNS...");
      await snsClient.send(new SubscribeCommand({
        TopicArn: topicArn,
        Protocol: 'email',
        Endpoint: email
      }));
      console.log("✅ Subscribed email to SNS");
    } else {
      console.log("⚠️ Skipping email subscription (no email provided)");
    }

    // 4. Save to DynamoDB
    console.log("Saving to DynamoDB...");
    const params = {
      TableName: "Users",
      Item: {
        username: { S: username },
        email: { S: email || "no-email-provided" },
        kmsKeyId: { S: kmsKeyId },
        snsTopicArn: { S: topicArn },
        createdAt: { S: new Date().toISOString() }
      }
    };

    await dynamoClient.send(new PutItemCommand(params));
    console.log("✅ User saved to DynamoDB");

    console.log("=== SUCCESS ===");
    return event;  // MUST return event for Cognito

  } catch (err) {
    console.error("❌ ERROR:", err);
    console.error("Error stack:", err.stack);
    // Return event anyway to not block user registration
    return event;
  }
};