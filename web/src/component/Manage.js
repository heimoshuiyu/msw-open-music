import { useState } from "react";
import { useNavigate } from "react-router";

function Manage(props) {
  let navigate = useNavigate();

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
        pattern: [".wav", ".mp3"],
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
      <p>Hi, {props.user.username}</p>
      {props.user.role === 0 && (
        <button
          onClick={() => {
            navigate("/manage/login");
          }}
        >
          Login
        </button>
      )}
      {props.user.role !== 0 && (
        <button onClick={() => props.setUser({})}>Logout</button>
      )}
      <hr />
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
