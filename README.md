# Encryption Service

Data encryption and tokenization service for KodraPay platform.

## Overview

The Encryption Service provides AES-256-GCM encryption for sensitive data and tokenization for payment card information. It helps protect payment link amounts, customer data, and card details.

## Features

- **AES-256-GCM Encryption** - Industry-standard authenticated encryption
- **Data Encryption/Decryption** - Encrypt and decrypt arbitrary data
- **Card Tokenization** - Securely tokenize payment card information
- **Token Retrieval** - Retrieve tokenized card details
- **In-Memory Token Storage** - Fast access to tokenized data (production should use database)

## Configuration

### Environment Variables

Create a `.env` file:

```bash
PORT=7016
ENCRYPTION_KEY=your-32-byte-encryption-key-here
```

### Encryption Key Setup

**CRITICAL:** The encryption service requires a secure encryption key.

1. **Generate a secure 32-byte encryption key:**
   ```bash
   # Generate hex-encoded 32-byte key
   openssl rand -hex 32
   ```

2. **Set the encryption key in your environment:**
   ```bash
   export ENCRYPTION_KEY=$(openssl rand -hex 32)
   ```

3. **For Docker deployment:** Update `docker-compose.yml`:
   ```yaml
   environment:
     ENCRYPTION_KEY: kodrapay-super-secret-key-2024-aes256
   ```

   **IMPORTANT: Change this key in production!**

## API Endpoints

### Health Check
```
GET /health
```

### Encrypt Data
```
POST /api/v1/encrypt
Headers: Content-Type: application/json
Body: {
  "data": "sensitive information"
}

Response: {
  "encrypted_data": "base64-encoded-encrypted-data"
}
```

### Decrypt Data
```
POST /api/v1/decrypt
Headers: Content-Type: application/json
Body: {
  "encrypted_data": "base64-encoded-encrypted-data"
}

Response: {
  "data": "decrypted sensitive information"
}
```

### Tokenize Card
```
POST /api/v1/keys/tokenize
Headers: Content-Type: application/json
Body: {
  "pan": "4111111111111111",
  "cvv": 123,
  "reference": 12345
}

Response: {
  "token_id": 1001,
  "last4": "1111",
  "brand": "visa"
}
```

### Get Token
```
GET /api/v1/keys/:id

Response: {
  "token_id": 1001,
  "last4": "1111",
  "brand": "visa"
}
```

## Usage Examples

### Encrypting Payment Link Amounts

```bash
# Encrypt amount for payment link
curl -X POST http://localhost:7016/api/v1/encrypt \
  -H "Content-Type: application/json" \
  -d '{"data": "50000"}'

# Response: {"encrypted_data": "...base64..."}
```

### Decrypting Payment Link Amounts

```bash
# Decrypt amount from payment link
curl -X POST http://localhost:7016/api/v1/decrypt \
  -H "Content-Type: application/json" \
  -d '{"encrypted_data": "...base64..."}'

# Response: {"data": "50000"}
```

### Tokenizing Card Data

```bash
# Tokenize card PAN
curl -X POST http://localhost:7016/api/v1/keys/tokenize \
  -H "Content-Type: application/json" \
  -d '{
    "pan": "4111111111111111",
    "cvv": 123,
    "reference": 12345
  }'

# Response: {"token_id": 1001, "last4": "1111", "brand": "visa"}
```

## Encryption Algorithm

- **Algorithm:** AES-256-GCM (Galois/Counter Mode)
- **Key Size:** 256 bits (32 bytes)
- **Features:**
  - Authenticated encryption (provides confidentiality and integrity)
  - Random nonce for each encryption
  - Base64 encoding for transport

## Security Features

1. **Authenticated Encryption** - Detects tampering attempts
2. **Random Nonces** - Each encryption uses a unique nonce
3. **Key Management** - Encryption key loaded from environment
4. **No Logging** - Sensitive data is never logged

## Integration with Payment Links

To use encryption in payment links:

```javascript
// Frontend: Encrypt amount before creating payment link
const response = await fetch('http://localhost:7016/api/v1/encrypt', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ data: '50000' })
});
const { encrypted_data } = await response.json();

// Include encrypted_data in payment link
const link = `${baseUrl}?encrypted_amount=${encodeURIComponent(encrypted_data)}`;
```

```go
// Backend: Decrypt amount when processing payment link
func processPaymentLink(encryptedAmount string) {
    resp, err := http.Post(
        "http://encryption-service:7016/api/v1/decrypt",
        "application/json",
        strings.NewReader(fmt.Sprintf(`{"encrypted_data":"%s"}`, encryptedAmount)),
    )
    // Parse response and use decrypted amount
}
```

## Running the Service

### Local Development
```bash
cd encryption-service
export ENCRYPTION_KEY=$(openssl rand -hex 32)
go run ./cmd/encryption-service
```

### Docker
```bash
docker-compose up encryption-service
```

### Build
```bash
go build -o bin/encryption-service ./cmd/encryption-service
./bin/encryption-service
```

## Supported Card Brands

The tokenization service detects the following card brands:

- **Visa** - Starts with 4
- **Mastercard** - Starts with 5
- **American Express** - Starts with 37
- **Discover** - Starts with 6

## Production Considerations

### Security Best Practices

1. **Key Management**
   - Use a secure key management service (AWS KMS, HashiCorp Vault, etc.)
   - Rotate encryption keys periodically
   - Never commit keys to version control
   - Use different keys for each environment

2. **Token Storage**
   - Current implementation uses in-memory storage
   - For production, implement database persistence
   - Consider token expiration and cleanup

3. **HTTPS**
   - Always use HTTPS/TLS in production
   - Protect encrypted data in transit

4. **Access Control**
   - Limit access to encryption service
   - Use API keys or service authentication
   - Implement rate limiting

### Performance

- Encryption: ~100μs per operation
- In-memory token lookup: ~1μs
- Concurrent operations supported via mutex

## Troubleshooting

### "encryption failed" error
- Check that ENCRYPTION_KEY is set
- Verify key is at least 32 bytes
- Ensure data field is not empty

### "decryption failed" error
- Verify encrypted_data is valid base64
- Ensure same encryption key is used
- Check for data corruption

### "ciphertext too short" error
- Encrypted data is corrupted or incomplete
- Verify complete base64 string is sent

## Future Enhancements

- [ ] Database persistence for tokens
- [ ] Token expiration and cleanup
- [ ] Key rotation support
- [ ] Hardware Security Module (HSM) integration
- [ ] Audit logging for encryption operations
- [ ] Support for additional encryption algorithms
- [ ] Batch encryption/decryption endpoints
- [ ] Token metadata storage
