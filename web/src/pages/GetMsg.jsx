import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useLocation } from "react-router-dom";
import axios from "axios";


import "../index.css"


function Msg() {
  const location = useLocation();
  const { sender, messages } = location.state || {};

  const [images, setImages] = useState([]);
  // const x = window.RUNTIME_CONFIG.BACKEND_URL;
  const x = process.env.REACT_APP_BACKEND_URL;

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
    <div>
      <h2>Messages from: {sender}</h2>

      <button onClick={getMessages}>Decrypt & Load Images</button>

      <div style={{ display: "flex", flexWrap: "wrap", gap: "10px", marginTop: "20px" }}>
        {images.map((src, index) => (
          <img
            key={index}
            src={src}
            alt={`Decrypted ${index}`}
            style={{ width: "200px", borderRadius: "8px", boxShadow: "0 0 8px #aaa" }}
          />
        ))}
      </div>
    </div>
  );
}

export { Msg };