import { useNavigate } from "react-router-dom";
import { useState } from "react";

function Register(props) {
  let navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [password2, setPassword2] = useState("");
  const [role, setRole] = useState("");
  return (
    <div>
      <h2>Register</h2>
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
      />
      <label htmlFor="password2">Confirm Password</label>
      <input
        type="password"
        id="password2"
        value={password2}
        onChange={(e) => setPassword2(e.target.value)}
      />
      <label htmlFor="role">Role</label>
      <select value={role} onChange={(e) => setRole(e.target.value)}>
        <option value="">Select a role</option>
        <option value="2">User</option>
        <option value="1">Admin</option>
      </select>
      <button
        onClick={() => {
          if (!username || !password || !password2 || !role) {
            alert("Please fill out all fields");
          } else if (password !== password2) {
            alert("Passwords do not match");
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
                  props.setUser(data.user);
                  navigate("/login");
                }
              });
          }
        }}
      >
        Register
      </button>
    </div>
  );
}

export default Register;
