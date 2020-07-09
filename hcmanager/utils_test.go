package hcmanager


import (
	"testing"
	"fmt"
	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

func benchmarkConvertBytesToGb (input bmh.Capacity,b *testing.B){
	var inGb bmh.Capacity
	for i:=0; i<b.N; i++ {
		inGb =	ConvertBytesToGb(input)
	}
	fmt.Printf("Size %d GB\n",inGb)
}

func BenchmarkConvertBytesToGb(b *testing.B)  { benchmarkConvertBytesToGb(21474825484, b) }
func BenchmarkConvertBytesToGbNegative(b *testing.B) { benchmarkConvertBytesToGb(-21474825484, b) }
