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
  <div className="relative min-h-screen bg-slate-950 overflow-hidden flex items-center justify-center">
    {/* Background */}
    <div className="absolute top-0 left-0 w-96 h-96 bg-cyan-500 blur-3xl opacity-20 animate-pulse"></div>
    <div className="absolute bottom-0 right-0 w-96 h-96 bg-pink-500 blur-3xl opacity-20 animate-pulse"></div>

    <div className="relative z-10 w-full max-w-lg px-6">
      <div className="bg-slate-900/40 backdrop-blur-xl border border-slate-800/50 rounded-2xl shadow-2xl p-8">

        <h1 className="text-3xl font-bold text-center bg-gradient-to-r from-cyan-400 via-purple-400 to-pink-400 bg-clip-text text-transparent mb-6">
          Send Message
        </h1>

        {/* Username search */}
        <div className="relative mb-6">
          <input
            type="text"
            placeholder="Search username"
            value={search}
            onChange={e => {
              setSearch(e.target.value);
              setShowDropdown(true);
            }}
            onFocus={() => setShowDropdown(true)}
            className="w-full px-4 py-3 bg-slate-900/60 border border-slate-700 rounded-lg text-white placeholder-gray-500 focus:outline-none focus:border-cyan-500 transition"
          />

          {showDropdown && filteredUsers.length > 0 && (
            <ul className="absolute z-20 mt-2 w-full bg-slate-900 border border-slate-700 rounded-lg max-h-40 overflow-y-auto">
              {filteredUsers.map(user => (
                <li
                  key={user}
                  onClick={() => {
                    setSelectedUser(user);
                    setSearch(user);
                    setShowDropdown(false);
                  }}
                  className="px-4 py-2 cursor-pointer text-gray-300 hover:bg-slate-800 hover:text-white transition"
                >
                  {user}
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Selected user */}
        <p className="text-sm text-gray-400 mb-6">
          Selected user:{" "}
          <span className="text-white font-medium">
            {selectedUser || "None"}
          </span>
        </p>

        {/* Image upload */}
        <div className="mb-6">
          <label className="block mb-2 text-gray-300 font-medium">
            Add Image
          </label>
          <input
            type="file"
            onChange={handleFileChange}
            className="block w-full text-sm text-gray-400
              file:mr-4 file:py-2 file:px-4
              file:rounded-lg file:border-0
              file:bg-slate-800 file:text-white
              hover:file:bg-slate-700 transition"
          />

          {filePreview && (
            <img
              src={filePreview}
              alt="Preview"
              className="mt-4 rounded-lg max-h-64 object-contain border border-slate-700"
            />
          )}
        </div>

        {/* Send button */}
        <button
          onClick={handleSendMessage}
          disabled={uploading}
          className="relative w-full py-3 rounded-lg font-semibold text-white transition-all disabled:opacity-50 overflow-hidden"
        >
          <div className="absolute inset-0 bg-gradient-to-r from-cyan-500 via-purple-500 to-pink-500"></div>
          <div className="absolute inset-0 bg-gradient-to-r from-cyan-600 via-purple-600 to-pink-600 opacity-0 hover:opacity-100 transition"></div>
          <span className="relative">
            {uploading ? "Sending..." : "Send Message"}
          </span>
        </button>
      </div>
    </div>
  </div>
);
}
export { SendMsg };
