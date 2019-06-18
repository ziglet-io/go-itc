package itc

// Section 5.2 Reduce the size of the event tree by norming
func (event *Event) Norm() *Event {
    var e *Event

    // Case 1: Norm(n) -> n
    if event.IsLeaf {
        e = event.Copy()
    }

    // Case 2: Norm((n,m,m)) -> (n+m) where m is an integer
    if !event.IsLeaf && event.Left.IsLeaf && event.Right.IsLeaf && event.Left.Value == event.Right.Value {
        e = &Event{
            IsLeaf:true,
            Value:event.Value + event.Left.Value,
            Left:nil,
            Right:nil,
        }
    }

    // Case 3: Norm((n,e1,e2)) -> (n+m,e1.Sink(m),e2.Sink(m)) where m = Min(Min(e1),Min(e2))
    if !event.IsLeaf {
        m := Min(event.Left.Min().Value,event.Right.Min().Value)

        e := event.Copy()
        el := event.Left.Copy()
        er := event.Right.Copy()

        e.Value += m
        el = event.Left.Sink(m)
        er = event.Right.Sink(m)
        e.IsLeaf = false
        e.Right=er
        e.Left=el
    }

    return e
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
    e := event.Copy()

    if !e.IsLeaf {
        ml := event.Left.Min()
        mr := event.Right.Min()

        e.Value += Min(ml.Value,mr.Value)

    } else {
        // Nothing to do, return e
    }

    return nil
}

// Section 5.2 Get an event that represents the maximum event value in a tree of values
func (event *Event) Max() *Event {
    e := event.Copy()

    if !e.IsLeaf {
        ml := e.Left.Max()
        mr := e.Right.Max()

        if e.Value >= ml.Value && e.Value >= mr.Value {
            e = &Event{
                Value:e.Value,
                IsLeaf:true,
                Right:nil,
                Left:nil,
            }
        } else if ml.Value >= mr.Value {
            e = &Event{
                Value:ml.Value,
                IsLeaf:true,
                Left:nil,
                Right:nil,
            }
        } else {
            e = &Event{
                Value:mr.Value,
                IsLeaf:true,
                Left:nil,
                Right:nil,
            }
        }
    } else {
        // Nothing, return e
    }

    return e
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
        e.Value = Max(event1.Value,event2.Value)
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
        event1.IsLeaf = false
        event1.Left = &a
        event1.Right = &b

        return event1.Join(event2)
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
        left := event1.Left.Join(event2.Lift(event2.Value - event1.Value))
        right := event1.Right.Join(event2.Lift(event2.Value - event1.Value))

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

// Simple min function over uint32
func Min(a uint32,b uint32) uint32 {
    if a>b {
        return b
    } else {
        return a
    }
}

// Simple max function over uint32
func Max(a uint32,b uint32) uint32 {
    if a>b {
        return a
    } else {
        return b
    }
}
