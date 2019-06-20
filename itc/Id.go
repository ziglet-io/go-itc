package itc

import (
    "fmt"
    "strings"
)

// Section 5.3.2 splits Ids used in the Fork operation
// Interestingly, this does NOT return a valid tree of Ids but instead two separate values
func (id *Id) Split() (*Id,*Id) {
    id1 := &Id{}
    id2 := &Id{}

    // Case 1 : split(0) -> (0,0)
    if id.IsLeaf && id.Value == 0 {
        id1.IsLeaf = true
        id2.IsLeaf = true

        return id1,id2
    }

    // Case 2: split(1) -> ((1,0),(0,1))
    if id.IsLeaf && id.Value == 1 {
        id1.IsLeaf = false
        id1.Left = &Id{
            IsLeaf: true,
            Value:  1,
        }
        id1.Right = &Id{
            IsLeaf:true,
            Value:0,
        }
        id2.Left = &Id{
            IsLeaf:true,
            Value:0,
        }
        id2.Right = &Id{
            IsLeaf:true,
            Value:1,
        }
        return id1,id2
    }

    // Case 3: split((0,i)) -> ((0,i1),(0,i2)) where (i1,i2) = split(i)
    if !id.IsLeaf && id.Left.IsLeaf && id.Left.Value ==0 {
        ida,idb := id.Right.Split()

        id1.IsLeaf = false
        id1.Left = NewId(0)
        id1.Right = ida

        id2.IsLeaf = false
        id2.Left = NewId(0)
        id2.Right = idb

        return id1,id2
    }

    // Case 4: split((i,0)) -> ((i1,0),(i2,0)) where (i1,i2) = split(i)
    if !id.IsLeaf && id.Right.IsLeaf && id.Right.Value ==0 {
        ida,idb := id.Left.Split()

        id1.Left = ida
        id1.Right.IsLeaf = true

        id2.Left = idb
        id2.Right.IsLeaf = true

        return id1,id2
    }

    // Case 5: split((i1,i2)) -> ((i1,0),(0,i2))
    if !id.IsLeaf {
        id1.Left = id.Left
        id1.Right.Value = 0
        id1.Right.IsLeaf = true

        id2.Left.Value = 0
        id2.Left.IsLeaf = true
        id2.Right = id.Right

        return id1,id2
    }

    return id1,id2

}

// Section 5.3.3 produce the sum of identifiers for use in Join
func (id1 *Id) Sum(id2 *Id) *Id{

    // Case 1: Sum(0,i) -> i
    if id1.IsLeaf && id1.Value==0 {
        return id2
    }

    // Case 2: Sum(i,0) -> i
    if id2.IsLeaf && id2.Value==0 {
        return id1
    }

    // TODO Seems like there are cases missing
    // (1,1) -> ??? This SHOULDN'T happen by the definition of split
    // (0,0) -> ??? Same

    // Case 3: Sum((l1,r1),(l2,r2)) -> Norm((Sum(l1,l2),Sum(r1,r2)))
    // AKA Sum(i1,i2)
    if !id1.IsLeaf && !id2.IsLeaf {
        fmt.Println(id1.Print())
        fmt.Println(id2.Print())

        i := &Id{
            Left: id1.Left.Sum(id2.Left),
            Right: id1.Right.Sum(id2.Right),
        }

        fmt.Println(i.Print())

        return i.Norm()
    }

    return nil
}

// Section 5.2 - There can be many representations of the same function. Reduce the function to a smaller representation.
func (id *Id) Norm() *Id {
    i := id.Copy()

    // Case 1: Norm((0,0)) -> 0
    if !id.IsLeaf && id.Left.IsLeaf && id.Left.Value==0 && id.Right.IsLeaf && id.Right.Value == 0 {
        i.IsLeaf=true
        i.Value=0
        i.Left = nil
        i.Right = nil

        return i
    }

    // Case 2: Norm((1,1)) -> 1
    if !id.IsLeaf && id.Left.IsLeaf && id.Left.Value==1 && id.Right.IsLeaf && id.Right.Value ==1 {
        i.Value = 1
        i.IsLeaf = true
        i.Left = nil
        i.Right = nil

        return i
    }

    // Case 3: Norm(i) -> i
    return i

}

// Produce a shallow copy of the Id
func (id *Id) Copy() *Id {
    return &Id{
        IsLeaf: id.IsLeaf,
        Value: id.Value,
        Left: id.Left,
        Right: id.Right,
    }
}

// Print a pretty version with parens
func (id *Id) Print() string {
    var sb strings.Builder

    if id.IsLeaf {
        sb.WriteString(fmt.Sprint(id.Value))
    } else {
        sb.WriteString("(")
        sb.WriteString(id.Left.Print())
        sb.WriteString(",")
        sb.WriteString(id.Right.Print())
        sb.WriteString(")")
    }

    return sb.String()
}

func NewId(v uint32) *Id {
    var id *Id
    if v != 0 {
        id = &Id{
            IsLeaf:true,
            Value: 1,
        }
    }   else {
        id = &Id {
            IsLeaf: true,
            Value: 0,
        }
    }
    return id
}

