import { useNavigate } from "react-router-dom";
import { useContext, useState } from "react";
import { Tr, tr, langCodeContext } from "../translate";

function Login(props) {
  let navigate = useNavigate();
  let [username, setUsername] = useState("");
  let [password, setPassword] = useState("");
  const { langCode } = useContext(langCodeContext);

  function login() {
    if (!username || !password) {
      alert(tr("Please enter username and password", langCode));
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
      <h2>{Tr("Login")}</h2>
      <label htmlFor="username">{Tr("Username")}</label>
      <input
        type="text"
        id="username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />
      <label htmlFor="password">{Tr("Password")}</label>
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
        <button onClick={login}>{Tr("Login")}</button>
        <button
          onClick={() => {
            navigate("/manage/register");
          }}
        >
          {Tr("Register")}
        </button>
      </span>
    </div>
  );
}

export default Login;
