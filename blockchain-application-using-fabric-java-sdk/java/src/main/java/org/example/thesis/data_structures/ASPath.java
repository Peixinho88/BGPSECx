package main.java.org.example.thesis.data_structures;

import java.io.Serializable;
import java.util.ArrayList;

public class ASPath implements Serializable {

	/**
	 * Generated serial version for this class
	 */
	private static final long serialVersionUID = -7555577381287587086L;
	String path;
	int pathLength;
	int numComplaints;
	ArrayList<String> complaintList;
	boolean verified;

	public ASPath(String path, int pathLength) {
		this.path = path;
		this.verified = false;
		this.pathLength = pathLength;
		this.numComplaints = 0;
		this.complaintList = new ArrayList<String>();
	}

	public String getPath() {
		return this.path;
	}

	public void setPath(String path) {
		this.path = path;
		this.pathLength = path.split("-").length;
	}

	public boolean getVerified() {
		return this.verified;
	}

	public void setVerified(boolean verified) {
		this.verified = verified;
	}

	public void addASToPath(String asNum) {
		this.path += "-" + asNum;
		this.pathLength++;
	}

	public int getPathLength() {
		return this.pathLength;
	}

	public void setPathLength(int pathLength) {
		this.pathLength = pathLength;
	}

	public int getNumComplaints() {
		return this.numComplaints;
	}

	public void setNumComplaints(int numComplaints) {
		this.numComplaints = numComplaints;
	}

	public ArrayList<String> getComplaintList() {
		return this.complaintList;
	}

	public void setComplaintList(ArrayList<String> complaintList) {
		this.complaintList = complaintList;
	}

	public void addComplaint(String asNum) {
		this.complaintList.add(asNum);
	}

	/**
	 * 
	 * @param currentPath  - saved path on the routing table
	 * @param receivedPath - received path to compare
	 * @return 1 to update the routing table, 0 to leave the value unchanged
	 */
	public static int bestPathEvaluation(ASPath currentPath, ASPath receivedPath) {
		// A minha ideia é fazer a comparação por path length, mas a cada
		// 2 complaints contar como se o path tivesse uma posição a mais

		int currentPathRealLength = currentPath.getPathLength() + (currentPath.getNumComplaints() % 2);
		int receivedPathRealLength = receivedPath.getPathLength() + (receivedPath.getNumComplaints() % 2);

		System.out.println("Path to insert " + receivedPath.getPath());
		System.out.println("Path in the Routing Table" + currentPath.getPath());

		if (!currentPath.verified) {
			return 1;
		} else if (receivedPath.getPath().startsWith(currentPath.getPath())) {
			return 1;
		} else if (currentPathRealLength < receivedPathRealLength) {
			return 0;
		} else if (currentPathRealLength > receivedPathRealLength) {
			return 1;
		} else if (currentPath.getNumComplaints() < receivedPath.getNumComplaints()) {
			return 0;
		} else {
			return 1;
		}
	}

}
