package main.java.org.example.thesis.auxiliary_blockchain_functions;

import java.util.Random;

import main.java.org.example.client.CAClient;
import main.java.org.example.client.ChannelClient;
import main.java.org.example.client.FabricClient;
import main.java.org.example.thesis.data_structures.CAOrg;
import main.java.org.example.thesis.data_structures.ConfigStructure;
import main.java.org.example.user.UserContext;
import main.java.org.example.util.Util;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.EventHub;
import org.hyperledger.fabric.sdk.Orderer;
import org.hyperledger.fabric.sdk.Peer;

//TODO: copy the chaincode deployment and instatitation from the DeployInstatiateChaincode.java to here

/**
 * This class allows the creation of a user context (client and admin user), the
 * channel required for communication and a peer.
 */
public class BlockchainContext {

	FabricClient fabClient;
	Channel channel;
	ChannelClient channelClient;
	UserContext adminUserContext;
	CAClient caClient;
	ConfigStructure cs;
	int ixpOrg;
	int currPeer;

	public BlockchainContext(ConfigStructure cs, int ixpID, int currPeer) {
		this.fabClient = null;
		this.channel = null;
		this.channelClient = null;
		this.adminUserContext = null;
		this.caClient = null;
		this.cs = cs;
		this.ixpOrg = ixpID - 1;
		this.currPeer = currPeer;
	}

	public void createUserContext() {
		try {
			Random rd = new Random();
			CAOrg randOrgCA = cs.getCaOrgs()[rd.nextInt(cs.getCaOrgs().length)];

			Util.cleanUp();
			// String caUrl = randOrgCA.getCaOrgURL();
			String caUrl = cs.getCaOrgs()[ixpOrg].getCaOrgURL();
			this.caClient = new CAClient(caUrl, null);
			// Enroll Admin to Org1MSP
			this.adminUserContext = new UserContext();
			this.adminUserContext.setName(cs.getAdminInfo().getAdminName());
			this.adminUserContext.setAffiliation(cs.getOrgs()[ixpOrg].getOrgName());
			this.adminUserContext.setMspId(cs.getOrgs()[ixpOrg].getOrgMSP());
			caClient.setAdminUserContext(adminUserContext);
			this.adminUserContext = caClient.enrollAdminUser(cs.getAdminInfo().getAdminName(),
					cs.getAdminInfo().getAdminPW());
			// Register user
			String name = "user" + System.currentTimeMillis();
			String eSecret = this.caClient.registerUser(name, cs.getOrgs()[ixpOrg].getOrgName());
			this.adminUserContext = this.caClient.enrollUser(this.adminUserContext, eSecret);

		} catch (Exception e) {
			e.printStackTrace();
		}
	}

	public void createChannel() {

		try {
			//System.out.println("User context name: " + this.adminUserContext.getName());
			this.fabClient = new FabricClient(adminUserContext);
			this.channelClient = fabClient.createChannelClient(cs.getChannelInfo().getChannelName());
			this.channel = channelClient.getChannel();

		} catch (Exception e) {
			e.printStackTrace();
		}

	}

	public void createPeer() {

		try {

			Random rd = new Random();
			int randOrdNum = rd.nextInt(cs.getOrderers().size()); // TODO: MAYBE USE JUST THE FIRST ORDERER

			// Maybe randomize the peers that are being used?
			// Peer peer = fabClient.getInstance().newPeer(Config.ORG1_PEER_0,
			// Config.ORG1_PEER_0_URL);
			Peer peer = fabClient.getInstance().newPeer(
					cs.getOrgs()[ixpOrg].getMemberPeers().get(currPeer).getPeerName(),
					cs.getOrgs()[ixpOrg].getMemberPeers().get(currPeer).getPeerURL());
			EventHub eventHub = fabClient.getInstance().newEventHub("eventhub01", "grpc://localhost:7053");
			// Orderer orderer = fabClient.getInstance().newOrderer(Config.ORDERER_NAME,
			// Config.ORDERER_URL);
			Orderer orderer = fabClient.getInstance().newOrderer(cs.getOrderers().get(randOrdNum).getPeerName(),
					cs.getOrderers().get(randOrdNum).getPeerURL());

			this.channel.addPeer(peer);
			this.channel.addEventHub(eventHub);
			this.channel.addOrderer(orderer);
			this.channel.initialize();

		} catch (Exception e) {
			e.printStackTrace();
		}

	}

	public FabricClient getFabClient() {
		return this.fabClient;
	}

	public Channel getChannel() {
		return this.channel;
	}

	public ChannelClient getChannelClient() {
		return this.channelClient;
	}
}
