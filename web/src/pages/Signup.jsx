import { useState } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";


import "../index.css"

export function SignUp() {
  const [formData, setFormData] = useState({
    username: "",
    email: "",
    password: ""
  });
  const [message, setMessage] = useState("");
  const [isSuccess, setIsSuccess] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
    const [focusedField, setFocusedField] = useState(null);
  const navigate = useNavigate();
  const x = window.RUNTIME_CONFIG.BACKEND_URL;
  // const x = process.env.REACT_APP_BACKEND_URL;

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    const { username, email, password } = formData;

    if (!username || !email || !password) {
      setIsSuccess(false);
      setMessage("All fields are required.");
      return;
    }

    setIsLoading(true);
    setMessage("");

    let config = {
      headers: {
        "Content-Type": "application/json"
        
      }
    }

    try {
      
      
      const res = await axios.post(`${x}/signup`, {
        "username": formData.username,
        "email": formData.email,
        "password": formData.password
      }, {
                        withCredentials: true,
                        headers: {
                        "Content-Type": "application/json"}
                    });

      setIsSuccess(true);
      setMessage("Account created successfully!");
      localStorage.setItem("pendingVerification", formData.username);
      setFormData({ username: "", email: "", password: "" }); // clear
       
      
      setTimeout(() => navigate("/verify"), 1000);

    } catch (error) {
      setIsSuccess(false);
      if (error.response) setMessage(error.response.data || "Signup failed");
      else setMessage("Cannot connect to server");
    } finally {
      setIsLoading(false);
    }
  };

  return (
  <div className="background">
    <div className="tab">
      <text>hello</text>
    </div>
    <div className="test1"></div>
    <div className="test2"></div>
    
    <div className="blackbox">
      <div className="signUp-box">
        
       
          <div className="neon-border"></div>
           <div className="bubble">
            <p className="bubbleText">Sign up!!</p>
          </div>
        <div className="kitty"></div>
        
        
        <div className="container-signup">
          <div className="username">
            <div className="star"></div>
            <div className="star1"></div>
            <input
              type="text"
              name="username"
              placeholder="Username"
              value={formData.username}
              onChange={handleChange}
              onFocus={() => setFocusedField('username')}
              onBlur={() => setFocusedField(null)}
              disabled={isLoading}
              className="field"  
            />
          </div>

          <div className="email">
            <div className="star"></div>
            <div className="star1"></div>
            <input
              type="email"
              name="email"
              placeholder="Email"
              value={formData.email}
              onChange={handleChange}
              onFocus={() => setFocusedField('email')}
              onBlur={() => setFocusedField(null)}
              disabled={isLoading}
              className="field"  
            />
          </div>

          <div className="password">
            <div className="star1"></div>
            <div className="star"></div>
            <input 
              type="password"
              name="password"
              placeholder="Password"
              value={formData.password}
              onChange={handleChange}
              onFocus={() => setFocusedField('password')}
              onBlur={() => setFocusedField(null)}
              disabled={isLoading}
              className="field"
            />
          </div>

          <div className="signUp">
            <div className="star"></div>
            <div className="star1"></div>
            <button
              onClick={handleSubmit}
              disabled={isLoading}
              className="field"
            >
              {isLoading ? 'Creating Account' : 'Create Account'}
            </button>
          </div>

          {message && (
            <div className={`message ${isSuccess ? 'success' : 'error'}`}>
              {message}
            </div>
          )}

          <div className="login">
            Already have an account?{' '}
            <span onClick={() => navigate('/')}>
              Log in
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
);
}








