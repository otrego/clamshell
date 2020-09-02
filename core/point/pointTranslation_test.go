package point

import (
	"fmt"
	"testing"
)

// Two sample correct results
// New(8, 20), SGF string "iu"
// New(13, 5), SGF string "nf"

func TestTranslate(t *testing.T) {
	
	fmt.Println()
	fmt.Println("[*** Go Game Point Translation Tests ***]")
	
	// Test Point Struct
	type PointStrct struct {
		x int64
		y int64
	}
	// Test SGF Point
	var SGFCrct01 string
	// Constant value to convert ascii to letter and int
	const aValue = int64('a')
	// Test *Point
	var x1 int64 = 8
	var y1 int64 = 20
	Pnt01 := New(x1, y1)
	// Test Control Point, is non-pointer type so can locally verify
	PntCrct01 := PointStrct{
		x: Pnt01.x,
		y: Pnt01.y,
	}
	// Convert *Point to non-pointer type for local verification
	PntTest01 := PointStrct{
		x: Pnt01.x,
		y: Pnt01.y,
	}
	
	// *** To Test/Verify *Point creation for below tests
	if PntTest01 != PntCrct01 {
		fmt.Println()
		t.Errorf("Point creation fail.\n"+
			"got: %v \n"+
			"want: %v \n", Pnt01, PntCrct01)
	} else {
		fmt.Println("Point creation PASS. ")
	}
	
	// *** To Test ToSGF method
	// Test SGF
	SGF01 := Pnt01.ToSGF()
	// Test control SGF, as a string "xy", carrying above *Point values
	SGFCrct01 = string(rune((PntCrct01.x)+aValue)) + string(rune((
		PntCrct01.y)+aValue))
	// To interject a different value, comment out the above
	/*SGFCrct01 = "ju"*/
	
	if SGF01 != SGFCrct01 {
		fmt.Println()
		t.Errorf("To SGF translation fail.\n"+
			"got: %v \n"+
			"want: %v \n", SGF01, SGFCrct01)
	} else {
		fmt.Println("To SGF translation PASS. ")
	}
	
	// *** To Test *Point method, carrying above values
	// Test *Point
	Pnt02 := NewFromSGF(SGF01)
	// To interject a different value, comment out the above
	/*Pnt02 := PointStrct{
		x: 1,
		y: 24,
	}*/
	// Convert *Point to non-pointer type for local verification
	PntTest02 := PointStrct{
		x: Pnt02.x,
		y: Pnt02.y,
	}
	
	if PntTest02 != PntCrct01 {
		fmt.Println()
		t.Errorf("To Point translation fail.\n"+
			"got: %v \n"+
			"want: %v \n", Pnt02, PntCrct01)
	} else {
		fmt.Println("To Point translation PASS. ")
	}
}
