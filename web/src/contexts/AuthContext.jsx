// contexts/AuthContext.jsx
import { createContext, useContext, useState, useEffect } from "react";
import axios from "axios";

const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  // const x = window.RUNTIME_CONFIG.BACKEND_URL;
  const x = process.env.REACT_APP_BACKEND_URL;

  // Check if user is authenticated on app load
  useEffect(() => {
    checkAuth();
  }, []);



  const checkAuth = async () => {
    try {
      const res = await axios.get(`${x}/auth/check`, {
        withCredentials: true,
      });
      if (res.data.authenticated) {
        setUser({ username: res.data.username });
      }
    } catch (error) {
      setUser(null);
    } finally {
      setLoading(false);
    }
  };
  

  const logout = async () => {
    
    try {
      await axios.post(`${x}/logout`,{}, { withCredentials: true, headers:{
        "Content-Type": "application/json"

      }});
      setUser(null);
    } catch (error) {
      console.error("Logout failed:", error);
    }
  };

  return (
    <AuthContext.Provider value={{ user, loading, logout, checkAuth }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}