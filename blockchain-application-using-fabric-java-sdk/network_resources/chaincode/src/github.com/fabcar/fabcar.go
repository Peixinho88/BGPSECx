/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	//"io"
	//"github.com/hyperledger/fabric-chaincode-go/blob/master/shim/stub.go" //- KEEP THIS ONE JUST IN CASE. SEEMS TO BE THE UPDATED SHIM

	//@ts-ignore
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	//sc "github.com/hyperledger/fabric-protos-go/peer"
)

//LAST MODIFICATIONS:
//IF ON LINE 646 (NOT SURE HOW TO CHECK FOR EMPTY TREE)
//IF AND ELSE ON INSERTPATHONTREE METHOD
//
//STILL HAVE TO TEST ALL OF THIS
//['/home/peixinho/Desktop/Test_files/rrc00.20180501.0000.parsed']

//To check the docker log
//docker logs dev-peer0.org1.example.com-fabcar-1.0

// Define the Smart Contract structure
type SmartContract struct {
}

// Defines a key/value pair to store the info from the iterator returned by GetStateByPartialCompositeKey
type pairKeyValue struct {
	Key   string
	Value []byte
}

// Define the structure for an AS Path, which will be the values stored on the key/value store (I think)
type ASPath struct {
	Path string `json:"path"`
	//maybe something else? don't know yet
}

//Read more about this and test it
type BGPTree struct {
	Value    string    `json:"value"`
	Status   string    `json:"status"`
	Children []BGPTree `json:"children"`
}

//Structure for a complaint which stores the prefix and the path that are being complained, and the list of ASes
//who filed the complaint
type Complaint struct {
	ASList string `json:"asList"`
	//timestamp, maybe? Think about this
}

//Converts a string array to int array
func sliceItoa(si []int) ([]string, error) {
	sa := make([]string, 0, len(si))
	for _, i := range si {
		a := strconv.Itoa(i)
		sa = append(sa, a)
	}
	return sa, nil
}

//Remove duplicates from a slice
func unique(s []int) []int {
	seen := make(map[int]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

//Checks if an IP addres is IPv4 or IPv6
func checkIPAddressType(ip string) int {
	if net.ParseIP(ip) == nil {
		return -1 //Error
	}
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			return 0 //IPv4
		case ':':
			return 1 //IPv6
		}
	}
	return -1
}

//Returns a list with all possible octets that contain the given one (for ipv4)
func octetCombinationsIPv4(octet string) []int {
	octetValue, _ := strconv.Atoi(octet)
	mask := 255

	var possibilities []int
	for i := 0; i <= 8; i++ {
		aux := octetValue & mask
		possibilities = append(possibilities, aux)
		mask = mask << 1
	}

	return unique(possibilities)
}

//Returns a list with all possible octets that contain the given one (for ipv6)
func octetCombinationsIPv6(octet string) []int {
	//octetValue, _ := strconv.Atoi(octet)
	hexOctet := "0x" + octet
	octetValue64, _ := strconv.ParseInt(hexOctet, 0, 64)
	octetValue := int(octetValue64)
	mask := 65535

	var possibilities []int
	for i := 0; i <= 16; i++ {
		aux := octetValue & mask
		possibilities = append(possibilities, aux)
		mask = mask << 1
	}

	return unique(possibilities)
}

//Gets the separated octets and mask, and returns the string which is whole prefix
func reconstructPrefix(prefixParts []string, ipType int) string {
	var octetsOnly []string
	var octets string
	for i := 0; i < len(prefixParts)-1; i++ {
		octetsOnly = append(octetsOnly, prefixParts[i])
	}

	if ipType == 0 { //IPv4
		octets = strings.Join(octetsOnly, ".")
	} else if ipType == 1 { //IPv6
		octets = strings.Join(octetsOnly, ":")
		//octets = octets + "::"
	}
	finalPrefix := octets + "/" + prefixParts[len(prefixParts)-1]
	return finalPrefix
}

//Completes the missing octets in an IPv6 address (with zeros)
func handleIPv6Address(ip string) string {
	if strings.Contains(ip, "::") {
		s := strings.Split(ip, "::")
		ss := strings.Split(s[1], "/")
		fmt.Println(ss)
		sss := strings.Split(ss[0], ":")
		newS := []string{s[0]}
		if strings.Contains(string(s[1][0]), "/") { //:: is right before the mask
			fmt.Println("dentro do if dos :: no fim")
			for i := 0; i < 8-len(strings.Split(s[0], ":")); i++ {
				newS = append(newS, "0")
			}
		} else if s[0] == "" { //:: is on the beginning on the address
			fmt.Println("dentro do if dos :: no inicio")
			newS = []string{}
			fmt.Println(len(sss))
			for i := 0; i < 8-len(sss); i++ {
				fmt.Println("dentro do ciclo de adicionar zeros")
				fmt.Println(i)
				newS = append(newS, "0")
			}
			newS = append(newS, ss[0])
		} else { //:: is in the middle of the address somewhere
			fmt.Println("dentro do if dos :: no meio")
			x := len(strings.Split(newS[0], ":")) + len(sss)
			for i := 0; i < 8-x; i++ {
				newS = append(newS, "0")
			}
			newS = append(newS, ss[0])
		}

		finalS := strings.Join(newS, ":")
		finalS = finalS + "/" + ss[1]

		return finalS
	}
	return "erro"
}

