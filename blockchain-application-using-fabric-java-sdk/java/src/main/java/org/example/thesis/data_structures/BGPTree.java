package main.java.org.example.thesis.data_structures;

import java.util.ArrayList;

//TODO: EVERYTHING SEEMS TO BE WORKING FINE. FURTHER TESTING REQUIRED. CLEAN THE CODE. START DOING THE SHIT THAT I SHOULD'VE BEEN DOING A LONG TIME AGO.
//TODO: ADAPT THE TREE PRINTING FUNCTION TO SHOW EACH NODE'S STATUS

public class BGPTree {

	public String value;
	public String status;
	public ArrayList<BGPTree> children;

	public BGPTree() {
		this.value = "";
		this.status = "";
		this.children = new ArrayList<BGPTree>();
	}

	public BGPTree(String value, String status) {
		this.value = value;
		this.status = status;
		this.children = new ArrayList<BGPTree>();
	}

	public BGPTree(String value, String status, ArrayList<BGPTree> children) {
		this.value = value;
		this.status = status;
		this.children = children;
	}

	public String getValue() {
		return value;
	}

	public void setValue(String value) {
		this.value = value;
	}

	public String getStatus() {
		return status;
	}

	public void setStatus(String status) {
		this.status = status;
	}

	public ArrayList<BGPTree> getChildren() {
		return children;
	}

	public void setChildren(ArrayList<BGPTree> children) {
		this.children = children;
	}

	// Returns the number of leaf nodes the tree has (Checked on golang, checked
	// here)
	public static int nodeCounting(BGPTree tree) {
		int childrenCounter = 0;

		if (tree.equals(null) || tree.value.equals("")) {
			return 0;
		}

		if (tree.children.size() == 0) {
			return 1;
		}

		for (int i = 0; i < tree.children.size(); i++) {
			childrenCounter += nodeCounting(tree.children.get(i));
		}

		return childrenCounter;
	}

	// Pulled out of my ass, don't use this
	public static void printTree(BGPTree tree, boolean firstTime) {

		if (firstTime == true)
			System.out.println(tree.value);

		if (tree.children.size() == 0)
			return;

		for (int i = 0; i < tree.children.size(); i++) {
			System.out.print(tree.children.get(i).value + "     ");
		}
		System.out.println();
		for (int j = 0; j < tree.children.size(); j++) {
			printTree(tree.children.get(j), false);
		}
	}

	// Prints the paths of the tree, one on each line
	public static void queryTree(BGPTree tree, String path, boolean firstTime) {
		if (firstTime) {
			path += tree.value;
			firstTime = false;
		} else {
			path += "-";
			path += tree.value;
		}

		if (tree.children.size() == 0) {
			System.out.println(path);
			return;
		} else {
			for (int i = 0; i < tree.children.size(); i++) {
				queryTree(tree.children.get(i), path, firstTime);
			}
			return;
		}
	}

	// Returns a pathList which has a different path separated by "-" in each of the
	// array's positions
	public static String[] queryTreeImproved(BGPTree tree, String path, String[] pathList, int index,
			boolean firstTime) {
		if (tree.value.equals("")) {
			return pathList;
		}
		if (firstTime) {
			path += tree.value;
			firstTime = false;
		} else {
			path += "-";
			path += tree.value;
		}
		if (tree.children.size() == 0) {
			// pathList[index] = path; //<- This is how it previously was and it worked.
			// Replace the line below with this one
			pathList[index] = tree.status + ":" + path;
			return pathList;
		} else {
			if (tree.children.size() == 1) {
				pathList = queryTreeImproved(tree.children.get(0), path, pathList, index, firstTime);
			} else {
				int excessLeaves = 0;
				for (int i = 0; i < tree.children.size(); i++) {
					if (i != 0) {
						excessLeaves += nodeCounting(tree.children.get(i - 1));
						pathList = queryTreeImproved(tree.children.get(i), path, pathList, index + excessLeaves,
								firstTime);
					} else {
						pathList = queryTreeImproved(tree.children.get(i), path, pathList, index, firstTime);
					}
				}
			}
		}
		return pathList;
	}

