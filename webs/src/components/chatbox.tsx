import  { useState } from "react";

const Chatbox = () => {
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
    <div className="flex flex-col w-1/3 h-[500px] bg-white border border-gray-300 rounded-md shadow-md">
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
  );
};

export default Chatbox;