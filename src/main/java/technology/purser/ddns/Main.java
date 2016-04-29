package technology.purser.ddns;

import java.util.Scanner;
import java.io.File;

import com.google.gson.Gson;


class Main{
    public static void main(String[] args) throws Exception{
        
        String fileName = "config.json";
        
        if(args.length > 0){
            fileName = args[0];
        }
        
        ConfigFile cf = ConfigFile.read(fileName);
        
        ApiClient ac = new ApiClient(cf.getApiEmail(), cf.getApiKey(), cf.validateCerts());

        ac.getIP();
        
        String ipv4 = ac.getIPV4();
        String ipv6 = ac.getIPV6();

        Zone zone = ac.getZones().getZone(cf.getDomain());
        
        RecordList recordList = ac.getRecords(zone.getID());
        
        Record ipv4Record = recordList.getRecord( cf.getV4Host() + "." + cf.getDomain() );
        
        
        // if we have a ipv4 address
        if(ipv4.length() > 0){
            ac.updateIP(ipv4Record, zone, ipv4);
        }
        
        //if we have a ipv6 address
        if(ipv6.length() > 0){
        
            for(String str: cf.getIPV6().keySet()){
                
                Record rec = recordList.getRecord(str + "." + cf.getDomain() );

                String ip = IP.join( ipv6, cf.getIPV6().get(str) );

                ac.updateIP(rec, zone, ip);
            
            }
        }
    }
}