//Checks whether a given prefix already exists on the blockchain list of broader prefixes
func checkSubnet(APIstub shim.ChaincodeStubInterface, prefix string, ipType int) (bool, bool) {
	resultsIterator, _ := APIstub.GetStateByPartialCompositeKey("MASK", []string{})
	//100.10.10.10/24 = 100.10.10.0/24
	//100.10.10.30/25 = 100.10.10.0/25
	//Loop to go through every IP prefix stored
	for resultsIterator.HasNext() {
		queryResponse, _ := resultsIterator.Next()

		_, compositeKeyParts, _ := APIstub.SplitCompositeKey(queryResponse.GetKey())

		//Reconstruct according to the IP type that is stored
		var reconstructedPrefix string
		if len(compositeKeyParts) == 5 {
			reconstructedPrefix = reconstructPrefix(compositeKeyParts, 0)
		} else if len(compositeKeyParts) == 9 {
			reconstructedPrefix = reconstructPrefix(compositeKeyParts, 1)
		}

		//reconstructedPrefix := reconstructPrefix(compositeKeyParts, ipType)
		foundIP, foundPrefix, _ := net.ParseCIDR(reconstructedPrefix)
		foundMask, _ := strconv.Atoi(strings.Split(foundPrefix.String(), "/")[1])
		givenIP, givenPrefix, _ := net.ParseCIDR(prefix)
		givenMask, _ := strconv.Atoi(strings.Split(givenPrefix.String(), "/")[1])

		fmt.Println("found prefix: " + foundPrefix.String())
		fmt.Println("given prefix: " + givenPrefix.String())

		//Given prefix and found prefix are the same - WON'T BE INSERTED BUT AS PATH CAN BE UPDATED
		if givenPrefix.String() == foundPrefix.String() {
			return true, true
		}

		//Given prefix is a supernet of one already in the blockchain - CAN'T BE INSERTED
		if givenPrefix.Contains(foundIP) && givenMask < foundMask {
			return true, false
		}

		//Given prefix is a subtnet of one already in the blockchain - CAN BE INSERTED IN THE BLOCKCHAIN (IP LIST ONLY)
		if foundPrefix.Contains(givenIP) {
			return false, true
		}
	}
	//Given prefix isn't a subnet or a supernet of a stored one - CAN BE INSERTED IN THE BLOCKCHAIN (BOTH IP AND MASK LISTS)
	return false, false
}

//Returns a list with all the subnets that match a given IPv4 prefix
func getAllPossiblePrefixesIPv4(APIstub shim.ChaincodeStubInterface, prefix string) []pairKeyValue {

	ip, ipNetwork, _ := net.ParseCIDR(prefix) //se não for usar tenho de remover o ip

	ipAndMask := strings.Split(ipNetwork.String(), "/") //[10.10.220.0] [17]
	octets := strings.Split(ipAndMask[0], ".")          //[10] [10] [220] [0]
	//compositePrefix := append(octets, ipAndMask[1])     ////[10] [10] [220] [0] / [7]
	//finalCompositeKey, _ := APIstub.CreateCompositeKey("IP", compositePrefix)

	//TODO: check if ip is ipv4 or ipv6
	mask, _ := strconv.Atoi(ipAndMask[1])
	var zero []string = []string{"0"} //TODO: mudar esta porcaria para não tar a criar um array sem razão nenhuma
	var bestMatches []pairKeyValue

	//TODO: test this for the last two ifs (first two are right, so it's probably all right, but you never know)
	if mask < 8 {
		//Verification for combinations of the first octet
		intCombos := octetCombinationsIPv4(octets[0])
		combos, _ := sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	} else if mask >= 8 && mask < 16 {
		//Verification for combinations of the second octet
		intCombos := octetCombinationsIPv4(octets[1])
		combos, _ := sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], combos[i], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the first octet
		intCombos = octetCombinationsIPv4(octets[0])
		combos, _ = sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	} else if mask >= 16 && mask < 24 {
		//Verification for combinations of the third octet
		intCombos := octetCombinationsIPv4(octets[2])
		combos, _ := sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], combos[i], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the second octet
		intCombos = octetCombinationsIPv4(octets[1])
		combos, _ = sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], combos[i], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the first octet
		intCombos = octetCombinationsIPv4(octets[0])
		combos, _ = sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	} else if mask >= 24 && mask <= 32 {
		//Verification for combinations of the fourth octet
		intCombos := octetCombinationsIPv4(octets[3])
		combos, _ := sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], combos[i]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the third octet
		intCombos = octetCombinationsIPv4(octets[2])
		combos, _ = sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], combos[i], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the second octet
		intCombos = octetCombinationsIPv4(octets[1])
		combos, _ = sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], combos[i], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the first octet
		intCombos = octetCombinationsIPv4(octets[0])
		combos, _ = sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	}
	return bestMatches
}

