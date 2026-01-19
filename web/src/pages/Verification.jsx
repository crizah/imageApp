import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";


import "../index.css"

export function Verification() {
    
    const [code, setCode] = useState('');
    const [message, setMessage] = useState('');
    const [username, setUsername] = useState('');
    const navigate = useNavigate();
    // const x = window.RUNTIME_CONFIG.BACKEND_URL;
    const x = process.env.REACT_APP_BACKEND_URL;
    useEffect(() => {

    

    const pendingUser = localStorage.getItem("pendingVerification");
    if (pendingUser) {
      setUsername(pendingUser);
    } else {
      // No pending verification, redirect to signup
      navigate("/signup");
     }
    }, [navigate]);
  

    const handleVerify = async (e) => {
        e.preventDefault();
           
    
        try {
            const res = await axios.post(`${x}/verify`, {
                username : username,
                verificationCode: code
                }, {
                        withCredentials: true,
                        headers: {
                        "Content-Type": "application/json"}
                    });

            setMessage(res.data.message || 'Verification successful!');
            setTimeout(() => navigate("/"), 1000);
        } catch (err) {
            setMessage(err.response?.data?.message || 'Verification failed. Please try again.');
        }
    };

    return (
        <div className="verification-container">
            <h2>Email Verification</h2>
            <form onSubmit={handleVerify} className="verification-form">
                <input
                    type="text"
                    placeholder="Enter verification code"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    required
                />
                <button type="submit">Verify</button>
            </form>
            {message && <p className="verification-message">{message}</p>}
        </div>
    );

}