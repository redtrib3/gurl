/*
 -> Github: https://github.com/redtrib3
  gurl or Go-url
*/
package main

import (
	"flag"
	"fmt"
	"github.com/redtrib3/gurl/pkg/requests"
	"os"
)



func main() {
	var options requests.RequestFlags

	flag.StringVar(&options.URL, "u", "", "URL to remote/local endpoint")
	flag.Var(&options.Headers, "H", "Specify Header separated by colon (multiple Headers are allowed)")
	flag.StringVar(&options.Outfile, "o", "", "Save response to a file.")
	flag.StringVar(&options.RequestType, "m", "GET", "Specify the request Method (GET, POST, PUT, DELETE, PATCH)")
	flag.StringVar(&options.Data, "data", "", "Specify POST data (form-data/JSON)")
	flag.BoolVar(&options.PrettyPrint, "pprint", false, "Pretty print JSON response")
	flag.BoolVar(&options.RawRequest, "raw-request", false, "Print request in raw format (with request headers and body)")
	flag.StringVar(&options.UploadPath, "upload-file", "", "Upload file to remote endpoint. (default method - PUT)")
    flag.BoolVar(&options.Colorize,"c",false,"prints colored/syntax highlighted response body.")
    flag.BoolVar(&options.AutoRedirect,"redirect", false, "Follow redirects (disabled by default) ")
    flag.StringVar(&options.Proxy,"proxy","","Specify Proxy URI in format -> [protocol]://host[:port] ")
    flag.BoolVar(&options.IsQuiet, "quiet",false,"Supress explicit Warning/Info messages.")
    
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage %s [options...]  \n", os.Args[0])
        
        fmt.Printf("\n%24v %36v %38v\n",requests.Colorize("flag",requests.Cyan),requests.Colorize("defaults", requests.Cyan),requests.Colorize("description",requests.Cyan))
        fmt.Printf("%10v %22v %25v","----","--------","----------\n")
        flag.VisitAll(func(f *flag.Flag) {
            var value string
            if len(f.Value.String()) >= 9 {
                value = f.Value.String()[0:5] + "..."
            }else {
                value = f.Value.String()
            }
            
            fmt.Fprintf(os.Stderr, "%-5v -%-12v %-5v %-10v %-5v %3v \n","", f.Name, "|", value, "|", f.Usage) // f.Name, f.Value
        })
    }
    
	flag.Parse()
    
    if options.URL == "" {
        flag.Usage()
        fmt.Println(requests.Colorize("\n[gURL] Error:",requests.Red),"flag '-u' required but not specified")
        return
    }
    
    options.URL = requests.UrlFix(options.URL, options.IsQuiet)
    
	switch options.RequestType {
	case "GET":
		if options.Data != "" {
			flag.Usage()
			fmt.Println(requests.Colorize("\n[gURL] Error:", requests.Red), "Cannot send body with GET, use flag -m to change request method.")
			return
		}

		requests.MakeRequest(options)

	case "POST":
		if options.Data == "" && options.UploadPath == "" {
			requests.PrintQuiet("Sending with empty body.", "warning", options.IsQuiet)
		} 

		requests.MakeRequest(options)

	case "PUT", "HEAD", "DELETE", "PATCH":
		requests.MakeRequest(options)

	default:
		requests.PrintQuiet("Unknown request method type " + "\""+options.RequestType+"\"" + ". Sending GET.","warning",options.IsQuiet)
		options.RequestType = "GET"
		requests.MakeRequest(options)
	}
}

