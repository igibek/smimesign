# How-to

## 1. Initialization steps

1. We first need to generate asymmetric key inside AWS. You can follow the instruction here on how to create asymmetric key for signing purpose: https://docs.aws.amazon.com/kms/latest/developerguide/asymm-create-key.html

2. Generate Certificate Signing Request (CSR) using openssl. The CSR is going to be used to request for S/MIME certificate from Public Certificate Authority.

`openssl req -nodes --newkey rsa:4096 -keyout discard.key -out smime.csr`

This above command will generate CSR in interactive mode. Please fill out all relevant information. 
For Command Name (CN) field you will need to provide email address for which S/MIME will be requested. 

Note: we do not need the generate key-pair discard.key because we are going to use AWS KMS to sign our CSR.

3. Sign the generated `smime.csr` with AWS KMS. There is python script aws-kms-sign-csr.py that will help with signing the CSR with key stored inside the AWS KMS.
`.\aws-kms-sign-csr.py --keyid <key-id> --signalgo=ECDSA --hashalgo=sha256 --region=<aws-region> smime.csr > signed-smime.csr`

4. Buy the S/MIME certificate from Public Certificate Authority
## 2. Signing steps

## 3. Verification steps