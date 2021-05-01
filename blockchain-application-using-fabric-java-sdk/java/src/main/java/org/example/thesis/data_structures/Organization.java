package main.java.org.example.thesis.data_structures;

import java.util.ArrayList;

import main.java.org.example.thesis.data_structures.Peer;

public class Organization {

	private String orgName;
	private String orgMSP;
	private String orgUsrBasePath;
	private String orgUsrAdminPK;
	private String orgUsrAdminCert;
	private ArrayList<Peer> memberPeers;

	public Organization(String name, String msp, String usrBasePath, String usrAdminPK, String usrAdminCert,
			String orgCA) {
		this.orgName = name;
		this.orgMSP = msp;
		this.orgUsrBasePath = usrBasePath;
		this.orgUsrAdminPK = usrAdminPK;
		this.orgUsrAdminCert = usrAdminCert;
		this.memberPeers = new ArrayList<Peer>();
	}

	public Organization() {
		this.orgName = new String();
		this.orgMSP = new String();
		this.orgUsrBasePath = new String();
		this.orgUsrAdminPK = new String();
		this.orgUsrAdminCert = new String();
		this.memberPeers = new ArrayList<Peer>();
	}

	public String getOrgName() {
		return this.orgName;
	}

	public void setOrgName(String orgName) {
		this.orgName = orgName;
	}

	public String getOrgMSP() {
		return this.orgMSP;
	}

	public void setOrgMSP(String orgMSP) {
		this.orgMSP = orgMSP;
	}

	public String getOrgUsrBasePath() {
		return this.orgUsrBasePath;
	}

	public void setOrgUsrBasePath(String orgUsrBasePath) {
		this.orgUsrBasePath = orgUsrBasePath;
	}

	public String getOrgUsrAdminPK() {
		return this.orgUsrAdminPK;
	}

	public void setOrgUsrAdminPK(String orgUsrAdminPK) {
		this.orgUsrAdminPK = orgUsrAdminPK;
	}

	public String getOrgUsrAdminCert() {
		return this.orgUsrAdminCert;
	}

	public void setOrgUsrAdminCert(String orgUsrAdminCert) {
		this.orgUsrAdminCert = orgUsrAdminCert;
	}

	public ArrayList<Peer> getMemberPeers() {
		return this.memberPeers;
	}

	public void setMemberPeers(ArrayList<Peer> memberPeers) {
		this.memberPeers = memberPeers;
	}

	public boolean addPeer(Peer peer) {
		return this.memberPeers.add(peer);
	}

}