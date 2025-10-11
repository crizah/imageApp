import React, { useState, useEffect } from 'react';
import './App.css';
import { Amplify } from 'aws-amplify';
import { withAuthenticator, Button, Heading } from '@aws-amplify/ui-react';
import awsconfig from './aws-exports';
import { BrowserRouter as Router, Routes, Route, useNavigate } from "react-router-dom";


import axios from "axios";

Amplify.configure(awsconfig);


function Send(userA, userB, attatchment, bucket_name){

  
  // trigger lamda function 

  // get key from KMS
  
// encrypts attatchment

  // send encrypted to  s3


  // send metadata to dynamoDb  (userA, encrypted, key(get from KMS, allowed access), allowed recipients (add B)) (search with key, if already exists, update allowed recipents, else make the entrtry)


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
    formData.append('file', file);  // Send actual file
    formData.append('recipient', selectedUser);
    formData.append('sender', username);

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
    // Reset form...
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
  // const hasSaved = React.useRef(false);

  // React.useEffect(() => {
  //   if (user && !hasSaved.current) { 
  //     const username = user?.username;
  //     const userId = user?.userId;
  //     // const email = user?.attributes?.email || user?.email;
  //     if (username && userId) {
  //       saveUser(username, userId);
  //       hasSaved.current = true; // mark as saved for this session
  //     }
      
  //   }
  // }, [user]);

  return ( 
    <div> 
      <h1>Signed in as {user?.username}</h1> 
      <button onClick={() => navigate('/sendmsg')}>Send Message</button> 
      <button onClick={signOut}>Sign Out</button> 
    </div> 
  ); 
}


function App({ signOut, user }) {
  console.log('Full user object:', user);
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage signOut={signOut} user={user} />} />
        <Route path="/sendmsg" element={<SendMsg username= {user?.username}/>} />
      </Routes>
    </Router>
  );
}

export default withAuthenticator(App, {
  signUpAttributes: ['email'],
});
