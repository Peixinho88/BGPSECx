package main.java.org.example.thesis.auxiliary_blockchain_functions;

public class TransactionProfile {

	private String threadName;
	private BlockchainContext bcContext;
	private String ipPrefix;
	private String path;
	private int executionNum;

	public TransactionProfile(String name, BlockchainContext bcContext, String prefix, String path, int execNum) {
		this.threadName = name;
		this.bcContext = bcContext;
		this.ipPrefix = prefix;
		this.path = path;
		this.executionNum = execNum;
	}

	public String getThreadName() {
		return this.threadName;
	}

	public void setThreadName(String threadName) {
		this.threadName = threadName;
	}

	public BlockchainContext getBcContext() {
		return this.bcContext;
	}

	public void setBcContext(BlockchainContext bcContext) {
		this.bcContext = bcContext;
	}

	public String getIpPrefix() {
		return this.ipPrefix;
	}

	public void setIpPrefix(String ipPrefix) {
		this.ipPrefix = ipPrefix;
	}

	public String getPath() {
		return this.path;
	}

	public void setPath(String path) {
		this.path = path;
	}

	public int getExecutionNum() {
		return this.executionNum;
	}

	public void setExecutionNum(int executionNum) {
		this.executionNum = executionNum;
	}
}