//Returns a list with all the subnets that match a given IPv6 prefix (this is awful, don't open this)
func getAllPossiblePrefixesIPv6(APIstub shim.ChaincodeStubInterface, prefix string) []pairKeyValue {

	ip, _, _ := net.ParseCIDR(prefix) //se não for usar tenho de remover o ip

	ipAndMask := strings.Split(prefix, "/")    //[10.10.220.0] [17]
	octets := strings.Split(ipAndMask[0], ":") //[10] [10] [220] [0]
	//compositePrefix := append(octets, ipAndMask[1])     ////[10] [10] [220] [0] / [7]
	//finalCompositeKey, _ := APIstub.CreateCompositeKey("IP", compositePrefix)

	//TODO: check if ip is ipv4 or ipv6
	mask, _ := strconv.Atoi(ipAndMask[1])
	var zero []string = []string{"0"} //TODO: mudar esta porcaria para não tar a criar um array sem razão nenhuma
	var bestMatches []pairKeyValue

	//TODO: test this for the last two ifs (first two are right, so it's probably all right, but you never know)
	if mask < 16 {
		//Verification for combinations of the first octet
		intCombos := octetCombinationsIPv6(octets[0])
		var combos []string
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		//combos, _ := sliceItoa(intCombos)
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	} else if mask >= 16 && mask < 32 {
		//Verification for combinations of the second octet
		intCombos := octetCombinationsIPv6(octets[1])
		var combos []string
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the first octet
		intCombos = octetCombinationsIPv6(octets[0])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	} else if mask >= 32 && mask < 48 {
		//Verification for combinations of the third octet
		intCombos := octetCombinationsIPv6(octets[2])
		var combos []string
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the second octet
		intCombos = octetCombinationsIPv6(octets[1])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the first octet
		intCombos = octetCombinationsIPv6(octets[0])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	} else if mask >= 48 && mask < 64 {
		//Verification for combinations of the fourth octet
		intCombos := octetCombinationsIPv6(octets[3])
		var combos []string
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], combos[i], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the third octet
		intCombos = octetCombinationsIPv6(octets[2])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the second octet
		intCombos = octetCombinationsIPv6(octets[1])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the first octet
		intCombos = octetCombinationsIPv6(octets[0])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	} else if mask >= 64 && mask < 80 {
		//Verification for combinations of the fifth octet
		intCombos := octetCombinationsIPv6(octets[4])
		var combos []string
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], octets[3], combos[i], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the fourth octet
		intCombos = octetCombinationsIPv6(octets[3])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], combos[i], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the third octet
		intCombos = octetCombinationsIPv6(octets[2])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the second octet
		intCombos = octetCombinationsIPv6(octets[1])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the first octet
		intCombos = octetCombinationsIPv6(octets[0])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	} else if mask >= 80 && mask < 96 {
		//Verification for combinations of the sixth octet
		intCombos := octetCombinationsIPv6(octets[5])
		var combos []string
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], octets[3], octets[4], combos[i], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the fifth octet
		intCombos = octetCombinationsIPv6(octets[4])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], octets[3], combos[i], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the fourth octet
		intCombos = octetCombinationsIPv6(octets[3])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], combos[i], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the third octet
		intCombos = octetCombinationsIPv6(octets[2])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the second octet
		intCombos = octetCombinationsIPv6(octets[1])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the first octet
		intCombos = octetCombinationsIPv6(octets[0])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	} else if mask >= 96 && mask < 112 {
		//Verification for combinations of the seventh octet
		intCombos := octetCombinationsIPv6(octets[6])
		var combos []string
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], octets[3], octets[4], octets[5], combos[i], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the sixth octet
		intCombos = octetCombinationsIPv6(octets[5])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], octets[3], octets[4], combos[i], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the fifth octet
		intCombos = octetCombinationsIPv6(octets[4])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], octets[3], combos[i], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the fourth octet
		intCombos = octetCombinationsIPv6(octets[3])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], combos[i], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the third octet
		intCombos = octetCombinationsIPv6(octets[2])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the second octet
		intCombos = octetCombinationsIPv6(octets[1])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the first octet
		intCombos = octetCombinationsIPv6(octets[0])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	} else if mask >= 112 && mask <= 128 {
		//Verification for combinations of the eighth octet
		intCombos := octetCombinationsIPv6(octets[7])
		var combos []string
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], octets[3], octets[4], octets[5], octets[6], combos[i]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the seventh octet
		intCombos = octetCombinationsIPv6(octets[7])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], octets[3], octets[4], octets[5], octets[6], combos[i]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the sixth octet
		intCombos = octetCombinationsIPv6(octets[5])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], octets[3], octets[4], combos[i], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the fifth octet
		intCombos = octetCombinationsIPv6(octets[4])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], octets[3], combos[i], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the fourth octet
		intCombos = octetCombinationsIPv6(octets[3])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], octets[2], combos[i], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the third octet
		intCombos = octetCombinationsIPv6(octets[2])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], octets[1], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the second octet
		intCombos = octetCombinationsIPv6(octets[1])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{octets[0], combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}

		//Verification for combinations of the first octet
		intCombos = octetCombinationsIPv6(octets[0])
		for c := 0; c < len(intCombos); c++ {
			combos = append(combos, fmt.Sprintf("%x", intCombos[c]))
		}
		for i := 0; i < len(combos); i++ {
			iter, err := APIstub.GetStateByPartialCompositeKey("IP", []string{combos[i], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0], zero[0]})
			if err != nil {
				fmt.Println("Error in GetState: " + err.Error())
			}

			//Copies the iterator to a more manageable structure
			var kvArray []pairKeyValue
			for iter.HasNext() {
				aux, _ := iter.Next()
				kvAux := pairKeyValue{aux.GetKey(), aux.GetValue()}
				kvArray = append(kvArray, kvAux)
			}

			for j := 0; j < len(kvArray); j++ {
				fmt.Println("dentro do for do hasNext")
				next := kvArray[j]
				key := next.Key
				_, splitPrefix, _ := APIstub.SplitCompositeKey(key)
				recoveredPrefix := reconstructPrefix(splitPrefix, checkIPAddressType(ip.String()))
				_, bcIP, _ := net.ParseCIDR(recoveredPrefix)
				maskBC, _ := strconv.Atoi(strings.Split(recoveredPrefix, "/")[1])
				fmt.Println("maskBC: " + strconv.Itoa(maskBC))
				fmt.Println("mask: " + strconv.Itoa(mask))

				if maskBC <= mask && bcIP.Contains(ip) { // /7   /3  |  /7
					bestMatches = append(bestMatches, next)
				}
			}
		}
		return bestMatches
	}
	return bestMatches
}

//Returns the number of leaf nodes the tree has (Checked and works properly)
func nodeCounting(tree BGPTree) int {
	var childrenCounter int = 0

	if len(tree.Children) == 0 {
		return 1
	}

	for i := 0; i < len(tree.Children); i++ {
		childrenCounter += nodeCounting(tree.Children[i])
	}

	return childrenCounter
}

//Prints the paths of the tree, one on each line (Checked)
func queryTree(tree BGPTree, path string, firstTime bool) {
	if firstTime {
		path += tree.Value
		firstTime = false
	} else {
		path += "-"
		path += tree.Value
	}

	if len(tree.Children) == 0 {
		fmt.Println(path)
		return
	} else {
		for i := 0; i < len(tree.Children); i++ {
			queryTree(tree.Children[i], path, firstTime)
		}
		return
	}
}

