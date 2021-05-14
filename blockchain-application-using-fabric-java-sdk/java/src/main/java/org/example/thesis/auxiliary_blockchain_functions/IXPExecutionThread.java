package main.java.org.example.thesis.auxiliary_blockchain_functions;

import main.java.org.example.thesis.data_structures.ASPath;
import main.java.org.example.thesis.data_structures.ConfigStructure;
import main.java.org.example.thesis.data_structures.PrefixAnnouncement;
import main.java.org.example.thesis.data_structures.TransactionStructure;
import main.java.org.example.thesis.simulators.IXPsimulator;
import main.java.org.example.util.ChaincodeEventCapture;
import main.java.org.example.util.ConsoleColors;
import main.java.org.example.util.CsvWriter;
import main.java.org.example.util.SDKEventHandler;

import java.io.EOFException;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;
import java.lang.ClassNotFoundException;
import java.net.Socket;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.Date;
import java.util.HashMap;
import java.util.Hashtable;
import java.util.Iterator;
import java.util.HashSet;
import java.util.Set;
import java.util.Vector;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.Semaphore;
// import java.util.logging.Level;
// import java.util.logging.Logger;
import java.util.concurrent.atomic.AtomicInteger;

import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.BlockInfo.EnvelopeInfo;

public class IXPExecutionThread<K, V> implements Runnable {
	private Thread t;
	private String threadName;
	private Socket socket;
	private Hashtable<PrefixAnnouncement, ASPath> routingTable;
	private int ixpID;
	private Semaphore fileSemaphore;
	private Semaphore peerSemaphore;
	private Semaphore csvSemaphore;
	private int numConnections = 0;
	private CsvWriter csvWriter;
	private ConfigStructure cs;
	private ConcurrentLinkedQueue<TransactionProfile> transactionQueue;
	private ConcurrentLinkedQueue<TransactionStructure> verifyingQueue;
	private AtomicInteger sequenceNumber;

	public IXPExecutionThread(Socket socket, Hashtable<PrefixAnnouncement, ASPath> routingTable, int ixpID,
			int numConnections, CsvWriter writer, Semaphore fileSem, Semaphore peerSem, Semaphore csvSem,
			ConfigStructure cs, ConcurrentLinkedQueue<TransactionProfile> transactionQueue,
			ConcurrentLinkedQueue<TransactionStructure> verifyingQueue, AtomicInteger sequenceNumber) {
		this.transactionQueue = transactionQueue;
		this.verifyingQueue = verifyingQueue;
		this.numConnections = numConnections;
		this.csvWriter = writer;
		this.ixpID = ixpID;
		this.threadName = this.ixpID + " connection " + this.numConnections; // TODO: this might be wrong/incomplete
		this.socket = socket;
		this.routingTable = routingTable;
		this.fileSemaphore = fileSem;
		this.peerSemaphore = peerSem;
		this.csvSemaphore = csvSem;
		this.cs = cs;
		this.sequenceNumber = sequenceNumber;
	}

