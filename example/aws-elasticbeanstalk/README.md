# AWS Elastic Beanstalk example

To run this example, you need to
[install EB CLI](https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/eb-cli3-install.html).

1. Set Uptrace DSN in `application.go`.
2. Deploy the app to AWS:

```shell
eb create my-env
```

3. Open the app:

```shell
eb open
```
