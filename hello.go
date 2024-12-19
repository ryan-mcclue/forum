package main

import (
  "fmt"
  "os"
  "bufio"
  "strings"
  "regexp"
  "strconv"
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

  var children []Child

  r := regexp.MustCompile(" ([0-9]+) (.*) bags?")

  if a[1] != " no other bags" {
    for _, child := range strings.Split(a[1], ",") {
      b := r.FindStringSubmatch(child)
      amt, err := strconv.Atoi(b[1])
      if (err != nil) {
        panic(err)
      }
      children = append(children, Child{
        colour: b[2],
        amount: amt, 
      })
    }
  } else {
    fmt.Println("no child bags");
  }

  return Rule{
    container: container,
    children: children,
  }
}

func solve_puzzle(file_path string) {
  file, err := os.Open(file_path)
  if (err != nil) {
    panic(err)
  }
  defer file.Close()

  rules := make(map[Colour][]Child)

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()

    rule := line_to_rule(line)
    rules[rule.container] = rule.children
  }
  if err := scanner.Err(); err != nil {
    panic(err)
  }

  fmt.Printf("%s\n", file_path)
}

func main() {
  for _, file_path := range os.Args[1:] {
    solve_puzzle(file_path)
  }

  fmt.Printf("%#v\n", line_to_rule("light red bags contain 1 bright white bag, 2 muted yellow bags."));
}



