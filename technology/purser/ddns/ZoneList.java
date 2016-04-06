class ZoneList{
    private Zone[] result;
    public Zone getZone(String name){
        for(Zone z: result){
            if(z.getName().equalsIgnoreCase(name.trim())){
                return z;
            }
        }
        return null;
    }
}
