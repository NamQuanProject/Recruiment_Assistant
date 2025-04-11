package aiservices

import (
	"fmt"
)

func GetCVAnalysisPrompt(jobTitle string, jobDetails string, textBlocksStr string, evaluationReference map[string]any) string {
	return fmt.Sprintf(`Analyze the CV for the job title "%s" with the following job details: "%s".

	Here are the text blocks extracted from the CV:
	%s
	
	Here is the evaluation reference from another service:
	%s
	
	Identify both strong and weak areas in the CV that match or could be improved to better match the job requirements. For each area, provide:
	1. The exact text from the CV (must match one of the text blocks above)
	2. The page number where it appears
	3. The x, y coordinates and dimensions of the text on the page (must match the position of the text block above and try to generarte coordinates so that they are not overlapping and are in the same area, also try to generate coordinates that make the highligted cover full area of the text block and don't cover other words or any special characters)
	4. A description of why this area is strong or weak and how it could be improved (if weak)
	5. The type of the area ("strong" or "weak")

	Return the results in the following JSON format:
	{
	"areas": [
		{
		"text": "Area text",
		"page": 1,
		"x": 100,
		"y": 200,
		"width": 200,
		"height": 50,
		"description": "This area is strong/weak because...",
		"type": "strong" or "weak"
		}
	]
	}

	Make sure to use the exact text and coordinates from the text blocks provided.
	Identify at most 5 strong areas and 5 weak areas if possible.
	For strong areas, focus on relevant skills, experience, and achievements that match the job requirements.
	For weak areas, focus on the roles (whether it fit with the job description or the job), projects (related or not), and experience (enough or not), provide decent constructive feedback on how to improve them.`, jobTitle, jobDetails, textBlocksStr, evaluationReference)
}
