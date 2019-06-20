package itc_test

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/ipfs/go-ipfs/thirdparty/assert"
	"github.com/saltatory/go-itc/itc"
	"testing"
)

// Split

func TestIdSplitZero(t *testing.T) {
	i := &itc.Id{
		IsLeaf: true,
		Value:  0,
	}
	a, b := i.Split()

	expectedA := &itc.Id{
		IsLeaf: true,
		Value:  0,
	}
	expectedB := expectedA.Copy()

	assert.True(proto.Equal(a, expectedA), t)
	assert.True(proto.Equal(b, expectedB), t)
}

func TestIdPrint(t *testing.T) {
	i := &itc.Id{
		Left: &itc.Id{
			IsLeaf: true,
			Value:  1,
		},
		Right: &itc.Id{
			IsLeaf: true,
			Value:  0,
		},
	}
	fmt.Println(i.Print())
}

func TestIdSplitOne(t *testing.T) {
	i := &itc.Id{
		IsLeaf: true,
		Value:  1,
	}
	a, b := i.Split()

	expectedA := &itc.Id{
		IsLeaf: false,
		Left: &itc.Id{
			IsLeaf: true,
			Value:  1,
		},
		Right: &itc.Id{
			IsLeaf: true,
			Value:  0,
		},
	}
	expectedB := &itc.Id{
		IsLeaf: false,
		Left: &itc.Id{
			IsLeaf: true,
			Value:  0,
		},
		Right: &itc.Id{
			IsLeaf: true,
			Value:  1,
		},
	}

	assert.True(proto.Equal(a, expectedA), t)
	assert.True(proto.Equal(b, expectedB), t)
}

func TestIdSplitZeroOne(t *testing.T) {
	i := &itc.Id{
		IsLeaf: false,
		Value:  0,
		Left: &itc.Id{
			IsLeaf: true,
			Value:  0,
		},
		Right: &itc.Id{
			IsLeaf: true,
			Value:  1,
		},
	}

	// Money shot
	a, b := i.Split()

	// TODO Building up Id trees is annoying
	expectedB := &itc.Id{
		IsLeaf: false,
		Left: &itc.Id{
			IsLeaf: true,
			Value:  0,
		},
		Right: &itc.Id{
			IsLeaf: false,
			Left:   itc.NewId(0),
			Right:  itc.NewId(1),
		},
	}
	expectedA := &itc.Id{
		IsLeaf: false,
		Left:   &itc.Id{IsLeaf: true},
		Right: &itc.Id{
			IsLeaf: false,
			Left:   itc.NewId(1),
			Right:  itc.NewId(0),
		},
	}

	assert.True(proto.Equal(a, expectedA), t)
	assert.True(proto.Equal(b, expectedB), t)
}

// Sum
func TestIdSumSimple(t *testing.T) {
	id1 := itc.NewId(0)
	id2 := itc.NewId(1)

	result := id1.Sum(id2)

	assert.True(proto.Equal(id2,result),t)
}

func TestIdSumTree(t *testing.T){
	id1 := itc.NewId(1)
	id2 := &itc.Id{
		IsLeaf:false,
	}
	id2.Left,id2.Right = id1.Split()

	result := id2.Left.Sum(id2.Right)
	fmt.Println(result.Print())
	assert.True(proto.Equal(result,itc.NewId(1)),t)
}

func TestIdSumDual(t *testing.T){
	id1 := itc.NewId(1)
	id2 := itc.NewId(0)
	result := id1.Sum(id2)
	assert.True(proto.Equal(result,itc.NewId(1)),t)
}

// Norm
func TestIdNormZero(t *testing.T) {
	id := &itc.Id{
		IsLeaf: false,
		Left: itc.NewId(0),
		Right: itc.NewId(0),
	}
	assert.True(proto.Equal(id.Norm(),itc.NewId(0)),t)
}

func TestIdNormOne(t *testing.T) {
	id := &itc.Id{
		IsLeaf: false,
		Left: itc.NewId(1),
		Right: itc.NewId(1),
	}
	assert.True(proto.Equal(id.Norm(),itc.NewId(1)),t)
}

func TestIdNormId(t *testing.T){
	id := itc.NewId(1)

	assert.True(proto.Equal(id.Norm(),id),t)
}
