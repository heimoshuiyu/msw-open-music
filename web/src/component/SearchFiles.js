import { useState, useEffect } from "react";
import FilesTable from "./FilesTable";

function SearchFiles(props) {
  const [files, setFiles] = useState([]);
  const [filename, setFilename] = useState("");
  const [offset, setOffset] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const limit = 10;

  function searchFiles() {
    // check empty filename
    if (filename === "") {
      return;
    }
    setIsLoading(true);
    fetch("/api/v1/search_files", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        filename: filename,
        limit: limit,
        offset: offset,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        const files = data.files ? data.files : [];
        setFiles(files);
      })
      .catch((error) => {
        alert("search_files error: " + error);
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

  useEffect(() => searchFiles(), [offset]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="page">
      <h3>Search Files</h3>
      <div className="search_toolbar">
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
        <button
          onClick={() => {
            searchFiles();
          }}
        >
          {isLoading ? "Loading..." : "Search"}
        </button>
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
