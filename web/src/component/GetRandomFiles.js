import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQuery } from "./Common";
import FilesTable from "./FilesTable";

function GetRandomFiles(props) {
  const [files, setFiles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [tags, setTags] = useState([]);
  const navigator = useNavigate();
  const query = useQuery();
  const selectedTag = query.get("t") || "";

  function getRandomFiles() {
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

  function getRandomFilesWithTag() {
    setIsLoading(true);
    fetch("/api/v1/get_random_files_with_tag", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: parseInt(selectedTag),
      }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setFiles(data.files);
        }
      })
      .catch((error) => {
        alert("get_random_files_with_tag error: " + error);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  function refresh() {
    if (selectedTag === "") {
      getRandomFiles();
    } else {
      getRandomFilesWithTag();
    }
  }

  function getTags() {
    fetch("/api/v1/get_tags")
      .then((response) => response.json())
      .then((data) => {
        setTags(data.tags);
      })
      .catch((error) => {
        alert("get_tags error: " + error);
      });
  }

  useEffect(() => {
    getTags();
  }, []);

  useEffect(() => {
    refresh();
  }, [selectedTag]);

  return (
    <div className="page">
      <div className="search_toolbar">
        <button className="refresh" onClick={() => refresh(setFiles)}>
          {isLoading ? "Loading..." : "Refresh"}
        </button>
        <select
          className="tag_select"
          onChange={(event) => {
            navigator(`/?t=${event.target.value}`);
          }}
          value={selectedTag}
        >
          <option value="">All</option>
          {tags.map((tag) => (
            <option key={tag.id} value={tag.id}>
              {tag.name}
            </option>
          ))}
        </select>
      </div>
      <FilesTable setPlayingFile={props.setPlayingFile} files={files} />
    </div>
  );
}

export default GetRandomFiles;
