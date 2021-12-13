package main

import (
	"github.com/jtrim-ons/ckmeans/pkg/ckmeans"
	"fmt"
	"io"
	"os"
)

func main() {
	var val float64
	var vals []float64
	nClusters := -1

	for {
		_, err := fmt.Scanf("%f\n", &val)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}
		if nClusters == -1 {
		    nClusters = int(val)
        } else {
            vals = append(vals, val)
        }
	}

    clusters, err := ckmeans.Ckmeans(vals, nClusters)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for i, cluster := range clusters {
        for _, val := range cluster {
            fmt.Println(i, val)
        }
    }
}
