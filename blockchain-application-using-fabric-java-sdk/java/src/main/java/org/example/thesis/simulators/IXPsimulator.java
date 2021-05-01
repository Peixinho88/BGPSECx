package main.java.org.example.thesis.simulators;

import java.util.Hashtable;
import java.util.Scanner;
import java.util.concurrent.ConcurrentLinkedQueue;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Semaphore;
import java.util.concurrent.atomic.AtomicInteger;

import main.java.org.example.thesis.auxiliary_blockchain_functions.IXPExecutionThread;
import main.java.org.example.thesis.auxiliary_blockchain_functions.RequestThread;
import main.java.org.example.thesis.auxiliary_blockchain_functions.TransactionProfile;
import main.java.org.example.thesis.data_structures.ASPath;
import main.java.org.example.thesis.data_structures.PrefixAnnouncement;
import main.java.org.example.thesis.data_structures.TransactionStructure;
import main.java.org.example.util.CsvWriter;
import main.java.org.example.util.GenerateConfig;

import java.io.EOFException;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.ObjectInputStream;
import java.lang.ClassNotFoundException;
import java.net.InetAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.net.UnknownHostException;

/*
 * Class that represents an IXP (internet exchange point).
 * All the logic of what goes in the blockchain and what 
 * doesn't (every verification) is here.
 */
public class IXPsimulator<K, V> {

	// Routing table of an IXP
	private static Hashtable<PrefixAnnouncement, ASPath> routingTable;
	// Static ServerSocket variable
	private static ServerSocket server;
	// Socket server port on which it will listen
	private static int dynamicPort = 0;
	// IXP id
	private static int ixpID;
	// Number of connections to this IXP
	private static int numConnections = 1;
	// Semaphore for access to the routing table file
	private static Semaphore fileSem = new Semaphore(1);
	// Semaphore for access to peer number
	private static Semaphore peerSem = new Semaphore(1);
	// Semaphore for csv file access
	private static Semaphore csvSem = new Semaphore(1);
	// Current peer
	public static volatile int currPeer = 0;
	// Number of client ASes
	public static int numAS = 2;
	//Create file writer for testing purposes
	public static volatile CsvWriter csvWriter;
	//Multi-thread counter
	private static AtomicInteger sequenceNumber = new AtomicInteger(0);

	public static void main(String[] args) throws UnknownHostException, IOException, ClassNotFoundException {

		Scanner sc = new Scanner(System.in); // nextInt() might need a nextLine() to consume the \n

		// Get the id for the IXP
		boolean receivedIXPID = false;
		System.out.println("Enter IXP id");
		while (!receivedIXPID) {
			if (sc.hasNextInt()) {
				ixpID = sc.nextInt();
				receivedIXPID = true;
			} else {
				System.out.println("IXP id must be a number");
			}
		}
		// Get the port for the IXP to receive connections
		boolean receivedPort = false;
		System.out.println("Enter port number");
		while (!receivedPort) {
			if (sc.hasNextInt()) {
				dynamicPort = sc.nextInt();
				receivedPort = true;
			} else {
				System.out.println("Port must be a number");
			}
		}

		// Initialize the routing table
		routingTable = new Hashtable<PrefixAnnouncement, ASPath>();

		try {
			String fileName = "routingTableFile" + ixpID + ".txt";
			File routingTableFile = new File(fileName);
			if (!routingTableFile.exists()) {
				System.out.println("Routing table file doesn't exist. Creating new one: " + fileName);
				routingTableFile.createNewFile();
			}
			FileInputStream fileIn = new FileInputStream(routingTableFile);
			ObjectInputStream in = new ObjectInputStream(fileIn);
			routingTable = (Hashtable<PrefixAnnouncement, ASPath>) in.readObject();
			in.close();
			fileIn.close();
		} catch (EOFException e) {
			System.out.println("Finished reading routing table from file.");
		} catch (ClassNotFoundException e) {
			e.printStackTrace();
		} catch (FileNotFoundException e) {
			System.out.println("File does not exist.");
		} catch (IOException e) {
			e.printStackTrace();
		}

		// Connection with the
		// ASes--------------------------------------------------------

		// Create the socket server object
		server = new ServerSocket(dynamicPort, 0, InetAddress.getLocalHost());
		// Stopping condition to terminate the server
		boolean terminated = false;
		Socket socket = null;

		// Generate network configuration
		GenerateConfig gc = new GenerateConfig();

		// Create a writer for the csv files
		csvWriter = new CsvWriter();

		ExecutorService execServReceivedConnections = Executors.newCachedThreadPool();
		ExecutorService execServTransactions = Executors.newCachedThreadPool();
		ConcurrentLinkedQueue<TransactionProfile> transactionQueue = new ConcurrentLinkedQueue<TransactionProfile>();
		ConcurrentLinkedQueue<TransactionStructure> verifyingQueue = new ConcurrentLinkedQueue<TransactionStructure>();

		// Loop to create fixed number of clients
		for (int i = 0; i < numAS; i++) {

			System.out.println("Waiting for client request");

			// Creating socket and waiting for client connection
			socket = server.accept();

			// Receives client connections, reads lines to be inserted in the blockchain and
			// updates the local routing table
			execServReceivedConnections.submit(new IXPExecutionThread<PrefixAnnouncement, ASPath>(socket, routingTable, ixpID, numConnections, 
												csvWriter, fileSem, peerSem, csvSem, gc.getConfig(), transactionQueue, verifyingQueue, sequenceNumber));

			numConnections++;
		}

		while (!terminated) {
			while (!transactionQueue.isEmpty()) {
				execServTransactions.submit(new RequestThread(csvWriter, csvSem, transactionQueue.poll(), verifyingQueue, sequenceNumber)); // TODO: PENSAR NISTO
			}

			// TODO: maybe ask if the IXP should terminate (think about this)

		}

		//Close the ServerSocket object
		System.out.println("Shutting down Socket server!!");
		server.close();
		sc.close();
		csvWriter.closeWriter();
	}
}
