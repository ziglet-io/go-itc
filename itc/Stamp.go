package itc

import "github.com/gogo/protobuf/proto"

const GrowIncrement uint32 = 1000

// Section 5 Define the seed stamp, THE starting value
func SeedStamp() *Stamp {
    return &Stamp{
        Id: NewId(1),
        Event: NewEvent(0),
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
    stamp := NewStamp(id,event)

    return stamp
}

// Section 5.3.4 During Advance, attempt to simplify the event tree
func (stamp *Stamp) Fill() *Event {
    // Case 1: fill(0,e) -> e
    if stamp.Id.IsLeaf && stamp.Id.Value == 0 {
        return stamp.Event
    }

    // Case 2: fill(1,e) -> max(e)
    if stamp.Id.IsLeaf && stamp.Id.Value ==1 {
        return stamp.Event.Max()
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

    // Case 4: fill((1,ir),(n,el,er)) -> norm((n,max(max(el),min(erprime)),erprime)) where erprime = fill(ir,er)
    if !stamp.Id.IsLeaf && stamp.Id.Left.Value ==1 && !stamp.Event.IsLeaf {
        s := &Stamp{
            Id:stamp.Id.Right.Copy(),
            Event:stamp.Event.Right.Copy(),
        }

        erprime := s.Fill()

        e := &Event{
            IsLeaf: false,
            Value: stamp.Event.Value,
            Left: Max(stamp.Event.Left.Max(),erprime.Min()),
            Right: erprime,
        }

        return e.Norm()
    }

    // Case 5: fill((il,1),(n,el,er)) -> norm((n,elprime,max(max(er),min(elprime)))) where elprime=fill(il,el)
    if !stamp.Id.IsLeaf && stamp.Id.Right.Value ==1 && !stamp.Event.IsLeaf {
        s := &Stamp{
            Id:stamp.Id.Left.Copy(),
            Event:stamp.Event.Left.Copy(),
        }

        elprime := s.Fill()

        e := &Event{
            IsLeaf: false,
            Value: stamp.Event.Value,
            Left: elprime,
            Right: Max(stamp.Event.Right.Max(),elprime.Min()),
        }

        return e.Norm()
    }

    // Error
    return nil
}

// Section 5.3.4 During Advance, when fill is not possible, grow the event tree
// Returns a new event tree and a cost for that tree
func (stamp *Stamp) Grow() (*Event,uint32) {

    // Case 1: grow(1,n) -> (n+1,0)
    if stamp.Id.IsLeaf && stamp.Id.Value ==1 && stamp.Event.IsLeaf {
        e := &Event{
            IsLeaf:true,
            Value:stamp.Event.Value+1,
        }
        return e,0
    }

    // Case 2: grow(i,n) -> (eprime,c + N) where (eprime,c) = grow(i,(n,0,0))
    if !stamp.Id.IsLeaf && stamp.Event.IsLeaf {
        el := &Event{
            IsLeaf: true,
            Value: 0,
        }
        er := el.Copy()
        e := &Event{
            IsLeaf: false,
            Value: stamp.Event.Value,
            Left: el,
            Right: er,
        }

        s := &Stamp{
            Id: stamp.Id,
            Event: e,
        }

        eprime,c := s.Grow()

        return eprime,c+GrowIncrement
    }

    // Case 3: grow((0,ir),(n,el,er)) -> ((n,el,erprime),cr+1) where (erprime,cr) = grow(ir,er)
    if !stamp.Id.IsLeaf && stamp.Id.Left.IsLeaf && stamp.Id.Left.Value == 0 && !stamp.Event.IsLeaf {
        s := &Stamp{
            Id: stamp.Id.Right,
            Event: stamp.Event.Right,
        }
        erprime,cr := s.Grow()

        e := &Event{
            IsLeaf: false,
            Value: stamp.Event.Value,
            Left: stamp.Event.Left.Copy(),
            Right: erprime,
        }

        return e,cr+1
    }

    // Case 4: grow((il,0),(n,el,er)) -> ((n,elprime,er),cl+1) where (elprime,cl) = grow(il,el)
    if !stamp.Id.IsLeaf && stamp.Id.Right.IsLeaf && stamp.Id.Right.Value == 0 && !stamp.Event.IsLeaf {
        s := &Stamp{
            Id: stamp.Id.Left,
            Event: stamp.Event.Left,
        }
        elprime,cl := s.Grow()

        e := &Event{
            IsLeaf: false,
            Value: stamp.Event.Value,
            Left: elprime,
            Right: stamp.Event.Right.Copy(),
        }

        return e,cl+1
    }

    // Case 5: grow((il,ir),(n,el,er)) ->
    // ((n,elprime,er),cl+1)    if cl < cr
    // ((n,el,erprime),cr+1)    if cl >= cr
    // where (elprime,cl) = grow(il,el)
    // and (erprime,cr) = grow(ir,er)
    if !stamp.Id.IsLeaf && !stamp.Event.IsLeaf {
        sl := &Stamp{
            Id: stamp.Id.Left,
            Event: stamp.Event.Left,
        }
        elprime,cl := sl.Grow()

        sr := &Stamp{
            Id: stamp.Id.Right,
            Event : stamp.Event.Right,
        }
        erprime,cr := sr.Grow()

        if cl < cr {
            e := &Event{
                IsLeaf: false,
                Value: stamp.Event.Value,
                Left: elprime,
                Right: stamp.Event.Right,
            }

            return e,cl+1
        } else {
            e := &Event{
                IsLeaf: false,
                Value: stamp.Event.Value,
                Left: stamp.Event.Left,
                Right: erprime,
            }

            return e,cr+1
        }
    }

    return nil,0
}

// Section 5.3.4 Advance
// Called "Event" in the document but renamed to avoid name collision
func (stamp *Stamp) Advance() *Stamp {
    e := stamp.Fill()

    if !proto.Equal(e,stamp.Event){
        return stamp.Copy()
    } else {
        e,_ := stamp.Grow()
        return NewStamp(stamp.Id,e)
    }
}

// Produce a shallow copy of the stamp
func (stamp *Stamp) Copy() *Stamp {
    s := Stamp{
        Id: stamp.Id,
        Event: stamp.Event,
    }

    return &s
}

// Create a new stamp with premade Id and Event
func NewStamp(id *Id, event *Event) *Stamp {
    return &Stamp{
        Id: id.Copy(),
        Event: event.Copy(),
    }
}