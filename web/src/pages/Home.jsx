import { useNavigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";



import "../index.css"

function HomePage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  
  const handleLogout = async () => {
    await logout();
    navigate("/");
  };


  return (
    <div>
      <h1>Signed in as {user?.username}</h1>
      <button onClick={() => navigate(`/sendmsg`)}>Send Message</button>
      <button onClick={() => navigate(`/notifs`)}>Messages</button>
      <button onClick={handleLogout}>Sign Out</button>
    </div>
  );
}

export { HomePage };


