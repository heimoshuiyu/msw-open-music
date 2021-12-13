import { useNavigate } from "react-router";
import Database from "./Database";

function Manage(props) {
  let navigate = useNavigate();

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
        <button
          onClick={() => {
            fetch("/api/v1/logout")
              .then((res) => res.json())
              .then((data) => {
                if (data.error) {
                  alert(data.error);
                } else {
                  props.setUser(data.user);
                }
              });
          }}
        >
          Logout
        </button>
      )}
      <hr />
      <button onClick={() => navigate("/manage/tags")}>Tags</button>
      <Database />
    </div>
  );
}

export default Manage;
