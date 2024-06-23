# Demo for OpenTelemetry, Simplifying Observability for Cloud-Native Environments

The golang-response-echoer is a simple web server designed to demonstrate the capabilities of OpenTelemetry for tracing and monitoring web applications. This server responds dynamically based on incoming request parameters, integrating OpenTelemetry to trace the sequence and status of these requests.

## Prerequisites

Ensure that you have Docker installed on your machine. If not, you can download it from Docker's official site.

### Setup Instructions

1. Clone the repository:

```bash
git clone https://github.com/NoNickeD/otel-demo.git

cd golang-echoer-server
```

2. Build the Docker image:

```bash
docker build -t golang-response-echoer .
```

3. Run the container:

```bash
docker run -p 8080:8080 golang-response-echoer
```

## Usage

The server can handle various commands sent via HTTP requests. Below are examples of how you can interact with the server using `curl`:

**Echo Status Code Sequence**

- Initialize the sequence:

```bash
curl -i "http://localhost:8080/?echo_code=200-400-500&init=1"
```

**Echo Current Time**

- Get the current server time:

```bash
curl "http://localhost:8080/?echo_time"
```

**Echo Request Headers**

- Send a custom header and echo it back:

```bash
curl -H "Custom-Header: Value" "http://localhost:8080/?echo_env"
```

**Echo Request Body**

- Send a request body and echo it back:

```bash
curl -d "This is a test body" "http://localhost:8080/?echo_body"
```

## Visualizing Trace Data

Once the server is running, OpenTelemetry will capture and visualize the trace data, which can be viewed to understand the flow of requests and identify any issues or performance bottlenecks.

```json
{"level":"info","msg":"Responded with HTTP status code: 400","time":"2024-06-23T18:58:44Z"}
{
	"Name": "handleStatusCodeSequence",
	"SpanContext": {
		"TraceID": "f78e898a6dd2b20d877b7ee52a6ccf93",
		"SpanID": "85c5cabc39d13c9b",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"Parent": {
		"TraceID": "f78e898a6dd2b20d877b7ee52a6ccf93",
		"SpanID": "91baf8b3c925c74a",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"SpanKind": 1,
	"StartTime": "2024-06-23T18:58:44.714478553Z",
	"EndTime": "2024-06-23T18:58:44.714557261Z",
	"Attributes": null,
	"Events": null,
	"Links": null,
	"Status": {
		"Code": "Unset",
		"Description": ""
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 0,
	"Resource": [
		{
			"Key": "service.name",
			"Value": {
				"Type": "STRING",
				"Value": "unknown_service:golang-echoer-server"
			}
		},
		{
			"Key": "telemetry.sdk.language",
			"Value": {
				"Type": "STRING",
				"Value": "go"
			}
		},
		{
			"Key": "telemetry.sdk.name",
			"Value": {
				"Type": "STRING",
				"Value": "opentelemetry"
			}
		},
		{
			"Key": "telemetry.sdk.version",
			"Value": {
				"Type": "STRING",
				"Value": "1.27.0"
			}
		}
	],
	"InstrumentationLibrary": {
		"Name": "echo-server",
		"Version": "",
		"SchemaURL": ""
	}
}
{
	"Name": "handleRequests",
	"SpanContext": {
		"TraceID": "f78e898a6dd2b20d877b7ee52a6ccf93",
		"SpanID": "91baf8b3c925c74a",
		"TraceFlags": "01",
		"TraceState": "",
		"Remote": false
	},
	"Parent": {
		"TraceID": "00000000000000000000000000000000",
		"SpanID": "0000000000000000",
		"TraceFlags": "00",
		"TraceState": "",
		"Remote": false
	},
	"SpanKind": 1,
	"StartTime": "2024-06-23T18:58:44.714445845Z",
	"EndTime": "2024-06-23T18:58:44.714564261Z",
	"Attributes": null,
	"Events": null,
	"Links": null,
	"Status": {
		"Code": "Unset",
		"Description": ""
	},
	"DroppedAttributes": 0,
	"DroppedEvents": 0,
	"DroppedLinks": 0,
	"ChildSpanCount": 1,
	"Resource": [
		{
			"Key": "service.name",
			"Value": {
				"Type": "STRING",
				"Value": "unknown_service:golang-echoer-server"
			}
		},
		{
			"Key": "telemetry.sdk.language",
			"Value": {
				"Type": "STRING",
				"Value": "go"
			}
		},
		{
			"Key": "telemetry.sdk.name",
			"Value": {
				"Type": "STRING",
				"Value": "opentelemetry"
			}
		},
		{
			"Key": "telemetry.sdk.version",
			"Value": {
				"Type": "STRING",
				"Value": "1.27.0"
			}
		}
	],
	"InstrumentationLibrary": {
		"Name": "echo-server",
		"Version": "",
		"SchemaURL": ""
	}
}
```
