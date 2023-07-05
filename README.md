# gurl - go cmdline request-response tool


## Installation

If you do not have Go installed, Install a prebuilt binary from  [releases](https://github.com/redtrib3/gurl/releases/latest). <br>
If you have the GO version >= 1.16  installed:
 
    ```
    $ git clone https://github.com/redtrib3/gurl.git 
    $ cd gurl
    $ go get
    $ go build

    ```

    Or

    ```
    go install github.com/redtrib3/gurl@latest

    ```

## Usage:

Information on usage can be found by using --help/-h flag.
`./gurl --help`

  `Usage gurl [options...]  `
```
      flag               defaults              description
      ----               --------               ----------
      -H            |                |     Specify Header separated by colon (multiple Headers are allowed) 
      -c            |     false      |     prints colored/syntax highlighted response body. 
      -data         |                |     Specify POST data (form-data/JSON) 
      -m            |     GET        |     Specify the request Method (GET, POST, PUT, DELETE, PATCH) 
      -o            |                |     Save response to a file. 
      -pprint       |     false      |     Pretty print JSON response 
      -proxy        |                |     Specify Proxy URI in format -> [protocol]://host[:port]  
      -raw-request  |     false      |     Print request in raw format (with request headers and body) 
      -redirect     |     false      |     Follow redirects (disabled by default)  
      -u            |                |     URL to remote/local endpoint 
      -upload-file  |                |     Upload file to remote endpoint. (default method - PUT)
```
### Example usages:

Send a GET request:\
` gurl -u https://example.com/ `   <br>

* Send a GET request with Header: <br>
` gurl -u https://example.com/test?uname=123 -H "X-Custom-Header 123" ` <br>
* Multiple Headers:<br>
` gurl -u http://example.com -H "Header1:test" -H "Header2:test" ` <br>

Send POST requests : <br>

`gurl -u https://login.com/ -m POST -data "username=test&password=summer123!" ` <br>
`gurl -u https://login123.com/ -m POST -data '{"username":"test", "password": "pass123"}'` <br>

Other examples: <br>

`gurl -u https://transfer.sh/test -upload-file ~/test.txt` <br>

### Raw requests:<br>

![gurl_chal](https://github.com/redtrib3/gurl/assets/68897241/b7e58421-4988-4473-a617-1e4859d48a29)

### Syntax Highlighting:

HTML Syntax Highlight: <br>
![htmlhigh](https://github.com/redtrib3/gurl/assets/68897241/53788d95-7904-427b-a5b0-e7cc0952f040)

<br>JSON Syntax Highlight:<br>
![jsonhigh](https://github.com/redtrib3/gurl/assets/68897241/9feea2cf-d926-4047-9350-508526800a68)
