# Fix PakaiLink RSA Public Key

File `pakailink_rsa_public_key.pem` di webhook-v2 saat ini masih berisi private key. File ini harus diganti dengan public key.

## Cara Generate Public Key dari Private Key

Jika Anda memiliki private key (PKCS#8 format), gunakan perintah berikut untuk generate public key:

```bash
cd /opt/webhook-v2

# Generate public key dari private key (jika Anda punya private key)
openssl rsa -in pkcs8_rsa_private_key.pem -pubout -out pakailink_rsa_public_key.pem

# Atau jika private key dalam format PKCS#8
openssl pkey -in pkcs8_rsa_private_key.pem -pubout -out pakailink_rsa_public_key.pem
```

## Atau Request Public Key dari PakaiLink

Jika Anda tidak memiliki private key, minta public key langsung dari PakaiLink support.

## Format Public Key

Public key harus dalam format PEM:

```
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...
...
-----END PUBLIC KEY-----
```

**BUKAN** format private key:
```
-----BEGIN PRIVATE KEY-----
...
-----END PRIVATE KEY-----
```

atau

```
-----BEGIN RSA PRIVATE KEY-----
...
-----END RSA PRIVATE KEY-----
```

## Setelah Update

Setelah file public key sudah benar:

```bash
# Set permission
chmod 644 pakailink_rsa_public_key.pem

# Restart webhook service
cd /opt/webhook-v2
docker compose restart webhook-v2

# Check logs
docker compose logs -f webhook-v2
```

