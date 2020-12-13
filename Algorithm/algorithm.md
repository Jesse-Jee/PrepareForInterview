# 剑指offer
## 38 字符串的排列
回溯 + 剪枝 dfs  
时间复杂度O(N^2)  
空间复杂度O（N）
输入一个字符串，打印出该字符串中字符的所有排列。

```go
 func permutation(s string) []string {
      length := len(s)
      if length == 0 {
          return nil
      }
      str := []byte(s)
      result := []string{}
      dfs(str, 0, length - 1, &result)
      return result
  }
 func dfs(str []byte, i, l int, result *[]string) {
      if i == l {
          *result = append(*result, string(str))
          return
      }
      visited := make([]bool, 26)
      for k := i; k <= l; k++ {
          if !visited[str[k] - 'a']{
               visited[str[k]-'a'] = true
               str[i], str[k] = str[k], str[i]
               dfs(str, i+1, l, result)
               str[i], str[k] = str[k], str[i]
          }
         
      }
 } 
```

## 46 把数字翻译成字符串
给定一个数字，我们按照如下规则把它翻译为字符串：0 翻译成 “a” ，1 翻译成 “b”，……，11 翻译成 “l”，……，25 翻译成 “z”。
一个数字可能有多个翻译。请编程实现一个函数，用来计算一个数字有多少种不同的翻译方法。

时间复杂度 O（logN）
空间复杂度 因为s所以O（logN）

方法一、动规
```go
    func translateNum(num int) int {
        str := strconv.Itoa(num)
        result, curpre, pre := 1, 0, 0
        for i := 0; i < len(str); i++ {
           curpre, pre, result = pre, result, 0
            result += pre
            if i == 0 {
                continue
            }
            s := str[i-1:i+1]
            if s <= "25" && s >= "10" {
                result += curpre
            }
        }
        return result
    }
```
方法二、整除法
时间复杂度O（logN）
空间复杂度O（1）
```go
 func translateNum(num int) int {
     if num < 10 {
         return 1
     }
 
     var result int
     if num % 100 <= 25 && num % 100 >= 10 {
         result += translateNum(num/100)
         result += translateNum(num/10)
     }else {
         result += translateNum(num/10)
     }
 
     return result
 }
```

## 29 顺时针打印矩阵
输入一个矩阵，按照从外向里以顺时针的顺序依次打印出每一个数字。  
时间复杂度 O（N）
空间复杂度 O（N）
```go
    func spiralOrder(matrix [][]int) []int {
        var upper, bottom, left, right int
        bottom = len(matrix)
        if bottom ==  0 {
            return nil
        }
        right = len(matrix[0])
        num := bottom * right
        result := make([]int, num)
        upper, bottom, left, right = 0, bottom - 1, 0, right - 1 
    
        for k := 0; k < num; {
            for i := left; i <= right; i++ {
                result[k] =  matrix[upper][i]
                k++
            }
            upper++
            if k > num ||upper > bottom {
                break
            }
    
            for i := upper; i <= bottom; i++ {
                result[k] =  matrix[i][right]
                k++
            }
            right--
            if  k > num || left > right {
                break
            }
    
    
            for i := right; i >= left; i-- {
                result[k] =  matrix[bottom][i]
                k++
            }
            bottom--
            if  k > num || upper > bottom {
                break
            }
    
            for i := bottom; i >= upper; i-- {
                result[k] = matrix[i][left]
                k++
            }
            left++
            if  k > num || left > right {
                break
            }
    
        }
    
        return result
    }
```

## 24 反转链表
定义一个函数，输入一个链表的头节点，反转该链表并输出反转后链表的头节点。  

方法一、暴力解法，空间大

```go
    /**
     * Definition for singly-linked list.
     * type ListNode struct {
     *     Val int
     *     Next *ListNode
     * }
     */
    func reverseList(head *ListNode) *ListNode {
        if head == nil {
            return nil
        }
    
        var result = new(ListNode)
        result.Next = nil
        result.Val = head.Val
        head = head.Next
        for head != nil {
             tmp := ListNode{
                Val : head.Val,
                Next: result,
            }
            result = &tmp
            head = head.Next
        }
        return result
    }
```
方法二、双指针
```go
    /**
     * Definition for singly-linked list.
     * type ListNode struct {
     *     Val int
     *     Next *ListNode
     * }
     */
    func reverseList(head *ListNode) *ListNode {
    	var cur *ListNode
        for head != nil {
            tmp := head.Next
            head.Next = cur
            cur = head
            head = tmp
        }
    
        return  cur
    }
```

