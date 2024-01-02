# webhook-router

![](docs/images/webhookd-arch.png)

Webhook-Router is a tool designed to receive, transform, and forward webhooks. It acts as a middleware, receiving incoming webhooks from various sources, applying specified JSON transformations, and then forwarding them to their designated endpoints. This tool is helps integrating disparate systems, streamlining data flow, and customizing webhook handling.

_Receivers_ handle the actual webhook side of things, doing basic sanity checking and validation. Next the _Receivers_ can hand the message over to an optional chain of _Transformations_ which can customize the received message further. At the end of the _Transformation_ chain the message is handed over to the _Dispatchers_ which route the webhooks to the final destinations.

## Quickstart

First start the `webhook-router`:

```bash
docker-compose up -d
```

Then we call `webhook-router` at the endpoint `/log` with an example payload:

```bash
curl  \
	-H "Content-Type: application/json" \
	-v 'http://localhost:8080/echo' -d '
{                                                                                                            
  "schemaId": "Microsoft.Insights/activityLogs",              
  "data": {                                                                      
    "status": "Activated",                                                
    "context": {                                                           
      "activityLog": {                                                    
        "channels": "Admin, Operation",                                    
        "correlationId": "a1be61fd-37ur-ba05-b827-cb874708babf",
        "eventSource": "ResourceHealth",      
        "eventTimestamp": "2018-09-04T23:09:03.343+00:00",
        "eventDataId": "2b37e2d0-7bda-4de7-ur8c6-1447d02265b2",                                                                                      
        "level": "Informational",               
        "operationName": "Microsoft.Resourcehealth/healthevent/Activated/action",                                                                    
        "operationId": "2b37e2d0-7bda-489f-81c6-1447d02265b2",
        "properties": {                                                       
          "title": "Virtual Machine health status changed to unavailable",
          "details": "Virtual machine has experienced an unexpected event",
          "currentHealthStatus": "Unavailable",                    
          "previousHealthStatus": "Available",             
          "type": "Downtime",
          "cause": "PlatformInitiated"                    
          },
        "resourceId": "/subscriptions/<subscription Id>/resourceGroups/<resource group>/providers/Microsoft.Compute/virtualMachines/<resource name>",
        "resourceGroupName": "<resource group>",
        "resourceProviderName": "Microsoft.Resourcehealth/healthevent/action",
        "status": "Active",
        "subscriptionId": "<subscription Id>",
        "submissionTimestamp": "2018-09-04T23:11:06.1607287+00:00",
        "resourceType": "Microsoft.Compute/virtualMachines"
      }
    }
  }
}
'
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /echo HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 2102
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Type: application/json
< X-Webhookd-Time-To-Dispatch: 50.130875ms
< X-Webhookd-Time-To-Process: 50.368541ms
< X-Webhookd-Time-To-Receive: 236.292µs
< X-Webhookd-Time-To-Transform: 708ns
< Date: Tue, 02 Jan 2024 15:01:12 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```

Now you can check the logs of the `webhook-echo` container.

