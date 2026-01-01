package stats

import (
	"sort"

	"github.com/bobparsons/rootcamp/internal/types"
)

type ProgressStats struct {
	Total      int
	Completed  int
	Percentage float64
}

type LevelStats struct {
	Level string
	Stats ProgressStats
}

type ModuleStats struct {
	Module string
	Stats  ProgressStats
}

type OverallProgress struct {
	Overall  ProgressStats
	ByLevel  []LevelStats
	ByModule []ModuleStats
}

func CalculateProgress(
	lessons []types.Lesson,
	progressMap map[string]*types.UserProgress,
) OverallProgress {
	overall := ProgressStats{Total: len(lessons)}
	levelMap := make(map[string]*ProgressStats)
	moduleMap := make(map[string]*ProgressStats)

	for _, lesson := range lessons {
		if prog, exists := progressMap[lesson.ID]; exists && prog.Completed {
			overall.Completed++
		}

		if lesson.Level != "" {
			if _, exists := levelMap[lesson.Level]; !exists {
				levelMap[lesson.Level] = &ProgressStats{}
			}
			levelMap[lesson.Level].Total++
			if prog, exists := progressMap[lesson.ID]; exists && prog.Completed {
				levelMap[lesson.Level].Completed++
			}
		}

		if lesson.Module != "" {
			if _, exists := moduleMap[lesson.Module]; !exists {
				moduleMap[lesson.Module] = &ProgressStats{}
			}
			moduleMap[lesson.Module].Total++
			if prog, exists := progressMap[lesson.ID]; exists && prog.Completed {
				moduleMap[lesson.Module].Completed++
			}
		}
	}

	overall.Percentage = calculatePercentage(overall.Completed, overall.Total)

	byLevel := make([]LevelStats, 0, len(levelMap))
	for level, stats := range levelMap {
		stats.Percentage = calculatePercentage(stats.Completed, stats.Total)
		byLevel = append(byLevel, LevelStats{Level: level, Stats: *stats})
	}
	sort.Slice(byLevel, func(i, j int) bool {
		order := map[string]int{"beginner": 1, "intermediate": 2, "advanced": 3, "expert": 3}
		return order[byLevel[i].Level] < order[byLevel[j].Level]
	})

	byModule := make([]ModuleStats, 0, len(moduleMap))
	for module, stats := range moduleMap {
		stats.Percentage = calculatePercentage(stats.Completed, stats.Total)
		byModule = append(byModule, ModuleStats{Module: module, Stats: *stats})
	}
	sort.Slice(byModule, func(i, j int) bool {
		return byModule[i].Module < byModule[j].Module
	})

	return OverallProgress{
		Overall:  overall,
		ByLevel:  byLevel,
		ByModule: byModule,
	}
}

func calculatePercentage(completed, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(completed) / float64(total) * 100
}

func RenderProgressBar(percentage float64) string {
	const barWidth = 60
	filled := min(int(percentage/100*barWidth), barWidth)

	bar := ""
	for i := range barWidth {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
}