## 25 合并两个有序链表
输入两个递增排序的链表，合并这两个链表并使新链表中的节点仍然是递增排序的。
```go
    /**
     * Definition for singly-linked list.
     * type ListNode struct {
     *     Val int
     *     Next *ListNode
     * }
     */
    func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
        result := new(ListNode)
        tmp := result
        for l1 != nil && l2 != nil {
            if l1.Val < l2.Val {
                tmp.Next = l1
                l1 = l1.Next
            }else{
                tmp.Next = l2
                l2 = l2.Next
            }
            tmp = tmp.Next
        }
    
        if l1 != nil {
            tmp.Next = l1
        }
        if l2 != nil {
            tmp.Next = l2
        }
        return result.Next
    }
```

## 3 重复数组中的重复数字

找出数组中重复的数字。
在一个长度为 n 的数组 nums 里的所有数字都在 0～n-1 的范围内。数组中某些数字是重复的，但不知道有几个数字重复了，也不知道每个数字重复了几次。
请找出数组中任意一个重复的数字。

```go
    func findRepeatNumber(nums []int) int {
        count := make([]int, cap(nums))
        for i := 0;i < len(nums); i++{
            count[nums[i]]++
            if count[nums[i]] > 1 {
                return nums[i]
            }
        }
        return -1
    }
```


## 9 用两个栈实现队列

```go
    type CQueue struct {
      in stack
      out stack
    }
    
    type stack []int
    
    func (s *stack) Push(value int) {
        *s = append(*s, value)
    }
    
    func (s *stack) Pop() int {
        l := len(*s)
        n := (*s)[l-1]
        *s = (*s)[:l-1]
        return n
    }
    
    
    
    func Constructor() CQueue {
       return CQueue{}
    }
    
    
    func (this *CQueue) AppendTail(value int)  {
        this.in.Push(value)
    }
    
    
    func (this *CQueue) DeleteHead() int {
       if len(this.out) != 0 {
           return this.out.Pop()
       }else if len(this.in) != 0 {
           for len(this.in) != 0 {
               this.out.Push(this.in.Pop())
           }
           return this.out.Pop()
       } 
       return -1
    }
    
    
    /**
     * Your CQueue object will be instantiated and called as such:
     * obj := Constructor();
     * obj.AppendTail(value);
     * param_2 := obj.DeleteHead();
     */
```

## 24 连续子数组的最大和
输入一个整型数组，数组中的一个或连续多个整数组成一个子数组。求所有子数组的和的最大值。
要求时间复杂度为O(n)  

```go
    func maxSubArray(nums []int) int {
         result := nums[0]
         for i := 1; i < len(nums); i++{
                nums[i] = max(nums[i-1]+nums[i], nums[i])
                result = max(result, nums[i])
         } 
         return result  
    }
    
    func max (a,b int) int {
        if a > b {
            return a
        }
        return b
    }
```

## 48 最长不含重复字符的子字符串
请从字符串中找出一个最长的不包含重复字符的子字符串，计算该最长子字符串的长度。  
使用动规或滑动窗口  

方法一、滑动窗口
时间复杂度O（N^2）
空间复杂度O（1）

```go
    func lengthOfLongestSubstring(s string) int {
        l := len(s)
        if l == 0 {
            return 0
        }
        left, maxLen := 0, 1
        m := make(map[byte]bool, 26)
        m[s[0]] = true
    
        for right := 1; right < l; {
            if _, ok := m[s[right]]; !ok {
                m[s[right]] = true
                right++  
            }else{
                delete(m, s[left])
                left++
            }
            tmpLen := right - left
            if tmpLen > maxLen {
                maxLen = tmpLen
            }
        }
        return maxLen
    }
```


## 68 二叉树最近公共祖先II
给定一个二叉树, 找到该树中两个指定节点的最近公共祖先

