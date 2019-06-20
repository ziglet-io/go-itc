package itc

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/ipfs/go-ipfs/thirdparty/assert"
	"github.com/saltatory/go-itc/itc"
	"testing"
)

func TestEventSimpleTree(t *testing.T){
	e := &itc.Event{
		IsLeaf: false,
		Value: 3,
		Left: itc.NewEvent(1),
		Right: itc.NewEvent(1),
	}

	expected := itc.NewEvent(4)
	assert.True(proto.Equal(e.Norm(),expected),t)

}

// Min
func TestEventMinSimple(t *testing.T){
	e := itc.NewEvent(3)
	assert.True(proto.Equal(e.Min(),e),t)
}

func TestEventMinTree(t *testing.T){
	e := &itc.Event{
		IsLeaf: false,
		Value: 1,
		Left: itc.NewEvent(2),
		Right: itc.NewEvent(3),
	}

	assert.True(proto.Equal(e.Min(),itc.NewEvent(1 + 2)),t)
}

// Max
func TestEventMaxSimple(t *testing.T){
	e := itc.NewEvent(3)
	assert.True(proto.Equal(e.Max(),e),t)
}

func TestEventMaxTree(t *testing.T){
	e := &itc.Event{
		IsLeaf: false,
		Value: 1,
		Left: itc.NewEvent(2),
		Right: itc.NewEvent(3),
	}

	fmt.Println(e.Print())

	assert.True(proto.Equal(e.Max(),itc.NewEvent(1 + 3)),t)
}

// Leq
func TestEventLeqFlat(t *testing.T){
	e1 := itc.NewEvent(1)
	e2 := itc.NewEvent(0)
	assert.True(e2.Leq(e1),t)
}

func TestEventLeqOneTree(t *testing.T){
	e1 := itc.NewEvent(1)
	e2 := &itc.Event{
		IsLeaf:false,
		Value:2,
		Right: itc.NewEvent(3),
		Left: itc.NewEvent(5),
	}

	assert.True(e1.Leq(e2),t)

	assert.False(e2.Leq(e1),t) // Change order
}

// Norm
func TestEventNormBasic(t *testing.T){
	e := itc.NewEvent(3)
	assert.True(proto.Equal(e,e.Norm()),t)
}

func TestEventNormEqual(t *testing.T){
	e := &itc.Event{
		IsLeaf:false,
		Value:2,
		Left: itc.NewEvent(3),
		Right: itc.NewEvent(3),
	}
	assert.True(proto.Equal(e.Norm(),itc.NewEvent(5)),t)
}

func TestEventNormTree(t *testing.T){
	l := itc.NewEvent(2)
	r := itc.NewEvent(3)
	e := &itc.Event{
		IsLeaf:false,
		Value:4,
		Left: l,
		Right: r,
	}

	expected := &itc.Event{
		IsLeaf: false,
		Value: 6,
		Left: itc.NewEvent(0),
		Right: itc.NewEvent(1),
	}

	assert.True(proto.Equal(e.Norm(),expected),t)
}

// Join
func TestEventJoinFlat(t *testing.T){
	l := itc.NewEvent(0)
	r := itc.NewEvent(1)
	expected := itc.NewEvent(1)

	assert.True(proto.Equal(l.Join(r),expected),t)
}

