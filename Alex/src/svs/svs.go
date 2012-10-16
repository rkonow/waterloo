package svs

import (
	"math"
//	"fmt"
	// "sort"
	"time"
)

type SetList struct {
	Sets []*InvertedSet
}

func (h *SetList) Len() int {
	return len(h.Sets)
}
func (h *SetList) Less(i, j int) bool {
	return len(h.Sets[i].Content) < len(h.Sets[j].Content)
}
func (h *SetList) Swap(i, j int) {
	h.Sets[i], h.Sets[j] = h.Sets[j], h.Sets[i]
}

type InvertedSet struct {
	max_docid   int
	min_docid   int
	Content     []int
	Frequencies []int
}

// Methods required by sort.Interface.
func (h *InvertedSet) Len() int {
	return len(h.Content)
}
func (h *InvertedSet) Less(i, j int) bool {
	return h.Content[i] < h.Content[j]
}
func (h *InvertedSet) Swap(i, j int) {
	h.Content[i], h.Content[j] = h.Content[j], h.Content[i]
	h.Frequencies[i], h.Frequencies[j] = h.Frequencies[j], h.Frequencies[i]
}

func (h *InvertedSet) AddDocument(doc int, freq int) {
	h.Content = append(h.Content, doc)
	h.Frequencies = append(h.Frequencies, freq)
	// fmt.Println("doc is",doc)
	// fmt.Println("doc added:",h.Content[len(h.Content)-1])
	if h.max_docid < doc {
		h.max_docid = doc
	}
	if h.min_docid > doc {
		h.min_docid = doc
	}
}

func NewInvertedSet() *InvertedSet {
	result := new(InvertedSet)
	result.max_docid = 0
	result.min_docid = 1<<31 - 1
	result.Content = make([]int, 0)
	result.Frequencies = make([]int, 0)
	return result
}

func (h *InvertedSet) GetMaxDoc() int {
	return h.max_docid
}
func (h *InvertedSet) GetMinDoc() int {
	return h.min_docid
}
func (h *InvertedSet) GetSize() int {
	return len(h.Content) * 8
}
func (h *InvertedSet) GetLength() int {
	return len(h.Content)
}

func exponentialSearch(a []int, value int,low int, high int) int {
	i:= low
	j:=2
	start := low
	for i<high && a[i] < value {
		//fmt.Println("i = ",i)
		start = i
		i+=j
		j*=2;
	}

	 // fmt.Println("*i = ",i)
	 // fmt.Println("*i-1 = ",start)
	if (i >= len(a)) {
		i = len(a)
		if (a[i-1] < value)  { 
			return -1
		}
	}
	// for k:=start;k<i;k++ {
	// 	if (a[k] == value) {
	// 		return k
	// 	}
	// }
	found := binarySearch(a,value,start,i)
	return found
}

func binarySearch(a []int, value int, low int, high int) int {
	// fmt.Println("entered with low = ", low , " high = " , high)
	old_high := high 
    mid := low+1
    for low <= high {
        mid = (low + high) / 2
       // fmt.Println("mid = ",mid)
        if (mid >= old_high) {
        //	fmt.Println("returning :", old_high)
        	return old_high
        }
        if a[mid] > value {
            high = mid - 1
        } else if a[mid] < value {
            low = mid + 1
        } else {
    	//    fmt.Println("returning found:", mid)
            return mid
        }
    }
    //fmt.Println("returning :", mid)
    return mid
}

func Intersect(terms []InvertedSet, y0 int, y1 int, freq int, result map[int]int) (float64, map[int]int) {
	//fmt.Println(len(terms))
	//total_terms := len(terms[0].Content)
	//fmt.Println(total_terms)
	//total_out := 0
	l := make([]int,len(terms))
	for i := 0;i<len(terms);i++ {
		l[i] = 0
	}

	t0 := time.Now()
	new_max := int(math.Min(float64(terms[0].GetMaxDoc()), float64(terms[1].GetMaxDoc())))
	new_min := int(math.Max(float64(terms[0].GetMinDoc()), float64(terms[1].GetMinDoc())))
	//for i:=0;i<len(terms.Sets);i++ {
	// 	fmt.Println("i=",i," len =",len(terms.Sets[i].Content))
	// }	
	for i:= 2;i<len(terms);i++ {
		new_max = int(math.Min(float64(new_max),float64(terms[i].GetMaxDoc())))
		new_min = int(math.Max(float64(new_min),float64(terms[i].GetMinDoc())))
	}

	// for i:=0 ; i < len(terms[0].Content);i++ {
	// 	result.Content[i] = terms[0].Content[i]
	// }
	//result.Frequencies = terms[0].Frequencies
	//deleted := 0

	cont := true
	for s := 1; s < len(terms); s++ {
		if cont == false { 
			break
		}
		for _,key := range terms[0].Content {
			// fmt.Println("pos",pos,"key=",key)
			if len(result) == 0 {
				cont = false
				break
			}
			_, ok := result[key]
			if !ok {
				continue
			}
			if key < new_min {
				delete(result, key)
				continue
			}
			if key > new_max {
				delete(result, key)
				continue
			}
			if key < y0 {
				delete(result, key)
				continue
			}
			if key > y1 {
				delete(result, key)
				continue
			}
			//found := sort.Search(len(terms[s].Content), func(found int) bool { return terms[s].Content[found] >= key })
			//found := sort.SearchInts(terms[s].Content, key)
			found := exponentialSearch(terms[s].Content,key,l[s],len(terms[s].Content))
	//		fmt.Println("found=",found)
			if found != -1 { 
				if found < len(terms[s].Content) && terms[s].Content[found] == key {
					l[s] = found
					if (terms[s].Frequencies[found] < freq) {
						delete(result, key)
						continue
					}
				} else {
					l[s] = found
					delete(result, key)
					continue
				}
			} else {
				delete(result, key)
				continue
			}
		}
	}
	t1 := time.Now()
	t_final := float64(t1.Sub(t0).Seconds())
	return t_final, result
}




type InvertedSetFreq struct {
	max_docid   int
	min_docid   int
	Content     []int
	Frequencies []int
}

// Methods required by sort.Interface.
func (h *InvertedSetFreq) Len() int {
	return len(h.Content)
}
func (h *InvertedSetFreq) Less(i, j int) bool {
	return h.Frequencies[i] < h.Frequencies[j]
}
func (h *InvertedSetFreq) Swap(i, j int) {
	h.Content[i], h.Content[j] = h.Content[j], h.Content[i]
	h.Frequencies[i], h.Frequencies[j] = h.Frequencies[j], h.Frequencies[i]
}

func (h *InvertedSetFreq) AddDocument(doc int, freq int) {
	h.Content = append(h.Content, doc)
	h.Frequencies = append(h.Frequencies, freq)
	// fmt.Println("doc is",doc)
	// fmt.Println("doc added:",h.Content[len(h.Content)-1])
	if h.max_docid < doc {
		h.max_docid = doc
	}
	if h.min_docid > doc {
		h.min_docid = doc
	}
}

func NewInvertedSetFreq() *InvertedSetFreq {
	result := new(InvertedSetFreq)
	result.max_docid = 0
	result.min_docid = 1<<31 - 1
	result.Content = make([]int, 0)
	result.Frequencies = make([]int, 0)
	return result
}

