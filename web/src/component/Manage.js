import { useState } from "react";

function Manage() {
  const [token, setToken] = useState("");
  const [walkPath, setWalkPath] = useState("");

  function updateDatabase() {
    fetch("/api/v1/walk", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        token: token,
        root: walkPath,
        pattern: [".wav", ".mp3", ".flac", ".ogg", ".aac", ".mka"],
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
      });
  }

  return (
    <div>
      <h2>Manage</h2>
      <input
        type="text"
        value={token}
        placeholder="token"
        onChange={(e) => setToken(e.target.value)}
      />
      <input
        type="text"
        value={walkPath}
        placeholder="walk path"
        onChange={(e) => setWalkPath(e.target.value)}
      />
      <button
        onClick={() => {
          updateDatabase();
        }}
      >
        Update Database
      </button>
    </div>
  );
}

export default Manage;