//Returns a pathList which has a different path separated by "-" in each of the array's positions (Checked)
func queryTreeImproved(tree BGPTree, path string, pathList []string, index int, firstTime bool) []string {
	if tree.Value == "" {
		return pathList
	}
	if firstTime {
		path += tree.Value
		firstTime = false
	} else {
		path += "-"
		path += tree.Value
	}
	if len(tree.Children) == 0 {
		pathList[index] = path //<- This is how it previously was and it worked. Replace the line below with this one
		//pathList[index] = tree.Status + ":" + path
		return pathList
	} else {
		if len(tree.Children) == 1 {
			pathList = queryTreeImproved(tree.Children[0], path, pathList, index, firstTime)
		} else {
			var excessLeaves int = 0
			for i := 0; i < len(tree.Children); i++ {
				if i != 0 {
					excessLeaves += nodeCounting(tree.Children[i-1])
					pathList = queryTreeImproved(tree.Children[i], path, pathList, index+excessLeaves, firstTime)
				} else {
					pathList = queryTreeImproved(tree.Children[i], path, pathList, index, firstTime)
				}
			}
		}
	}
	return pathList
}

//Verifies a path based on a given one
func verifyPath(tree BGPTree, path string) bool {
	splitPath := strings.Split(path, "-")
	var newPath string

	newSplitPath := make([]string, (len(splitPath) - 1))
	for i := 1; i < len(splitPath); i++ {
		newSplitPath[i-1] = splitPath[i]
	}
	if len(newSplitPath) == 1 {
		newPath = newSplitPath[0]
	} else {
		newPath = strings.Join(newSplitPath, "-")
	}

	if tree.Value == "" && len(splitPath) == 1 {
		return true
	}

	//esta nas duas ultimas posicoes da arvore e vai verificar se pode inserir o path recebido
	//(se o ramo for maior que o path recebido, ele e inserido como irmão caso não seja igual, penso que sem problema)
	if len(splitPath) == 2 {
		if tree.Value == splitPath[0] {
			if len(tree.Children) == 0 {
				//pode inserir logo o splitPath[1]
				return true
			} else {
				//tem de verificar se algum dos filhos e igual
				for i := 0; i < len(tree.Children); i++ {
					if tree.Children[i].Value == splitPath[1] {
						//nao insere porque ja existe um filho igual
						return false
					} else if i == (len(tree.Children) - 1) {
						//todos os filhos sao diferentes e pode inserir como irmao deles
						return true
					}
				}
				return false
			}
		} else {
			return false
		}
		//se o path recebido tiver mais do que duas posicoes, a recursao tem de continuar e ir reduzindo o path
	} else {
		childrenLength := len(tree.Children)
		if tree.Value == splitPath[0] && childrenLength > 0 {
			//ciclo que percorre os filhos do no em questao e vai chamando a recursao para cada um deles
			for i := 0; i < childrenLength; i++ {
				verified := verifyPath(tree.Children[i], newPath)
				if verified {
					return true
				}
			}
		}
	}
	return false
}

//Inserts all paths on a tree, being Verified or Unverified (differentiates them)
func insertSinglePathOnTree(tree BGPTree, path string) BGPTree {
	splitPath := strings.Split(path, "-")
	var newPath string

	newSplitPath := make([]string, (len(splitPath) - 1))
	for i := 1; i < len(splitPath); i++ {
		newSplitPath[i-1] = splitPath[i]
	}
	if len(newSplitPath) == 1 {
		newPath = newSplitPath[0]
	} else {
		newPath = strings.Join(newSplitPath, "-")
	}

	if len(splitPath) == 1 {
		return BGPTree{Value: splitPath[0], Children: []BGPTree{}}
	} else {
		tree.Value = splitPath[0]
		tree.Children = append(tree.Children, insertSinglePathOnTree(tree, newPath))
		return tree
	}
}

//TODO: (17/11/2019) Need to finish the function to call this one
//This function was only inserting the verified paths (important, actually). Test it to see if it works
func insertVerifiedPathOnTree(tree BGPTree, path string) (BGPTree, bool) {

	splitPath := strings.Split(path, "-")
	var newPath string

	newSplitPath := make([]string, (len(splitPath) - 1))
	for i := 1; i < len(splitPath); i++ {
		newSplitPath[i-1] = splitPath[i]
	}
	if len(newSplitPath) == 1 {
		newPath = newSplitPath[0]
	} else {
		newPath = strings.Join(newSplitPath, "-")
	}

	//esta nas duas ultimas posicoes da arvore e vai verificar se pode inserir o path recebido
	//(se o ramo for maior que o path recebido, ele e inserido como irmão caso não seja igual, penso que sem problema)
	if len(splitPath) == 2 {
		if tree.Value == splitPath[0] {
			if len(tree.Children) == 0 {
				//pode inserir logo o splitPath[1]
				newChild := BGPTree{Value: splitPath[len(splitPath)-1], Children: []BGPTree{}}
				tree.Children = append(tree.Children, newChild)
				return tree, true
			} else {
				//tem de verificar se algum dos filhos e igual
				for i := 0; i < len(tree.Children); i++ {
					if tree.Children[i].Value == splitPath[1] {
						//nao insere porque ja existe um filho igual
						return tree, false
					} else if i == (len(tree.Children) - 1) {
						//todos os filhos sao diferentes e pode inserir como irmao deles
						newChild := BGPTree{Value: splitPath[len(splitPath)-1], Children: []BGPTree{}}
						tree.Children = append(tree.Children, newChild)
						return tree, true
					}
				}
				return tree, false
			}
		} else {
			return tree, false
		}
		//se o path recebido tiver mais do que duas posicoes, a recursao tem de continuar e ir reduzindo o path
	} else { //I CHANGED THIS, MIGHT NOT WORK
		childrenLength := len(tree.Children)
		if tree.Value == splitPath[0] && childrenLength > 0 {
			var auxTree BGPTree
			var actuallyModded bool
			//ciclo que percorre os filhos do no em questao e vai chamando a recursao para cada um deles
			for i := 0; i < childrenLength; i++ {
				insertAux, modAux := insertVerifiedPathOnTree(tree.Children[i], newPath)
				auxTree.Children = append(auxTree.Children, insertAux)
				if modAux {
					actuallyModded = true
				}
			}
			tree.Children = auxTree.Children
			return tree, actuallyModded
		}
		return tree, false
	}
	return tree, false
}

