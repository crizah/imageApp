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
  <div className="relative min-h-screen bg-slate-950 overflow-hidden flex items-center justify-center">
    {/* Background glows */}
    <div className="absolute top-0 left-0 w-96 h-96 bg-purple-500 rounded-full blur-3xl opacity-20 animate-pulse"></div>
    <div className="absolute bottom-0 right-0 w-96 h-96 bg-cyan-500 rounded-full blur-3xl opacity-20 animate-pulse"></div>

    <div className="relative z-10 w-full max-w-md px-6">
      <div className="bg-slate-900/40 backdrop-blur-xl border border-slate-800/50 rounded-2xl shadow-2xl p-8 transition-all hover:scale-[1.02]">

        <h1 className="text-2xl font-bold text-center bg-gradient-to-r from-cyan-400 via-purple-400 to-pink-400 bg-clip-text text-transparent mb-2">
          Welcome back
        </h1>

        <p className="text-center text-gray-400 mb-8">
          Signed in as <span className="text-white font-semibold">{user?.username}</span>
        </p>

        <div className="space-y-4">
          <button
            onClick={() => navigate(`/sendmsg`)}
            className="w-full py-3 rounded-lg font-semibold text-white bg-gradient-to-r from-cyan-500 to-purple-500 hover:from-cyan-600 hover:to-purple-600 transition-all"
          >
            Send Message
          </button>

          <button
            onClick={() => navigate(`/notifs`)}
            className="w-full py-3 rounded-lg font-semibold text-white bg-gradient-to-r from-purple-500 to-pink-500 hover:from-purple-600 hover:to-pink-600 transition-all"
          >
            Messages
          </button>

          <button
            onClick={handleLogout}
            className="w-full py-3 rounded-lg font-semibold text-red-400 border border-red-500/40 hover:bg-red-500/10 transition-all"
          >
            Sign Out
          </button>
        </div>
      </div>
    </div>
  </div>
);
}

export { HomePage };


