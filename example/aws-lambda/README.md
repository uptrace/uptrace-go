# AWS Lambda example for OpenTelemetry and Uptrace

- [Documentation](https://uptrace.dev/opentelemetry/instrumentations/go-aws-lambda.html)

To build AWS lambda function:

```shell
cd function
./build.sh
```

To deploy the function using Terraform:

```shell
cd deploy
terraform init
terraform apply
```

Then open the gateway proxy to trigger the lambda function.
