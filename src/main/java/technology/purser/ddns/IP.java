package technology.purser.ddns;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.Arrays;
import java.io.UnsupportedEncodingException;


class IP{

    public static String join(String ip, String mac){
        
        mac = mac.trim();
        
        byte[] hostBytes = new byte[8];
        
        // Not a mac-address
        if( mac.length() < 12 ){
            InetAddress address = null;
            try{
                address = InetAddress.getByName(mac);
            }catch(UnknownHostException e){
                return null;
            }
            
            byte[] fullAddress = address.getAddress();
            
            // Is ipv6
            if(fullAddress.length > 5){
                for(int i=0; i<8; i++){
                    hostBytes[i] = fullAddress[i+8];
                }
            
            }else{
                return null;
            }
        
        
        }else{
        
            char[] tmpHost = mac.toCharArray();
    
            char[] host = new char[16];
        
            // inserting ff:fe in the middle of the host part
            host[6] = 'f';
            host[7] = 'f';
            host[8] = 'f';
            host[9] = 'e';

            int pos = 0;

            for(int i=0; i<tmpHost.length; i++){
            
                if(pos == 6){
                    pos += 4;
                }
            
                char elem = tmpHost[i];
                
                // convert to lower case
                if( (elem > 60) && (elem < 80) ){
                    elem += 32;
                }
            
                // remove colon
                if( elem == ':' ){
                    continue;
                }
            
                host[pos] = elem;
                pos++;
            }

        
            SingleByte firstByte = new SingleByte( new String( Arrays.copyOf(host, 2) ) );
            
            boolean[] firstByteAsBin = firstByte.toBinaryArray();
        
            // Flipping the 7th bit to create the ipv6 address from a mac address
            if(firstByteAsBin[6]){
                firstByteAsBin[6] = false;
            }else{
                firstByteAsBin[6] = true;
            }
        
            // Converting the first byte back to hex
            SingleByte modifiedFirstByte = new SingleByte( firstByteAsBin );
        
            String modifiedByteAsHex = modifiedFirstByte.toHexString();
        
            host[0] = modifiedByteAsHex.charAt(0);
            host[1] = modifiedByteAsHex.charAt(1);
        
            // convert hex to bytes
            for(int i=0; i < 8; i++){
        
                String hexit = host[i*2] + "" + host[(i*2)+1];
                                            // we use Integer.parseInt instead of Byte.parseByte to make sure
                                            // that it overflows when it reaches 127
                hostBytes[i] = (byte) ( Integer.parseInt( hexit , 16 ) );
            }
        }
        // Done creating host-part of address
    
        InetAddress address = null;
        try{
            address = InetAddress.getByName(ip);
        }catch(UnknownHostException e){
            return null;
        }
        
        byte[] content = address.getAddress();
    
        if(content.length < 5){
            return null;
        }

        for(int i=8; i < 16; i++){
            content[i] = hostBytes[i-8];
        }
        
        try{
            address = InetAddress.getByAddress(content);
        }catch(UnknownHostException e){
            return null;
        }
        
        return address.getHostAddress();
    
    }

}
