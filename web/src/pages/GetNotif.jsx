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
    <div>
      <button onClick={fetchMSG}>
        Get Notification
      </button>

      <div className="senders">
        {senders.map((sender) => (
          <button
  key={sender}
  onClick={() => navigate("/msg", { state: { sender, messages: senderMap.get(sender) } })}
>
  {sender}
</button>
        ))}
      </div>
    </div>
  );
}

export { Notif };

