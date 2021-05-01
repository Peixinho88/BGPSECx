package main.java.org.example.thesis.data_structures;

/**
 * Class that represents a peer at the application level with all the required
 * information in it
 */
public class Peer {

	private String peerName;
	private String peerURL;
	private Organization ownerOrg;

	public Peer(String name, String url, Organization org) {
		this.peerName = name;
		this.peerURL = url;
		this.ownerOrg = org;
	}

	public Peer(Organization org) {
		this.peerName = new String();
		this.peerURL = new String();
		this.ownerOrg = org;
	}

	public Peer(String name, String url) {
		this.peerName = name;
		this.peerURL = url;
		this.ownerOrg = new Organization();
	}

	public Peer() {
		this.peerName = new String();
		this.peerURL = new String();
		this.ownerOrg = new Organization();
	}

	public String getPeerName() {
		return this.peerName;
	}

	public void setPeerName(String peerName) {
		this.peerName = peerName;
	}

	public String getPeerURL() {
		return this.peerURL;
	}

	public void setPeerURL(String peerURL) {
		this.peerURL = peerURL;
	}

	public Organization getOwnerOrg() {
		return this.ownerOrg;
	}

	public void setOwnerOrg(Organization ownerOrg) {
		this.ownerOrg = ownerOrg;
	}
}