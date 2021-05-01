package main.java.org.example.thesis.data_structures;

public class ChannelInfo {

	private String channelName;
	private String channelConfigPath;

	public ChannelInfo(String name, String configPath) {
		this.channelName = name;
		this.channelConfigPath = configPath;
	}

	public String getChannelName() {
		return this.channelName;
	}

	public void setChannelName(String channelName) {
		this.channelName = channelName;
	}

	public String getChannelConfigPath() {
		return this.channelConfigPath;
	}

	public void setChannelConfigPath(String channelConfigPath) {
		this.channelConfigPath = channelConfigPath;
	}
}