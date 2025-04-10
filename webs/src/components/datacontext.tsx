import React, { createContext, useState, ReactNode } from "react";

export interface CandidateData {
  name: string;
  rank: number;
  [key: string]: any; // Allow additional properties
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

export interface CriteriaJson {
  MainCategory: Category[];
  SubCategory: Category[];
}

interface DataContextType {
  sharedData: CandidateData | null;
  setSharedData: (data: CandidateData) => void;
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
  const [sharedData, setSharedData] = useState<CandidateData | null>(null);
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