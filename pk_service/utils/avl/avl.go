package avl

import (
	"bytes"
	"fmt"
	"server/shared/queue"
)

type AvlTree struct {
	root       *AVLTreeNode
	nodeMp     map[queue.UserId]*AVLTreeNode // 记录名字和
	debugPrint bool
}
type AVLTreeNode struct {
	name   string // debug
	height int32  //节点高
	value  interface{}
	parent *AVLTreeNode // 父节点
	left   *AVLTreeNode //节点左儿子
	right  *AVLTreeNode //节点右儿子
	size   int32        // 子树大小
}

// 插入谁之前
func (a *AvlTree) Add(other, you queue.UserId, value interface{}) error {
	var val *AVLTreeNode
	var valEx bool

	newVal := &AVLTreeNode{
		height: 1,
		value:  value,
		size:   1,
	}
	if a.root == nil {
		a.root = newVal
		a.nodeMp = map[queue.UserId]*AVLTreeNode{}
	} else {
		if node, ok := a.nodeMp[other]; ok {
			if val, valEx = a.nodeMp[you]; valEx { // 如果 已经存在 先删除 在 加入
				if a.debugPrint {
					println("before after-----------------")
					println(a.String())
				}
				a.root = val.delete()
				// debug print 这个时为了调试用的
				if a.debugPrint {
					for j := 0; j < len(a.nodeMp); j++ {
						q := a.nodeMp[queue.UserId(j)]
						if q == nil || q == val {
							continue
						}
						if q.parent != nil && q.parent.left != q && q.parent.right != q {
							panic("")
						}
						if node.getBF(false) == 2 {
							println("----------------")
							println(a.String())
							panic("")
						}
					}
					println("delete after-----------------")
					println(a.String())
				}
			}
			a.root = a.root.insertByRank(node.getRank(), newVal)
		} else {
			return fmt.Errorf("没有other")
		}
	}
	a.nodeMp[you] = newVal
	return nil
}

// 遍历一段区间
func (t *AvlTree) Seg(low, high int32) []*AVLTreeNode {
	// TODO: 参数合法化
	return t.root.seg(low, high, make([]*AVLTreeNode, 0))
}
func (a *AvlTree) GetRank(id queue.UserId) (int32, error) {
	if node, ok := a.nodeMp[id]; ok {
		return node.getRank(), nil
	}
	return 0, fmt.Errorf("User not Exist")
}
func (n *AvlTree) DisplayNodesInOrder() {
	n.root.displayNodesInOrder()
}

// 树形结构图
func (t *AvlTree) String() string {
	if t.root == nil {
		return ""
	}
	return string(t.print(t.root, 0, true))
}

// 插入谁之前
func (a *AvlTree) Delete(you queue.UserId) {
	if a.root == nil {
		panic("")
	}
	if node, ok := a.nodeMp[you]; ok {
		a.root = node.delete()
		delete(a.nodeMp, you)
	}
}

//getSize 返回子树大小
func getSize(a *AVLTreeNode) int32 {
	if a == nil {
		return 0
	}
	return a.size
}

//getHeight 返回节点的高
func getHeight(a *AVLTreeNode) int32 {
	if a == nil {
		return 0
	}
	return a.height
}

//max 获取2数中的较大值
func max(x, y int32) int32 {
	if x < y {
		return y
	}
	return x
}

//findMax 查找最大节点
func (a *AVLTreeNode) findMax() *AVLTreeNode {
	if a == nil {
		return nil
	}
	if a.right != nil {
		return a.right.findMax()
	}
	return a
}

//findMin 查找最小值
func (a *AVLTreeNode) findMin() *AVLTreeNode {
	if a == nil {
		return nil
	}
	if a.left != nil {
		return a.left.findMin()
	}
	return a
}

//getBF 获取节点左右子树高度差绝对值
//将二叉树上节点的左子树高度减去右子树高度取绝对值(Balance Factor)
func (a *AVLTreeNode) getBF(isAbs bool) int32 {
	var lh, rh int32
	if a.left != nil {
		lh = getHeight(a.left)
	}
	if a.right != nil {
		rh = getHeight(a.right)
	}
	bf := lh - rh
	if isAbs && bf < 0 {
		bf = bf * -1
	}
	return bf
}

//leftRotation 左旋
//a为最小失衡子数的根节点
//问题：在右子树插入右孩子导致AVL失衡
//return 新的平衡树的根节点
func (a *AVLTreeNode) leftRotation() *AVLTreeNode {
	tmpNode := a.right
	replaceFather(a, tmpNode)
	a.right = tmpNode.left
	if tmpNode.left != nil {
		tmpNode.left.parent = a // tmd, 一定要记得更新父亲
	}

	tmpNode.parent = a.parent
	tmpNode.left = a
	a.parent = tmpNode

	a.flushHS()
	tmpNode.flushHS()
	return tmpNode
}

