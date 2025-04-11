import React, { createContext, useState, ReactNode } from "react";

export interface CandidateData {
  full_name: string;
  worked_for: string;
  experience_level: string;
  authenticity: number;
  final_score: number;
  path_to_cv: string;
  path_to_evaluation: string;
}

export interface ScoringGuide {
  Range: string;
  Comment: string;
}

export interface Category {
  Description: string;
  EvaluationStrategy: string;
  ScoringScale: number;
  ScoringGuided: ScoringGuide[];
}
export interface SharedData {
  list: CandidateData[];
}

export interface CriteriaJson {
  MainCategory: Category[];
  SubCategory: Category[];
}

interface DataContextType {
  sharedData: SharedData | null;
  setSharedData: React.Dispatch<React.SetStateAction<SharedData | null>>;
  criteriaJson: CriteriaJson | null; // Add criteriaJson to the context
  setCriteriaJson: (data: CriteriaJson) => void; // Function to update criteriaJson
}

export const DataContext = createContext<DataContextType>({
  sharedData: null,
  setSharedData: () => {},
  criteriaJson: null,
  setCriteriaJson: () => {},
});

export const DataProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [sharedData, setSharedData] = useState<SharedData | null>(null);
  const [criteriaJson, setCriteriaJson] = useState<CriteriaJson | null>(null); // State for criteriaJson
  const updateCriteriaJson = (data: CriteriaJson) => {
    console.log("Updating Criteria JSON in context:", data); // Log the updated data
    setCriteriaJson(data);
  };
  return (
    <DataContext.Provider value={{ sharedData, setSharedData, criteriaJson, setCriteriaJson: updateCriteriaJson }}>
      {children}
    </DataContext.Provider>
  );
};