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
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.logging.Level;
import java.util.logging.Logger;

import main.java.org.example.client.ChannelClient;
import main.java.org.example.client.FabricClient;
import main.java.org.example.user.UserContext;
import main.java.org.example.util.GenerateConfig;
import main.java.org.example.util.Util;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.Enrollment;
import org.hyperledger.fabric.sdk.Orderer;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.ProposalResponse;
import org.hyperledger.fabric.sdk.TransactionRequest.Type;
import org.hyperledger.fabric.sdk.security.CryptoSuite;

/**
 * 
 * @author Balaji Kadambi
 *
 */

public class DeployInstantiateChaincode {

	public static void main(String[] args) {
		try {

			GenerateConfig gc = new GenerateConfig();

			CryptoSuite cryptoSuite = CryptoSuite.Factory.getCryptoSuite();

			int orgNumber = gc.getConfig().getOrgs().length;
			FabricClient fabClient = null;
			Channel mychannel = null;
			Orderer orderer = null;
			boolean firstTime = true;
			Collection<ProposalResponse> response = null;

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

				if (firstTime) {
					fabClient = new FabricClient(orgAdmin);
					mychannel = fabClient.getInstance().newChannel(gc.config.getChannelInfo().getChannelName());
					orderer = fabClient.getInstance().newOrderer(gc.config.getOrderers().get(0).getPeerName(),
							gc.config.getOrderers().get(0).getPeerURL());
					firstTime = false;
				}
				fabClient.getInstance().setUserContext(orgAdmin);

				int peerNumber = gc.getConfig().getOrgs()[i].getMemberPeers().size();
				List<Peer> orgPeers = new ArrayList<Peer>();
				for (int j = 0; j < peerNumber; j++) {

					Peer peer = fabClient.getInstance().newPeer(
							gc.config.getOrgs()[i].getMemberPeers().get(j).getPeerName(),
							gc.config.getOrgs()[i].getMemberPeers().get(j).getPeerURL());

					mychannel.addOrderer(orderer);
					mychannel.addPeer(peer);
					orgPeers.add(peer);
				}
				mychannel.initialize();

				response = fabClient.deployChainCode(gc.config.getCcInfo().getCcName(),
						gc.config.getCcInfo().getCcPath(), gc.config.getCcInfo().getCcRootDir(),
						Type.GO_LANG.toString(), gc.config.getCcInfo().getCcVersion(), orgPeers);

				for (ProposalResponse res : response) {
					// Logger.getLogger(DeployInstantiateChaincode.class.getName()).log(Level.INFO,
					// 		gc.config.getCcInfo().getCcName() + "- Chain code deployment " + res.getStatus());
				}
			}
			ChannelClient channelClient = new ChannelClient(mychannel.getName(), mychannel, fabClient);

			String[] arguments = { "" };
			// TODO: se não funcionar, tirar o ultimo parâmetro e meter a null (isto é para
			// a parte dos endorsement policies)
			response = channelClient.instantiateChainCode(gc.config.getCcInfo().getCcName(),
					gc.config.getCcInfo().getCcVersion(), gc.config.getCcInfo().getCcPath(), Type.GO_LANG.toString(),
					"init", arguments, null); // "/root/Desktop/chaincodeendorsementpolicyAllMembers.yaml");

			for (ProposalResponse res : response) {
				// Logger.getLogger(DeployInstantiateChaincode.class.getName()).log(Level.INFO,
				// 		gc.config.getCcInfo().getCcName() + "- Chain code instantiation " + res.getStatus());
			}
		} catch (Exception e) {
			e.printStackTrace();
		}
	}

}