//rightRotation 右旋
//a为最小失衡树的根节点
//问题：在左子树上插入左孩子导致AVL树失衡
//return 新的平衡树的根节点
func (a *AVLTreeNode) rightRotation() *AVLTreeNode {
	tmpNode := a.left
	// 更改前的节点父亲只向新节点
	replaceFather(a, tmpNode)
	tmpNode.parent = a.parent

	a.left = tmpNode.right
	if tmpNode.right != nil {
		tmpNode.right.parent = a
	}
	tmpNode.right = a
	a.parent = tmpNode

	a.flushHS()
	tmpNode.flushHS()
	return tmpNode
}

//rightLeftRotation  右左双旋转
//问题：通常因为在右子树上插入左孩子导致AVL失衡
//解发：先右旋后左旋调整
//return 新的平衡树根节点
func (a *AVLTreeNode) rightLeftRotation() *AVLTreeNode {
	a.right = a.right.rightRotation()
	return a.leftRotation()
}

//leftRightRotation  左右双选择
//问题：通常因为在左子树上插入右孩子导致AVL失衡
//解发：先左旋后右旋调整
//return 新的平衡树根节点
func (a *AVLTreeNode) leftRightRotation() *AVLTreeNode {
	a.left = a.left.leftRotation()
	return a.rightRotation()
}

func insertLeft(tree *AVLTreeNode, val *AVLTreeNode) *AVLTreeNode {
	if tree == nil {
		tree = val
	} else {
		tree.right = insertLeft(tree.right, val)
		tree.right.parent = tree
		tree = tree.rebalance()
		//if tree.getBF(true) == 2 { //在右子树插入新节点后avl树失衡
		//	//情况4：v插入到右子树的右孩子节点，只需要进行一次左旋转
		//	tree = tree.leftRotation()
		//}
	}
	tree.flushHS()
	return tree
}

//更新高度和大小
func (tree *AVLTreeNode) flushHS() {
	tree.size = getSize(tree.left) + getSize(tree.right) + 1
	tree.height = max(getHeight(tree.left), getHeight(tree.right)) + 1
}

// 得到 rank 的树节点
func (root *AVLTreeNode) find(rank int32) *AVLTreeNode {
	if rank < 0 {
		panic("")
	}
	l := getSize(root.left)
	if rank == l+1 {
		// 当前点就是
		return root
	}
	if rank <= l {
		return root.left.find(rank)
	} else {
		return root.right.find(rank - l - 1)
	}
}

func (root *AVLTreeNode) getRank() int32 {
	rank := getSize(root.left) + 1
	var dfs func(a *AVLTreeNode, sum int32, isLeft bool) int32
	dfs = func(a *AVLTreeNode, sum int32, isLeft bool) int32 {
		if a == nil {
			return sum // 根节点之上
		}
		if !isLeft {
			sum += getSize(a.left) + 1 // 查询的点在右边
		}
		if a.parent != nil {
			sum = dfs(a.parent, sum, a.parent.left == a)
		}
		return sum
	}

	if root.parent != nil {
		return dfs(root.parent, rank, root.parent.left == root) // 从下往上累加
	}
	return rank
}
func (root *AVLTreeNode) insertByRank(rank int32, val *AVLTreeNode) *AVLTreeNode {
	if rank < 0 {
		panic("")
	}
	l := getSize(root.left)
	if rank == l+1 {
		// 当前点就是
		root.left = insertLeft(root.left, val)
		root.left.parent = root
		root.flushHS()
		// 当前调整
		if root.getBF(true) == 2 {
			root = root.leftRightRotation()
		}
		return root
	}
	if rank <= l {
		root.left = root.left.insertByRank(rank, val)
		// 左
		if root.getBF(true) == 2 {
			if rank-1 <= getSize(root.left.left) { // 左
				root = root.rightRotation() //  左 左
			} else {
				root = root.leftRightRotation() // 左 右
			}
		}
	} else {
		rank = rank - l - 1
		root.right = root.right.insertByRank(rank, val)
		if root.getBF(true) == 2 {
			if rank-1 <= getSize(root.right.left) { // 左
				root = root.rightLeftRotation()
			} else {
				root = root.leftRotation()
			}
		}
	}
	root.flushHS()
	return root
}

