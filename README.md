# Incremental Merkle Tree

This repository contains a Go implementation of an Incremental Merkle Tree (IMTs).

## The Algorithm

IMTs are perfect (balanced) binary trees. They allow us to recalcluate tree roots in polynomial time when new leaves are added. This is achieved by utilizing two, constant-sized slices of digests:
* the zero digest slice which is created on initialization and never updated; and
* the left-node digest slice which is built up from left-node digests as we add leaves and calculate digests towards the root.

The following depicts what the IMT would look like if we replaced `hash(left, right)` with `cat(left, right)` for the purposes of visualization.

On creation, all the leaves are zeroes. At this point, we have a slice which contains a single digest for every level of the tree. The tree has a height of 3 so our slice has 3 elements.
```mermaid
flowchart
subgraph level: 2
6(0000)
end
subgraph level: 1
4(00)
5(00)
end
subgraph level: 0
0(0)
1(0)
2(0)
3(0)
end
6(0000) --- 4(00)
6(0000) --- 5(00)
4(00) --- 0(0)
4(00) --- 1(0)
5(00) --- 2(0)
5(00) --- 3(0)
```

When we add a leaf (A), we recalculate the tree from the bottom-up. As we do this, we maintain a list of the digests of the left-sided nodes (A, and A0).
```mermaid
flowchart TD
subgraph Recalculated
6(A000) --- 4(00)
end
subgraph Added
4(A0) --- 0(A)
end
6(A000) --- 5(00)
4(A0) --- 1(0)
5(00) --- 2(0)
5(00) --- 3(0)
```

Adding another leaf (B) will cause the same subtree to be recalculated. We will rely on the list of left-sided node digests calculated from the previous step when leaf A was added.

```mermaid
flowchart TD
subgraph Recalculated
6(AB00) --- 4(00)
end
4(A0) --- 0(A)
6(AB00) --- 5(00)
subgraph Added
4(AB) --- 1(B)
end
5(00) --- 2(0)
5(00) --- 3(0)
```

After adding 4 leaves (A, B, C, and D), the tree has been completely recalculated.
```mermaid
flowchart TD
6(ABCD) --- 4(AB)
6(ABCD) --- 5(CD)
4(AB) --- 0(A)
4(AB) --- 1(B)
5(CD) --- 2(C)
5(CD) --- 3(D)
```
