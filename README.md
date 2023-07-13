# Primaza Tools

This repository contains a set of companion tools for [Primaza](https://github.com/primaza/primaza).

## primaza-mon

Use this tool to identify links between workloads and services in a Primaza tenant.
The tool can output the result as JSON, Mermaid, or HTML.

Download the binary or compile it yourself. Then execute `primaza-mon --help` for information on how to use the tool.

### Example

You can setup a test environment using the script `hack/create_env.sh`.
Then execute the following command:

```console
primaza-mon get connections primaza-mytenant -o html > primaza-mytenant-graph.html
```

The result would be similar to the following:

```mermaid
graph TD;
        accTitle: self-demo;
        catalog[/catalog\] --> catalog-rds{{catalog-rds}} --> rds-postgres[\rds-postgres/];
        orders[/orders\] --> orders-dynamo{{orders-dynamo}} --> dynamo[\dynamo/];
        catalog[/catalog\] --> sqs-catalog{{sqs-catalog}} --> sqs-queue-reader[\sqs-queue-reader/] --> my-queue[my-queue];
        orders-events-consumer[/orders-events-consumer\] --> sqs-orders-event{{sqs-orders-event}} --> sqs-queue-writer[\sqs-queue-writer/] --> my-queue[my-queue];

        click catalog-rds call callback()
        click orders-dynamo call callback()
        click sqs-catalog call callback()
        click sqs-orders-event call callback()
        click dynamo call callback()
        click sqs-queue-reader call callback()
        click sqs-queue-writer call callback()
        click rds-postgres call callback()
        click catalog call callback()
        click orders call callback()
        click orders-events-consumer call callback()
        click my-queue call callback()
```
