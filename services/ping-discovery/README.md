# Ping Discovery Service

This service will inspect a Google Cloud project for 
[Agones Latency Ping endpoints](https://agones.dev/site/docs/guides/ping-service/), and return one for each region 
that Agones is installed.

The Service will choose an endpoint at random on each request for each region, assuming there are more than one. 
This is to distribute the load amongst clusters.

## API

<table>
    <thead>
        <td>Endpoint</td>
        <td>Input</td>
        <td>Return</td>
        <td>Description</td>
    </thead>
    <tbody>
        <tr>
            <td><code>GET /list</code></td>
            <td> None </td>
            <td>
                <pre>
{
    "asia-east1": {
        "Name": "agones-ping-udp-service",
        "Namespace": "agones-system",
        "Region": "asia-east1",
        "Address": "104.155.211.151",
        "Protocol": "UDP"
    },
    "europe-west1": {
        "Name": "agones-ping-udp-service",
        "Namespace": "agones-system",
        "Region": "europe-west1",
        "Address": "34.22.151.131",
        "Protocol": "UDP"
    },
    "us-central1": {
        "Name": "agones-ping-udp-service",
        "Namespace": "agones-system",
        "Region": "us-central1",
        "Address": "35.227.137.95",
        "Protocol": "UDP"
    }
}
                </pre>
            </td>
            <td>
                Map of region, where the key is the region name, and a singular
                endpoint for the UDP ping service for each region as the value.
            </td>
        </tr>
    </tbody>
</table>

## Running locally

When running locally, make sure you have [gcloud](https://cloud.google.com/sdk/docs/install) installed,
and a default project authenticated and configured, so that the binary can determine the project it should be 
scanning for
Ping endpoints.

```shell
go run main.go
```

## Building image

```shell
docker build . -t ping-discovery
```

Note: The docker image will fail locally, since it has no access Google Cloud Application Default Credentials.
