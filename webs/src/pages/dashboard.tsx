import  { useContext } from "react";
import { DataContext } from "../components/datacontext";
import Display from "../components/personDisplay";
import NavBar from "../components/navbar";
import Footer from "../components/footer";
import Logo1 from '../assets/Ellipse 1.png';
import Logo2 from '../assets/Ellipse 2.png';
import Logo3 from '../assets/Ellipse 3.png';
import Logo4 from '../assets/Ellipse 4.png';
const DashboardPage = () => {
  const { sharedData } = useContext(DataContext); // Access sharedData from context
  console.log("Shared Data in DashboardPage:", sharedData); // Log the sharedData for debugging
  return (
    <section id="leaderboard" className="min-h-screen relative">
      <NavBar />
      <h1 className="text-center w-full text-3xl font-bold text-red-400 mt-32">
        Candidates Leaderboard
      </h1>
      <div className="w-5/6 flex flex-row text-gray-700 justify-between mx-auto mt-4 bg-secondary rounded-md shadow-md border-[0.5px] border-gray-400">
        <p className="h-full w-1/8 text-xl text-center py-3 font-semibold">Rank</p>
        <div className="flex flex-row justify-between h-full w-full">
          <p className="font-semibold text-xl h-full w-1/5 text-start py-3">Full Name</p>
          <p className="font-semibold text-xl h-full w-1/5 text-center py-3">Worked For</p>
          <p className="font-semibold text-xl h-full w-1/5 text-center py-3">Experience Level</p>
          <p className="font-semibold text-xl h-full w-1/5 text-center py-3">Authenticity</p>
          <p className="font-semibold text-xl h-full w-1/5 text-center py-3">Final Score</p>
        </div>
      </div>
      <div className="flex flex-col gap-2 w-5/6 mx-auto mt-4 pb-40">
        {sharedData?.list.map((candidate, index) => (<>
          <Display evalPath={candidate.path_to_evaluation} key={index} rank={index + 1} name={candidate.full_name} score={candidate.final_score} authenticity={candidate.authenticity} companies={candidate.worked_for} experience={candidate.experience_level} />
          </>
        ))}
      </div>
      <div className="absolute bottom-0 w-full ">
      <img src={Logo1} alt="Logo 1" className="fixed bottom-0 right-0 translate-y-2/5 -z-10" />
      <img src={Logo2} alt="Logo 2" className="fixed top-0 right-0 -translate-x-[200px] -z-10" />
      <img src={Logo3} alt="Logo 3" className="fixed bottom-0 left-0 translate-x-[30px] -z-10" />
      <img src={Logo4} alt="Logo 4" className="fixed top-0 left-0 translate-y-[30px] -z-10" />
        <Footer />
      </div>
    </section>
  );
};

export default DashboardPage;