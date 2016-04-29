package technology.purser.ddns;

import java.util.HashMap;
import com.google.gson.Gson;
import java.util.Scanner;
import java.io.File;

class ConfigFile{
    private String apiKey;
    private String apiEmail;
    
    private String domain;
    
    private String v4Host;
    
    private boolean validateCerts;
    
    private HashMap<String, String> ipv6;
    
    public HashMap<String, String> getIPV6(){
        return ipv6;
    }
    
    public void print(){
        for(String s: ipv6.keySet()){
            System.out.println(s);
        }
    }
    
    public boolean validateCerts(){
        return validateCerts;
    }
    
    public String getApiKey(){
        return apiKey;
    }
    public String getApiEmail(){
        return apiEmail;
    }
    public String getDomain(){
        return domain;
    }
    public String getV4Host(){
        return v4Host;
    }
    
    public static ConfigFile read(String filename){
        try{
            Scanner sc = new Scanner(new File(filename));

            String file = "";
        
            while(sc.hasNextLine()){
                file += sc.nextLine();
            }
            
            Gson gson = new Gson();
            
            ConfigFile cf = gson.fromJson(file, ConfigFile.class);
            
            return cf;
            
        }catch(Exception ex){}
        
        return null;
    }
}
