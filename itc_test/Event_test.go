package itc

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/ipfs/go-ipfs/thirdparty/assert"
	"github.com/saltatory/go-itc/itc"
	"testing"
)

func TestEventNormN(t *testing.T){
	e := &itc.Event{
		IsLeaf: true,
		Value: 3,
	}
	assert.True(proto.Equal(e.Norm(),e),t)
}

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


