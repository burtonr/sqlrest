# sqlrest
A simple GoLang RESTful api proxy to a database

# Background
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

## _Notes_ 
### API
There should be a separate route for each of the special functions performed on the database. 
Examples could include (but not limited to):
* `/procedure`
* `/query`
* `/insert`
* `/delete`
* `/update`

This follows the familiar REST practices as well as allowing special handling of requests based on the route.

* Versioning
    * The API routes will be versioned. `localhost:8080/v1/query`
    * The handler files will include a version if they are of a previous version
        * When a new version is developed, the previous file and the exported `func` will have the version appended to it
        * `queryHandler.go` -> `queryHandler.v1.go`
        * `ExecuteQuery()` -> `ExecuteQueryV1()`

### Useage 
* How to define the database to be used for the connection/query/request?
  * Would the caller build the database into the requested SQL statement? `SELECT * FROM [SomeDB].[dbo].[TheTable]`
  * Would it be part of the request? `{"database": "SomeDB", "query": "SELECT * FROM TheTable"}`
  * Something else?

### Validation
* There should be some way to know the requested SQL statement is valid SQL syntax before sending it to the database. This will help avoid unnecessary connections and return a useful error.
* Possible future implementation could look for potential SQL injection attacks (maybe)

## _Future_
### Security
* Security may be required to limit access to certain functions (no deletes), or certain databases, possibly even certain tables