```bash
docker logs webhook-echo
Listening on ports 8888 for http, and 9999 for https.
{
    "path": "/",
    "headers": {
        "host": "webhook-echo:8888",
        "user-agent": "Go-http-client/1.1",
        "content-length": "2102",
        "content-type": "application/json",
        "accept-encoding": "gzip"
    },
    "method": "POST",
    "body": "\n{                                                                                                            \n  \"schemaId\": \"Microsoft.Insights/activityLogs\",              \n  \"data\": {                                                                      \n    \"status\": \"Activated\",                                                \n    \"context\": {                                                           \n      \"activityLog\": {                                                    \n        \"channels\": \"Admin, Operation\",                                    \n        \"correlationId\": \"a1be61fd-37ur-ba05-b827-cb874708babf\",\n        \"eventSource\": \"ResourceHealth\",      \n        \"eventTimestamp\": \"2018-09-04T23:09:03.343+00:00\",\n        \"eventDataId\": \"2b37e2d0-7bda-4de7-ur8c6-1447d02265b2\",                                                                                      \n        \"level\": \"Informational\",               \n        \"operationName\": \"Microsoft.Resourcehealth/healthevent/Activated/action\",                                                                    \n        \"operationId\": \"2b37e2d0-7bda-489f-81c6-1447d02265b2\",\n        \"properties\": {                                                       \n          \"title\": \"Virtual Machine health status changed to unavailable\",\n          \"details\": \"Virtual machine has experienced an unexpected event\",\n          \"currentHealthStatus\": \"Unavailable\",                    \n          \"previousHealthStatus\": \"Available\",             \n          \"type\": \"Downtime\",\n          \"cause\": \"PlatformInitiated\"                    \n          },\n        \"resourceId\": \"/subscriptions/<subscription Id>/resourceGroups/<resource group>/providers/Microsoft.Compute/virtualMachines/<resource name>\",\n        \"resourceGroupName\": \"<resource group>\",\n        \"resourceProviderName\": \"Microsoft.Resourcehealth/healthevent/action\",\n        \"status\": \"Active\",\n        \"subscriptionId\": \"<subscription Id>\",\n        \"submissionTimestamp\": \"2018-09-04T23:11:06.1607287+00:00\",\n        \"resourceType\": \"Microsoft.Compute/virtualMachines\"\n      }\n    }\n  }\n}\n",
    "fresh": false,
    "hostname": "webhook-echo",
    "ip": "::ffff:192.168.176.3",
    "ips": [],
    "protocol": "http",
    "query": {},
    "subdomains": [],
    "xhr": false,
    "os": {
        "hostname": "dd1b9f9bb339"
    },
    "connection": {},
    "json": {
        "schemaId": "Microsoft.Insights/activityLogs",
        "data": {
            "status": "Activated",
            "context": {
                "activityLog": {
                    "channels": "Admin, Operation",
                    "correlationId": "a1be61fd-37ur-ba05-b827-cb874708babf",
                    "eventSource": "ResourceHealth",
                    "eventTimestamp": "2018-09-04T23:09:03.343+00:00",
                    "eventDataId": "2b37e2d0-7bda-4de7-ur8c6-1447d02265b2",
                    "level": "Informational",
                    "operationName": "Microsoft.Resourcehealth/healthevent/Activated/action",
                    "operationId": "2b37e2d0-7bda-489f-81c6-1447d02265b2",
                    "properties": {
                        "title": "Virtual Machine health status changed to unavailable",
                        "details": "Virtual machine has experienced an unexpected event",
                        "currentHealthStatus": "Unavailable",
                        "previousHealthStatus": "Available",
                        "type": "Downtime",
                        "cause": "PlatformInitiated"
                    },
                    "resourceId": "/subscriptions/<subscription Id>/resourceGroups/<resource group>/providers/Microsoft.Compute/virtualMachines/<resource name>",
                    "resourceGroupName": "<resource group>",
                    "resourceProviderName": "Microsoft.Resourcehealth/healthevent/action",
                    "status": "Active",
                    "subscriptionId": "<subscription Id>",
                    "submissionTimestamp": "2018-09-04T23:11:06.1607287+00:00",
                    "resourceType": "Microsoft.Compute/virtualMachines"
                }
            }
        }
    }
}
::ffff:192.168.176.3 - - [02/Jan/2024:15:01:12 +0000] "POST / HTTP/1.1" 200 4216 "-" "Go-http-client/1.1"
```

### Call Auth0 protected API path