	public void run() {

		boolean terminated = false;
		boolean addedHeaders = false;
		int clientExecution = 0;
		int lastBlock = 3;

		// Streams to communicate with the AS
		ObjectInputStream ois = null;
		ObjectOutputStream oos = null;
		String receivedMessage = null;

		try {
			// Read from socket to ObjectInputStream object
			ois = new ObjectInputStream(socket.getInputStream());
			// Create ObjectOutputStream object
			oos = new ObjectOutputStream(socket.getOutputStream());
		} catch (IOException e) {
			e.printStackTrace();
		}

		// Keeps listening indefinitely until it receives 'exit' call or program
		// terminates
		while (!terminated) {

			//Restart atomic integer that is used to check the blockchain insertion thread count
			this.sequenceNumber.set(0);

			// Tries to connect to the AS
			try {

				// Convert ObjectInputStream object to String
				receivedMessage = (String) ois.readObject();

				if (receivedMessage != null) {
					System.out.println("Message Received: " + receivedMessage);
					// Write object to Socket
					oos.writeObject("Connected to IXP " + ixpID);
				}

				// Blockchain and routing table
				// insertions~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

				// Behaves the same whether it receives one update manually or a stream of
				// updates from bgp file (options 1 and 2)
				String bgpEntry = new String();
				// Current line being read
				int lineNumber = 1;

				// Read number of lines on the file first (numLines are the lines still left)
				int numLines = (int) ois.readObject();
				int totalNumLines = numLines;
				//System.out.println("Total updates to be received: " + numLines);

				// Blockchain context to allow user communication with peers
				BlockchainContext bcContext = null;

				// Rotate the peers for each client (synchronized)
				try {
					this.peerSemaphore.acquire();
					if (IXPsimulator.currPeer == (cs.getOrgs()[ixpID - 1].getMemberPeers().size() - 1)) {
						IXPsimulator.currPeer = 0;
					} else {
						IXPsimulator.currPeer++;
					}
					bcContext = new BlockchainContext(cs, ixpID, IXPsimulator.currPeer);
					this.peerSemaphore.release();
				} catch (InterruptedException e) {
					e.printStackTrace();
				}

				// All the steps to create a working peer that can interact with the blockchain
				bcContext.createUserContext();
				bcContext.createChannel();
				bcContext.createPeer();

				// Bgp entries counters
				int bcEntries = 0;
				int rtEntries = 0;
				int totalEntriesRT = 0;
				int maliciousEntries = 0;

				// Start timer to count the time it takes to process all transactions
				long startTime = System.nanoTime();

				// Initial timestamp to measure each transaction elapsed time
				Date initialTimestamp = new Date();

				// EVENT LISTENER SET TEST----------------------------------------------------------------------------------
				String expectedEventName = "UpdateBGP";
				ConcurrentLinkedQueue<ChaincodeEventCapture> chaincodeEvents = new ConcurrentLinkedQueue<ChaincodeEventCapture>(); // Test
																																	// list
																																	// to
																																	// capture
				SDKEventHandler sdkEH = new SDKEventHandler();
				String chaincodeEventListenerHandle = null;
				try {
					chaincodeEventListenerHandle = sdkEH.setChaincodeEventListener(bcContext.getChannel(),
							expectedEventName, clientExecution, chaincodeEvents);
				} catch (InvalidArgumentException e) {
					e.printStackTrace();
				}
				// ---------------------------------------------------------------------------------------------------------

				// Loop to receive all the BGP updates from the AS
				while (numLines > 0) {
					/**System.out.println();
					System.out.println("Current numLines is: " + numLines);
					System.out.println();**/

					// Read line with BGP update from the AS
					bgpEntry = (String) ois.readObject();

					// If AS sends an exit, terminate the loop
					if (bgpEntry.equalsIgnoreCase("exit")) {
						break;
					}

					// If AS sends a print, terminate the loop
					if (bgpEntry.equalsIgnoreCase("print")) {
						break;
					}

					/**System.out.println();
					System.out.println("Inserting line " + lineNumber + ": " + bgpEntry);
					System.out.println();**/

					// Separating entry: [ip prefix] [path0] [path1] [...] ...
					String[] splitBGPEntry = bgpEntry.split(" ");
					// separatedPath is the entry without the ip prefix (just the full path)
					String[] separatedPath = Arrays.copyOfRange(splitBGPEntry, 1, splitBGPEntry.length);
					// path joined by '-' to be used by the functions in golang (by the blockchain)
					String rawPath = String.join("-", separatedPath);
					// IP prefix to be stored in the blockchain
					String ipPrefix = splitBGPEntry[0];

					// Check for loops and duplicates in the path
					String path;
					if ((path = checkForLoops(rawPath)).equals("0")) {
						lineNumber++;
						numLines--;
						continue;
					}

					// Add each transaction to the queue to be inserted in the
					// blockchain-------------------------
					TransactionProfile txProfile = new TransactionProfile("Transaction " + lineNumber, bcContext,
							ipPrefix, path, clientExecution);
					transactionQueue.add(txProfile);
					// -------------------------------------------------------------------------------------------
					
					// Update the line counters
					lineNumber++;
					numLines--;
				}

				clientExecution++;
				
				if (!(bgpEntry.equalsIgnoreCase("exit") || bgpEntry.equalsIgnoreCase("print"))) {

					// Create csv file to write the elapsed time values
					csvWriter.createFile("ixp" + ixpID + "_connection" + numConnections + "_execution" + clientExecution + ".csv");
					addedHeaders = false;
					
					// Add headers to the csv file (do it only the first time, hence the boolean)
					if (!addedHeaders) {
						csvWriter.addHeaders(totalNumLines, 1);
						addedHeaders = true;
					}

					//TODO: testar isto tudo muito bem
					//Waits for all the threads to finish
					int timer = 0;
					int lastSeqNum = 0;
					while (this.sequenceNumber.get() != totalNumLines) {
						
						//If a certain time passes, exit the cycle
						if(timer == 60) {
							break;
						}

						//Check if it got stuck (if it did increase timer, if not restart it)
						if(lastSeqNum == this.sequenceNumber.get()) {
							timer++;
						} else {
							timer = 0;
						}

						//Update value
						lastSeqNum = this.sequenceNumber.get();
						
						//Sleep 5 seconds and print the last terminated thread total
						try {
							Thread.sleep(5000);
						} catch (InterruptedException e) {
							e.printStackTrace();
						}
						System.out.println("THREAD CYCLE: " + this.sequenceNumber.get() + "/" + totalNumLines);
						System.out.println("TIMER: " + timer);
					}

					//EVENT LISTENER CAPTURE----------------------------------------------------------------------------------------------------
					boolean eventDone = false;
					try {
						eventDone = sdkEH.waitForChaincodeEvent(200, bcContext.getChannel(), chaincodeEvents, chaincodeEventListenerHandle);
					} catch (InvalidArgumentException e) {
						System.out.println("Error waiting for event");
						e.printStackTrace();
					}
					//Logger.getLogger(IXPsimulator.class.getName()).log(Level.INFO, "eventDone: " + eventDone);
					//--------------------------------------------------------------------------------------------------------------------------

					//EVENT HANDLING------------------------------------------------------------------------------------------------------------
					EnvelopeInfo correctEnvelope;
					boolean fileUpdate = false;

					//Stuff for block commits
					HashMap<Long, Integer> blockCommits = new HashMap<Long, Integer>();
					CsvWriter commitCountWriter = new CsvWriter();
					commitCountWriter.createFile("BlockCommitCounting.csv");
					int numBlocks = (int) ((Math.floor((chaincodeEvents.size() * 2) / 10.0) + Math.ceil(((chaincodeEvents.size() * 2) % 10) / 10.0)));
					commitCountWriter.addHeaders(numBlocks, 2);

					System.out.println("NumBlocks is: " + numBlocks);
					
					//Loop to go through all the events and find the envelope that matches said event
					for (ChaincodeEventCapture chaincodeEventCapture : chaincodeEvents) {
						fileUpdate = true;
						
						//Get the correct envelope corresponding to the current event
						correctEnvelope = chaincodeEventCapture.getTx();
						if(correctEnvelope == null) {
							System.out.println("Error reading the events. Something went wrong.");
							continue;
						}

						//Count the commits per block and add them to the hashmap
						chaincodeEventCapture.getCommitsPerBlock(blockCommits);

						//Elapsed time to be written to the file
						long timeDiff = correctEnvelope.getTimestamp().getTime() - initialTimestamp.getTime();
						
						try {
							this.csvSemaphore.acquire();
							this.csvWriter.addValue(String.valueOf(timeDiff), false);
						} catch (InterruptedException e) {
							e.printStackTrace();
						}
						this.csvSemaphore.release();

						//Get the payload from the event and use it to check if it's verified, and get the path and the prefix
						String payload = new String(chaincodeEventCapture.getChaincodeEvent().getPayload());
						String[] splitPayload = payload.split("»");
						String[] splitAS = splitPayload[1].split("-");
						String announcerAS = splitAS[splitAS.length - 1];
						PrefixAnnouncement preAnn = new PrefixAnnouncement(announcerAS, splitPayload[0]);
						ASPath asPathUpdated = new ASPath(splitPayload[1], splitAS.length);
	
						boolean verified = false;

						//Check if the transaction was committed and if it is verified
						if(correctEnvelope.isValid() && (splitPayload[2].equals("V"))) {
							verified = true;
							bcEntries++;
							rtEntries++;
							//TODO: maybe add a "total entries to the routing table" and differentiate between those and the verified ones
							for (String eachAS : splitAS) {
								if(eachAS.equals("66666")) {
									maliciousEntries++;
									break;
								}
							}
						} else {
							verified = false;
						}
						asPathUpdated.setVerified(verified);
						routingTable.put(preAnn, asPathUpdated);
						totalEntriesRT++; //TODO: pensar se me interessa contabilizar isto

						chaincodeEvents.remove(chaincodeEventCapture);
					}
					System.out.println(Arrays.asList(blockCommits)); //TODO: for testing
					//--------------------------------------------------------------------------------------------------------------------------

					//End timer to count the time it takes to process all transactions
					long endTime = System.nanoTime();

					//Get difference of two nanoTime values
					long timeElapsed = endTime - startTime;

					//Write total time elapsed, plus all the counters, to the file, and the block commits to the second file
					try {
						this.csvSemaphore.acquire();
						this.csvWriter.addValue(String.valueOf(timeElapsed / 1000000) + "," 
												+ String.valueOf(totalEntriesRT) + "," + String.valueOf(bcEntries) + "," 
												+ String.valueOf(rtEntries) + "," + String.valueOf(maliciousEntries), true);
						
						for(long i = lastBlock; i < (numBlocks + lastBlock); i++) {
							int aux;
							if(blockCommits.get(i) == null) {
								aux = 0;
							} else {
								aux = blockCommits.get(i);
							}

							if(i == (numBlocks + lastBlock - 1)) {
								commitCountWriter.addValue(String.valueOf(aux), true);
							} else {
								commitCountWriter.addValue(String.valueOf(aux), false);
							}
						}

						lastBlock += numBlocks;

					} catch (InterruptedException e) {
						e.printStackTrace();
					}
					this.csvSemaphore.release();

				}

				// Finished receiving file lines from AS
				oos.writeObject("IXP has finished receiving!");

				// Terminate the server if client sends exit request
				if (bgpEntry.equalsIgnoreCase("exit")) {
					terminated = true;
				}

				// If client sends a print, print the contents of the routing table
				if (bgpEntry.equalsIgnoreCase("print")) {

					//Actual printing of the table in this next section
					Vector<PrefixAnnouncement> v = new Vector<PrefixAnnouncement> (routingTable.keySet());
					Collections.sort(v);
					Iterator<PrefixAnnouncement> it = v.iterator();
					System.out.println(
							"____________________________________________________________________________________________________________________________________________________________");
					System.out.println();
					System.out.println(String.format("%86s", "IXP " + ixpID + " ROUTING TABLE"));
					System.out.println(
							"____________________________________________________________________________________________________________________________________________________________");
					System.out.println();
					System.out.println(String.format("%18s %-47s %-74s %-15s", "Announcer", "| " + "IP Prefix",
							"| " + "AS Path", "| " + "Verification"));
					System.out.println(
							"____________________________________________________________________________________________________________________________________________________________");
					System.out.println();
					while (it.hasNext()) {
						PrefixAnnouncement element = (PrefixAnnouncement) it.next();
						String colorVerified = new String();
						if (routingTable.get(element).getVerified()) {
							colorVerified = ConsoleColors.GREEN + "Verified";
						} else {
							colorVerified = ConsoleColors.RED + "Not Verified";
						}
						System.out.println(String.format("%25s %-58s %-78s %-30s",
								ConsoleColors.YELLOW_BOLD + element.getAnnouncerAS(),
								ConsoleColors.RESET + "| " + ConsoleColors.WHITE_BOLD + element.getIpPrefix(),
								ConsoleColors.RESET + "| " + routingTable.get(element).getPath(),
								"| " + colorVerified + ConsoleColors.RESET));
					}
					System.out.println(
							"____________________________________________________________________________________________________________________________________________________________");
					System.out.println();
				}

				// Write the routing table to the file
				try {
					this.fileSemaphore.acquire();
					FileOutputStream fileOut = new FileOutputStream("routingTableFile" + this.ixpID + ".txt");
					ObjectOutputStream out = new ObjectOutputStream(fileOut);
					out.writeObject(routingTable);
					out.close();
					fileOut.close();
				} catch (FileNotFoundException e) {
					e.printStackTrace();
				} catch (IOException e) {
					e.printStackTrace();
				} catch (InterruptedException e) {
					e.printStackTrace();
				}
				this.fileSemaphore.release();

			} catch (SocketException e) {
				terminated = true;
				System.out.println("Socket connection error, for some reason...");
			} catch (EOFException e) {
				terminated = true;
				System.out.println("Finished reading from socket");
			} catch (ClassNotFoundException e) {
				e.printStackTrace();
			} catch (IOException e) {
				e.printStackTrace();
			} catch (ArrayIndexOutOfBoundsException e) {
				e.printStackTrace();
			}
			// -----------------------------------------------------------------
		}
		try {
			// close resources
			ois.close();
			oos.close();
			socket.close();
		} catch (IOException e) {
			System.out.println("Error closing the socket");
		}
	}

