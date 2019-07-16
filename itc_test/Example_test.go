package itc_test

import (
    "fmt"
    "github.com/ziglet.io/go-itc/itc"
    "testing"
)

func Print(message string, stamps ...*itc.Stamp) {
    fmt.Print(fmt.Sprintf("%20s\t",message))
    for _,s := range stamps {
        fmt.Print(fmt.Sprintf("%10s\t%10s\t",s.Id.Print(),s.Event.Print()))
    }
    fmt.Println()
}

func TestTwoNodeForkJoin(t *testing.T){
    // Seed
    s := itc.SeedStamp()
    Print("seed",s)

    // Two events on original node
    s = s.Advance()
    Print("advance",s)
    s = s.Advance()
    Print("advance",s)

    // Fork and advance both
    l,r := s.Fork()
    Print("fork",l,r)
    l = l.Advance()
    r = r.Advance()
    Print("both",l,r)

    // Join
    s = l.Join(r)
    Print("join",s)
}

// Test asymmetric advance
func TestExampleAsymmetricAdvance(t *testing.T){
    s := itc.SeedStamp()
    Print("seed",s)

    l,r := s.Fork()
    Print("fork",l,r)

    r = r.Advance()
    Print("r advance",l,r)
    r= r.Advance()
    Print("r advance",l,r)

    s = l.Join(r)
    Print("join",s)
}

// What happens when a single node starts writing, a second node forks the first and writes,
// and they later join?
func TestExampleUnknownSplitJoin(t *testing.T){
    a := itc.SeedStamp()

    // Node A writes
    a = a.Advance()
    Print("a advance",a)

    // Node B forks unknown to A
    b,c := a.Fork()
    Print("fork",a,b,c)

    // A advances a couple
    a = a.Advance()
    Print("a advance",a,b,c)
    a = a.Advance()
    Print("a advance",a,b,c)

    // B advances a couple
    b = b.Advance()
    Print("b advance",a,b,c)
    b = b.Advance()
    Print("b advance",a,b,c)

    // The later join
    // This will fail because one cannot join a region with overlapping domains (e.g. a parent with 1/2 of a child)
    //d := a.Join(b)
    //Print("join a,b",a,b,c,d)
}