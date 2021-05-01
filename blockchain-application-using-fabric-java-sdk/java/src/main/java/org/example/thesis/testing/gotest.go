package main

import (
	"fmt"
	"net"
	"strings"
)

// BGPTree...
type BGPTree struct {
	Value    string    `json:"value"`
	Status   string    `json:"status"`
	Children []BGPTree `json:"children"`
}

func main() {
	var brokenPrefix []string = []string{"100", "200", "69", "134", "8", "0", "0", "12", "64"}
	fmt.Println(reconstructPrefix(brokenPrefix, 1))

	a, b, _ := net.ParseCIDR(reconstructPrefix(brokenPrefix, 1))
	fmt.Println(a.String())
	fmt.Println(b.String())

	//var tree BGPTree = BGPTree{Value: "A", Children: []BGPTree{}}

	/*
		tree, _ = insertVerifiedPathOnTree(tree, "A-B")
		tree, _ = insertVerifiedPathOnTree(tree, "A-D")
		tree, _ = insertVerifiedPathOnTree(tree, "A-B-C")
		tree, _ = insertVerifiedPathOnTree(tree, "A-B-C-E")


		tree, _ = insertVerifiedPathOnTree(tree, "A-E")
		tree = insertAllPathsOnTree(tree, "A-B-C")
		tree = insertAllPathsOnTree(tree, "A-L-K-S")
		tree = insertAllPathsOnTree(tree, "A-L-K-S-T")
		tree = insertAllPathsOnTree(tree, "A-R")
		tree = insertAllPathsOnTree(tree, "H-G-J")
		tree = insertAllPathsOnTree(tree, "A-L-K-X-U-V")

		queryTree(tree, "", true)

		fmt.Println(verifyPath(tree, "A-B-C-E-F"))
		fmt.Println(verifyPath(tree, "A-R"))
		fmt.Println(verifyPath(tree, "A-L-K-X-U-V-T"))
		fmt.Println(verifyPath(tree, "A-B-G-H"))
		fmt.Println(verifyPath(tree, "A-D-K-L"))
		fmt.Println(verifyPath(tree, "A-D-K"))
		fmt.Println(verifyPath(tree, "F"))

	*/

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

//Verifies a path based on a given tree
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
