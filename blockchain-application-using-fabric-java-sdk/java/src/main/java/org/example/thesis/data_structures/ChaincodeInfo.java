package main.java.org.example.thesis.data_structures;

public class ChaincodeInfo {

	private String ccName;
	private String ccPath;
	private String ccRootDir;
	private String ccVersion;

	public ChaincodeInfo(String name, String path, String rootDir, String version) {
		this.ccName = name;
		this.ccPath = path;
		this.ccRootDir = rootDir;
		this.ccVersion = version;
	}

	public String getCcName() {
		return this.ccName;
	}

	public void setCcName(String ccName) {
		this.ccName = ccName;
	}

	public String getCcPath() {
		return this.ccPath;
	}

	public void setCcPath(String ccPath) {
		this.ccPath = ccPath;
	}

	public String getCcRootDir() {
		return this.ccRootDir;
	}

	public void setCcRootDir(String ccRootDir) {
		this.ccRootDir = ccRootDir;
	}

	public String getCcVersion() {
		return this.ccVersion;
	}

	public void setCcVersion(String ccVersion) {
		this.ccVersion = ccVersion;
	}
}