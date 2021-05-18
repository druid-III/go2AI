// mapAsGraph project main.go
package main

import (
	"fmt"
	"math/rand"
	"time"

	gc "./sOrGraph"
)

// -------------------------------------------------------------------------------
func main() {

	wkr := func(n uint64) error {

		fmt.Println(n)

		return nil
	}

	//............................................ обязательная инициализация карт
	o_gn := gc.GraphInit()
	/*
	   ...................................................................................
	         полный граф на 5 нод - пентаграмма
	   ...................................................................................
	*/

	//.......................... присоеденяем заполненный узел графа к колоде карт
	o_gn.AddNode(0, []uint64{1, 2, 3, 4}, [4]float32{0., 0.1, 2., 17.}, 1)
	o_gn.AddNode(1, []uint64{0, 2, 3, 4}, [4]float32{0., 0.1, 2., 17.}, 1)
	o_gn.AddNode(2, []uint64{1, 0, 3, 4}, [4]float32{0., 0.1, 2., 17.}, 1)
	o_gn.AddNode(3, []uint64{1, 2, 0, 4}, [4]float32{0., 0.1, 2., 17.}, 1)
	o_gn.AddNode(4, []uint64{1, 2, 3, 0}, [4]float32{0., 0.1, 2., 17.}, 1)

	o_gn.PrintGraphNetInfo()
	o_gn.PrintNodeInfo(2)
	fmt.Println("=====================================================In Depth=1=")
	o_gn.REC_StepsInDepth(0, 2, wkr)
	fmt.Println("======================================================In Wide=1=")
	o_gn.StepsInWide(0, 2, wkr)
	fmt.Println("========================================================Chain=1=")
	o_gn.ChainReaction(0, wkr)
	N, cp := o_gn.CheckOnAcyclicity(0)
	fmt.Println(N, cp)
	/*
	   ...................................................................................
	                 *->|2|->|5|->|8|->|11|
	            S    |
	        *->|0|->|1|->|3|->|6|->|9|->|12|
	        |        |
	        |        *->|4|->|7|->|10|->|13|
	        |                            |
	        *----------------------------*
	   ...................................................................................
	*/
	//............................................ обязательная инициализация карт
	o_gn = gc.GraphInit()
	//.......................... присоеденяем заполненный узел графа к колоде карт
	o_gn.AddNode(0, []uint64{1}, [1]float32{0.}, 1)
	o_gn.AddNode(1, []uint64{2, 3, 4}, [1]float32{0.}, 1)
	o_gn.AddNode(2, []uint64{5}, [1]float32{0.}, 1)
	o_gn.AddNode(3, []uint64{6}, [1]float32{0.}, 1)
	o_gn.AddNode(4, []uint64{7}, [1]float32{0.}, 1)
	o_gn.AddNode(5, []uint64{8}, [1]float32{0.}, 1)
	o_gn.AddNode(6, []uint64{9}, [1]float32{0.}, 1)
	o_gn.AddNode(7, []uint64{10}, [1]float32{0.}, 1)
	o_gn.AddNode(8, []uint64{11}, [1]float32{0.}, 1)
	o_gn.AddNode(9, []uint64{12}, [1]float32{0.}, 1)
	o_gn.AddNode(10, []uint64{13}, [1]float32{0.}, 1)
	o_gn.AddNode(11, []uint64{}, [1]float32{0.}, 1)
	o_gn.AddNode(12, []uint64{}, [1]float32{0.}, 1)
	o_gn.AddNode(13, []uint64{0}, [1]float32{0.}, 1)

	o_gn.PrintGraphNetInfo()
	o_gn.PrintNodeInfo(2)
	fmt.Println("=====================================================In Depth=2=")
	o_gn.REC_StepsInDepth(0, 3, wkr)
	fmt.Println("======================================================In Wide=2=")
	o_gn.StepsInWide(0, 2, wkr)
	fmt.Println("========================================================Chain=2=")
	o_gn.ChainReaction(0, wkr)
	N, cp = o_gn.CheckOnAcyclicity(0)
	fmt.Println(N, cp)
	/*
	   ...................................................................................
	                 *->|2|--------*
	            S    |             |
	        *->|0|->|1|->|4|->|7|<-*
	        |        |
	        |        *->|3|->|12|
	        |                 |
	        *-----------------*
	   ...................................................................................
	*/
	//............................................ обязательная инициализация карт
	o_gn = gc.GraphInit()
	//.......................... присоеденяем заполненный узел графа к колоде карт
	o_gn.AddNode(0, []uint64{1}, [1]int{0}, 1)
	o_gn.AddNode(1, []uint64{2, 4, 3}, [1]float32{0.}, 1)
	o_gn.AddNode(2, []uint64{7}, [1]float32{0.}, 1)
	o_gn.AddNode(3, []uint64{12}, [1]string{"folio"}, 1)
	o_gn.AddNode(4, []uint64{7}, [1]float32{0.}, 1)
	o_gn.AddNode(12, []uint64{0}, [1]float32{0.}, 1)

	o_gn.PrintGraphNetInfo()
	o_gn.PrintNodeInfo(1)
	fmt.Println("=====================================================In Depth=3=")
	o_gn.REC_StepsInDepth(0, 1, wkr)
	fmt.Println("======================================================In Wide=3=")
	o_gn.StepsInWide(0, 3, wkr)
	fmt.Println("========================================================Chain=3=")
	o_gn.ChainReaction(0, wkr)

	//	o_gn = gc.GraphInit()
	o_gn.PrintNodeInfo(7)
	o_gn.AddConnectionsIfExist(7, []uint64{2}, true)

	N, cp = o_gn.CheckOnAcyclicity(0)
	fmt.Println(N, cp)
}

type targetFunction4P2PTransition func(uint) uint // for pool of functions

//........................................................................... /\/ .
var m_tfPool = map[string]targetFunction4P2PTransition{
	//----------------------------------------------------------------------- /\/ -
	"void": func(uint) uint { return 0 },
	//----------------------------------------------------------------------- /\/ -
	"random": func(uint) uint {

		rand.Seed(time.Now().UnixNano())

		return uint(rand.Uint32())
	},
	//----------------------------------------------------------------------- /\/ -
	"nextWiht": func(uint) uint {

		return 0
	},
}
