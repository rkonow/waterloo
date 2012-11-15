package invClustering

import (
	"container/heap"
	"fmt"
	"io/ioutil"
	"os"
	//"io"
	// "runtime"
	"sort"
	//"strconv"
	"bytes"
	"encoding/gob"
	"math/rand"
	"strings"
	//"stemmer"
)

type Posting struct {
	max_docid int
	min_docid int
	max_freq  int
	min_freq  int
	Content   map[int]int
}
type DocPosting struct {
	Terms map[string]int
}

type InvList struct {
	Postings   map[string]*Posting
	Dictionary map[string]int
	DocPost    map[int]*DocPosting
	TermMap    []string
	Doclen     []int
	Numdocs    int
	Numterms   int
}

func (h *InvList) GenerateQueries(start int, end int, n int, thresh int) {

	for j := start; j < end; j++ {
		aux := fmt.Sprintf("queries.%d", j)
		fo, err := os.Create(aux)
		if err != nil {
			panic("Error Opening the file")
		}
		defer fo.Close()

		i := 0
		s := ""
		for i < n {
			for k := 0; k < j; {
				p := rand.Intn(len(h.TermMap) - 1)
				if h.Dictionary[h.TermMap[p]] >= thresh {
					k++
					s = fmt.Sprintf("%s ", h.TermMap[p])
					fo.Write([]byte(s))
				}
			}
			i++
			s = fmt.Sprintf("\n")
			fo.Write([]byte(s))
		}

	}
}

func newInvlist() *InvList {
	i := new(InvList)
	i.Postings = make(map[string]*Posting, 0)
	i.Dictionary = make(map[string]int, 0)
	i.DocPost = make(map[int]*DocPosting, 0)
	i.Doclen = make([]int, 0)
	return i

}

func (h *InvList) GetDocsFromTerm(term string) int {
	total := 0
	if h.Postings[term] != nil {
	total +=  len(h.Postings[term].Content)
	} else {
	total = 100000000
	}
	return total
}

func (h *InvList) TopKQuery(queries []string, topn int) int {
	pq := make(PriorityQueue, 0, len(h.TermMap))
	nitem := 0
	nsize := 0
	result := make(map[int]int, 0)
	//fmt.Println(queries)
	for q := range queries {
		if h.Postings[queries[q]] != nil {
			nsize += len(h.Postings[queries[q]].Content)
			//	fmt.Println(len(h.Postings[queries[q]].Content))
			for s := range h.Postings[queries[q]].Content {
				//			fmt.Println(h.Postings[queries[q]].Content)
				result[s] += h.Postings[queries[q]].Content[s]
			}
		}
	}

	for r := range result {
		item := &Item{
			value:    r,
			priority: result[r],
		}
		nitem++
		heap.Push(&pq, item)
	}

	for i := 0; i < topn; i++ {
		if nitem <= i {
			break
		}
		item := heap.Pop(&pq).(*Item)
		fmt.Sprintf("%d", item.value)

	}
	//	fmt.Println(nsize)
	return nsize
}
func LoadFiles(file string) *InvList {
	data, _ := ioutil.ReadFile(file)

	docs := strings.Split(string(data), "\n")
	fmt.Println("Number of docs:" + string(len(docs)))
	docs = docs[:len(docs)-1]
	master := newInvlist()
	for i := range docs {
		//fmt.Println(docs[i])
		pctg := (float64(i) / float64(len(docs))) * 100
		// if i%100 == 0 {
		// 	fmt.Println(pctg)
		// }
		if (int(pctg))%10 == 0 && int(pctg) > 0 {
			fmt.Println(pctg)
		}
		if docs[i] != "" {
			master.readFile(docs[i], i)
		}
	}
	master.Numdocs = len(docs)
	master.Numterms = len(master.Postings)
	master.TermMap = make([]string, 0)
	for i := range master.Dictionary {
		master.TermMap = append(master.TermMap, i)
	}
	sort.Strings(master.TermMap)
	return master
}

