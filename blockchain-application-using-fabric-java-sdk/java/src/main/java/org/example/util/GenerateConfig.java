package main.java.org.example.util;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Scanner;

import main.java.org.example.thesis.data_structures.AdminInfo;
import main.java.org.example.thesis.data_structures.CAOrg;
import main.java.org.example.thesis.data_structures.ChaincodeInfo;
import main.java.org.example.thesis.data_structures.ChannelInfo;
import main.java.org.example.thesis.data_structures.ConfigStructure;
import main.java.org.example.thesis.data_structures.Organization;
import main.java.org.example.thesis.data_structures.Peer;

public class GenerateConfig {

	public ConfigStructure config;

	public GenerateConfig(String s) {

		System.out.println(s);

		Scanner sc = new Scanner(System.in);

		System.out.println("Input type? 1 - automatic; 2 - manual (not finished)");
		int inputType = 0;
		boolean error = true;
		if (sc.hasNextInt()) {
			while (error) {
				inputType = sc.nextInt();
				if (inputType < 1 || inputType > 2) {
					System.out.println("Input is an incorrect number. Must be 1 or 2");
				} else {
					error = false;
				}
			}
		}

		// Get the number of orgs in the network
		int orgNumber = 0;
		boolean receivedOrgNum = false;
		System.out.println("How many Orgs in the network?");
		while (!receivedOrgNum) {
			if (sc.hasNextInt()) {
				orgNumber = sc.nextInt();
				receivedOrgNum = true;
			} else {
				System.out.println("Input must be a number");
			}
		}

		int[] peerNumber = new int[orgNumber];
		for (int i = 1; i <= orgNumber; i++) {
			// Get the number of peers in each org
			boolean receivedPeerNum = false;
			System.out.println("How many peers in org " + i + "?");
			while (!receivedPeerNum) {
				if (sc.hasNextInt()) {
					peerNumber[i - 1] = sc.nextInt();
					receivedPeerNum = true;
				} else {
					System.out.println("Input must be a number");
				}
			}
		}

		// Get the number of orgs in the network
		int caOrgNumber = 0;
		boolean receivedCAOrgNum = false;
		System.out.println("How many CA Orgs in the network?");
		while (!receivedCAOrgNum) {
			if (sc.hasNextInt()) {
				caOrgNumber = sc.nextInt();
				receivedCAOrgNum = true;
			} else {
				System.out.println("Input must be a number");
			}
		}

		Organization[] orgs = new Organization[orgNumber];
		CAOrg[] caOrgs = new CAOrg[caOrgNumber];

		if (inputType == 1) {
			// Create the Organization structure automatically with the peers already
			// defined inside them
			for (int i = 1; i <= orgNumber; i++) {
				orgs[i - 1] = new Organization();
				orgs[i - 1].setOrgName("org" + i);
				orgs[i - 1].setOrgMSP("Org" + i + "MSP");
				orgs[i - 1].setOrgUsrBasePath("crypto-config" + File.separator + "peerOrganizations" + File.separator
						+ "org" + i + ".example.com" + File.separator + "users" + File.separator + "Admin@org" + i
						+ ".example.com" + File.separator + "msp");
				orgs[i - 1].setOrgUsrAdminPK("crypto-config" + File.separator + "peerOrganizations" + File.separator
						+ "org" + i + ".example.com" + File.separator + "users" + File.separator + "Admin@org" + i
						+ ".example.com" + File.separator + "msp" + File.separator + "keystore");
				orgs[i - 1].setOrgUsrAdminCert("crypto-config" + File.separator + "peerOrganizations" + File.separator
						+ "org" + i + ".example.com" + File.separator + "users" + File.separator + "Admin@org" + i
						+ ".example.com" + File.separator + "msp" + File.separator + "admincerts");

				for (int j = 0; j < peerNumber[i - 1]; j++) {
					Peer p = new Peer(orgs[i - 1]);
					p.setPeerName("peer" + j + "." + orgs[i - 1].getOrgName() + ".example.com");
					p.setPeerURL("grpc://localhost:" + (7051 + ((i - 1) * 1000) + (j * 5)));
					orgs[i - 1].addPeer(p);
				}

			}

			// Generate automatic CA Org info
			for (int i = 1; i <= caOrgNumber; i++) {
				caOrgs[i - 1] = new CAOrg("CA_Org" + i, "http://localhost:" + (7054 + ((i - 1) * 1000)));
			}

			// Generate automatic Orderer info //TODO: put this in a cycle if I want more
			// orderers
			Peer orderer = new Peer("orderer.example.com", "grpc://localhost:7050");
			ArrayList<Peer> orderers = new ArrayList<Peer>();
			orderers.add(orderer);

			// Generate automatic admin info
			AdminInfo adminInfo = new AdminInfo("admin", "adminpw");

			// Generate automatic channel info
			ChannelInfo channelInfo = new ChannelInfo("mychannel", "config/channel.tx");

			// Generate automatic chaincode info
			ChaincodeInfo ccInfo = new ChaincodeInfo("fabcar", "github.com/fabcar", "chaincode", "1");

			// Put all the generated information in the config structure
			this.config = new ConfigStructure(orgs, caOrgs, orderers, adminInfo, channelInfo, ccInfo);

		} else if (inputType == 2) {
			sc.nextLine();

			// Get the Organization user base path
			String orgUsrBasePath = null;
			boolean receivedLine1 = false;
			System.out.println("What is the organization user base path?");
			while (!receivedLine1) {
				if (sc.hasNext()) {
					orgUsrBasePath = sc.nextLine();
					receivedLine1 = true;
				} else {
					System.out.println("Input must be a path");
				}
			}

			// Get the Organization user admin pk
			String orgUsrAdminPK = null;
			boolean receivedLine2 = false;
			System.out.println("What is the organization user admin private keystore?");
			while (!receivedLine2) {
				if (sc.hasNext()) {
					orgUsrAdminPK = sc.nextLine();
					receivedLine2 = true;
				} else {
					System.out.println("Input must be the name of a keystore");
				}
			}

			// Get the Organization user admin certificate
			String orgUsrAdminCert = null;
			boolean receivedLine3 = false;
			System.out.println("What is the organization user admin certificate?");
			while (!receivedLine3) {
				if (sc.hasNext()) {
					orgUsrAdminCert = sc.nextLine();
					receivedLine3 = true;
				} else {
					System.out.println("Input must be the name of a certificate");
				}
			}

			// Get the Organization CA
			String orgCA = null;
			boolean receivedLine4 = false;
			System.out.println("What is the name of the organization CA?");
			while (!receivedLine4) {
				if (sc.hasNext()) {
					orgCA = sc.nextLine();
					receivedLine4 = true;
				} else {
					System.out.println("Input must be the name of a CA");
				}
			}

			// TODO: test this
			for (int i = 1; i <= orgNumber; i++) {
				orgs[i - 1].setOrgName("org" + i);
				orgs[i - 1].setOrgMSP("Org" + i + "MSP");
				orgs[i - 1].setOrgUsrBasePath(orgUsrBasePath.replaceAll("org", "org" + i));
				orgs[i - 1].setOrgUsrAdminPK(orgUsrAdminPK.replaceAll("org", "org" + i));
				orgs[i - 1].setOrgUsrAdminCert(orgUsrAdminCert.replaceAll("org", "org" + i));
			}

		}

		sc.close();
	}

