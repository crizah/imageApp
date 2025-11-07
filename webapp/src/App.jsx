import React, { useState, useEffect } from 'react';
import './App.css';
import { Amplify } from 'aws-amplify';
import { withAuthenticator, Button, Heading } from '@aws-amplify/ui-react';
import awsconfig from './aws-exports';
import { BrowserRouter as Router, Routes, Route, useNavigate } from "react-router-dom";
import {v4 as uuidv4} from 'uuid';

import axios from "axios";

Amplify.configure(awsconfig);








// function Messages() {
//   const [messages, setMessages] = useState([]);
//   const [loading, setLoading] = useState(true);
//   const [selectedMessage, setSelectedMessage] = useState(null);
//   const [imageData, setImageData] = useState(null);
//   const [username, setUsername] = useState('');

//   useEffect(() => {
//     fetchMessages();
//   }, []);

//   async function fetchMessages() {
//     try {
//       const user = await Auth.currentAuthenticatedUser();
//       const currentUsername = user.username;
//       setUsername(currentUsername);

//       const response = await axios.get(
//         `http://localhost:8080/messages?username=${currentUsername}`
//       );

//       setMessages(response.data.messages || []);
//       setLoading(false);
//     } catch (err) {
//       console.error('Error fetching messages:', err);
//       setLoading(false);
//     }
//   }

//   async function viewMessage(messageID) {
//     try {
//       setSelectedMessage(messageID);
//       setImageData(null);

//       const response = await axios.get(
//         `http://localhost:8080/message?messageID=${messageID}&username=${username}`
//       );

//       setImageData(response.data);
//     } catch (err) {
//       console.error('Error viewing message:', err);
//       alert('Failed to load message: ' + err.message);
//     }
//   }

//   function closeMessage() {
//     setSelectedMessage(null);
//     setImageData(null);
//     // Refresh messages to update read status
//     fetchMessages();
//   }

//   if (loading) {
//     return <div>Loading messages...</div>;
//   }

//   return (
//     <div style={{ padding: '20px' }}>
//       <h1>Your Messages</h1>

//       {messages.length === 0 ? (
//         <p>No messages yet.</p>
//       ) : (
//         <div>
//           {messages.map((msg) => (
//             <div
//               key={msg.messageID}
//               onClick={() => viewMessage(msg.messageID)}
//               style={{
//                 border: '1px solid #ccc',
//                 padding: '15px',
//                 margin: '10px 0',
//                 cursor: 'pointer',
//                 backgroundColor: msg.status === 'unread' ? '#e3f2fd' : '#fff',
//                 borderRadius: '5px',
//               }}
//             >
//               <div style={{ display: 'flex', justifyContent: 'space-between' }}>
//                 <div>
//                   <strong>From: {msg.sender}</strong>
//                   <p style={{ margin: '5px 0', color: '#666' }}>
//                     File: {msg.fileName}
//                   </p>
//                   {msg.timestamp && (
//                     <p style={{ margin: '5px 0', fontSize: '12px', color: '#999' }}>
//                       {new Date(msg.timestamp).toLocaleString()}
//                     </p>
//                   )}
//                 </div>
//                 {msg.status === 'unread' && (
//                   <span
//                     style={{
//                       backgroundColor: '#2196F3',
//                       color: 'white',
//                       padding: '5px 10px',
//                       borderRadius: '12px',
//                       fontSize: '12px',
//                       height: 'fit-content',
//                     }}
//                   >
//                     NEW
//                   </span>
//                 )}
//               </div>
//             </div>
//           ))}
//         </div>
//       )}

//       {/* Message Viewer Modal */}
//       {selectedMessage && (
//         <div
//           style={{
//             position: 'fixed',
//             top: 0,
//             left: 0,
//             right: 0,
//             bottom: 0,
//             backgroundColor: 'rgba(0,0,0,0.8)',
//             display: 'flex',
//             alignItems: 'center',
//             justifyContent: 'center',
//             zIndex: 1000,
//           }}
//           onClick={closeMessage}
//         >
//           <div
//             style={{
//               backgroundColor: 'white',
//               padding: '20px',
//               borderRadius: '10px',
//               maxWidth: '90%',
//               maxHeight: '90%',
//               overflow: 'auto',
//             }}
//             onClick={(e) => e.stopPropagation()}
//           >
//             {imageData ? (
//               <div>
//                 <h2>From: {imageData.sender}</h2>
//                 <p>File: {imageData.fileName}</p>
//                 <img
//                   src={`data:image/jpeg;base64,${imageData.imageData}`}
//                   alt="Message"
//                   style={{ maxWidth: '100%', marginTop: '20px' }}
//                 />
//                 <button
//                   onClick={closeMessage}
//                   style={{
//                     marginTop: '20px',
//                     padding: '10px 20px',
//                     cursor: 'pointer',
//                   }}
//                 >
//                   Close
//                 </button>
//               </div>
//             ) : (
//               <p>Loading image...</p>
//             )}
//           </div>
//         </div>
//       )}
//     </div>
//   );
// }



function Messages({ receiver }) {
  const [senders, setSenders] = useState([]);
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

        
     
      </Routes>
    </Router>
  );
}

export default withAuthenticator(App, {
  signUpAttributes: ['email'],
});
