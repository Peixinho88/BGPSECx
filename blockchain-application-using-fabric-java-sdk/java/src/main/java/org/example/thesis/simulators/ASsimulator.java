package main.java.org.example.thesis.simulators;

import java.io.BufferedInputStream;
import java.io.BufferedReader;
import java.io.EOFException;
import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;
import java.net.InetAddress;
import java.net.Socket;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.nio.file.Paths;
import java.util.Scanner;

/*
 * Class that represents an AS (autonomous system)
 */
public class ASsimulator {

	private static int asNum;
	private static int port;

	public static void main(String[] args)
			throws UnknownHostException, IOException, ClassNotFoundException, InterruptedException {

		Scanner sc = new Scanner(System.in);

		// Get the ID for the AS
		boolean receivedID = false;
		System.out.println("Enter AS number");
		while (!receivedID) {
			if (sc.hasNextInt()) {
				asNum = sc.nextInt();
				receivedID = true;
				sc.nextLine(); // to consume the '\n' that nextInt() ignores
			} else {
				System.out.println("AS must be a number");
			}
		}

		// Get the port for the IXP to receive connections
		boolean receivedPort = false;
		System.out.println("Enter port number");
		while (!receivedPort) {
			if (sc.hasNextInt()) {
				port = sc.nextInt();
				receivedPort = true;
				sc.nextLine(); // to consume the '\n' that nextInt() ignores
			} else {
				System.out.println("Port must be a number");
			}
		}

		// get the localhost IP address, if server is running on some other IP, you need
		// to use that
		InetAddress host = InetAddress.getLocalHost();
		Socket socket = null;
		ObjectOutputStream oos = null;
		ObjectInputStream ois = null;
		boolean terminated = false;

		// Establish socket connection to server
		socket = new Socket(host.getHostName(), port);
		oos = new ObjectOutputStream(socket.getOutputStream());
		ois = new ObjectInputStream(socket.getInputStream());

		while (!terminated) {

			// Send a message to server to check connectivity
			System.out.println("Sending request to Socket Server");
			oos.writeObject("Connecting to IXP. AS number " + asNum);

			// read the server response message
			String message = (String) ois.readObject();
			System.out.println("Message: " + message);

			// Do other stuff
			// here----------------------------------------------------------------------------

			boolean correctInput = false;
			int inputType = 0;

			while (!correctInput) {
				System.out.println("Input type? ");
				if (sc.hasNext()) {
					if (sc.hasNextInt()) {
						inputType = sc.nextInt();
						if (inputType != 1 && inputType != 2 && inputType != 3) {
							System.out.println("Input is an incorrect number. Must be 1, 2 or 3.");
						} else {
							correctInput = true;
						}
						sc.nextLine(); // to consume the '\n' that nextInt() ignores
					} else if (sc.hasNextLine()) {
						String s = new String(sc.nextLine().toString());
						if (s.equalsIgnoreCase("exit")) {
							correctInput = true;
							terminated = true;
							oos.writeObject(1);
							oos.writeObject("exit");
						} else {
							System.out.println("Input is not a number. Must be 1, 2 or 3.");
						}
					}
				}
			}

			// Input type 1 (manual insertion of one update <ip_prefix as1 as2
			// asN>)................
			if (inputType == 1) {
				int linie = 1;
				oos.writeObject(linie);
				System.out.println("Enter a pair <ip prefix : as path> "); // example: 10.0.0.1/24 A B C
				String bgpUpdate = new String();
				if (sc.hasNextLine()) {
					bgpUpdate = sc.nextLine();
				}

				oos.writeObject(bgpUpdate);
			}
			// .....................................................................................

			// Input type 2 (path to a file with multiple bgp
			// entries)..............................
			else if (inputType == 2) {
				System.out.println("Enter path to BGP update file ");
				String bgpFilePath = new String();
				if (sc.hasNextLine()) {
					boolean correctFile = false;
					while(!correctFile) {
						try {
							bgpFilePath = sc.nextLine();
							String filePath = Paths.get(bgpFilePath).toString();
							File bgpFile = new File(filePath.trim());
							int numLines = count(filePath);
						
							// Sends the number of lines on the file first
							oos.writeObject(numLines);
							BufferedReader reader = new BufferedReader(new FileReader(bgpFile));
							correctFile = true;
							String line = new String();
							while (numLines > 0) {
								line = reader.readLine();
								oos.writeObject(line);
								numLines--;
							}
							reader.close();
						} catch (EOFException e) {
							System.out.println("File has finished reading.");
						} catch (SocketException e) {
							System.out.println("Socket error, why?");
						} catch (FileNotFoundException e) {
							System.out.println("File doesn't exist, try again.");
						}
					}
				}

			}
			// .....................................................................................

			// Input type 3 (print the routing table on the IXP)
			else if (inputType == 3) {
				int linie = 1;
				oos.writeObject(linie);
				oos.writeObject("print");
			}
			// .....................................................................................

			// Is is finished?
			System.out.println((String) ois.readObject());

			// -----------------------------------------------------------------------------------------------
		}
		// Close streams
		ois.close();
		oos.close();
		// Close other resources
		sc.close();
		socket.close();
	}

	public static int count(String filename) throws IOException {
		InputStream is = new BufferedInputStream(new FileInputStream(filename));
		try {
			byte[] c = new byte[1024];
			int count = 0;
			int readChars = 0;
			boolean endsWithoutNewLine = false;
			while ((readChars = is.read(c)) != -1) {
				for (int i = 0; i < readChars; ++i) {
					if (c[i] == '\n')
						++count;
				}
				endsWithoutNewLine = (c[readChars - 1] != '\n');
			}
			if (endsWithoutNewLine) {
				++count;
			}
			return count;
		} finally {
			is.close();
		}
	}

}
