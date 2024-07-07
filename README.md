# Stuber

Stabber is a tool for creating mocks. It allows you to easily generate mocks from YAML and JSON files, as well as retrieve the history of requests to the mock. Another powerful feature is the ability to filter incoming requests. 


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