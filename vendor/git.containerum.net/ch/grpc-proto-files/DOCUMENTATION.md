# Protocol Documentation
<a name="top"/>

## Table of Contents

- [uuid.proto](#uuid.proto)
    - [UUID](#.UUID)
  
  
  
  

- [timestamp.proto](#timestamp.proto)
    - [Timestamp](#google.protobuf.Timestamp)
  
  
  
  

- [duration.proto](#duration.proto)
    - [Duration](#google.protobuf.Duration)
  
  
  
  

- [auth_types.proto](#auth_types.proto)
    - [AccessObject](#.AccessObject)
    - [ResourcesAccess](#.ResourcesAccess)
    - [StoredToken](#.StoredToken)
    - [StoredTokenForUser](#.StoredTokenForUser)
  
  
  
  

- [empty.proto](#empty.proto)
    - [Empty](#google.protobuf.Empty)
  
  
  
  

- [auth.proto](#auth.proto)
    - [CheckTokenRequest](#.CheckTokenRequest)
    - [CheckTokenResponse](#.CheckTokenResponse)
    - [CreateTokenRequest](#.CreateTokenRequest)
    - [CreateTokenResponse](#.CreateTokenResponse)
    - [DeleteTokenRequest](#.DeleteTokenRequest)
    - [DeleteUserTokensRequest](#.DeleteUserTokensRequest)
    - [ExtendTokenRequest](#.ExtendTokenRequest)
    - [ExtendTokenResponse](#.ExtendTokenResponse)
    - [GetUserTokensRequest](#.GetUserTokensRequest)
    - [GetUserTokensResponse](#.GetUserTokensResponse)
    - [UpdateAccessRequest](#.UpdateAccessRequest)
    - [UpdateAccessRequestElement](#.UpdateAccessRequestElement)
  
  
  
    - [Auth](#.Auth)
  

- [Scalar Value Types](#scalar-value-types)



<a name="uuid.proto"/>
<p align="right"><a href="#top">Top</a></p>

## uuid.proto



<a name=".UUID"/>

### UUID
Represents UUID in standart format


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) |  | @inject_tag: binding:&#34;uuid4&#34; |





 

 

 

 



<a name="timestamp.proto"/>
<p align="right"><a href="#top">Top</a></p>

## timestamp.proto



<a name="google.protobuf.Timestamp"/>

### Timestamp
A Timestamp represents a point in time independent of any time zone
or calendar, represented as seconds and fractions of seconds at
nanosecond resolution in UTC Epoch time. It is encoded using the
Proleptic Gregorian Calendar which extends the Gregorian calendar
backwards to year one. It is encoded assuming all minutes are 60
seconds long, i.e. leap seconds are &#34;smeared&#34; so that no leap second
table is needed for interpretation. Range is from
0001-01-01T00:00:00Z to 9999-12-31T23:59:59.999999999Z.
By restricting to that range, we ensure that we can convert to
and from  RFC 3339 date strings.
See [https://www.ietf.org/rfc/rfc3339.txt](https://www.ietf.org/rfc/rfc3339.txt).

# Examples

Example 1: Compute Timestamp from POSIX `time()`.

Timestamp timestamp;
timestamp.set_seconds(time(NULL));
timestamp.set_nanos(0);

Example 2: Compute Timestamp from POSIX `gettimeofday()`.

struct timeval tv;
gettimeofday(&amp;tv, NULL);

Timestamp timestamp;
timestamp.set_seconds(tv.tv_sec);
timestamp.set_nanos(tv.tv_usec * 1000);

Example 3: Compute Timestamp from Win32 `GetSystemTimeAsFileTime()`.

FILETIME ft;
GetSystemTimeAsFileTime(&amp;ft);
UINT64 ticks = (((UINT64)ft.dwHighDateTime) &lt;&lt; 32) | ft.dwLowDateTime;

A Windows tick is 100 nanoseconds. Windows epoch 1601-01-01T00:00:00Z
is 11644473600 seconds before Unix epoch 1970-01-01T00:00:00Z.
Timestamp timestamp;
timestamp.set_seconds((INT64) ((ticks / 10000000) - 11644473600LL));
timestamp.set_nanos((INT32) ((ticks % 10000000) * 100));

Example 4: Compute Timestamp from Java `System.currentTimeMillis()`.

long millis = System.currentTimeMillis();

Timestamp timestamp = Timestamp.newBuilder().setSeconds(millis / 1000)
.setNanos((int) ((millis % 1000) * 1000000)).build();


Example 5: Compute Timestamp from current time in Python.

timestamp = Timestamp()
timestamp.GetCurrentTime()

# JSON Mapping

In JSON format, the Timestamp type is encoded as a string in the
[RFC 3339](https://www.ietf.org/rfc/rfc3339.txt) format. That is, the
format is &#34;{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z&#34;
where {year} is always expressed using four digits while {month}, {day},
{hour}, {min}, and {sec} are zero-padded to two digits each. The fractional
seconds, which can go up to 9 digits (i.e. up to 1 nanosecond resolution),
are optional. The &#34;Z&#34; suffix indicates the timezone (&#34;UTC&#34;); the timezone
is required, though only UTC (as indicated by &#34;Z&#34;) is presently supported.

For example, &#34;2017-01-15T01:30:15.01Z&#34; encodes 15.01 seconds past
01:30 UTC on January 15, 2017.

In JavaScript, one can convert a Date object to this format using the
standard [toISOString()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/toISOString]
method. In Python, a standard `datetime.datetime` object can be converted
to this format using [`strftime`](https://docs.python.org/2/library/time.html#time.strftime)
with the time format spec &#39;%Y-%m-%dT%H:%M:%S.%fZ&#39;. Likewise, in Java, one
can use the Joda Time&#39;s [`ISODateTimeFormat.dateTime()`](
http://joda-time.sourceforge.net/apidocs/org/joda/time/format/ISODateTimeFormat.html#dateTime())
to obtain a formatter capable of generating timestamps in this format.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| seconds | [int64](#int64) |  | Represents seconds of UTC time since Unix epoch 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z inclusive. |
| nanos | [int32](#int32) |  | Non-negative fractions of a second at nanosecond resolution. Negative second values with fractions must still have non-negative nanos values that count forward in time. Must be from 0 to 999,999,999 inclusive. |





 

 

 

 



<a name="duration.proto"/>
<p align="right"><a href="#top">Top</a></p>

## duration.proto



<a name="google.protobuf.Duration"/>

### Duration
A Duration represents a signed, fixed-length span of time represented
as a count of seconds and fractions of seconds at nanosecond
resolution. It is independent of any calendar and concepts like &#34;day&#34;
or &#34;month&#34;. It is related to Timestamp in that the difference between
two Timestamp values is a Duration and it can be added or subtracted
from a Timestamp. Range is approximately &#43;-10,000 years.

# Examples

Example 1: Compute Duration from two Timestamps in pseudo code.

Timestamp start = ...;
Timestamp end = ...;
Duration duration = ...;

duration.seconds = end.seconds - start.seconds;
duration.nanos = end.nanos - start.nanos;

if (duration.seconds &lt; 0 &amp;&amp; duration.nanos &gt; 0) {
duration.seconds &#43;= 1;
duration.nanos -= 1000000000;
} else if (durations.seconds &gt; 0 &amp;&amp; duration.nanos &lt; 0) {
duration.seconds -= 1;
duration.nanos &#43;= 1000000000;
}

Example 2: Compute Timestamp from Timestamp &#43; Duration in pseudo code.

Timestamp start = ...;
Duration duration = ...;
Timestamp end = ...;

end.seconds = start.seconds &#43; duration.seconds;
end.nanos = start.nanos &#43; duration.nanos;

if (end.nanos &lt; 0) {
end.seconds -= 1;
end.nanos &#43;= 1000000000;
} else if (end.nanos &gt;= 1000000000) {
end.seconds &#43;= 1;
end.nanos -= 1000000000;
}

Example 3: Compute Duration from datetime.timedelta in Python.

td = datetime.timedelta(days=3, minutes=10)
duration = Duration()
duration.FromTimedelta(td)

# JSON Mapping

In JSON format, the Duration type is encoded as a string rather than an
object, where the string ends in the suffix &#34;s&#34; (indicating seconds) and
is preceded by the number of seconds, with nanoseconds expressed as
fractional seconds. For example, 3 seconds with 0 nanoseconds should be
encoded in JSON format as &#34;3s&#34;, while 3 seconds and 1 nanosecond should
be expressed in JSON format as &#34;3.000000001s&#34;, and 3 seconds and 1
microsecond should be expressed in JSON format as &#34;3.000001s&#34;.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| seconds | [int64](#int64) |  | Signed seconds of the span of time. Must be from -315,576,000,000 to &#43;315,576,000,000 inclusive. Note: these bounds are computed from: 60 sec/min * 60 min/hr * 24 hr/day * 365.25 days/year * 10000 years |
| nanos | [int32](#int32) |  | Signed fractions of a second at nanosecond resolution of the span of time. Durations less than one second are represented with a 0 `seconds` field and a positive or negative `nanos` field. For durations of one second or more, a non-zero value for the `nanos` field must be of the same sign as the `seconds` field. Must be from -999,999,999 to &#43;999,999,999 inclusive. |





 

 

 

 



<a name="auth_types.proto"/>
<p align="right"><a href="#top">Top</a></p>

## auth_types.proto



<a name=".AccessObject"/>

### AccessObject



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| label | [string](#string) |  |  |
| id | [string](#string) |  | @inject_tag: binding:&#34;uuid4&#34; |
| access | [string](#string) |  |  |






<a name=".ResourcesAccess"/>

### ResourcesAccess



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| namespace | [.AccessObject](#..AccessObject) | repeated |  |
| volume | [.AccessObject](#..AccessObject) | repeated |  |






<a name=".StoredToken"/>

### StoredToken



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token_id | [.UUID](#..UUID) |  |  |
| user_agent | [string](#string) |  |  |
| platform | [string](#string) |  |  |
| fingerprint | [string](#string) |  |  |
| user_id | [.UUID](#..UUID) |  |  |
| user_role | [string](#string) |  |  |
| user_namespace | [string](#string) |  |  |
| user_volume | [string](#string) |  |  |
| rw_access | [bool](#bool) |  |  |
| user_ip | [string](#string) |  | @inject_tag: binding:&#34;ip&#34; |
| part_token_id | [.UUID](#..UUID) |  |  |
| created_at | [.google.protobuf.Timestamp](#..google.protobuf.Timestamp) |  |  |
| life_time | [.google.protobuf.Duration](#..google.protobuf.Duration) |  |  |






<a name=".StoredTokenForUser"/>

### StoredTokenForUser



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token_id | [.UUID](#..UUID) |  |  |
| user_agent | [string](#string) |  |  |
| ip | [string](#string) |  | @inject_tag: binding:&#34;ip&#34; |
| created_at | [string](#string) |  |  |





 

 

 

 



<a name="empty.proto"/>
<p align="right"><a href="#top">Top</a></p>

## empty.proto



<a name="google.protobuf.Empty"/>

### Empty
A generic empty message that you can re-use to avoid defining duplicated
empty messages in your APIs. A typical example is to use it as the request
or the response type of an API method. For instance:

service Foo {
rpc Bar(google.protobuf.Empty) returns (google.protobuf.Empty);
}

The JSON representation for `Empty` is empty JSON object `{}`.





 

 

 

 



<a name="auth.proto"/>
<p align="right"><a href="#top">Top</a></p>

## auth.proto



<a name=".CheckTokenRequest"/>

### CheckTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| user_agent | [string](#string) |  |  |
| finger_print | [string](#string) |  |  |
| user_ip | [string](#string) |  | @inject_tag: binding:&#34;ip&#34; |






<a name=".CheckTokenResponse"/>

### CheckTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access | [.ResourcesAccess](#..ResourcesAccess) |  |  |
| user_id | [.UUID](#..UUID) |  |  |
| user_role | [string](#string) |  |  |
| token_id | [.UUID](#..UUID) |  |  |
| part_token_id | [.UUID](#..UUID) |  |  |






<a name=".CreateTokenRequest"/>

### CreateTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_agent | [string](#string) |  |  |
| fingerprint | [string](#string) |  |  |
| user_id | [.UUID](#..UUID) |  |  |
| user_ip | [string](#string) |  | @inject_tag: binding:&#34;ip&#34; |
| user_role | [string](#string) |  |  |
| rw_access | [bool](#bool) |  |  |
| access | [.ResourcesAccess](#..ResourcesAccess) |  |  |
| part_token_id | [.UUID](#..UUID) |  |  |






<a name=".CreateTokenResponse"/>

### CreateTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| refresh_token | [string](#string) |  |  |






<a name=".DeleteTokenRequest"/>

### DeleteTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token_id | [.UUID](#..UUID) |  |  |
| user_id | [.UUID](#..UUID) |  |  |






<a name=".DeleteUserTokensRequest"/>

### DeleteUserTokensRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [.UUID](#..UUID) |  |  |






<a name=".ExtendTokenRequest"/>

### ExtendTokenRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| refresh_token | [string](#string) |  |  |
| fingerprint | [string](#string) |  |  |






<a name=".ExtendTokenResponse"/>

### ExtendTokenResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| refresh_token | [string](#string) |  |  |






<a name=".GetUserTokensRequest"/>

### GetUserTokensRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [.UUID](#..UUID) |  |  |






<a name=".GetUserTokensResponse"/>

### GetUserTokensResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tokens | [.StoredTokenForUser](#..StoredTokenForUser) | repeated |  |






<a name=".UpdateAccessRequest"/>

### UpdateAccessRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| users | [.UpdateAccessRequestElement](#..UpdateAccessRequestElement) | repeated |  |






<a name=".UpdateAccessRequestElement"/>

### UpdateAccessRequestElement



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [.UUID](#..UUID) |  |  |
| access | [.ResourcesAccess](#..ResourcesAccess) |  |  |





 

 

 


<a name=".Auth"/>

### Auth
The Auth API project is an OAuth authentication server that is used to authenticate users.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateToken | [CreateTokenRequest](#CreateTokenRequest) | [CreateTokenResponse](#CreateTokenRequest) |  |
| CheckToken | [CheckTokenRequest](#CheckTokenRequest) | [CheckTokenResponse](#CheckTokenRequest) |  |
| ExtendToken | [ExtendTokenRequest](#ExtendTokenRequest) | [ExtendTokenResponse](#ExtendTokenRequest) |  |
| UpdateAccess | [UpdateAccessRequest](#UpdateAccessRequest) | [google.protobuf.Empty](#UpdateAccessRequest) |  |
| GetUserTokens | [GetUserTokensRequest](#GetUserTokensRequest) | [GetUserTokensResponse](#GetUserTokensRequest) |  |
| DeleteToken | [DeleteTokenRequest](#DeleteTokenRequest) | [google.protobuf.Empty](#DeleteTokenRequest) |  |
| DeleteUserTokens | [DeleteUserTokensRequest](#DeleteUserTokensRequest) | [google.protobuf.Empty](#DeleteUserTokensRequest) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