	public void start() throws UnknownHostException, IOException, ClassNotFoundException {
		//System.out.println("Starting " + threadName);
		if (t == null) {
			t = new Thread(this, threadName);
			t.start();
		}
	}

	/**
	 * Checks for loops in an AS path and if it has none returns a new path, with
	 * adjacent duplicates already removed
	 * 
	 * @param path - AS path separated by '-'
	 * @return new path without duplicates or '0' if it has any loops
	 */
	public static String checkForLoops(String path) {
		String newPath = checkForDuplicates(path);
		String[] separatedASes = newPath.split("-");

		Set<String> duplicates = new HashSet<String>();

		for (String s : separatedASes) {
			if (!duplicates.add(s)) {
				return "0";
			}
		}
		return newPath;
	}

	/**
	 * Checks for adjacent duplicate ASes in a path
	 * 
	 * @param path - AS path separated by '-'
	 * @return a new path separated by '-' without any adjacent duplicates
	 */
	public static String checkForDuplicates(String path) {
		String[] separatedASes = path.split("-");
		ArrayList<String> newPath = new ArrayList<String>();

		newPath.add(separatedASes[0]);
		for (int i = 1; i < separatedASes.length; i++) {
			if (!separatedASes[i].equals(separatedASes[i - 1])) {
				newPath.add(separatedASes[i]);
			}
		}

		String correctPath = String.join("-", newPath);

		return correctPath;

	}
}