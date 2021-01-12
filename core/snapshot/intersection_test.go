package snapshot

func TestIntersection_TopLayerUnicodeChar(t *testing.T) {
	testCases := []struct{
		in *Intersection
		exp string
	}{
		{
			in: &Intersection{
				Base: symbol.TopLeft
			}
		}
	}
}
