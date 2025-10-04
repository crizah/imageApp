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





function SendMsg() {


  // allow to select receoiver from list of usernames

  // place to upload a file(goes to s3)

  const [file, setFile] = useState(null);

    function handleChange(e) {
        console.log(e.target.files);
        setFile(URL.createObjectURL(e.target.files[0]));
    }

 

  const [users, setUsers] = useState([]);
  const [search, setSearch] = useState('');
  const [selectedUser, setSelectedUser] = useState('');
  const [showDropdown, setShowDropdown] = useState(false);

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
            <input type="file" onChange={handleChange} />
            {file && <img src={file} alt="Uploaded preview" />}

      </div>
    </div>
  );
}


async function getUsers(){
  try{
    const users = await axios.get(
      "https://5dx7ydfhxe.execute-api.eu-north-1.amazonaws.com/production/resource"
    );


    return users.data.usernames || [];
    
  }catch (err){
    console.log(err)
  }
}


async function saveUser(username, userID) {
  try {
    const payload = {
      username,
      userID
    };

    const res = await axios.post(
      "https://tusnmpvawj.execute-api.eu-north-1.amazonaws.com/prod/res", 
      payload
    );

    console.log("User saved:", res.data);
  } catch (err) {
    console.error("Error saving user:", err.response ? err.response.data : err.message);
  }
}









function HomePage({ signOut, user }) { 
  const navigate = useNavigate(); 
  const hasSaved = React.useRef(false);

  React.useEffect(() => {
    if (user && !hasSaved.current) { 
      const username = user?.username;
      const userId = user?.userId;
      // const email = user?.attributes?.email || user?.email;
      if (username && userId) {
        saveUser(username, userId);
        hasSaved.current = true; // mark as saved for this session
      }
      
    }
  }, [user]);

  return ( 
    <div> 
      <h1>Signed in as {user?.username}</h1> 
      <button onClick={() => navigate('/sendmsg')}>Send Message</button> 
      <button onClick={signOut}>Sign Out</button> 
    </div> 
  ); 
}


function App({ signOut, user }) {
  // console.log('Full user object:', user);
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage signOut={signOut} user={user} />} />
        <Route path="/sendmsg" element={<SendMsg />} />
      </Routes>
    </Router>
  );
}

export default withAuthenticator(App, {
  signUpAttributes: ['email'],
});
