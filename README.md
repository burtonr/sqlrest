# sqlrest
A simple GoLang RESTful api proxy to a database

# About
The idea came from the ~need~ want to turn full size APIs into a group of serverless functions. It was quickly discovered that the connections to a database take around 5-7 seconds. No good for a serverless function!

The thought is to create an ultra-minimalist API that acts as a proxy in front of a database in order to maintain the connection, as well as handle the pooling while the functions can focus on being functional.

The function would build the required SQL statment to be executed, and send it to SqlRest for execution on the remote DB. SqlRest does nothing more than passing the query (or procedure name) to the connected database server for execution. The response will depend on the route/command. Select/Query would return a json array of column names and row results
Example Results:
```
data: {
  [
      [
        "Column1",
        "Column2"
      ],
      [
        "Result1_1",
        "Result1_2"
      ],
      [
        "Result2_1",
        "Result2_2"
      ]
    ]
}
```
Errors will be returned with an appropriate HTTP response code and a "message" property containing the reason of the error

Example Error Result:
```
{
  "message": "Query must contain at least 1 'SELECT' statement for 'Query' operation"
}
```
Also, errors encountered in the database will include an "error" property with the error being passed directly from the database
Example DB Error Response:
```
{
  "error": {
    "Number": 102,
    "State": 1,
    "Class": 15,
    "Message": "Incorrect syntax near 'Blah:'.",
    "ServerName": "efe87d4ca854",
    "ProcName": "",
    "LineNo": 1
  },
  "message": "Error returned from database"
}
```

# Run with Docker
Pull the image from [Docker Hub](https://hub.docker.com):

    docker pull burtonr/sqlrest

#### Environment variables that need to be set
|Name | Value |
|-----|-------|
|DATABASE_USERNAME  | The username to log in to the SQL Server with |
|DATABASE_PASSWORD  | The password associated with the user |
|DATABASE_SERVER | The IP address, or hostname, of the SQL server. Do not include the instance, or port number, we got that covered for you (assuming 1433 (default)) :) |

Run the image with the following command (replacing the environment variables with your own)

    docker run -d -p 5050:5050 -e DATABASE_USERNAME=sa -e DATABASE_PASSWORD=secretSauce2! -e DATABASE_SERVER=172.17.0.2 --name sqlrest burtonr/sqlrest:latest


# Usage
The API exposes the following endpoints:
* `/connect`
* `v1/procedure`
* `v1/query`
* `v1/insert`
* `v1/delete`
* `v1/update`

_(the `v1` is an example of the API version that will update as breaking changes happen)_

#### Connect
Sending a `GET` request to this endpoint forces the API to attempt to reconnect to the database using the environment variables provided.

There is a process that runs every 2 minutes to ping the database that will reconnect if it fails. Use this endpoint if you don't want to wait for that process to run

If it is already connected, it will return with a success, otherwise it will attempt the connection and return either a `200` or `500`

#### Procedure
_Not yet implemented_

#### Query
Sending a `POST` request to this endpoint to execute a SQL query and get the results back

There are some basic syntax checks. 
* There must be a field `query` in the request
* There must be at least 1 `SELECT` command 
* _told you it was basic..._

The query passed in is sent to the database directly with no modifications. **Note:** SQL connects to the `master` database by default, so be sure to include a `USE` statement
> `USE Database_Name; SELECT 1 FROM Table_Name` 

_or_ use the full object name in the table definition

>`SELECT 1 FROM [Database_Name].[dbo].[Table_Name]`

See above for examples of error and success responses

#### Insert
_Not yet implemented_

#### Delete
_Not yet implemented_

#### Update
Sending a `POST` request to thie endpoint to execute a SQL update

There are some basic syntax checks.
* There must be a field `update` in the request
* There must be at least 1 `UPDATE` command
* There must be at least 1 `WHERE` clause
  * This is for your own protection!

The function that handles executing updates will first create a transaction, then execute the command. If there is an error, or something goes wrong (`panic`), the transaction will roll back. 

The request and command passed in follow the same rules as the `Query` endpoint. Be sure to include the database name

No results are returned with this command. To get the updated values, you will need to `Query` again.


# Security
* Yes, this is the very definition of SQL Injection, and it's intentional
  * This is intended to act like a database, but easier for serverless by using http protocol and managing the connection
  * Treat it like database access and secure the network around it
___

## _Developer Notes_ 
## API
This (mostly) follows the familiar REST practices as well as special handling of requests based on the route.

* Versioning
    * The API routes will be versioned. `localhost:8080/v1/query`
    * The handler files will include a version if they are of a previous version
        * When a new version is developed, the previous file and the exported `func` will have the version appended to it
        * `queryHandler.go` -> `queryHandler.v1.go`
        * `ExecuteQuery()` -> `ExecuteQueryV1()`

## Validation
* There should be some way to know the requested SQL statement is valid SQL syntax before sending it to the database. This will help avoid unnecessary connections and return a useful error.
* Possible future implementation could look for potential SQL injection attacks (maybe)

## _Future_
## Security
* Look in to implementing HMAC security or possibly API keys
  * Something that's easily added to a serverless function, Docker container, or other API, but still secure
* Security may be required to limit access to certain functions (no deletes), or certain databases, possibly even certain tables

## Considerations
* Should the `Update` and `Insert` calls return the modified/created entry?
  * I like that from a uasability perspective so you don't need to make 2 calls
  * On the otherhand, what about performance? If I (caller) don't need the result, why wait for it?
