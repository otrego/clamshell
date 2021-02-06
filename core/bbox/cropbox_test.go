package bbox

import (
	"reflect"
	"testing"

	"github.com/otrego/clamshell/core/point"
)

func TestCropBoxFromPreset(t *testing.T) {
	testCases := []struct {
		desc           string
		croppingPreset CroppingPreset
		boardSize      int
		exp            *CropBox
	}{
		{
			desc:           "All",
			croppingPreset: All,
			boardSize:      19,
			exp: &CropBox{
				BBox:         &BoundingBox{tl: point.New(0, 0), br: point.New(19, 19)},
				OriginalSize: 19,
			},
		},
		{
			desc:           "Left",
			croppingPreset: Left,
			boardSize:      19,
			exp: &CropBox{
				BBox:         &BoundingBox{tl: point.New(0, 0), br: point.New(19, 10)},
				OriginalSize: 19,
			},
		},
		{
			desc:           "Right",
			croppingPreset: Right,
			boardSize:      19,
			exp: &CropBox{
				BBox:         &BoundingBox{tl: point.New(0, 8), br: point.New(19, 19)},
				OriginalSize: 19,
			},
		},
		{
			desc:           "Top",
			croppingPreset: Top,
			boardSize:      19,
			exp: &CropBox{
				BBox:         &BoundingBox{tl: point.New(0, 0), br: point.New(10, 19)},
				OriginalSize: 19,
			},
		},
		{
			desc:           "Bottom",
			croppingPreset: Bottom,
			boardSize:      19,
			exp: &CropBox{
				BBox:         &BoundingBox{tl: point.New(8, 0), br: point.New(19, 19)},
				OriginalSize: 19,
			},
		},
		{
			desc:           "TopLeft",
			croppingPreset: TopLeft,
			boardSize:      19,
			exp: &CropBox{
				BBox:         &BoundingBox{tl: point.New(0, 0), br: point.New(10, 11)},
				OriginalSize: 19,
			},
		},
		{
			desc:           "TopRight",
			croppingPreset: TopRight,
			boardSize:      19,
			exp: &CropBox{
				BBox:         &BoundingBox{tl: point.New(0, 7), br: point.New(10, 19)},
				OriginalSize: 19,
			},
		},
		{
			desc:           "BottomLeft",
			croppingPreset: BottomLeft,
			boardSize:      19,
			exp: &CropBox{
				BBox:         &BoundingBox{tl: point.New(8, 0), br: point.New(19, 11)},
				OriginalSize: 19,
			},
		},
		{
			desc:           "BottomRight",
			croppingPreset: BottomRight,
			boardSize:      19,
			exp: &CropBox{
				BBox:         &BoundingBox{tl: point.New(8, 7), br: point.New(19, 19)},
				OriginalSize: 19,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := CropBoxFromPreset(tc.croppingPreset, tc.boardSize)

			if err != nil {
				return
			}

			if !reflect.DeepEqual(got, tc.exp) {
				t.Errorf("got %v%v, expected %v%v", got.BBox.TopLeft(), got.BBox.BotRight(), tc.exp.BBox.TopLeft(), tc.exp.BBox.BotRight())
			}

		})
	}

}