后序遍历
时间复杂度O（N）
空间复杂度O（N）
```go
    /**
     * Definition for TreeNode.
     * type TreeNode struct {
     *     Val int
     *     Left *ListNode
     *     Right *ListNode
     * }
     */
     func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
         if root == nil {
             return nil
         }
    
         if root.Val == p.Val || root.Val == q.Val {
             return root
         }
    
         left := lowestCommonAncestor(root.Left, p, q)
         right := lowestCommonAncestor(root.Right, p, q)
    
         if left != nil && right != nil {
             return root
         }
         if left == nil {
             return right
         }
         return left
    }
```

## 68 二叉搜索树的最近公共祖先I

方法一、迭代法：
```go
    func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
        pVal, qVal := p.Val, q.Val
        node := root 
        for node != nil {
            parentVal := node.Val
            switch{
                case pVal > parentVal && qVal > parentVal:
                    node = node.Right
                case pVal < parentVal && qVal < parentVal:
                    node = node.Left
                default:
                    return node
            }
        }
        return nil
    }
```

方法二、递归法：
```go
    func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
        rootVal, pVal, qVal := root.Val, p.Val, q.Val
        switch {
            case pVal > rootVal && qVal > rootVal:
                return lowestCommonAncestor(root.Right,p,q)
            case pVal < rootVal && qVal < rootVal:
                return lowestCommonAncestor(root.Left,p,q)
            default:
                return root
        }
    }
```

## 检查链表是否有环
```go
    /**
     * Definition for singly-linked list.
     * type ListNode struct {
     *     Val int
     *     Next *ListNode
     * }
     */
    func hasCycle(head *ListNode) bool {
        fast, slow := head, head
        for fast != nil && fast.Next != nil {
            fast = fast.Next.Next
            slow = slow.Next
            if fast == slow {
                return true
            }
        }
        return false
    }
```

## 22. 链表中倒数第k个节点
输入一个链表，输出该链表中倒数第k个节点。为了符合大多数人的习惯，本题从1开始计数，即链表的尾节点是倒数第1个节点。例如，一个链表有6个节点，
从头节点开始，它们的值依次是1、2、3、4、5、6。这个链表的倒数第3个节点是值为4的节点。

双指针：
```go
    func getKthFromEnd(head *ListNode, k int) *ListNode {
        fast, slow := head, head
        for i := 0; i < k ; i++{
            if fast == nil {
                return nil
            }
            fast = fast.Next
        }
        
        for fast != nil {
            fast = fast.Next
            slow = slow.Next
        }
        
        return slow
    }
```

# 二叉树层序遍历
bfs递归

```go
    var result [][]int
    func levelOrder(root *TreeNode) [][]int {
        if root == nil {
            return nil
        }
        result = [][]int{}
        bfs(root, 0)
        return result
    }
    
    func bfs(root *TreeNode, level int) {
        if root == nil {
            return
        }
    
        if level == len(result) {
            result = append(result, []int{})
        } 
    
        result[level] = append(result[level], root.Val)
        bfs(root.Left, level+1)
        bfs(root.Right, level+1)
    }
```
非递归

```go
    func levelOrder(root *TreeNode) [][]int {
        result := [][]int{}
        if root == nil {
           return result
        }
    
        q := []*TreeNode{root}
        for i := 0; len(q) > 0; i++{
            result = append(result, []int{})
            p := []*TreeNode{}
            for j := 0; j < len(q); j++{
                node := q[j]
                result[i] = append(result[i], node.Val)
                if node.Left != nil {
                    p = append(p, node.Left)
                }
                if node.Right != nil {
                    p = append(p, node.Right)
                }
            }
            q = p
        }   
        return result
    }
```

# 07 重建二叉树
输入某二叉树的前序遍历和中序遍历的结果，请重建该二叉树。假设输入的前序遍历和中序遍历的结果中都不含重复的数字。  
例如，给出  
前序遍历 preorder = [3,9,20,15,7]  
中序遍历 inorder = [9,3,15,20,7]  

