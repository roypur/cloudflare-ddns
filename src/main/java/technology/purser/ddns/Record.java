package technology.purser.ddns;

class Record{
    private String name;
    private String id;
    private String type;
    private int ttl;
    
    private String content;
    
    public String getName(){
        return name;
    }
    public String getID(){
        return id;
    }
    public String getType(){
        return type;
    }
    
    public int getTTL(){
        return ttl;
    }
    public void setContent(String content){
        this.content = content;
    }
}
