// movie_history_table.go — table-formatted output for unified history
package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/movie-cli-v4/db"
)

func printHistoryTableUnified(records []unifiedRecord) {
	idW := 6
	typeW := 14
	statusW := 8
	dateW := 19
	detailW := 40

	fmt.Println()
	fmt.Printf("  %-*s │ %-*s │ %-*s │ %-*s │ %-*s\n",
		idW, "ID",
		typeW, "Type",
		statusW, "Status",
		dateW, "Date",
		detailW, "Detail")

	fmt.Printf("  %s─┼─%s─┼─%s─┼─%s─┼─%s\n",
		strings.Repeat("─", idW),
		strings.Repeat("─", typeW),
		strings.Repeat("─", statusW),
		strings.Repeat("─", dateW),
		strings.Repeat("─", detailW))

	for i := range records {
		status := "OK"
		if r.IsReverted {
			status = "Reverted"
		}

		prefix := r.Source[0:1] // "m" or "a"
		idStr := fmt.Sprintf("%s-%d", prefix, r.ID)

		fmt.Printf("  %-*s │ %-*s │ %-*s │ %-*s │ %-*s\n",
			idW, idStr,
			typeW, truncate(r.Type, typeW),
			statusW, status,
			dateW, truncate(r.Timestamp, dateW),
			detailW, truncate(r.Detail, detailW))
	}

	fmt.Printf("  %s─┴─%s─┴─%s─┴─%s─┴─%s\n",
		strings.Repeat("─", idW),
		strings.Repeat("─", typeW),
		strings.Repeat("─", statusW),
		strings.Repeat("─", dateW),
		strings.Repeat("─", detailW))

	fmt.Printf("\n  Total: %d records\n\n", len(records))
}
