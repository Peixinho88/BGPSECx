package main.java.org.example.thesis.auxiliary_blockchain_functions;

import main.java.org.example.config.Config;
import main.java.org.example.thesis.data_structures.PrefixAnnouncement;
import main.java.org.example.thesis.data_structures.TransactionStructure;
import main.java.org.example.thesis.simulators.IXPsimulator;
import main.java.org.example.util.ConsoleColors;
import main.java.org.example.util.CsvWriter;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.Semaphore;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.logging.Level;
import java.util.logging.Logger;

import static java.nio.charset.StandardCharsets.UTF_8;

import org.hyperledger.fabric.sdk.ChaincodeID;
import org.hyperledger.fabric.sdk.ChaincodeResponse.Status;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.exception.ProposalException;
import org.hyperledger.fabric.sdk.ProposalResponse;
import org.hyperledger.fabric.sdk.TransactionProposalRequest;

/**
 * Class to define a thread that takes a request for an insertion on the
 * blockchain and handles it individually until it is completed.
 */
public class RequestThread implements Runnable {

	private final String CHAINCODE_NAME = "fabcar";

	private Thread t;
	private String threadName;
	private BlockchainContext bcContext;
	private CsvWriter writer;
	private Semaphore csvSem;
	private String ipPrefix;
	private String path;
	private int executionNum;
	private ConcurrentLinkedQueue<TransactionStructure> verifyingQueue;
	private AtomicInteger sequenceNumber;

	private static final byte[] EXPECTED_EVENT_DATA = "!".getBytes(UTF_8);
	private static final String EXPECTED_EVENT_NAME = "event";

	public RequestThread(CsvWriter writer, Semaphore csvSem, TransactionProfile txProfile,
			ConcurrentLinkedQueue<TransactionStructure> verifyingQueue, AtomicInteger sequenceNumber) {
		this.writer = writer;
		this.csvSem = csvSem;
		this.verifyingQueue = verifyingQueue;
		this.threadName = txProfile.getThreadName();
		this.bcContext = txProfile.getBcContext();
		this.ipPrefix = txProfile.getIpPrefix();
		this.path = txProfile.getPath();
		this.executionNum = txProfile.getExecutionNum();
		this.sequenceNumber = sequenceNumber;
	}

