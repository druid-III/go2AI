package sOrGraph

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type Nothing struct{} //                                                          тип структура из ничего - имеет "0"й размер ! - это не структура указателей на Nil

var Null Nothing

type orGraphsNodes struct {
	conID     map[uint64]Nothing //                                               колода уникальных (! повторы исключены !) беззнаковых чисел - направленных из этого узла ребер
	data      interface{}        //                                               данные узла (ноды) - любой тип, любой размер
	flagColor uint64             //                                               цвет флага - 18_446_744_073_709_551_615 комбинаций(that is, hexadecimal 0xFFFFFFFFFFFFFFFF) - НЕ ОБЯЗАН СОВПАДАТЬ С ID хотя количество комбинаций совпадает
}

type orGraphsNetwork map[uint64]orGraphsNodes //                                  тип колода уникальных ID с произвольными структурами(как колода разных карт)

func MapCopy(src map[interface{}]interface{}) map[interface{}]interface{} {

	targ := make(map[interface{}]interface{})

	for k, v := range src {
		targ[k] = v
	}

	return targ
}

// -------------------------------------------------------------------------------
func GraphInit() orGraphsNetwork {

	pt := map[uint64]orGraphsNodes{}

	return pt
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) NodeExists(index uint64) bool {

	if _, ok := gm[index]; ok {

		return true

	} else {

		return false
	}
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) ConnectionInNodeExists(index uint64, conn uint64) bool {

	if _, ok := gm[index].conID[conn]; ok {

		return true

	} else {

		return false
	}
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) ConnectionsAsSlice(index uint64) []uint64 {

	s4c := make([]uint64, len(gm[index].conID))
	k := 0

	for i, _ := range gm[index].conID { //                                         перебираем колоду пустых карт соединений(направленных ребер) и добавляем их ID как ЗНАЧЕНИЯ слайса

		s4c[k] = i
		k++
	}
	return s4c
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) AddConnections(index uint64, adID []uint64) {

	for _, v := range adID { //                                                    перебираем срез и его значения станут индексами колоды пустых карт соединений(направленных ребер)
		gm[index].conID[v] = Null
	}

	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) RemoveConnection(index uint64, rmID []uint64) {

	for _, v := range rmID { //                                                    перебираем срез и его значения станут индексами колоды пустых карт которые удаляются
		delete(gm[index].conID, v)
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) AddConnectionsIfExist(index uint64, adID []uint64, strictly bool) error {

	if gm.NodeExists(index) { //                                                  если есть запрошенный узел...

		for _, v := range adID {

			if gm.NodeExists(v) { //                                               ...и если есть узлы для запрошенных связей

				gm[index].conID[v] = Null
			} else {

				if strictly { //                                                   если флаг выводить ошибку поднят

					return errors.New("can't find a node for connection")
				}
			}
		}
	} else {

		if strictly { //                                                           если флаг выводить ошибку поднят
			return errors.New("can't find a node for work")
		}
	}
	return nil
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) RemoveConnectionsIfExist(index uint64, rmID []uint64, strictly bool) error {

	if gm.NodeExists(index) { //                                                  если есть запрошенный узел...

		for _, v := range rmID { //                                               перебираем срез и его значения станут индексами колоды пустых карт которые удаляются
			delete(gm[index].conID, v)
		}

	} else {

		if strictly { //                                                          если флаг выводить ошибку поднят
			return errors.New("can't find a node for work")
		}
	}
	return nil
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) AddNode(index uint64, conID []uint64, data interface{}, flagColor uint64) {

	var temp orGraphsNodes //                                                     новый узел с "0"-ми полями

	temp.conID = map[uint64]Nothing{}

	for _, v := range conID { //                                                  перебираем срез и его значения станут индексами колоды пустых карт соединений(направленных ребер)
		temp.conID[v] = Null
	}
	temp.data = data
	temp.flagColor = flagColor

	gm[index] = temp

	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) RemoveNode(index uint64) {

	delete(gm, index)

	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) AddNodeAndCheck(index uint64, conID []uint64, data interface{}, flagColor uint64) {

	var temp orGraphsNodes //                                                     новый узел с "0"-ми полями

	temp.conID = map[uint64]Nothing{}

	for _, v := range conID { //                                                  перебираем срез и его значения станут индексами колоды пустых карт соединений(направленных ребер)
		temp.conID[v] = Null
	}
	temp.data = data
	temp.flagColor = flagColor

	gm[index] = temp

	for i, _ := range gm[index].conID { //                                        перебираем колоду карт связей

		if !gm.NodeExists(i) { //                                                 переходим по каждой, и если такого узла нет то затираем связь
			delete(gm[index].conID, i)
		}
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) RemoveNodeAndCheck(index uint64) {

	delete(gm, index)

	for i, _ := range gm { //                                                     перебераем ВСЕ узлы

		for k, _ := range gm[i].conID { //                                        перебираем ВСЕ их соединения

			if k == index {
				delete(gm[i].conID, index) //                                     и затираем все соединения на этот узел
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) WedgeNodeAndCheck(index uint64, nodes map[uint64][]uint64) {

	for n, v := range nodes {

		for _, c := range v {

			delete(gm[n].conID, c)    //                                           разрываем перечисленные в срезе соединения в текущем узле из перебираемых...
			gm[index].conID[c] = Null //                                           ...но ВСЕ(каждой из перебираемых нод) разорванные связи теперь будут вести из этого узла...
		}
		gm[n].conID[index] = Null //                                               ...затем в каждую ноду добавляем соединение с этим узлом
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) EradicateNodeAndCheck(index uint64) {

	for i, _ := range gm { //                                                     перебераем ВСЕ узлы

		for k, _ := range gm[i].conID { //                                        перебираем ВСЕ их соединения

			if k == index {

				delete(gm[i].conID, index) //                                     затираем все соединения на этот узел...

				for j, _ := range gm[index].conID {

					gm[i].conID[j] = Null //                                      ...в i-й узел ранее указывавший на эту ноду добавляем все связи из нее
				}
			}
		}
	}

	delete(gm, index)

	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) VerifiAllConnection() {

	for i, _ := range gm { //                                                     перебераем все узлы

		for k, _ := range gm[i].conID { //                                        перебираем все их соединения

			if !gm.NodeExists(k) {
				delete(gm[i].conID, k) //                                         и затираем связь с несуществующим узлом
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) AssociationForGraph(ex orGraphsNetwork) {

	for i, v := range ex {

		if !gm.NodeExists(i) { //                                                  если есть общая нода...

			gm[i] = v

		} else {

			gm.AddConnections(i, ex.ConnectionsAsSlice(i))
		}
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) IntersectionForGraph(ex orGraphsNetwork) {

	for i, _ := range ex {

		if !gm.NodeExists(i) { //                                                 если есть общая нода...

			delete(gm, i)

		} else {

			for j, _ := range ex[i].conID {

				if !gm.ConnectionInNodeExists(i, j) { //                          если есть общая дуга...

					delete(gm[i].conID, j)
				}
			}
		}
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) DifferenceForGraph(ex orGraphsNetwork) {

	for i, _ := range ex {

		if gm.NodeExists(i) { //                                                  если есть общая нода...

			delete(gm, i)

		}
	}
	gm.VerifiAllConnection()

	return
}

// -------------------------------------------------------------------------------
/*
func (gm orGraphsNetwork) MultiplicationForGraph(ex orGraphsNetwork) {

	return
}
*/
// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) ReadValues(index uint64) interface{} {

	return gm[index].data
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) SetValues(index uint64, data interface{}) {

	temp := gm[index] //                                                          присваивать весь узел придется в любом случае - мы работаем с копиями, а инициализация нового узла даст нулевые поля
	temp.data = data
	gm[index] = temp //                                                           копия во внешний узел

	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) SetFlag(index uint64, flg uint64) {

	temp := gm[index] //                                                          присваивать весь узел придется в любом случае - мы работаем с копиями, а инициализация нового узла даст нулевые поля
	temp.flagColor = flg
	gm[index] = temp //                                                           копия во внешний узел

	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) ResetAllFlags() {

	for i, _ := range gm { //                                                     перебираем колоду карт связей

		temp := gm[i] //                                                          присваивать весь узел придется в любом случае - мы работаем с копиями, а инициализация нового узла даст нулевые поля
		temp.flagColor = 0
		gm[i] = temp //                                                           копия во внешний узел
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) SetAllValuesFromSlice(s_data []interface{}) { //        только однородные данные

	for i, v := range s_data { //                                                 перебираем все доступные ячейки с данными, индекс станет ключем карты, а значение полем значения

		temp := gm[uint64(i)] //                                                  !срез может иметь отрицательную индексацию!
		temp.data = v
		gm[uint64(i)] = temp //                                                   !срез может иметь отрицательную индексацию!
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) SetAllValuesFromMap(m_data map[uint64]interface{}) { // поддерживает разнородные данные

	for i, v := range m_data { //                                                 перебираем все доступные ячейки с данными, индекс станет ключем карты, а значение полем значения

		temp := gm[i]
		temp.data = v
		gm[i] = temp
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) NextNodeID(index uint64,
	TF func(...interface{}) uint64, farg ...interface{}) uint64 {

	return TF(farg) //                                                             просто заворачиваем "целевую функцию перехода к следующему узлу" привязывая ее к объекту
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) REC_StepsInDepth(index, setColor uint64, work4node func(uint64) error) {

	gm.ResetAllFlags() //                                                          все узлы непосещены

	var REC func(uint64)

	REC = func(n uint64) { // >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> R

		gm.SetFlag(n, setColor) //                                                 посетили узел
		work4node(n)            //                                                 обработали узел

		for i, _ := range gm[n].conID { //                                         перебираем все карты связей, идекс и есть их ID

			if (gm[i].flagColor == 0) && gm.NodeExists(i) { //                     если узел непосещен, то...

				REC(i)
			}
		}
		return
	} // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< R
	if gm.NodeExists(index) {

		REC(index)
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) StepsInWide(index, setColor uint64, work4node func(uint64) error) {

	gm.ResetAllFlags() //                                                         все узлы непосещены

	bucket4Curr := map[uint64]map[uint64]Nothing{} //                             корзина колод карт текущего уровня обработки - над каждым узлом совершается элементарное действие обработчика
	bucket4Next := map[uint64]map[uint64]Nothing{} //                             корзина колод карт для следующего уровня обработки

	bucket4Curr[index] = gm[index].conID //                                       единственный элемент в корзине на старте - соединения единственного входного узла, причем индекс колоды и узла одинаковы!

	if gm.NodeExists(index) {

		gm.SetFlag(index, setColor) //                                            этот узел(с ID колоды) уже посещен
		work4node(index)            //                                            обработали текущий узел(с ID колоды)
	}

	for 0 != len(bucket4Curr) {

		for _, v := range bucket4Curr { //                                        перебираем колоды в корзине текущего K-слоя

			for n, _ := range v { //                                              перебираем все узлы текущей колоды

				if gm.NodeExists(n) {
					gm.SetFlag(n, setColor) //                                    этот узел(с ID колоды) уже посещен
					work4node(n)            //                                    обработали текущий узел(с ID колоды)
				}
			}

			for n, _ := range v { //                                              перебираем все узлы текущей колоды

				var temp = map[uint64]Nothing{}

				for c, _ := range gm[n].conID { //                                перебираем все связи текущего узла

					if gm[c].flagColor == 0 { //                                  если ведут на непосещенный узел...

						temp[c] = Null
					}
				}

				if 0 != len(temp) { //                                            не пустую ли N-ную колоду предстоит создать

					bucket4Next[n] = map[uint64]Nothing{}
					bucket4Next[n] = temp
				}
			}
			bucket4Curr = bucket4Next                     //                      затирает (заменяет)
			bucket4Next = map[uint64]map[uint64]Nothing{} //                      опустошает
		}
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) ChainReaction(index uint64, work4node func(uint64) error) {

	gm.ResetAllFlags() //                                                         все узлы непосещены

	if gm.NodeExists(index) {

		gm.SetFlag(index, 1) //                                                   этот узел уже посещен
		work4node(index)     //                                                   обработали узел вхождения

	}

	deck4Curr := map[uint64]Nothing{} //                                          колода карт текущего уровня обработки - над каждым узлом совершается элементарное действие обработчика
	deck4Next := map[uint64]Nothing{} //                                          колода карт для следующего уровня обработки

	deck4Curr = gm[index].conID //                                                единственный элемент в колоде на старте - соединения единственного входного узла

	for 0 != len(deck4Curr) {

		var temp = map[uint64]Nothing{}

		for n, _ := range deck4Curr { //                                          перебираем все узлы входной колоды

			if gm.NodeExists(n) {

				gm.SetFlag(n, 1) //                                               этот узел(с ID из колоды вхождения) уже посещен
				work4node(n)     //                                               обработали текущий узел(с ID из колоды вхождения)
			}
		}

		for n, _ := range deck4Curr { //                                          перебираем все узлы входной колоды

			for c, _ := range gm[n].conID { //                                    перебираем все связи текущего узла

				if gm[c].flagColor == 0 { //                                      если ведут на непосещенный узел...

					temp[c] = Null //                                             ...добавляем во временный буфер c ID нашей связи
				}
			}
		}

		if 0 != len(temp) { //                                                    не пустую ли колоду предстоит создать
			deck4Next = map[uint64]Nothing{}
			deck4Next = temp
		}
		deck4Curr = deck4Next            //                                       затирает (заменяет)
		deck4Next = map[uint64]Nothing{} //                                       опустошает
	}
	return
}

// -------------------------------------------------------------------------------
func (gm orGraphsNetwork) CheckOnAcyclicity(index uint64) (bool, map[uint64]map[uint64]Nothing) {

	//...................................................................\/.static
	var cn uint64

	flag := false
	cm := map[uint64]Nothing{}
	outV := map[uint64]map[uint64]Nothing{}

	MC := func(src map[uint64]Nothing) map[uint64]Nothing {

		targ := make(map[uint64]Nothing)

		for k, v := range src {
			targ[k] = v
		}
		return targ
	}
	//.................................................................../\.static
	var REC func(uint64)

	REC = func(n uint64) { // >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> R

		gm.SetFlag(n, 1) //                                                        посетили узел
		cm[n] = Null

		for i, _ := range gm[n].conID { //                                         перебираем все карты связей, идекс и есть их ID

			if (gm[i].flagColor == 0) && gm.NodeExists(i) { //                     если узел непосещен, то...

				REC(i)
			}

			if gm[i].flagColor == 1 { //                                           если узел уже посещен, то...

				flag = true //                                                     отметим это
				outV[cn] = MC(cm)
				cn = cn + 1 //                                                     счетчик циклов
			}
		}

		delete(cm, n)
		gm.SetFlag(n, 2) //                                                        ...и посетили, и в цикле учли

		return
	} // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< R
	gm.ResetAllFlags() //                                                          все узлы непосещены

	if gm.NodeExists(index) {

		REC(index)
	}
	return flag, outV //                                                            "false, nil" if cycle no exist
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ INFO +
func (gm orGraphsNetwork) PrintNodeInfo(index uint64) {

	fmt.Println("Info about node with ID: ", index)
	fmt.Println("Flag is:                 ", gm[index].flagColor)
	fmt.Println("Data type is:            ", reflect.TypeOf(gm[index].data))
	fmt.Println("|")
	fmt.Println("con", gm[index].conID)
	for i, _ := range gm[index].conID { //                                         перебираем все соединения

		fmt.Println("+-->", i)

	}
	return
}

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ INFO +
func (gm orGraphsNetwork) PrintGraphNetInfo() {

	tc := map[string]Nothing{}

	connum := 0
	badcon := 0
	edge := 0

	dropp := map[uint64]Nothing{}

	for i, _ := range gm { //                                                     перебираем все узлы

		tc[fmt.Sprintf("%T", gm[i].data)] = Null

		if len(gm[i].conID) > 0 {
			for k, _ := range gm[i].conID { //                                    перебираем все их, узлов, связи

				connum = connum + 1 //                                            считаем связи

				_, ok := gm[k]

				if !ok {
					badcon = badcon + 1 //                                        считаем пустые связи
				}

				if len(gm[k].conID) > 0 { //                                      если множество связей в указываемом узле не пустое
					_, yes := gm[k].conID[i] //                                   ...и есть дуга на этот узел...

					_, ouch := dropp[k]
					_, stop := dropp[i]

					if yes && !ouch && !stop {
						edge = edge + 1 //
						dropp[k] = Null
						dropp[i] = Null
					}
				}
			}
		}
	}

	fmt.Println("info about", reflect.TypeOf(gm))
	fmt.Println("with size of memory ", unsafe.Sizeof(gm), "of bytes")
	fmt.Println("nodes number:       ", len(gm))
	fmt.Println("with data types:    ")

	for it, _ := range tc {
		fmt.Println(it)
	}

	fmt.Println("connection(arc) number:  ", connum)
	fmt.Println("two way connection(edge) number:  ", edge)
	fmt.Println("...of them - empty: ", badcon)

	return
}
