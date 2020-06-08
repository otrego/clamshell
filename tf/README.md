Our infratructure is managed with Terraform.

Install terraform
`brew install terraform`
`brew install google-cloud-sdk`

## Credentials

###Set up the terraform project locally
`terraform init`

### Plan updates to the infrastructure based on local changes.
`terraform --plan --out=out.tf`

### Apply changes to infrastructure.
`terraform apply "out.tf"`