```go
    func buildTree(preorder []int, inorder []int) *TreeNode {
        for k := range inorder {
            if preorder[0] == inorder[k] {
                return &TreeNode{
                    Val: inorder[k],
                    Left: buildTree(preorder[1:k+1], inorder[0:k]),
                    Right:buildTree(preorder[k+1:],inorder[k+1:]),
                }
            }
        }
        return nil
    }
```


# 买卖股票全家桶
## 121 买卖股票的最佳时机
给定一个数组，它的第 i 个元素是一支给定股票第 i 天的价格。
如果你最多只允许完成一笔交易（即买入和卖出一支股票一次），设计一个算法来计算你所能获取的最大利润。
注意：你不能在买入股票前卖出股票。

定义两个变量，一个存当前最小值，一个记录当前最大利润。
如果当前值比最小值小，就赋值给最小值变量。如果比最小值大，就检查是否比当前最大利润大。如果比最大利润大，就赋值给最大利润变量。返回最大利润。

时间复杂度 O(N)
```go
    func maxProfit(prices []int) int {
        if len(prices) < 2 {
            return 0
        }
        maxP, min := 0, prices[0]
        for i := range prices {
            if prices[i] < min {
                min = prices[i]
            }else if prices[i] - min > maxP {
                maxP = prices[i] - min
            }
        } 
        return maxP
    }
```

## 122 买卖股票的最佳时机2
允许多次买卖，但再买之前必须先卖掉手上的。

给定一个数组，它的第 i 个元素是一支给定股票第 i 天的价格。        
设计一个算法来计算你所能获取的最大利润。你可以尽可能地完成更多的交易（多次买卖一支股票）。       
注意：你不能同时参与多笔交易（你必须在再次购买前出售掉之前的股票）。          

方法一：
一次遍历法。
时间复杂度O(n)
```go
    func maxProfit(prices []int) int {
        // 定义利润和，只要当天比前一天有利润，就计入利润。
        if len(prices) < 2 {
            return 0
        }
        maxP := 0
        for i := 1; i < len(prices); i++ {
            if prices[i] > prices[i-1] {
                maxP += prices[i]-prices[i-1]
            }
        }
        return maxP
    }
```

方法二：
动态规划

```go
    func maxProfit(prices []int) int {
        // 定义状态
        // dp[i][0] = 第i天手上没有股票最大利润
        //  dp[i][0] 有两种情况：
        //    前一天手里也没有：
        //    前一天手里有，今天卖  
        //    dp[i][0] = max(dp[i-1][0], dp[i-1][1] + prices[i])    
        // dp[i][1] = 第i天手上有股票最大利润
        //    前一天手上没有股票，今天买：dp[i-1][0] - prices[i]
        //    前一天手头上有股票，今天没买： dp[i-1][1]
        //    dp[i][1] = max(dp[i-1][1], dp[i-1][0] - prices[i])
        // 最终就是求 dp[i][0] 和dp[i][1]的大小
        // 肯定0 大
        // 为了减少空间，使用两个变量标识即可
         l := len(prices)
        if l < 2 {
            return 0
        }
        hold, unhold := -prices[0], 0
        for i := 1; i < l; i++{
            hold = max(unhold-prices[i], hold)
            unhold = max(hold+prices[i], unhold)
        }
        return unhold
    }
    
    func max(a, b int) int {
        if a > b {
            return a
        }
        return b
    }
```

##  309 买卖股票的最佳时机含冷冻期
给定一个整数数组，其中第 i 个元素代表了第 i 天的股票价格 。​
设计一个算法计算出最大利润。在满足以下约束条件下，你可以尽可能地完成更多的交易（多次买卖一支股票）:
你不能同时参与多笔交易（你必须在再次购买前出售掉之前的股票）。
卖出股票后，你无法在第二天买入股票 (即冷冻期为 1 天)。

动态规划方程
```go
    func maxProfit(prices []int) int {
        // dp[i][s]
        // 状态分为：持有股票非冷冻；无股票冷冻；无股票非冷冻
        //dp[i][0]; dp[i][1];dp[i][2]
        l := len(prices)
        if l < 2 {
            return 0
        }
        holdno, unholdyes, unholdno := -prices[0], 0, 0
        for i := 1; i < l; i++ {
            new1 := max(unholdno-prices[i], holdno)
            new2 := holdno+prices[i]
            new3 := max(unholdno, unholdyes)
            holdno, unholdyes, unholdno = new1, new2, new3
        }
        return max(unholdyes,unholdno)
    }
    
    
    func max(a,b int) int {
        if a > b {
            return a
        }
        return b
    }
```

