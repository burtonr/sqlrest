# Wait for the SQL Server to come up
sleep 5s

# Run the setup script to create the DB and the schema in the DB
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P someSecret42! -d master -i sql/schema.sql

# Run the setup script to generate the stored procedures
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P someSecret42! -d master -i sql/procedures.sql

# Import the data from the csv files
/opt/mssql-tools/bin/bcp Flights.dbo.Airlines in "airlines.csv" -c -t',' -S localhost -U sa -P someSecret42!
/opt/mssql-tools/bin/bcp Flights.dbo.Airports in "airports.csv" -c -t',' -S localhost -U sa -P someSecret42!
/opt/mssql-tools/bin/bcp Flights.dbo.Routes in "routes.csv" -c -t',' -S localhost -U sa -P someSecret42!