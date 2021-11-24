import { useState, useEffect } from "react";
import FilesTable from "./FilesTable";

function SearchFiles(props) {
  const [files, setFiles] = useState([]);
  const [filename, setFilename] = useState("");
  const [offset, setOffset] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const limit = 10;

  function searchFiles() {
    if (
      filename === "" &&
      (props.folder === undefined || props.folder.id === undefined)
    ) {
      return;
    }
    const folder = props.folder ? props.folder : {};
    const url = folder.id
      ? "/api/v1/get_files_in_folder"
      : "/api/v1/search_files";
    setIsLoading(true);
    fetch(url, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        filename: filename,
        limit: limit,
        offset: offset,
        folder_id: folder.id,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        const files = data.files ? data.files : [];
        setFiles(files);
      })
      .catch((error) => {
        alert("get_files_in_folder error: " + error);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

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

  useEffect(() => searchFiles(), [offset, props.folder]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="page">
      <h3>Search Files</h3>
      <div className="search_toolbar">
        {!props.folder && (
          <input
            onChange={(event) => setFilename(event.target.value)}
            onKeyDown={(event) => {
              if (event.key === "Enter") {
                searchFiles();
              }
            }}
            type="text"
            placeholder="Enter filename"
          />
        )}
        <button
          disabled={!!props.folder}
          onClick={() => {
            searchFiles();
          }}
        >
          {isLoading ? "Loading..." : "Search"}
        </button>
        {props.folder && props.folder.foldername && (
          <button onClick={searchFiles}>{props.folder.foldername}</button>
        )}
        <button onClick={lastPage}>Last page</button>
        <button disabled>
          {offset} - {offset + files.length}
        </button>
        <button onClick={nextPage}>Next page</button>
      </div>
      <FilesTable setPlayingFile={props.setPlayingFile} files={files} />
    </div>
  );
}

export default SearchFiles;
