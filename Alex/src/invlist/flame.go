package invClustering

import (
	//"fmt"
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"math"
	"os"
	"runtime"
)

type Flame struct {
	simtype   int
	N         int
	KMAX      int
	K         int
	Graph     [][]int
	Dists     [][]float64
	Weights   [][]float64
	Nncounts  []int
	Obtypes   []int
	Fuzzyship [][]float64
	Cscount   int
	Clusters  []IntArray
	Count     int
}

type IntArray struct {
	Array []int
}

type IndexFloat struct {
	Index int
	Value float64
}

var EPSILON float64 = 1E-9

func Euclidean(x []float64, y []float64, m int) float64 {
	d := float64(0)

	for i := 0; i < m; i++ {
		if x[i] == y[i] {
			continue
		}
		d += (x[i] - y[i]) * (x[i] - y[i])
		//	fmt.Println("x[",i,"]",x[i])
	}

	//	fmt.Println("returning:", float64(math.Sqrt(d)))
	return float64(math.Sqrt(d))
}

func Cosine(x []float64, y []float64, m int) float64 {
	r := float64(0)
	x2 := float64(0)
	y2 := float64(0)
	for i := 0; i < m; i++ {
		r += x[i] * y[i]
		x2 += x[i] * x[i]
		y2 += y[i] * y[i]
	}
	return (r/math.Sqrt(x2*y2) + EPSILON)
}
func Pearson(x []float64, y []float64, m int) float64 {
	r := float64(0)
	x2 := float64(0)
	y2 := float64(0)
	xavg := float64(0)
	yavg := float64(0)

	if m == 0 {
		return float64(0)
	}
	for i := 0; i < m; i++ {
		xavg += x[i]
		yavg += y[i]
	}
	xavg = xavg / float64(m)
	yavg = yavg / float64(m)
	for i := 0; i < m; i++ {
		r += (x[i] - xavg) * (y[i] - yavg)
		x2 += (x[i] - xavg) * (x[i] - xavg)
		y2 += (y[i] - yavg) * (y[i] - yavg)
	}
	return r / (math.Sqrt(x2*y2) + EPSILON)
}

func PearsonDist(x []float64, y []float64, m int) float64 {
	return 1 - Pearson(x, y, m)
}

func CosineDist(x []float64, y []float64, m int) float64 {
	return 1 - Cosine(x, y, m)
}
func UCPearson(x []float64, y []float64, m int) float64 {
	r := float64(0)
	x2 := float64(0)
	y2 := float64(0)
	xavg := float64(0)
	yavg := float64(0)

	for i := 0; i < m; i++ {
		xavg += x[i]
		yavg += y[i]
	}
	xavg = xavg / float64(m)
	yavg = yavg / float64(m)
	for i := 0; i < m; i++ {
		r += x[i] * y[i]
		x2 += (x[i] - xavg) * (x[i] - xavg)
		y2 += (y[i] - yavg) * (y[i] - yavg)
	}
	return r / (math.Sqrt(x2*y2) + EPSILON)
}

func Covariance(x []float64, y []float64, m int) float64 {

	r := float64(0)
	xavg := float64(0)
	yavg := float64(0)
	for i := 0; i < m; i++ {
		xavg += x[i]
		yavg += y[i]
	}
	xavg = xavg / float64(m)
	yavg = yavg / float64(m)
	for i := 0; i < m; i++ {
		r += (x[i] - xavg) * (y[i] - yavg)
	}
	return r / float64(m-1)
}

func DistCovariance(x []float64, y []float64, m int) float64 {
	return float64(1) - Covariance(x, y, m)
}

func DotProduct(x []float64, y []float64, m int) float64 {
	r := float64(0)
	for i := 0; i < m; i++ {
		r += x[i] * y[i]
	}
	if m == 0 {
		return 0
	}
	return r / float64(m)
}
func DotProductDist(x []float64, y []float64, m int) float64 {
	return 1 - DotProduct(x, y, m)
}

