package technology.purser.ddns;

import java.io.InputStream;
import java.util.HashMap;

import java.util.Arrays;
import java.net.InetAddress;
import java.net.URL;
import java.io.IOException;
import java.net.URISyntaxException;
import org.apache.http.conn.routing.HttpRoute;
import org.apache.http.client.HttpClient;
import org.apache.http.HttpResponse;
import org.apache.http.message.BasicHttpRequest;
import org.apache.http.impl.client.HttpClientBuilder;
import org.apache.http.impl.client.CloseableHttpClient;
import org.apache.http.conn.routing.HttpRoutePlanner;
import org.apache.http.HttpRequest;
import org.apache.http.protocol.HttpContext;
import org.apache.http.HttpException;
import org.apache.http.HttpHost;
import org.apache.http.HttpEntity;

class SimpleHttpClient{
    
    private InetAddress ip;
    private boolean ssl;
    
    private HashMap<String, String> headers;
    
    public SimpleHttpClient(InetAddress ip , boolean ssl){
        this.ip = ip;
        this.ssl = ssl;
        headers = new HashMap<String, String>();
    }
    
    public void addHeader(String k, String v){
        headers.put(k, v);
    }
    
    
    private class RoutePlanner implements HttpRoutePlanner{
        public HttpRoute determineRoute(HttpHost host, HttpRequest req, HttpContext context) throws HttpException{
            HttpHost hostAddress = new HttpHost(ip);
            HttpRoute route = new HttpRoute(host, null, ssl);
            return route;
        }
    }
    
    public byte[] exec(URL url, String method) throws IOException, URISyntaxException{
    
        System.out.println("Doing request!");
        
        
    
        BasicHttpRequest req = new BasicHttpRequest(method, url.toURI().toString());
        
        for(String s: headers.keySet()){
            req.setHeader(s, headers.get(s));
        }
        
        HttpHost host = new HttpHost(ip, url.getPort() );

        HttpClientBuilder builder = HttpClientBuilder.create();
        
        builder.setRoutePlanner(new RoutePlanner());
        
        CloseableHttpClient client = builder.build();

        HttpResponse resp = client.execute(host, req);
        
        HttpEntity entity = resp.getEntity();
        
        
                  // set default max to 50 bytes
        byte[] content = new byte[50];
        
        
        InputStream is = entity.getContent();
        
        System.out.println("Connected");
        
        int next = is.read();
        
        int pos = 0;
        
        while(next > 0){
            if( pos >= content.length ){
                // doubling the length
                content = Arrays.copyOf(content, content.length + 50);
            
            }
            content[pos] = (byte)next;
            next = is.read();
            pos++;
        }
        
        content = Arrays.copyOf(content, pos);

        client.close();
        
        System.out.println("Request done!\n");
        
        
        
        return content;
    }
        
}
    
