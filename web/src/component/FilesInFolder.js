import { useParams } from "react-router";
import { useState, useEffect } from "react";
import FilesTable from "./FilesTable";

function FilesInFolder(props) {
  let params = useParams();
  const [files, setFiles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [offset, setOffset] = useState(0);
  const [newFoldername, setNewFoldername] = useState("");
  const limit = 10;

  function refresh() {
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
        if (data.error) {
          alert(data.error);
        } else {
          setFiles(data.files);
          setNewFoldername(data.files[0].foldername);
        }
      })
      .catch((error) => alert(error))
      .finally(() => {
        setIsLoading(false);
      });
  }

  useEffect(() => {
    refresh();
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

  function updateFoldername() {
    setIsLoading(true);
    fetch("/api/v1/update_foldername", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        id: parseInt(params.id),
        foldername: newFoldername,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          refresh();
        }
      })
      .catch((error) => alert(error))
      .finally(() => {
        setIsLoading(false);
      });
  }

  function resetFoldername() {
    setIsLoading(true);
    fetch("/api/v1/reset_foldername", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        id: parseInt(params.id),
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          refresh();
        }
      })
      .catch((error) => alert(error))
      .finally(() => {
        setIsLoading(false);
      });
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
      <div>
        <input
          type="text"
          value={newFoldername}
          onChange={(e) => setNewFoldername(e.target.value)}
        />
        <div>
          <button onClick={() => updateFoldername()}>Save</button>
          <button onClick={() => resetFoldername()}>Reset</button>
        </div>
      </div>
    </div>
  );
}

export default FilesInFolder;
