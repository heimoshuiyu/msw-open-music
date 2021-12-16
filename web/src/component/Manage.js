import { useNavigate } from "react-router";
import Database from "./Database";

function Manage(props) {
  let navigate = useNavigate();

  return (
    <div className="page">
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
        <div className="horizontal">
          <button
            onClick={() => {
              navigate(`/manage/users/${props.user.id}`);
            }}
          >
            Edit
          </button>
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
        </div>
      )}
      <hr />
      <div className="horizontal">
        <button onClick={() => navigate("/manage/tags")}>Tags</button>
        <button onClick={() => navigate("/manage/users")}>Users</button>
        <button onClick={() => navigate("/manage/feedbacks")}>Feedbacks</button>
      </div>
      <Database />
    </div>
  );
}

export default Manage;
