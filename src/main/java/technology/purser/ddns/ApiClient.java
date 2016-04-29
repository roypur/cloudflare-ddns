package technology.purser.ddns;

import java.io.InputStream;
import java.util.Arrays;


import java.net.InetAddress;
import java.net.Inet4Address;
import java.net.Inet6Address;
import java.net.URL;

import java.io.IOException;
import java.net.URISyntaxException;

import com.google.gson.Gson;

import technology.purser.http.SimpleHttpClient;
import technology.purser.http.Request;
import technology.purser.http.Response;


class ApiClient{
    private final String endpoint = "https://api.cloudflare.com/client/v4/";
    
    private final String ipLookupURL = "https://icanhazip.com";
    
    private final String apiEmail;
    private final String apiKey;
    
    private String ipv4Address;
    private String ipv6Address;
    
    private boolean validateCerts;

    public ApiClient(String apiEmail, String apiKey, boolean validateCerts){
        this.apiEmail = apiEmail;
        this.apiKey = apiKey;
        ipv4Address = "";
        ipv6Address = "";
    }
    
    
    public String getIPV6(){
        return ipv6Address;
    }
    
    public String getIPV4(){
        return ipv4Address;
    }
    
    public ZoneList getZones() throws Exception, URISyntaxException{
    
        URL url = new URL(endpoint + "zones/");
        
        Gson gson = new Gson();

        SimpleHttpClient client = new SimpleHttpClient();
        
        Request req = new Request(url);
        
        req.putHeader("Content-Type", "application/json");
        req.putHeader("X-Auth-Key", apiKey);
        req.putHeader("X-Auth-Email", apiEmail);
        
        Response resp = client.exec(req);
        
        ZoneList zones = gson.fromJson(resp.getBody(), ZoneList.class);
        
        return zones;
    
    }
    
    public void updateIP(Record r, Zone z, String ip) throws Exception, URISyntaxException{
    
        Gson gson = new Gson();
    
        URL url = new URL(endpoint + "zones/" + z.getID() + "/dns_records/" + r.getID());
    
        r.setContent(ip);
        
        String jsonString = gson.toJson(r);
        
        SimpleHttpClient client = new SimpleHttpClient();
        
        Request req = new Request(url);
        req.setBody(jsonString);
        req.setMethod("PUT"); 
           
        req.putHeader("X-Auth-Key", apiKey);
        req.putHeader("X-Auth-Email", apiEmail);
        
        Response resp = client.exec(req);
        
    }

    public RecordList getRecords(String id) throws Exception, URISyntaxException{
    
        URL url = new URL(endpoint + "zones/" + id + "/dns_records");
        
        Gson gson = new Gson();

        SimpleHttpClient client = new SimpleHttpClient();
        
        Request req = new Request(url);
        
        req.putHeader("X-Auth-Key", apiKey);
        req.putHeader("X-Auth-Email", apiEmail);
        
        Response resp = client.exec(req);
        
        RecordList records = gson.fromJson(resp.getBody(), RecordList.class);
        
        return records;
    
    }
    
    public void getIP()throws Exception{
    
        URL url = new URL(ipLookupURL);
        
        SimpleHttpClient client = new SimpleHttpClient();
        
        Request req = null;
        Response resp = null;
        
        try{
        
            client.useV6();
        
            req = new Request(url);
        
            resp = client.exec(req);
        
            ipv6Address = resp.getBody().trim();
        }catch(Exception e){
            System.out.println(e.getMessage());
        }
        
        try{
            client.useV4();
        
            resp = client.exec(req);
        
            ipv4Address = resp.getBody().trim();
        }catch(Exception e){
            System.out.println(e.getMessage());
        }
        
    }
    
}
