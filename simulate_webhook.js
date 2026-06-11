const crypto = require('crypto');

const orderId = process.argv[2];
if (!orderId) {
    console.error("Please provide an Order ID. Usage: node simulate_webhook.js <order_id>");
    process.exit(1);
}

const statusCode = "200";
// You should change this to match the actual order's gross amount if you want points to be awarded accurately, 
// though the signature just needs to match whatever is sent in the payload.
const grossAmount = "50000.00"; 
const serverKey = "test-server-key"; // Change this if you set MIDTRANS_SERVER_KEY in .env

// Calculate SHA512 Signature
const payload = orderId + statusCode + grossAmount + serverKey;
const signatureKey = crypto.createHash('sha512').update(payload).digest('hex');

const webhookData = {
    transaction_time: new Date().toISOString(),
    transaction_status: "settlement",
    transaction_id: "mock-midtrans-" + Math.floor(Math.random() * 100000),
    status_code: statusCode,
    signature_key: signatureKey,
    payment_type: "credit_card",
    order_id: orderId,
    gross_amount: grossAmount,
    fraud_status: "accept",
    approval_code: "123456"
};

console.log("Sending Webhook Payload:");
console.log(webhookData);

fetch('http://localhost:3002/api/payments/midtrans-callback', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify(webhookData)
})
.then(res => res.json())
.then(data => {
    console.log("\nResponse from Webhook API:");
    console.log(data);
})
.catch(err => {
    console.error("\nError calling Webhook API:", err);
});
