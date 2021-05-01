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
	"net"
	"strings"
)

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

func insertVerifiedPathOnTree(tree BGPTree, path string) BGPTree {

	splitPath := strings.Split(path, "-")
	var newPath string

	if len(path) > 1 {
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
	}

	//esta nas duas ultimas posicoes da arvore e vai verificar se pode inserir o path recebido
	//(se o ramo for maior que o path recebido, ele e inserido como irmão caso não seja igual, penso que sem problema)
	if len(splitPath) == 2 {
		if tree.Value == splitPath[0] {
			if len(tree.Children) == 0 {
				//pode inserir logo o splitPath[1]
				newChild := BGPTree{Value: splitPath[len(splitPath)-1], Children: []BGPTree{}}
				tree.Children = append(tree.Children, newChild)
				return tree
			} else {
				//tem de verificar se algum dos filhos e igual
				for i := 0; i < len(tree.Children); i++ {
					if tree.Children[i].Value == splitPath[1] {
						//nao insere porque ja existe um filho igual
						return tree
					} else if i == (len(tree.Children) - 1) {
						//todos os filhos sao diferentes e pode inserir como irmao deles
						newChild := BGPTree{Value: splitPath[len(splitPath)-1], Children: []BGPTree{}}
						tree.Children = append(tree.Children, newChild)
						return tree
					}
				}
				return tree
			}
		} else {
			return tree
		}
		//se o path recebido tiver mais do que duas posicoes, a recursao tem de continuar e ir reduzindo o path
	} else {
		childrenLength := len(tree.Children)
		if tree.Value == splitPath[0] && childrenLength > 0 {
			var auxTree BGPTree
			//ciclo que percorre os filhos do no em questao e vai chamando a recursao para cada um deles
			for i := 0; i < childrenLength; i++ {
				auxTree.Children = append(auxTree.Children, insertVerifiedPathOnTree(tree.Children[i], newPath))
			}
			tree.Children = auxTree.Children
			return tree
		}
		return tree
	}
	return tree
	/*}
	return tree*/
}

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

func main() {

	prefixToCompare := "10.10.0.0/19"
	pathToCompare := "X"

	tests := [5]string{"10.10.220.0/17", "10.0.0.0/8", "10.10.0.0/20", "10.10.0.0/18", "142.10.220.0/16"}

	storedTrees := []BGPTree{
		BGPTree{Value: "A", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "D", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "B", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "F", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "I", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "R", Status: "V", Children: []BGPTree{}},
		BGPTree{Value: "T", Status: "V", Children: []BGPTree{}},
	}

	var prefixes map[string]BGPTree
	prefixes = make(map[string]BGPTree)

	for i := 0; i < 5; i++ {
		prefix := tests[i]
		_, realPrefix, _ := net.ParseCIDR(prefix)
		prefixes[realPrefix.String()] = storedTrees[i]
	}

	//insertVerifiedPathOnTree(storedTree, "A")
	//storedTree = insertVerifiedPathOnTree(storedTree, "A-B")
	//storedTree = insertVerifiedPathOnTree(storedTree, "A-C")
	//storedTree = insertVerifiedPathOnTree(storedTree, "A-B-D")

	/*var buffer bytes.Buffer
	var path string
	var allPaths = make([]string, nodeCounting(storedTree))
	allPaths = queryTreeImproved(storedTree, path, allPaths, 0, true)

	//IP prefix contains multiple paths that have to be printed individually
	buffer.WriteString("10.0.0.0/24")
	buffer.WriteString("»")
	for i := 0; i < len(allPaths); i++ {
		buffer.WriteString(allPaths[i])
		buffer.WriteString(";")
	}
	buffer.WriteString("\n")

	fmt.Println(buffer.String())*/
}
