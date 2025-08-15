# Hardware crypto storage

## Trusted Platform Module

Non-volatile storage is about 64KB in size

TPM has two main use-cases:

- secure key generation
- remote system attestation

Platform Configuration Registers (PCRs)

Create primary object and save password in TPM

```sh
# Create a primary key in the owner hierarchy
tpm2_createprimary -C o -c primary.ctx

# Generate a random password and seal it in the TPM
echo -n "MySecretPassword123" > password.txt
tpm2_create -C primary.ctx -i password.txt -u seal.pub -r seal.priv -L 0x0004:0,1,2,3,4,5,6,7

# Persist the sealed object at a fixed handle
tpm2_load -C primary.ctx -u seal.pub -r seal.priv -c seal.ctx
tpm2_evictcontrol -C o -c seal.ctx 0x81000000
```

Retrieve the password

```sh
# Unseal the password when needed
tpm2_unseal -c 0x81000000 -o password2.txt

# View the password (for demonstration only - be careful!)
cat password2.txt

# Immediately clear from memory
shred -u password.txt
```

- `shred` - overwrite a file to hide its contents, and optionally delete it

## Apple Secure Enclave
