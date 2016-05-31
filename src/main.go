package main;

import("fmt"
       "net/http"
       "net"
       "encoding/hex"
       "strings"
       "io/ioutil")

const ipEndpoint = "https://icanhazip.com/"

var ip6 net.IP
var ip4 net.IP

func main(){
    var conf config = getConfig()

    var resp *http.Response
    var err error
    
    resp,err = GetV6(ipEndpoint)
    
    if err == nil{
        
        content, err := ioutil.ReadAll(resp.Body)
        if err == nil{        
            ip6 = net.ParseIP( string( trim(content) ) )
        }else{
            fmt.Println("IPv6 address not found!")
        }
        
    }
    
    resp,err = GetV4(ipEndpoint)
    
    if err == nil{
        
        content, err := ioutil.ReadAll(resp.Body)
        if err == nil{
            ip4 = net.ParseIP( string( trim(content) ) )
        }else{
            fmt.Println("IPv4 address not found!")
        }
    }

    records, err := getRecords(conf.ApiEmail, conf.ApiKey, conf.Domain)

    var hasIp6 bool = (ip6 != nil)
    var hasIp4 bool = (ip4 != nil)
    
    var success bool = true
    
    if err == nil{

        for host,mac := range conf.Ipv6{
                
            for _,rec := range records.Result{
                    
                //lower case fqdn
                var fqdn string = strings.ToLower( strings.TrimSpace(host) + "." + strings.TrimSpace(conf.Domain) )
                    
                var recordName string = strings.TrimSpace( strings.ToLower(rec.Name) )
                    
                    
                if (recordName == fqdn) && (rec.Type == "A") && hasIp4{
                    rec.Content = ip4.String()
                    err = update(conf.ApiEmail, conf.ApiKey, rec)
                    
                    if err != nil{
                        fmt.Println("Failed to update A record " + fqdn)
                        success = false
                    }
                    
                }          
                    
                if (rec.Type == "AAAA") && ( recordName == fqdn ) && hasIp6 {
                
                    rec.Content = joinIP(mac)
                    
                    err = update(conf.ApiEmail, conf.ApiKey, rec)

                    if err != nil{
                        fmt.Println("Failed to update AAAA record " + fqdn)
                        success = false
                    }
                }
                    
            }
        }
        if success{
            fmt.Println("ddns-update: SUCCESS")
        }else{
            fmt.Println("ddns-update: FAIL")
        }
    }else{
        fmt.Println("Failed communicating with cloudflare")
    }
}

func joinIP(mac string)(string){
    
    var tmp string = mac
    
    mac = ""
    
    // remove colons and insert fffe in the middle
    for _,v := range tmp{

        if isHex(v){
            mac += string(v)
        }
        
        if len(mac) == 6{
            mac += "fffe"
        }
    }

    macAddr,err := hex.DecodeString(mac)

    if err == nil{
        macAddr[0] = macAddr[0] ^ 2
    }else{
        return "Invalid mac-address"
    }

    var fullIP net.IP = make(net.IP, 16)
    
    // merge mac and ip
    for k,v := range ip6{
        if k<8{
            fullIP[k] = v
        }else{
            fullIP[k] = macAddr[k-8]
        }
    }
    
    return net.IP(fullIP).String()

}

func isHex(r rune)(bool){
    // Upper and lower case HEX ascii values
    return ( ( (r < 58) && (r > 47) ) || ( (r > 64) && (r < 71) ) || ( (r > 96) && (r < 103) ) )
}
