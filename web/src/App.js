import { HashRouter as Router, Routes, Route, NavLink } from "react-router-dom";
import "./App.css";

import GetRandomFiles from "./component/GetRandomFiles";
import SearchFiles from "./component/SearchFiles";
import SearchFolders from "./component/SearchFolders";
import FilesInFolder from "./component/FilesInFolder";
import Manage from "./component/Manage";
import ManageUser from "./component/ManageUser";
import FileInfo from "./component/FileInfo";
import Share from "./component/Share";
import Login from "./component/Login";
import Register from "./component/Register";
import Tags from "./component/Tags";
import EditTag from "./component/EditTag";
import EditReview from "./component/EditReview";
import AudioPlayer from "./component/AudioPlayer";
import UserStatus from "./component/UserStatus";
import ReviewPage from "./component/ReviewPage";
import UserProfile from "./component/UserProfile";
import FeedbackPage from "./component/FeedbackPage";
import { useState } from "react";

function App() {
  const [playingFile, setPlayingFile] = useState({});
  const [user, setUser] = useState({});
  return (
    <div className="base">
      <Router>
        <header className="header">
          <h3 className="title">
            <img src="favicon.png" alt="logo" className="logo" />
            <span className="title-text">MSW Open Music Project</span>
            <UserStatus user={user} setUser={setUser} />
          </h3>
          <nav className="nav">
            <NavLink to="/" className="nav-link">
              Feeling luckly
            </NavLink>
            <NavLink to="/files" className="nav-link">
              Files
            </NavLink>
            <NavLink to="/folders" className="nav-link">
              Folders
            </NavLink>
            <NavLink to="/manage" className="nav-link">
              Manage
            </NavLink>
          </nav>
        </header>
        <main>
          <Routes>
            <Route
              index
              path="/"
              element={<GetRandomFiles setPlayingFile={setPlayingFile} />}
            />
            <Route
              path="/files"
              element={<SearchFiles setPlayingFile={setPlayingFile} />}
            />
            <Route
              path="/folders"
              element={<SearchFolders setPlayingFile={setPlayingFile} />}
            />
            <Route
              path="/folders/:id"
              element={<FilesInFolder setPlayingFile={setPlayingFile} />}
            />
            <Route
              path="/manage"
              element={<Manage user={user} setUser={setUser} />}
            />
            <Route
              path="/manage/feedbacks"
              element={<FeedbackPage user={user} />}
            />
            <Route
              path="/manage/login"
              element={<Login user={user} setUser={setUser} />}
            />
            <Route
              path="/manage/register"
              element={<Register user={user} setUser={setUser} />}
            />
            <Route
              path="/manage/tags"
              element={<Tags user={user} />}
            />
            <Route
              path="/manage/tags/:id"
              element={<EditTag user={user} />}
            />
            <Route
              path="/manage/reviews/:id"
              element={<EditReview user={user} />}
            />
            <Route
              path="/manage/users"
              element={<ManageUser user={user} setUser={setUser} />}
            />
            <Route
              path="/manage/users/:id"
              element={<UserProfile user={user} setUser={setUser} />}
            />
            <Route
              path="/files/:id"
              element={<FileInfo setPlayingFile={setPlayingFile} />}
            />
            <Route
              path="/files/:id/share"
              element={<Share setPlayingFile={setPlayingFile} />}
            />
            <Route
              path="/files/:id/review"
              element={<ReviewPage user={user} setPlayingFile={setPlayingFile} />}
            />
          </Routes>
        </main>
        <footer>
          <AudioPlayer
            playingFile={playingFile}
            setPlayingFile={setPlayingFile}
          />
        </footer>
      </Router>
    </div>
  );
}

export default App;
