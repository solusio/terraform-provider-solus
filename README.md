SOLUS Terraform plugin
=========================

Used to manage SOLUS application.

SOLUS is a virtual infrastructure management solution that facilitates
choice, simplicity, and performance for ISPs and enterprises. Offer blazing
fast, on-demand VMs, a simple API, and an easy-to-use self-service control
panel for your customers to unleash your full potential for growth.

[Official site](https://www.solus.io/)

Development
-----------

```shell
make init
```

Run unit tests:

```shell
make test
```

Run acceptance tests:

```shell
cp .testacc.env.example .testacc.env
vim .testacc.env
make testacc
```
