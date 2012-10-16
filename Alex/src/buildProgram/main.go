package main

import (
	"fmt"
	"hashlist"
	"invlist"
	"math"
	"math/rand"
	"os"
	"sort"
	"svs"
	"time"
	"treap"

	//	"flame"
	//	"unsafe"
	//	"reflect"
	//"time"
	//"io/ioutil"
	//"math"
	//"strings"

//"math/rand"

//	"gomatrix"
)

func StringLess(p, q interface{}) bool {
	return p.(string) < q.(string)
}

func IntLess(p, q interface{}) bool {
	return p.(int) < q.(int)
}

func main() {

	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %s <input>\n", args[0])
		return
	}

	//t := invClustering.LoadFiles(args[1])
	//t.Write("test_large.invlist")
	t := invClustering.Load("test_large.invlist.serial")
	tr := make([]*treap.Tree, 0)
	hs := make([]*hashlist.InvertedHash, 0)
	ss := make([]*svs.InvertedSet, 0)
	term := 0
	count := 0
	total_height := float64(0)
	total_balanced := float64(0)
	total_len := float64(0)

	total_len_f_real := float64(0)
	total_len_f_delta := float64(0)
	total_len_f_delta_nozero := float64(0)
	total_len_k_real := float64(0)
	total_len_k_delta := float64(0)
	size_freq_stupid :=0
	count2:=0 


	for i := range t.Postings {

		//fmt.Println("Adding term:",i)
		size := len(t.Postings[i].Content)
		if size > 10   {
			tr = append(tr, treap.NewTree(IntLess))
			hs = append(hs, hashlist.NewInvertedHash())
			ss = append(ss, svs.NewInvertedSet())
			for k := range t.Postings[i].Content {
				tr[term].Insert(k, t.Postings[i].Content[k])
				hs[term].AddDocument(k, t.Postings[i].Content[k])
				ss[term].AddDocument(k, t.Postings[i].Content[k])

			}

			// for k := 0; k < 1000; k++ {
			// 	randnum := rand.Intn(5000)
			// 	randfreq := rand.Intn(100)
			// 	tr[term].Insert(randnum, randfreq)
			// 	hs[term].AddDocument(randnum, randfreq)
			// 	ss[term].AddDocument(randnum, randfreq)
			// }
			sort.Sort(ss[term])
			balanced := math.Ceil((math.Log2(float64(len(t.Postings[i].Content)))))
			height := tr[term].GetHeightTree(tr[term].GetRoot())
			total_height += float64(height)
			total_balanced += balanced
			total_len += float64(size)
			//tr[term].PrintPriorities(tr[term].GetRoot(),0,acum)
			tr[term].CompressPriority(tr[term].GetRoot())
			tr[term].CompressKeys(tr[term].GetRoot())
			acum := make([]int,0)
	//		tr[term].PrintDeltas(tr[term].GetRoot(),0,acum)
			for k:=0;k< len(ss[term].Content);k++ {
				acum = append(acum,tr[term].GetKeyPriorityDelta(tr[term].GetRoot(),ss[term].Content[k]))
			}
			acum_real := make([]int,0)
	//		tr[term].PrintDeltas(tr[term].GetRoot(),0,acum)
			for k:=0;k< len(ss[term].Content);k++ {
				acum_real = append(acum_real,tr[term].GetKeyPriority(tr[term].GetRoot(),ss[term].Content[k]))
			}
			acumkey := make([]int,0)
			for k:=0;k< len(ss[term].Content);k++ {
				acumkey = append(acumkey,tr[term].GetKeyDelta(tr[term].GetRoot(),ss[term].Content[k]))
			}

			size_freq_real := 0
			size_freq_delta := 0
			size_key_real := 0
			size_key_delta :=0

			for k:=0;k<len(ss[term].Content);k++ {
				f_real := acum_real[k]
				f_delta := acum[k]
				k_real := ss[term].Content[k]
				k_delta := acumkey[k]
				// fmt.Println("k_real:",k_real)
				// fmt.Println("k_delta",k_delta)
				size_freq_stupid += 32
				size_freq_real += int(math.Ceil((math.Log2(math.Abs(float64(f_real))+2))))
				size_freq_delta += int(math.Ceil(math.Log2(math.Abs(float64(f_delta))+2)))
				//fmt.Println(f_delta)
				if f_delta != 0 {
					//fmt.Println(f_delta)
					total_len_f_delta_nozero += float64(int(math.Ceil(math.Log2(math.Abs(float64(f_delta))+1))))
				}
				//fmt.Println(int(math.Ceil(math.Log2(math.Abs(float64(f_delta))+1))))
				size_key_real += int(math.Ceil(math.Log2(math.Abs(float64(k_real))+2)))
				size_key_delta += int(math.Ceil(math.Log2(math.Abs(float64(k_delta))+2)))
				count2++
			}

			total_len_f_real += float64(size_freq_real)
			total_len_f_delta += float64(size_freq_delta)
			total_len_k_real += float64(size_key_real)
			total_len_k_delta += float64(size_key_delta)
			count++
			term++
			
		}

	}
	// fmt.Println("AVG HEIGHT:",float64(total_balanced)/float64(count),"\t AVG TREAP:",float64(float64(total_height)/float64(count)))
	// fmt.Println("AVG LEN: ",float64(total_len)/float64(count))
	fmt.Println("f_stupid,f_real,f_delta,f_delta_no_zero,k_real,k_delta,count2,count,height,htreap,len")
	fmt.Println(size_freq_stupid/8,total_len_f_real/8,total_len_f_delta/8,total_len_f_delta_nozero/8,total_len_k_real/8,total_len_k_delta/8,count2,count,total_balanced/float64(count),total_height/float64(count),total_len/float64(count))
	//fmt.Println("f_stupid ", size_freq_stupid)
	// fmt.Println("f_real ", total_len_f_real)
	// fmt.Println("f_delta ", total_len_f_delta)
	// fmt.Println("f_delta_no_zero ", total_len_f_delta_nozero)


	// fmt.Println("k_real ", total_len_k_real)
	// fmt.Println("k_delta ", total_len_k_delta)
	

	// fmt.Println("AVG f_stupid ", float64(size_freq_stupid)/float64(count2))
	// fmt.Println("AVG f_real ", float64(total_len_f_real)/float64(count2))
	// fmt.Println("AVG f_delta ", float64(total_len_f_delta)/float64(count2))
	// fmt.Println("AVG f_delta_no_zero ", float64(total_len_f_delta_nozero)/float64(count2))


	// fmt.Println("AVG k_real ", float64(total_len_k_real)/float64(count2))
	// fmt.Println("AVG k_delta ", float64(total_len_k_delta)/float64(count2))

	// fmt.Println("AVG f_stupid/term ", float64(size_freq_stupid)/float64(count))
	// fmt.Println("AVG f_real/term ", float64(total_len_f_real)/float64(count))
	// fmt.Println("AVG f_delta/term ", float64(total_len_f_delta)/float64(count))
	// fmt.Println("AVG f_delta_no_zero/term ", float64(total_len_f_delta_nozero)/float64(count))


	// fmt.Println("AVG k_real/term", float64(total_len_k_real)/float64(count))
	// fmt.Println("AVG k_delta/term ", float64(total_len_k_delta)/float64(count))


	// fmt.Println("Naive keys/term", count2*32/count)
	// fmt.Println("Naive Freqs/term", count2*32/count)	
	// fmt.Println("Naive sum/term", count2*32*2/count)	

	for k := 2; k < 10; k++ {
		hashtime := float64(0)
		svstime := float64(0)
		treaptime := float64(0)
		
		for j := 0; j < 100000; j++ {
			termquerySet := new(svs.SetList)
			termquerySet.Sets = make([]*svs.InvertedSet, 0)

			termquerytreap := new(treap.TreapList)
			termquerytreap.Treap = make([]*treap.Tree, 0)

			termqueryHash := new(hashlist.InvertedHashList)
			termqueryHash.Hashings = make([]*hashlist.InvertedHash, 0)
			for i := 0; i < k; i++ {
				randnumber := rand.Intn(len(hs) - 1)

				termqueryHash.Hashings = append(termqueryHash.Hashings, hs[randnumber])

				termquerySet.Sets = append(termquerySet.Sets, ss[randnumber])

				termquerytreap.Treap = append(termquerytreap.Treap, tr[randnumber])
			}

			sort.Sort(termqueryHash)
			sort.Sort(termquerytreap)
			sort.Sort(termquerySet)
			min_freq := 2
			t_treap := time.Now()
			// r_treap := treap.Intersect(termquerytreap, 0, 1000000000, min_freq)
			treap.Intersect(termquerytreap, 0, 1000000000, min_freq)
			t_treap2 := time.Now()
			t_treap_f := float64(t_treap2.Sub(t_treap).Seconds())
			treaptime += t_treap_f

			t0 := time.Now()
			// r_hash := hashlist.Intersect(termqueryHash, 0, 1000000000, min_freq)
			hashlist.Intersect(termqueryHash, 0, 1000000000, min_freq)
			t1 := time.Now()
			t_hash := float64(t1.Sub(t0).Seconds())
			hashtime += t_hash

			testcopy := make([]svs.InvertedSet, k)
			for i := 0; i < k; i++ {
				testcopy[i] = *termquerySet.Sets[i]
			}

			result := make(map[int]int, len(testcopy[0].Content))
			for i := 0; i < len(testcopy[0].Content); i++ {
				result[testcopy[0].Content[i]] = testcopy[0].Frequencies[i]
			}
			t_svs_partial, _ := svs.Intersect(testcopy, 0, 1000000000, min_freq, result)

			//t_svs := float64(t3.Sub(t2).Seconds())
			svstime += t_svs_partial
			//fmt.Println("Result =",len(t.Content))
			//fmt.Println("Result2 =",len(t2.Content))
			// if (r != nil) {
			// if len(r_hash.Content) != len(test) {
			// 	fmt.Println("j = ", j)
			// 	fmt.Println(r_hash.Content)
			// 	fmt.Println("--------")
			// 	fmt.Println(test)
			// 	// 		// hashtime -= t_hash
			// 	// 		// svstime -= t_svs
			// 	// 		// treaptime -= t_treap_f
			// }
			//}
			// count_treap := 0
			// if (r_treap != nil ) {
			// 	count_treap = r_treap.GetCount()
			// }
			// count_hash := 0
			// if (r_hash != nil) {
			// 	count_hash = len(r_hash.Content)
			// }
			// fmt.Println("treap size:",count_treap)
			// fmt.Println("hash size:",count_hash)
		}

		fmt.Println(k,hashtime,svstime,treaptime)
		// fmt.Println(k," hash time:", hashtime)
		// fmt.Println(k," svs time:", svstime)
		// fmt.Println(k," treap time:", treaptime)
	}
	

	// fmt.Println("N = ", len(t.DocPost))
	// fmt.Println("M = ", len(t.TermMap))
	// doc_vector := make([][]float64, len(t.DocPost))
	// //max_freq := float64(0)
	// for i := range t.DocPost {
	// 	doc_vector[i] = make([]float64, len(t.TermMap))
	// 	for j := range t.DocPost[i].Terms {
	// 		//if (float64(t.DocPost[i].Terms[j]) > 1) {
	// 		//	fmt.Println("setting position ",sort.SearchStrings(t.TermMap, j),"to ", float64(t.DocPost[i].Terms[j]))
	// 		doc_vector[i][sort.SearchStrings(t.TermMap, j)] = float64(t.DocPost[i].Terms[j])
	// 		//}
	// 		//	fmt.Println("term =",j,"value = ", t.DocPost[i].Terms[j])
	// 	}
	// 	//fmt.Println("i=",i,doc_vector[i])
	// }
	//	f := invClustering.
	//f := flame.NewFlame()

	// juice := make([]*float32, len(doc_vector))
	// for i := range juice {
	// 	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&doc_vector[i])))
	// 	juice[i] = (*float32)(unsafe.Pointer(sliceHeader.Data))
	// }
	// sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&juice)))
	// flame.Flame_SetDataMatrix(f,(**float32)(unsafe.Pointer(sliceHeader.Data)),len(t.DocPost),len(t.TermMap),1)
	// flame.Flame_DefineSupports(f, 10, 100)
	// flame.Flame_LocalApproximation(f,10,100)
	// flame.Flame_MakeClusters(f,0.5)
	// ia := f.GetClusters().GetArray()
	// // for i := 0; i < f.GetClusters().GetSize(); i++ {
	// // 	fmt.Println("array[",i,"]=",ia[i])
	// // }
	// var ia2 []int
	// sliceHeader = (*reflect.SliceHeader)(unsafe.Pointer(&ia2))
	// sliceHeader.Len = f.GetClusters().GetSize()
	// sliceHeader.Cap = f.GetClusters().GetSize()
	// sliceHeader.Data = (uintptr)(unsafe.Pointer(ia))
	// for i := range ia2 {
	// 	fmt.Println("array[", i, "]=", ia2[i])
	// }
	//fmt.Println(doc_vector)
	//fmt.Println(t.TermMap)
	// f := invClustering.CreateFlame()
	// fmt.Println("Euclidean Distance")
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.Euclidean)
	// f.Write("clusters_euclidean2")

	// fmt.Println("Cosine")
	// f = invClustering.CreateFlame()
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.Cosine)
	// f.Write("clusters_cosine2")

	// fmt.Println("Cosine Distance")
	// f = invClustering.CreateFlame()
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.CosineDist)
	// f.Write("clusters_cosinedist2")

	// fmt.Println("Pearson")
	// f = invClustering.CreateFlame()
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.Pearson)
	// f.Write("clusters_pearson2")

	// fmt.Println("Pearson Distance")
	// f = invClustering.CreateFlame()
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.PearsonDist)
	// f.Write("clusters_pearsondist2")

	// fmt.Println("Dot Product")
	// f = invClustering.CreateFlame()
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.DotProduct)
	// f.Write("clusters_dot2")

	// fmt.Println("Dot Product Distance")
	// f = invClustering.CreateFlame()
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.DotProductDist)
	// f.Write("clusters_dotdist2")

	// fmt.Println("Covariance")
	// f = invClustering.CreateFlame()
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.Covariance)
	// f.Write("clusters_covariance2")

	// fmt.Println("Covariance Distance")
	// f = invClustering.CreateFlame()
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.DistCovariance)
	// f.Write("clusters_covariancedist2")

	// fmt.Println("Jaccard ")
	// f = invClustering.CreateFlame()
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.Jaccard)
	// f.Write("clusters_jaccard2")

	// fmt.Println("Jaccard Pearson ")
	// f = invClustering.CreateFlame()
	// f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.JaccardPearson)
	// f.Write("clusters_jaccard_pearson2")

	//clusters := []string{"clusters_euclidean", "clusters_cosine", "clusters_cosinedist", "clusters_pearson", "clusters_pearsondist", "clusters_dot", "clusters_dotdist", "clusters_covariance", "clusters_covariancedist", "clusters_jaccard", "clusters_jaccard_pearson"}
	// for c := range clusters {
	// 	f := invClustering.LoadClusters(clusters[c] + ".serial")
	// 	f.DefineSupports(1, 1)
	// 	f.LocalApproximation(1, 1)
	// 	f.MakeClusters(0.9)
	// 	reps := make(map[int]int)
	// 	avg := float64(0)
	// 	nonzero := 0
	// 	num_clusters := 0
	// 	num_clusters = len(f.Clusters)
	// 	max_cluster_size := 0
	// 	for i := 0; i < len(f.Clusters); i++ {
	// 		if len(f.Clusters[i].Array) != 0 {
	// 			if max_cluster_size < len(f.Clusters[i].Array) {
	// 				max_cluster_size = len(f.Clusters[i].Array)
	// 			}
	// 			avg += float64(len(f.Clusters[i].Array))
	// 			nonzero++
	// 		}
	// 		for j := 0; j < len(f.Clusters[i].Array); j++ {
	// 			reps[f.Clusters[i].Array[j]]++
	// 		}
	// 	}
	// 	avg = avg / float64(nonzero)
	// 	repet := 0
	// 	for i := range reps {
	// 		if reps[i] > 1 {
	// 			repet++
	// 		}
	// 	}
	// 	rep_per_cluster := float64(repet) / float64(num_clusters)
	// 	fmt.Println(clusters[c], "\t", num_clusters, "\t", avg, "\t", repet, "\t", rep_per_cluster, "\t", max_cluster_size)
	// }
	// t.GenerateQueries(2, 11, 1000, 1)

	// for i := 0; i < len(clusters); i++ {
	// 	f := invClustering.LoadClusters(clusters[i] + ".serial")
	// 	fmt.Println("-------------------------------")
	// 	fmt.Println(clusters[i])
	// 	f.DefineSupports(1, 1)
	// 	f.LocalApproximation(1, 1)
	// 	f.MakeClusters(0.5)
	// 	// //fmt.Println(f)
	// 	clusts := make([]*invClustering.Cluster, 0)
	// 	//fmt.Println(f.Clusters)
	// 	for i := 0; i < len(f.Clusters); i++ {
	// 		if len(f.Clusters[i].Array) != 0 {
	// 			if len(f.Clusters[i].Array) > 128 {
	// 				clusts_aux := t.NewFromCluster(f.Clusters[i].Array, int(math.Sqrt(float64(len(f.Clusters[i].Array))))+1)
	// 				for j := 0; j < len(clusts_aux); j++ {
	// 					clusts = append(clusts, clusts_aux[j])
	// 				}
	// 			} else {
	// 				clusts_aux := t.NewFromCluster(f.Clusters[i].Array, 128)
	// 				for j := 0; j < len(clusts_aux); j++ {
	// 					clusts = append(clusts, clusts_aux[j])
	// 				}
	// 			}

	// 		}
	// 	}
	// 	//fmt.Println(clusts)
	// 	termCluster := invClustering.GenerateTermClusters(clusts)
	// 	termCluster.Stats()
	// 	fmt.Println("# clusters = ", len(clusts))
	// 	//num_cluster := make([]int, 0)

	// 	// for m := 10; m <= 20; m += 20 {
	// 	// 	for l := 100; l <= 200; l += 100 {
	// 	// c := invClustering.GenerateClusters(t, l)
	// 	// tc := invClustering.GenerateTermClusters(c)

	// 	sol_cluster3 := make([]int, 0)
	// 	sol_cluster4 := make([]int, 0)
	// 	sol_invlist := make([]int, 0)
	// 	for i := 2; i < 10; i++ {
	// 		s := fmt.Sprintf("queries.%d", i)
	// 		data, _ := ioutil.ReadFile(s)
	// 		//fmt.Println(s)
	// 		d := strings.Split(string(data), "\n")
	// 		// sum_cluster1 := 0
	// 		// sum_cluster2 := 0
	// 		sum_cluster3 := 0
	// 		sum_cluster4 := 0
	// 		sum_invlist := 0
	// 		for dd := range d {
	// 			strings.Trim(d[dd], " ")
	// 			q := strings.Split(d[dd], " ")
	// 			if len(q) > 2 {
	// 				//	t0 := time.Now()
	// 				// n, s := tc.TopKQuery(q, 10, m)
	// 				n2, s2 := termCluster.TopKQuery(q, 10, 1)
	// 				// fmt.Println("n2 =",n2)
	// 				// fmt.Println("s2 =",s2)
	// 				sum_cluster3 += n2
	// 				sum_cluster4 += s2
	// 				//	t1 := time.Now()
	// 				// sum_cluster1 += n
	// 				// sum_cluster2 += s
	// 				//	fmt.Println(n2, s2)
	// 				//fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

	// 				//	t0 = time.Now()
	// 				size := t.TopKQuery(q, 10)
	// 				sum_invlist += size
	// 				//fmt.Println("sum = " , sum_invlist)
	// 				//	t1 = time.Now()
	// 				//fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
	// 				//fmt.Println(size)
	// 				// fmt.Println(i,"\t",size,"\t",n2,"\t",s2)
	// 			}
	// 		}
	// 		// sol_cluster = append(sol_cluster, sum_cluster1)
	// 		// sol_cluster2 = append(sol_cluster2, sum_cluster2)

	// 		sol_cluster3 = append(sol_cluster3, sum_cluster3)
	// 		sol_cluster4 = append(sol_cluster4, sum_cluster4)
	// 		//	fmt.Println("appending",sum_invlist)
	// 		sol_invlist = append(sol_invlist, sum_invlist)
	// 	}
	// 	for i := 0; i < len(sol_cluster4); i++ {
	// 		fmt.Println("|q|=", i+2, "\t", sol_cluster4[i], "\t", sol_invlist[i])
	// 	}
	// 	// sum_total_cluster1 := float64(0.0)
	// 	// sum_total_cluster2 := float64(0.0)

	// 	// sum_total_cluster3 := float64(0.0)
	// 	// sum_total_cluster4 := float64(0.0)

	// 	// sum_total_invlist := float64(0.0)
	// 	// for i := range sol_cluster {
	// 	// 	sum_total_cluster1 += float64(sol_cluster[i])
	// 	// 	sum_total_cluster2 += float64(sol_cluster2[i])

	// 	// 	sum_total_cluster3 += float64(sol_cluster3[i])
	// 	// 	sum_total_cluster4 += float64(sol_cluster4[i])
	// 	// 	sum_total_invlist += float64(sol_invlist[i])
	// 	// }
	// 	//fmt.Println("Inverted List:",sum_total_invlist)
	// 	// fmt.Println("Fake Cluster 1:",sum_total_cluster1)
	// 	// fmt.Println("Fake Cluster 1.1:",sum_total_cluster2)

	// 	// fmt.Println("REAL Cluster:",sum_total_cluster3)
	// 	// fmt.Println("REAL Cluster2:",sum_total_cluster4)
	// 	// 	}

	// 	// }
	// 	fmt.Println("----------------------------")
	// }
}
