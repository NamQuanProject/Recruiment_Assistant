import NavBar from "../components/navbar";
import Display from "../components/personDisplay";

const DashboardPage = () => {
  const Candidates = [
    { name: "John Doe", score: 85, companies: "Meta, Amazon", experience: "4 years", rank:  1 },
    { name: "Jane Smith", score: 90, companies: "Meta, Amazon", experience: "4 years", rank: 2  },
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
  return (
    <section id="leaderboard">
      <NavBar/>
      <h1 className="text-center w-full text-3xl font-bold text-red-400 mt-32 ">Candidates Leaderboard</h1>
      <div className="w-5/6 flex flex-row justify-between mx-auto mt-4 bg-secondary rounded-md shadow-md border-[0.5px] border-gray-400 ">
        <p className="h-full w-1/8 text-xl text-center py-3 font-semibold">Rank</p>
        <div className="flex flex-row justify between h-full w-full">
          <p className="font-semibold text-xl h-full w-1/5 text-start py-3 ">Full Name</p>
          <p className=" font-semibold text-xl h-full w-1/5 text-center py-3">Worked For</p>
          <p className=" font-semibold text-xl h-full w-1/5 text-center py-3">Experience Level</p>
          <p className=" font-semibold text-xl h-full w-1/5 text-center py-3">Authenticity</p>
          <p className=" font-semibold text-xl h-full w-1/5 text-center py-3">Final Score</p>
        </div>
      </div>
      <div className="flex flex-col gap-2 w-5/6 max-h-screen mx-auto mt-4">
        {Candidates.map((candidate, index) => (
          <Display companies={candidate.companies} rank={candidate.rank} name={candidate.name} score={candidate.score} experience={candidate.experience} key={index} />
        ))}
      </div>
    </section>
  );
}

export default DashboardPage;