func Jaccard(x []float64, y []float64, m int) float64 {
	r := float64(0)
	for i := 0; i < m; i++ {
		if x[i] != 0 && y[i] != 0 {
			r++
		}
	}
	return r / float64(m)
}
func JaccardPearson(x []float64, y []float64, m int) float64 {
	return Jaccard(x, y, m) * Pearson(x, y, m)
}
func PartialQuickSort(data []IndexFloat, first int, last int, part int) {
	lower := first + 1
	upper := last
	var pivot float64
	var val IndexFloat
	if first >= last {
		return
	}
	val = data[first]
	data[first] = data[(first+last)/2]
	data[(first+last)/2] = val
	pivot = data[first].Value

	for lower <= upper {
		for lower <= last && data[lower].Value < pivot {
			lower++
		}
		for pivot < data[upper].Value {
			upper--
		}
		if lower < upper {
			val = data[lower]
			data[lower] = data[upper]
			data[upper] = val
			upper--
		}
		lower++
	}
	val = data[first]
	data[first] = data[upper]
	data[upper] = val
	if first < upper-1 {
		PartialQuickSort(data, first, upper-1, part)
	}
	if upper >= part {
		return
	}
	if upper+1 < last {
		PartialQuickSort(data, upper+1, last, part)
	}
}

func CreateFlame() *Flame {
	f := new(Flame)
	return f
}

func (self *Flame) SetMatrix(data [][]float64, n int, m int, d func([]float64, []float64, int) float64) {
	MAX := math.Sqrt(float64(n)) + float64(10)
	//vals := make([]IndexFloat, n)

	if MAX >= float64(n) {
		MAX = float64(n - 1)
	}
	self.N = n
	self.KMAX = int(MAX)
	id := make([]int, 0)
	self.Graph = make([][]int, n)
	self.Dists = make([][]float64, n)
	self.Weights = make([][]float64, n)
	self.Nncounts = make([]int, n)
	self.Obtypes = make([]int, n)
	//	self.Fuzzyship = make([][]float64,n)
	for i := 0; i < n; i++ {
		self.Graph[i] = make([]int, int(MAX))
		self.Dists[i] = make([]float64, int(MAX))
		self.Weights[i] = make([]float64, int(MAX))
		id = append(id, i)
	}
	//fmt.Println(id)
	self.WorkParallel(id, data, n, MAX, 6, d)
}

func (self *Flame) WorkParallel(id []int, data [][]float64, n int, MAX float64, threads int, d func([]float64, []float64, int) float64) {
	queue := make(chan int)
	ncpu := threads
	/*	ncpu := runtime.NumCPU()
		if threads < ncpu {
			ncpu = threads
		}*/
	runtime.GOMAXPROCS(ncpu)
	for i := 0; i < ncpu; i++ {
		go self.Worker(i, queue, data, n, MAX, d)
	}
	// master: give work
	for i := range id {
		queue <- id[i]

	}
	for n := 0; n < ncpu; n++ {
		queue <- -1
	}
}

func (self *Flame) Worker(id int, queue chan int, data [][]float64, n int, MAX float64, d func([]float64, []float64, int) float64) {
	var inv int
	for {
		vals := make([]IndexFloat, n)
		//fmt.Println(len(queue))
		inv = <-queue
		if inv == -1 {
			break
		}
		//fmt.Println("Worker ", id ," has id  ", inv )
		for j := 0; j < n; j++ {
			vals[j].Index = j
			vals[j].Value = d(data[inv], data[j], len(data[j]))
			//fmt.Println("value = ", vals[j].Value)
		}
		PartialQuickSort(vals, 0, n-1, int(MAX+1))
		for j := 0; j < int(MAX); j++ {
			self.Graph[inv][j] = vals[j+1].Index
			self.Dists[inv][j] = vals[j+1].Value
		}
	}
}

func (self *Flame) SetDataMatrix(data [][]float64, n int, m int, dt func([]float64, []float64, int) float64) {
	self.simtype = n
	self.SetMatrix(data, n, m, dt)
	// for i:=0;i<self.N;i++	{
	// 	fmt.Println("i = ",i,self.Dists[i])
	// 	//fmt.Println("j = ",len(self.Weights[i]))
	// }
}

