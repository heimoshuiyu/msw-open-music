import * as React from 'react';
import {useParams} from "react-router";
import {useState, useEffect} from "react";
import {useNavigate} from "react-router-dom";
import {useQuery} from "./Common";
import FilesTable from "./FilesTable";
import {Tr} from "../translate";

function FilesInFolder(props) {
  let params = useParams();
  const query = useQuery();
  const navigator = useNavigate();
  const [files, setFiles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const offset = parseInt(query.get("o")) || 0;
  const [newFoldername, setNewFoldername] = useState("");
  const [folderPath, setFolderPath] = useState("");
  const limit = 10;

  function refresh() {
    setIsLoading(true);
    fetch("/api/v1/get_files_in_folder", {
      method: "POST",
      headers: {"Content-Type": "application/json"},
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
          setFolderPath(data.folder);
          if (data.files.length > 0) {
            setNewFoldername(data.files[0].foldername);
          }
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
    navigator(`/folders/${params.id}?o=${offset + limit}`);
  }

  function lastPage() {
    const offsetValue = offset - limit;
    if (offsetValue < 0) {
      return;
    }
    navigator(`/folders/${params.id}?o=${offsetValue}`);
  }

  function updateFoldername() {
    setIsLoading(true);
    fetch("/api/v1/update_foldername", {
      method: "POST",
      headers: {"Content-Type": "application/json"},
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
      headers: {"Content-Type": "application/json"},
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
      <h3>{Tr("Files in Folder")}</h3>
      <div className="search_toolbar">
        <button onClick={lastPage}>{Tr("Last page")}</button>
        <button disabled>
          {isLoading
            ? Tr("Loading...")
            : `${offset} - ${offset + files.length}`}
        </button>
        <button onClick={nextPage}>{Tr("Next page")}</button>
      </div>
      <FilesTable setPlayingFile={props.setPlayingFile} files={files} />
      <span>{folderPath}</span>
      <div>
        <input
          type="text"
          value={newFoldername}
          onChange={(e) => setNewFoldername(e.target.value)}
        />
        <div>
          <button onClick={() => updateFoldername()}>{Tr("Save")}</button>
          <button onClick={() => resetFoldername()}>{Tr("Reset")}</button>
        </div>
      </div>
    </div>
  );
}

export default FilesInFolder;
