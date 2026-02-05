package tui

import (
	"fmt"
	"testing"
)

var (
	benchmarkResultItems    []WorktreeItem
	benchmarkResultPointers []*WorktreeItem
)

// BenchmarkPruneSearchNested mimics the current logic in enterPruneMode
func BenchmarkPruneSearchNested(b *testing.B) {
	// Setup scenarios
	scenarios := []struct {
		numItems int
		numStale int
	}{
		{50, 5},
		{100, 10},
		{1000, 100},
		{5000, 500},
		{10000, 1000},
	}

	for _, sc := range scenarios {
		b.Run(fmt.Sprintf("Items-%d_Stale-%d", sc.numItems, sc.numStale), func(b *testing.B) {
			// Prepare data
			items := make([]WorktreeItem, sc.numItems)
			for i := 0; i < sc.numItems; i++ {
				items[i] = WorktreeItem{Name: fmt.Sprintf("wt-%d", i), Path: "/tmp", IsMain: false}
			}

			staleNames := make([]string, sc.numStale)
			for i := 0; i < sc.numStale; i++ {
				staleNames[i] = fmt.Sprintf("wt-%d", i*2)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var found []WorktreeItem
				for _, name := range staleNames {
					for _, item := range items {
						if item.Name == name && !item.IsMain {
							found = append(found, item)
							break
						}
					}
				}
				benchmarkResultItems = found
			}
		})
	}
}

// BenchmarkPruneSearchMap mimics the optimized logic using a map
func BenchmarkPruneSearchMap(b *testing.B) {
	// Setup scenarios
	scenarios := []struct {
		numItems int
		numStale int
	}{
		{50, 5},
		{100, 10},
		{1000, 100},
		{5000, 500},
		{10000, 1000},
	}

	for _, sc := range scenarios {
		b.Run(fmt.Sprintf("Items-%d_Stale-%d", sc.numItems, sc.numStale), func(b *testing.B) {
			// Prepare data
			items := make([]WorktreeItem, sc.numItems)
			for i := 0; i < sc.numItems; i++ {
				items[i] = WorktreeItem{Name: fmt.Sprintf("wt-%d", i), Path: "/tmp", IsMain: false}
			}

			staleNames := make([]string, sc.numStale)
			for i := 0; i < sc.numStale; i++ {
				staleNames[i] = fmt.Sprintf("wt-%d", i*2)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var found []WorktreeItem

				// Build map
				itemMap := make(map[string]WorktreeItem, len(items))
				for _, item := range items {
					itemMap[item.Name] = item
				}

				for _, name := range staleNames {
					if item, ok := itemMap[name]; ok {
						if !item.IsMain {
							found = append(found, item)
						}
					}
				}
				benchmarkResultItems = found
			}
		})
	}
}

// BenchmarkPruneSearchMapPointer mimics the optimized logic using a map of pointers
func BenchmarkPruneSearchMapPointer(b *testing.B) {
	// Setup scenarios
	scenarios := []struct {
		numItems int
		numStale int
	}{
		{50, 5},
		{100, 10},
		{1000, 100},
		{5000, 500},
		{10000, 1000},
	}

	for _, sc := range scenarios {
		b.Run(fmt.Sprintf("Items-%d_Stale-%d", sc.numItems, sc.numStale), func(b *testing.B) {
			// Prepare data
			items := make([]WorktreeItem, sc.numItems)
			for i := 0; i < sc.numItems; i++ {
				items[i] = WorktreeItem{Name: fmt.Sprintf("wt-%d", i), Path: "/tmp", IsMain: false}
			}

			staleNames := make([]string, sc.numStale)
			for i := 0; i < sc.numStale; i++ {
				staleNames[i] = fmt.Sprintf("wt-%d", i*2)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var found []*WorktreeItem

				// Build map
				itemMap := make(map[string]*WorktreeItem, len(items))
				for i := range items {
					itemMap[items[i].Name] = &items[i]
				}

				for _, name := range staleNames {
					if item, ok := itemMap[name]; ok {
						if !item.IsMain {
							found = append(found, item)
						}
					}
				}
				benchmarkResultPointers = found
			}
		})
	}
}
