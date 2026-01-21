import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useLocation } from "react-router-dom";
import axios from "axios";


import "../index.css"


function Msg() {
  const location = useLocation();
  const { sender, messages } = location.state || {};

  const [images, setImages] = useState([]);
  const x = window.RUNTIME_CONFIG.BACKEND_URL;
  // const x = process.env.REACT_APP_BACKEND_URL;

  const getMessages = async () => {
    try {
      const res = await axios.post(`${x}/files`, {
        msgs: messages,
      }, {withCredentials: true, headers:{
        "Content-Type": "application/json"
      }});

      //array of base64-encoded image strings
      const fileData = res.data.files || [];

      // Convert base64 strings to object URLs for <img> display
      const urls = fileData.map((b64) => {
        // decode base64 → binary
        const binary = atob(b64);
        const bytes = new Uint8Array(binary.length);
        for (let i = 0; i < binary.length; i++) {
          bytes[i] = binary.charCodeAt(i);
        }

        const blob = new Blob([bytes], { type: "image/jpeg" }); // adjust MIME type if needed
        return URL.createObjectURL(blob);
      });

      setImages(urls);
    } catch (error) {
      alert("Error: " + error);
    }
  };




  return (
  <div className="relative min-h-screen bg-slate-950 overflow-hidden">
    {/* Background */}
    <div className="absolute top-0 left-0 w-96 h-96 bg-cyan-500 blur-3xl opacity-20 animate-pulse"></div>
    <div className="absolute bottom-0 right-0 w-96 h-96 bg-purple-500 blur-3xl opacity-20 animate-pulse"></div>

    {/* Stars */}
    <div className="absolute inset-0 overflow-hidden">
      {[...Array(25)].map((_, i) => (
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

    <div className="relative z-10 max-w-6xl mx-auto px-6 py-10">
      <div className="bg-slate-900/40 backdrop-blur-xl border border-slate-800/50 rounded-2xl shadow-2xl p-8">

        <h2 className="text-3xl font-bold bg-gradient-to-r from-cyan-400 via-purple-400 to-pink-400 bg-clip-text text-transparent mb-6">
          Messages from {sender}
        </h2>

        <button
          onClick={getMessages}
          className="mb-8 px-6 py-3 rounded-lg font-semibold text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:from-purple-600 hover:to-pink-600 transition"
        >
          Decrypt & Load Images
        </button>

        {images.length === 0 && (
          <p className="text-gray-400 text-sm">
            No decrypted images yet
          </p>
        )}

        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6">
          {images.map((src, index) => (
            <div
              key={index}
              className="relative group rounded-xl overflow-hidden border border-slate-700 bg-slate-800/60"
            >
              <img
                src={src}
                alt={`Decrypted ${index}`}
                className="w-full h-48 object-cover transition-transform duration-300 group-hover:scale-105"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent opacity-0 group-hover:opacity-100 transition"></div>
            </div>
          ))}
        </div>
      </div>
    </div>
  </div>
);

}

export { Msg };