// 4个指针
// 1. 左右儿子
// 2. 左右儿子的父指针
// 3. 替换节点的父指针
// 4. 原来节点的儿子指针
func (node *AVLTreeNode) delete() *AVLTreeNode {
	// 1. 利用map找到当前点
	up := node
	// 1. 有两个儿子
	if node.left != nil && node.right != nil {
		prev := node.left.findMax() // 前驱节点
		// 更新的起点
		up = prev.parent
		if up == node {
			up = prev // 特判就是 prev 是 node 的做儿子
		}
		// 前驱节点有儿子 (只能有 左儿子)
		if prev.right != nil {
			panic("找到的前驱节点不能有右儿子")
		}
		if prev.left != nil && node != prev.parent {
			prev.parent.right = prev.left
			prev.left.parent = prev.parent
		}

		// 1.将儿子节点都只向prev
		if node != prev.parent { //  前驱节点不是是 它 的左子树 需要将它的儿子指向父亲
			prev.left = node.left
			node.left.parent = prev
		}

		prev.right = node.right
		node.right.parent = prev
		replaceFather(prev, nil) // 如果这一步不设就会成还
		prev.parent = node.parent
		replaceFather(node, prev)

		// 向上递归调整
		// 删除
	} else if node.left != nil {
		// 代替
		up = node.left
		up.right = node.right
		up.parent = node.parent
		replaceFather(node, up)
	} else if node.right != nil {
		up = node.right
		up.left = node.left
		up.parent = node.parent
		replaceFather(node, up)
	} else {
		// 叶子节点 删除
		up = node.parent
		replaceFather(node, nil)
	}

	return upBalance(up)
}

func upBalance(node *AVLTreeNode) *AVLTreeNode {
	if node == nil {
		return nil
	}
	node.flushHS()
	node = node.rebalance()
	node.flushHS()
	if node.parent == nil {
		return node
	}
	return upBalance(node.parent)
}

func (node *AVLTreeNode) rebalance() *AVLTreeNode {
	// 如果右子树的高度比左子树的高度大于2
	if node.getBF(false) == -2 {
		// 如果 node.Right 的右子树的高度比node.Right的左子树高度大
		// 直接对node进行左旋转
		// 否则先对 node.Right进行右旋转然后再对node进行左旋转
		lz := getHeight(node.right.left)
		rz := getHeight(node.right.right)

		if rz > lz || (lz == rz && lz <= rz) {
			node = node.leftRotation()
		} else {
			node = node.rightLeftRotation()
		}
		// 如果左子树的高度比右子树的高度大2
	} else if node.getBF(false) == 2 {
		// 如果node.Left的左子树高度大于node.Left的右子树高度
		// 那么就直接对node进行右旋
		// 否则先对node.Left进行左旋，然后对node进行右旋
		lz := getHeight(node.left.left)
		rz := getHeight(node.left.right)
		if lz > rz || (lz == rz && lz >= rz) {
			node = node.rightRotation()
		} else {
			node = node.leftRightRotation()
		}
	}
	return node
}

func replaceFather(node *AVLTreeNode, val *AVLTreeNode) {
	if node.parent != nil {
		if node.parent.left == node {
			node.parent.left = val
		} else if node.parent.right == node {
			node.parent.right = val
		}
	}
}

// Displays nodes left-depth first (used for debugging)
func (n *AVLTreeNode) displayNodesInOrder() {
	if n.left != nil {
		n.left.displayNodesInOrder()
	}
	fmt.Print(n.value, " ")
	if n.right != nil {
		n.right.displayNodesInOrder()
	}
}

// 求得这个排名是区间的所有值
func (a *AVLTreeNode) seg(low, high int32, ans []*AVLTreeNode) []*AVLTreeNode {
	l := getSize(a.left)
	// high == l + 1 是向右情况下
	if low <= l && a.left != nil {
		ans = a.left.seg(low, high, ans)
	}
	l++ //还有当前点
	// 当前排名在范围内
	if low <= l && l <= high {
		ans = append(ans, a)
	}
	if a.right != nil && (low-l > 0 || high-l > 0) {
		ans = a.right.seg(low-l, high-l, ans)
	}
	return ans
}

// 打印树
func (t *AvlTree) print(node *AVLTreeNode, level int, tail bool) []byte {
	var b bytes.Buffer

	if node == nil {
	} else {
		b.Write(t.print(node.left, level+1, false))

		b.Write(bytes.Repeat([]byte("    "), level))
		if tail {
			b.Write([]byte("└── "))
		} else {
			b.Write([]byte("┌── "))
		}

		b.Write([]byte(fmt.Sprint(node.value)))
		b.Write([]byte(`
`))

		b.Write(t.print(node.right, level+1, true))
	}
	return b.Bytes()
}
