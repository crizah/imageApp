import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";


import "../index.css"

export function Verification() {
    
    const [code, setCode] = useState('');
    const [message, setMessage] = useState('');
    const [username, setUsername] = useState('');
    const navigate = useNavigate();
    // const x = window.RUNTIME_CONFIG.BACKEND_URL;
    const x = process.env.REACT_APP_BACKEND_URL;
    useEffect(() => {

    

    const pendingUser = localStorage.getItem("pendingVerification");
    if (pendingUser) {
      setUsername(pendingUser);
    } else {
      // No pending verification, redirect to signup
      navigate("/signup");
     }
    }, [navigate]);
  

    const handleVerify = async (e) => {
        e.preventDefault();
           
    
        try {
            const res = await axios.post(`${x}/verify`, {
                username : username,
                verificationCode: code
                }, {
                        withCredentials: true,
                        headers: {
                        "Content-Type": "application/json"}
                    });

            setMessage(res.data.message || 'Verification successful!');
            setTimeout(() => navigate("/"), 1000);
        } catch (err) {
            setMessage(err.response?.data?.message || 'Verification failed. Please try again.');
        }
    };

    return (
  <div className="relative flex items-center justify-center min-h-screen overflow-hidden bg-slate-950">
    {/* Glow blobs */}
    <div className="absolute top-0 left-0 w-96 h-96 bg-cyan-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-pulse"></div>
    <div
      className="absolute bottom-0 right-0 w-96 h-96 bg-purple-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-pulse"
      style={{ animationDelay: "3s" }}
    ></div>

    {/* Floating particles */}
    <div className="absolute inset-0 overflow-hidden">
      {[...Array(15)].map((_, i) => (
        <div
          key={i}
          className="absolute w-1 h-1 bg-white rounded-full opacity-30 animate-pulse"
          style={{
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 5}s`,
          }}
        />
      ))}
    </div>

    {/* Card */}
    <div className="relative z-10 w-full max-w-md px-6">
      <div className="bg-slate-900/40 backdrop-blur-xl rounded-2xl shadow-2xl p-8 border border-slate-800/50 transition-all duration-300 hover:scale-[1.02]">

        {/* Icon */}
        <div className="flex justify-center mb-8">
          <div className="relative">
            <div className="absolute inset-0 bg-gradient-to-r from-cyan-400 to-purple-500 rounded-full blur-md opacity-75"></div>
            <div className="relative w-16 h-16 bg-gradient-to-br from-cyan-400 via-purple-500 to-pink-500 rounded-full flex items-center justify-center">
              <svg
                className="w-8 h-8 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
                />
              </svg>
            </div>
          </div>
        </div>

        {/* Title */}
        <h2 className="text-3xl font-bold text-center bg-gradient-to-r from-cyan-400 via-purple-400 to-pink-400 bg-clip-text text-transparent mb-2">
          Verify Your Email
        </h2>
        <p className="text-center text-gray-400 mb-8 text-sm">
          Enter the verification code sent to your email
        </p>

        {/* Form */}
        <form onSubmit={handleVerify} className="space-y-6">
          <div className="relative group">
            <div className="absolute inset-0 bg-gradient-to-r from-cyan-500 to-purple-500 rounded-lg blur opacity-0 group-hover:opacity-25 transition-opacity"></div>
            <div className="relative">
              <input
                type="text"
                placeholder="Verification Code"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                required
                className="w-full px-4 py-3 bg-slate-900/50 border border-slate-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-cyan-500 transition-all duration-300 backdrop-blur-sm tracking-widest text-center"
              />
            </div>
          </div>

          <button
            type="submit"
            className="relative w-full py-3 overflow-hidden rounded-lg font-semibold text-white transition-all duration-300 group"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-cyan-500 via-purple-500 to-pink-500"></div>
            <div className="absolute inset-0 bg-gradient-to-r from-cyan-600 via-purple-600 to-pink-600 opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span className="relative flex items-center justify-center gap-2">
              Verify Email
            </span>
          </button>
        </form>

        {/* Message */}
        {message && (
          <div className="mt-6 p-3 rounded-lg text-center text-sm font-medium bg-slate-800/60 border border-slate-700 text-gray-300">
            {message}
          </div>
        )}
      </div>

      {/* Corner glows */}
      <div className="absolute -top-4 -right-4 w-24 h-24 bg-purple-500/20 rounded-full blur-2xl"></div>
      <div className="absolute -bottom-4 -left-4 w-32 h-32 bg-cyan-500/20 rounded-full blur-2xl"></div>
    </div>
  </div>
);


}