import { useEffect, useState } from "react";
import FilesTable from "./FilesTable";

function GetRandomFiles(props) {
  const [files, setFiles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);

  function refresh(setFiles) {
    setIsLoading(true);
    fetch("/api/v1/get_random_files")
      .then((response) => response.json())
      .then((data) => {
        setFiles(data.files);
      })
      .catch((error) => {
        alert("get_random_files error: " + error);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  useEffect(() => {
    refresh(setFiles);
  }, []);
  return (
    <div className="page">
      <div className="search_toolbar">
        <button className="refresh" onClick={() => refresh(setFiles)}>
          {isLoading ? "Loading..." : "Refresh"}
        </button>
      </div>
      <FilesTable setPlayingFile={props.setPlayingFile} files={files} />
    </div>
  );
}

export default GetRandomFiles;
