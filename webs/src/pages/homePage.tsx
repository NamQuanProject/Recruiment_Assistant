import React from "react";
import HomepageImg from "../assets/Homepage.png";
import { motion } from "framer-motion";
import { CheckCircle } from "lucide-react";
import {Link } from "react-router-dom";
// import Footer from "../components/footer";
import Logo from "../assets/logo.png";

const features = [
  "Instant CV Screening with AI",
  "Smart Candidate Matching",
  "Save Hours of Manual Review",
  "Boost Hiring Accuracy",
];

const Homepage: React.FC = () => {
  return (
    <main className="min-h-screen max-h-screen bg-gradient-to-br from-[#1c3156] to-[#0e1a2f] text-white p-6 overflow-y-auto">
      <section className="max-w-6xl mx-auto grid grid-cols-1 md:grid-cols-2 gap-10 items-center">
    
        <div className="mt-6 space-y-6">
          <div className="flex mb-6">
            <img src={Logo} alt="Logo" className="w-32 h-16" />
          </div>
          <motion.h1
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8 }}
            className="text-3xl md:text-6xl font-bold"
          >
            AI-Empowered Recruiment System
          </motion.h1>

          <motion.p
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.2 }}
            className="text-lg md:text-xl text-slate-300"
          >
            Revolutionize your recruitment process with smart, instant CV analysis.
            Reduce bias, speed up hiring, and make data-driven decisions.
          </motion.p>

          <ul className="space-y-2">
            {features.map((feature, idx) => (
              <motion.li
                key={idx}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ duration: 0.4, delay: 0.3 + idx * 0.1 }}
                className="flex items-center gap-2 text-slate-200"
              >
                <CheckCircle className="text-green-400" size={20} /> {feature}
              </motion.li>
            ))}
          </ul>
            <Link to="/input">
          <button className="mt-6 px-6 py-3 text-lg rounded-2xl shadow-lg cursor-pointer bg-green-500 hover:bg-green-600 transition">
            Get Started Free
          </button> 
            </Link>
        </div>

        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
          transition={{ duration: 0.8, delay: 0.3 }}
          className="rounded-xl overflow-hidden shadow-2xl"
        >
          <img
            src={HomepageImg}
            alt="AI CV Screening"
            className="w-full h-auto"
          />
        </motion.div>
      </section>
    </main>
  );
};

export default Homepage;
