import React from 'react';
import {Routes, Route,} from 'react-router-dom';
import './index.css';
import InputPage from './pages/inputPage';
import DashboardPage from './pages/dashboard';
import CandidateDetailPage from './pages/candidateDetailPage';

const App = () => {
  return (
    <Routes>
        <Route path="/" element={<InputPage />} />
        <Route path="/dashboard" element={<DashboardPage />} />
        <Route path="/candidate/:rank" element={<CandidateDetailPage />} /> 
        {/* <Route path="/display"/> */}
    </Routes>
  );
}

export default App;
