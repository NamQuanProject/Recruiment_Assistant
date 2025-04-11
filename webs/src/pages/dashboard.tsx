import React, { useContext } from "react";
import { DataContext } from "../components/datacontext";
import Display from "../components/personDisplay";
const DashboardPage = () => {
  const { sharedData } = useContext(DataContext); // Access sharedData from context
  console.log("Shared Data in DashboardPage:", sharedData); // Log the sharedData for debugging
  return (
    <section id="leaderboard">
      <h1 className="text-center w-full text-3xl font-bold text-red-400 mt-32">
        Candidates Leaderboard
      </h1>
      <div className="w-5/6 flex flex-row justify-between mx-auto mt-4 bg-secondary rounded-md shadow-md border-[0.5px] border-gray-400">
        <p className="h-full w-1/8 text-xl text-center py-3 font-semibold">Rank</p>
        <div className="flex flex-row justify-between h-full w-full">
          <p className="font-semibold text-xl h-full w-1/5 text-start py-3">Full Name</p>
          <p className="font-semibold text-xl h-full w-1/5 text-center py-3">Worked For</p>
          <p className="font-semibold text-xl h-full w-1/5 text-center py-3">Experience Level</p>
          <p className="font-semibold text-xl h-full w-1/5 text-center py-3">Authenticity</p>
          <p className="font-semibold text-xl h-full w-1/5 text-center py-3">Final Score</p>
        </div>
      </div>
      <div className="flex flex-col gap-2 w-5/6 max-h-screen mx-auto mt-4">
        {sharedData?.list.map((candidate, index) => (
          <Display evalPath={candidate.path_to_evaluation} key={index} rank={index + 1} name={candidate.full_name} score={candidate.final_score} authenticity={candidate.authenticity} companies={candidate.worked_for} experience={candidate.experience_level} />
        ))}
      </div>
    </section>
  );
};

export default DashboardPage;