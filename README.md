# nfon-api-client
NFON administration portal REST API client

The [NFON AG](https://nfon.com) offers a very extensive REST API with almost the same scope as the portal UI. But using the API is a bit complicated, this small package should help with the usage.

> [!WARNING]
> **Attention:** The use is expressly at your own risk. 

### Requirements

You need to have
1. the main URL to the API (i.e. `https://api09.nfon.com`)
2. your API public and (i.e. `NFON-ABCDEF-123456`)
3. your API secret and (i.e. `gAp3bTxxUev5JkxOcBdeC5Absm7J84jp6mEhJZd3XiLdjzoGSF`)
4. your account name (i.e. `K1234`)

> [!NOTE]
> If you do not have this information, please contact your NFON partner for assistance.

> [!WARNING]
> **Never share your `API secret`!**

## Usage

Install the package:
```bash
bash$ go get github.com/art-pub/nfon-api-client
```

```go
    ...
	apiConfig := nfonapiclient.ApiConfig{
		BaseURL: "https://api09.nfon.com"
		Public: "NFON-ABCDEF-123456"
    		Secret = "gAp3bTxxUev5JkxOcBdeC5Absm7J84jp6mEhJZd3XiLdjzoGSF"
	}
    
	var req = http.Request{
		Method:     http.MethodGet,
		Body:       io.NopCloser(strings.NewReader("")),
		RequestURI: "/api/version", // the account name is neccessary for some endpoints and part of the URL path
	}

	_, s, failed := nfonapiclient.Request(&req, apiConfig, false)

	if !failed {
		println("Status is " + string(s))
	} else {
		println("Something failed. Please check your API parameters")
	}
    ...
```

Executed this:

```json
Status is {"href":"/api/version","links":[],"data":[{"name":"version","value":"1.17.5.0"},{"name":"host","value":"api09.nfon.com"},{"name":"buildTime","value":"2023-09-14 16:38"}]}
```

### Result Parser
The result is JSON-like, but difficult to read or process. So as a little help there is a parser for simple results and for multiple results in this package.

#### Simple Result Parser
For querying the version, the code can be extended from above:
```go
    ...
	_, s, failed := nfonapiclient.Request(&req, apiConfig, false)

	if !failed {
		parsed := nfonapiclient.SingleresultParser(s)
	    println("version is '" + parsed.DataMap["version"] + " - " + parsed.DataMap["buildTime"] + " - " + parsed.DataMap["host"] + "'")
	} else {
		println("Something failed. Please check your API parameters")
	}
    ...
```
If one executes this, then one receives the following output:

```version is '1.17.5.0 - 2023-09-14 16:38 - api09.nfon.com'```

The data result is parsed into a simple `map[string]string` and can be used for further processing.

#### Multi Result Parser

If you query an endpoint with many results, the multiparser helps to continue working with the result. The usage is as simple as with the single parser:
```go
    ...
    

	apiConfig := nfonapiclient.ApiConfig{
		BaseURL: "https://api09.nfon.com"
		Public: "NFON-ABCDEF-123456"
    		Secret = "gAp3bTxxUev5JkxOcBdeC5Absm7J84jp6mEhJZd3XiLdjzoGSF"
	}
    account := "K4711"

    var req = http.Request{
		Method:     http.MethodGet,
		Body:       io.NopCloser(strings.NewReader("")),
		RequestURI: "/api/customers/" + account + "/phone-books?_pagesize=3",
	}

	_, s, failed := nfonapiclient.Request(&req, apiConfig, false)

	if !failed {
		println("API returns " + string(s))
	} else {
		println("Something failed. Please check your API parameters")
	}

	d := nfonapiclient.MultiResultParser(s)
	fmt.Printf("%v\n", d)
	if d.Total > 0 {
		fmt.Printf("First result is: %s: %s", d.Items[1].DataMap["displayName"], d.Items[1].DataMap["displayNumber"])
	}

    // first dataset is linked as d.LinksMap["first"]
    // last  dataset is linked as d.LinksMap["last"]
    // next  dataset is linked as d.LinksMap["next"]
    ...
```

Output:

`First result is: Foo Bar: +49 (4711) 0815`

### Good to know

#### Datasets and Pagination

Endpoints that return more than one record will return a maximum of 100 records on the first request. The result contains the following information:

Href: Path of the current request

Total: Amount of all datasets (not pages!)

Offset: Offset starting with 0

Size: Amount of maximum results in the response. You can set the amount in the request with the parameter `pageSize=XXX` with `XXX` being max. 100.

Query: You filter the results with the additional parameter `_q`, i.e. `/api/customers/K1234/phone-books?_q=SomeName`.

Links: Array of links including the first, the next and the last URL to retrieve all data. See `LinksMap["first"]`, `LinksMap["last"]` and `LinksMap["next"]` in the example above.

> [!IMPORTANT]
> **Please note:** You have to iterate through all those links to retrieve all data. Just repeat with the `next` given `Href` until your current `Href` (path of the current request) matches the `last` entry.

> [!IMPORTANT]
> If the `last` entry is empty, you already have all data in the current response.


## Links

* [Latest NFON API Documentation (zip)](https://cdn.cloudya.com/API_Documentation.zip)
* [PHP client for NFON API](https://github.com/art-pub/nfon-api-client-php)
* [node.js client for NFON API](https://www.npmjs.com/package/nfon)
