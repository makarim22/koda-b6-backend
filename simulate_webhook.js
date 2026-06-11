const crypto = require('crypto');
const fs = require('fs');
const path = require('path');

const orderId = process.argv[2];
if (!orderId) {
    console.error("Please provide an Order ID. Usage: node simulate_webhook.js <order_id> [gross_amount]");
    process.exit(1);
}

// Parse custom gross amount if provided, otherwise default to 50000.00
let grossAmount = "50000.00";
if (process.argv[3]) {
    const parsed = parseFloat(process.argv[3]);
    if (!isNaN(parsed)) {
        grossAmount = parsed.toFixed(2);
    } else {
        grossAmount = process.argv[3];
    }
}

// Load configurations from .env dynamically
let serverKey = "test-server-key";
let serverPort = "8081";

try {
    const envPath = path.join(__dirname, '.env');
    if (fs.existsSync(envPath)) {
        const envContent = fs.readFileSync(envPath, 'utf8');
        const envVars = {};
        envContent.split(/\r?\n/).forEach(line => {
            const trimmedLine = line.trim();
            if (trimmedLine && !trimmedLine.startsWith('#')) {
                const parts = trimmedLine.split('=');
                if (parts.length >= 2) {
                    const key = parts[0].trim();
                    const value = parts.slice(1).join('=').trim();
                    envVars[key] = value;
                }
            }
        });
        if (envVars['MIDTRANS_SERVER_KEY']) {
            serverKey = envVars['MIDTRANS_SERVER_KEY'];
        }
        if (envVars['SERVER_PORT']) {
            serverPort = envVars['SERVER_PORT'];
        }
    }
} catch (e) {
    console.warn("Could not read .env file, using default server key and port:", e.message);
}

const statusCode = "200";

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

console.log(`Using Midtrans Server Key: ${serverKey.substring(0, 10)}...`);
console.log("Sending Webhook Payload to port " + serverPort + ":");
console.log(webhookData);

const targetUrl = `http://localhost:${serverPort}/api/payments/midtrans-callback`;

fetch(targetUrl, {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify(webhookData)
})
.then(res => {
    if (!res.ok) {
        return res.text().then(text => {
            throw new Error(`HTTP error ${res.status}: ${text}`);
        });
    }
    return res.json();
})
.then(data => {
    console.log("\nResponse from Webhook API:");
    console.log(data);
})
.catch(err => {
    console.error("\nError calling Webhook API:", err.message || err);
});