First rename `.env.example` to `.env` and fill in your Auth0 audience and domain credentials. Next follow [this guide](https://auth0.com/docs/secure/tokens/access-tokens/get-access-tokens) to get a JWT access token. This token can now be used for the endpoint `api`.

```bash
curl  \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer YOUR_JWT" \
	-v 'http://localhost:8080/api' -d '
{                                                                                                            
  "schemaId": "Microsoft.Insights/activityLogs",              
  "data": {                                                                      
    "status": "Activated",                                                
    "context": {                                                           
      "activityLog": {                                                    
        "channels": "Admin, Operation",                                    
        "correlationId": "a1be61fd-37ur-ba05-b827-cb874708babf",
        "eventSource": "ResourceHealth",
        "eventTimestamp": "2018-09-04T23:09:03.343+00:00",
        "eventDataId": "2b37e2d0-7bda-4de7-ur8c6-1447d02265b2",                                                                                      
        "level": "Informational",               
        "operationName": "Microsoft.Resourcehealth/healthevent/Activated/action",                                                                    
        "operationId": "2b37e2d0-7bda-489f-81c6-1447d02265b2",
        "properties": {                                                       
          "title": "Virtual Machine health status changed to unavailable",
          "details": "Virtual machine has experienced an unexpected event",
          "currentHealthStatus": "Unavailable",                    
          "previousHealthStatus": "Available",             
          "type": "Downtime",
          "cause": "PlatformInitiated"                    
          },
        "resourceId": "/subscriptions/<subscription Id>/resourceGroups/<resource group>/providers/Microsoft.Compute/virtualMachines/<resource name>",
        "resourceGroupName": "<resource group>",
        "resourceProviderName": "Microsoft.Resourcehealth/healthevent/action",
        "status": "Active",
        "subscriptionId": "<subscription Id>",
        "submissionTimestamp": "2018-09-04T23:11:06.1607287+00:00",
        "resourceType": "Microsoft.Compute/virtualMachines"
      }
    }
  }
}
'
*   Trying 127.0.0.1:8080...
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /api HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.87.0
> Accept: */*
> Content-Type: application/json
> Authorization: Bearer YOUR_JWT
> Content-Length: 2096
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< X-Webhookd-Time-To-Dispatch: 412.375µs
< X-Webhookd-Time-To-Process: 551.375µs
< X-Webhookd-Time-To-Receive: 138.334µs
< X-Webhookd-Time-To-Transform: 375ns
< Date: Tue, 02 Jan 2024 15:35:57 GMT
< Content-Length: 0
<
* Connection #0 to host localhost left intact
```

The payload can now be seen in the logs of `webhook-router`.

```bash
❯ docker logs -f webhook-router
{"time":"2024-01-02T15:45:59.570269508Z","level":"INFO","msg":"Server listening on http://localhost:8080"}
{"time":"2024-01-02T15:46:03.522909426Z","level":"INFO","msg":"Webhook request received on path /api"}
2024/01/02 15:46:03
{
  "schemaId": "Microsoft.Insights/activityLogs",
  "data": {
    "status": "Activated",
    "context": {
      "activityLog": {
        "channels": "Admin, Operation",
        "correlationId": "a1be61fd-37ur-ba05-b827-cb874708babf",
        "eventSource": "ResourceHealth",
        "eventTimestamp": "2018-09-04T23:09:03.343+00:00",
        "eventDataId": "2b37e2d0-7bda-4de7-ur8c6-1447d02265b2",
        "level": "Informational",
        "operationName": "Microsoft.Resourcehealth/healthevent/Activated/action",
        "operationId": "2b37e2d0-7bda-489f-81c6-1447d02265b2",
        "properties": {
          "title": "Virtual Machine health status changed to unavailable",
          "details": "Virtual machine has experienced an unexpected event",
          "currentHealthStatus": "Unavailable",
          "previousHealthStatus": "Available",
          "type": "Downtime",
          "cause": "PlatformInitiated"
          },
        "resourceId": "/subscriptions/<subscription Id>/resourceGroups/<resource group>/providers/Microsoft.Compute/virtualMachines/<resource name>",
        "resourceGroupName": "<resource group>",
        "resourceProviderName": "Microsoft.Resourcehealth/healthevent/action",
        "status": "Active",
        "subscriptionId": "<subscription Id>",
        "submissionTimestamp": "2018-09-04T23:11:06.1607287+00:00",
        "resourceType": "Microsoft.Compute/virtualMachines"
      }
    }
  }
}
```

## Config files

Config files for `webhook-router` are YAML files consisting of four top-level sections. An [example config file](config.yaml) is included with this repository. The five top-level sections are:

### receivers

```yaml
    receivers:
      passthrough: "passthrough://"
```

The `receivers` section is a dictionary of "named" receiver configuations. This allows the actual [webhook configurations (described below)](#webhooks) to signal their respective receivers using the dictionary "name" as a simple short-hand.

### transformations

```yaml
    transformations:
      passthrough: "passthrough://"
      azure-maintenance: "azure-maintenance://"
```

The `transformations` section is a dictionary of "named" tranformation configuations. This allows the actual [webhook configurations (described below)](#webhooks) to signal their respective transformations using the dictionary "name" as a simple short-hand.

### dispatchers

```yaml
    dispatchers:
      log: "log://"
      slack: "slack://?webhook=https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
```

The `dispatchers` section is a dictionary of "named" dispatcher configuations. This allows the actual [webhook configurations (described below)](#webhooks) to signal their respective dispatchers using the dictionary "name" as a simple short-hand.

### webhooks

```
    webhooks:
      - endpoint: "/echo"
        receiver: "passthrough"
        transformations:
          - "passthrough"
        dispatchers:
          - "log"
      - endpoint: "/api"
        receiver: "passthrough"
        transformations:
          - "azure-maintenance"
        dispatchers:
          - "slack"
```

The `webhooks` section is a list of dictionaries. These are the actual webhook endpoints that clients (out there on the internet) will access.

* **endpoint** This is the path that a client will access. It _is_ the webhook URI that clients will send requests to.
* **receiver** The named receiver (defined in the `receivers` section) that the webhook will use to process requests.
* **transformations** An optional list of named transformations (defined in the `transformations` section) that the webhook process the message body with.
* **dispatchers** The list of named dispatchers (defined in the `dispatchers` section) that the webhook will relay a successful request to.

## Components

### Receivers

#### Passthrough

As the name suggests the `Passthrough` receiver applies no checks or validations at all. It will happily accept anything you send to it and relay it on to the dispatcher defined for that webhook. It is defined as a URI string in the form of:

```
passthrough://
```

This receiver exists primarily for debugging purposes.

### Transformations

#### Passthrough

Similar to the [passthrough receiver](#passthrough) this transformation applies no changes and simply passes the webhook from the receiver to the dispatcher.

#### Azure-Maintenance

This transformation accepts Azure Service Health alerts and creates a Slack maintenance alert message from it.

### Dispatchers

#### Log

The `Log` dispatcher will send messages to the logging facility. 

```
log://
```

#### Echo

The `Null` dispatcher will send messages to an echo server. An example can be found in the [docker-compose.yaml](docker-compose.yaml).

```
echo: "echo://webhook-echo:8888"
```

#### Slack

The `Slack` dispatcher sends messages to the incoming webhook URL of a Slack channel.

```
slack: "slack://?webhook=https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
```


## To do

* Write tests across the project
* Log more information, add debug logging
* Investigate using [slog-slack](https://github.com/samber/slog-slack) instead of making HTTP requests
* Think about whether how to serve multiple webhooks 
* JSON schema validation


