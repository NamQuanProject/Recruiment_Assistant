import { Link, useNavigate } from "react-router-dom";
import { useContext } from "react";
import { DataContext, SharedData } from "./datacontext";
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
//	evalID := "20250410_165023"
// 	cvID := "20250410_013723_0065"
// 	question := "List all questions that I asked you please."
  const evalId = evalPath.match(/evaluation_(.*?)\//)?.[1];
  const {setSharedData} = useContext(DataContext); // Access setSharedData from context
  const handleClick = async () => {
    const evalId = evalPath.match(/evaluation_(.*?)\//)?.[1]; // Extract evalId from evalPath
  
    try {
      // First POST request to initialize the chatbot
      const initRequestBody = {
        eval_id: evalId, // Use evalId for the request
      };
  
      const initResponse = await fetch("http://localhost:8081/ai/chatbot/init", {
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
  
        const hlCVResponse = await fetch("http://localhost:8080/getHlCV", {
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
          <p className="text-gray-700 h-full w-1/5 text-center py-3">{authenticity}%</p>
          <p className="text-gray-700 h-full w-1/5 text-center py-3">{score}%</p>
        </div>
      </div>
    </div>
  );
};

export default Display;