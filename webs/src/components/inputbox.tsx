import React, { useState } from "react";

const IPBox = () => {
  const [showAdditionalInfo, setShowAdditionalInfo] = useState(false);

  const handleCheckboxChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setShowAdditionalInfo(event.target.checked);
  };
  const jobData = {
    description:
      "We are seeking a Software Engineer to design, develop, and maintain software solutions. Responsibilities include writing clean, efficient code, debugging issues, collaborating with cross-functional teams, and optimizing performance. The ideal candidate has experience with C++, Python, or Java, knowledge of data structures and algorithms, and a problem-solving mindset.",
    requirements: [
      "Proficiency in one or more programming languages",
      "Strong understanding of software development principles",
      "Experience with version control (e.g., Git)",
      "Ability to work independently and in a team",
    ],
  };

  // const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
  //   const files = event.target.files;
  //   if (files) {
  //     const fileArray = Array.from(files);
  //     fileArray.forEach((file) => {
  //       console.log(file.name);
  //       console.log(files.length);
  //     });
  //   }
  // };
  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files) {
      console.log("Number of files uploaded:", files.length);
  
      // Log details of each file
      Array.from(files).forEach((file) => {
        console.log("File Name:", file.name);
        console.log("File Type:", file.type);
      });
  
      const formData = new FormData();
      Array.from(files).forEach((file) => {
        formData.append("files", file);
      });
  
      try {
        const response = await fetch("http://localhost:5000/upload", {
          method: "POST",
          body: formData,
        });
  
        if (response.ok) {
          console.log("Files uploaded successfully!");
        } else {
          console.error("Failed to upload files.");
        }
      } catch (error) {
        console.error("Error uploading files:", error);
      }
    }
  };

  return (
    <div className="relative flex flex-col mx-auto mt-40 w-5/6 h-[300px] bg-primary shadow-md border-[0.5px] border-gray-400 rounded-sm ">
      <h1 className="absolute top-0 left-0 -translate-y-10 text-bold text-2xl">REQUIRED INPUT</h1>
      <div className="flex flex-row justify-between h-full">
        <div className="h-full w-1/2">
          <div className="m-8 h-4/5 bg-white shadow-md px-3 py-2">
            <h2 className="text-center bg-white z-10 pb-2 ">Job Description</h2>
            <textarea className="w-full h-4/5 overflow-y-auto border-2" placeholder="Enter your Job Description...">{jobData.description}</textarea>
          </div>
        </div>

        <div className="h-full w-1/2 p-4 flex flex-col justify-between">
          <div className="flex flex-col h-full m-4 pb-2 ">
            <div className="flex flex-col h-full bg-white">
            <h2 className="text-center sticky top-0 bg-white z-10 p-2 border-b-2">Target Resumes</h2>
            <label
              htmlFor="resume-upload"
              className="mt-8 mx-auto w-[200px] text-center bg-blue-500 text-white px-4 py-2 rounded cursor-pointer hover:bg-blue-600"
            >
              Upload Resumes
            </label>
            <input
              id="resume-upload"
              type="file"
              accept=".pdf, .zip, .docx, .txt"
              data-webkitdirectory=""
              multiple
              className="hidden"
              onChange={handleFileUpload}
            />
            </div>
          </div>
          <div className="mx-4 h-2/5 mb-3 bg-white p-4">
            <label className="">
              <input
                type="checkbox"
                className="form-checkbox"
                onChange={handleCheckboxChange}
              />
              <span className="ml-2">
                <a
                  href="https://example.com/preferred-resumes"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-500 underline hover:text-blue-700"
                >
                  Adding your preferred resumes as examples
                </a>
              </span>
              <span className="block ml-5">This would allow us to know more about your preference</span>
            </label>
          </div>
        </div>
      </div>
    </div>
  );
};

export default IPBox;