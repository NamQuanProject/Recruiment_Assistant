import { useParams, useLocation } from "react-router-dom";
import Navbar from "../components/navbar";
import { useContext, useState } from "react";
import { DataContext } from "../components/datacontext";
import ReactMarkdown from "react-markdown";

const CandidateDetailPage = () => {
  const { rank } = useParams<{ rank: string }>(); // Get the rank from the URL
  const { sharedData } = useContext(DataContext); // Access sharedData from context
  const location = useLocation(); // Access the route state
  const { hlCVData } = location.state || {}; // Extract hlCVData from the state
  // Find the candidate based on the rank
  const candidate = sharedData?.list.find(
    (c, index) => index + 1 === parseInt(rank || "", 10)
  );

  if (!candidate) {
    return <p>Candidate not found</p>;
  }

  // Extract cv_id from path_to_evaluation
  const cvId = candidate.path_to_cv.match(/cvs\\(.*?)\.pdf/)?.[1];
  console.log("Extracted cvId:", cvId); // Logs "0025" // Log the extracted cv_id for debugging

  // Chatbox state and logic
  const [messages, setMessages] = useState<{ sender: string; text: string }[]>([]);
  const [input, setInput] = useState("");

  const handleSendMessage = async () => {
    if (input.trim() === "") return;

    // Add user message to chat
    setMessages((prev) => [...prev, { sender: "User", text: input }]);

    try {
      // Prepare the request body
      const requestBody = {
        cv_id: cvId, // Use the extracted cv_id
        question: input, // The user's input as the question
      };
      console.log("Request body:", requestBody); // Log the request body for debugging

      // Send the POST request to the backend
      const response = await fetch("http://localhost:8081/ai/chatbot/ask", {
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
    <section id="candidate-detail" className="mt-32 mb-20">
      <Navbar />

      <div className="flex flex-col justify-between h-[1200px] w-5/6 mx-auto bg-primary border border-gray-300 rounded-md shadow-md py-6 px-4">
        <div className="w-full gap-4 flex justify-between h-7/10 mx-auto">
          <iframe
             src={`http://localhost:8080/${hlCVData?.highlighted_pdf_path.replace(/\\/g, "/")}`}
            title="Resume"
            className="w-3/5 h-full border object-contain"
          ></iframe>
          <p className="w-2/5 h-full bg-white border-2 overflow-y-auto">Criteria</p>
        </div>

        {/* Chatbox */}
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