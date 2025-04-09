import { useParams } from 'react-router-dom';
import Navbar from '../components/navbar';
import PdfImage from '../assets/pdf.png'; 
import { useState } from 'react';// only for testing
// Adjust the path to your PDF file

const CandidateDetailPage = () => {
  const { rank } = useParams<{ rank: string }>(); // Get the rank from the URL

  // Mock data (replace with API call or context if needed)
  const Candidates = [
    { name: "John Doe", score: 85, companies: "Meta, Amazon", experience: "4 years", rank:  1, resume: "/resumes/0001.pdf", },
    { name: "Jane Smith", score: 90, companies: "Meta, Amazon", experience: "4 years", rank: 2, resume: "/resumes/0002.pdf" , },
    { name: "Alice Johnson", score: 78, companies: "Meta, Amazon", experience: "4 years", rank:  3 },
    { name: "Bob Brown", score: 92, companies: "Meta, Amazon", experience: "4 years", rank: 4 },
    { name: "Charlie Davis", score: 88, companies: "Meta, Amazon", experience: "4 years", rank:  5 },
    { name: "David Wilson", score: 95, companies: "Meta, Amazon", experience: "4 years", rank: 6 },
    { name: "Eva Martinez", score: 80, companies: "Meta, Amazon", experience: "4 years", rank: 7 },
    { name: "Frank Garcia", score: 87, companies: "Meta, Amazon", experience: "4 years", rank: 8 },
    { name: "Grace Lee", score: 91, companies: "Meta, Amazon", experience: "4 years", rank: 9 },
    { name: "Henry Walker", score: 84, companies: "Meta, Amazon", experience: "4 years", rank: 10 },
    { name: "Isabella Hall", score: 89, companies: "Meta, Amazon", experience: "4 years", rank: 11 },
  ]

  const candidate = Candidates.find((c) => c.rank === parseInt(rank || '', 10));

  if (!candidate) {
    return <p>Candidate not found</p>;
  }
  const [selectedContent, setSelectedContent] = useState('Overview'); // State to manage selected content
  const Contents = [
    { name: "Overview", content: "Content A" },
    { name: "Unverified Information", content: "Content B" },
    { name: "Suggested Questionnaire", content: "Content C" },
  ]
  return (
    <section id="candidate-detail" className="mt-32 mb-20">
      <Navbar/>
      {/* <iframe
        src={candidate.resume}
        title="Resume"
        className="w-5/6 h-[600px] border"
      ></iframe> */} {/*  to display the PDF resume */}
      <div className="flex flex-col justify-between h-[1200px]  w-2/3 mx-auto  bg-primary border border-gray-300 rounded-md shadow-md py-6">
        <img src={PdfImage} alt="PDF" className="object-contain w-5/6 h-2/3 mx-auto border-2 bg-gray-400" /> {/*  to display the PDF resume */}
        <div className="w-5/6 h-3/10 mx-auto bg-white rounded-md shadow-md">
          <div className="flex justify-between w-full h-1/6  rounded-t-md">
          <button
              className={`w-1/3 h-full  cursor-pointer rounded-tl-md text-black font-normal text-lg border-r-[0.5px] ${
                selectedContent === "Overview" ? "button-active" : "button-primary"
              }`}
              onClick={() => setSelectedContent("Overview")}
            >
              Overview
            </button>
            <button
              className={`w-1/3 h-full cursor-pointer text-black font-normal text-lg ${
                selectedContent === "Unverified Information" ? "button-active" : "button-primary"
              }`}
              onClick={() => setSelectedContent("Unverified Information")}
            >
              Unverified Information
            </button>
            <button
              className={`w-1/3 h-full cursor-pointer rounded-tr-md text-black font-normal text-lg border-l-[0.5px] ${
                selectedContent === "Suggested Questionnaire" ? "button-active" : "button-primary"
              }`}
              onClick={() => setSelectedContent("Suggested Questionnaire")}
            >
              Suggested Questionnaire
            </button>
          </div>
          <p className="text-center text-2xl font-semibold mt-4">
            { Contents.find((content) => content.name === selectedContent)?.content }
          </p>
        </div>
      </div>
    </section> 
  );
};

export default CandidateDetailPage;