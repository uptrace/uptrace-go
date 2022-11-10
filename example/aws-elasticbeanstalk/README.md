# AWS Elastic Beanstalk example

To run this example, you need to
[install EB CLI](https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/eb-cli3-install.html).

1. Set Uptrace DSN in `application.go`.

```go
	uptrace.ConfigureOpentelemetry(
		// FIXME
		uptrace.WithDSN("https://<token>@uptrace.dev/<project_id>"),
	)
```

2. Deploy the app to AWS Elastic Beanstalk:

```shell
eb create my-env
```

3. Open the app and click on the trace URL:

```shell
eb open
```
