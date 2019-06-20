package itc_test

import (
    "github.com/gogo/protobuf/proto"
    "github.com/ipfs/go-ipfs/thirdparty/assert"
    "github.com/saltatory/go-itc/itc"
    "testing"
)

// Seed Stamp
func TestStampSeedStamp(t *testing.T){
    s := &itc.Stamp{
        Id: itc.NewId(1),
        Event: itc.NewEvent(0),
    }
    assert.True(proto.Equal(itc.SeedStamp(),s),t)
}

// Leq
func TestStampFlat(t *testing.T){
    s1 := itc.SeedStamp()
    s2 := itc.SeedStamp()
    s2.Event.Value = 1

    assert.True(s1.Leq(s2),t)
}

func TestStampTreeRight(t *testing.T){
    s1 := itc.SeedStamp()
    e := &itc.Event{
        IsLeaf:false,
        Value:1,
        Left:itc.NewEvent(3),
        Right:itc.NewEvent(5),
    }
    s2 := itc.NewStamp(itc.NewId(0),e.Norm())

    assert.True(s1.Leq(s2),t)
}

func TestStampTreeLeft(t *testing.T){
    s2 := itc.SeedStamp()
    e := &itc.Event{
        IsLeaf:false,
        Value:0,
        Left:itc.NewEvent(3),
        Right:itc.NewEvent(5),
    }
    s1 := itc.NewStamp(itc.NewId(0),e)

    // Should be false because although the top value is 0, the lifted value is 3
    assert.False(s1.Leq(s2),t)

    // Test the comparison of lifted values
    s2.Event.Value = 1
    s1.Event.Value = 0
    s1.Event.Left.Value = 1
    s1.Event.Right.Value = 1
    assert.True(s1.Leq(s2),t)
}

// TODO Should test Case 4 but I'm not certain I understand it