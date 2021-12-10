import { useParams } from "react-router";
import { useState, useEffect } from "react";
import FilesTable from "./FilesTable";

function FilesInFolder(props) {
  let params = useParams();
  const [files, setFiles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [offset, setOffset] = useState(0);
  const limit = 10;

  useEffect(() => {
    setIsLoading(true);
    fetch("/api/v1/get_files_in_folder", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        folder_id: parseInt(params.id),
        offset: offset,
        limit: limit,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        setFiles(data.files ? data.files : []);
      })
      .catch((error) => alert(error))
      .finally(() => {
        setIsLoading(false);
      });
  }, [params.id, offset]);

  function nextPage() {
    setOffset(offset + limit);
  }

  function lastPage() {
    const offsetValue = offset - limit;
    if (offsetValue < 0) {
      return;
    }
    setOffset(offsetValue);
  }

  return (
    <div className="page">
      <h3>Files in Folder</h3>
      <div className="search_toolbar">
        <button onClick={lastPage}>Last page</button>
        <button disabled>
          {isLoading ? "Loading..." : `${offset} - ${offset + files.length}`}
        </button>
        <button onClick={nextPage}>Next page</button>
      </div>
      <FilesTable setPlayingFile={props.setPlayingFile} files={files} />
    </div>
  );
}

export default FilesInFolder;
