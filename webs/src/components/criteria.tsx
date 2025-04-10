import React, { useContext } from "react";
import { DataContext } from "./datacontext";

const Criteria = ({ criteriaData, loading }: { criteriaData: any; loading: boolean }) => {
 // Log the criteria data for debugging
  const { criteriaJson } = useContext(DataContext);
  console.log("Criteria JSON in Criteria component:", criteriaJson); // Log the criteriaJson
   // Access criteriaJson from context
  return (
    <div className="flex justify-center items-center mt-12 w-5/6 h-[600px] mx-auto bg-primary shadow-md text-color-black p-8 border-[0.5px] border-gray-400 rounded-sm">
      <div className="w-full h-full bg-white">
        <h2 className="text-center max-h-1/6 font-bold py-4 text-2xl border-b-2 h-20">
          Suggested Criteria for Assessing a Resume
        </h2>
        <div className="p-4 h-5/6 overflow-y-auto">
          {loading ? (
            // Show pulse effect while loading
            <div className="animate-pulse h-full flex flex-col space-y-4">
              <div className="h-1/4 bg-gray-300 rounded"></div>
              <div className="h-1/4 bg-gray-300 rounded"></div>
              <div className="h-1/4 bg-gray-300 rounded"></div>
              <div className="h-1/4 bg-gray-300 rounded"></div>
            </div>
          ) : (
            <>
              {/* Render Main Categories */}
              <h3 className="font-bold text-xl mb-4">Main Categories</h3>
              {criteriaJson?.MainCategory?.map((item: any, index: number) => (
                <div key={index} className="mb-6">
                  <h4 className="font-semibold text-lg">{item.Description}</h4>
                  <p className="text-sm mb-2">
                    <strong>Evaluation Strategy:</strong> {item.EvaluationStrategy}
                  </p>
                  <p className="text-sm mb-2">
                    <strong>Scoring Scale:</strong> {item.ScoringScale}
                  </p>
                  <div className="text-sm">
                    <strong>Scoring Guided:</strong>
                    <ul className="list-disc pl-5">
                      {item.ScoringGuided.map((guide: any, guideIndex: number) => (
                        <li key={guideIndex}>
                          <strong>Range:</strong> {guide.Range},{" "}
                          <strong>Comment:</strong> {guide.Comment}
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              ))}

              {/* Render Sub Categories */}
              <h3 className="font-bold text-xl mb-4">Sub Categories</h3>
              {criteriaJson?.SubCategory?.map((item: any, index: number) => (
                <div key={index} className="mb-6">
                  <h4 className="font-semibold text-lg">{item.Description}</h4>
                  <p className="text-sm mb-2">
                    <strong>Evaluation Strategy:</strong> {item.EvaluationStrategy}
                  </p>
                  <p className="text-sm mb-2">
                    <strong>Scoring Scale:</strong> {item.ScoringScale}
                  </p>
                  <div className="text-sm">
                    <strong>Scoring Guided:</strong>
                    <ul className="list-disc pl-5">
                      {item.ScoringGuided.map((guide: any, guideIndex: number) => (
                        <li key={guideIndex}>
                          <strong>Range:</strong> {guide.Range},{" "}
                          <strong>Comment:</strong> {guide.Comment}
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              ))}
            </>
          )}
        </div>
      </div>
    </div>
  );
};

export default Criteria;