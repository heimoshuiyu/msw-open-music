import {
  HashRouter as Router,
  Routes,
  Route,
  NavLink,
} from "react-router-dom";
import "./App.css";

import GetRandomFiles from "./component/GetRandomFiles";
import SearchFiles from "./component/SearchFiles";
import SearchFolders from "./component/SearchFolders";
import FilesInFolder from "./component/FilesInFolder";
import Manage from "./component/Manage";
import Share from "./component/Share";
import AudioPlayer from "./component/AudioPlayer";
import { useState } from "react";

function App() {
  const [playingFile, setPlayingFile] = useState({});
  return (
    <div className="base">
      <Router>
        <header className="header">
          <h3 className="title">
            <img src="favicon.png" alt="logo" className="logo" />
            <span className="title-text">MSW Open Music Project</span>
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
            <Route path="/manage" element={<Manage />} />
            <Route
              path="/files/:id/share"
              element={<Share setPlayingFile={setPlayingFile} />}
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
