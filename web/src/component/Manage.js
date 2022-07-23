import { useNavigate } from "react-router";
import Database from "./Database";

import { Tr, langCodeContext, LANG_OPTIONS } from "../translate";
import { useContext } from "react";

function Manage(props) {
  let navigate = useNavigate();
  const { langCode, setLangCode } = useContext(langCodeContext);

  return (
    <div className="page">
      <h2>{Tr("Manage")}</h2>
      <p>
        {Tr("Hi")}, {props.user.username}
      </p>

      <select
        value={langCode}
        onChange={(event) => {
          setLangCode(event.target.value);
        }}
      >
        {Object.keys(LANG_OPTIONS).map((code) => {
          const langOption = LANG_OPTIONS[code];
          return (
            <option value={code} key={code}>
              {langOption.name}
            </option>
          );
        })}
      </select>

      {props.user.role === 0 && (
        <div>
          <button
            onClick={() => {
              navigate("/manage/login");
            }}
          >
            {Tr("Login")}
          </button>
          <button
            onClick={() => {
              navigate("/manage/register");
            }}
          >
            {Tr("Register")}
          </button>
        </div>
      )}
      {props.user.role !== 0 && (
        <div className="horizontal">
          <button
            onClick={() => {
              navigate(`/manage/users/${props.user.id}`);
            }}
          >
            {Tr("Profile")}
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
            {Tr("Logout")}
          </button>
        </div>
      )}
      <hr />
      <div className="horizontal">
        <button onClick={() => navigate("/manage/tags")}>{Tr("Tags")}</button>
        <button onClick={() => navigate("/manage/users")}>{Tr("Users")}</button>
        <button onClick={() => navigate("/manage/feedbacks")}>
          {Tr("Feedbacks")}
        </button>
      </div>
      <Database />
      <p>
        <a
          href="https://github.com/heimoshuiyu/msw-open-music"
          target="_blank"
          rel="noreferrer"
        >
          {Tr("View source code on Github")}
        </a>
      </p>
    </div>
  );
}

export default Manage;
