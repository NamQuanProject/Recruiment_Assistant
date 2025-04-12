import React, { useState, useContext } from "react";
import { useNavigate } from "react-router-dom"; // Import useNavigate for navigation
import { DataContext } from "./datacontext";
import { Loader2 } from 'lucide-react'; 
import Select from "react-select"
import jsPDF from "jspdf";
const jobNames = [
  {'value': 'account-executive', 'label': 'Account Executive'}, {'value': 'account-manager', 'label': 'Account Manager'}, {'value': 'accountants', 'label': 'Accountants'}, {'value': 'accounting-professional', 'label': 'Accounting Professional'}, {'value': 'accounting-clerk', 'label': 'Accounting Clerk'}, {'value': 'application-developer', 'label': 'Application Developer'}, {'value': 'architect', 'label': 'Architect'}, {'value': 'art-director', 'label': 'Art Director'}, {'value': 'artificial-intelligence-engineer', 'label': 'Artificial Intelligence Engineer'}, {'value': 'assistant-manager', 'label': 'Assistant Manager'}, {'value': 'auto-mechanic', 'label': 'Auto Mechanic'}, {'value': 'back-end-developer', 'label': 'Back End Developer'}, {'value': 'travel-agent', 'label': 'Travel Agent'}, {'value': 'treasurer', 'label': 'Treasurer'}, {'value': 'underwriter', 'label': 'Underwriter'}, {'value': 'unity-developer', 'label': 'Unity Developer'}, {'value': 'ux-designer', 'label': 'Ux Designer'}, {'value': 'veterinarian', 'label': 'Veterinarian'}, {'value': 'veterinary-assistant', 'label': 'Veterinary Assistant'}, {'value': 'video-editor', 'label': 'Video Editor'}, {'value': 'warehouse-associate', 'label': 'Warehouse Associate'}, {'value': 'warehouse-manager', 'label': 'Warehouse Manager'}, {'value': 'web-developer', 'label': 'Web Developer'}, {'value': 'welder', 'label': 'Welder'}
]
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
  const { setSharedData } = useContext(DataContext); 
  
  const handleNavigateToDashboard = () => {
    navigate("/dashboard"); // Navigate to the /dashboard route
  };

  const customStyles = {
    menuList: (provided:any) => ({
      ...provided,
      maxHeight: 200, // limit height
      overflowY: 'auto', // make scrollable
    }),
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
      if (jobDescription.trim() !== "") {
        // Convert the job description to a PDF file
        const pdf = new jsPDF();
        pdf.text("Job Description", 10, 10); // Add a title
        pdf.text(jobDescription, 10, 20); // Add the job description text
        const pdfBlob = pdf.output("blob"); // Convert the PDF to a Blob
  
        // Append the PDF Blob to the FormData
        formData.append("pdf_file", pdfBlob, "job_description.pdf");
      } else if (uploadedFile) {
        // Append the uploaded PDF file if no job description is typed
        formData.append("pdf_file", uploadedFile);
      } else {
        alert("Please provide a job description or upload a file.");
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
    setSubmitCvClicked(true); // Set the submit button clicked state
    try {
      setRankLoading(true); // Set loading state
      const response = await fetch("http://localhost:8080/submitCVs", {
        method: "POST",
        body: formData,
      });
      console.log(rankLoading) ;
      if (response.ok) {
        const responseData = await response.json(); // Parse the initial response
        console.log("Response Data:", responseData);
        console.log("cc");
        // Fetch the JSON file from the path provided in the response
        const jsonResponse = await fetch(`http://localhost:8080/${responseData.final_out_path.replace(/\\/g, "/")}`);
        if (jsonResponse.ok) {
          const jsonData = await jsonResponse.json(); // Parse the JSON file
          console.log("Fetched JSON Data:", jsonData);
          
          // Store the JSON data in the shared context
          setSharedData(jsonData);
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
    finally {
      setRankLoading(false); // Reset loading state
    }
  };
  return (
  <div className="relative flex flex-col mx-auto mt-40 w-5/6 h-[400px] bg-primary shadow-md border-[0.5px] border-gray-400 rounded-sm ">
    <h1 className="absolute top-0 left-0 -translate-y-10 font-bold text-2xl text-gray-800">REQUIRED INPUT</h1>
      <div className="flex flex-row justify-between h-full">
          <div className="h-full w-1/2 p-4">
          <div className="h-full bg-white flex flex-col justify-between border-[0.1px] border-gray-400 shadow-md px-3 py-2 rounded-sm">
            <h2 className="text-center text-xl font-semibold max-h-1/6 bg-white z-10 pb-2">Job Description</h2>
            <div className="flex justify-between max-h-1/8 mb-4">
              <div className="flex h-full justify-center mb-4 r">
            <button
                className={`px-4 py-2 rounded-l border-[1.5px] border-black cursor-pointer ${!isUploadMode ? "button-primary text-white" : "button-active"}`}
                onClick={() => {
                  setIsUploadMode(false);
                  setUploadedFile(null);
                  setUploadedFileName(null); 
                }}
              >
                Type Description
              </button>
              <button
                className={`px-4 py-2 rounded-r border-t-[1.5px] border-b-[1.5px] border-r-[1.5px] cursor-pointer border-black ${isUploadMode ? "button-primary text-white" : "button-active"}`}
                onClick={() => {
                  setIsUploadMode(true);
                  setJobDescription("");
                }}
              >
                Upload PDF
              </button>
              </div>
              <Select
              options={jobNames}
              placeholder="Job Name"
              styles={customStyles}
              classNames={{
                control: (state) =>
                  state.isFocused ? 'h-full' : 'w-full h-full',
              }}
              className="w-2/5 h-full border-[1.5px] rounded-sm "
              value={jobNames.find((option) => option.value === jobName) || null}
              onChange={(e) => setJobName(e ? e.value : "")}
              ></Select>
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
                  className="mx-auto w-[150px] text-center px-4 py-2  bg-blue-400 text-white shadow-xl rounded cursor-pointer hover:bg-blue-600 hover:scale-105 transition duration-300"
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
              className="mt-4 max-h-1/8 w-full bg-blue-500 text-white px-4 py-2 shadow-xl rounded cursor-pointer hover:bg-blue-600 hover:tracking-wide transition duration-300"
              onClick={handleJobDescriptionSubmit}
            >
              Submit Job Description
            </button>
          </div>
        </div>

        <div className="h-full w-1/2 p-4">
            <div className="flex flex-col justify-between h-full bg-white shadow-md border-[0.1px] border-gray-400 px-3 py-2 rounded-sm">
            <h2 className="text-center text-xl font-semibold max-h-1/6 bg-white z-10 pb-2 ">Candidate Resumes</h2>

            <div className="flex flex-col h-3/5 items-center justify-center p-2 border-2 border-gray-500 box-border">
            <label
              htmlFor="resume-upload"
              className="mt-8 mx-auto w-[200px] text-center bg-blue-400  text-white px-4 py-2 rounded cursor-pointer shadow-xl hover:bg-blue-600 hover:scale-105 transition duration-300"
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
            </div>
            {(!submitCvClicked)? <button className="my-auto flex flex-row w-1/3 cursor-pointer h-[40px] bg-blue-500 mx-auto items-center px-4 py-2 rounded-md shadow-xl hover:bg-blue-600 hover:scale-105 transition duration-300"
                onClick={
                  handleUploadResume // Set the submit button clicked state
                }>
                <span className="mx-auto w-full text-white">
                    Submit Resumes
                </span>
            </button>

            : <> {rankLoading ? (
                  <div className="my-auto flex flex-row max-w-1/2 gap-2 h-[40px] bg-blue-600 opacity-40 mx-auto items-center px-4 py-2 rounded-md shadow-xl hover:shadow-xl">
                    <Loader2 className="mx-auto animate-spin max-w-1/5 h-full text-white" />
                    <span className="mx-auto w-full text-white">
                      Preparing Results...
                    </span>
                  </div>
                ) : (
                  <button  className="my-auto flex flex-row w-1/3 h-[40px] cursor-pointer bg-blue-500 mx-auto items-center px-4 py-2 rounded-md shadow-xl hover:scale-105 hover:shadow-xl hover:bg-blue-500 transition duration-300"
                  onClick={handleNavigateToDashboard}>
                  <span className="mx-auto w-full text-white">
                  👉 Candidates Rank 
                  </span>
                  </button>
                )}
                </>
                }
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