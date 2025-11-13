# How-to

## 1. Initialization steps

1. We first need to generate asymmetric key inside AWS. You can follow the instruction here on how to create asymmetric key for signing purpose: https://docs.aws.amazon.com/kms/latest/developerguide/asymm-create-key.html
    - Select as Key Spec: ECC_NIST_P384
    - Allow key users only following permissions:
        - SIGN
        - VERIFY
        - GetPublicKey
        - DescribeKey
2. Generate Certificate Signing Request (CSR) using openssl. The CSR is going to be used to request for S/MIME certificate from Public Certificate Authority.

`openssl req -nodes --newkey rsa:4096 -keyout discard.key -out smime.csr`

This above command will generate CSR in interactive mode. Please fill out all relevant information. 
For Command Name (CN) field you will need to provide email address for which S/MIME will be requested. 

Note: we do not need the generated key-pair discard.key because we are going to use AWS KMS to sign our CSR.

3. Sign the generated `smime.csr` with AWS KMS. There is python script aws-kms-sign-csr.py that will help with signing the CSR with key stored inside the AWS KMS.
`.\aws-kms-sign-csr.py --keyid <key-id> --signalgo=ECDSA --hashalgo=sha256 --region=<aws-region> smime.csr > signed-smime.csr`

4. Buy the S/MIME certificate from Public Certificate Authority (e.g. SSL.com). When they ask, you will need to provide signed-smime.csr.

5. Public CA will send you the email to verify that you indeed own the EMAIL listed in CSR file. 
6. Download a generated certificate

## 2. Signing steps
For signing purposes we use a custom tool called QCSIGN (fork of https://github.com/github./smimesign) tool that can use AWS KMS to sign the commits.

1. Build the tool using `go build` command
2. Add the generated `qcsign` into the PATH.
3. Create config file: `$HOME/.qcsign/config.yaml`. You can check out provided example `config.yaml` file 
4. Configure git to use qcsign to sign the commits
    1. `git config gpg.format x509` - this tells git to use x509 format to sign the commits
    2. `git config gpg.x509.program qcsign` - this tells git to use qcsign tool to sign the commits
    3. `qcsign --list-keys` - this will list the available keys for signing. Right now it uses CERT_PATH value from `$HOME/.qcsign/config.yaml` file
        - copy ID value, e.g. bfa1f3e8a45351edda8dd94922efaddfbdc0dd3b
    4. `git config user.signingkey <ID from previous step>`, e.g. `git config user.signingkey bfa1f3e8a45351edda8dd94922efaddfbdc0dd3b` - this will tell git which certificate to use for signing
5. Now for each git commit command you can add -S option, as here `git commit -S -m "message"` - this will tell to sign the commit

## 3. Verification steps


## Debugging:

Run the command to view the certificates: ```git cat-file commit HEAD | sed -n '/BEGIN/, /END/p' | sed 's/^ //g' | sed 's/gpgsig //g' | sed 's/SIGNED MESSAGE/PKCS7/g' | openssl pkcs7 -print -print_certs -text```