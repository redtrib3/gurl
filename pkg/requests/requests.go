package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
	UrlParser "net/url"  
)

type HeaderList []string

func (h *HeaderList) String() string {
	return strings.Join(*h, ", ")
}

func (h *HeaderList) Set(value string) error {
	*h = append(*h, value)
	return nil
}

const (
    Black        = "\033[38;5;232m"
    Red          = "\033[38;5;196m"   
    Green        = "\033[32m"
    Yellow       = "\033[33m" 
    Cyan         = "\033[38;5;85m" 
    Reset        = "\033[0m"
    GreenBg      = "\033[42m" 
    CyanBg       = "\033[46m"
)

// functions for cmdline colors
func Colorize(text, ansi string) string {
    return ansi + text + Reset
}


type RequestFlags struct {
	URL         string
	RequestType string
	Data        string
	Outfile     string
	Headers     HeaderList
	PrettyPrint bool
	RawRequest  bool
	UploadPath  string
	Colorize     bool
	AutoRedirect bool
	Proxy       string
}

func isJSON(jstring string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(jstring), &js) == nil
}


func UrlFix(url string) string{
    
    urlparts := strings.SplitN(url, "://", 2)
    if len(urlparts) < 2 {
        fmt.Println(Colorize("[gURL] Warning:", Cyan),"Protocol not specified/detected in url, Using HTTP.")
        url = "http://" + url
        return url
    } 
    
    return url
}

func MakeRequest(args RequestFlags) {
	
    var req *http.Request
    var err error
    var reqData *strings.Reader
    var transport *http.Transport
    
    
    url     := args.URL
    reqtype := args.RequestType
    data    := args.Data
    outfile := args.Outfile
    headers := args.Headers
    pprint  := args.PrettyPrint
    rawreq  := args.RawRequest
    uploadpath := args.UploadPath
    colorize := args.Colorize
    autoRedirect := args.AutoRedirect
    //proxyURI := args.Proxy
    
     
    if args.Proxy != "" {
    
        proxy, _ := UrlParser.Parse(args.Proxy)
        transport = &http.Transport{
            DisableCompression: true,
            Proxy: http.ProxyURL(proxy),
        }
        
    } else {
    
        transport = &http.Transport{
            DisableCompression: true,
            Proxy: nil,
        }
    }
    
	// client confs
    client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if autoRedirect {
				return nil // Allow redirects
			} else {
				return http.ErrUseLastResponse // Prevent redirects
			}
		},
	}
    
    switch reqtype {
        case "GET", "HEAD", "DELETE":
            req, err = http.NewRequest(reqtype, url, nil)
    
        case "POST", "PATCH", "PUT":
    
            if uploadpath != "" {
                
                //reading data
                file, err := os.Open(uploadpath)
                if err != nil {
                    fmt.Println(Colorize("[gURL] Error:", Red),"Failed opening file",uploadpath)
                    return
                }
                
                fileContent, err := ioutil.ReadAll(file)
                if err != nil {
                        fmt.Println(Colorize("[gURL] Error:", Red),"Failed reading file content.")
                    return
                }
            
                //fmt.Println(string(fileContent))
                reqData = strings.NewReader(string(fileContent))
    
            } else{
                reqData = strings.NewReader(data)  
            }
            
            req, err = http.NewRequest(reqtype, url, reqData)
                
        
        default:
            fmt.Println(Colorize("[gURL] Unknown request type.", Yellow))
            return
   }
    
    if err != nil{
            fmt.Println(Colorize("[gURL] Error occured: ", Red), err)
    }

    //set headers
    if isJSON(data){
        req.Header.Set("Content-Type","application/json")        
    }else{
        req.Header.Set("Content-Type","application/x-www-form-urlencoded")
    }
    
    req.Header.Set("User-Agent","gurl/1.1")
    req.Header.Set("Accept","*/*")    
    
    if len(headers) > 0 {
    
        // set custom headers if provided.
        for _,header := range headers{
            header = strings.ReplaceAll(header," ","")
            headAndValue := strings.Split(header, ":")
            
            if len(headAndValue) >= 2 && headAndValue[0] != ""{
                req.Header.Set(headAndValue[0], headAndValue[1])
            } else{
                fmt.Println(Colorize("[gURL] Warning:", Cyan),"Invalid header detected, not sent -> ", "\""+headAndValue[0]+"\"")   
                
            }
        }
        
    }
    
    // common response handling
    response, err := client.Do(req)
    if err != nil{
        err := err.Error()
        switch
        {
            case strings.Contains(err, "connection refused"):
                fmt.Println(Colorize("[gURL] Error: ", Red),"Connection refused to URL")
                os.Exit(0)
            
            default:
                fmt.Println(Colorize("[gURL] Error occured: ",Red), err)
                return
        }
    }
    

    defer response.Body.Close()
 
 
    //handle HEAD request 
    if reqtype == "HEAD"{
    
        tempSlice := make([]string, 0, len(response.Header))
        
        fmt.Println(response.Proto, response.Status)
        for header, values := range response.Header {
            tempSlice = append(tempSlice, header+": "+values[0])
    	}
    	
        sort.Strings(tempSlice)
    	for _, value := range tempSlice{
    	    fmt.Println(value)
    	}
	    return
	}
	
	// --raw-request to print the request headers too.
	
	
	if rawreq {
	
        fmt.Printf("\n%s %s Request %s\n",CyanBg, Black,Reset)
        rawRequest := fmt.Sprintf("\n%s %s %s\r\n", Colorize(req.Method,Cyan), req.URL.RequestURI(), req.Proto)
    	for header, values := range req.Header {
    		for _, value := range values {
    			rawRequest += Colorize(header, Green) + ": " + value + "\r\n"
    		}
	    }
    	fmt.Println(rawRequest)
        fmt.Println(data,"\n")
        
        
        fmt.Printf("%s %s Response %s\n\n",GreenBg,Black,Reset)
        rawResponse := fmt.Sprintf("%s %s\r\n", response.Status, response.Proto)
        for header, values := range response.Header {
            for _, value := range values {
                rawResponse += Colorize(header, Green) + ": " + value + "\r\n"
            }         
        } 	
        fmt.Println(rawResponse)    
        
        
	}
    
	
    //extract response   
    body, _ := ioutil.ReadAll(response.Body)

    // pretty print json if required
    if isJSON(string(body)) {
	    
	    //if pprint is mentioned
	    if pprint {
    		var prettyJson bytes.Buffer
    		err := json.Indent(&prettyJson, body, "", "  ")
    		if err != nil {
	    		fmt.Println(Colorize("[gURL] Error indenting JSON.", Red))
    			fmt.Println(string(body))
	    		return
            }
		    if colorize{
		        fmt.Println(JSONhighlight(prettyJson.String()))
		    } else {
		        fmt.Println(prettyJson.String())
		    }
		    
		
	    } else {
	        if colorize {
	            fmt.Println(JSONhighlight(string(body)))
	        } else {
	            fmt.Println(string(body))
	        }
		    
	    }
    } else {
        if colorize{
            fmt.Println(HTMLhighlight(string(body)))  
        } else {
            fmt.Println(string(body))
        }
    	 
    }
        
    // writing outfile
    if outfile != "" {
        reqBytes:= []byte(string(body))
        err := ioutil.WriteFile(outfile,reqBytes,0644)
        if err != nil{
            fmt.Println(Colorize("[gURL] Error writing outfile" + outfile, Red))
            return
        }
        
        fmt.Println(Colorize("[gURL] Info:", Cyan),"Data written to "+outfile)
    }
    
}