	public GenerateConfig() {

		// Network configuration file path
		String filePath1 = System.getProperty("user.dir") + "/networkConfig.txt";

		BufferedReader br;
		String line;
		int orgNumber = 0;
		int[] orgIDNumber = null;
		int caOrgNumber = 0;
		// int ordNumber = 0;
		int[] peerNumber = null;
		int k = 0;
		try {
			br = new BufferedReader(new FileReader(new File(filePath1)));

			orgNumber = Integer.parseInt(br.readLine().split(":")[1]);
			orgIDNumber = new int[orgNumber];
			while ((line = br.readLine()).split(":")[0].equals("OrgID")) {
				orgIDNumber[k] = Integer.parseInt(line.split(":")[1]); // specific number of the org, which is now
																		// defined in the networkConfig.txt
				k++;
			}
			caOrgNumber = Integer.parseInt(line.split(":")[1]);
			// ordNumber = Integer.parseInt(br.readLine().split(":")[1]);
			peerNumber = new int[orgNumber];
			k = 0;
			while ((line = br.readLine()) != null) {
				peerNumber[k] = Integer.parseInt(line.split(":")[1]);
				k++;
			}
		} catch (FileNotFoundException e) {
			e.printStackTrace();
		} catch (IOException e) {
			e.printStackTrace();
		}

		Organization[] orgs = new Organization[orgNumber];
		CAOrg[] caOrgs = new CAOrg[caOrgNumber];
		ArrayList<Peer> orderers = new ArrayList<Peer>();

		// Create the Organization structure automatically with the peers already
		// defined inside them
		for (int i = 0; i < orgNumber; i++) {
			orgs[i] = new Organization();
			orgs[i].setOrgName("org" + orgIDNumber[i]);
			orgs[i].setOrgMSP("Org" + orgIDNumber[i] + "MSP");
			orgs[i].setOrgUsrBasePath("crypto-config" + File.separator + "peerOrganizations" + File.separator + "org"
					+ orgIDNumber[i] + ".example.com" + File.separator + "users" + File.separator + "Admin@org"
					+ orgIDNumber[i] + ".example.com" + File.separator + "msp");
			orgs[i].setOrgUsrAdminPK("crypto-config" + File.separator + "peerOrganizations" + File.separator + "org"
					+ orgIDNumber[i] + ".example.com" + File.separator + "users" + File.separator + "Admin@org"
					+ orgIDNumber[i] + ".example.com" + File.separator + "msp" + File.separator + "keystore");
			orgs[i].setOrgUsrAdminCert("crypto-config" + File.separator + "peerOrganizations" + File.separator + "org"
					+ orgIDNumber[i] + ".example.com" + File.separator + "users" + File.separator + "Admin@org"
					+ orgIDNumber[i] + ".example.com" + File.separator + "msp" + File.separator + "admincerts");

			for (int j = 0; j < peerNumber[i]; j++) {
				Peer p = new Peer(orgs[i]);
				p.setPeerName("peer" + j + "." + orgs[i].getOrgName() + ".example.com");
				p.setPeerURL("grpc://localhost:" + (7051 + ((orgIDNumber[i] - 1) * 1000) + (j * 5)));
				orgs[i].addPeer(p);
			}
		}

		// Generate automatic CA Org info
		for (int i = 0; i < caOrgNumber; i++) {
			caOrgs[i] = new CAOrg("CA_Org" + orgIDNumber[i],
					"http://localhost:" + (7054 + ((orgIDNumber[i] - 1) * 1000)));
		}

		// Generate automatic Orderer info //TODO: put this in a cycle if I want more
		// orderers
		Peer orderer = new Peer("orderer.example.com", "grpc://localhost:7050"); // TODO: change this to the orderer
																					// machine IP
		orderers.add(orderer);

		// Generate automatic admin info
		AdminInfo adminInfo = new AdminInfo("admin", "adminpw");

		// Generate automatic channel info
		ChannelInfo channelInfo = new ChannelInfo("mychannel", "config/channel.tx");

		// Generate automatic chaincode info
		ChaincodeInfo ccInfo = new ChaincodeInfo("fabcar", "github.com/fabcar", "chaincode", "1");

		// Put all the generated information in the config structure
		this.config = new ConfigStructure(orgs, caOrgs, orderers, adminInfo, channelInfo, ccInfo);

	}

	public ConfigStructure getConfig() {
		return this.config;
	}

}