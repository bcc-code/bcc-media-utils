package mediainfoserver

type MediaInfoResult struct {
	CreatingLibrary CreatingLibrary `json:"creatingLibrary,omitempty"`
	Media           Media           `json:"media,omitempty"`
}

type CreatingLibrary struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	URL     string `json:"url,omitempty"`
}
type Extra struct {
	DelaySDTI        string `json:"Delay_SDTI,omitempty"`
	IntraDcPrecision string `json:"intra_dc_precision,omitempty"`
	Locked           string `json:"Locked,omitempty"`
	BlockAlignment   string `json:"BlockAlignment,omitempty"`
}

type Track struct {
	Type                           string `json:"@type,omitempty"`
	VideoCount                     string `json:"VideoCount,omitempty"`
	AudioCount                     string `json:"AudioCount,omitempty"`
	OtherCount                     string `json:"OtherCount,omitempty"`
	FileExtension                  string `json:"FileExtension,omitempty"`
	Format                         string `json:"Format,omitempty"`
	FormatCommercialIfAny          string `json:"Format_Commercial_IfAny,omitempty"`
	FormatVersion                  string `json:"Format_Version,omitempty"`
	FormatProfile                  string `json:"Format_Profile,omitempty"`
	FormatSettings                 string `json:"Format_Settings,omitempty"`
	FileSize                       string `json:"FileSize,omitempty"`
	Duration                       string `json:"Duration,omitempty"`
	OverallBitRate                 string `json:"OverallBitRate,omitempty"`
	FrameRate                      string `json:"FrameRate,omitempty"`
	FrameCount                     string `json:"FrameCount,omitempty"`
	StreamSize                     string `json:"StreamSize,omitempty"`
	FooterSize                     string `json:"FooterSize,omitempty"`
	EncodedDate                    string `json:"Encoded_Date,omitempty"`
	FileModifiedDate               string `json:"File_Modified_Date,omitempty"`
	FileModifiedDateLocal          string `json:"File_Modified_Date_Local,omitempty"`
	EncodedApplicationCompanyName  string `json:"Encoded_Application_CompanyName,omitempty"`
	EncodedApplicationName         string `json:"Encoded_Application_Name,omitempty"`
	EncodedApplicationVersion      string `json:"Encoded_Application_Version,omitempty"`
	EncodedLibraryName             string `json:"Encoded_Library_Name,omitempty"`
	EncodedLibraryVersion          string `json:"Encoded_Library_Version,omitempty"`
	StreamOrder                    string `json:"StreamOrder,omitempty"`
	ID                             string `json:"ID,omitempty"`
	FormatLevel                    string `json:"Format_Level,omitempty"`
	FormatSettingsBVOP             string `json:"Format_Settings_BVOP,omitempty"`
	FormatSettingsMatrix           string `json:"Format_Settings_Matrix,omitempty"`
	FormatSettingsGOP              string `json:"Format_Settings_GOP,omitempty"`
	FormatSettingsPictureStructure string `json:"Format_Settings_PictureStructure,omitempty"`
	FormatSettingsWrapping         string `json:"Format_Settings_Wrapping,omitempty"`
	CodecID                        string `json:"CodecID,omitempty"`
	BitRateMode                    string `json:"BitRate_Mode,omitempty"`
	BitRate                        string `json:"BitRate,omitempty"`
	Width                          string `json:"Width,omitempty"`
	Height                         string `json:"Height,omitempty"`
	SampledWidth                   string `json:"Sampled_Width,omitempty"`
	SampledHeight                  string `json:"Sampled_Height,omitempty"`
	PixelAspectRatio               string `json:"PixelAspectRatio,omitempty"`
	DisplayAspectRatio             string `json:"DisplayAspectRatio,omitempty"`
	FrameRateNum                   string `json:"FrameRate_Num,omitempty"`
	FrameRateDen                   string `json:"FrameRate_Den,omitempty"`
	ColorSpace                     string `json:"ColorSpace,omitempty"`
	ChromaSubsampling              string `json:"ChromaSubsampling,omitempty"`
	BitDepth                       string `json:"BitDepth,omitempty"`
	ScanType                       string `json:"ScanType,omitempty"`
	ScanOrder                      string `json:"ScanOrder,omitempty"`
	CompressionMode                string `json:"Compression_Mode,omitempty"`
	Delay                          string `json:"Delay,omitempty"`
	DelayDropFrame                 string `json:"Delay_DropFrame,omitempty"`
	DelaySource                    string `json:"Delay_Source,omitempty"`
	DelayOriginal                  string `json:"Delay_Original,omitempty"`
	DelayOriginalDropFrame         string `json:"Delay_Original_DropFrame,omitempty"`
	DelayOriginalSource            string `json:"Delay_Original_Source,omitempty"`
	TimeCodeFirstFrame             string `json:"TimeCode_FirstFrame,omitempty"`
	TimeCodeSource                 string `json:"TimeCode_Source,omitempty"`
	GopOpenClosed                  string `json:"Gop_OpenClosed,omitempty"`
	GopOpenClosedFirstFrame        string `json:"Gop_OpenClosed_FirstFrame,omitempty"`
	BufferSize                     string `json:"BufferSize,omitempty"`
	ColourDescriptionPresent       string `json:"colour_description_present,omitempty"`
	ColourDescriptionPresentSource string `json:"colour_description_present_Source,omitempty"`
	ColourPrimaries                string `json:"colour_primaries,omitempty"`
	ColourPrimariesSource          string `json:"colour_primaries_Source,omitempty"`
	TransferCharacteristics        string `json:"transfer_characteristics,omitempty"`
	TransferCharacteristicsSource  string `json:"transfer_characteristics_Source,omitempty"`
	MatrixCoefficients             string `json:"matrix_coefficients,omitempty"`
	MatrixCoefficientsSource       string `json:"matrix_coefficients_Source,omitempty"`
	Extra                          Extra  `json:"extra,omitempty"`
	Typeorder                      string `json:"@typeorder,omitempty"`
	FormatSettingsEndianness       string `json:"Format_Settings_Endianness,omitempty"`
	Channels                       string `json:"Channels,omitempty"`
	SamplesPerFrame                string `json:"SamplesPerFrame,omitempty"`
	SamplingRate                   string `json:"SamplingRate,omitempty"`
	SamplingCount                  string `json:"SamplingCount,omitempty"`
	VideoDelay                     string `json:"Video_Delay,omitempty"`
	SubType                        string `json:"Type,omitempty"`
	TimeCodeLastFrame              string `json:"TimeCode_LastFrame,omitempty"`
	TimeCodeSettings               string `json:"TimeCode_Settings,omitempty"`
	TimeCodeStripped               string `json:"TimeCode_Stripped,omitempty"`
	MuxingMode                     string `json:"MuxingMode,omitempty"`
}
type Media struct {
	Ref   string  `json:"@ref,omitempty"`
	Track []Track `json:"track,omitempty"`
}
