import * as React from 'react';
import { useNavigate } from "react-router-dom";
import { useContext, useState } from "react";
import { tr, Tr, langCodeContext } from "../translate";

function Register() {
  let navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [password2, setPassword2] = useState("");
  const [role, setRole] = useState("");
  const { langCode } = useContext(langCodeContext);

  function register() {
    if (!username || !password || !password2 || !role) {
      alert(tr("Please fill out all fields", langCode));
    } else if (password !== password2) {
      alert(tr("Password do not match", langCode));
    } else {
      fetch("/api/v1/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          username: username,
          password: password,
          role: parseInt(role),
        }),
      })
        .then((res) => res.json())
        .then((data) => {
          if (data.error) {
            alert(data.error);
          } else {
            navigate("/manage/login");
          }
        });
    }
  }

  return (
    <div className="page">
      <h2>{Tr("Register")}</h2>
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
      />
      <label htmlFor="password2">{Tr("Confirm Password")}</label>
      <input
        type="password"
        id="password2"
        value={password2}
        onChange={(e) => setPassword2(e.target.value)}
        onKeyPress={(e) => {
          if (e.key === "Enter") {
            e.preventDefault();
            register();
          }
        }}
      />
      <label htmlFor="role">{Tr("Role")}</label>
      <select value={role} onChange={(e) => setRole(e.target.value)}>
        <option value="">{tr("Select a role", langCode)}</option>
        <option value="2">{tr("User", langCode)}</option>
        <option value="1">{tr("Admin", langCode)}</option>
      </select>
      <button onClick={register}>{Tr("Register")}</button>
    </div>
  );
}

export default Register;
