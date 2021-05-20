# ASCII canvas (Sketch Challenge)

The canvas project fulfills the requirements of the Sketch Challenge as it will be explained below.

## Design choices

- The canvas project has an idempotent RESTful API. That means, if the system receives a duplicate request it will not fail the second time, it will just be ignored.
- Since only the server part of the challenge is being implemented it was required for the server to actually render the ASCII image. However, in a real client-server scenario the rendering part of the process, usually the most expensive one, could be left for the client side. That way the load in the server side would be smaller overall.

## Project structure

- **cmd/canvas**: contains the main executable of the project.
- **config**: contains the general configuration of the project. It follows the [12 factor](https://12factor.net/config) configuration design.
- **devops**: contains the configuration for all related matters that are outside the project scope, but are required for it be run.
- **devops/db/migrations**: contains the migrations required to be run in the DB that wants to be used to store the project information.
- **internal**: contains the Domain Driven Design approach to implement the requirements of the Sketch Challenge.
- **internal/app**: contains the application layer, it uses [Command Query Separation (CQS)](https://en.wikipedia.org/wiki/Command%E2%80%93query_separation) to implement the use cases for the project.
- **internal/domain**: contains the domain layer.
- **internal/infra/ascii**: contains the ASCII renderer used by the project to transform a canvas into an ASCII representation of it.
- **internal/infra/http**: contains the HTTP handlers for the endpoints that will use the command and query handlers from the application layer.
- **internal/infra/sql**: contains the SQL repositories used to store the canvas information.

## How to run the project

The following dependencies are required for the project to be executed. `make`, `docker`, `docker-compose`, and `go`.

### Start the DB

First step to run the project binary you need to start the DB and run the migrations stored in the `devops/db/migrations` folder. In order to do that you need to run:

```bash
make start-infra
```

### Run the binary

Now that the DB is up and running you can start the execution of the project binary:

```bash
make run
```

### Stop the DB

You can stop the DB without losing its contents with the following:
```bash
make stop-infra
```

Alternatively, you can also remove all the content of the DB with:
```bash
make remove-infra
```

## How to test the project

Besides the mandatory dependencies to run the project we will need the following extra dependencies to run test and run linting checks in the project. `golangci-lint`, and `moq`. They can be installed with the following:

```bash
make tools
```

To run the linting checks it is similar to the previous command:

```bash
make lint
```

### Unit test

In order to run the unit test you just need the following command:

```bash
make test
```

### Integration test

Before running the integration test you need the DB to be up and running. Please check the `Start the DB` section above.

```bash
make test-integration
```

## Configuration

As mentioned above the project follows the [12 factor](https://12factor.net/config) configuration design. Therefore, it uses environment variables in order to configure several aspects of the project:

* `CANVAS_HEIGHT`: sets the default height of the canvases created.
* `CANVAS_WIDTH`: sets the default width of the canvases created.
* `HOST`: sets the host name to listen for HTTP requests.
* `PORT`: sets the port to listen for HTTP requests.
* `DB_URL`: sets the connection URL for the DB (expects it to be PostgreSQL).
* `DB_SSL_MODE`: sets the SSL mode for the DB connection.
* `DB_BINARY_PARAMETERS`: sets the binary parameters for the DB connection.

## Usage

### Create a new canvas

By sending a POST request to the `/canvas` endpoint with a body containing a JSON with the following structure:
```json
{
  "id": "02d1170b-67ce-4d19-ae99-acc9ef03c808"
}
```

#### Example

```bash
$ curl -i -X POST "http://localhost:8080/canvas" -d '{"id":"02d1170b-67ce-4d19-ae99-acc9ef03c808"}'
HTTP/1.1 201 Created
Location: http://localhost:8080/canvas/02d1170b-67ce-4d19-ae99-acc9ef03c808
Date: Mon, 17 May 2021 13:13:18 GMT
Content-Length: 0
```

### Draw rectangle on an existing canvas

By sending a POST request to the `/canvas/{canvasID}` endpoint with a JSON body with the following structure:

*Note: `{canvasID}` needs to be replaced by the ID of the canvas where you want to draw the rectangle*

```json
{
  "type": "draw_rectangle",
  "rectangle": {
    "id": "2fb51cca-c789-4938-9d66-948c16a4d42f",
    "point": {
      "x": 5,
      "y": 5
    },
    "height": 3,
    "width": 5,
    "filler": "X",
    "outline": "0"
  }
}
```

#### Example

```bash
$ curl -i -X POST "http://localhost:8080/canvas/02d1170b-67ce-4d19-ae99-acc9ef03c808" -d '{"type":"draw_rectangle","rectangle":{"id":"2fb51cca-c789-4938-9d66-948c16a4d42f","point":{"x":5,"y":5},"height":3,"width":5,"filler":"X","outline":"0"}}'
HTTP/1.1 200 OK
Date: Mon, 17 May 2021 13:16:25 GMT
Content-Length: 0
```


### Perform a flood fill operation on an existing canvas

By sending a POST request to the `/canvas/{canvasID}` endpoint with a JSON body with the following structure:

*Note: `{canvasID}` needs to be replaced by the ID of the canvas where you want to perform the flood fill operation*

```json
{
  "type": "add_fill",
  "fill": {
    "id": "2c2daf0d-97b1-4274-a9ca-06d3c7b167cf",
    "point": {
      "x": 0,
      "y": 0
    },
    "filler": "-"
  }
}
```

#### Example

```bash
$ curl -i -X POST "http://localhost:8080/canvas/02d1170b-67ce-4d19-ae99-acc9ef03c808" -d '{"type":"add_fill","fill":{"id":"2c2daf0d-97b1-4274-a9ca-06d3c7b167cf","point":{"x":0,"y":0},"filler":"-"}}'
HTTP/1.1 200 OK
Date: Mon, 17 May 2021 13:19:42 GMT
Content-Length: 0
```

### Render

By sending a GET request to the `/canvas/{canvasID}` you will receive an ASCII rendering of the canvas as a response.

*Note: `{canvasID}` needs to be replaced by the ID of the canvas that you want to render*

#### Example

```bash
$ curl "http://localhost:8080/canvas/02d1170b-67ce-4d19-ae99-acc9ef03c808"
--------------.......-----------
--------------.......-----------
--------------.......-----------
00000000------.......-----------
0      0------.......-----------
0    XXXXX----.......-----------
00000XXXXX----------------------
-----XXXXX----------------------
--------------------------------
--------------------------------
--------------------------------
--------------------------------
```
