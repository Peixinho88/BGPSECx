package main.java.org.example.util;

import org.hyperledger.fabric.sdk.BlockEvent;
import org.hyperledger.fabric.sdk.BlockInfo.EnvelopeInfo;
import org.hyperledger.fabric.sdk.ChaincodeEvent;

public class ChaincodeEventCapture {
  private final String handle;
  private final BlockEvent blockEvent;
  private final ChaincodeEvent chaincodeEvent;

  public ChaincodeEventCapture(String handle, BlockEvent blockEvent, ChaincodeEvent chaincodeEvent) {
    this.handle = handle;
    this.blockEvent = blockEvent;
    this.chaincodeEvent = chaincodeEvent;
  }

  /**
   * @return the handle
   */
  public String getHandle() {
    return handle;
  }

  /**
   * @return the blockEvent
   */
  public BlockEvent getBlockEvent() {
    return blockEvent;
  }

  /**
   * @return the chaincodeEvent
   */
  public ChaincodeEvent getChaincodeEvent() {
    return chaincodeEvent;
  }

  public EnvelopeInfo getTx() {

    String queryID = this.chaincodeEvent.getTxId();

    for (EnvelopeInfo info : this.blockEvent.getEnvelopeInfos()) {
      if(queryID.equals(info.getTransactionID())) {
        return info;
      }
    }
    return null;
  }
}