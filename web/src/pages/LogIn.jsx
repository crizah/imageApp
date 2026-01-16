import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

import "../index.css";

export function Login() {
  const navigate = useNavigate();
  
  const x = window.RUNTIME_CONFIG.BACKEND_URL;

  const [formData, setFormData] = useState({
    identifier: "",
    password: ""
  });
  const [message, setMessage] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setMessage("");

    // const isEmail = formData.identifier.includes("@");

    try {
      await axios.post(
        `${x}/login`, formData,
        {
          withCredentials: true,
          headers: {
            "Content-Type": "application/json"
          }
        }
      );

      navigate("/home");
    } catch (error) {
      if (error.response) {
        setMessage(error.response.data || "Login failed");
      } else {
        setMessage("Cannot connect to server");
      }
    } finally {
      setIsLoading(false);
    }
  };

return (
  <div className="relative flex items-center justify-center min-h-screen overflow-hidden bg-slate-950">
    <div className="absolute top-0 left-0 w-96 h-96 bg-blue-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-pulse"></div>
    <div className="absolute top-0 right-0 w-96 h-96 bg-indigo-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-pulse" style={{ animationDelay: '2s' }}></div>
    <div className="absolute bottom-0 left-1/2 w-96 h-96 bg-violet-500 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-pulse" style={{ animationDelay: '4s' }}></div>
    
    <div className="absolute inset-0 overflow-hidden">
      {[...Array(20)].map((_, i) => (
        <div
          key={i}
          className="absolute w-1 h-1 bg-white rounded-full opacity-30 animate-pulse"
          style={{
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 5}s`
          }}
        />
      ))}
    </div>

    <div className="relative z-10 w-full max-w-md px-6">
      <div className="bg-slate-900/40 backdrop-blur-xl rounded-2xl shadow-2xl p-8 border border-slate-800/50 transform transition-all duration-300 hover:scale-[1.02]">
        <div className="flex justify-center mb-8">
          <div className="relative">
            <div className="absolute inset-0 bg-gradient-to-r from-blue-400 to-indigo-500 rounded-full blur-md opacity-75"></div>
            <div className="relative w-16 h-16 bg-gradient-to-br from-blue-400 via-indigo-500 to-violet-500 rounded-full flex items-center justify-center">
              <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
              </svg>
            </div>
          </div>
        </div>

        <h2 className="text-3xl font-bold text-center bg-gradient-to-r from-blue-400 via-indigo-400 to-violet-400 bg-clip-text text-transparent mb-2">
          Welcome Back
        </h2>
        <p className="text-center text-gray-400 mb-8 text-sm">Login to continue your journey</p>

        <div className="space-y-5">
          <div className="relative group">
            <div className={`absolute inset-0 bg-gradient-to-r from-blue-500 to-indigo-500 rounded-lg blur opacity-0 group-hover:opacity-25 transition-opacity`}></div>
            <div className="relative">
              <input
                type="text"
                placeholder="Username or Email"
                value={formData.identifier}
                onChange={(e) => setFormData({ ...formData, identifier: e.target.value })}
                disabled={isLoading}
                className="w-full px-4 py-3 bg-slate-900/50 border border-slate-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-blue-500 transition-all duration-300 backdrop-blur-sm"
              />
              <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                <svg className="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
            </div>
          </div>

          <div className="relative group">
            <div className={`absolute inset-0 bg-gradient-to-r from-indigo-500 to-violet-500 rounded-lg blur opacity-0 group-hover:opacity-25 transition-opacity`}></div>
            <div className="relative">
              <input
                type="password"
                placeholder="Password"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                disabled={isLoading}
                className="w-full px-4 py-3 bg-slate-900/50 border border-slate-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500 transition-all duration-300 backdrop-blur-sm"
              />
              <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                <svg className="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
            </div>
          </div>

          <button
            onClick={handleSubmit}
            disabled={isLoading}
            className="relative w-full py-3 mt-6 overflow-hidden rounded-lg font-semibold text-white transition-all duration-300 group disabled:opacity-50"
          >
            <div className="absolute inset-0 bg-gradient-to-r from-blue-500 via-indigo-500 to-violet-500"></div>
            <div className="absolute inset-0 bg-gradient-to-r from-blue-600 via-indigo-600 to-violet-600 opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <span className="relative flex items-center justify-center gap-2">
              {isLoading ? (
                <>
                  <svg className="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Logging in...
                </>
              ) : (
                'Login'
              )}
            </span>
          </button>

          {message && (
            <div className="p-3 rounded-lg text-center text-sm font-medium bg-red-500/20 text-red-400 border border-red-500/50">
              {message}
            </div>
          )}
        </div>

        <div className="mt-6 text-center">
          <p className="text-gray-400 text-sm">
            New user?{' '}
            <span
              onClick={() => navigate('/signup')}
              className="text-transparent bg-gradient-to-r from-blue-400 to-indigo-400 bg-clip-text font-semibold cursor-pointer hover:from-blue-300 hover:to-indigo-300 transition-all"
            >
              Sign up
            </span>
          </p>
        </div>
      </div>

      <div className="absolute -top-4 -right-4 w-24 h-24 bg-indigo-500/20 rounded-full blur-2xl"></div>
      <div className="absolute -bottom-4 -left-4 w-32 h-32 bg-blue-500/20 rounded-full blur-2xl"></div>
    </div>
  </div>
);
}