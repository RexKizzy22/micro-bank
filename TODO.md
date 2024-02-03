# TODO

- Implement Log Out API
- Fix TransferTx bug
- Write tests for transfer API
- Complete other tests for Users API
- Write TestLoginUserAPI
- Write tests for entry.sql.go queries
- Write tests for transfer.sql.go queries
- Write custom logger to customize the first 3 redis queue log to ensure consistent logging format

## Deployment Strategy

1. Using GCP App Engine

   - Create app.yaml
   - Create service account on GCP
   - Deploy

2. Using kubernetes
   - Use AWS
   - Provision cluster
   - Provision node group
   - Update kubeconfig to use AWS credentials
   - Apply deployment, service, ingress and ingress class
   - Purchase domain from Route53
   - Create DNS record type of A record
   - Set up SSL/TLS certificate with letsencrypt

## Configure Automatic TLS Certificate Management

- Install cert manager (kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.0/cert-manager.yaml)
- Configure ACME protocol with ClusterIssuer (kubectl apply -f issuer.yaml)
