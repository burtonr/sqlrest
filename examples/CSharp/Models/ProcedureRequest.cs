
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

public class ProcedureRequest
    {
        [JsonProperty("name")]
        public string Name { get; set; }

        [JsonProperty("parameters")]
        public JObject Parameters { get; set; }

        [JsonProperty("executeOnly")]
        public bool ExecuteOnly { get; set; } = true;
    }