
using Newtonsoft.Json;

public class QueryRequest
    {
        [JsonProperty("query")]
        public string Query { get; set; }
    }