func (self *Flame) DefineSupports(knn int, thd float64) {

	// for i:=0;i<self.N;i++	{
	// 		fmt.Println("i = ",i,self.Dists[i])
	// 		//fmt.Println("j = ",len(self.Weights[i]))
	// 	}
	var k int
	n := self.N
	kmax := self.KMAX
	density := make([]float64, n)
	d, sum, sum2, fmin, fmax := 0.000, 0.000, 0.000, 0.000, 0.000

	if knn > kmax {
		knn = kmax
	}
	self.K = knn
	for i := 0; i < n; i++ {
		k = knn
		d = self.Dists[i][knn-1]
		//fmt.Println("i,knn-1",i,knn-1,"=>",d)

		for j := knn; j < kmax; j++ {
			if float64(self.Dists[i][j]) == float64(d) {
				k++
			} else {
				break
			}
		}
		//fmt.Println(k)
		self.Nncounts[i] = k
		sum = 0.5 * float64(k) * (float64(k) + 1.0)

		for j := 0; j < k; j++ {
			self.Weights[i][j] = (float64(k) - float64(j)) / sum
			//	fmt.Println(self.Weights[i][j])	
		}

		sum = 0.0
		for j := 0; j < k; j++ {
			sum += self.Dists[i][j]
		}
		density[i] = 1.0 / (sum + EPSILON)
		//	fmt.Println("density = " , density[i])
	}
	// fmt.Println(self.Weights)
	// fmt.Println(density)

	sum = 0.0
	sum2 = 0.0
	for i := 0; i < n; i++ {
		sum += density[i]
		sum2 += density[i] * density[i]
	}
	sum = sum / float64(n)
	thd = sum + thd*math.Sqrt(sum2/float64(n)-sum*sum)
	//fmt.Println("thd = ", thd)
	for i := 0; i < n; i++ {
		self.Obtypes[i] = 0
	}
	self.Cscount = 0
	for i := 0; i < n; i++ {
		k = self.Nncounts[i]
		fmax = 0.0
		fmin = density[i] / density[self.Graph[i][0]]
		//fmt.Println(fmin)
		for j := 1; j < k; j++ {
			d = density[i] / density[self.Graph[i][j]]
			if d > fmax {
				fmax = d
			}
			if d < fmin {
				fmin = d
			}
			//fmt.Println(self.Obtypes)
			if self.Obtypes[self.Graph[i][j]] == 0 {
				fmin = 0.0
			}
		}
		//fmt.Println(fmin)
		if fmin >= 1.0 {
			self.Cscount = self.Cscount + 1
			//	fmt.Println("Cscount = ", self.Cscount)
			self.Obtypes[i] = 1
		} else if fmax <= 1.0 && density[i] < thd {
			self.Obtypes[i] = 2
		}
	}
	// for i:=0;i<self.N;i++	{
	// 	fmt.Println("i = ",i,self.Weights[i])
	// 	fmt.Println("j = ",len(self.Weights[i]))
	// }
	//fmt.Println(self.Obtypes)
}

func (self *Flame) LocalApproximation(steps int, episilon float64) {
	k, n := self.N, self.N
	m := self.Cscount
	// fmt.Println("Cscount=", self.Cscount)
	// fmt.Println("n =", n)
	// fmt.Println("k =", k)

	self.Fuzzyship = make([][]float64, n)
	for i := 0; i < n; i++ {
		self.Fuzzyship[i] = make([]float64, m+1)
	}

	fuzzyship2 := make([][]float64, n)
	even := 0
	var dev float64

	k = 0
	for i := 0; i < n; i++ {
		fuzzyship2[i] = make([]float64, m+1)
		if self.Obtypes[i] == 1 {
			// fmt.Println("k=",k)
			// fmt.Println("size=",len(self.Fuzzyship[i]))
			// fmt.Println("m = ",m)
			self.Fuzzyship[i][k] = 1.0
			fuzzyship2[i][k] = 1.0
			k++
		} else if self.Obtypes[i] == 2 {
			self.Fuzzyship[i][m] = 1.0
			fuzzyship2[i][m] = 1.0
		} else {
			for j := 0; j <= m; j++ {
				// fmt.Println("len of fuzzyship0",len(fuzzyship0)," j = ", j)
				// fmt.Println("len of fuzzyship0[",i,"]=",len(fuzzyship0[i])," j = ", j)

				self.Fuzzyship[i][j] = float64(1.0) / (float64(m) + float64(1.0))
				fuzzyship2[i][j] = float64(1.0) / (float64(m) + float64(1.0))
			}
		}
	}
	var fuzzy2 [][]float64
	for t := 0; t < steps; t++ {
		dev = 0.0
		for i := 0; i < n; i++ {
			knn := self.Nncounts[i]
			//ids := self.Graph[i]
			//self.Weights[i] := self.Weights[i]
			sum := 0.0

			if self.Obtypes[i] != 0 {
				continue
			}
			if even == 0 {
				self.Fuzzyship[i] = fuzzyship2[i]
				fuzzy2 = self.Fuzzyship
			}
			for j := 0; j <= m; j++ {
				self.Fuzzyship[i][j] = 0.0
				for k := 0; k < knn; k++ {
					self.Fuzzyship[i][j] += float64(self.Weights[i][k]) * fuzzy2[self.Graph[i][k]][j]
				}
				// fmt.Println("j = ", j)
				// fmt.Println("i = ", i)

				dev += (self.Fuzzyship[i][j] - fuzzy2[i][j]) * (self.Fuzzyship[i][j] - fuzzy2[i][j])
				sum += self.Fuzzyship[i][j]
			}
			for j := 0; j <= m; j++ {
				self.Fuzzyship[i][j] = self.Fuzzyship[i][j] / sum
			}
			//fmt.Println("i=",i,"->",self.Fuzzyship[i])

		}
		even = (even + 1) % 2
		if dev < episilon {
			//			fmt.Println("BREAK!")
			break
		}
	}
	for i := 0; i < n; i++ {
		knn := self.Nncounts[i]
		// self.Graph[i] := self.Graph[i]
		// self.Weights[i] := self.Weights[i]
		fuzzy2 := fuzzyship2
		for j := 0; j <= m; j++ {
			self.Fuzzyship[i][j] = float64(0.0)
			for k := 0; k < knn; k++ {
				self.Fuzzyship[i][j] += self.Weights[i][k] * fuzzy2[self.Graph[i][k]][j]
			}
			dev += (self.Fuzzyship[i][j] - fuzzy2[i][j]) * (self.Fuzzyship[i][j] - fuzzy2[i][j])
		}
	}
}

