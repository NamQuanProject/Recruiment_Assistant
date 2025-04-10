import { useParams } from 'react-router-dom';
import Navbar from '../components/navbar';
import PdfImage from '../assets/pdf.png';
import Chatbox from '../components/chatbox';
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

  //chatbox
  const [messages, setMessages] = useState<{ sender: string; text: string }[]>([]); // Chat history
    const [input, setInput] = useState(""); // User input
  
    const handleSendMessage = () => {
      if (input.trim() === "") return;
  
      // Add user message to chat
      setMessages((prev) => [...prev, { sender: "User", text: input }]);
  
      // Simulate bot response
      setTimeout(() => {
        setMessages((prev) => [
          ...prev,
          { sender: "Bot", text: `You said: "${input}"` }, // Simple bot response
        ]);
      }, 500);
  
      // Clear input
      setInput("");
    };
  return (
    <section id="candidate-detail" className="mt-32 mb-20">
      <Navbar/>

      <div className="flex flex-col justify-between h-[1200px]  w-5/6 mx-auto  bg-primary border border-gray-300 rounded-md shadow-md py-6 px-4">
        <div className="w-full gap-4 flex justify-between h-7/10 mx-auto">
          <iframe
          src={candidate.resume}
          title="Resume"
          className="w-3/5 h-full border object-contain"
          ></iframe>  
          <p className="w-2/5 h-full bg-white border-2 overflow-y-auto">Criteria </p>
        </div>
      
        {/* Feedback */}
        <div className="w-full h-3/10 mt-4 mx-auto bg-white rounded-md shadow-md">
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

        <div className="flex flex-col mt-6 mx-auto w-full h-[500px] bg-white border border-gray-300 rounded-md shadow-md">
        {/* Chat History */}
        <div className="flex-1 overflow-y-auto p-4">
        {messages.map((message, index) => (
          <div
            key={index}
            className={`mb-2 ${
              message.sender === "User" ? "text-right" : "text-left"
            }`}
          >
            <span
              className={`inline-block px-4 py-2 rounded-lg ${
                message.sender === "User"
                  ? "bg-blue-500 text-white"
                  : "bg-gray-200 text-black"
              }`}
            >
              {message.text}
            </span>
          </div>
        ))}
      </div>

         {/* Input Area */}
        <div className="flex items-center border-t border-gray-300 p-2">
        <input
          type="text"
          className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none"
          placeholder="Type your message..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSendMessage()}
        />
          <button
          className="ml-2 px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
          onClick={handleSendMessage}
        >
          Send
         </button>
          </div>
        </div>
      </div>
    </section> 
  );
};

export default CandidateDetailPage;