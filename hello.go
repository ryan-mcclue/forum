package main

import (
  "fmt"
  "os"
  "strings"
  "regexp"
)

// alias, not new incompatible type
type Colour = string 

type Child struct {
  colour Colour
  amount int
}

type Rule struct {
  container Colour
  children []Child
} 

func line_to_rule(line string) Rule {
  a := strings.Split(line, "contain")
  container := strings.TrimSpace(strings.Split(a[0], "bags")[0])
  fmt.Printf("%#v\n", container)

  fmt.Println("children bags:")
  if a[1] != " no other bags" {
    for _, child := range strings.Split(a[1], ",") {
      r := regexp.MustCompile(" ([0-9]+) (.*) bags?")
      fmt.Printf("%#v\n", r.FindStringSubmatch(child))
    }

  }

  return Rule{}
  return Rule{
    container: "ryan",
    children: []Child{
      Child{
        colour: "blue",
        amount: 10,
      },
    },
  }
}

func solve_puzzle(file_name string) {
  fmt.Printf("%s\n", file_name)
}

func main() {
  for _, file_name := range os.Args[1:] {
    solve_puzzle(file_name)
  }

  fmt.Printf("%#v\n", line_to_rule("light red bags contain 1 bright white bag, 2 muted yellow bags."));
}



