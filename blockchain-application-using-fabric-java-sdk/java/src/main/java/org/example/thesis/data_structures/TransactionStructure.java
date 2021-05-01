package main.java.org.example.thesis.data_structures;

public class TransactionStructure {

    private String txID;
    private String prefix;
    private String path;
    private boolean verified;

    public TransactionStructure(String txID, String prefix, String path, boolean verified) {
        this.txID = txID;
        this.prefix = prefix;
        this.path = path;
        this.verified = verified;
    }

    public String getPrefix() {
        return prefix;
    }

    public void setPrefix(String prefix) {
        this.prefix = prefix;
    }

    public String getPath() {
        return path;
    }

    public void setPath(String path) {
        this.path = path;
    }

    public boolean isVerified() {
        return verified;
    }

    public void setVerified(boolean verified) {
        this.verified = verified;
    }

    public String getTxID() {
        return txID;
    }

    public void setTxID(String txID) {
        this.txID = txID;
    }

}