import { useEffect, useState } from "react";
import { useParams } from "react-router";
import FilesTable from "./FilesTable";

function Share(props) {
  let params = useParams();
  const [file, setFile] = useState([]);
  useEffect(() => {
    fetch("/api/v1/get_file_info", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        id: parseInt(params.id),
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        setFile([data]);
      })
      .catch((error) => {
        alert("get_file_info error: " + error);
      });
  }, [params]);
  return (
    <div className="page">
      <h3>Share with others!</h3>
      <p>
        👇 Click the filename below to enjoy music!
        <br />
      </p>
      <p>
        Share link: <a href={window.location.href}>{window.location.href}</a>
      </p>
      <FilesTable setPlayingFile={props.setPlayingFile} files={file} />
    </div>
  );
}

export default Share;
