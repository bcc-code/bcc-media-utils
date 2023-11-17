package cantemo

import "time"

type Item struct {
	ItemType        string `json:"item_type"`
	MetadataSummary struct {
		ArchiveStatus string      `json:"archive_status"`
		Duration      string      `json:"duration"`
		StartTimecode string      `json:"start_timecode"`
		User          string      `json:"user"`
		Added         string      `json:"added"`
		Type          string      `json:"type"`
		Filename      string      `json:"filename"`
		Format        string      `json:"format"`
		Dimension     string      `json:"dimension"`
		UserFullName  interface{} `json:"user_full_name"`
	} `json:"metadata_summary"`
	MediaDetails struct {
		StartTimecode struct {
			Vidispine string `json:"vidispine"`
			Dropframe bool   `json:"dropframe"`
		} `json:"start_timecode"`
	} `json:"media_details"`
	Id             string `json:"id"`
	SystemMetadata struct {
		StartTimeCode      string    `json:"startTimeCode"`
		OriginalAudioCodec string    `json:"originalAudioCodec"`
		DurationSeconds    string    `json:"durationSeconds"`
		StartSeconds       string    `json:"startSeconds"`
		OriginalWidth      string    `json:"originalWidth"`
		DurationTimeCode   string    `json:"durationTimeCode"`
		MimeType           string    `json:"mimeType"`
		OriginalHeight     string    `json:"originalHeight"`
		MediaType          string    `json:"mediaType"`
		OriginalFilename   string    `json:"originalFilename"`
		Title              string    `json:"title"`
		Created            time.Time `json:"created"`
		ShapeTag           string    `json:"shapeTag"`
		OriginalVideoCodec string    `json:"originalVideoCodec"`
		OriginalFormat     string    `json:"originalFormat"`
		User               string    `json:"user"`
	} `json:"system_metadata"`
	Archived bool `json:"archived"`
	Online   bool `json:"online"`
	Previews struct {
		Shapes []struct {
			Uri              string `json:"uri"`
			Closed           bool   `json:"closed"`
			Growing          bool   `json:"growing"`
			Displayname      string `json:"displayname"`
			Id               string `json:"id"`
			DurationTimecode struct {
				Vidispine string `json:"vidispine"`
				Dropframe bool   `json:"dropframe"`
			} `json:"duration_timecode"`
		} `json:"shapes"`
		AudioTracks []struct {
			Name     string `json:"name"`
			Uri      string `json:"uri"`
			Language string `json:"language"`
		} `json:"audio_tracks"`
	} `json:"previews"`
	Subtitles []struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"subtitles"`
	Locked           bool        `json:"locked"`
	AutoPurgeSeconds interface{} `json:"auto_purge_seconds"`
	OriginalShape    struct {
		Id          string `json:"id"`
		DisplayName string `json:"display_name"`
		Closed      bool   `json:"closed"`
	} `json:"original_shape"`
	Permissions struct {
		GenericReadPermission   bool `json:"generic_read_permission"`
		GenericWritePermission  bool `json:"generic_write_permission"`
		UriReadPermission       bool `json:"uri_read_permission"`
		UriWritePermission      bool `json:"uri_write_permission"`
		MetadataReadPermission  bool `json:"metadata_read_permission"`
		MetadataWritePermission bool `json:"metadata_write_permission"`
	} `json:"permissions"`
	LatestVersion int `json:"latest_version"`
}

type ItemMetadata struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	GroupName string `json:"group_name"`
	Metadata  struct {
		Id     int `json:"id"`
		Fields []struct {
			Name            string      `json:"name"`
			Created         time.Time   `json:"created"`
			Updated         time.Time   `json:"updated"`
			FieldType       string      `json:"field_type"`
			Position        *int        `json:"position"`
			Label           string      `json:"label"`
			Language        string      `json:"language,omitempty"`
			ReadAccess      bool        `json:"read_access"`
			WriteAccess     bool        `json:"write_access"`
			Pattern         string      `json:"pattern,omitempty"`
			ExternalId      string      `json:"external_id"`
			Default         string      `json:"default,omitempty"`
			Description     *string     `json:"description,omitempty"`
			Sortable        bool        `json:"sortable,omitempty"`
			IsSystemfield   bool        `json:"is_systemfield"`
			IsSingleton     bool        `json:"is_singleton"`
			SystemfieldName string      `json:"systemfield_name"`
			Reusable        bool        `json:"reusable,omitempty"`
			Readonly        bool        `json:"readonly"`
			VerticalMode    bool        `json:"vertical_mode"`
			Autoset         bool        `json:"autoset"`
			MinInclusive    float64     `json:"min_inclusive,omitempty"`
			MaxInclusive    float64     `json:"max_inclusive,omitempty"`
			Required        bool        `json:"required"`
			Url             string      `json:"url"`
			IsConditional   bool        `json:"is_conditional,omitempty"`
			Value           *string     `json:"value,omitempty"`
			Hideifnotset    interface{} `json:"hideifnotset,omitempty"`
			Representative  interface{} `json:"representative,omitempty"`
			ChoicesUrl      string      `json:"choices_url,omitempty"`
			MinOccurance    interface{} `json:"min_occurance"`
			MaxOccurance    int         `json:"max_occurance,omitempty"`
			Choices         []struct {
				Id         int       `json:"id"`
				Created    time.Time `json:"created"`
				Updated    time.Time `json:"updated"`
				Key        string    `json:"key"`
				Value      string    `json:"value"`
				Position   int       `json:"position"`
				Definition int       `json:"definition"`
			} `json:"choices,omitempty"`
			IsDropframe bool `json:"is_dropframe,omitempty"`
			Fields      []struct {
				Name            string      `json:"name"`
				Created         time.Time   `json:"created"`
				Updated         time.Time   `json:"updated"`
				FieldType       string      `json:"field_type"`
				Position        *int        `json:"position"`
				Label           string      `json:"label"`
				Language        string      `json:"language"`
				ReadAccess      bool        `json:"read_access"`
				WriteAccess     bool        `json:"write_access"`
				Pattern         string      `json:"pattern"`
				ExternalId      string      `json:"external_id"`
				Default         string      `json:"default"`
				Description     string      `json:"description"`
				Sortable        bool        `json:"sortable"`
				IsSystemfield   bool        `json:"is_systemfield"`
				IsSingleton     bool        `json:"is_singleton"`
				SystemfieldName string      `json:"systemfield_name"`
				Reusable        bool        `json:"reusable"`
				Readonly        bool        `json:"readonly"`
				VerticalMode    bool        `json:"vertical_mode"`
				Autoset         bool        `json:"autoset"`
				MinInclusive    float64     `json:"min_inclusive"`
				MaxInclusive    float64     `json:"max_inclusive"`
				Required        bool        `json:"required"`
				Url             string      `json:"url"`
				ChoicesUrl      string      `json:"choices_url,omitempty"`
				MinOccurance    interface{} `json:"min_occurance"`
				MaxOccurance    int         `json:"max_occurance,omitempty"`
				IsConditional   bool        `json:"is_conditional"`
				Value           *string     `json:"value"`
				Hideifnotset    interface{} `json:"hideifnotset"`
				Representative  interface{} `json:"representative"`
				Choices         []struct {
					Id         int       `json:"id"`
					Created    time.Time `json:"created"`
					Updated    time.Time `json:"updated"`
					Key        string    `json:"key"`
					Value      string    `json:"value"`
					Position   int       `json:"position"`
					Definition int       `json:"definition"`
				} `json:"choices,omitempty"`
				IsDropframe bool `json:"is_dropframe,omitempty"`
			} `json:"fields,omitempty"`
			Reference     bool   `json:"reference,omitempty"`
			Collapsed     bool   `json:"collapsed,omitempty"`
			DisplayName   string `json:"display_name,omitempty"`
			VidispineName string `json:"vidispine_name,omitempty"`
		} `json:"fields"`
		Created    time.Time `json:"created"`
		Updated    time.Time `json:"updated"`
		Name       string    `json:"name"`
		ExternalId string    `json:"external_id"`
		Access     int       `json:"access"`
	} `json:"metadata"`
}

type SearchResult struct {
	Results []struct {
		VidispineId         string      `json:"vidispine_id"`
		PortalArchiveStatus []string    `json:"portal_archive_status,omitempty"`
		DurationSeconds     []string    `json:"durationSeconds"`
		Created             []time.Time `json:"created"`
		PortalItemtype      []string    `json:"portal_itemtype"`
		IsOnline            bool        `json:"is_online"`
		Type                string      `json:"type"`
		Title               []string    `json:"title"`
		HasAnnotations      bool        `json:"has_annotations"`
		LockExists          bool        `json:"lock_exists"`
		Type1               string      `json:"_type"`
	} `json:"results"`
	Aggregations struct {
		IsArchived struct {
			Buckets []struct {
				Key      bool `json:"key"`
				DocCount int  `json:"doc_count"`
			} `json:"buckets"`
		} `json:"is_archived"`
		PortalItemtype struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"portal_itemtype"`
		IsOnline struct {
			Buckets []struct {
				Key      bool `json:"key"`
				DocCount int  `json:"doc_count"`
			} `json:"buckets"`
		} `json:"is_online"`
		CurrentFacetFilters struct {
			Meta struct {
			} `json:"meta"`
			IsArchived struct {
				Buckets []struct {
					Key      bool `json:"key"`
					DocCount int  `json:"doc_count"`
				} `json:"buckets"`
			} `json:"is_archived"`
			PortalItemtype struct {
				Buckets []struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
				} `json:"buckets"`
			} `json:"portal_itemtype"`
			IsOnline struct {
				Buckets []struct {
					Key      bool `json:"key"`
					DocCount int  `json:"doc_count"`
				} `json:"buckets"`
			} `json:"is_online"`
			Type struct {
				Buckets []struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
				} `json:"buckets"`
			} `json:"type"`
			User struct {
				Buckets []struct {
					Key      string `json:"key"`
					DocCount int    `json:"doc_count"`
				} `json:"buckets"`
			} `json:"user"`
		} `json:"current_facet_filters"`
		Type struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"type"`
		User struct {
			Buckets []struct {
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"user"`
	} `json:"aggregations"`
	HasNext        bool `json:"has_next"`
	HasPrevious    bool `json:"has_previous"`
	HasOtherPages  bool `json:"has_other_pages"`
	Next           int  `json:"next"`
	Previous       int  `json:"previous"`
	Hits           int  `json:"hits"`
	FirstOnPage    int  `json:"first_on_page"`
	LastOnPage     int  `json:"last_on_page"`
	ResultsPerPage int  `json:"results_per_page"`
	Page           int  `json:"page"`
	Pages          int  `json:"pages"`
}
