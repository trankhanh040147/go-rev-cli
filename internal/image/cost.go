package image

import "github.com/trankhanh040147/revcli/internal/gemini"

// CalculateCost calculates the cost for an image generation request
// based on token usage and resolution
func CalculateCost(promptTokens, completionTokens int32, resolution string, hasImageOutput, hasTextOutput bool) *gemini.CostBreakdown {
	if promptTokens < 0 || completionTokens < 0 {
		return &gemini.CostBreakdown{}
	}

	// Calculate input cost: (promptTokens / 1,000,000) * $2.00
	inputCost := (float64(promptTokens) / 1_000_000.0) * InputPricePerMillionTokens

	var outputCost float64

	// Calculate output cost based on what was generated
	if hasImageOutput {
		// For images, use fixed pricing based on resolution
		switch resolution {
		case Resolution4K:
			outputCost = ImagePrice4K
		case Resolution1K, Resolution2K:
			outputCost = ImagePrice1K
		default:
			// Fallback: calculate based on token count if resolution unknown
			outputCost = (float64(completionTokens) / 1_000_000.0) * OutputImagePricePerMillionTokens
		}
	} else if hasTextOutput {
		// For text output: (completionTokens / 1,000,000) * $12.00
		outputCost = (float64(completionTokens) / 1_000_000.0) * OutputTextPricePerMillionTokens
	}

	// If both image and text, we need to handle separately
	// The API typically returns image tokens in completionTokens, but we use fixed pricing for images
	// If there's text output, we'd need to know how many tokens were text vs image
	// For now, we prioritize image pricing if image exists, otherwise use text pricing

	totalCost := inputCost + outputCost

	// Calculate VND costs
	inputCostVND := inputCost * USDToVNDExchangeRate
	outputCostVND := outputCost * USDToVNDExchangeRate
	totalCostVND := totalCost * USDToVNDExchangeRate

	return &gemini.CostBreakdown{
		InputCost:     inputCost,
		OutputCost:    outputCost,
		TotalCost:     totalCost,
		InputCostVND:  inputCostVND,
		OutputCostVND: outputCostVND,
		TotalCostVND:  totalCostVND,
	}
}
