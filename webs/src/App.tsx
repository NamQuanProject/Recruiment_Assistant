// import React from 'react';
import {Routes, Route,} from 'react-router-dom';
import './index.css';
import InputPage from './pages/inputPage';
import DashboardPage from './pages/dashboard';
import CandidateDetailPage from './pages/candidateDetailPage';
import { DataProvider } from './components/datacontext';
// import NavBar from './components/navbar';
import HomePage from './pages/homePage';
const App = () => {
  return (
    <DataProvider>
    <Routes>

        <Route path="/" element={<HomePage />} />
        <Route path="/input" element={<InputPage />} />
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/candidate/:rank" element={<CandidateDetailPage />} /> 
        {/* <Route path="/display"/> */}
    </Routes>
    </DataProvider>
  );
}

export default App;
