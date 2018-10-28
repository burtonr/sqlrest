using System;
using System.Net;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Runtime.Serialization.Json;
using System.Security.Cryptography;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;

public class SqlRest
{
    private static HttpClient client;
    private static readonly string apiKey = "sqlrestTestKey";
    private static readonly string realm = "testing-func";

    public string Ping()
    {
        client = new HttpClient();
        var serializer = new DataContractJsonSerializer(typeof(ConnectResponse));

        SetSqlRestHeaders(null);

        var response = client.GetStringAsync("http://localhost:5050/ping");

        return response.Result;
    }

    public ConnectResponse Connect()
    {
        client = new HttpClient();
        var serializer = new DataContractJsonSerializer(typeof(ConnectResponse));

        SetSqlRestHeaders(null);

        var resp = client.GetStreamAsync("http://localhost:5050/connect");
        var connectResponse = serializer.ReadObject(resp.Result) as ConnectResponse;

        return connectResponse;
    }

    public async Task<QueryResponse> Query(QueryRequest query)
    {
        client = new HttpClient();
        var serializer = new DataContractJsonSerializer(typeof(QueryResponse));

        var request = new StringContent(JsonConvert.SerializeObject(query), Encoding.UTF8, "application/json");

        SetSqlRestHeaders(request);

        var resp = await client.PostAsync("http://localhost:5050/v1/query", request);
        if (resp.StatusCode != HttpStatusCode.OK)
        {
            Console.WriteLine("None 200 response code");
            return null;
        }

        var response = serializer.ReadObject(await resp.Content.ReadAsStreamAsync()) as QueryResponse;

        return response;
    }

    public async Task<QueryResponse> Procedure(ProcedureRequest proc)
    {
        client = new HttpClient();
        var serializer = new DataContractJsonSerializer(typeof(QueryResponse));

        var request = new StringContent(JsonConvert.SerializeObject(proc), Encoding.UTF8, "application/json");

        SetSqlRestHeaders(request);

        var resp = await client.PostAsync("http://localhost:5050/v1/procedure", request);
        if (resp.StatusCode != HttpStatusCode.OK)
        {
            Console.WriteLine("None 200 response code");
            return null;
        }

        var response = serializer.ReadObject(await resp.Content.ReadAsStreamAsync()) as QueryResponse;

        return response;
    }

    public async void SetSqlRestHeaders(HttpContent content)
    {
        client.DefaultRequestHeaders.Accept.Clear();
        client.DefaultRequestHeaders.Accept.Add(
            new MediaTypeWithQualityHeaderValue("application/json"));
        client.DefaultRequestHeaders.Add("User-Agent", ".NET Foundation Repository Reporter");

        var hashed = "";

        if (content != null)
        {
            var request = await content.ReadAsStringAsync();
            hashed = GetHashString(request);
        }

        var nonce = Guid.NewGuid();

        var timeStamp = DateTimeOffset.Now.ToUnixTimeMilliseconds();

        client.DefaultRequestHeaders.TryAddWithoutValidation("Authorization", $"{realm}:{hashed}:{nonce}:{timeStamp}");
    }

    private static string GetHashString(string requestBody)
    {
        System.Text.ASCIIEncoding encoding = new System.Text.ASCIIEncoding();
        byte[] keyByte = encoding.GetBytes(apiKey);

        var hmacsha256 = new HMACSHA256(keyByte);

        byte[] messageBytes = encoding.GetBytes(requestBody);
        byte[] hashmessage = hmacsha256.ComputeHash(messageBytes);
        return ByteToString(hashmessage);
    }

    private static string ByteToString(byte[] buff)
    {
        string sbinary = "";

        for (int i = 0; i < buff.Length; i++)
        {
            sbinary += buff[i].ToString("x2"); // hex format
        }
        return (sbinary);
    }
}