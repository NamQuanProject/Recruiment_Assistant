import { useParams, useLocation } from "react-router-dom";
import Navbar from "../components/navbar";
import { useContext, useState, useEffect } from "react";
import { DataContext } from "../components/datacontext";
import ReactMarkdown from "react-markdown";
import Footer from "../components/footer";
import Logo1 from "../assets/Ellipse 1.png";
import Logo2 from "../assets/Ellipse 2.png";  
import Logo3 from "../assets/Ellipse 3.png";
import Logo4 from "../assets/Ellipse 4.png";

const CandidateDetailPage = () => {
  const { rank } = useParams<{ rank: string }>(); // Get the rank from the URL
  const { sharedData } = useContext(DataContext); // Access sharedData from context
  const location = useLocation(); // Access the route state
  const { hlCVData } = location.state || {}; // Extract hlCVData from the state
  // Find the candidate based on the rank
  const candidate = sharedData?.list.find(
    (_, index) => index + 1 === parseInt(rank || "", 10)
  );
  console.log("Candidate data:", sharedData); // Log the candidate data for debugging
  const [criteriaData, setCriteriaData] = useState<any>(null); // State to store fetched JSON data
  const [loading, setLoading] = useState(true); // State to manage loading
  const { criteriaJson } = useContext(DataContext);
  console.log("Criteria JSON from context:", criteriaJson);

  useEffect(() => {
    const fetchCriteriaData = async () => {
      if (candidate?.path_to_evaluation) {
        try {
          const API_URL ="http://apigateway23.onrender.com"; // Use environment variable or default URL
          // const response = await fetch(`http://localhost:8080/${candidate.path_to_evaluation.replace(/\\/g, "/")}`);
          const response = await fetch(`${API_URL}/${candidate.path_to_evaluation.replace(/\\/g, "/")}`);
          // console.log("Fetching criteria data from:", `${API_URL}/${candidate.path_to_evaluation.replace(/\\/g, "/")}`);
          if (response.ok) {
            const data = await response.json();
            setCriteriaData(data); // Store the fetched JSON data
          } else {
            console.error("Failed to fetch criteria data. Status:", response.status);
          }
        } catch (error) {
          console.error("Error fetching criteria data:", error);
        } finally {
          setLoading(false); // Stop loading
        }
      }
    };

    fetchCriteriaData();
  }, [candidate?.path_to_evaluation]);

  if (!candidate) {
    return <p>Candidate not found</p>;
  }

  // Extract cv_id from path_to_evaluation
  const cvId = candidate.path_to_cv.match(/cvs[\\/](.*?)\.pdf/)?.[1];
  console.log("Extracted cvId:", cvId); // Logs "0025" // Log the extracted cv_id for debugging
  console.log("Path:", candidate.path_to_cv)
  // Chatbox state and logic
  const [messages, setMessages] = useState<{ sender: string; text: string }[]>([]);
  const [input, setInput] = useState("");
  const [showInput, setShowInput] = useState("");
  const handleSendMessage = async () => {
    if (input.trim() === "") return;
    // Add user message to chat
    setMessages((prev) => [...prev, { sender: "User", text: input }]);

    try {
      setShowInput("");
      // Prepare the request body
      const requestBody = {
        cv_id: cvId, // Use the extracted cv_id
        question: input, // The user's input as the question
      };
      console.log("Request body:", requestBody); // Log the request body for debugging

      // Send the POST request to the backend
      const AI_URL = "https://aiservice23.onrender.com"; // Use environment variable or default URL
      // const AI_URL = "http://localhost:8081";
      // const response = await fetch("http://localhost:8081/ai/chatbot/ask", {
      //   method: "POST",
      //   headers: { "Content-Type": "application/json" },
      //   body: JSON.stringify(requestBody),
      // });
      const response = await fetch(`${AI_URL}/ai/chatbot/ask`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(requestBody),
      });
      if (response.ok) {
        const data = await response.json();
        console.log("Chatbot response:", data);

        // Add bot response to chat
        setMessages((prev) => [
          ...prev,
          { sender: "Bot", text: data.answer || "No response from chatbot" },
        ]);
      } else {
        console.error("Failed to get chatbot response. Status:", response.status);
        setMessages((prev) => [
          ...prev,
          { sender: "Bot", text: "Failed to get a response from the chatbot." },
        ]);
      }
    } catch (error) {
      console.error("Error communicating with chatbot:", error);
      setMessages((prev) => [
        ...prev,
        { sender: "Bot", text: "An error occurred while communicating with the chatbot." },
      ]);
    }

    // Clear input
    setInput("");
  };

  
  return (
    <section id="candidate-detail" className="mt-32 h-screen mb-20 pb-20">
      <Navbar />

      <div className="flex flex-col justify-between h-[1600px] w-5/6 mx-auto bg-primary border border-gray-300 rounded-md shadow-md py-6 px-4">
        <div className="w-full gap-4 flex justify-between h-7/10 mx-auto">
          <iframe
            // src={`http://localhost:8080/${hlCVData?.highlighted_pdf_path.replace(/\\/g, "/")}`}
             src={`https://apigateway23.onrender.com/${hlCVData?.highlighted_pdf_path.replace(/\\/g, "/")}`}
            // src={`${import.meta.env.VITE_API_URL}/${hlCVData?.highlighted_pdf_path.replace(/\\/g, "/")}`}
            title="Resume"
            className="w-3/5 h-full border object-contain"
          ></iframe>

              <div className="w-2/5 h-full bg-white border-2 overflow-y-auto p-4">
            {loading ? (
              <p>Loading criteria...</p>
            ) : criteriaData && criteriaJson ? (
              <div>
                <h3 className="font-bold text-xl mb-2">I. Final Score</h3>
                <p className="text-lg text-gray-700 mb-4">Based on our evaluation, the final score is {criteriaData.FinalScore.toFixed(2)}</p>
                <h3 className="font-bold text-xl mb-2">II. Authenticity</h3>
                <p className="text-lg text-gray-700 mb-4">Based on provided link, this CV is assessed to have {criteriaData.Authenticity/10*100}% authenticity which indicates the creditability of this CV</p>

              
                {/* Display Personal Info */}
                <h3 className="font-bold text-xl mb-2">III. Personal Information</h3>
                <div className="text-sm text-gray-700 mb-4">
                  <p><strong>Full Name:</strong> {criteriaData.PersonalInfo.FullName}</p>
                  <p><strong>Email:</strong> {criteriaData.PersonalInfo.Email}</p>
                  <p><strong>Phone Number:</strong> {criteriaData.PersonalInfo.PhoneNumber}</p>
                  <p><strong>Address:</strong> {criteriaData.PersonalInfo.Address || "Not provided"}</p>
                  <p><strong>LinkedIn:</strong> <a href={`https://${criteriaData.PersonalInfo.LinkedIn}`} target="_blank" rel="noopener noreferrer" className="text-blue-500 underline">{criteriaData.PersonalInfo.LinkedIn}</a></p>
                  <p><strong>GitHub:</strong> {criteriaData.PersonalInfo.GitHub || "Not provided"}</p>
    
                </div>

                {/* Display Evaluation */}
                <h3 className="font-bold text-xl mb-2">IV. Evaluation Summary</h3>
                <p className="font-semibold text-base">-- MainCategories -- </p>

                {criteriaData.Evaluation.slice(0, criteriaJson?.MainCategory.length).map((item: any, index: number) => (
                  <div key={index} className="mb-4">
                    <h4 className="font-semibold text-base">📌 {index+1 }. {item.category} </h4>
                    <p className="text-sm font-bold pl-6">🌟 Score: {item.score}/{(index < criteriaJson?.MainCategory.length) ? criteriaJson?.MainCategory[0].ScoringScale : criteriaJson?.SubCategory[0].ScoringScale}</p>
                    <p className="text-sm text-gray-700 pl-6">🔍 Explanation: 
                      <br/>
                      {item.explanation}</p>
                      
                  </div>
                ))}
                <p className="font-semibold text-base">-- SubCategories --</p>

                {criteriaData.Evaluation.slice(criteriaJson?.MainCategory.length).map((item: any, index: number) => (
                  <div key={index} className="mb-4">
                    <h4 className="font-semibold text-base">📌 {index+1}. {item.category} </h4>
                    <p className="text-sm font-bold pl-6">🌟 Score: {item.score}/{(index < criteriaJson?.SubCategory.length) ? criteriaJson?.SubCategory[0].ScoringScale : criteriaJson?.SubCategory[0].ScoringScale}</p>
                    <p className="text-sm text-gray-700 pl-6">🔍 Explanation: {item.explanation}</p>
                      
                  </div>
                ))}
                
                </div>
              ) : (
              <p>No criteria data available.</p>
              )}
          </div>
        </div>

        {/* Chatbox */}

        <div className="flex flex-col mx-auto w-full h-[28%]  bg-white border-gray-300 rounded-md shadow-md">
          {/* Chat History */}
          <h3 className="text-lg font-bold text-center py-2 border-b border-gray-300 bg-blue-400">Confused about our Evaluation? Ask the AI Chatbot!</h3>
          <div className="flex-1 overflow-y-auto p-4">
            {messages.map((message, index) => (
              <div
                key={index}
                className={`mb-2 ${
                  message.sender === "User" ? "text-right" : "text-left"
                }`}
              >
                {message.sender === "User" ? (
                  <span className="inline-block px-4 py-2 rounded-lg bg-blue-500 text-white">
                    {message.text}
                  </span>
                ) : (
                  <div className="inline-block px-4 py-2 rounded-lg bg-gray-200 text-black">
                    <ReactMarkdown>{message.text}</ReactMarkdown>
                  </div>
                )}
              </div>
            ))}
          </div>

          {/* Input Area */}
          <div className="flex items-center border-t border-gray-300 p-2">
            <input
              type="text"
              className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none"
              placeholder="Type your message..."
              value={showInput }
              onChange={(e) => {setInput(e.target.value); setShowInput(e.target.value)}}
              onKeyDown={(e) => e.key === "Enter"  && handleSendMessage()}
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
      <img src={Logo1} alt="Logo 1" className="fixed bottom-0 right-0 translate-y-2/5 -z-10" />
      <img src={Logo2} alt="Logo 2" className="fixed top-0 right-0 -translate-x-[200px] -z-10" />
      <img src={Logo3} alt="Logo 3" className="fixed bottom-0 left-0 translate-x-[30px] -z-10" />
      <img src={Logo4} alt="Logo 4" className="fixed top-0 left-0 translate-y-[30px] -z-10" />
      <Footer />
    </section>
  );
};

export default CandidateDetailPage;