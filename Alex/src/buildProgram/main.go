package main

import (
	"fmt"
	"invlist"
	"os"
	"treap"
	"math"
	//"sort"

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
	/*	N := 1000
			M := 1000
			data := make([][]float64, N)
			for i := 0; i < N; i++ {
				data[i] = make([]float64, M)
				for j := 0; j < M; j++ {
					data[i][j] = rand.NormFloat64()*100000000 + 1000
				}
			}
			//fmt.Println(data)
			f := invClustering.CreateFlame()
			f.SetMatrix(data, N, M, invClustering.Euclidean)
			f.DefineSupports(3, 100)
			f.LocalApproximation(10, 1E-9)
			f.MakeClusters(0.1)
			f.Write("clusters")
			fmt.Println("cluster size ", len(f.Clusters))
			clusts := make(map[int]int,0)
			for i := 0; i < len(f.Clusters); i++ {
			fmt.Println("Cluster with ", len(f.Clusters[i].Array), " elements")
			for j := 0; j < len(f.Clusters[i].Array); j++ {
				clusts[f.Clusters[i].Array[j]]++
				fmt.Printf(",%d", f.Clusters[i].Array[j])
			}
		}
		fmt.Println("")

		for i:=0;i<len(clusts);i++ {
				if (clusts[i] != 1) { fmt.Println("[",i,"] ->",clusts[i]) }
			}
		/*/
	 args := os.Args
	// if len(args) != 2 {
	// 	fmt.Printf("Usage: %s <input>\n", args[0])
	// 	return
	// }

	t := invClustering.LoadFiles(args[1])
	t.Write("test_large.invlist")
	//t := invClustering.Load("test.invlist.serial")
	tr := make([]*treap.Tree,0)
	term := 0
	count := 0
	total_height := float64(0)
	total_balanced := float64(0)
	total_len := float64(0)
	for i := range t.Postings {
		tr = append(tr,treap.NewTree(IntLess))
		//fmt.Println("Adding term:",i)
		size := len(t.Postings[i].Content)
		if (size > 10) {
			for k := range t.Postings[i].Content {
				//fmt.Println(k)
				tr[term].Insert(k,t.Postings[i].Content[k])
			}
			balanced :=  math.Ceil((math.Log2(float64(len(t.Postings[i].Content)))))
			height := tr[term].GetHeightTree(tr[term].GetRoot())
			total_height += float64(height)
			total_balanced +=balanced
			total_len += float64(size)
			count++
			// fmt.Println("Length",len(t.Postings[i].Content))
			// fmt.Println("balanced = ",balanced)
			// fmt.Println("treap = ",height)
			// fmt.Println("diference:",int(height) -int(balanced))

		}
		term++
	}
	fmt.Println("AVG HEIGHT:",float64(total_balanced)/float64(count),"\t AVG TREAP:",float64(float64(total_height)/float64(count)))
	fmt.Println("AVG LEN: ",float64(total_len)/float64(count))
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
