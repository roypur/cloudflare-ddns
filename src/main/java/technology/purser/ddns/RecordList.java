package technology.purser.ddns;
class RecordList{
    private Record[] result;
    public Record getRecord(String name){
        for(Record r: result){
            if(r.getName().equalsIgnoreCase(name.trim())){
                return r;
            }
        }
        return null;
    }
}