func (h *InvList) readFile(file string, doc_id int) {
	dataDoc, _ := ioutil.ReadFile(file)
	Terms := strings.Split(string(dataDoc), " ")
	cnt := 0
	for i := range Terms {
		if len(Terms[i]) > 4 {
			Terms[i] = strings.ToLower(Terms[i])
			Terms[i] = strings.Trim(Terms[i], " /@#$*+-.,\n()?!{}|;:@\\'\"")
			Terms[i] = strings.Replace(Terms[i], ".", "", -1)
			Terms[i] = strings.Replace(Terms[i], ",", "", -1)
			Terms[i] = strings.Replace(Terms[i], "(", "", -1)
			Terms[i] = strings.Replace(Terms[i], ")", "", -1)
			Terms[i] = strings.Replace(Terms[i], "-", "", -1)
			Terms[i] = strings.Replace(Terms[i], ";", "", -1)
			Terms[i] = strings.Replace(Terms[i], ":", "", -1)
			Terms[i] = strings.Replace(Terms[i], " ", "", -1)
			Terms[i] = strings.Replace(Terms[i], "?", "", -1)
			Terms[i] = strings.Replace(Terms[i], "!", "", -1)
			Terms[i] = strings.Replace(Terms[i], "@", "", -1)
			Terms[i] = strings.Replace(Terms[i], "'", "", -1)
			Terms[i] = strings.Replace(Terms[i], "&", "", -1)
			Terms[i] = strings.Replace(Terms[i], "_", "", -1)
			Terms[i] = strings.Replace(Terms[i], "{", "", -1)
			Terms[i] = strings.Replace(Terms[i], "}", "", -1)
			Terms[i] = strings.Replace(Terms[i], "/", "", -1)
			Terms[i] = strings.Replace(Terms[i], "\\", "", -1)
			Terms[i] = strings.Replace(Terms[i], "\"", "", -1)
			Terms[i] = strings.Replace(Terms[i], "+", "", -1)
			Terms[i] = strings.Replace(Terms[i], "*", "", -1)
			Terms[i] = strings.Replace(Terms[i], "#", "", -1)
			Terms[i] = strings.Replace(Terms[i], "$", "", -1)
			Terms[i] = strings.Replace(Terms[i], "\n", "", -1)

			if len(Terms[i]) > 4 {
				Terms[i] = string(Stem([]byte(Terms[i])))
				h.processTerm(string(Terms[i]), doc_id)
				cnt++
			}
		}
	}
	h.Doclen = append(h.Doclen, cnt)
}

func (h *InvList) TermsDoc(doc_id int) map[string]int {
	return h.DocPost[doc_id].Terms
}

func (h *InvList) FreqInDoc(term string, doc_id int) int {
	return h.DocPost[doc_id].Terms[term]
}

func (h *InvList) TermId(term string) int {
	return sort.SearchStrings(h.TermMap, term)
}

func (h *InvList) NewFromCluster(documents []int, thd int) []*Cluster {
	//i := new(InvList)
	c := make([]*Cluster, 0)
	c = append(c, new(Cluster))
	// for i:= 0 ; i < len(documents)/thd+1;i++ {
	// 	c[i] = new(Cluster)
	// }
	count := 0
	for i := range documents {

		//fmt.Println("document = ",documents[i])
		for t := range h.DocPost[documents[i]].Terms {
			//doclen += h.DocPost[documents[i]].Terms[t]
			//new_h.processTerm(t, documents[i])
			//fmt.Println("index = ", index)
			// if h.DocPost[documents[i]].Terms[t] < 2 {
			// 	// add this document elsewhere
			// 	continue
			// }
			if count%thd == 0 {
				c = append(c, new(Cluster))
			}
			count++
			index := count / thd
			c[index].add(t, documents[i], h.DocPost[documents[i]].Terms[t])
		}
		//new_h.Doclen = append(new_h.Doclen, doclen)
	}
	return c
}

func (h *InvList) processTerm(term string, doc_id int) {
	h.Dictionary[term]++
	if h.DocPost[doc_id] == nil {
		h.DocPost[doc_id] = new(DocPosting)
		h.DocPost[doc_id].Terms = make(map[string]int, 0)
		h.DocPost[doc_id].Terms[term] = 1
	} else {
		h.DocPost[doc_id].Terms[term]++
	}

	if h.Postings[term] == nil {
		h.Postings[term] = new(Posting)
		h.Postings[term].Content = make(map[int]int, 0)
		h.Postings[term].Content[doc_id] = 1
	} else {
		h.Postings[term].Content[doc_id]++
	}
}

func Load(file string) *InvList {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Error Opening the file")
	}
	var input bytes.Buffer
	input.Write([]byte(data))
	dec := gob.NewDecoder(&input)
	var inv *InvList
	err = dec.Decode(&inv)
	return inv
}

func (h *InvList) Write(file string) {
	fo, err := os.Create(file + ".serial")
	if err != nil {
		panic("Error Opening the file")
	}
	defer fo.Close()
	var output bytes.Buffer
	enc := gob.NewEncoder(&output)
	err = enc.Encode(h)
	fo.Write(output.Bytes())
}
