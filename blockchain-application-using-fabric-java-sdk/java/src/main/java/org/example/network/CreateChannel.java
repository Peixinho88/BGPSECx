/****************************************************** 
 *  Copyright 2018 IBM Corporation 
 *  Licensed under the Apache License, Version 2.0 (the "License"); 
 *  you may not use this file except in compliance with the License. 
 *  You may obtain a copy of the License at 
 *  http://www.apache.org/licenses/LICENSE-2.0 
 *  Unless required by applicable law or agreed to in writing, software 
 *  distributed under the License is distributed on an "AS IS" BASIS, 
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. 
 *  See the License for the specific language governing permissions and 
 *  limitations under the License.
 */
package main.java.org.example.network;

import java.io.File;
import java.util.Collection;
import java.util.Iterator;
import java.util.logging.Level;
import java.util.logging.Logger;

import main.java.org.example.client.FabricClient;
import main.java.org.example.user.UserContext;
import main.java.org.example.util.GenerateConfig;
import main.java.org.example.util.Util;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.ChannelConfiguration;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.Orderer;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.security.CryptoSuite;

/**
 * 
 * @author Balaji Kadambi
 *
 */

public class CreateChannel {

	public static void main(String[] args) {
		try {

			GenerateConfig gc = new GenerateConfig();

			CryptoSuite.Factory.getCryptoSuite();
			Util.cleanUp();

			int orgNumber = gc.getConfig().getOrgs().length;
			FabricClient fabClient = null;
			Channel mychannel = null;
			Orderer orderer = null;
			boolean firstTime = true;

			System.out.println("orgNumber " + orgNumber);

			for (int i = 0; i < orgNumber; i++) {

				// Construct Channel
				UserContext orgAdmin = new UserContext();
				File pkFolder = new File(gc.config.getOrgs()[i].getOrgUsrAdminPK());
				File[] pkFiles = pkFolder.listFiles();
				File certFolder = new File(gc.config.getOrgs()[i].getOrgUsrAdminCert());
				File[] certFiles = certFolder.listFiles();
				Enrollment enrollOrgAdmin = Util.getEnrollment(gc.config.getOrgs()[i].getOrgUsrAdminPK(),
						pkFiles[0].getName(), gc.config.getOrgs()[i].getOrgUsrAdminCert(), certFiles[0].getName());
				orgAdmin.setEnrollment(enrollOrgAdmin);
				orgAdmin.setMspId(gc.config.getOrgs()[i].getOrgMSP());
				orgAdmin.setName(gc.config.getAdminInfo().getAdminName());

				// Create a new channel if it hasn't been created before
				if (firstTime) {
					fabClient = new FabricClient(orgAdmin);
					orderer = fabClient.getInstance().newOrderer(gc.config.getOrderers().get(0).getPeerName(),
							gc.config.getOrderers().get(0).getPeerURL());
					ChannelConfiguration channelConfiguration = new ChannelConfiguration(
							new File(gc.config.getChannelInfo().getChannelConfigPath()));
					byte[] channelConfigurationSignatures = fabClient.getInstance()
							.getChannelConfigurationSignature(channelConfiguration, orgAdmin);
					mychannel = fabClient.getInstance().newChannel(gc.config.getChannelInfo().getChannelName(), orderer,
							channelConfiguration, channelConfigurationSignatures);
					firstTime = false;
				}
				fabClient.getInstance().setUserContext(orgAdmin);
				mychannel = fabClient.getInstance().getChannel("mychannel");

				int peerNumber = gc.getConfig().getOrgs()[i].getMemberPeers().size();
				for (int j = 0; j < peerNumber; j++) {

					Peer peer = fabClient.getInstance().newPeer(
							gc.config.getOrgs()[i].getMemberPeers().get(j).getPeerName(),
							gc.config.getOrgs()[i].getMemberPeers().get(j).getPeerURL());

					mychannel.joinPeer(peer);
				}
				mychannel.addOrderer(orderer);
				mychannel.initialize();
			}

			// Logger.getLogger(CreateChannel.class.getName()).log(Level.INFO, "Channel created " + mychannel.getName());
			Collection peers = mychannel.getPeers();
			Iterator peerIter = peers.iterator();
			while (peerIter.hasNext()) {
				Peer pr = (Peer) peerIter.next();
				// Logger.getLogger(CreateChannel.class.getName()).log(Level.INFO, pr.getName() + " at " + pr.getUrl());
			}

		} catch (Exception e) {
			e.printStackTrace();
		}
	}

}
