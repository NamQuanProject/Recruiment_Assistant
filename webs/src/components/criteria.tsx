import  { useContext } from "react";
import { DataContext } from "./datacontext";

const Criteria = ({  loading }: { loading: boolean }) => {
  const { criteriaJson } = useContext(DataContext); // Access criteriaJson from context
  console.log("Criteria JSON in Criteria component:", criteriaJson); // Log the criteriaJson

  return (
    <div className="flex justify-center items-center my-20 w-5/6 h-[700px] mx-auto bg-primary shadow-md text-color-black p-8 border-[0.5px] border-gray-400 rounded-sm">
      <div className="w-full h-full bg-white border-[0.5px] border-gray-400 rounded-sm shadow-md p-4">
        <h2 className="text-center h-1/8 font-bold py-4 text-2xl ">
          Criteria for Assessing a Resume
        </h2>
        <div className="h-7/8 overflow-y-auto border-1 border-gray-300 p-3">
          {loading ? (
            // Show pulse effect while loading
            <div className="animate-pulse h-full flex flex-col space-y-4">
              <div className="h-1/4 bg-gray-400 rounded"></div>
              <div className="h-1/4 bg-gray-400 rounded"></div>
              <div className="h-1/4 bg-gray-400 rounded"></div>
              <div className="h-1/4 bg-gray-400 rounded"></div>
            </div>
          ) : Array.isArray(criteriaJson?.MainCategory) && criteriaJson?.MainCategory.length > 0 ||
          Array.isArray(criteriaJson?.SubCategory) && criteriaJson?.SubCategory.length > 0 ? (
            <>
              {/* Render Main Categories */}
              <h3 className="font-bold text-2xl mb-2">-- Main Categories --</h3>
              {criteriaJson?.MainCategory?.map((item: any, index: number) => (
                <div key={index} className="mb-2">
                  <h4 className="font-semibold text-lg leading-6 my-2">
                    📌 {index + 1}. {item.Description}
                  </h4>
                  <p className="text-base my-2 pl-6">
                    <strong>🔍Evaluation Strategy:</strong>
                    <br />
                    {item.EvaluationStrategy}
                  </p>
                  <p className="text-base my-2 pl-6">
                    <strong>🌟Scoring Scale:</strong> {item.ScoringScale}
                  </p>
                  <div className="text-sm pl-6">
                    <ul className="list-disc pl-6">
                      {item.ScoringGuided.map((guide: any, guideIndex: number) => (
                        <li key={guideIndex} className="text-sm my-[2px]">
                          <strong>Range</strong> {guide.Range}: {guide.Comment}
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              ))}

              {/* Render Sub Categories */}
              <h3 className="font-bold text-2xl mt-10 mb-2">-- Sub Categories --</h3>
              {criteriaJson?.SubCategory?.map((item: any, index: number) => (
                <div key={index} className="mb-2">
                  <h4 className="font-semibold text-lg leading-6 my-2">
                    📌 {index + 1}. {item.Description}
                  </h4>
                  <p className="text-base my-2 pl-6">
                    <strong>🔍Evaluation Strategy:</strong>
                    <br />
                    {item.EvaluationStrategy}
                  </p>
                  <p className="text-base my-2 pl-6">
                    <strong>🌟Scoring Scale:</strong> {item.ScoringScale}
                  </p>
                  <div className="text-sm pl-6">
                    <ul className="list-disc pl-6">
                      {item.ScoringGuided.map((guide: any, guideIndex: number) => (
                        <li key={guideIndex}>
                          <strong>Range</strong> {guide.Range}: {guide.Comment}
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              ))}
            </>
          ) : (
            // Fallback UI when no data is available
            <div className="flex flex-col items-center justify-center h-full">
              <p className="text-lg font-semibold text-gray-600">
                No criteria data available.
              </p>
              <p className="text-sm text-gray-500">
                Please upload a resume or provide input to generate criteria.
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Criteria;