	public void run() {

		// Copy the code from the IXP to here and see if the execution still
		// behaves the same way. Check if insertions are all correct and in order
		// TODO: it's creating a user context and channel on every iteration (MIGHT BE A
		// PROBLEM). Might have to put this outside the while loop

		// Start timer to count the time it takes to process each individual transaction
		long startTime = System.nanoTime();

		boolean isVerified = false;
		String txID = null;
		ChaincodeID ccid = ChaincodeID.newBuilder().setName(CHAINCODE_NAME).build();
		Map<String, byte[]> tm2 = new HashMap<>();
		
		try {
			
			// THIS PART BELOW IS COMMUNICATING WITH THE CHAINCODE IN GOLANG fabcar.go - DO
			// NOT FORGET!
			// This is where I define the function name and the appropriate arguments (check
			// my fabcar functions)
			TransactionProposalRequest request = bcContext.getFabClient().getInstance().newTransactionProposalRequest();
			// Chaincode name along with peer names and org names are defined in the
			// Config.java file
			
			request.setChaincodeID(ccid);
			request.setFcn("announceVerifiedTreePath");
			System.out.println();
			System.out.println("___________________________________________________");
			System.out.println();
			System.out.println(ConsoleColors.YELLOW + "IP PREFIX: " + ipPrefix + " PATH: " + path + ConsoleColors.RESET);
			System.out.println("___________________________________________________");
			System.out.println();
			String[] arguments = { ipPrefix.trim(), path.trim(), String.valueOf(this.executionNum) }; // IP prefix / AS path
																									// / ExecutionNum
			// System.out.println("PREFIX: " + arguments[0] + " | " + "PATH: " +
			// arguments[1]);
			request.setArgs(arguments);
			request.setProposalWaitTime(120000);

			// ALl of this within the map is not necessary for the request, I think
			tm2.put("HyperLedgerFabric", "TransactionProposalRequest:JavaSDK".getBytes(UTF_8));
			tm2.put("method", "TransactionProposalRequest".getBytes(UTF_8));
			tm2.put("result", ":)".getBytes(UTF_8));
			tm2.put(EXPECTED_EVENT_NAME, EXPECTED_EVENT_DATA); // This line I guess is just an example of how to put
															// extra
															// stuff in the request

			try {
				request.setTransientMap(tm2);
			} catch (InvalidArgumentException e) {
				e.printStackTrace();
			}
			
			Collection<ProposalResponse> responses = bcContext.getChannelClient().sendTransactionProposal(request);
			
			
			for (ProposalResponse res : responses) {
				//System.out.println(res.toString());
				Status status = res.getStatus(); // if it worked, status is "SUCCESS"
				String responseString = new String(res.getChaincodeActionResponsePayload());
				
				isVerified = responseString.charAt(0) == 'V';
				txID = res.getTransactionID();

				//Check to see if second transaction needs to happen
				String[] params = responseString.split(" \\| ");
				String[] paramsFinal = new String[params.length + 1];

				//Copy all the params and the new one (execution number) to the new request params array
				for (int i = 0; i < params.length; i++) {
					paramsFinal[i] = params[i];
				}
				paramsFinal[paramsFinal.length - 1] = String.valueOf(this.executionNum);

				if(paramsFinal.length > 3) {
					TransactionProposalRequest finalRequest = bcContext.getFabClient().getInstance().newTransactionProposalRequest();
					finalRequest.setChaincodeID(ccid);
					finalRequest.setFcn("updateVerifiedPath");
					/*String[] newArgs = new String[params.length-1];
					for (int i = 0; i < params.length; i++) {
						newArgs[i-1] = params[i];
					}*/
					finalRequest.setArgs(paramsFinal);
					finalRequest.setProposalWaitTime(120000);
					try {
						finalRequest.setTransientMap(tm2);
					} catch (InvalidArgumentException e) {
						e.printStackTrace();
					}
					
					try {
						Collection<ProposalResponse> newResponses = bcContext.getChannelClient().sendTransactionProposal(finalRequest);
					
						for (ProposalResponse newRes : newResponses) {
							Status newStatus = newRes.getStatus(); // if it worked status is "SUCCESS"
							Logger.getLogger(IXPsimulator.class.getName()).log(Level.INFO,
							"Invoked updateVerifiedPath on " + CHAINCODE_NAME + ". Status - " + newStatus);
						}

					} catch (ProposalException e) {
						e.printStackTrace();
					} catch (InvalidArgumentException e) {
						e.printStackTrace();
					}
				}

				Logger.getLogger(IXPsimulator.class.getName()).log(Level.INFO,
						"Invoked announceVerifiedTreePath on " + CHAINCODE_NAME + ". Status - " + status);
			}

		} catch (ProposalException e) {
			e.printStackTrace();
			TransactionStructure txStruct = new TransactionStructure(txID, ipPrefix, path, isVerified);
			verifyingQueue.add(txStruct);
		} catch (InvalidArgumentException e) {
			e.printStackTrace();
			TransactionStructure txStruct = new TransactionStructure(txID, ipPrefix, path, isVerified);
			verifyingQueue.add(txStruct);
		} catch (NullPointerException e) {
			e.printStackTrace();
			TransactionStructure txStruct = new TransactionStructure(txID, ipPrefix, path, isVerified);
			verifyingQueue.add(txStruct);
		} catch (RuntimeException e) {
			e.printStackTrace();
			TransactionStructure txStruct = new TransactionStructure(txID, ipPrefix, path, isVerified);
			verifyingQueue.add(txStruct);
		} /*catch (InterruptedException e) {
			e.printStackTrace();
			TransactionStructure txStruct = new TransactionStructure(txID, ipPrefix, path, isVerified);
			verifyingQueue.add(txStruct);
		}*/

		this.sequenceNumber.getAndIncrement();

		return;
	}

	public void start() {
		try{
			//System.out.println("Starting " + threadName);
			if (t == null) {
				t = new Thread(this, threadName);
				t.start();
			}
		} catch (Exception e){
			System.out.println("Misterious error in the thread...");
		}
	}
}