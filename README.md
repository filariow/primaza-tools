# Primaza Tools

This repository contains a set of companion tools for [Primaza](https://github.com/primaza/primaza).

## primaza-mon

Use this tool print links between workloads and services in a Primaza tenant.
The tool can output the result as JSON, Mermaid, or HTML.

Download the binary or compile it yourself. Then execute `primaza-mon --help` for information on how to use the tool.

### Example

You can setup a test environment using the script `hack/create_env.sh` and then execute the following command

```console
primaza-mon get connections primaza-mytenant -o html > primaza-mytenant-graph.html
```

The result would be similar to the following:

```mermaid
graph TD;
	accTitle: self-demo;
	catalog --> catalog-rds;
	catalog-rds --> rds-postgres;
	orders --> orders-dynamo;
	orders-dynamo --> dynamo;
	catalog --> sqs-catalog;
	sqs-catalog --> sqs-queue-reader;
	orders-events-consumer --> sqs-orders-event;
	sqs-orders-event --> sqs-queue-writer;

	click catalog-rds call callback()
	click orders-dynamo call callback()
	click sqs-catalog call callback()
	click sqs-orders-event call callback()
	click rds-postgres call callback()
	click dynamo call callback()
	click sqs-queue-reader call callback()
	click sqs-queue-writer call callback()
	click catalog call callback()
	click orders call callback()
	click orders-events-consumer call callback()
```
