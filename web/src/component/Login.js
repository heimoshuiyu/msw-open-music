import { useNavigate } from "react-router-dom";
import { useState } from "react";

function Login(props) {
  let navigate = useNavigate();
  let [username, setUsername] = useState("");
  let [password, setPassword] = useState("");

  function login() {
    if (!username || !password) {
      alert("Please enter username and password");
      return;
    }
    fetch("/api/v1/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        username,
        password,
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
          return;
        }
        props.setUser(data.user);
        navigate("/");
      });
  }

  return (
    <div className="page">
      <h2>Login</h2>
      <label htmlFor="username">Username</label>
      <input
        type="text"
        id="username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <label htmlFor="password">Password</label>
      <input
        type="password"
        id="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        onKeyPress={(e) => {
          if (e.key === "Enter") {
            e.preventDefault();
            login();
          }
        }}
      />
      <span>
        <button onClick={login}>Login</button>
        <button
          onClick={() => {
            navigate("/manage/register");
          }}
        >
          Register
        </button>
      </span>
    </div>
  );
}

export default Login;
