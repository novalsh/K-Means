package main

import (
	"fmt"
	"math"
	"time"
)

type TodoItem struct {
	Name     string
	Deadline time.Time
	Level    float64
}

type Cluster struct {
	Centroid TodoItem
	Items    []TodoItem
}

func main() {
	// Memasukkan data Todo Item
	todoItems := []TodoItem{}
	for {
		fmt.Println("Masukkan nama Todo Item (tekan 'q' untuk keluar):")
		var name string
		fmt.Scanln(&name)
		if name == "q" {
			break
		}

		fmt.Println("Masukkan deadline (format: yy-mm-dd):")
		var deadlineStr string
		fmt.Scanln(&deadlineStr)
		deadline, err := time.Parse("2006-01-02", deadlineStr)
		if err != nil {
			fmt.Println("Format tanggal tidak valid. Coba lagi.")
			continue
		}

		fmt.Println("Masukkan level (float):")
		var level float64
		fmt.Scanln(&level)

		todoItem := TodoItem{
			Name:     name,
			Deadline: deadline,
			Level:    level,
		}

		todoItems = append(todoItems, todoItem)
	}

	// Inisialisasi jumlah cluster
	k := 2

	// Inisialisasi centroid awal secara acak
	centroids := initializeCentroids(todoItems, k)

	// Melakukan iterasi hingga konvergensi
	for {
		// Melakukan clustering pada data items
		clusters := clusterItems(todoItems, centroids)

		// Menghitung centroid baru berdasarkan rata-rata items pada setiap cluster
		newCentroids := calculateNewCentroids(clusters)

		// Memeriksa apakah centroid sudah konvergen (tidak berubah)
		if centroidsConverged(centroids, newCentroids) {
			break
		}

		// Mengupdate centroid dengan yang baru
		centroids = newCentroids
	}

	// Menampilkan hasil clustering
	clusters := clusterItems(todoItems, centroids)
	for i, cluster := range clusters {
		fmt.Printf("Cluster %d:\n", i+1)
		for _, item := range cluster.Items {
			fmt.Printf("- %s (Deadline: %s, Level: %.2f)\n", item.Name, item.Deadline.Format("2006-01-02"), item.Level)
		}
		fmt.Println()
	}
}

// Menginisialisasi centroid awal secara acak
func initializeCentroids(items []TodoItem, k int) []TodoItem {
	centroids := make([]TodoItem, k)
	for i := 0; i < k; i++ {
		centroids[i] = items[i]
	}
	return centroids
}

// Melakukan clustering pada data items
func clusterItems(items []TodoItem, centroids []TodoItem) []Cluster {
	clusters := make([]Cluster, len(centroids))
	for _, item := range items {
		minDistance := math.MaxFloat64
		clusterIndex := 0
		for i, centroid := range centroids {
			distance := euclideanDistance(item, centroid)
			if distance < minDistance {
				minDistance = distance
				clusterIndex = i
			}
		}
		clusters[clusterIndex].Items = append(clusters[clusterIndex].Items, item)
	}
	return clusters
}

// Menghitung jarak Euclidean antara dua Todo Items
func euclideanDistance(item1, item2 TodoItem) float64 {
	return math.Sqrt(math.Pow(item1.Level-item2.Level, 2))
}

// Menghitung centroid baru berdasarkan rata-rata items pada setiap cluster
func calculateNewCentroids(clusters []Cluster) []TodoItem {
	newCentroids := make([]TodoItem, len(clusters))
	for i, cluster := range clusters {
		var sumLevel float64
		for _, item := range cluster.Items {
			sumLevel += item.Level
		}
		centroid := TodoItem{
			Name:     fmt.Sprintf("Centroid %d", i+1),
			Deadline: time.Time{},
			Level:    sumLevel / float64(len(cluster.Items)),
		}
		newCentroids[i] = centroid
	}
	return newCentroids
}

// Memeriksa apakah centroid sudah konvergen (tidak berubah)
func centroidsConverged(centroids, newCentroids []TodoItem) bool {
	for i := range centroids {
		if centroids[i].Level != newCentroids[i].Level {
			return false
		}
	}
	return true
}
