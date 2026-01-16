import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from "./contexts/AuthContext";
import { Msg } from './pages/GetMsg';
import {HomePage} from './pages/Home';
import { SendMsg } from './pages/SendMsg';
import { Notif } from './pages/GetNotif';
import { SignUp } from './pages/Signup';  
import { Login } from './pages/LogIn';
import { ProtectedRoute } from './components/ProtectedRoute';
import { Verification } from './pages/Verification';


import './App.css';


function App() {
 
  return (
<AuthProvider>
  {/* <BrowserRouter> */}
    <Router>
      <Routes>
        <Route
          path={`/`}
          element={<Login />}
        />
        <Route
          path={`/verify`}
          element={<Verification />}
        />
       
       <Route
            path="/home"
            element={
              <ProtectedRoute>
                <HomePage />
              </ProtectedRoute>
            }
        />

        <Route
          path={`/signup`}
          element={<SignUp />}
        />

        <Route
            path="/sendmsg"
            element={
              <ProtectedRoute>
                <SendMsg />
              </ProtectedRoute>
            }
          />

       


        <Route
            path="/notifs"
            element={
              <ProtectedRoute>
                <Notif />
              </ProtectedRoute>
            }
          />

        

        <Route
            path="/msg"
            element={
              <ProtectedRoute>
                <Msg />
              </ProtectedRoute>
            }
          />

        

    
          
      
      </Routes>
    </Router>
  {/* </BrowserRouter> */}
</AuthProvider>
  );
}

export default {App};