## 714 买卖股票最佳时机含手续费
给定一个整数数组 prices，其中第 i 个元素代表了第 i 天的股票价格 ；非负整数 fee 代表了交易股票的手续费用。
你可以无限次地完成交易，但是你每笔交易都需要付手续费。如果你已经购买了一个股票，在卖出它之前你就不能再继续购买股票了。
返回获得利润的最大值。
注意：这里的一笔交易指买入持有并卖出股票的整个过程，每笔交易你只需要为支付一次手续费。

动态规划

```go
    func maxProfit(prices []int, fee int) int {
        // dp[i][0] 最大利润没有股票
        //  max(dp[i-1][1]+prices[i],dp[i-1][0])
        // dp[i][1] 最大利润有股票
        //  max(dp[i-1][0]-prices[i],dp[i-1][1])
        l := len(prices)
        if l < 2 {
            return 0
        }
        hold, unhold := -prices[0], 0
        for i := 1; i < l; i++ {
            hold = max(unhold - prices[i], hold)
            unhold = max(hold+prices[i]-fee, unhold)
        }
        return unhold
    }
    
    func max(a,b int) int{
        if a > b {
            return a
        }
        return b
    }
```

# 226 翻转二叉树
方法一：递归       
时间复杂度O(N)
空间复杂度O(N)   
```go
    /**
     * Definition for a binary tree node.
     * type TreeNode struct {
     *     Val int
     *     Left *TreeNode
     *     Right *TreeNode
     * }
     */
    func invertTree(root *TreeNode) *TreeNode {
        if root == nil {
            return root
        } 
        left := invertTree(root.Left)
        right := invertTree(root.Right)
        root.Right = left
        root.Left = right
        return root
    }
```

方法二：迭代          
时间复杂度O(N)
空间复杂度O(N) 
```go
    /**
     * Definition for a binary tree node.
     * type TreeNode struct {
     *     Val int
     *     Left *TreeNode
     *     Right *TreeNode
     * }
     */
    func invertTree(root *TreeNode) *TreeNode {
        if root == nil {
            return root
        } 
        q :=[]*TreeNode{root}
    
        for len(q) > 0 {
            node := q[0]
            q = q[1:]
            node.Left, node.Right = node.Right, node.Left
            if node.Left != nil {
                q = append(q,node.Left)
            }
            if node.Right != nil {
                q = append(q, node.Right)
            }
        }
        return root
    }
```



# 爬楼梯全家桶
## 70 爬楼梯
假设你正在爬楼梯。需要 n 阶你才能到达楼顶。                 
每次你可以爬 1 或 2 个台阶。你有多少种不同的方法可以爬到楼顶呢？                 
注意：给定 n 是一个正整数。                 

时间复杂度O(N)
```go
func climbStairs(n int) int {
    p, q, r := 0,0,1

    for i := 1; i <= n; i++ {
           p = q
           q = r
           r = p + q
    }

    return r
    
}
```

## 746 使用最小花费爬楼梯
数组的每个索引作为一个阶梯，第 i个阶梯对应着一个非负数的体力花费值 cost[i](索引从0开始)。                 
每当你爬上一个阶梯你都要花费对应的体力花费值，然后你可以选择继续爬一个阶梯或者爬两个阶梯。                   
您需要找到达到楼层顶部的最低花费。在开始时，你可以选择从索引为 0 或 1 的元素作为初始阶梯。                

时间复杂度O(N)
```go
func minCostClimbingStairs(cost []int) int {  
    for i := len(cost)-3; i >= 0; i-- {
        cost[i] += min(cost[i+1], cost[i+2])
    }
    return min(cost[0],cost[1])
}
func min(a,b int) int {
    if a < b {
        return a
    }
    return b
}
```

