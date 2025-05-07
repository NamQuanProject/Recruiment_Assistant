import { useNavigate } from "react-router-dom";
import { useState } from "react";

interface Props {
  rank: number;
  name: string;
  score: number;
  authenticity: number;
  companies: string;
  experience: string;
  evalPath: string;
}

const Display = ({ rank, name, score, experience, authenticity, companies, evalPath }: Props) => {
  const navigate = useNavigate(); // Use navigate to programmatically redirect
  const [isLoading, setIsLoading] = useState(false); // State to track loading

  const handleClick = async () => {
    if (isLoading) return; // Prevent multiple clicks while loading

    const evalId = evalPath.match(/evaluation_(.*?)\//)?.[1]; // Extract evalId from evalPath

    try {
      setIsLoading(true); // Set loading to true

      // First POST request to initialize the chatbot
      const initRequestBody = {
        eval_id: evalId, // Use evalId for the request
      };

      const AI_URL = "https://aiservice23.onrender.com"; // Use environment variable or default URL
      // const AI_URL = "http://localhost:8081";
      // const initResponse = await fetch("http://localhost:8081/ai/chatbot/init", {
      //   method: "POST",
      //   headers: { "Content-Type": "application/json" },
      //   body: JSON.stringify(initRequestBody),
      // });
      const initResponse = await fetch(`${AI_URL}/ai/chatbot/init`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(initRequestBody),
      });

      if (initResponse.ok) {
        const initData = await initResponse.json();
        console.log("Chatbot initialized:", initData);

        // Second POST request to fetch highlighted CV
        const hlCVRequestBody = {
          index: rank-1, // Use rank as the index for the request
        };

        // const hlCVResponse = await fetch("http://localhost:8080/getHlCV", {
        //   method: "POST",
        //   headers: { "Content-Type": "application/json" },
        //   body: JSON.stringify(hlCVRequestBody),
        // });
        const API_URL = "https://apigateway23.onrender.com"; // Use environment variable or default URL
        // const API_URL = "http://localhost:8080";
        const hlCVResponse = await fetch(`${API_URL}/getHlCV`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(hlCVRequestBody),
        });

        if (hlCVResponse.ok) {
          const hlCVData = await hlCVResponse.json();
          console.log("Highlighted CV data:", hlCVData);

          // Navigate to the candidate detail page after both requests
          navigate(`/candidate/${rank}`, { state: { hlCVData } });
        } else {
          console.error("Failed to fetch highlighted CV. Status:", hlCVResponse.status);
        }
      } else {
        console.error("Failed to initialize chatbot. Status:", initResponse.status);
      }
    } catch (error) {
      console.error("Error during handleClick execution:", error);
    } finally {
      setIsLoading(false); // Set loading to false
    }
  };

  return (
    <div
      onClick={handleClick} // Trigger the POST request on click
      className={`block ${isLoading ? "cursor-progress" : "cursor-pointer"}`}
    >
      {isLoading ? (
        <div className="flex cursor-progress justify-center items-center h-[50px] w-full bg-gray-200">
          <p className="text-gray-700 text-lg">Preparing detailed feedback...</p>
        </div>
      ) : (
        <div className="relative flex flex-row justify-between h-[50px] w-full bg-thirdary shadow-md rounded-md border-[0.5px] border-gray-400 hover:bg-blue-300 transition duration-300 ease-in-out">
          <p className="text-gray-700 text-lg h-full w-1/8 text-center py-3">{rank}</p>
          <div className="flex flex-row justify-between h-full w-full">
            <h3 className="text-lg font-semibold h-full w-1/5 text-start py-3">
              {name.length > 20 ? name.slice(0, 17) + "..." : name}
            </h3>
            <p className="text-gray-700 h-full w-1/5 text-center py-3">
              {(companies === "") ? "No companies" : `${companies.length > 30 ? companies.slice(0, 25) + "..." : companies}`}
            </p>
            <p className="text-gray-700 h-full w-1/5 text-center py-3">{(experience == "0" || experience == "" || experience.toLowerCase().includes("year")) ? "Entry level" : `${experience} years`}</p>
            <div className="h-full w-1/5 text-center py-3 flex flex-col justify-center items-center">
              <div className="w-3/5 bg-gray-300 rounded-full h-2">
                <div
                  className={`h-2 rounded-full ${
                    (authenticity / 10) * 100 <= 33
                      ? "bg-red-400"
                      : (authenticity / 10) * 100 <= 66
                      ? "bg-yellow-400"
                      : "bg-green-400"
                  }`}
                  style={{ width: `${(authenticity / 10) * 100}%` }} // Set width dynamically
                ></div>
              </div>
              <p className="text-sm text-gray-700 mt-1 w-full">{(authenticity / 10 * 100).toFixed(1)}%</p>
            </div>
            <p className="text-gray-700 font-semibold text-lg h-full w-1/5 text-center py-3">{(score).toFixed(1)}</p>
          </div>
        </div>
      )}
    </div>
  );
};

export default Display;