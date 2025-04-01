import requests
from bs4 import BeautifulSoup
from tqdm import tqdm
import json
import re
import unicodedata
from fuzzywuzzy import process
import html


def clean_unicode_strings(data):
    if isinstance(data, str):
        import re
    
        pattern = r'\\u([0-9a-fA-F]{4})'
        
        def unicode_replace(match):
            code = int(match.group(1), 16)
            return chr(code)
            
        return re.sub(pattern, unicode_replace, data)
    elif isinstance(data, list):
        return [clean_unicode_strings(item) for item in data]
    elif isinstance(data, dict):
        return {key: clean_unicode_strings(value) for key, value in data.items()}
    return data

source_link = "https://business.linkedin.com/talent-solutions/resources/how-to-hire-guides/"
response = requests.get(source_link)

soup = BeautifulSoup(response.content, "html.parser")

base_url = "https://business.linkedin.com"
how_to_hire_links = [
    base_url + a["href"] if a["href"].startswith("/") else a["href"]
    for a in soup.find_all("a", href=True)
    if "how-to-hire-guides/" in a["href"]
]



def process_each_link(example_link):
    # INTRODUCTION_PART 
    example_response = requests.get(example_link)
    example_soup = BeautifulSoup(example_response.content, "html.parser")
    headline = example_soup.find(class_="banner-headline headline-42")
    subheadline = example_soup.find(class_="banner-subheadline subheadline-17")
    table_section = example_soup.find(class_="table-content")
    td_elements = [td.text.strip() for td in table_section.find_all("td")] if table_section else []
    li_elements = []
    ul_section = example_soup.find(class_="rich-text rich-text-mode-standard rich-text-align-center rich-text__padding--none rich-text__width--narrow")
    if ul_section:
        li_elements = [li.text.strip() for li in ul_section.find_all("li")]
        p_elements = [p.text.strip() for p in ul_section.find_all("p")]

        p_elements = p_elements + li_elements
        p_elements = [item.replace("\xa0", " ") for item in p_elements[1:]]
        p_elements = list(dict.fromkeys(p_elements))
    td_elements = [item for item in td_elements if item]

    # JOB_DESCRIPTION_PART
    jd_link = example_link + "/job-description"
    jd_response = requests.get(jd_link)
    jd_soup = BeautifulSoup(jd_response.content, "html.parser")

    class_name = "rich-text rich-text-mode-standard rich-text-align-center rich-text__padding--none rich-text__width--narrow"
    jd_section = jd_soup.find(class_=class_name)

    categories = ["Objectives of this role", "Responsibilities", "Required skills and qualifications", "Preferred skills and qualifications"]
    job_description = {category: [] for category in categories}
    current_category = None
    for element in jd_section.find_all(["h3", "ul", "h2"]):
        if element.name in ["h3", "h2"]:
            
            best_match, score = process.extractOne(element.get_text(strip=True), categories)
            if score > 70:  
                current_category = best_match
        elif element.name == "ul" and current_category:
            job_description[current_category].extend(
                [li.get_text(strip=True) for li in element.find_all("li")]
            )

    
    # for category, items in job_description.items():
    #     print(f"**{category}:**")
    #     for item in items:
    #         print(f"  - {item}")
    #     print("\n")



    # JOB INTERVIEW PART - STRONG POINTS MARKS - MARKING QUITERIA
    interview_questions_links = example_link + "/interview-questions"
    interview_response = requests.get(interview_questions_links)
    interview_soup = BeautifulSoup(interview_response.content, "html.parser")


    class_name_marking_quiteria = "banner-headline headline-42"

    marking_questions_type = interview_soup.find_all("h2", class_="banner-headline headline-42")

    marking_questions_type = [item.text.strip() for item in marking_questions_type]


    class_name_for_each_questions = "banner-subheadline headline-34"
    class_name_for_each_questions_part = "rich-text rich-text-mode-standard rich-text-align-center rich-text__padding--none rich-text__width--narrow" # Including "Whats to listen  for and why this is matters"
    string = "rich-text rich-text-mode-standard rich-text-align-left rich-text__padding--none rich-text__width--narrow"
    questions = interview_soup.find_all("p", class_=class_name_for_each_questions)
    center_class = "rich-text rich-text-mode-standard rich-text-align-center rich-text__padding--none rich-text__width--narrow"
    left_class = "rich-text rich-text-mode-standard rich-text-align-left rich-text__padding--none rich-text__width--narrow"


    questions_parts = [
        part for part in interview_soup.find_all("div")
        if part.get("class") and (center_class in " ".join(part.get("class")) or left_class in " ".join(part.get("class")))
    ]


    output_data = []

    questions_text = [q.get_text(strip=True) for q in questions]

    if questions_text and questions_text[-1] == "Related job titles":
        questions.pop()

    
    for i, question in enumerate(questions):
        question_entry = {
            "question": question.get_text(strip=True),
            "why_this_matters": "",
            "what_to_listen_for": []
        }

        if i * 2 < len(questions_parts):
            question_entry["why_this_matters"] = questions_parts[i * 2].get_text(strip=True)


            question_entry["why_this_matters"] = question_entry["why_this_matters"].replace("Why this matters:", "").strip()

            listen_for_section = questions_parts[i * 2 + 1]
            listen_for_title = listen_for_section.find("strong")

            if listen_for_title and "What to listen for" in listen_for_title.get_text():
                listen_for_points = listen_for_section.find_all("li")
                question_entry["what_to_listen_for"] = [
                    point.get_text(strip=True) for point in listen_for_points
                ]

        output_data.append(question_entry)



    

    job_name = example_link.split("/")[-1]

    job_name = clean_unicode_strings(job_name)
    description = clean_unicode_strings(subheadline.text if subheadline else "Not found")
    purpose = [clean_unicode_strings(p) for p in p_elements]
    skills_requirements = [clean_unicode_strings(td) for td in td_elements]
    jd = clean_unicode_strings(job_description)
    interview_ques = [clean_unicode_strings(q) for q in output_data]

    each_category = {
        "job_name": job_name,
        "description": description,
        "purpose": purpose,
        "skills_requirements": skills_requirements,
        "job_description": jd,
        "interview_questions": interview_ques
    }
    return each_category


jobs_guideds = []
for example_link in tqdm(how_to_hire_links):
    result = process_each_link(example_link)
    jobs_guideds.append(result)

with open("jobs_guideds.json", "w") as f:
    json.dump(jobs_guideds, f, indent=4, ensure_ascii=False)