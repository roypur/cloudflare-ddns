package technology.purser.ddns;

class Main{
    public static void main(String[] args){
        
        
        /*
            curl -X GET "https://api.cloudflare.com/client/v4/zones/" \                                             
            -H "X-Auth-Email: spam@royolav.net" \
            -H "X-Auth-Key: f3a64b8f9fe09f14e46f92c6f0952901b125a" \
            -H "Content-Type: application/json" | json_reform
        */
        System.out.println(expandIP("a:b::e:f:4:a"));
        
    }
    private static String expandIP(String ipToExpand){
        char[] ip = new char[39];
        
        ipToExpand = ipToExpand + "-";
        
        int offset = 0;
        
        char[] currentBlock = new char[4];
        
        int posInBlock = 0;
        int posInIPArray = 0;
        
        for(int j=0; j<currentBlock.length; j++){
            currentBlock[j] = '-';
        }
        
        for(int i=0; i < ipToExpand.length(); i++){
        
            char lastChar = ipToExpand.charAt(i);
            
            if( ( lastChar == ':' ) || ( lastChar == '-' ) ){
                if(posInBlock > 0){
                    while(posInBlock < 4){
                        ip[posInIPArray] = '0';
                        posInIPArray++;
                        posInBlock++;
                    }
                    for(int j=0; j < currentBlock.length; j++){
                        if(currentBlock[j] == '-'){
                            break;
                        }
                        ip[posInIPArray] = currentBlock[j];
                        posInIPArray++;
                    }
                }else{
                    offset = posInIPArray;
                    
                }
                posInBlock = 0;
                
                if(lastChar == ':'){
                    ip[posInIPArray] = ':';
                }
                
                posInIPArray++;
                
                
                for(int j=0; j<currentBlock.length; j++){
                    currentBlock[j] = '-';
                }
                
            
            }else{
                currentBlock[posInBlock] = ipToExpand.charAt(i);
                posInBlock++;
            }
        }
        
        char[] tmpIP = new char[39];
        
        for(int i=0; i < tmpIP.length; i++){
            tmpIP[i] = '0';
            if(i % 5 == 4){
                tmpIP[i] = ':';
            }
        }

        for(int i=0; i < offset; i++){
            tmpIP[i] = ip[i];
        }

        int j = tmpIP.length - 1;
        
        for(int i = posInIPArray-2; i > offset; i--){
            tmpIP[j] = ip[i];
            j--;
        }
        
        return new String(tmpIP);
    }
    
}
