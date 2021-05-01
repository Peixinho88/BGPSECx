package main.java.org.example.util;

import java.util.EnumSet;
import java.util.Vector;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.regex.Pattern;
import java.util.concurrent.ConcurrentLinkedQueue;

import com.google.protobuf.InvalidProtocolBufferException;

import org.hyperledger.fabric.sdk.BlockEvent;
import org.hyperledger.fabric.sdk.BlockInfo;
import org.hyperledger.fabric.sdk.ChaincodeEvent;
import org.hyperledger.fabric.sdk.ChaincodeEventListener;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.BlockInfo.EnvelopeInfo;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.exception.ProposalException;

public class SDKEventHandler {

    private boolean hasEvents;

    public SDKEventHandler() {
        this.hasEvents = false;
    }

    /**
     * Creates the chaincode event listener and defines what to do when it gets
     * caught
     * 
     * @param channel
     * @param expectedEventName
     * @param chaincodeEvents
     * @return
     * @throws InvalidArgumentException
     */
    public String setChaincodeEventListener(Channel channel, String expectedEventName, int executionNum,
            ConcurrentLinkedQueue<ChaincodeEventCapture> chaincodeEvents) throws InvalidArgumentException {

        ChaincodeEventListener chaincodeEventListener = new ChaincodeEventListener() {

            @Override
            public void received(String handle, BlockEvent blockEvent, ChaincodeEvent chaincodeEvent) {
                String s = new String(chaincodeEvent.getPayload());
                if(s.split("»")[3].equals(String.valueOf(executionNum))) {
                    chaincodeEvents.add(new ChaincodeEventCapture(handle, blockEvent, chaincodeEvent));
                }
                
                String eventHub = blockEvent.getPeer().toString();
                if (eventHub != null) {
                    eventHub = blockEvent.getPeer().getName();
                } else {
                    eventHub = blockEvent.getEventHub().getName();
                }
                // Here put what you want to do when receive chaincode event
                // System.out.println("RECEIVED CHAINCODE EVENT with handle: " + handle + ", chaincodeId: "
                //         + chaincodeEvent.getChaincodeId() + ", chaincode event name: " + chaincodeEvent.getEventName()
                //         + ", transactionId: " + chaincodeEvent.getTxId() + ", event Payload: "
                //         + new String(chaincodeEvent.getPayload()) + ", from eventHub: " + eventHub);

            }
        };
        // chaincode events.
        String eventListenerHandle = channel.registerChaincodeEventListener(Pattern.compile(".*"),
                Pattern.compile(Pattern.quote(expectedEventName)), chaincodeEventListener);
        
        if(chaincodeEvents.size() > 0) {
            this.hasEvents = true;
        }
        
        return eventListenerHandle;
    }

    /**
     * Waits for an event from the Chaincode
     * 
     * @param timeout
     * @param channel
     * @param chaincodeEvents
     * @param chaincodeEventListenerHandle
     * @return
     * @throws InvalidArgumentException
     */
    public boolean waitForChaincodeEvent(Integer timeout, Channel channel, ConcurrentLinkedQueue<ChaincodeEventCapture> chaincodeEvents, 
                                            String chaincodeEventListenerHandle) throws InvalidArgumentException {
        boolean eventDone = false;
        if (chaincodeEventListenerHandle != null) {

            int numberEventsExpected = channel.getEventHubs().size()
                    + channel.getPeers(EnumSet.of(Peer.PeerRole.EVENT_SOURCE)).size();
            // Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
            //         "numberEventsExpected: " + numberEventsExpected);
            // just make sure we get the notifications
            if (timeout.equals(0)) {
                // get event without timer
                /*
                 * while (chaincodeEvents.size() != numberEventsExpected) { // do nothing }
                 */
                eventDone = true;
            } else {
                // get event with timer
                for (int i = 0; i < timeout; i++) {
                    if (chaincodeEvents.size() == numberEventsExpected) {
                        eventDone = true;
                        break;
                    } else {
                        try {
                            // double j = i;
                            // j = j / 10;
                            // Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,j + "
                            // second");
                            Thread.sleep(100); // wait for the events for one tenth of second.
                        } catch (InterruptedException e) {
                            e.printStackTrace();
                        }
                    }
                }
            }

            // Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
            //         "chaincodeEvents.size(): " + chaincodeEvents.size());

            // unregister event listener
            channel.unregisterChaincodeEventListener(chaincodeEventListenerHandle);
            int i = 1;
            // arrived event handling
            for (ChaincodeEventCapture chaincodeEventCapture : chaincodeEvents) {
                Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO, "Event number. " + i);
                Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
                        "event capture object: " + chaincodeEventCapture.toString());
                Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
                        "Event Handle: " + chaincodeEventCapture.getHandle());
                Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
                        "Event TxId: " + chaincodeEventCapture.getChaincodeEvent().getTxId());
                Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
                        "Event Name: " + chaincodeEventCapture.getChaincodeEvent().getEventName());
                Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
                        "Event Payload: " + new String(chaincodeEventCapture.getChaincodeEvent().getPayload())); // byte
                Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
                        "Event ChaincodeId: " + chaincodeEventCapture.getChaincodeEvent().getChaincodeId());
                BlockEvent blockEvent = chaincodeEventCapture.getBlockEvent();
                try {
                    Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
                            "Event Channel: " + blockEvent.getChannelId());
                } catch (InvalidProtocolBufferException e) {
                    e.printStackTrace();
                }
                Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
                        "Event Hub: " + blockEvent.getEventHub());

                i++;
            }

        } 
        else {
            Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO,
                    "chaincodeEvents.isEmpty(): " + chaincodeEvents.isEmpty());
        }
        Logger.getLogger(SDKEventHandler.class.getName()).log(Level.INFO, "eventDone: " + eventDone);
        return eventDone;
    }

    public void blockInfoTest(Channel channel, BlockEvent blockEvent) {

        BlockInfo blockInfo = null;

        try {
            blockInfo = channel.queryBlockByNumber(blockEvent.getBlockNumber());
        } catch (InvalidArgumentException e) {
            e.printStackTrace();
        } catch (ProposalException e) {
            e.printStackTrace();
        }

        //TODO: for each para percorrer os envelopeInfos e ir ver individualmente se as transacções estão válidas
        for (EnvelopeInfo info : blockInfo.getEnvelopeInfos()) {
            System.out.println("THIS TX IS: " + info.isValid());
        }
    }

    public boolean hasEvents() {
        return hasEvents;
    }

    public void setHasEvents(boolean hasEvents) {
        this.hasEvents = hasEvents;
    }

}