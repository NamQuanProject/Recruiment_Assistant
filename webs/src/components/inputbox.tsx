import React, { useState, useContext } from "react";
import { useNavigate } from "react-router-dom"; // Import useNavigate for navigation
import { DataContext } from "./datacontext";
import { Loader2 } from 'lucide-react'; 

const IPBox = ({ setCriteriaData, setLoading }: { setCriteriaData: (data: any) => void; setLoading: (loading: boolean) => void; }) => {
  const [jobDescription, setJobDescription] = useState(""); // State for typed job description
  const [uploadedFile, setUploadedFile] = useState<File | null>(null); // State for uploaded file
  const [jobName, setJobName] = useState("");
  const [isUploadMode, setIsUploadMode] = useState(false); // Toggle between typing and uploading
  const [uploadedFileName, setUploadedFileName] = useState<string | null>(null); // State for uploaded file name
  const [uploadedResumes, setUploadedResumes] = useState<File | null>(null);
  const [resumeFileName, setResumeFileName] = useState<string | null>(null); // Store the file name
  const { setCriteriaJson } = useContext(DataContext); // Access setCriteriaJson from context
  const [rankLoading, setRankLoading] = useState(false); // State for loading
  const [submitCvClicked, setSubmitCvClicked] = useState(false); // State for submit button clicked
  const navigate = useNavigate(); // Initialize useNavigate
  const { setSharedData } = useContext(DataContext); // Access setSharedData from context
  const handleNavigateToDashboard = () => {
    navigate("/dashboard"); // Navigate to the /dashboard route
  };
  const handleFileUpload = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files && files.length > 0) {
      const file = files[0];
      if (file.type !== "application/pdf" && file.type !== "text/plain" && file.type !== "application/vnd.openxmlformats-officedocument.wordprocessingml.document") {
        alert("Only PDF, TXT, or DOCX files are allowed!");
        return;
      }
      setUploadedFile(file);
      setUploadedFileName(file.name); // Set the uploaded file name
      setJobDescription(""); // Clear the typed description when switching to upload mode
      console.log("Uploaded file:", file.name);
    }
  };

  const handleJobDescriptionSubmit = async () => {
    try {
      setLoading(true); 
      const formData = new FormData();
  
      // Append job name
      formData.append("job_name", jobName);
  
      // Append the uploaded PDF file
      if (uploadedFile) {
        formData.append("pdf_file", uploadedFile);
      } else {
        alert("Please upload a PDF file.");
        return;
      }
  
      // Send the form data to the backend
      const response = await fetch("http://localhost:8080/submitJD", {
        method: "POST",
        body: formData,
      });
  
      if (response.ok) {
        // Parse the JSON response
        const responseData = await response.json();
  
        console.log("Message:", responseData.message); // Log the success message
        console.log("Parsed JSON Path:", responseData.path); // Log the path to the parsed JSON file
  
        // Fetch the JSON file from the path
        const jsonResponse = await fetch(`http://localhost:8080/${responseData.path.replace(/\\/g, '/')}`);
        if (jsonResponse.ok) {
          const criteriaJson = await jsonResponse.json();
          console.log("Criteria JSON:", criteriaJson); // Log the parsed JSON data
  
          // Pass the parsed JSON data to the parent component
          setCriteriaJson(criteriaJson);
        } else {
          console.error("Failed to fetch the parsed JSON file.");
        }
      } else {
        console.error("Failed to submit job description. Status:", response.status);
        const errorText = await response.text();
        console.error("Error response:", errorText);
      }
    } catch (error) {
      if (error instanceof Error) {
        console.error("Error submitting job description:", error.message);
      } else {
        console.error("Error submitting job description:", error);
      }
    } finally {
      setLoading(false); // Reset loading state
    }
  };
  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    console.log("Selected files:", files); // Log the selected files
    if (files && files.length > 0) {
      const file = files[0];
      const validTypes = ["application/pdf", "application/x-zip-compressed", "application/zip"]; // Allowed file types
       // Log the selected file type
      if (!validTypes.includes(file.type)) {
        alert("Only PDF or ZIP files are allowed!");
        return;
      }

      setUploadedResumes(file); // Store the selected file
      setResumeFileName(file.name); // Store the file name
      console.log("Selected file:", file.name);
    }
  };
  const handleUploadResume = async () => {
    if (!uploadedResumes) {
      alert("Please select a file before submitting.");
      return;
    }
  
    console.log("Uploading file:", uploadedResumes.name);
  
    const formData = new FormData();
    formData.append("file", uploadedResumes);
  
    try {
      const response = await fetch("http://localhost:8080/submitCVs", {
        method: "POST",
        body: formData,
      });
  
      if (response.ok) {
        const responseData = await response.json(); // Parse the initial response
        console.log("Response Data:", responseData);
  
        // Fetch the JSON file from the path provided in the response
        const jsonResponse = await fetch(`http://localhost:8080/${responseData.final_out_path.replace(/\\/g, "/")}`);
        if (jsonResponse.ok) {
          const jsonData = await jsonResponse.json(); // Parse the JSON file
          console.log("Fetched JSON Data:", jsonData);
  
          // Store the JSON data in the shared context
          setSharedData(jsonData);
          navigate("/dashboard"); // Navigate to the DashboardPage
        } else {
          console.error("Failed to fetch the JSON file. Status:", jsonResponse.status);
        }
      } else {
        console.error("Failed to upload resume. Status:", response.status);
        const errorText = await response.text();
        console.error("Error response:", errorText);
      }
    } catch (error) {
      console.error("Error uploading resume:", error);
    }
  };
  return (
  <div className="relative flex flex-col mx-auto mt-40 w-5/6 h-[400px] bg-primary shadow-md border-[0.5px] border-gray-400 rounded-sm ">
    <h1 className="absolute top-0 left-0 -translate-y-10 text-bold text-2xl">REQUIRED INPUT</h1>
      <div className="flex flex-row justify-between h-full">
          <div className="h-full w-1/2 p-4">
          <div className="h-full bg-white flex flex-col justify-between shadow-md px-3 py-2">
            <h2 className="text-center max-h-1/6 bg-white z-10 pb-2">Job Description</h2>
            <div className="flex justify-between max-h-1/8 mb-4">
              <div className="flex h-full justify-center mb-4">
            <button
                className={`px-4 py-2 rounded-l ${!isUploadMode ? "bg-blue-500 text-white" : "bg-gray-300"}`}
                onClick={() => {
                  setIsUploadMode(false);
                  setUploadedFile(null);
                  setUploadedFileName(null); 
                }}
              >
                Type Description
              </button>
              <button
                className={`px-4 py-2 rounded-r ${isUploadMode ? "bg-blue-500 text-white" : "bg-gray-300"}`}
                onClick={() => {
                  setIsUploadMode(true);
                  setJobDescription("");
                }}
              >
                Upload PDF
              </button>
              </div>
              <textarea
              className="w-2/5 h-full border-2 border-gray-500 box-border p-2 mb-4 overflow-hidden"
              placeholder="Enter the Job Name..."
              value={jobName}
              onChange={(e) => setJobName(e.target.value)}
              ></textarea>
            </div>

            {!isUploadMode ? (
              <textarea
                className="flex h-3/5 w-full border-2 border-gray-500 box-border"
                placeholder="Enter your Job Description..."
                value={jobDescription}
                onChange={(e) => setJobDescription(e.target.value)}
              ></textarea>
            ) : (
              <div className="flex flex-col h-3/5 items-center justify-center p-2 border-2 border-gray-500 box-border">
                <label
                  htmlFor="jd-upload"
                  className="mx-auto w-[150px] text-center bg-blue-500 text-white px-4 py-2 rounded cursor-pointer hover:bg-blue-600"
                >
                  Upload PDF
                </label>
                <input
                  id="jd-upload"
                  type="file"
                  accept=".pdf, .txt, .docx"
                  className="hidden"
                  onChange={handleFileUpload}
                />
                <p className={`mt-4 text-sm ${uploadedFile ? "text-gray-600" : "text-transparent"  }`}>Uploaded File: {uploadedFileName}</p>
              </div>
            )}
            <button
              className="mt-4 max-h-1/8 w-full bg-blue-500 text-white px-4 py-2 rounded cursor-pointer hover:bg-blue-600"
              onClick={handleJobDescriptionSubmit}
            >
              Submit Job Description
            </button>
          </div>
        </div>

        <div className="h-full w-1/2 p-4 flex flex-col justify-between">
          <div className="flex flex-col h-full">
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
              accept=".pdf, .zip, .x-zip-compressed, application/zip"
              data-webkitdirectory=""
              multiple
              className="hidden"
              onChange={handleFileSelect}
            />
            <p className={`mt-4 mx-auto text-sm ${resumeFileName ? "text-gray-600" : "text-transparent"  }`}>Uploaded File: {resumeFileName}</p>
            {(!submitCvClicked)? <button className="mt-20 flex flex-row w-1/3 h-[40px] bg-blue-300 mx-auto items-center px-4 py-2 rounded-md shadow-xl hover:shadow-xl hover:bg-blue-500 transition duration-300"
                onClick={
                  handleUploadResume // Set the submit button clicked state
                }>
                <span className="mx-auto w-full text-white">
                    Submit Resumes
                </span>
            </button>

            : <> {!rankLoading ? (
                  <div className="mt-20 flex flex-row w-1/3 h-[40px] bg-blue-300 mx-auto items-center px-4 py-2 rounded-md shadow-xl hover:shadow-xl">
                    <Loader2 className="mx-auto animate-spin w-1/5 h-full text-white" />
                    <span className="mx-auto w-full text-white">
                      Preparing Results...
                    </span>
                  </div>
                ) : (
                  <button  className="mt-20 flex flex-row w-1/3 h-[40px] bg-blue-300 mx-auto items-center px-4 py-2 rounded-md shadow-xl hover:shadow-xl hover:bg-blue-500 transition duration-300"
                  onClick={handleNavigateToDashboard}>
                  <span className="mx-auto w-full text-white">
                    Candidates Rank
                  </span>
                  </button>
                )}
                </>
                }
            </div>
          </div>
          {/* <div className="mx-4 h-2/5 mb-3 bg-white p-4">
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
          </div> */}
        </div>
      </div>
    </div>
  );
};

export default IPBox;