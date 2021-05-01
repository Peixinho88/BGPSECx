package main.java.org.example.thesis.data_structures;

public class AdminInfo {

	private String adminName;
	private String adminPW;

	public AdminInfo(String name, String pw) {
		this.adminName = name;
		this.adminPW = pw;
	}

	public String getAdminName() {
		return this.adminName;
	}

	public void setAdminName(String adminName) {
		this.adminName = adminName;
	}

	public String getAdminPW() {
		return this.adminPW;
	}

	public void setAdminPW(String adminPW) {
		this.adminPW = adminPW;
	}
}