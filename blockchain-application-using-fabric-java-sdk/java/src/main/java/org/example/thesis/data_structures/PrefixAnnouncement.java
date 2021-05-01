package main.java.org.example.thesis.data_structures;

import java.io.Serializable;
import java.lang.Comparable;

/**
 * Class that represents an announcement to be stored in the local routing table
 * of an IXP (to be used as the key of an entry, with the value being the
 * announced path)
 */
public class PrefixAnnouncement implements Serializable, Comparable<PrefixAnnouncement> {

	/**
	 * Generated serial version for this class
	 */
	private static final long serialVersionUID = 6918761133277840751L;
	private String announcerAS; // Last AS on the path
	private String ipPrefix; // IP prefix that's being announced

	public PrefixAnnouncement(String as, String ip) {
		this.announcerAS = as;
		this.ipPrefix = ip;
	}

	public String getAnnouncerAS() {
		return announcerAS;
	}

	public void setAnnouncerAS(String announcerAS) {
		this.announcerAS = announcerAS;
	}

	public String getIpPrefix() {
		return ipPrefix;
	}

	public void setIpPrefix(String ipPrefix) {
		this.ipPrefix = ipPrefix;
	}

	@Override
	public int hashCode() {

		String joined = this.announcerAS + this.ipPrefix;
		int ascii = 0;

		for (int i = 0; i < joined.length(); i++) { // while counting characters if less than the length add one
			char character = joined.charAt(i); // start on the first character
			ascii += (int) character; // convert the first character
		}

		return ascii;
	}

	@Override
	public boolean equals(Object obj) {
		PrefixAnnouncement e = null;
		if (obj instanceof PrefixAnnouncement) {
			e = (PrefixAnnouncement) obj;
		}
		if (this.getAnnouncerAS().equals(e.getAnnouncerAS()) && this.getIpPrefix().equals(e.getIpPrefix())) {
			return true;
		} else {
			return false;
		}
	}

	@Override
	public int compareTo(PrefixAnnouncement pa) {
		int prefix = this.getIpPrefix().compareTo(pa.getIpPrefix());
		if (prefix == 0) {
			return this.getAnnouncerAS().compareTo(pa.getAnnouncerAS());
		}
		return prefix;
	}

}