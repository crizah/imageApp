import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";
import axios from "axios";


import "../index.css"


function Notif() {
  const [senders, setSenders] = useState([]);
  const [senderMap, setSenderMap] = useState(new Map());
  const {user} = useAuth();
  const receiver = user?.username;
  const navigate = useNavigate();
  // const x = window.RUNTIME_CONFIG.BACKEND_URL;
  const x = process.env.REACT_APP_BACKEND_URL;
  const fetchMSG = async () => {
    try {
      const res = await axios.post(`${x}/notifs`, {
        username: receiver,  
      }, {withCredentials: true, headers: {
                        "Content-Type": "application/json"}
                    });

      const count = res.data.count;
      const msgs = res.data.msgs; // array of objects
      // iterate msgs and group them based on sender
      // map[string][object]
      const m = new Map();
      for (const msg of msgs) {
        if (!m.has(msg.Sender)) {
          m.set(msg.Sender, [msg]); 
        } else {
           m.get(msg.Sender).push(msg); 
        }
      }

      setSenders([...m.keys()]);
      setSenderMap(m);


      console.log(m.size);
      // based on m.size, thats the number of box elements with the text as the sender

      alert(`got msgs ${count}`);
    } catch (error) {
      alert("Error while getting messages: " + error.message);
    }
  };

  

return (
  <div className="relative min-h-screen bg-slate-950 overflow-hidden flex items-center justify-center">
    {/* Cosmic background */}
    <div className="absolute top-0 left-0 w-96 h-96 bg-purple-500 blur-3xl opacity-20 animate-pulse"></div>
    <div className="absolute bottom-0 right-0 w-96 h-96 bg-cyan-500 blur-3xl opacity-20 animate-pulse"></div>

    {/* Stars */}
    <div className="absolute inset-0 overflow-hidden">
      {[...Array(20)].map((_, i) => (
        <div
          key={i}
          className="absolute w-1 h-1 bg-white rounded-full opacity-40 animate-pulse"
          style={{
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 5}s`,
          }}
        />
      ))}
    </div>

    <div className="relative z-10 w-full max-w-md px-6">
      <div className="bg-slate-900/40 backdrop-blur-xl border border-slate-800/50 rounded-2xl shadow-2xl p-8">

        <h1 className="text-3xl font-bold text-center bg-gradient-to-r from-cyan-400 via-purple-400 to-pink-400 bg-clip-text text-transparent mb-6">
          Incoming Signals
        </h1>

        <button
          onClick={fetchMSG}
          className="w-full mb-6 py-3 rounded-lg font-semibold text-white bg-gradient-to-r from-cyan-500 to-purple-500 hover:from-cyan-600 hover:to-purple-600 transition"
        >
          Scan for Messages
        </button>

        <div className="space-y-3">
          {senders.map((sender) => (
            <button
              key={sender}
              onClick={() =>
                navigate("/msg", {
                  state: { sender, messages: senderMap.get(sender) },
                })
              }
              className="w-full text-left px-4 py-3 rounded-lg bg-slate-800/60 border border-slate-700 text-gray-300 hover:bg-slate-700/60 hover:text-white transition"
            >
              <span className="font-medium">{sender}</span>
            </button>
          ))}

          {senders.length === 0 && (
            <p className="text-center text-gray-500 text-sm">
              No signals detected yet
            </p>
          )}
        </div>
      </div>
    </div>
  </div>
);
}

export { Notif };

