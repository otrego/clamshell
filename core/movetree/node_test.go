package movetree

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTraverse(t *testing.T) {
	testCases := []struct {
		desc string
		n    Node
		exp  []int
	}{
		{
			desc: "leaf node",
			n:    *NewNode(),
			exp:  []int{0},
		},
		{
			desc: "small tree",
			n: Node{
				varNum: 0,
				Children: []*Node{
					&Node{
						varNum: 1,
						Children: []*Node{
							&Node{
								varNum: 2,
							},
						},
					},
					&Node{
						varNum: 1,
						Children: []*Node{
							&Node{
								varNum: 2,
								Children: []*Node{
									&Node{
										varNum: 53,
									},
								},
							},
						},
					},
				},
			},
			exp: []int{0, 1, 1, 2, 2, 53},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := make([]int, 0)
			tc.n.Traverse(func(n *Node) {
				got = append(got, n.varNum)
			})
			if !cmp.Equal(got, tc.exp) {
				t.Errorf("got %v, expected %v", got, tc.exp)
			}
		})
	}
}

func TestTraverseMainBranch(t *testing.T) {

	testCases := []struct {
		desc string
		n    Node
		exp  []int
	}{
		{
			desc: "leaf node",
			n:    *NewNode(),
			exp:  []int{0},
		},
		{
			desc: "all 0th variation",
			n: Node{
				varNum: 0,
				Children: []*Node{
					&Node{
						varNum: 0,
						Children: []*Node{
							&Node{
								varNum: 0,
							},
						},
					},
					&Node{
						varNum: 0,
						Children: []*Node{
							&Node{
								varNum: 0,
								Children: []*Node{
									&Node{
										varNum: 0,
									},
								},
							},
						},
					},
				},
			},
			exp: []int{0, 0, 0, 0, 0, 0},
		},
		{
			desc: "small tree",
			n: Node{
				varNum: 0,
				Children: []*Node{
					&Node{
						varNum: 0,
						Children: []*Node{
							&Node{
								varNum: 1,
							},
						},
					},
					&Node{
						varNum: 0,
						Children: []*Node{
							&Node{
								varNum: 1,
								Children: []*Node{
									&Node{
										varNum: 0,
									},
								},
							},
						},
					},
				},
			},
			exp: []int{0, 0, 0},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := make([]int, 0)
			tc.n.TraverseMainBranch(func(n *Node) {
				got = append(got, n.varNum)
			})
			if !cmp.Equal(got, tc.exp) {
				t.Errorf("got %v, expected %v", got, tc.exp)
			}
		})
	}
}
