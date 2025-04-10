import { Link, useNavigate } from "react-router-dom";

interface Props {
  rank: number;
  name: string;
  score: number;
  companies: string;
  experience: string;
}

const Display = ({ rank, name, score, experience, companies }: Props) => {
  const navigate = useNavigate(); // Use navigate to programmatically redirect
//	evalID := "20250410_165023"
// 	cvID := "20250410_013723_0065"
// 	question := "List all questions that I asked you please."
  const handleClick = async () => {
    const requestBody = {
      eval_id: `20250410_165023`, // Use rank or another unique identifier for eval_id
    };

    try {
      const response = await fetch("http://localhost:8081/ai/chatbot/init", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(requestBody),
      });

      if (response.ok) {
        const data = await response.json();
        console.log("Chatbot initialized:", data);

        // Navigate to the candidate detail page after the request
        navigate(`/candidate/${rank}`);
      } else {
        console.error("Failed to initialize chatbot. Status:", response.status);
      }
    } catch (error) {
      console.error("Error initializing chatbot:", error);
    }
  };

  return (
    <div
      onClick={handleClick} // Trigger the POST request on click
      className="block cursor-pointer"
    >
      <div className="relative flex flex-row justify-between h-[50px] w-full bg-thirdary shadow-md rounded-md border-[0.5px] border-gray-400 hover:bg-blue-300 transition duration-300 ease-in-out">
        <p className="text-gray-700 text-lg h-full w-1/8 text-center py-3">{rank}</p>
        <div className="flex flex-row justify-between h-full w-full">
          <h3 className="text-lg font-semibold h-full w-1/5 text-start py-3">{name}</h3>
          <p className="text-gray-700 h-full w-1/5 text-center py-3">{companies}</p>
          <p className="text-gray-700 h-full w-1/5 text-center py-3">{experience}</p>
          <p className="text-gray-700 h-full w-1/5 text-center py-3">{score}%</p>
        </div>
      </div>
    </div>
  );
};

export default Display;