package main.java.org.example.thesis.data_structures;

public class CAOrg {

	private String caOrgName;
	private String caOrgURL;

	public CAOrg(String caOrgName, String caOrgURL) {
		this.caOrgName = caOrgName;
		this.caOrgURL = caOrgURL;
	}

	public String getCaOrgName() {
		return this.caOrgName;
	}

	public void setCaOrgName(String caOrgName) {
		this.caOrgName = caOrgName;
	}

	public String getCaOrgURL() {
		return this.caOrgURL;
	}

	public void setCaOrgURL(String caOrgURL) {
		this.caOrgURL = caOrgURL;
	}

}