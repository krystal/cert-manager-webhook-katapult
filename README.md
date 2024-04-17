# Cert Manager Webhook for Katapult

This is a Cert Manager Webhook service to faciliate DNS challenges with Katapult's DNS platform.

## Installation

To install, just runs the following to install or upgrade on your cluster.

```
helm upgrade --install oci://ghcr.io/krystal/charts/cert-manager-webhook-katapult -n cert-manager
```

You'll need to add a secret containing an API key for your Katapult account. This is referenced by the `Issuer` or `ClusterIssuer` which uses this webhook.

## Example issuer

```yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: letsencrypt-dns-issuer
spec:
  acme:
    email: demo@example.com
    server: https://acme-v02.api.letsencrypt.org/directory
    # Use this instead for staging.
    # server: https://acme-staging-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-dns-issuer-secret
    solvers:
      - dns01:
          webhook:
            groupName: acme.katapult.io
            solverName: katapult
            config:
              apiToken:
                name: katapult-secret
                key: token
```