# 2 两数相加
给出两个 非空 的链表用来表示两个非负的整数。其中，它们各自的位数是按照 逆序 的方式存储的，并且它们的每个节点只能存储一位数字。               
如果，我们将这两个数相加起来，则会返回一个新的链表来表示它们的和。                   
您可以假设除了数字 0 之外，这两个数都不会以 0 开头。                   

时间复杂度O(N)               
```go
    /**
     * Definition for singly-linked list.
     * type ListNode struct {
     *     Val int
     *     Next *ListNode
     * }
     */
    func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
        var result *ListNode
        var tmp *ListNode
        var c int
    
        for l1 != nil || l2 != nil {
            n1, n2 := 0, 0
    
            if l1 != nil {
                n1 = l1.Val
                l1 = l1.Next
            }
    
            if l2 != nil {
                n2 = l2.Val
                l2 = l2.Next
            }
    
            sum := n1 + n2+ c
            sum, c = sum%10, sum/10
    
            if result == nil {
                result = &ListNode{Val:sum}
                tmp = result
            }else {
                tmp.Next = &ListNode{Val:sum}
                tmp = tmp.Next
            }
        }
        if c > 0 {
            tmp.Next = &ListNode{Val:c}
        }
    
        return result
    }
```

# 剑指offer 55-I 二叉树的高度
输入一棵二叉树的根节点，求该树的深度。从根节点到叶节点依次经过的节点（含根、叶节点）形成树的一条路径，最长路径的长度为树的深度。

迭代方式
时间复杂度O(N)
```go
    func maxDepth(root *TreeNode) int {
        if root == nil{
            return 0
        }
       var max int
       q := []*TreeNode{root}
    
       for len(q) > 0 {
           p := []*TreeNode{}
           max++
            for i := 0; i < len(q);i++{
                if q[i].Left != nil {
                    p = append(p, q[i].Left)
                }
                if q[i].Right != nil {
                    p = append(p, q[i].Right)
                }
            }
            q = p
       } 
    
        return max
    }
```

递归方式：
```go
var max int
func maxDepth(root *TreeNode) int {
    max = 0
    dfs(root, 0)
    return max

}
func dfs(root *TreeNode, level int) {
    if root == nil {
         if level > max {
                max=level
        }
        return
    }
   

    dfs(root.Left, level+1)
    dfs(root.Right,level+1)
}
```

# 100相同的树
给定两个二叉树，编写一个函数来检验它们是否相同。                
如果两个树在结构上相同，并且节点具有相同的值，则认为它们是相同的。               


时间复杂度O(min(m,n))

```go
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func isSameTree(p *TreeNode, q *TreeNode) bool {
    if p == nil && q == nil {
        return true
    }
    if p == nil || q == nil {
        return  false
    }
    if p.Val != q.Val {
        return false
    }
    return isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right)
}
```

# 实现LRU
```go
    package main
    
    /**
     * lru design
     * @param operators int整型二维数组 the ops
     * @param k int整型 the k
     * @return int整型一维数组
    */
    
    type LinkNode struct {
        key, value int
        prev, next *LinkNode
    }
    
    type LRUCache struct {
        Cap int
        data map[int]*LinkNode
        head, tail *LinkNode
        
    }
    
    func Constructor(capcity int) LRUCache {
        head := &LinkNode{-1,-1,nil,nil}
        tail := &LinkNode{-1,-1,nil,nil}
        head.next = tail
        tail.prev = head
        return LRUCache{Cap: capcity, data: make(map[int]*LinkNode),head: head,tail: tail}
    }
    
    func (l *LRUCache) AddNode(node *LinkNode){
        node.prev = l.head
        node.next = l.head.next
        l.head.next = node
        node.next.prev = node
    }
    
    func (l *LRUCache) RemoveNode(node *LinkNode) {
        node.prev.next = node.next
        node.next.prev = node.prev
    }
    
    func (l *LRUCache) MoveToHead(node *LinkNode) {
        l.RemoveNode(node)
        l.AddNode(node)
    }
    
    
    func (l *LRUCache) Set(key, value int) {
        d := l.data
        if node, ok := d[key]; ok{
            node.value = value
            l.MoveToHead(node)
        } else {
            newNode := &LinkNode{key: key,value: value, prev: nil, next: nil}
            if len(d) >= l.Cap {
                delete(d, l.tail.prev.key)
                l.RemoveNode(l.tail.prev)
            }
            d[key] = newNode
            l.AddNode(newNode)
        }
    }
    
    func (l *LRUCache) Get(key int) int {
        d := l.data
        if _,ok := d[key]; !ok {
            return -1
        }
        node := d[key]
        l.MoveToHead(node)
        return node.value
    }
    
    
    
    func LRU( operators [][]int ,  k int ) []int {
        // write code here
        lru := Constructor(k)
        result := make([]int, 0)
        for i := 0; i < len(operators); i++ {
            operator := operators[i]
            if operator[0] ==  1 {
               lru.Set(operator[1],operator[2])
            }else if operator[0] == 2 {
                result = append(result, lru.Get(operator[1]))
            }
        }
        
        return result
        
    }
```

