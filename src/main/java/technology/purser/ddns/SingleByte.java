package technology.purser.ddns;

public final class SingleByte{
    
    private final short data;
        
    private final boolean[] binArray;
        
    public SingleByte(String hex){
        this.data = Short.parseShort(hex, 16);
        this.binArray = new boolean[8];
        byteToBin();
    }
          
    public SingleByte(short data){
        this.data = data;
        this.binArray = new boolean[8];
        byteToBin();
    }
    
    public SingleByte(boolean[] data){
        
        this.binArray = data;
        
        short tmpData = 0;
        
        short multiplier = 128;

        for(int i=0; i < 8; i++){
            if( binArray[i] ){
                tmpData += multiplier;
            }
            multiplier = (short) ( multiplier / 2 );
        }
        
        this.data = tmpData;
    }
    
    public boolean[] toBinaryArray(){
        return binArray;
    }
    
    public String toHexString(){
    
        short tmpData = (short) ( data + 256 );
        //To prevent leading zeros
        
        String hexString = Integer.toHexString(tmpData).substring(1);
    
        return hexString;
    }
        
    private void byteToBin(){

        short tmpData = data;

        for(int i=0; i<8; i++){
            binArray[i] = false;
        }
    
        for(int i=0; i<255; i++){
            
            if(tmpData > 0){
                iterate();
                tmpData--;
            }else{
                break;
            }
        }
    }

    private void iterate(){
    
        for(int i=7; i >= 0; i--){
            if(binArray[i]){
                binArray[i] = false;
            }else{
                binArray[i] = true;
                break;
            }
        }
    }
}
