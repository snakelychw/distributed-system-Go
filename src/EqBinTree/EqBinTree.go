package main

import (
  "golang.org/x/tour/tree"
  "fmt"
  )

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, c chan int){ 
  if(t != nil){
    Walk(t.Left, c)
    c <- t.Value
    Walk(t.Right, c)
  } else {
    return
  }
  
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool{
  res := true
  c1 := make(chan int)
  c2 := make(chan int)
  go Walk(t1, c1)
  go Walk(t2, c2)
  for i:= 0; i<10; i++{
    if(<-c1 != <-c2){
      res = false
    }
  }
  return res
}

func main() {
  root := tree.New(1)
  rootie := tree.New(2)
  fmt.Print(Same(root, rootie))
}

