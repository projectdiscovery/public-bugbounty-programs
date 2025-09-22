package main

import (
  "encoding/json"
  "fmt"
  "net"
  "os"
)

type Program struct {
  Name    string   `json:"name"`
  URL     string   `json:"url"`
  Bounty  bool     `json:"bounty"`
  Domains []string `json:"domains"`
  Network []string `json:"network"`
}
type Root struct {
  Programs []Program `json:"programs"`
}

func main() {
  f, err := os.Open("chaos-bugbounty-list.json")
  if err != nil { fmt.Println(err); os.Exit(1) }
  defer f.Close()

  var r Root
  if err := json.NewDecoder(f).Decode(&r); err != nil {
    fmt.Printf("json decode error: %v\n", err); os.Exit(1)
  }

  var bad []string
  for _, p := range r.Programs {
    for _, n := range p.Network {
      if n == "" { continue }
      if ip := net.ParseIP(n); ip != nil { continue }
      if _, _, err := net.ParseCIDR(n); err == nil { continue }
      bad = append(bad, fmt.Sprintf("%s -> %s", p.Name, n))
    }
  }
  if len(bad) > 0 {
    fmt.Println("Invalid network entries:")
    for _, b := range bad { fmt.Println(" -", b) }
    os.Exit(2)
  }
}
