package itc

import (
    "fmt"
    "strings"
)

// Section 5.2 Reduce the size of the event tree by norming
func (event *Event) Norm() *Event {
    // Case 1: Norm(n) -> n
    if event.IsLeaf {
        return event.Copy()
    }

    // Case 2: Norm((n,m,m)) -> (n+m) where m is an integer
    if !event.IsLeaf && event.Left.IsLeaf && event.Right.IsLeaf && event.Left.Value == event.Right.Value {
        return &Event{
            IsLeaf:true,
            Value:event.Value + event.Left.Value,
            Left:nil,
            Right:nil,
        }
    }

    // Case 3: Norm((n,e1,e2)) -> (n+m,e1.Sink(m),e2.Sink(m)) where m = Min(Min(e1),Min(e2))
    if !event.IsLeaf {
        e := Event{}
        m := Min(event.Left,event.Right)

        e.Value += event.Value + m.Value
        el := event.Left.Sink(m.Value)
        er := event.Right.Sink(m.Value)
        e.IsLeaf = false
        e.Right=er
        e.Left=el

        return &e
    }

    return nil
}

// Section 5.2 Lift the entire event tree by a constant value - used in norming
func (event *Event) Lift(m uint32) *Event {
    // Case 1: (n) -> (n+m)
    // Case 2: (n,e1,e2) -> (n+m, e1,e2)
    e := event.Copy()
    e.Value += m

    return e
}

// Section 5.2 Sink the entire event tree by a constant value - used in norming
func (event *Event) Sink(m uint32) *Event {
    // Case 1: (n) -> (n-m)
    // Case 2: (n,e1,e2) -> (n-m,e1,e2)
    e := event.Copy()
    e.Value -= m

    return e
}

// Section 5.2 Get an event that represents the minimum event value in a tree of values
func (event *Event) Min() *Event {

    if !event.IsLeaf {
        e := Min(event.Left.Min(),event.Right.Min())
        e.Value += event.Value
        return e
    } else {
        return event.Copy()
    }

}

// Section 5.2 Get an event that represents the maximum event value in a tree of values
func (event *Event) Max() *Event {

    if !event.IsLeaf {
        e := Max(event.Left.Max(),event.Right.Max())
        e.Value += event.Value
        return e
    } else {
        return event.Copy()
    }

}

// Section 5.3.1 Less or equals - relation that defines the partial order
func (event1 *Event) Leq(event2 *Event) bool {
    // Case 1: Leq(n1,n2) -> n1 <= n2
    if event1.IsLeaf && event2.IsLeaf {
        return event1.Value <= event2.Value
    }

    // Case 2: Leq(n1,(n2,l2,r2)) -> n1 <= n2
    if event1.IsLeaf && !event2.IsLeaf {
        return event1.Value <= event2.Value
    }

    // Case 3: Leq((n1,l1,r2),n2) -> n1<=n2 AND Leq(l1.Lift(n1),n2) AND Leq(r1.Lift(n1),n2)
    if !event1.IsLeaf && event2.IsLeaf {
        return event1.Value <= event2.Value && event1.Lift(event1.Value).Leq(event2) && event1.Right.Lift(event1.Value).Leq(event2)
    }

    return false
}

// Section 5.3.3 Join the event trees used in joining the Stamps
func (event1 *Event) Join(event2 *Event) *Event {
    e := Event{}

    // Case 1: join(n1,n2) -> max(n1,n2)
    if event1.IsLeaf && event2.IsLeaf {
        e.IsLeaf = true
        e.Value = Max(event1.Max(),event2.Max()).Value
        return &e
    }

    // Case 2: join(n1,(n2,l2,r2)) -> join((n1,0,0),(n2,l2,r2))
    if event1.IsLeaf && !event2.IsLeaf {
        a := Event{
            IsLeaf:true,
            Value:0,
        }
        b:= Event{
            IsLeaf:true,
            Value:0,
        }
        top := &Event{
            IsLeaf: false,
            Value: event1.Value,
            Left: &a,
            Right: &b,
        }

        return top.Join(event2)
    }

    // Case 3: join((n1,l1,r1),n2) -> join((n1,l1,r1),(n2,l2,r2))
    if !event1.IsLeaf && event2.IsLeaf {
        a := Event{
            IsLeaf:true,
            Value:0,
        }
        b:= Event{
            IsLeaf:true,
            Value:0,
        }
        event2.IsLeaf = false
        event2.Left = &a
        event2.Right = &b

        return event1.Join(event2)
    }

    // Case 4: join((n1,l1,r1),(n2,l2,r2)) -> join((n2,l2,r2),(n1,l1,r1)) if n1 > n2
    if !event1.IsLeaf && !event2.IsLeaf && event1.Value > event2.Value {
        return event2.Join(event1)
    }

    // Case 5: join((n1,l1,r1),(n2,l2,r2)) -> Norm((n1,join(l1,l2.Lift(n2-n1),join(r1,r2.Lift(n2-n1)))
    if !event1.IsLeaf && !event2.IsLeaf {
        left := event1.Left.Join(event2.Left.Lift(event2.Value - event1.Value))
        right := event1.Right.Join(event2.Right.Lift(event2.Value - event1.Value))

        event := Event{
            IsLeaf: false,
            Value:event1.Value,
            Left:left,
            Right:right,
        }
        return event.Norm()
    }

    return nil
}

// Shallow copy an event
func (event *Event) Copy() *Event {
    e := Event{
        IsLeaf: event.IsLeaf,
        Value: event.Value,
        Right: event.Right,
        Left: event.Left,
    }

    return &e
}

func Max(a *Event, b *Event) *Event {
    am := a.Max()
    bm := b.Max()
    if am.Value > bm.Value {
        return am
    } else {
        return bm
    }
}

func Min(a *Event, b *Event) *Event {
    am := a.Min()
    bm := b.Min()
    if am.Value < bm.Value {
        return am
    } else {
        return bm
    }
}

func (event *Event) Print() string {
    var sb strings.Builder

    if event.IsLeaf {
        sb.WriteString(fmt.Sprint(event.Value))
    } else {
        sb.WriteString(fmt.Sprint(event.Value))
        sb.WriteString(",")
        sb.WriteString("(")
        sb.WriteString(event.Left.Print())
        sb.WriteString(",")
        sb.WriteString(event.Right.Print())
        sb.WriteString(")")
    }

    return sb.String()
}

func NewEvent(value uint32) *Event{
    return &Event{
        IsLeaf: true,
        Value: value,
        Left: nil,
        Right: nil,
    }
}