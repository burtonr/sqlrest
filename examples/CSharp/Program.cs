using System;
using System.Collections.Generic;
using System.Net;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Runtime.Serialization;
using System.Runtime.Serialization.Json;
using System.Security.Cryptography;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

namespace apitest
{
    class Program
    {
        static void Main(string[] args)
        {
            var acceptedInputs = "Accepted inputs are: 'ping', 'connect', 'query', and 'procedure'";

            if (args.GetLength(0) < 1 || string.IsNullOrWhiteSpace(args[0]))
            {
                Console.WriteLine("You must provide a method to execute");
                Console.WriteLine(acceptedInputs);
                return;
            }
            switch (args[0])
            {
                case "ping":
                    Ping();
                    break;
                case "connect":
                    Connect();
                    break;
                case "query":
                    Query();
                    break;
                case "procedure":
                    Procedure();
                    break;
                default:
                    Console.WriteLine("Unknown input");
                    Console.WriteLine(acceptedInputs);
                    break;
            }


        }

        static void Ping()
        {
            var sqlRest = new SqlRest();
            var result = sqlRest.Ping();
            Console.WriteLine(result);
        }

        static void Connect()
        {
            var sqlRest = new SqlRest();
            var result = sqlRest.Connect();
            Console.WriteLine("Message: " + result.Message);
        }

        static void Query()
        {
            var sqlRest = new SqlRest();
            var query = new QueryRequest
            {
                Query = "SELECT TOP 3 * FROM Flights.dbo.Airlines"
            };

            var results = sqlRest.Query(query).Result;

            WriteResults(results);
        }

        static void Procedure()
        {
            var sqlRest = new SqlRest();
            var req = new ProcedureRequest
            {
                Name = "Flights.dbo.AirportsByAirline",
                Parameters = JObject.FromObject(new {airlineId = 109}),
                ExecuteOnly = false
            };

            var results = sqlRest.Procedure(req).Result;

            WriteResults(results);
        }

        private static void WriteResults(QueryResponse response)
        {
            if (!string.IsNullOrWhiteSpace(response.Message))
            {
                Console.WriteLine("Error Happened!");
                Console.WriteLine(response.Message);
                Console.WriteLine(response.Error);
                return;
            }

            for (var i = 0; i < response.Data.Count; i++)
            {
                for (var j = 0; j < response.Columns.Count; j++)
                {
                    Console.WriteLine(response.Columns[j] + ": " + response.Data[i][j]);
                }
                Console.WriteLine();
            }
        }
    }
}
