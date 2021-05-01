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
package main.java.org.example.chaincode.invocation;

import static java.nio.charset.StandardCharsets.UTF_8;

import java.util.Collection;
import java.util.Random;
import java.util.Scanner;
import java.util.logging.Level;
import java.util.logging.Logger;

import main.java.org.example.client.CAClient;
import main.java.org.example.client.ChannelClient;
import main.java.org.example.client.FabricClient;
import main.java.org.example.user.UserContext;
import main.java.org.example.util.ConsoleColors;
import main.java.org.example.util.GenerateConfig;
import main.java.org.example.util.Util;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.EventHub;
import org.hyperledger.fabric.sdk.Orderer;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.ProposalResponse;

/**
 * 
 * @author Balaji Kadambi
 *
 */

public class QueryChaincode {

	private static final byte[] EXPECTED_EVENT_DATA = "!".getBytes(UTF_8);
	private static final String EXPECTED_EVENT_NAME = "event";

	public static void main(String args[]) {
		try {
			GenerateConfig gc = new GenerateConfig();
			Random rand = new Random();
			int caOrgNum = rand.nextInt(gc.getConfig().getCaOrgs().length);
			int orgNum = rand.nextInt(gc.getConfig().getOrgs().length);
			int ordNum = rand.nextInt(gc.getConfig().getOrderers().size());
			int peerNum = rand.nextInt(gc.getConfig().getOrgs()[orgNum].getMemberPeers().size());

			Util.cleanUp();
			String caUrl = gc.getConfig().getCaOrgs()[caOrgNum].getCaOrgURL();
			CAClient caClient = new CAClient(caUrl, null);
			// Enroll Admin to Org1MSP
			UserContext adminUserContext = new UserContext();
			adminUserContext.setName(gc.getConfig().getAdminInfo().getAdminName());
			adminUserContext.setAffiliation(gc.getConfig().getOrgs()[orgNum].getOrgName());
			adminUserContext.setMspId(gc.getConfig().getOrgs()[orgNum].getOrgMSP());
			caClient.setAdminUserContext(adminUserContext);
			adminUserContext = caClient.enrollAdminUser(gc.getConfig().getAdminInfo().getAdminName(),
					gc.getConfig().getAdminInfo().getAdminPW());

			FabricClient fabClient = new FabricClient(adminUserContext);

			ChannelClient channelClient = fabClient
					.createChannelClient(gc.getConfig().getChannelInfo().getChannelName());
			Channel channel = channelClient.getChannel();
			Peer peer = fabClient.getInstance().newPeer(
					gc.getConfig().getOrgs()[orgNum].getMemberPeers().get(peerNum).getPeerName(),
					gc.getConfig().getOrgs()[orgNum].getMemberPeers().get(peerNum).getPeerURL());
			EventHub eventHub = fabClient.getInstance().newEventHub("eventhub01", "grpc://localhost:7053");
			Orderer orderer = fabClient.getInstance().newOrderer(gc.getConfig().getOrderers().get(ordNum).getPeerName(),
					gc.getConfig().getOrderers().get(ordNum).getPeerURL());
			channel.addPeer(peer);
			channel.addEventHub(eventHub);
			channel.addOrderer(orderer);
			channel.initialize();
			
			Scanner sc = new Scanner(System.in);

			boolean terminated = false;
			int inputType = 0;
			System.out.println();
			System.out.println("Input type? ");
			while (!terminated) {
				if (sc.hasNext()) {
					if (sc.hasNextInt()) {
						inputType = sc.nextInt();
						if (inputType != 1 && inputType != 2) {
							System.out.println("Input is an incorrect number. Must be 1, 2 or 3.");
						} else {
							terminated = true;
						}
						sc.nextLine(); // to consume the '\n' that nextInt() ignores
					}
				}
			}

			if (inputType == 1) {
				System.out.println();
				System.out.println("Prefix to be queried?");
				String[] args1 = { sc.nextLine() };
				Logger.getLogger(QueryChaincode.class.getName()).log(Level.INFO, "Querying for prefix - " + args1[0]);

				Collection<ProposalResponse> responses1Query = channelClient.queryByChainCode("fabcar",
						"queryAnnouncementOnTree", args1);
				
				for (ProposalResponse pres : responses1Query) {
					String stringResponse = new String(pres.getChaincodeActionResponsePayload());
					// System.out.println(stringResponse); //TODO: usar isto se precisar de ver o
					// que está a ser guardado na blockchain
					String[] separatedPaths = stringResponse.split(";");

					if (!stringResponse.equals("")) {
						System.out.println();
						System.out.println("_________________________________________________________________");
						System.out.println();
						System.out.println(String.format("%23s %-49s", "IP Prefix", "| AS Path"));
						System.out.println("_________________________________________________________________");
						System.out.println();
						for (int i = 0; i < separatedPaths.length; i++) {
							System.out.println(String.format("%30s %-60s", ConsoleColors.YELLOW_BOLD + args1[0],
									ConsoleColors.RESET + "| " + ConsoleColors.WHITE_BOLD + separatedPaths[i]
											+ ConsoleColors.RESET));
						}
						System.out.println("_________________________________________________________________");
						System.out.println();
					} else {
						System.out.println("IP prefix doesn't exist in the Blockchain");
					}
				}
			} else if (inputType == 2) {
				Logger.getLogger(QueryChaincode.class.getName()).log(Level.INFO, "Querying for all prefixes");
				Collection<ProposalResponse> responsesQuery = channelClient.queryByChainCode("fabcar",
						"queryAllTreeAnnouncements", null);
				if (responsesQuery.isEmpty()) {
					System.out.println("The blockchain is empty!");
				} else {
					for (ProposalResponse pres : responsesQuery) {
						// 123.0.12.14/12»Q-W-E;\n10.0.0.0/24»A-B-C;A-T-L;A-R-K-O;\n
						String stringResponse = new String(pres.getChaincodeActionResponsePayload());

						if (!stringResponse.equals("")) {
							//System.out.println(stringResponse); //TODO: usar isto se precisar de ver o
							// que está a ser guardado na blockchain
							String[] allPaths = stringResponse.split("\n"); // [123.0.12.14/12»Q-W-E;]
																			// [10.0.0.0/24»A-B-C;A-T-L;A-R-K-O;]
							Logger.getLogger(QueryChaincode.class.getName()).log(Level.INFO,
									"All IP prefixes and paths");
							System.out.println();
							System.out.println("_____________________________________________________________________________________");
							System.out.println();
							System.out.println(String.format("%43s %-49s", "IP Prefix", "| AS Path"));
							System.out.println("_____________________________________________________________________________________");
							System.out.println();
							for (int i = 0; i < allPaths.length; i++) {
								try{
									String[] prefixAndPath = allPaths[i].split("»"); // [123.0.12.14/12 ] [ Q-W-E;]
									String[] separatedPaths = prefixAndPath[1].split(";"); // [Q-W-E]

									for (int j = 0; j < separatedPaths.length; j++) {
										System.out.println(
												String.format("%50s %-60s", ConsoleColors.YELLOW_BOLD + prefixAndPath[0],
														ConsoleColors.RESET + "| " + ConsoleColors.WHITE_BOLD
																+ separatedPaths[j] + ConsoleColors.RESET));
									}
								} catch(ArrayIndexOutOfBoundsException e){
									//System.out.println("Saltou: " + allPaths[i]);
									e.printStackTrace();
								}
							}
							System.out.println("_____________________________________________________________________________________");
							System.out.println();
						} else {
							System.out.println("Blockchain is empty.");
						}
					}
				}
			}

			// Prints blockchain height and block hashes
			System.out.println(channelClient.getChannel().queryBlockchainInfo(peer).getBlockchainInfo().toString());

			sc.close();

		} catch (Exception e) {
			e.printStackTrace();
		} 
	}

}
