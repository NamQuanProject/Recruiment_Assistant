import React from "react";

const Criteria = () => {
  const criteriaData = [
    {
      title: "Technical Skills & Tools (30%)",
      description:
        "Proficiency in data science tools and programming languages (e.g., Python, SQL, Tableau, Excel, SAS). Experience with statistical analysis, data mining, pattern recognition, and predictive modeling. Bonus for machine learning, AI, and big data frameworks.",
    },
    {
      title: "Relevant Experience & Projects (25%)",
      description:
        "Demonstrated experience in data science roles, ideally 7+ years. Evidence of analytical experiments, statistical models, and production-ready solutions. Experience working with unstructured and structured datasets.",
    },
    {
      title: "Problem-Solving & Analytical Abilities (20%)",
      description:
        "Proven track record of solving complex business problems using data. Examples of data-driven decision-making and impactful insights. Ability to interpret data trends and create actionable recommendations.",
    },
    {
      title: "Communication & Stakeholder Collaboration (15%)",
      description:
        "Experience presenting insights and findings to non-technical audiences. Effective communication skills for working with cross-functional teams. Experience translating complex data into business value.",
    },
    {
      title: "Education & Certifications (10%)",
      description:
        "Degree in statistics, mathematics, or related discipline. Relevant professional certifications in data science or analytics. Bonus for advanced degrees (e.g., Master’s, Ph.D.).",
    },
  ];
  return (
    <div className="flex justify-center items-center mt-12 w-5/6 h-[600px] mx-auto bg-primary shadow-md text-color-black p-8 border-[0.5px] border-gray-400 rounded-sm">
      <div className="w-full h-full bg-white">
      <h2 className="text-center font-bold py-4 text-2xl border-b-2 h-20">Suggested Criteria for Assessing a Resume</h2>
      <p className=" p-4 overflow-y-auto">
        {criteriaData.map((item, index) => (
            <div key={index} className="mb-4">
              <h3 className="font-semibold text-lg">{item.title}</h3>
              <p className="text-sm">{item.description}</p>
            </div>
       ))}
      </p>
      </div>
    </div>
  );
};

export default Criteria;