package main

import (
	"fmt"
//	"hashlist"
	"invlist"
	"math"
//	"math/rand"
	"os"
	"sort"
//	"svs"
//	"time"
//	"treap"

	//"flame"
	//	"unsafe"
	//	"reflect"
	//"time"
	"io/ioutil"
	//"math"
	"strings"

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

	fmt.Println("N = ", len(t.DocPost))
	fmt.Println("M = ", len(t.TermMap))
	doc_vector := make([][]float64, len(t.DocPost))
	for i := range t.DocPost {
		doc_vector[i] = make([]float64, len(t.TermMap))
		for j := range t.DocPost[i].Terms {
			doc_vector[i][sort.SearchStrings(t.TermMap, j)] = float64(t.DocPost[i].Terms[j])
		}
	}

	f := invClustering.CreateFlame()
	fmt.Println("Euclidean Distance")
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.Euclidean)
	f.Write("clusters_euclidean2")

	fmt.Println("Cosine")
	f = invClustering.CreateFlame()
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.Cosine)
	f.Write("clusters_cosine2")

	fmt.Println("Cosine Distance")
	f = invClustering.CreateFlame()
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.CosineDist)
	f.Write("clusters_cosinedist2")

	fmt.Println("Pearson")
	f = invClustering.CreateFlame()
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.Pearson)
	f.Write("clusters_pearson2")

	fmt.Println("Pearson Distance")
	f = invClustering.CreateFlame()
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.PearsonDist)
	f.Write("clusters_pearsondist2")

	fmt.Println("Dot Product")
	f = invClustering.CreateFlame()
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.DotProduct)
	f.Write("clusters_dot2")

	fmt.Println("Dot Product Distance")
	f = invClustering.CreateFlame()
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.DotProductDist)
	f.Write("clusters_dotdist2")

	fmt.Println("Covariance")
	f = invClustering.CreateFlame()
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.Covariance)
	f.Write("clusters_covariance2")

	fmt.Println("Covariance Distance")
	f = invClustering.CreateFlame()
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.DistCovariance)
	f.Write("clusters_covariancedist2")

	fmt.Println("Jaccard ")
	f = invClustering.CreateFlame()
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.Jaccard)
	f.Write("clusters_jaccard2")

	fmt.Println("Jaccard Pearson ")
	f = invClustering.CreateFlame()
	f.SetMatrix(doc_vector, len(t.DocPost), len(t.TermMap), invClustering.JaccardPearson)
	f.Write("clusters_jaccard_pearson2")

	clusters := []string{"clusters_euclidean", "clusters_cosine", "clusters_cosinedist", "clusters_pearson", "clusters_pearsondist", "clusters_dot", "clusters_dotdist", "clusters_covariance", "clusters_covariancedist", "clusters_jaccard", "clusters_jaccard_pearson"}
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
	t.GenerateQueries(2, 11, 1000, 1)

	for i := 0; i < len(clusters); i++ {
		f := invClustering.LoadClusters(clusters[i] + ".serial")
		fmt.Println("-------------------------------")
		fmt.Println(clusters[i])
		f.DefineSupports(1, 1)
		f.LocalApproximation(1, 1)
		f.MakeClusters(0.5)
		// //fmt.Println(f)
		clusts := make([]*invClustering.Cluster, 0)
		//fmt.Println(f.Clusters)
		for i := 0; i < len(f.Clusters); i++ {
			if len(f.Clusters[i].Array) != 0 {
				if len(f.Clusters[i].Array) > 128 {
					clusts_aux := t.NewFromCluster(f.Clusters[i].Array, int(math.Sqrt(float64(len(f.Clusters[i].Array))))+1)
					for j := 0; j < len(clusts_aux); j++ {
						clusts = append(clusts, clusts_aux[j])
					}
				} else {
					clusts_aux := t.NewFromCluster(f.Clusters[i].Array, 128)
					for j := 0; j < len(clusts_aux); j++ {
						clusts = append(clusts, clusts_aux[j])
					}
				}

			}
		}
		//fmt.Println(clusts)
		termCluster := invClustering.GenerateTermClusters(clusts)
		termCluster.Stats()
		fmt.Println("# clusters = ", len(clusts))
		//num_cluster := make([]int, 0)

		// for m := 10; m <= 20; m += 20 {
		// 	for l := 100; l <= 200; l += 100 {
		// c := invClustering.GenerateClusters(t, l)
		// tc := invClustering.GenerateTermClusters(c)

		sol_cluster3 := make([]int, 0)
		sol_cluster4 := make([]int, 0)
		sol_invlist := make([]int, 0)
		//ratio_list := make([]float64,0)
		nc_list := make([]float64,0)
		n1_list := make([]float64,0)
		for i := 2; i < 10; i++ {
			s := fmt.Sprintf("queries.%d", i)
			data, _ := ioutil.ReadFile(s)
			//fmt.Println(s)
			d := strings.Split(string(data), "\n")
			// sum_cluster1 := 0
			// sum_cluster2 := 0
			amount_clusters := 0
			amount_items := 0
			sum_cluster3 := 0
			sum_cluster4 := 0
			sum_invlist := 0
			for dd := range d {
				strings.Trim(d[dd], " ")
				q := strings.Split(d[dd], " ")
				if len(q) == 2 {
					n2, s2 := termCluster.TopKQuery(q, 10, 1)
					nC := math.Min(float64(termCluster.GetDocsFromTerm(q[0])),float64(termCluster.GetDocsFromTerm(q[0])))			
					sum_cluster3 += n2
					sum_cluster4 += s2
					amount_clusters += n2;
					amount_items += s2
					size := t.TopKQuery(q, 10)
					n1 := math.Min(float64(t.GetDocsFromTerm(q[0])),float64(t.GetDocsFromTerm(q[1])))			
					n1_list = append(n1_list,n1)
					nc_list = append(nc_list,nC)
					sum_invlist += size				
				}
			}
			sol_cluster3 = append(sol_cluster3, sum_cluster3)
			sol_cluster4 = append(sol_cluster4, sum_cluster4)
			sol_invlist = append(sol_invlist, sum_invlist)
		}
		for i := 0; i < len(sol_cluster4); i++ {
			fmt.Println("|q|=", i+2, "\t", sol_cluster4[i], "\t", sol_cluster3[i],'\t', sol_invlist[i])
		}
		fmt.Println("----------------------------")
		for i:=0; i< len(n1_list);i++ {
			fmt.Println(n1_list[i],"\t",nc_list[i],"\t",n1_list[i]/nc_list[i],nc_list[i]/n1_list[1])
		}
		fmt.Println("----------------------------")
	}
}
