import { useContext, useEffect, useState } from "react";
import { useParams } from "react-router";
import FilesTable from "./FilesTable";
import { Tr, tr, langCodeContext } from "../translate";

function Share(props) {
  let params = useParams();
  const { langCode } = useContext(langCodeContext);

  const [file, setFile] = useState({});
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
        setFile(data);
      })
      .catch((error) => {
        alert("get_file_info error: " + error);
      });
  }, [params]);

  // change title
  useEffect(() => {
    const oldTitle = document.title;

    document.title = `${tr("Share", langCode)}🎵: ${
      file.filename
    } - MSW Open Music`;

    // set title back
    return () => {
      document.title = oldTitle;
    };
  }, [file]);

  return (
    <div className="page">
      <h3>{Tr("Share with others!")}</h3>
      <p>
        {Tr("Share link")}:{" "}
        <a href={window.location.href}>{window.location.href}</a>
      </p>
      <p>
        👇 {Tr("Click the filename below to enjoy music!")}
        <br />
      </p>
      <FilesTable setPlayingFile={props.setPlayingFile} files={[file]} />
    </div>
  );
}

export default Share;
