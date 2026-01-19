import { useState,  useEffect } from "react";
import { useAuth } from "../contexts/AuthContext";
// import { useNavigate, } from "react-router-dom";
import axios from "axios";
import {v4 as uuidv4} from 'uuid';


import "../index.css"


function SendMsg() {

  const [file, setFile] = useState(null);
  const [uploading, setUploading] = useState(false)



  const [filePreview, setFilePreview] = useState(null);
  const [users, setUsers] = useState([]);
  const [search, setSearch] = useState('');
  const [selectedUser, setSelectedUser] = useState('');
  const [showDropdown, setShowDropdown] = useState(false);
  // const x = window.RUNTIME_CONFIG.BACKEND_URL;
  const x = process.env.REACT_APP_BACKEND_URL;
  const { user, logout } = useAuth();
  const username = user?.username;



    const allowedTypes = [
    'image/jpeg',
    'image/png',
  ];

  useEffect(() => {
    async function fetchUsers() {
    
      try {
        const res = await axios.get(
          `${x}/usernames` , {
          withCredentials: true});
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
    formData.append('sender', user?.username);
    formData.append('msgID', msgID)
    const res = await axios.post(
      `${x}/upload`,
      formData,
      {
          withCredentials: true,
          
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

export { SendMsg };
