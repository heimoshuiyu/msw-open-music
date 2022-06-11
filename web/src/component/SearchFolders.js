import { useEffect, useState } from "react";
import FoldersTable from "./FoldersTable";

function SearchFolders() {
  const [foldername, setFoldername] = useState("");
  const [folders, setFolders] = useState([]);
  const [offset, setOffset] = useState(0);
  const [isLoading, setIsLoading] = useState(false);
  const limit = 10;

  function searchFolder() {
    if (foldername === "") {
      return;
    }
    setIsLoading(true);
    fetch("/api/v1/search_folders", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        foldername: foldername,
        limit: limit,
        offset: offset,
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        setFolders(data.folders ? data.folders : []);
      })
      .catch((error) => {
        alert("search_folders error: " + error);
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

  useEffect(() => searchFolder(), [offset]); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <div className="page">
      <h3>Search Folders</h3>
      <div className="search_toolbar">
        <input
          onChange={(event) => setFoldername(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              searchFolder();
            }
          }}
          type="text"
          placeholder="Enter folder name"
        />
        <button
          onClick={() => {
            setOffset(0);
            searchFolder();
          }}
        >
          {isLoading ? "Loading..." : "Search"}
        </button>
        <button onClick={lastPage}>Last page</button>
        <button disabled>
          {offset} - {offset + limit}
        </button>
        <button onClick={nextPage}>Next page</button>
      </div>
      <FoldersTable folders={folders} />
    </div>
  );
}

export default SearchFolders;