# 剑指offer04二维数组中的查找
在一个 n * m 的二维数组中，每一行都按照从左到右递增的顺序排序，每一列都按照从上到下递增的顺序排序。                   
请完成一个高效的函数，输入这样的一个二维数组和一个整数，判断数组中是否含有该整数。                   

时间复杂度O(m+n)
空间复杂度O(1)

````go
func findNumberIn2DArray(matrix [][]int, target int) bool {
    m := len(matrix)
    if m < 1 {
        return false
    }
    n := len(matrix[0])
    row, col := 0, n-1

    for row < m && col >= 0 {
        num := matrix[row][col]
        if num == target {
            return true
        }else if num > target {
            col--
        }else {
            row++
        }
    }
    return false
}
````

# 剑指offer06 从尾到头打印链表
输入一个链表的头节点，从尾到头反过来返回每个节点的值（用数组返回）。                  
```go
/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func reversePrint(head *ListNode) []int {
    if head == nil {
        return nil
    }
    r := make([]int, 0)
    for head != nil {
        r = append(r, head.Val)
        head = head.Next
    }
    start, end := 0, len(r)-1

    for start < end {
        r[start], r[end] = r[end], r[start]
        start++
        end--
    }

    return r
}
```

# 剑指offer05替换空格
请实现一个函数，把字符串 s 中的每个空格替换成"%20"。                  
```go
func replaceSpace(s string) string {
    if s == "" {
        return s
    }
    result := make([]byte, 0)
    for i := range s {
        if s[i] == ' ' {
            result = append(result, []byte{'%','2','0'}...)    
        }else {
            result = append(result, s[i])
        }
    }
    return string(result)
}
```

# 剑指offer28对称的二叉树
请实现一个函数，用来判断一棵二叉树是不是对称的。如果一棵二叉树和它的镜像一样，那么它是对称的。                 
````go
/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func isSymmetric(root *TreeNode) bool {
     if root == nil {
         return true
     }
     return bfs(root.Left, root.Right)
}

func bfs(a, b *TreeNode) bool {
    if a == nil && b == nil {
        return true
    }

    if (a == nil || b == nil) || (a.Val != b.Val) {
        return false
    }

    return bfs(a.Left, b.Right) && bfs(a.Right, b.Left)
}   
````


## 判断是否是平衡二叉树




## 旋转二叉树
## 二叉树前序遍历



# 数组
## 旋转数组最小数字
## 二维数组中的查找








#偏转顺序数组

#区间数组求交集


#十进制ip转换为32位整数


#位图


#方程求根


#leetcode 105  128  125

#原地链表排序


#编辑距离


# 堆排序



#一个长度为N的数组，里面的元素值在1-N之间（闭区间），找出重复元素。要求时间复杂度O（N），空间复杂度O（1）


#零钱兑换


#回文链表

#接雨水全家桶


#一堆木头，锯K段，找每段最大长度。


#合并K个有序数组

#64匹马，8个赛道，最少比多少次能找到最快的4匹


#排序字符串去重


#岛屿

#最短路径


#生成括号


#往有序循环列表中插入新节点


#给两个字符串形式的数 求和


#二叉树中是否存在一条路径和等于给定值

# 判定两棵树是相同的



#10亿数据内存够用的情况下，选取前100
#40亿数据内存不够的情况下找出中位数