//Inserts all paths on a tree
func insertAllPathsOnTree(tree BGPTree, path string) BGPTree {
	splitPath := strings.Split(path, "-")
	var newPath string

	//This if verifies how long the path that was passed is and reduces it by one position. The first AS
	//in the path (splitPath[0]) is compared and the rest (newPath) is passed to the recursion.
	newSplitPath := make([]string, (len(splitPath) - 1))

	for i := 1; i < len(splitPath); i++ {
		newSplitPath[i-1] = splitPath[i]
	}

	//var newPath string
	if len(newSplitPath) == 1 {
		newPath = newSplitPath[0]
	} else {
		newPath = strings.Join(newSplitPath, "-")
	}

	//Stopping condition 1
	if len(splitPath) == 1 {
		if len(tree.Children) != 0 {
			for i := 0; i < len(tree.Children); i++ {
				if tree.Children[i].Value == splitPath[0] {
					return tree
				}
			}
		}
		newChild := BGPTree{Value: splitPath[0], Children: []BGPTree{}}
		tree.Children = append(tree.Children, newChild)
		return tree
	}

	//Stopping condition 2
	if len(tree.Children) == 0 {
		tree.Children = append(tree.Children, insertSinglePathOnTree(tree, newPath))
		return tree

		//If the node has children and the received path has more than one position,
		//then I need to compare them with the current AS on the path
	} else {
		var auxTree BGPTree
		var equalChild bool
		//Loop to go through all the node's children
		for i := 0; i < len(tree.Children); i++ {
			//If one of the children is the same as the AS on the path, continue the recursion
			if tree.Children[i].Value == splitPath[1] {
				equalChild = true
				auxTree.Children = append(auxTree.Children, insertAllPathsOnTree(tree.Children[i], newPath))
				//Append the children which aren't equal to the current AS, so their branches aren't lost
			} else {
				auxTree.Children = append(auxTree.Children, tree.Children[i])
			}
		}
		tree.Children = auxTree.Children
		//If no children was the same as the current AS, then it needs to be inserted as their sibling
		if !equalChild {
			newChild := BGPTree{Value: splitPath[1], Children: []BGPTree{}}
			tree.Children = append(tree.Children, insertSinglePathOnTree(newChild, newPath))
		}
		return tree
	}

}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryAnnouncementOnTree" {
		return s.queryAnnouncementOnTree(APIstub, args)
	} else if function == "queryAllTreeAnnouncements" {
		return s.queryAllTreeAnnouncements(APIstub)
	} else if function == "queryAllTreeAnnouncementsAlt" {
		return s.queryAllTreeAnnouncementsAlt(APIstub)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "announceVerifiedTreePath" {
		return s.announceVerifiedTreePath(APIstub, args)
	} else if function == "updateVerifiedPath" {
		return s.updateVerifiedPath(APIstub, args)
	} else if function == "registerComplaint" {
		return s.registerComplaint(APIstub, args)
	} else if function == "queryComplaint" {
		return s.queryComplaint(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

//TODO: think about how the V/U paths should be presented and update the function <- done but could be better
func (s *SmartContract) queryAnnouncementOnTree(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("-------------------START OF QUERYANNOUNCEMENT-------------------")
	if len(args) != 1 { //IP prefix to be queried
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	ipNoMask, ipNetwork, _ := net.ParseCIDR(args[0])
	realIP := ipNetwork.String()
	ipAddrCheck := checkIPAddressType(ipNoMask.String())

	var finalCompositeKey string

	if ipAddrCheck == -1 {
		return shim.Error("Incorrect IP address")
	} else if ipAddrCheck == 0 {
		//Create the composite key to insert in the blockchain
		ipAndMask := strings.Split(realIP, "/")    //[10.10.220.0] [17]
		octets := strings.Split(ipAndMask[0], ".") //[10] [10] [220] [0]
		compositePrefix := append(octets, ipAndMask[1])
		finalCompositeKey, _ = APIstub.CreateCompositeKey("IP", compositePrefix)
	} else if ipAddrCheck == 1 {
		//Fill the missing octets with zeroes
		realIPFull := handleIPv6Address(realIP)

		//Split the ip and mask for the composite key creation
		ipAndMask := strings.Split(realIPFull, "/") //[10.10.220.0] [17]
		octets := strings.Split(ipAndMask[0], ":")  //[10] [10] [220] [0]

		//Finish the actual composite key
		compositePrefix := append(octets, ipAndMask[1])
		finalCompositeKey, _ = APIstub.CreateCompositeKey("IP", compositePrefix)
	}

	announcementAsBytes, _ := APIstub.GetState(finalCompositeKey)
	storedTree := BGPTree{}
	json.Unmarshal(announcementAsBytes, &storedTree)

	var path string
	var allPaths = make([]string, nodeCounting(storedTree))
	allPaths = queryTreeImproved(storedTree, path, allPaths, 0, true)

	var finalPath string
	for i := 0; i < len(allPaths); i++ {
		if i != (len(allPaths) - 1) {
			finalPath += allPaths[i] + ";"
		} else {
			finalPath += allPaths[i]
		}
	}

	fmt.Println("finalPath is: " + finalPath)
	fmt.Println("real tree is: ")
	queryTree(storedTree, "", true)
	fmt.Println("-------------------END OF QUERYANNOUNCEMENT-------------------")
	return shim.Success([]byte(finalPath))
}

//This initLedger instatiates the ledger with paths stored as trees
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	announcements := []BGPTree{
		BGPTree{Value: "A", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "D", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "B", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "F", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "I", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "R", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "T", Status: "V", Children: []BGPTree{}},
	}

	announcementAsBytes0, _ := json.Marshal(announcements[0])
	APIstub.PutState("10.0.0.0/24", announcementAsBytes0)
	fmt.Println("Added 10.0.0.0/24:", announcements[0].Value)

	announcementAsBytes1, _ := json.Marshal(announcements[1])
	APIstub.PutState("192.168.0.0/16", announcementAsBytes1)
	fmt.Println("Added 192.168.0.0/16:", announcements[1].Value)

	announcementAsBytes2, _ := json.Marshal(announcements[2])
	APIstub.PutState("127.10.0.0/24", announcementAsBytes2)
	fmt.Println("Added 127.10.0.0/24:", announcements[2].Value)

	announcementAsBytes3, _ := json.Marshal(announcements[3])
	APIstub.PutState("200.0.0.0/8", announcementAsBytes3)
	fmt.Println("Added 200.0.0.0/8:", announcements[3].Value)

	announcementAsBytes4, _ := json.Marshal(announcements[4])
	APIstub.PutState("134.0.0.0/16", announcementAsBytes4)
	fmt.Println("Added 134.0.0.0/16:", announcements[4].Value)

	announcementAsBytes5, _ := json.Marshal(announcements[5])
	APIstub.PutState("156.0.0.0/8", announcementAsBytes5)
	fmt.Println("Added 156.0.0.0/8:", announcements[5].Value)

	announcementAsBytes6, _ := json.Marshal(announcements[6])
	APIstub.PutState("211.0.0.0/24", announcementAsBytes6)
	fmt.Println("Added 211.0.0.0/24:", announcements[6].Value)

	return shim.Success(nil)
}

//TODO (01/10/2019): Last function, hopefully. Compares the received path with the ones in the blockchain and stores it if
//it's verified and if it is better than one of the ones already stored (if worse, gets discarded, but still need to think about this)
//Apparently all I had to do was copy the other one and change the sub-function that it called (still need to test it though)
func (s *SmartContract) announceVerifiedTreePath(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Print("INICIO DO ANNOUNCE VERIFIED")
	//Key, which is the IP prefix, value, which is the AS path, and the counter for the current execution
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	//Check if the AS Path is in the correct form
	var pathReceived []string = strings.Split(args[1], "-")
	for i := 0; i < len(pathReceived); i++ {
		if len(pathReceived[i]) == 0 {
			return shim.Error("AS Path on the incorrect form")
		}
	}

	//Get the "real" network from the received IP prefix
	ipNoMask, ipNetwork, _ := net.ParseCIDR(args[0])
	realIP := ipNetwork.String()
	ipAddrCheck := checkIPAddressType(ipNoMask.String())

	var possiblePrefixes []pairKeyValue
	var finalCompositeKey string
	var alternateCompositeKey string
	var supernetBool bool
	var subnetBool bool

	if ipAddrCheck == -1 {
		return shim.Error("Incorrect IP address")
	} else if ipAddrCheck == 0 {
		//Create the composite key to insert in the blockchain
		ipAndMask := strings.Split(realIP, "/")                                  //[10.10.220.0] [17]
		octets := strings.Split(ipAndMask[0], ".")                               //[10] [10] [220] [0]
		compositePrefix := append(octets, ipAndMask[1])                          //[10] [10] [220] [0] [17]
		finalCompositeKey, _ = APIstub.CreateCompositeKey("IP", compositePrefix) //IP1010220017

		//Check mask list to see if there is subnet of the given prefix
		supernetBool, subnetBool = checkSubnet(APIstub, realIP, ipAddrCheck)
		if supernetBool && !subnetBool {
			return shim.Success([]byte("U | Prefix is a supernet of a previously stored prefix. Insertion impossible."))
		}

		//Alternate composite key to insert on the blockchain list of broader prefixes
		alternateCompositeKey, _ = APIstub.CreateCompositeKey("MASK", compositePrefix)

		//Check the compatible prefixes from all possible combinations
		possiblePrefixes = getAllPossiblePrefixesIPv4(APIstub, realIP)
	} else if ipAddrCheck == 1 {
		//Fill the missing octets with zeroes
		realIPFull := handleIPv6Address(realIP)
		fmt.Println("string after handleipv6: " + realIPFull)
		//Split the ip and mask for the composite key creation
		ipAndMask := strings.Split(realIPFull, "/") //[10.10.220.0] [17]
		octets := strings.Split(ipAndMask[0], ":")  //[10] [10] [220] [0]

		//Finish the actual composite key
		compositePrefix := append(octets, ipAndMask[1])
		finalCompositeKey, _ = APIstub.CreateCompositeKey("IP", compositePrefix)

		//Check mask list to see if there is subnet of the given prefix
		supernetBool, subnetBool = checkSubnet(APIstub, realIPFull, ipAddrCheck)
		if supernetBool && !subnetBool {
			s := []byte(args[0] + "»" + args[1] + "»U»" + args[2])
			APIstub.SetEvent("UpdateBGP", s)
			return shim.Success([]byte("U | Prefix is a supernet of a previously stored prefix. Insertion impossible."))
		}

		//Alternate composite key to check whether there's a subnet of the given prefix already on the blockchain
		alternateCompositeKey, _ = APIstub.CreateCompositeKey("MASK", compositePrefix)

		//Check the compatible prefixes from all possible combinations
		possiblePrefixes = getAllPossiblePrefixesIPv6(APIstub, realIPFull)
	}

	var verifiedNotInserted bool
	var verified bool
	var foundTree bool
	var finalTree BGPTree
	fmt.Println("possible prefixes: " + strconv.Itoa(len(possiblePrefixes)))
	//If there is no entry for that key or for any broader network of it
	if len(possiblePrefixes) == 0 {
		//Insert it if it has only one AS on the path (that AS is the owner of the prefix)
		if len(pathReceived) == 1 {
			finalTree = BGPTree{Value: pathReceived[0], Status: "V", Children: []BGPTree{}}
			verified = true
		} else { //Otherwise, nothing happens, announcement is discarded
			s := []byte(args[0] + "»" + args[1] + "»U»" + args[2])
			APIstub.SetEvent("UpdateBGP", s)
			return shim.Success([]byte("U | No corresponding prefix in the blockchain. Invalid path. Blockchain wasn't altered."))
		}
	} else {
		iteration := 0 //For debug purposes only (can delete this after everything works)
		for _, pair := range possiblePrefixes {
			foundTreeAsBytes := pair.Value
			storedTree := BGPTree{}
			json.Unmarshal(foundTreeAsBytes, &storedTree)

			fmt.Print("This was the iteration number: ")
			fmt.Println(iteration)

			fmt.Println("given key: " + finalCompositeKey)

			//Verifications
			if verifyPath(storedTree, args[1]) {
				fmt.Println("key that verified the path: " + pair.Key)
				verified = true
			}
			if pair.Key == finalCompositeKey {
				fmt.Println("key that matched the given one: " + pair.Key)
				finalTree = storedTree
				foundTree = true
			}

			//Insert the path on the given prefix's tree
			if verified && foundTree {
				fmt.Println("tree was already in the blockchain")
				finalTree = insertAllPathsOnTree(finalTree, args[1]) //TODO: PROBLEM IS HERE! NOT INSERTING FOR SOME REASON
				break
			}
			iteration++
		}

		//Create new tree with only the given path
		if verified && !foundTree {
			fmt.Println("tree wasn't in the blockchain")
			finalTree = insertSinglePathOnTree(finalTree, args[1])
		}

		//Check if inserted path is exactly the same as in the blockchain
		if !verified && foundTree {
			fmt.Println("path is the same as in the tree")
			var allPathsLocal = make([]string, nodeCounting(finalTree))
			var path string
			allPathsLocal = queryTreeImproved(finalTree, path, allPathsLocal, 0, true)
			fmt.Printf("%+q", allPathsLocal)
			fmt.Println(args[1])
			for i := 0; i < len(allPathsLocal); i++ {
				if allPathsLocal[i] == args[1] {
					verifiedNotInserted = true
				}
			}
		}
	}

	if verified {
		//Marshal the new tree and put it on the blockchain
		finalTreeAsBytes, _ := json.Marshal(finalTree)
		fmt.Println("Path was verified!!!")
		fmt.Println("FINAL INSERTION. KEY: " + finalCompositeKey)
		fmt.Println("FINAL INSERTION. TREE: " + string(finalTreeAsBytes))

		//Put ip prefix on the "alternate" list with only the broader prefixes
		if !supernetBool && !subnetBool {
			/*putErr1 := APIstub.PutState(alternateCompositeKey, []byte("0"))
			if putErr1 != nil {
				return shim.Error("First put error: " + putErr1.Error())
			}*/
			byteReturn := append([]byte("V | "+alternateCompositeKey+" | "+finalCompositeKey+" | "+args[0]+" | "+args[1]+" | "), finalTreeAsBytes...)
			return shim.Success(byteReturn)
		}

		//Put found key and updated tree in the blockchain
		/*putErr2 := APIstub.PutState(finalCompositeKey, finalTreeAsBytes) //TODO: MIGHT HAVE TO CHANGE THE KEY FROM REALIP TO ARGS[0]
		if putErr2 != nil {
			return shim.Error("Second put error: " + putErr2.Error())
		}*/

		//s := []byte(args[0] + "»" + args[1])
		//APIstub.SetEvent("BGP update", s)
		byteReturn := append([]byte("V | "+finalCompositeKey+" | "+args[0]+" | "+args[1]+" | "), finalTreeAsBytes...)
		return shim.Success(byteReturn)
	} else {
		if verifiedNotInserted {
			s := []byte(args[0] + "»" + args[1] + "»V»" + args[2])
			APIstub.SetEvent("UpdateBGP", s)
			return shim.Success([]byte("V | Blockchain wasn't altered."))
		} else {
			s := []byte(args[0] + "»" + args[1] + "»U»" + args[2])
			APIstub.SetEvent("UpdateBGP", s)
			return shim.Success([]byte("U | Blockchain wasn't altered."))
		}
	}
}

//After the previous function finds the actual prefix to be updated, this one simply updates the ledger
func (s *SmartContract) updateVerifiedPath(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	for i := 0; i < len(args); i++ {
		fmt.Println("argument " + strconv.Itoa(i) + " is " + args[i])
	}

	//Must insert the update, as well as the alternate composite key in the MASK list
	if len(args) == 7 {
		//Put ip prefix on the "alternate" list with only the broader prefixes
		putErr1 := APIstub.PutState(args[1], []byte("0"))
		if putErr1 != nil {
			return shim.Error("First put error: " + putErr1.Error())
		}

		//Put found key and updated tree in the blockchain
		putErr2 := APIstub.PutState(args[2], []byte(args[5]))
		if putErr2 != nil {
			return shim.Error("Second put error: " + putErr2.Error())
		}

		//Set event (prefix»path»verified»executionNum)
		s := []byte(args[3] + "»" + args[4] + "»" + args[0] + "»" + args[6])
		APIstub.SetEvent("UpdateBGP", s)

		return shim.Success([]byte(args[0] + " | Blockchain was updated!"))

	} else {
		//Put found key and updated tree in the blockchain
		putErr := APIstub.PutState(args[1], []byte(args[4]))
		if putErr != nil {
			return shim.Error("Put error: " + putErr.Error())
		}

		//Set event (prefix»path»verified»executionNum)
		s := []byte(args[2] + "»" + args[3] + "»" + args[0] + "»" + args[5])
		APIstub.SetEvent("UpdateBGP", s)

		return shim.Success([]byte(args[0] + " | Blockchain was updated!"))
	}
}

//Is it done? Seems too simple. Ask teachers what else it needs to do
//TODO: also, still need to test it better
func (s *SmartContract) registerComplaint(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//Key, which is the IP prefix and the AS path, and value, which is the AS that is making the complaint
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//Check if the IP prefix and AS Path are in the correct form
	var prefixAndPath []string = strings.Split(args[0], ":")
	if prefixAndPath[0] == "" {
		return shim.Error("Incorrect IP prefix")
	}
	var pathReceived []string = strings.Split(prefixAndPath[1], "-")
	for i := 0; i < len(pathReceived); i++ {
		if len(pathReceived[i]) == 0 {
			return shim.Error("AS Path on the incorrect form")
		}
	}

	//Get the announcements for the same IP prefix that are already stored
	announcementAsBytes, _ := APIstub.GetState(args[0])
	complaint := Complaint{}
	if announcementAsBytes == nil {
		//if timestamp becomes an entry, add it here as well
		complaint.ASList = args[1]
		complaintAsBytes, _ := json.Marshal(complaint)
		APIstub.PutState(args[0], complaintAsBytes)
	} else {
		json.Unmarshal(announcementAsBytes, &complaint)
		var asList []string = strings.Split(complaint.ASList, "-")
		for i := 0; i < len(asList); i++ {
			if asList[i] == args[1] {
				//This means the AS already filed the complaint, as it is already on the list
				return shim.Success(nil)
			} else {
				if i == (len(asList) - 1) {
					//This means it reached the last position and the AS isn't on the list yet, so add it
					complaint.ASList += ("-" + args[1])
					complaintAsBytes, _ := json.Marshal(complaint)
					APIstub.PutState(args[0], complaintAsBytes)
				}
			}
		}
	}
	return shim.Success(nil)
}

func (s *SmartContract) queryComplaint(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if args[0] == "" {
		return shim.Error("Incorrect prefix and path")
	}
	complaintAsBytes, _ := APIstub.GetState(args[0])
	complaint := Complaint{}
	if complaintAsBytes == nil {
		return shim.Success(nil)
	} else {
		json.Unmarshal(complaintAsBytes, &complaint)
		println(args[0] + " -> " + complaint.ASList)
	}
	return shim.Success(nil)
}

//Apparently working. Should test it further and think about how this will work with the complaints also on the blockchain
func (s *SmartContract) queryAllTreeAnnouncements(APIstub shim.ChaincodeStubInterface) sc.Response {

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	resultsIterator, err := APIstub.GetStateByPartialCompositeKey("IP", []string{})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	//Loop to go through every IP prefix stored
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		_, splitPrefix, _ := APIstub.SplitCompositeKey(queryResponse.GetKey())
		var recoveredPrefix string
		if len(splitPrefix) == 5 {
			recoveredPrefix = reconstructPrefix(splitPrefix, 0)
		} else if len(splitPrefix) == 9 {
			recoveredPrefix = reconstructPrefix(splitPrefix, 1)
		} else {
			return shim.Error("Corrupted IP prefix...")
		}

		storedTree := BGPTree{}
		json.Unmarshal(queryResponse.Value[:], &storedTree)

		var path string
		var allPaths = make([]string, nodeCounting(storedTree))
		allPaths = queryTreeImproved(storedTree, path, allPaths, 0, true)

		//Just to abbreviate the IP prefix, in case it's an IPv6
		_, finalRecoveredPrefix, _ := net.ParseCIDR(recoveredPrefix)

		//IP prefix contains multiple paths that have to be printed individually
		if len(allPaths) > 1 {
			buffer.WriteString(finalRecoveredPrefix.String())
			buffer.WriteString("»")
			for i := 0; i < len(allPaths); i++ {
				buffer.WriteString(allPaths[i])
				buffer.WriteString(";")
			}
			buffer.WriteString("\n")
			//IP prefix only has one path
		} else {
			buffer.WriteString(finalRecoveredPrefix.String())
			buffer.WriteString("»")
			buffer.WriteString(allPaths[0]) //if this doesn't work, use queryResponse.Value
			buffer.WriteString("\n")
		}
	}

	return shim.Success(buffer.Bytes())
}

//Alternative that prints the MASK list, as well as the IP list
func (s *SmartContract) queryAllTreeAnnouncementsAlt(APIstub shim.ChaincodeStubInterface) sc.Response {

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	//Monkey business to run the iterator loop twice, once for each partial key list <- REMOVE WHEN CORRECT
	for x := 0; x < 2; x++ {
		var partialKey string
		if x == 0 {
			partialKey = "MASK"
		} else {
			partialKey = "IP"
		}
		//resultsIterator, err := APIstub.GetStateByRange("", "")
		resultsIterator, err := APIstub.GetStateByPartialCompositeKey(partialKey, []string{})
		//resultsIterator, err := APIstub.GetStateByPartialCompositeKey("IP", []string{})
		if err != nil {
			return shim.Error(err.Error())
		}
		defer resultsIterator.Close()

		fmt.Println("BEFORE THE LOOP TO GO THROUGH ALL BC ENTRIES")
		//Loop to go through every IP prefix stored
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
			if err != nil {
				return shim.Error(err.Error())
			}

			_, splitPrefix, _ := APIstub.SplitCompositeKey(queryResponse.GetKey())
			fmt.Println("INSIDE THE LOOP, BEFORE THE RECOVERED PREFIX IFS")
			var recoveredPrefix string
			if len(splitPrefix) == 5 {
				recoveredPrefix = reconstructPrefix(splitPrefix, 0)
			} else if len(splitPrefix) == 9 {
				recoveredPrefix = reconstructPrefix(splitPrefix, 1)
			} else {
				return shim.Error("Corrupted IP prefix...")
			}

			fmt.Println("CHECKED KEY IS:" + queryResponse.Key)

			storedTree := BGPTree{}
			json.Unmarshal(queryResponse.Value[:], &storedTree)

			var path string
			var allPaths = make([]string, nodeCounting(storedTree))
			allPaths = queryTreeImproved(storedTree, path, allPaths, 0, true)

			//Just to abbreviate the IP prefix, in case it's an IPv6
			_, finalRecoveredPrefix, _ := net.ParseCIDR(recoveredPrefix)
			fmt.Println("FINAL RECOVERED PREFIX IS: " + finalRecoveredPrefix.String())
			if x == 0 {
				buffer.WriteString(finalRecoveredPrefix.String())
				buffer.WriteString("»")
				buffer.WriteString("0")
				buffer.WriteString("\n")
			} else {

				//IP prefix contains multiple paths that have to be printed individually
				if len(allPaths) > 1 {
					buffer.WriteString(finalRecoveredPrefix.String())
					buffer.WriteString("»")
					for i := 0; i < len(allPaths); i++ {
						buffer.WriteString(allPaths[i])
						buffer.WriteString(";")
					}
					buffer.WriteString("\n")
					//IP prefix only has one path
				} else {
					buffer.WriteString(finalRecoveredPrefix.String())
					buffer.WriteString("»")
					buffer.WriteString(allPaths[0]) //if this doesn't work, use queryResponse.Value
					buffer.WriteString("\n")
				}
			}
		}
	}
	fmt.Printf("- queryAllTreeAnnouncements:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
