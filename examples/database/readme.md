# Example Database

This is an example database with some sample data retrieved from [OpenFlights](https://openflights.org/data)

Execute the following commands (in this directory) to get a local SQL Server running:

```bash
docker build -t sqlrest-data .
docker run -d -p 1433:1433 --name sqlrest-data sqlrest-data
```

## Credentials
username: sa
password: someSecret42!