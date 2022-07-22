import { useState, useEffect, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { useQuery } from "./Common";
import FilesTable from "./FilesTable";
import { Tr, tr, langCodeContext } from "../translate";

function SearchFiles(props) {
  const navigator = useNavigate();
  const [files, setFiles] = useState([]);
  const query = useQuery();
  const filename = query.get("q") || "";
  const [filenameInput, setFilenameInput] = useState(filename);
  const offset = parseInt(query.get("o")) || 0;
  const [isLoading, setIsLoading] = useState(false);
  const limit = 10;
  const { langCode } = useContext(langCodeContext);

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
    navigator(`/files?q=${filenameInput}&o=${offset + limit}`);
  }

  function lastPage() {
    const offsetValue = offset - limit;
    if (offsetValue < 0) {
      return;
    }
    navigator(`/files?q=${filenameInput}&o=${offsetValue}`);
  }

  useEffect(() => searchFiles(), [offset, filename]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="page">
      <h3>{Tr("Search Files")}</h3>
      <div className="search_toolbar">
        <input
          onChange={(event) => setFilenameInput(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              navigator(`/files?q=${filenameInput}&o=0`);
            }
          }}
          type="text"
          placeholder={tr("Enter filename", langCode)}
          value={filenameInput}
        />
        <button
          onClick={() => {
            navigator(`/files?q=${filenameInput}&o=0`);
          }}
        >
          {isLoading ? Tr("Loading...") : Tr("Search")}
        </button>
        <button onClick={lastPage}>{Tr("Last page")}</button>
        <button disabled>
          {offset} - {offset + files.length}
        </button>
        <button onClick={nextPage}>{Tr("Next page")}</button>
      </div>
      <FilesTable setPlayingFile={props.setPlayingFile} files={files} />
    </div>
  );
}

export default SearchFiles;
