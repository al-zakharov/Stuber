# Stuber

Stuber is a powerful tool for creating mocks. It allows you to:

* Easily generate mocks from YAML and JSON files.
* Retrieve the history of requests to the mock.
* Filter incoming requests.
* Dynamic mock body replacement.

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

or use docker-compose:

```yaml
services:
  stuber:
    image: hruuum/stuber:latest
    ports:
      - "8080:8080"
    volumes:
      - ./example:/example
    command: [ "up", "-f", "/example/collection.yaml" ]
```

# Configuration

To describe a mock, you must use the `collection` field and fill in the following fields:

* `http_method`: The type of HTTP method (e.g., GET, POST).
* `path`: The path for the request.
* `status`: The status code of the response.

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

# Save incoming requests

By default, all incoming requests are saved. However, if you want to filter incoming requests, there are three
mechanisms available. Note that filtering is specified for each request individually.

1. collect_params with query_param: Specify the name of an HTTP query parameter. By adding this parameter to your
   requests, you can filter incoming requests based on its value.

```yaml
collection:
  foo:
    http_method: GET
    path: /api/foo
    body_path: ./example/json/test.json
    status: 200
    collect_params:
      query_param: "x-request-id"
```

2. collect_params with json_path: Specify a JSONPath. The value of the field found using JSONPath in the incoming
   request body will be saved for filtering.

```yaml
collection:
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
    collect_params:
      json_path: "collect_key"
```

3. collect_params with path_param: Specify Pathparam. The value of the field found at the specified path_param will be
   saved for filtering.

```yaml
collection:
  qux:
    http_method: PUT
    path: /api/:id/qux/:userId
    body: |
      {
        "id": 3,
        "country": "Austria",
        "city": "Liezen",
        "cord": [
          511.115,
          312.213
        ],
        "is_good": true
      }
    status: 200
    collect_params:
      path_param: "userId"
```

# Retrieving Incoming Requests

You can retrieve incoming requests using the following options:

1. Retrieve all requests, you will receive array all requests:

```shell
curl http://localhost:8080/income_request/all
```

2. Retrieve the last request, you will receive one last request:

```shell
curl http://localhost:8080/income_request/last
```

3. Retrieve filtered requests:

```shell
curl 'http://localhost:8080/income_request?searchRequestParam=my_param' 
```

To retrieve filtered requests, use the `searchRequestParam` query parameter.
Pass the value you saved using `query_param`, `json_path` or `path_param` in the `searchRequestParam` query parameter to
get an array of all requests with this value.

# Dynamically substitute request body

Stuber supports dynamic request body substitution. To enable this, set the `dynamic_body` parameter to true in the YAML
configuration for the request.
For example:

```yaml
collection:
  quiz:
    http_method: GET
    path: /api/users
    status: 200
    dynamic_body: true
```

To set the response body that your mock should return, use the `/dynamic_body` endpoint. 
In the `body` of this request, provide the content you want the mock to return using the body parameter, and specify the request path that should return this body using the `path` parameter. For example:

```shell
curl -X POST http://localhost:8080/dynamic_body --data '{"path": "/api/users", "body": {"id": 5, "country": "Amsterdam"}}'
```

Now, the `/api/users` request will return the specified body `{"id": 5, "country": "Amsterdam"}`.

# Example YAML Configuration

```yaml
collection:
  foo:
    http_method: GET
    path: /api/foo
    body_path: ./example/json/test.json
    status: 200
    collect_params:
      query_param: "x-request-id"

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
    collect_params:
      json_path: "collect_key"

  baz:
    http_method: DELETE
    path: /api/baz
    status: 204

  qux:
    http_method: PUT
    path: /api/:id/qux/:userId
    body: |
      {
        "id": 3,
        "country": "Austria",
        "city": "Liezen",
        "cord": [
          511.115,
          312.213
        ],
        "is_good": true
      }
    status: 200
    collect_params:
      path_param: "userId"

  quiz:
    http_method: GET
    path: /api/users
    status: 200
    dynamic_body: true
```