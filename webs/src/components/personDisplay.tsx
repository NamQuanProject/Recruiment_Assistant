import { Link } from "react-router-dom";

interface Props {
  rank: number;
  name: string;
  score: number;
  companies: string;
  experience: string;
  
}

const Display = ({
  rank,
  name,
  score,
  experience,
  companies,
} : Props) => {
  return (
  <Link to={`/candidate/${rank}`} className="block"> {/* Navigate to the detail page */}
      <div className="relative flex flex-row justify-between h-[50px] w-full bg-thirdary shadow-md rounded-md border-[0.5px] border-gray-400  hover:bg-blue-300 transition duration-300 ease-in-out">
          <p className="text-gray-700 text-lg h-full w-1/8 text-center py-3">{rank}</p>
          <div className="flex flex-row justify between h-full w-full">
          <h3 className="text-lg font-semibold h-full w-1/5 text-start py-3 ">{name}</h3>
          <p className="text-gray-700 h-full w-1/5 text-center py-3">{companies}</p>
          <p className="text-gray-700 h-full w-1/5 text-center py-3">{experience}</p>
          <p className="text-gray-700 h-full w-1/5 text-center py-3">{score}%</p>
          <p className="text-gray-700 h-full w-1/5 text-center py-3">{score}%</p>
          </div>

      </div>
    </Link>
  )
}

export default Display;