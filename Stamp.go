package itc

// Section 5 Define the seed stamp, THE starting value
func SeedStamp() *Stamp {
    return &Stamp{
        Id: &Id{
            Value:1,
            IsLeaf:true,
            Left:nil,
            Right:nil,
        },
        Event: &Event{
            IsLeaf:true,
            Value:0,
            Left:nil,
            Right:nil,
        },
    }
}

// Section 5.3.1 Comparison - the basis for the partial order on stamps
func (s1 *Stamp) Leq(s2 *Stamp) bool {
    // Case 1: leq(n1,n2) -> n1 <= n2
    if s1.Event.IsLeaf && s2.Event.IsLeaf {
        return s1.Event.Value <= s2.Event.Value
    }

    // Case 2: leq(n1,(n2,l2,r2)) -> n1 <= n2 (only if normalized)
    if s1.Event.IsLeaf && !s2.Event.IsLeaf {
        return s1.Event.Value <= s2.Event.Value
    }

    // Case 3: leq((n1,l1,r1),n2) -> n1 <= n2 AND l1.Lift(n1).Leq(n2) AND r1.Lift(n1).Leq(n2)
    if !s1.Event.IsLeaf && s2.Event.IsLeaf {
        return s1.Event.Value <= s2.Event.Value && s1.Event.Left.Lift(s1.Event.Value).Leq(s2.Event) && s1.Event.Right.Lift(s1.Event.Value).Leq(s2.Event)
    }

    // Case 4: leq((n1,l1,r1),(n2,l2,r2)) -> n1 <= n2 AND l1.Lift(n1).Leq(l2.Lift(n2)) AND r1.Lift(n1).Leq(r2.Lift(n2)
    if !s1.Event.IsLeaf && !s2.Event.IsLeaf {
        return s1.Event.Value <= s2.Event.Value && s1.Event.Left.Lift(s1.Event.Value).Leq(s2.Event.Lift(s2.Event.Value)) && s1.Event.Right.Lift(s1.Event.Value).Leq(s2.Event.Right.Lift(s2.Event.Value))
    }

    // TODO better error handling?
    return false
}

// Section 5.3.2 Fork stamps for creating a new lineage of operations
func (stamp *Stamp) Fork() (*Stamp, *Stamp) {
    i1,i2 := stamp.Id.Split()

    e1 := stamp.Event.Copy()
    s1 := Stamp{
        Event:e1,
        Id:i1,
    }

    e2 := stamp.Event.Copy()
    s2 := Stamp{
        Event:e2,
        Id:i2,
    }
    return &s1,&s2
}

// Section 5.3.3 Join the lineages
func (s1 *Stamp) Join(s2 *Stamp) *Stamp {
    id := s1.Id.Sum(s2.Id)
    event := s1.Event.Join(s2.Event)
    stamp := Stamp{
        Id:id,
        Event:event,
    }

    return &stamp
}

func (stamp *Stamp) Fill() *Event {
    // Case 1: fill(0,e) -> e
    if stamp.Id.IsLeaf && stamp.Id.Value == 0 && !stamp.Event.IsLeaf{
        return stamp.Event
    }

    // Case 2: fill(1,e) -> max(e)
    if stamp.Id.IsLeaf && stamp.Id.Value ==1 && !stamp.Event.IsLeaf {
        event := stamp.Event.Max()
        return event
    }

    // Case 3: fill(i,n) -> n
    if !stamp.Id.IsLeaf && stamp.Event.IsLeaf {
        event := Event{
            IsLeaf:true,
            Value:stamp.Event.Value,
            Right:nil,
            Left:nil,
        }

        return &event
    }

    // Case 4: fill((1,ir),(n,el,er)) -> norm((n,max(max(el),min(eprimer)),eprimer)) where eprimer = fill(ir,er)
    if !stamp.Id.IsLeaf && stamp.Id.Left.Value==1 && !stamp.Event.IsLeaf {
        eprimer := &Event{}
        erprime := Stamp{
            Id:stamp.Id.Right,
            Event:stamp.Event.Right,
        }
        eprimer = erprime.Fill()

        var normable *Event

        normable.Value = stamp.Event.Value

        if stamp.Event.Left.Max().Value >= eprimer.Min().Value {
            normable.Left = stamp.Event.Left.Max()
        } else {
            normable.Right = eprimer.Min()
        }

        return normable.Norm()
    }

    // Case 5: fill((il,1),(n,el,er)) -> norm((n,elprime,max(max(er),min(elprime)))) where elprime=fill(il,el)

    // TODO
    return nil
}

func (stamp *Stamp) Grow() *Stamp {
    // TODO
}

// Section 5.3.4 Advance
// Called Event in the document but renamed to avoid name collision
func (stamp *Stamp) Advance() *Stamp {
    // TODO
}

// Produce a shallow copy of the stamp
func (stamp *Stamp) Copy() *Stamp {
    s := Stamp{
        Id: stamp.Id,
        Event: stamp.Event,
    }

    return &s
}