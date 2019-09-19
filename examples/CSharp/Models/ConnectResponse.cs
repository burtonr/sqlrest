using System.Runtime.Serialization;

[DataContract]
public class ConnectResponse
{
    [DataMember(Name = "message")]
    public string Message { get; set; }
}