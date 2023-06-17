package main

import (
    "fmt"
    "flag"
    "net/http"
    "io/ioutil"
    "strings"
    "os"
    "encoding/json"
    "bytes"
)


type headerSlice []string

func (h *headerSlice) String() string {
    return strings.Join(*h,", ")
}

func (h *headerSlice) Set(value string) error {
    *h = append(*h, value)
    return nil
}


func isJson(jstring string) bool {
    var js json.RawMessage
    return json.Unmarshal([]byte(jstring), &js) == nil
}

func Request(url, reqtype, data, outfile string, headers []string, prettyPrint bool){

    var req *http.Request
    var err error
    
    // starting client
    client := &http.Client{
        Transport: &http.Transport{
        DisableCompression: true,
        },
    }
    
    if reqtype == "GET" { 
        //creating request
        req, err = http.NewRequest("GET", url, nil)
        
    } else if reqtype == "POST" {
        //creating POST request

        postData := strings.NewReader(data)
        req, err = http.NewRequest("POST", url, postData)

    } else {    
        fmt.Println("[gURL] Unknown request type.")
        os.Exit(0)
    }
    
    if err != nil{
            fmt.Println("[gURL] Error occured: ",err)
    }


    //set headers
    if isJson(data){
        req.Header.Set("Content-Type","application/json")        
    }else{
        req.Header.Set("Content-Type","application/x-www-form-urlencoded")
    }
    
    req.Header.Set("User-Agent","gURL/0.0.1")
    req.Header.Set("Accept","*/*")    
    
    // set custom headers if provided.
    for _,header := range headers{
        header = strings.ReplaceAll(header," ","")
        headAndValue := strings.Split(header, ":")
        
        if len(headAndValue) >= 2 && headAndValue[0] != ""{
            req.Header.Set(headAndValue[0], headAndValue[1])
        } else{
            fmt.Println("[gURL] Warning:","Invalid header detected, not sent -> ", "\""+headAndValue[0]+"\"")   
            
        }
    }
    
    // response handling
    response, err := client.Do(req)
    if err != nil{
        err := err.Error()
        switch
        {
            case strings.Contains(err, "connection refused"):
                fmt.Println("[gURL] Error: ","Connection refused to URL")
                os.Exit(0)
            default:
                fmt.Println("[gURL] Error occured: ", err)
                return
        }
    }
    

    defer response.Body.Close()
 
    //extract response   
    body, _ := ioutil.ReadAll(response.Body)

    // pretty print json if required
    if isJson(string(body)) {
	    
	    if prettyPrint {
    		var prettyJson bytes.Buffer
    		err := json.Indent(&prettyJson, body, "", "  ")
    		if err != nil {
	    		fmt.Println("[gURL] Error indenting JSON.")
    			fmt.Println(string(body))
	    		return
	    	}
		
		fmt.Println(prettyJson.String())
		
	    } else {
		    fmt.Println(string(body))
	    }
    } else {
	    fmt.Println(string(body))
    }
        
    // writing outfile
    if outfile != "" {
        reqBytes:= []byte(string(body))
        err := ioutil.WriteFile(outfile,reqBytes,0644)
        if err != nil{
            fmt.Printf("\n[!] Error writing outfile %s",outfile)
            return
        }
        
        fmt.Println("[gURL] Data written to "+outfile)
    }
    
}

func main(){

    var url, outfile, reqtype, data string
    var headers headerSlice
    var prettyPrint bool
    
    //flags
    flag.StringVar(&url, "u", "", "URL") 
    flag.Var(&headers,"H","Header Seperated by colon \neg: -H 'X-Custom-Header: example' ")
    flag.StringVar(&outfile, "o","","Save response to a file.")

    flag.StringVar(&reqtype,"type","GET", "Specify the Request type (GET/POST/PUT/UPDATE)")
    flag.StringVar(&data, "data", "", "Specify POST Data (form/JSON) ")
    
    flag.BoolVar(&prettyPrint, "pprint", false, "Pretty print JSON Response") 
    
            
    flag.Parse()
    

    switch reqtype{
        case "GET":
        
            if data != ""{
                flag.Usage()
                fmt.Println("\n[gURL] Alert:","Cannot send body with GET, use -type to change request method.")
                return
            }
            
            Request(url, "GET", data, outfile, headers, prettyPrint)

        case "POST":
            if data == ""{
                fmt.Println("[gURL]  Sending with Empty body.")
            }
            Request(url, "POST", data, outfile, headers, prettyPrint) 
              
        default:
            Request(url,"GET", data, outfile, headers, prettyPrint)
    }       
    
}
