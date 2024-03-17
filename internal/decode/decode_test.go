package decode

import (
	"fmt"
	"k8s.io/klog"
	"math"
	"testing"
)

func TestDecode(t *testing.T) {
	//ctx := context.TODO()
	//file, _ := os.Open("/Users/ted/workspace/go/music_player/1.flac")
	//d, err := newFLACDecoder(ctx, file, 0, &music.Callback{
	//	CurTime: func(duration time.Duration) {
	//
	//	},
	//	DoneFn: func(status model.Status) {
	//
	//	},
	//})
	//if err != nil {
	//	t.Fatal(err)
	//}
	//d.Play()
	// -1 50%
	// -2 25%
	// -3 12.5%
	// -4 6.25%
	
	var slice [101]float64
	for idx := range slice {
		if idx == 0 {
			continue
		}
		for i := 660; i >= 0; i-- {
			cnt := math.Pow(2, -0.01*float64(i))
			if int(cnt/0.01) == idx {
				slice[idx] = -0.01 * float64(i)
				break
			}
		}
	}
	
	for idx, v := range slice {
		klog.Info(idx, math.Pow(2, v), v)
	}
	
	for _, v := range slice {
		
		fmt.Printf("%.2f,", v)
	}
	
}
