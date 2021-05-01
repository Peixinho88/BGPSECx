package main.java.org.example.thesis.data_structures;

import java.util.ArrayList;

public class ConfigStructure {

	private Organization[] orgs;
	private CAOrg[] caOrgs;
	private ArrayList<Peer> orderers;
	private AdminInfo adminInfo;
	private ChannelInfo channelInfo;
	private ChaincodeInfo ccInfo;

	public ConfigStructure(Organization[] orgs, CAOrg[] caOrgs, ArrayList<Peer> orderers, AdminInfo adminInfo,
			ChannelInfo channelInfo, ChaincodeInfo ccInfo) {
		this.orgs = orgs;
		this.caOrgs = caOrgs;
		this.orderers = orderers;
		this.adminInfo = adminInfo;
		this.channelInfo = channelInfo;
		this.ccInfo = ccInfo;
	}

	public Organization[] getOrgs() {
		return this.orgs;
	}

	public void setOrgs(Organization[] orgs) {
		this.orgs = orgs;
	}

	public CAOrg[] getCaOrgs() {
		return this.caOrgs;
	}

	public void setCaOrgs(CAOrg[] caOrgs) {
		this.caOrgs = caOrgs;
	}

	public ArrayList<Peer> getOrderers() {
		return this.orderers;
	}

	public void setOrderers(ArrayList<Peer> orderers) {
		this.orderers = orderers;
	}

	public AdminInfo getAdminInfo() {
		return this.adminInfo;
	}

	public void setAdminInfo(AdminInfo adminInfo) {
		this.adminInfo = adminInfo;
	}

	public ChannelInfo getChannelInfo() {
		return this.channelInfo;
	}

	public void setChannelInfo(ChannelInfo channelInfo) {
		this.channelInfo = channelInfo;
	}

	public ChaincodeInfo getCcInfo() {
		return this.ccInfo;
	}

	public void setCcInfo(ChaincodeInfo ccInfo) {
		this.ccInfo = ccInfo;
	}
}