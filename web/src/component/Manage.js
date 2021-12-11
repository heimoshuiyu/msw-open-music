import { useState } from "react";
import { useNavigate } from "react-router";

function Manage(props) {
  let navigate = useNavigate();

  const [walkPath, setWalkPath] = useState("");
  const [patternString, setPatternString] = useState("");

  function updateDatabase() {
    // split pattern string into array
    let patternArray = patternString.split(" ");
    // remove whitespace from array
    patternArray = patternArray.map((item) => item.trim());
    // remove empty strings from array
    patternArray = patternArray.filter((item) => item !== "");
    // add dot before item array
    patternArray = patternArray.map((item) => "." + item);

    fetch("/api/v1/walk", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        root: walkPath,
        pattern: patternArray,
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
      <h3>Update Database</h3>
      <input
        type="text"
        value={walkPath}
        placeholder="walk path"
        onChange={(e) => setWalkPath(e.target.value)}
      />
      <input
        type="text"
        value={patternString}
        placeholder="pattern wav flac mp3"
        onChange={(e) => setPatternString(e.target.value)}
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
