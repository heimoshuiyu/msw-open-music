import { useEffect, useState } from "react";
import { useParams } from "react-router";
import FoldersTable from "./FoldersTable";
import SearchFiles from "./SearchFiles";

function SearchFolders(props) {
  const [foldername, setFoldername] = useState("");
  const [folders, setFolders] = useState([]);
  const [folder, setFolder] = useState({});
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
        setIsLoading(false);
        let folders;
        if (data.folders) {
          folders = data.folders;
        } else {
          folders = [];
        }
        setFolders(folders);
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

  function viewFolder(folder) {
    setFolder(folder);
  }

  let params = useParams();
  useEffect(() => searchFolder(), [offset]); // eslint-disable-line react-hooks/exhaustive-deps
  useEffect(() => {
    if (params.id !== undefined) {
      setFolder({ id: parseInt(params.id) });
    }
  }, [params.id]);

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
        <button onClick={searchFolder}>
          {isLoading ? "Loading..." : "Search"}
        </button>
        <button onClick={lastPage}>Last page</button>
        <button disabled>
          {offset} - {offset + limit}
        </button>
        <button onClick={nextPage}>Next page</button>
      </div>
      <FoldersTable viewFolder={viewFolder} folders={folders} />
      <SearchFiles setPlayingFile={props.setPlayingFile} folder={folder} />
    </div>
  );
}

export default SearchFolders;
