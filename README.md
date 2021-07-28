SOLUS IO Terraform plugin
=========================

Used to manage SOLUS IO application.

SOLUS IO is a virtual infrastructure management solution that facilitates
choice, simplicity, and performance for ISPs and enterprises. Offer blazing
fast, on-demand VMs, a simple API, and an easy-to-use self-service control
panel for your customers to unleash your full potential for growth.

[Official site](https://www.solus.io/)

Development
-----------

```shell script
make init
```

Run unit tests:

```shell script
make test
```

Run acceptance tests:

```shell script
export SOLUSIO_BASE_URL="https://localhost/api/v1/"
export SOLUSIO_TOKEN="..."
export SOLUSIO_INSECURE=1

make testacc
```
