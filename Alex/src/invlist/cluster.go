package invClustering

import (
	"fmt"
	// //"invlist"
	"container/heap"
	//"fmt"
	//"math/rand"
)

type Cluster struct {
	n         int
	Terms     map[string]*Posting
	Documents []int
	id        int
	max_docid int
	min_docid int
	max_freq  int
}

type StatValues struct {
	max_freq  int
	max_docid int
	min_docid int
}
type TermCluster struct {
	TermMap    map[string][]*Cluster
	TermValues map[string]*StatValues
}

func (c *Cluster) add(term string, doc_id int, freq int) {
	if len(c.Terms) == 0 {
		c.Terms = make(map[string]*Posting, 0)
	}
	if c.Terms[term] == nil {
		c.Terms[term] = new(Posting)
		c.Terms[term].Content = make(map[int]int, 0)
		c.Terms[term].Content[doc_id] = freq
	} else {
		c.Terms[term].Content[doc_id] = freq
	}

	if len(c.Documents) == 0 {
		c.Documents = make([]int, 0)
		c.Documents = append(c.Documents, doc_id)
	} else {
		c.Documents = append(c.Documents, doc_id)
	}

	if c.min_docid >= doc_id {
		c.min_docid = doc_id
	}
	if c.max_docid <= doc_id {
		c.max_docid = doc_id
	}
	if c.max_freq <= freq {
		c.max_freq = freq
	}
	if c.Terms[term].max_docid <= doc_id {
		c.Terms[term].max_docid = doc_id
	}
	if c.Terms[term].min_docid >= doc_id {
		c.Terms[term].min_docid = doc_id
	}
	if c.Terms[term].max_freq <= freq {
		c.Terms[term].max_freq = freq
	}
	if c.Terms[term].min_freq >= freq {
		c.Terms[term].min_freq = freq
	}

}

func (t *TermCluster) add(term string, cluster *Cluster) {
	// fmt.Printf("Adding term %s to cluster %d \n", term, cluster.id)
	// fmt.Printf("%d \n",cluster.max_freq)
	// fmt.Printf("%d \n",cluster.max_docid)
	// fmt.Printf("%d \n",cluster.min_docid)

	if t.TermValues == nil {
		t.TermValues = make(map[string]*StatValues)
	}
	if t.TermValues[term] == nil {
		t.TermValues[term] = new(StatValues)
		t.TermValues[term].min_docid = 999999999
		t.TermValues[term].max_freq = -1
		t.TermValues[term].max_docid = -1
	}

	if t.TermValues[term].max_freq < cluster.max_freq {
		t.TermValues[term].max_freq = cluster.max_freq
	}
	if t.TermValues[term].max_docid < cluster.max_docid {
		t.TermValues[term].max_docid = cluster.max_docid
	}
	if t.TermValues[term].min_docid > cluster.min_docid {
		t.TermValues[term].min_docid = cluster.min_docid
	}
	t.TermMap[term] = append(t.TermMap[term], cluster)
}

func (t *TermCluster) Stats() {
	avg := float64(0)
	avg2 := float64(0)
	for i := range t.TermMap {
		avg += float64(len(t.TermMap[i]))
		for j := range t.TermMap[i] {
			avg2 += float64(len(t.TermMap[i][j].Documents))
		}
	}
	fmt.Println(avg2 / float64(len(t.TermMap)))
	fmt.Println(avg / float64(len(t.TermMap)))
	fmt.Println(len(t.TermMap))
}
func (t *TermCluster) TopKQuery(terms []string, topn int, thresh int) (int, int) {
	pq := make(PriorityQueue, 0, len(t.TermMap))
	nitems := 0
	nclusters := 0
	clusters_size := 0
	for i := range terms {
		// fmt.Printf("Term[%d] = %s \n", i, terms[i])
		for j := range t.TermMap[string(terms[i])] {
			if t.TermValues[string(terms[i])].max_freq < thresh {
				continue
			}
			// fmt.Printf("Term -> %s  Has a pointer to cluster: %d \n in Documents: ", terms[i], t.TermMap[terms[i]][j].id)
			nclusters++
			clusters_size += len(t.TermMap[terms[i]][j].Terms[terms[i]].Content)
			//len(t.TermMap[terms[i]][j].Documents)  //)
			//fmt.Println(len(t.TermMap[terms[i]][j].Documents))
			for k := range t.TermMap[terms[i]][j].Terms[terms[i]].Content {
				freq := t.TermMap[terms[i]][j].Terms[terms[i]].Content[k]
				c_id := t.TermMap[terms[i]][j].id
				item := &Item{
					value:    c_id,
					priority: freq,
				}
				heap.Push(&pq, item)
				nitems++
				// fmt.Printf("Document : %d with freq -> %d \n", k, t.TermMap[terms[i]][j].Terms[terms[i]].Content[k])
			}
		}
	}
	// results_cluster := make([]int, 0)
	// results_freqs := make([]int, 0)
	results_map := make(map[int]int)
	for i := 0; i < topn; i++ {
		if nitems <= i {
			break
		}
		item := heap.Pop(&pq).(*Item)
		if results_map[item.value] == 0 {
			results_map[item.value] = item.priority
		} else {
			results_map[item.value] += item.priority
		}
	}
	// fmt.Printf("Result MAP: \n ------------------ \n")
	// for i := range results_map {
	// 	fmt.Printf("Cluster: %d -> freq: %d \n", i, results_map[i])
	// 	results_cluster = append(results_cluster, i)
	// 	results_freqs = append(results_freqs, results_map[i])
	// }
	return nclusters, clusters_size
}

func GenerateTermClusters(c []*Cluster) *TermCluster {
	t := new(TermCluster)
	t.TermMap = make(map[string][]*Cluster, 0)
	for i := range c {
		for j := range c[i].Terms {
			// fmt.Printf("Cluster -> %d  Has term %s \n", i, j)
			t.add(j, c[i])
		}
	}
	return t
}

func GenerateClusters(list *InvList, n int) []*Cluster {

	clusters := make([]*Cluster, 0)
	t := new(Cluster)
	t.id = 0
	//t.Terms = make(map[string]*DocCluster,0)
	clusters = append(clusters, t)
	max_capacity := n
	c_id := 0
	terms := make([]string, 0)
	docs := make([]int, 0)
	freqs := make([]int, 0)

	for d := range list.DocPost {
		for i := range list.DocPost[d].Terms {
			terms = append(terms, i)
			docs = append(docs, d)
			freqs = append(freqs, list.Postings[i].Content[d])
		}
	}
	terms, docs, freqs = Shuffle(terms, docs, freqs)
	x := 0
	for x < len(docs) {
		i := terms[x]
		d := docs[x]
		f := freqs[x]
		if len(clusters[c_id].Documents) < max_capacity {
			//fmt.Printf("%s \t %d \t %d \n",i,d,f)
			clusters[c_id].add(i, d, f)

			x++
			//	fmt.Printf("Adding term %s from doc %d to Cluster %d \n", i, d, c_id)
		} else {
			tp := new(Cluster)
			tp.id = c_id
			c_id++
			clusters = append(clusters, tp)
			clusters[c_id].add(i, d, f)
			x++
		}
	}
	return clusters
}
