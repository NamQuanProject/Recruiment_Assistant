package aiservices

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	r := gin.Default()

	r.GET("/ai", func(c *gin.Context) {
		agent, err := NewAIAgent(Config{}, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create AI agent"})
			return
		}

		prompt := "List 3 popular cookie recipes."

		final_prompt := prompt

		result := agent.CallChatGemini(final_prompt)

		c.JSON(http.StatusOK, gin.H{
			"Question": prompt,
			"Response": result["Response"],
		})
		agent.Close()
	})

	r.GET("/ai/jd_category/", func(c *gin.Context) {
		fmt.Println("Route /ai/category is hit")

		structure, err := ReadJsonStructure("./internal/aiservices/jobs_guideds.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CV"})
			return
		}

		jobData := make(map[string]interface{})

		for jobCategory, accountDataRaw := range structure {
			accountData, ok := accountDataRaw.(map[string]interface{})
			if !ok {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid job data format"})
				return
			}
			jobData[jobCategory] = accountData
		}

		// Respond with all job data
		c.JSON(http.StatusOK, gin.H{
			"Response": jobData,
		})
	})

	r.GET("/ai/parsing", func(c *gin.Context) {

		cv_raw_text :=
			`
		--- Page 1 ---
		Omkar Pathak
		SOFTWARE ENGINEER · FULL STACK PYTHON DEVELOPER
		Pune, Maharashtra, India
		(+91) omkarpathak27@gmail.com (mailto:omkarpathak27@gmail.com) www.omkarpathak.in (http://www.omkarpathak.in) OmkarPathak (https://github.com/OmkarPathak) omkar-pathak-94473811b (https://www.linkedin.com/in/omkar-pathak-94473811b)
		8087996634 | | | |
		“Make the change that you want to see in the world.”
		Experience
		Schlumberger
		Pune, Maharashtra, India
		DATA ENGINEER
		July 2018 - Present
		• Responsible for implementing and managing an end-to-end CI/CD Pipeline with custom validations for Informatica migrations which
		brought migration time to 1.5 hours from 9 hours without any manual intervention
		• Enhancing, auditing and maintaining custom data ingestion framework that ingest around 1TB of data each day to over 70 business
		units
		• Working with L3 developer team to ensure the discussed Scrum PBI’s are delivered on time for data ingestions
		• Planning and Executing QA and Production Release Cycle activities
		Truso
		Pune, Maharashtra, India
		FULL STACK DEVELOPER INTERN
		June 2018 - July 2018
		• Created RESTful apis
		• Tried my hands on Angular 5/6
		• Was responsible for Django backend development
		Propeluss
		Pune, Maharashtra, India
		DATA ENGINEERING INTERN
		October 2017 - January 2018
		• Wrote various automation scripts to scrape data from various websites.
		• Applied Natural Language Processing to articles scraped from the internet to extract different entities in these articles using entity
		extraction algorithms and applying Machine Learning to classify these articles.
		• Also applied KNN with LSA for extracting relevant tags for various startups based on their works.
		GeeksForGeeks
		Pune, Maharashtra, India
		TECHNICAL CONTENT WRITER
		July 2017 - September 2017
		• Published 4 articles for the topics such as Data Structures and Algorithms and Python
		Softtestlab Technologies
		Pune, Maharashtra, India
		WEB DEVELOPER INTERN
		June 2017 - July 2017
		• Was responsible for creating an internal project for the company using PHP and Laravel for testing purposes
		• Worked on a live project for creating closure reports using PHP and Excel
		Projects
		Pyresparser (https://github.com/OmkarPathak/pyresparser)
		API/Python Package
		PERSONAL PROJECT
		July 2019 - Present
		• A simple resume parser used for extracting information from resumes
		• Extract information from thousands of resumes in just a few seconds
		• Author and maintainer of this project
		Garbage Level Monitoring System (https://github.com/OmkarPathak/Garbage-Level-Monitoring-System)
		IoT
		TEAM PROJECT
		October 2017 - May 2018
		• To find a economical and smarter alternative to current garbage problems
		• Users can monitor levels of all garbage bins from a global dashboard provided
		• Was responsible for Django backend development
		NOVEMBER 3, 2019 OMKAR PATHAK · RÉSUMÉ 1

		--- Page 2 ---
		Pygorithm (https://github.com/OmkarPathak/pygorithm)
		API / Python Package
		PERSONAL PROJECT
		July 2017 - Present
		• Author and maintainer of this project
		• An educational library to teach all the major algorithms
		• Got covered in Fosstack, (https://fosstack.com/algorithms-with-python/) FullStackFeed, (https://fullstackfeed.com/pygorithm-a-python-module-for-learning-all-major-algorithms/) Kleiber (https://kleiber.me/blog/2017/08/10/tutorial-decorator-primer/) and Tagged under Hotest Github Project on ITCodeMonkey (https://www.itcodemonkey.com/article/653.html)
		IoT
		Smart Surveillance System using Raspberry Pi and Face Recognition (https://github.com/OmkarPathak/Smart-Surveillance-System-using-Raspberry-Pi)
		PERSONAL PROJECT
		January 2017 - February 2017
		• Face Recognition using OpenCV and Python
		• Raspberry Pi was used as the data server
		• User notified if any suspicious activity detected in real time
		Password Strength Evaluator using Machine Learning (https://github.com/OmkarPathak/Password-Strength-Evaluator-using-Machine-Learning)
		Machine Learning
		PERSONAL PROJECT
		March 2017
		• SVM algorithm used for training and classification
		• Flask framework used
		• Self-generated dataset
		Education
		Marathwada Mitra Mandal’s College of Engineering
		Pune, Maharashtra, India
		B.E. IN COMPUTER ENGINEERING
		2014 - 2018
		• Aggregate 74%
		Skills
		Programming Languages:
		Python, C, PHP, C++, Shell Script
		Frontend Technologies:
		HTML, CSS, JavaScript, Angular 6/7
		Backend Technologies:
		Django, Flask (Python), Laravel (PHP)
		Operating Systems:
		Linux, Unix, Windows
		Databases:
		MySQL, SQLite, MongoDB
		Other:
		Git, NLP, Scikit-Learn, OpenCV, Cloud (GCP, Azure, DigitalOcean)
		Honors & Awards
		developer,
		Top rated Python
		2018 in Pune and Fifth in India at Github (http://git-awards.com/users/omkarpathak) India
		Writer,
		Quora Top
		2018 India
		2017-18’,
		Awarded ‘The Best Outgoing Student Award
		2018 MMCOE, Pune
		Prize,
		Won 2nd
		2018 in an Hackathon organized by MIT-ADT Persona Fest 2018 Pune
		Award,
		Best Paper
		in National Level Conference on “Emerging Trends in Computing , Analytics
		2018 MMCOE, Pune
		and Security - 2018”(NCETCAS-2018)
		Extracurricular Activities
		Contributor in Pune PyCon 2018
		PUNE, MAHARASHTRA, INDIA
		2018
		• Was a part of Website Designing and volunteering committee
		NOVEMBER 3, 2019 OMKAR PATHAK · RÉSUMÉ 2

		--- Page 3 ---
		Mentor at GirlScript Summer of Code 2019
		PUNE, MAHARASHTRA, INDIA
		2019
		• Mentored 4+ teams in various domains
		Organizing head for the National level technical event -
		Innovatus
		PUNE, MAHARASHTRA, INDIA
		2018
		• Organized project competitions
		Workshop on IoT and Python
		MMCOE, PUNE
		10 Jan 2017
		• Conducted a workshop (https://www.omkarpathak.in/2017/01/10/iot-workshop/) for second year students to give them a brief overview about IoT by completing three mini projects and taught
		them basics of Python programming language
		Publications
		Smart Surveillance System using Raspberry Pi and Face
		DOI10.17148/IJARCCE.2017.64117
		Recognition
		Garbage Level Monitoring System
		Interests
		•
		Competitive Programming
		•
		Photography
		•
		Sketching
		•
		Reading/Writing on Quora
		•
		Contributing to Open Source projects
		NOVEMBER 3, 2019 OMKAR PATHAK · RÉSUMÉ 3
		`

		prompt := "Parse the following CV: " + cv_raw_text
		parsed_response, err := GeminiParsingRawCVText(prompt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CV"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Question": prompt,
			"Response": parsed_response,
		})
	})

	r.GET("/ai/jd_category/:job_name", func(c *gin.Context) {
		jobName := c.Param("job_name")
		fmt.Printf("Route /ai/jd_category/%s is hit\n", jobName)

		structure, err := ReadJsonStructure("./internal/aiservices/jobs_guideds.json")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CV"})
			return
		}

		accountDataRaw, exists := structure[jobName]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job category not found"})
			return
		}

		accountData, ok := accountDataRaw.(map[string]interface{})
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid job data format"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Response": accountData,
		})
	})
	r.Run(":8080")
}
