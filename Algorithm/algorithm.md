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




#二叉树

##判断是否是平衡二叉树

##二叉树的最近公共祖先
##旋转二叉树
##二叉树前序遍历



#数组
##旋转数组最小数字
##二维数组中的查找








#偏转顺序数组

#区间数组求交集


#两个栈实现一个队列


#十进制ip转换为32位整数


#位图


#方程求根


#剑指offer，leetcode hot100

#leetcode 105  128  125

#原地链表排序


#编辑距离







#链表两数相加

#股票买卖全家桶

#爬楼梯全家桶

#令牌桶算法

#堆排序


#快排

#归并排序


#选择排序

#实现一个LRU


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

#如何判定两棵树是相同的
#二叉树的高度
#合并两个有序链表


#10亿数据内存够用的情况下，选取前100
#40亿数据内存不够的情况下找出中位数

