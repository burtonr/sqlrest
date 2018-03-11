USE Flights
GO

IF NOT EXISTS (SELECT * FROM sys.objects WHERE type = 'P' AND OBJECT_ID = OBJECT_ID('AirlinesByAirport'))
   exec('CREATE PROCEDURE AirlinesByAirport AS BEGIN SET NOCOUNT ON; END')
GO

ALTER PROCEDURE AirlinesByAirport
    @airportId INT
AS

  SELECT DISTINCT 1 AS 'Source', 0 AS 'Destination', a.*
  FROM Routes r
    JOIN Airlines a ON r.AirlineId = a.Id
  WHERE SourceAirportId = @airportId
  UNION
    SELECT DISTINCT 0 AS 'Source', 1 AS 'Destination', a.*
  FROM Routes r
    JOIN Airlines a ON r.AirlineId = a.Id
  WHERE DestinationAirportId = @airportId
GO

IF NOT EXISTS (SELECT * FROM sys.objects WHERE type = 'P' AND OBJECT_ID = OBJECT_ID('AirportsByAirline'))
   exec('CREATE PROCEDURE AirportsByAirline AS BEGIN SET NOCOUNT ON; END')
GO

ALTER PROCEDURE AirportsByAirline
    @airlineId INT
AS

  SELECT DISTINCT 1 AS 'Source', 0 AS 'Destination', ap.*
    FROM Routes r
    JOIN Airlines al ON r.AirlineId = al.Id
    JOIN Airports ap ON r.SourceAirportId = ap.Id
    WHERE al.Id = @airlineId
  UNION
  SELECT DISTINCT 0 AS 'Source', 1 AS 'Destination', ap.*
    FROM Routes r
    JOIN Airlines al ON r.AirlineId = al.Id
    JOIN Airports ap ON r.DestinationAirportId = ap.Id
    WHERE al.Id = @airlineId
GO

