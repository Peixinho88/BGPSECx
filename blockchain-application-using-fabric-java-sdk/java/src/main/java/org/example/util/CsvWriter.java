package main.java.org.example.util;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileWriter;
import java.io.IOException;

public class CsvWriter {

    StringBuilder sb;
    FileWriter writer;

    public CsvWriter() {
        this.sb = new StringBuilder();
    }

    public void createFile(String fileName) {
        try {
            this.writer = new FileWriter(new File(fileName), true);
        } catch (FileNotFoundException e) {
            System.out.println(e.getMessage());
        } catch(IOException e) {
            e.printStackTrace();
        }
    }

    public void addHeaders(int numLines) {
        for (int i = 0; i < numLines; i++) {
            this.sb.append("Tx" + (i+1));
            this.sb.append(',');
        }
        sb.append("TotalTime,");
        sb.append("BlockchainEntries,");
        sb.append("RoutingTableEntries,");
        sb.append("MaliciousEntries\n");
        try {
            this.writer.write(sb.toString());
            this.writer.flush();
        } catch (IOException e) {
            e.printStackTrace();
        }
        this.sb = new StringBuilder();
    }

    public synchronized void addValue(String value, boolean finalInsertion) {
        if(!finalInsertion) {
            this.sb.append(value + ",");
        } else {
            this.sb.append(value + "\n");
        }

        try {
            this.writer.write(sb.toString());
            this.writer.flush();
        } catch (IOException e) {
            e.printStackTrace();
        }
        this.sb = new StringBuilder();
    }

    public void closeWriter() {
        try {
            this.writer.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}