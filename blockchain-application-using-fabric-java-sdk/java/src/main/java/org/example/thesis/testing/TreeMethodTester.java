package main.java.org.example.thesis.testing;

import java.util.ArrayList;
import main.java.org.example.thesis.data_structures.BGPTree;

public class TreeMethodTester {

	public static void main(String[] args) {

		BGPTree testTree0 = new BGPTree();

		BGPTree testTree1 = new BGPTree("A", "V");

		BGPTree testTree2_1 = new BGPTree("B", "U");
		BGPTree testTree2_2 = new BGPTree("C", "U");
		ArrayList<BGPTree> aux2 = new ArrayList<BGPTree>();
		aux2.add(testTree2_1);
		aux2.add(testTree2_2);
		BGPTree testTree2 = new BGPTree("A", "U", aux2);

		ArrayList<BGPTree> aux3 = new ArrayList<BGPTree>(); // A has children B and C
		aux3.add(new BGPTree("B", "U"));
		aux3.add(new BGPTree("C", "U"));
		ArrayList<BGPTree> aux3_1 = new ArrayList<BGPTree>(); // B has children D
		aux3_1.add(new BGPTree("D", "U"));
		aux3.get(0).setChildren(aux3_1);
		ArrayList<BGPTree> aux3_2 = new ArrayList<BGPTree>(); // D has children E
		aux3_2.add(new BGPTree("E", "U"));
		aux3.get(0).getChildren().get(0).setChildren(aux3_2);
		BGPTree testTree3 = new BGPTree("A", "U", aux3);

		ArrayList<BGPTree> aux4 = new ArrayList<BGPTree>(); // A has children B and C
		aux4.add(new BGPTree("B", "U"));
		aux4.add(new BGPTree("C", "U"));
		ArrayList<BGPTree> aux4_1 = new ArrayList<BGPTree>(); // B has children D
		aux4_1.add(new BGPTree("D", "U"));
		aux4.get(0).setChildren(aux4_1);
		ArrayList<BGPTree> aux4_2 = new ArrayList<BGPTree>(); // D has children E and F
		aux4_2.add(new BGPTree("E", "U"));
		aux4_2.add(new BGPTree("F", "U"));
		aux4.get(0).getChildren().get(0).setChildren(aux4_2);
		ArrayList<BGPTree> aux4_3 = new ArrayList<BGPTree>(); // C has children G and H
		aux4_3.add(new BGPTree("G", "U"));
		aux4_3.add(new BGPTree("H", "U"));
		aux4.get(1).setChildren(aux4_3);
		BGPTree testTree4 = new BGPTree("A", "U", aux4);

		// Method 1 test (works correctly)
		System.out.println(BGPTree.nodeCounting(testTree0));
		System.out.println(BGPTree.nodeCounting(testTree1));
		System.out.println(BGPTree.nodeCounting(testTree2));
		System.out.println(BGPTree.nodeCounting(testTree3));
		System.out.println(BGPTree.nodeCounting(testTree4));

		// Method 2 test (works correctly)
		String path0 = new String();
		String path1 = new String();
		String path2 = new String();
		String path3 = new String();
		String path4 = new String();

		System.out.println("Tree 0");
		BGPTree.queryTree(testTree0, path0, true);
		System.out.println("Tree 1");
		BGPTree.queryTree(testTree1, path1, true);
		System.out.println("Tree 2");
		BGPTree.queryTree(testTree2, path2, true);
		System.out.println("Tree 3");
		BGPTree.queryTree(testTree3, path3, true);
		System.out.println("Tree 4");
		BGPTree.queryTree(testTree4, path4, true);

		System.out.println("The end");

		// Method 3 test (works with all the trees)
		String path4_1 = new String();
		String[] pathList = new String[BGPTree.nodeCounting(testTree1)];
		String[] finalPathList = BGPTree.queryTreeImproved(testTree1, path4_1, pathList, 0, true);

		for (int i = 0; i < finalPathList.length; i++) {
			System.out.println(finalPathList[i]);
		}

		System.out.println("Finally: ");

		// Method 4 test (test better)
		// TODO: the case where I'm trying to insert a path bigger than 2 on an existing
		// path with 2 or less, and
		// while having a fork on the 3rd position, is not working properly (example:
		// trying to insert "A-C-R-E-L"
		// While having "A-C-D", being "A" and "C" verified, and "D" unverified). Figure
		// out why and fix it!

		/*
		 * testTree1 = BGPTree.insertAllPathsOnTree(testTree1, "A-C", testTree1.status);
		 * testTree1 = BGPTree.insertAllPathsOnTree(testTree1, "A-B", testTree1.status);
		 * testTree1 = BGPTree.insertAllPathsOnTree(testTree1, "A-C-D",
		 * testTree1.status); testTree1 = BGPTree.insertAllPathsOnTree(testTree1,
		 * "A-C-R", testTree1.status); testTree1 =
		 * BGPTree.insertAllPathsOnTree(testTree1, "A-C-R-E", testTree1.status);
		 * testTree1 = BGPTree.insertAllPathsOnTree(testTree1, "A-C-R-E-L-K",
		 * testTree1.status);
		 */

		// TRYING WITHOUT THE LASTSTATUS PARAM
		testTree1 = BGPTree.insertAllPathsOnTree(testTree1, "A-C");
		testTree1 = BGPTree.insertAllPathsOnTree(testTree1, "A-B");
		testTree1 = BGPTree.insertAllPathsOnTree(testTree1, "A-C-D");
		testTree1 = BGPTree.insertAllPathsOnTree(testTree1, "A-C-R");
		testTree1 = BGPTree.insertAllPathsOnTree(testTree1, "A-C-R-E");
		testTree1 = BGPTree.insertAllPathsOnTree(testTree1, "A-C-R-E-L-K-F-X");
		System.out.println("New Tree ");
		String path5_1 = new String();
		String[] pathList5_1 = new String[BGPTree.nodeCounting(testTree1)];
		String[] finalPathList5_1 = BGPTree.queryTreeImproved(testTree1, path5_1, pathList5_1, 0, true);

		for (int i = 0; i < finalPathList5_1.length; i++) {
			System.out.println(finalPathList5_1[i]);
		}

		System.out.println();
		System.out.println("TEST WITH ANOTHER FUNCTION");
		BGPTree.printTree(testTree1, true);

	}
}
