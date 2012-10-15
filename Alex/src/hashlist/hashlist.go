package hashlist

import (
	"math"

//	"fmt"
)

type InvertedHashList struct {
	Hashings []*InvertedHash
}

func (s *InvertedHashList) Len() int {
	return len(s.Hashings)
}
func (s *InvertedHashList) Less(i, j int) bool {
	return len(s.Hashings[i].Content) < len(s.Hashings[j].Content)
}
func (s *InvertedHashList) Swap(i, j int) {
	s.Hashings[i], s.Hashings[j] = s.Hashings[j], s.Hashings[i]
}

type InvertedHash struct {
	max_docid int
	min_docid int
	Content   map[int]int
}

func NewInvertedHash() *InvertedHash {
	result := new(InvertedHash)
	result.max_docid = 0
	result.min_docid = 1<<31 - 1
	result.Content = make(map[int]int, 0)
	return result
}
func (h *InvertedHash) AddDocument(doc int, freq int) {
	h.Content[doc] = freq
	if h.max_docid < doc {
		h.max_docid = doc
	}
	if h.min_docid > doc {
		h.min_docid = doc
	}
}
func (h *InvertedHash) GetMaxDoc() int {
	return h.max_docid
}
func (h *InvertedHash) GetMinDoc() int {
	return h.min_docid
}
func (h *InvertedHash) GetSize() int {
	return len(h.Content) * 8
}
func (h *InvertedHash) GetLength() int {
	return len(h.Content)
}

func IntersectHash(a *InvertedHash, b *InvertedHash, y0 int, y1 int, freq int) *InvertedHash {
	new_max := int(math.Min(float64(a.GetMaxDoc()), float64(b.GetMaxDoc())))
	new_min := int(math.Max(float64(a.GetMinDoc()), float64(b.GetMinDoc())))
	result := NewInvertedHash()
	for key := range a.Content {
		// fmt.Println("a.key = ",key)
		// fmt.Println("b.value = ",b.Content[key])
		if key < new_min {
			continue
		}
		if key > new_max {
			continue
		}
		if key < y0 {
			continue
		}
		if key > y1 {
			continue
		}
		_,ok := b.Content[key]
		if ok == false {
			continue
		} else {
			if b.Content[key] < freq {
				continue
			} else {
				result.AddDocument(key, a.Content[key]+b.Content[key])
			}
		}
	}
	return result
}

func Intersect(terms *InvertedHashList, y0 int, y1 int, freq int) *InvertedHash {
	new_max := int(math.Min(float64(terms.Hashings[0].GetMaxDoc()), float64(terms.Hashings[1].GetMaxDoc())))
	new_min := int(math.Max(float64(terms.Hashings[0].GetMinDoc()), float64(terms.Hashings[1].GetMinDoc())))

	for i := 2; i < len(terms.Hashings); i++ {
		new_max = int(math.Min(float64(new_max), float64(terms.Hashings[i].GetMaxDoc())))
		new_min = int(math.Max(float64(new_min), float64(terms.Hashings[i].GetMinDoc())))
	}

	if new_min > new_max {
		return nil
	}

	res := IntersectHash(terms.Hashings[0], terms.Hashings[1], y0, y1, freq)
	//fmt.Println("res size=",len(res.Content))
	for i := 2; i < len(terms.Hashings); i++ {
		res = IntersectHash(res, terms.Hashings[i], y0, y1, freq)
		//	fmt.Println("res size=",len(res.Content))
	}
	//fmt.Println("res FINAL=",len(res.Content))
	return res
}
