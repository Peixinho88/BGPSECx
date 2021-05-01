package main.java.org.example.thesis.testing;

import main.java.org.example.util.GenerateConfig;

public class GenConfTest {
	public static void main(String[] args) {

		GenerateConfig cg = new GenerateConfig();

		for (int i = 0; i < cg.config.getOrgs().length; i++) {
			System.out.println(cg.config.getOrgs()[i].getOrgName());
			System.out.println(cg.config.getOrgs()[i].getOrgMSP());
			System.out.println(cg.config.getOrgs()[i].getOrgUsrBasePath());
			System.out.println(cg.config.getOrgs()[i].getOrgUsrAdminPK());
			System.out.println(cg.config.getOrgs()[i].getOrgUsrAdminCert());
			System.out.println(cg.config.getOrgs()[i].getMemberPeers().get(0).getPeerName());
			System.out.println(cg.config.getOrgs()[i].getMemberPeers().get(1).getPeerName());
		}
		System.out.println(cg.config.getOrgs()[2].getMemberPeers().get(2).getPeerName());

		System.out.println("----------------------------------------------------------");

		for (int i = 0; i < cg.config.getCaOrgs().length; i++) {
			System.out.println(cg.config.getCaOrgs()[i].getCaOrgName());
			System.out.println(cg.config.getCaOrgs()[i].getCaOrgURL());
		}

		System.out.println("----------------------------------------------------------");

		for (int i = 0; i < cg.config.getOrderers().size(); i++) {
			System.out.println(cg.config.getOrderers().get(i).getPeerName());
			System.out.println(cg.config.getOrderers().get(i).getPeerURL());
			System.out.println(cg.config.getOrderers().get(i).getOwnerOrg().getOrgName());
		}

		System.out.println("----------------------------------------------------------");

		System.out.println(cg.config.getAdminInfo().getAdminName());
		System.out.println(cg.config.getAdminInfo().getAdminPW());

		System.out.println("----------------------------------------------------------");

		System.out.println(cg.config.getChannelInfo().getChannelName());
		System.out.println(cg.config.getChannelInfo().getChannelConfigPath());

		System.out.println("----------------------------------------------------------");

		System.out.println(cg.config.getCcInfo().getCcName());
		System.out.println(cg.config.getCcInfo().getCcPath());
		System.out.println(cg.config.getCcInfo().getCcRootDir());
		System.out.println(cg.config.getCcInfo().getCcVersion());

	}
}