func (self *Flame) MakeClusters(thd float64) {
	var i int
	var j int
	N := self.N
	C := self.Cscount + 1
	// fmt.Println("N = ", N)
	// fmt.Println("C = ", C)

	var fmax float64

	vals := make([]IndexFloat, N)
	//fmt.Println(vals)
	for i = 0; i < N; i++ {
		vals[i].Index = i
		vals[i].Value = 0.0
		for j = 0; j < C; j++ {
			fs := self.Fuzzyship[i][j]
			if fs > EPSILON {
				vals[i].Value -= fs * math.Log(fs)
				//fmt.Println(EPSILON)
				//fmt.Println("Entre!!!",vals[i].Value)
			}
		}
	}
	imax := 0
	PartialQuickSort(vals, 0, N-1, N)
	self.Clusters = make([]IntArray, C)
	if thd < 0 || thd > 1.0 {
		for i := 0; i < N; i++ {

			id := vals[i].Index
			// fmt.Println("id = ",id)
			fmax = 0
			imax = -1
			for j := 0; j < C; j++ {
				if self.Fuzzyship[id][j] > fmax {
					imax = j
					fmax = self.Fuzzyship[id][j]
					//fmt.Println("fmax=",fmax)
				}
			}
			//fmt.Println("Adding id ", id, " to ", imax)
			self.Clusters[imax].Array = append(self.Clusters[imax].Array, id)
		}
	} else {
		for i := 0; i < N; i++ {
			id := vals[i].Index
			imax = -1
			for j := 0; j < C; j++ {
				if self.Fuzzyship[id][j] > thd || (j == C-1 && imax < 0) {
					imax = j
					self.Clusters[j].Array = append(self.Clusters[j].Array, id)
					// fmt.Println("Adding id 2 ", id, " to ", j)
				}
			}
		}
	}
	C = 0
	for i := 0; i < self.Cscount; i++ {
		if len(self.Clusters[i].Array) > 0 {
			self.Clusters[C] = self.Clusters[i]
			C++
		}
	}

	self.Clusters[C] = self.Clusters[self.Cscount]
	C++
	self.Count = C
	//	fmt.Println(self.Clusters)
}
func LoadClusters(file string) *Flame {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Error Opening the file")
	}
	var input bytes.Buffer
	input.Write([]byte(data))
	dec := gob.NewDecoder(&input)
	var inv *Flame
	err = dec.Decode(&inv)
	return inv
}

func (self *Flame) Write(file string) {
	fo, err := os.Create(file + ".serial")
	if err != nil {
		panic("Error Opening the file")
	}
	defer fo.Close()
	var output bytes.Buffer
	enc := gob.NewEncoder(&output)
	err = enc.Encode(self)
	fo.Write(output.Bytes())
}
