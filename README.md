# Stuber

Stuber is a powerful tool for creating mocks. It allows you to:
* Easily generate mocks from YAML and JSON files.
* Retrieve the history of requests to the mock.
* Filter incoming requests.

# Getting started
1. Place your YAML file with the description of your requests in the desired directory.
2. Run the following command:

```shell
docker run -d \
    -p 8080:8080 \
    -v $(pwd)/example:/example \
    hruuum/stuber:latest \
    up -f /example/collection.yaml
```

# Configuration
To describe a mock, you must use the `collection` field and fill in the following fields:
* http_method: The type of HTTP method (e.g., GET, POST).
* path: The path for the request.
* status: The status code of the response.


```yaml
collection:
  foo:
    http_method: GET
    path: /api/foo
    status: 200
```

You can also add a request body using one of the following options:
```yaml
collection:
  foo:
    http_method: GET
    path: /api/foo
    body_path: ./example/json/test.json
    status: 200

  bar:
    http_method: POST
    path: /api/bar
    body: |
      {
        "id": 2,
        "country": "Greece",
        "city": "Rodos",
        "cord": [
          245.245,
          542.542
        ],
        "is_good": true
      }
    status: 201
```
