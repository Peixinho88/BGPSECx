package main.java.org.example.thesis.testing;

import java.io.IOException;
import java.net.InetAddress;
import java.net.ServerSocket;
import java.net.UnknownHostException;

public class ServerTest {

	public static void main(String[] args) throws UnknownHostException, IOException {
		ServerSocket ss = new ServerSocket(0, 0, InetAddress.getLocalHost());

		System.out.println("port: " + ss.getLocalPort());
		System.out.println("ip address: " + ss.getInetAddress());
	}

}