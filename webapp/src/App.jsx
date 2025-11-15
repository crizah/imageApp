import React, { useState, useEffect } from 'react';
import './App.css';
import { Amplify } from 'aws-amplify';
import { withAuthenticator, Button, Heading } from '@aws-amplify/ui-react';
import awsconfig from './aws-exports';
import { BrowserRouter as Router, Routes, Route, useNavigate } from "react-router-dom";
import {v4 as uuidv4} from 'uuid';
import { useLocation } from "react-router-dom";

import axios from "axios";

Amplify.configure(awsconfig);










function Msg() {
  const location = useLocation();
  const { sender, messages } = location.state || {};

  const [images, setImages] = useState([]);

  const getMessages = async () => {
    try {
      const res = await axios.post("http://localhost:8080/files", {
        msgs: messages,
      });

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


function Messages({ receiver }) {
  const [senders, setSenders] = useState([]);
  const [senderMap, setSenderMap] = useState(new Map());
  const navigate = useNavigate();
  const fetchMSG = async () => {
    try {
      const res = await axios.post("http://localhost:8080/messages", {
        username: receiver,  
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











function SendMsg({username}) {
  const [file, setFile] = useState(null);
  const [uploading, setUploading] = useState(false)



  const [filePreview, setFilePreview] = useState(null);
  const [users, setUsers] = useState([]);
  const [search, setSearch] = useState('');
  const [selectedUser, setSelectedUser] = useState('');
  const [showDropdown, setShowDropdown] = useState(false);



    const allowedTypes = [
    'image/jpeg',
    'image/png',
  ];

  useEffect(() => {
    async function fetchUsers() {
      try {
        const res = await axios.get(
          "https://5dx7ydfhxe.execute-api.eu-north-1.amazonaws.com/production/resource"
        );
        setUsers(res.data.usernames || []);
        console.log("got users", res.data.usernames);
      } catch (err) {
        console.error("Error fetching users:", err);
      }
    }
    fetchUsers();
  }, []);



  function handleFileChange(e) {
    const selectedFile = e.target.files[0];
    if (allowedTypes.includes(selectedFile.type)) {
      setFile(selectedFile);
      setFilePreview(URL.createObjectURL(selectedFile));
    } else {
      alert('Invalid file type. Only images and PDFs are allowed.');
    }
  }



  async function handleSendMessage() {
  if (!file) {
    alert("Please select a file");
    return;
  }
  if (!selectedUser) {
    alert("Please select a recipient");
    return;
  }

  setUploading(true);





  

  try {
    const formData = new FormData();
    const msgID = uuidv4();
    formData.append('file', file);  // Send actual file
    formData.append('recipient', selectedUser);
    formData.append('sender', username);
    formData.append('msgID', msgID)
    const res = await axios.post(
      'http://localhost:8080/upload',
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      }
    );




    alert("Message sent successfully!");




    console.log("Notification message sent");
   
  } catch (error) {
    alert("Error: " + error.message);
  } finally {
    setUploading(false);
  }

  
}



  const filteredUsers = users.filter(u =>
    u.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div>
      <h1>Send Message</h1>
      <div>
        <input
          type="text"
          placeholder="Type username"
          value={search}
          onChange={e => {
            setSearch(e.target.value);
            setShowDropdown(true);
          }}
          onFocus={() => setShowDropdown(true)}
        />
        {showDropdown && filteredUsers.length > 0 && (
          <ul style={{ border: '1px solid #ccc', maxHeight: '150px', overflowY: 'auto' }}>
            {filteredUsers.map(user => (
              <li
                key={user}
                onClick={() => {
                  setSelectedUser(user);
                  setSearch(user);
                  setShowDropdown(false);
                }}
                style={{ cursor: 'pointer', padding: '5px' }}
              >
                {user}
              </li>
            ))}
          </ul>
        )}
      </div>
      <p>Selected user: {selectedUser}</p>
      
      <div>
        <h2>Add Image:</h2>
        <input type="file" onChange={handleFileChange} />
        {filePreview && <img src={filePreview} alt="Uploaded preview" style={{maxWidth: '300px'}} />}
        <button onClick={handleSendMessage} disabled={uploading}>
    {uploading ? 'Sending...' : 'Send Message'}
  </button>
      </div>
    </div>
  );
}





function HomePage({ signOut, user }) {
  const navigate = useNavigate();

  if (!user) return <div>Loading...</div>;

  const username = user.username;

  return (
    <div>
      <h1>Signed in as {username}</h1>
      <button onClick={() => navigate(`/sendmsg`)}>Send Message</button>
      <button onClick={() => navigate(`/messages`)}>Messages</button>
      <button onClick={signOut}>Sign Out</button>
    </div>
  );
}


function App({ signOut, user }) {
  if (!user) return <div>Loading...</div>;

  const username = user.username;

  return (
    <Router>
      <Routes>
       
       

    
        <Route
          path={`/`}
          element={<HomePage signOut={signOut} user={user} />}
        />
        <Route
          path={`/sendmsg`}
          element={<SendMsg username={username} />}
        />
        <Route
          path={`/messages`}
          element={<Messages receiver={username}/>}
        />
        <Route
          path = {`/msg`}
          element = {<Msg/>}
        />


        
     
      </Routes>
    </Router>
  );
}

export default withAuthenticator(App, {
  signUpAttributes: ['email'],
});
