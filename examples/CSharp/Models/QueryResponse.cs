using System.Collections.Generic;
using System.Runtime.Serialization;

[DataContract]
public class QueryResponse
{
    [DataMember(Name = "message")]
    public string Message { get; set; }

    [DataMember(Name = "error")]
    public string Error { get; set; }

    [DataMember(Name = "Columns")]
    public List<string> Columns { get; set; }

    [DataMember(Name = "Data")]
    public List<List<string>> Data { get; set; }
}