	// Inserts all paths on a tree, being Verified or Unverified (differentiates
	// them)
	public static BGPTree insertAllPathsOnTree(BGPTree tree, String path) {
		String[] splitPath = path.split("-");
		String newPath = null;

		// This if verifies how long the path that was passed is and reduces it by one
		// position. The first AS
		// in the path (splitPath[0]) is compared and the rest (newPath) is passed to
		// the recursion.
		if (splitPath.length > 1) {
			String[] newSplitPath = new String[splitPath.length - 1];

			for (int i = 1; i < splitPath.length; i++) {
				newSplitPath[i - 1] = splitPath[i];
			}

			if (newSplitPath.length == 1) {
				newPath = newSplitPath[0];
			} else {
				newPath = String.join("-", newSplitPath);
			}
		}

		// Status to be passed to the next nodes
		String status = tree.status;

		// This if covers the last case to be checked, since if the first position in
		// the AS path checks out,
		// the next one either already exists and nothing changes or doesn't and needs
		// to be inserted.
		if (splitPath.length == 2) {
			// Checks if the value of the node we're currently on matches the first one on
			// the received path
			if (tree.value.equals(splitPath[0])) {
				// If it does and it has no children, then splitPath[1] can be inserted on the
				// tree automatically
				if (tree.children.size() == 0) {
					// Can insert splitPath[1] immediately
					BGPTree newChild = new BGPTree(splitPath[splitPath.length - 1], status, new ArrayList<BGPTree>());
					tree.children.add(newChild);
					return tree;
					// If it has children, splitPath[1] has to be compared with the node's children
					// to see if it's already there
				} else {
					// Cycle to loop through all of the children
					for (int i = 0; i < tree.children.size(); i++) {
						// If it matches any one of them, it means that AS is already stored on the tree
						if (tree.children.get(i).value.equals(splitPath[1])) {
							return tree;
							// Else, none of the children match said AS and it can be inserted as their
							// brother
						} else if (i == (tree.children.size() - 1)) {
							BGPTree newChild = new BGPTree(splitPath[splitPath.length - 1], status,
									new ArrayList<BGPTree>());
							tree.children.add(newChild);
							return tree;
						}
					}
					return tree;
				}
				// If it doesn't match, just return the tree to the calls above (so the
				// recursion doesn't lose any branches)
			} else {
				return tree;
			}
			// If it's not checking the last two positions on the received path, the
			// recursion needs to keep going down the tree
		} else {
			int childrenLength = tree.children.size();
			boolean atLeastOne = false;
			// If current node's value matches the the first one on the received path
			if (tree.value.equals(splitPath[0])) {
				BGPTree auxTree = new BGPTree();
				// If it has no children, then the path needs to be inserted as "Unverified"
				if (childrenLength == 0) {
					BGPTree newChild = new BGPTree(splitPath[1], "U", new ArrayList<BGPTree>());
					tree.children.add(newChild);
					auxTree.children.add(insertAllPathsOnTree(tree.children.get(0), newPath));
					tree.children = auxTree.children;
					// If it has children, then the children need to be checked to see where to
					// proceed with the recursion
				} else {
					// Loop to go through all of the children and call the function recursively
					for (int i = 0; i < childrenLength; i++) {
						// Checks if there's one child node that's equal to the next value on the
						// received path
						if (splitPath[1].equals(tree.children.get(i).value)) {
							atLeastOne = true;
						}
						auxTree.children.add(insertAllPathsOnTree(tree.children.get(i), newPath));
					}

					// No child nodes are equal to the next one on the received path, so it needs to
					// be inserted as their brother and keep the recursion going there
					if (!atLeastOne && splitPath.length > 2) {
						BGPTree newChild = new BGPTree(splitPath[1], "U", new ArrayList<BGPTree>());
						tree.children.add(newChild);
						auxTree.children.add(insertAllPathsOnTree(tree.children.get(childrenLength), newPath));
					}
					tree.children = auxTree.children;
				}
			}
			return tree;
		}
	}

	// Inserts the verified paths on a tree
	public BGPTree insertVerifiedPathOnTree(BGPTree tree, String path) {

		String[] splitPath = path.split("-");
		String newPath = null;

		if (path.length() > 1) {
			String[] newSplitPath = new String[splitPath.length - 1];

			for (int i = 1; i < splitPath.length; i++) {
				newSplitPath[i - 1] = splitPath[i];
			}

			if (newSplitPath.length == 1) {
				newPath = newSplitPath[0];
			} else {
				newPath = String.join("-", newSplitPath);
			}
		}

		// esta nas duas ultimas posicoes da arvore e vai verificar se pode inserir o
		// path recebido
		// (se o ramo for maior que o path recebido, ele e inserido como irmão caso não
		// seja igual, penso que sem problema)
		if (splitPath.length == 2) {
			if (tree.getValue() == splitPath[0]) {
				if (tree.getChildren().size() == 0) {
					// If it has no children, can insert splitPath[1] right away
					BGPTree newChild = new BGPTree(splitPath[splitPath.length - 1], "", new ArrayList<BGPTree>());
					tree.children.add(newChild);
					return tree;
				} else {
					// Checks if any of the children is the same
					for (int i = 0; i < tree.getChildren().size(); i++) {
						if (tree.children.get(i).getValue().equals(splitPath[1])) {
							// Doesn't insert because there is already the same child
							return tree;
						} else if (i == tree.getChildren().size() - 1) {
							// Every child is different and can insert the new one as their brother
							BGPTree newChild = new BGPTree(splitPath[splitPath.length - 1], "",
									new ArrayList<BGPTree>());
							tree.children.add(newChild);
							return tree;
						}
					}
					return tree;
				}
			} else {
				return tree;
			}
			// se o path recebido tiver mais do que duas posicoes, a recursao tem de
			// continuar e ir reduzindo o path
		} else {
			int childrenLength = tree.getChildren().size();
			if (tree.getValue().equals(splitPath[0]) && childrenLength > 0) {
				BGPTree auxTree = new BGPTree();
				// ciclo que percorre os filhos do no em questao e vai chamando a recursao para
				// cada um deles
				for (int i = 0; i < childrenLength; i++) {
					auxTree.children.add(insertVerifiedPathOnTree(tree.children.get(i), newPath));
				}
				tree.setChildren(auxTree.getChildren());
				return tree;
			}
			return tree;
		}
	}
}
