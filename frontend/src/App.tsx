import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import SignUp from "./pages/auth/SignUp";
import Login from "./pages/auth/Login";
import Home from "./pages/Home";
import Chat from "./pages/chat/Chat";
import NotFound from "./pages/NotFound"
import PrivateRoutes from "./utils/PrivateRoutes";

// https://create-react-app.dev/docs/proxying-api-requests-in-development
// https://www.typescripttutorial.net/typescript-tutorial/typescript-access-modifiers/

class App extends React.Component<any,any> {
  constructor(props: any) {
    super(props);
  }

  render() {
    return (
      <div className="App">
        <BrowserRouter>
          <Routes>
            <Route element={<PrivateRoutes/>}>
              <Route element={<Home/>} path="/" />
              <Route element={<Chat/>} path="/chat" />
              <Route element={<Chat/>} path="/chat/:username" />

              <Route element={<NotFound/>} path="/*" />
            </Route>

            <Route path="/signup" element={<SignUp />} />
            <Route path="/login" element={<Login />} />
         </Routes>
        </BrowserRouter>
      </div>
    );
  }
}

export default App;
