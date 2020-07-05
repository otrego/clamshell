Our infrastructure is managed with Terraform.

# Setup

## Install terraform.
Note: version 0.12 is required at this time. 
`brew install terraform`
`brew install google-cloud-sdk`

## Credentials

You will need to have a Google Service Account configured in order for terraform to manage GCP.

[https://cloud.google.com/docs/authentication/production#command-line]
export GOOGLE_APPLICATION_CREDENTIALS="/Users/mdhoak/.gcloud/otrego-dev-8b549007c979.json"

##Set up the terraform project locally
`terraform init`

### Plan updates to the infrastructure based on local changes.
`terraform --plan --out=out.tfplan`

### Apply changes to infrastructure.
`terraform apply "out.tfplan"`
