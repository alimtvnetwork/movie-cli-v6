// types.go — shared option structs for functions with >3 parameters.
package cmd

import (
	"bufio"
	"net/http"
	"os"

	"github.com/alimtvnetwork/movie-cli-v4/cleaner"
	"github.com/alimtvnetwork/movie-cli-v4/db"
	"github.com/alimtvnetwork/movie-cli-v4/tmdb"
)

// MoveContext groups parameters for batch and interactive move flows.
type MoveContext struct {
	Files     []os.FileInfo
	SourceDir string
	Home      string
	Database  *db.DB
	Scanner   *bufio.Scanner
}

// CleanupContext groups parameters for popout folder cleanup operations.
type CleanupContext struct {
	BatchID  string
	Scanner  *bufio.Scanner
	Database *db.DB
}

// ScanServiceConfig groups parameters for post-scan services (REST, watch).
type ScanServiceConfig struct {
	ScanDir   string
	OutputDir string
	Database  *db.DB
	Creds     tmdbCredentials
}

// SuggestCollector groups parameters for suggestion collection helpers.
type SuggestCollector struct {
	ExistingIDs map[int]bool
	Client      *tmdb.Client
	Count       int
}

// StatsCounts groups the three media count values for stats rendering.
type StatsCounts struct {
	Movies int
	TV     int
	Total  int
}

// LsPage groups pagination parameters for list display.
type LsPage struct {
	Offset   int
	PageSize int
	Total    int
}

// RecursiveWalkOpts groups depth-control parameters for recursive directory walks.
type RecursiveWalkOpts struct {
	BaseParts int
	MaxDepth  int
}

// ThumbnailInput groups parameters for poster/thumbnail download functions.
type ThumbnailInput struct {
	PosterPath string
	OutputDir  string
	Client     *tmdb.Client
	Database   *db.DB
	Media      *db.Media
}

// HistoryLogInput groups parameters for saving move history to JSON log.
type HistoryLogInput struct {
	BasePath string
	Title    string
	FromPath string
	ToPath   string
	Year     int
}

// ScanLoopConfig groups parameters for the main scan processing loop.
type ScanLoopConfig struct {
	ScanDir   string
	BatchID   string
	Client    *tmdb.Client
	JSONItems *[]scanJSONItem
	UseJSON   bool
	UseTable  bool
	HasTMDb   bool
}

// ScanOutputOpts groups output format flags used during scan processing.
type ScanOutputOpts struct {
	UseTable bool
	UseJSON  bool
}

// DryRunCounters groups counter pointers for dry-run scan output.
type DryRunCounters struct {
	TotalFiles *int
	MovieCount *int
	TVCount    *int
}

// WatchState groups mutable state for watch-mode polling cycles.
type WatchState struct {
	Seen    map[string]bool
	Client  *tmdb.Client
	HasTMDb bool
}

// SuggestTypeInput groups parameters for type-based suggestion generation.
type SuggestTypeInput struct {
	MediaType string
	Database  *db.DB
	Client    *tmdb.Client
	Count     int
}

// BatchMovePreview groups parameters for batch move preview generation.
type BatchMovePreview struct {
	Files     []os.FileInfo
	SourceDir string
	MoviesDir string
	TVDir     string
}

// TrackMoveInput groups parameters for recording a file move operation.
type TrackMoveInput struct {
	SrcPath   string
	DestPath  string
	CleanName string
	FileInfo  os.FileInfo
	Database  *db.DB
	Result    cleaner.Result
}

// FindMoveMediaInput groups parameters for finding or creating media during moves.
type FindMoveMediaInput struct {
	SrcPath  string
	DestPath string
	FileInfo os.FileInfo
	Database *db.DB
	Result   cleaner.Result
}

// WalkEntryInput groups parameters for processing a single walk entry during popout discovery.
type WalkEntryInput struct {
	RootDir string
	Path    string
	Info    os.FileInfo
	Items   *[]popoutItem
	MaxDepth int
}

// FolderRemoveInput groups parameters for folder removal operations.
type FolderRemoveInput struct {
	DirPath string
	DirName string
	BatchID string
	Database *db.DB
}

// MediaRequest groups database context for REST media handlers.
type MediaRequest struct {
	Database *db.DB
	ID       int64
}

// MediaPatchRequest groups parameters for media PATCH REST handlers.
type MediaPatchRequest struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	Database *db.DB
	ID       int64
}

// MediaUpdateField groups parameters for a single media field update.
type MediaUpdateField struct {
	Key      string
	Val      interface{}
	Database *db.DB
	ID       int64
}

// UniqueFilter groups parameters for deduplicating search results.
type UniqueFilter struct {
	ExistingIDs map[int]bool
	Count       int
}

// RecursiveFileContext groups parameters for handling a file during recursive directory walks.
type RecursiveFileContext struct {
	Files   *[]videoFile
	Path    string
	ScanDir string
	Entry   os.DirEntry
	Opts    RecursiveWalkOpts
}

// TrackScanResult groups the result of scanning a single file for action tracking.
type TrackScanResult struct {
	FullPath  string
	InsertErr error
	Media     *db.Media
	MediaID   int64
}

// DiscoverGenreInput groups parameters for genre-based discovery in suggestions.
type DiscoverGenreInput struct {
	Sorted    []genreCount
	MediaType string
	TypeName  string
}

// FillRecoInput groups parameters for recommendation-based suggestion filling.
type FillRecoInput struct {
	Database  *db.DB
	MediaType string
}

// FinalizeScanInput groups parameters for post-scan finalization.
type FinalizeScanInput struct {
	JSONItems []scanJSONItem
	ScanDir   string
	OutputDir string
	Database  *db.DB
	Creds     tmdbCredentials
	Removed   int
	UseJSON   bool
}

// DryRunInput groups parameters for dry-run scan processing.
type DryRunInput struct {
	VideoFiles []videoFile
	UseJSON    bool
	UseTable   bool
}

// DryRunOutput groups mutable output pointers for dry-run scan results.
type DryRunOutput struct {
	JSONItems  *[]scanJSONItem
	TotalFiles *int
	MovieCount *int
	TVCount    *int
}

// RemoveStaleInput groups parameters for stale entry removal during scan.
type RemoveStaleInput struct {
	ExistingMedia []db.Media
	DiskPaths     map[string]bool
	BatchID       string
	Database      *db.DB
	Opts          ScanOutputOpts
}

// ProcessExistingInput groups parameters for processing existing media during scan.
type ProcessExistingInput struct {
	BatchID  string
	Client   *tmdb.Client
	Database *db.DB
	EM       *db.Media
	Opts     ScanOutputOpts
	VF       videoFile
	HasTMDb  bool
}

// HandleRescanInput groups parameters for rescanning a media entry.
type HandleRescanInput struct {
	BatchID  string
	Client   *tmdb.Client
	Database *db.DB
	EM       *db.Media
	Opts     ScanOutputOpts
}

// AppendUniqueInput groups parameters for appending unique search results.
type AppendUniqueInput struct {
	Results []tmdb.SearchResult
	DiscErr error
	Filter  UniqueFilter
}
