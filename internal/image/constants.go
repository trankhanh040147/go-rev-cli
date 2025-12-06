package image

// Model name for image generation
const (
	DefaultModel = "gemini-3-pro-image-preview"
)

// Supported aspect ratios
const (
	AspectRatio1x1  = "1:1"
	AspectRatio2x3  = "2:3"
	AspectRatio3x2  = "3:2"
	AspectRatio4x3  = "4:3"
	AspectRatio3x4  = "3:4"
	AspectRatio16x9 = "16:9"
	AspectRatio9x16 = "9:16"
)

// Default aspect ratio
const DefaultAspectRatio = AspectRatio1x1

// Supported resolutions
const (
	Resolution1K = "1K"
	Resolution2K = "2K"
	Resolution4K = "4K"
)

// Default resolution
const DefaultResolution = Resolution1K

// Default output filename
const DefaultOutputFile = "image.png"

// ValidAspectRatios contains all valid aspect ratio values
var ValidAspectRatios = []string{
	AspectRatio1x1,
	AspectRatio2x3,
	AspectRatio3x2,
	AspectRatio4x3,
	AspectRatio3x4,
	AspectRatio16x9,
	AspectRatio9x16,
}

// ValidResolutions contains all valid resolution values
var ValidResolutions = []string{
	Resolution1K,
	Resolution2K,
	Resolution4K,
}

// Pricing constants for Gemini 3 Pro Image Preview (per 1M tokens in USD)
// Reference: https://ai.google.dev/gemini-api/docs/pricing#gemini-3-pro-image-preview
const (
	// Input pricing
	InputPricePerMillionTokens = 2.00 // $2.00 per 1M tokens (text/image)

	// Output pricing
	OutputTextPricePerMillionTokens  = 12.00  // $12.00 per 1M tokens (text and thinking)
	OutputImagePricePerMillionTokens = 120.00 // $120.00 per 1M tokens (images)

	// Image-specific fixed pricing (based on resolution)
	ImagePrice1K = 0.134 // $0.134 per 1K/2K image (1120 tokens)
	ImagePrice4K = 0.24  // $0.24 per 4K image (2000 tokens)

	// Token counts for image outputs
	ImageTokens1K = 1120 // Tokens consumed for 1K/2K images
	ImageTokens4K = 2000 // Tokens consumed for 4K images

	// Currency conversion
	// USD to VND exchange rate (approximate, should be updated periodically)
	// Reference: https://www.exchange-rates.org/converter/usd-vnd
	USDToVNDExchangeRate = 25000.0 // 1 USD = ~25,000 VND (approximate rate)
)
