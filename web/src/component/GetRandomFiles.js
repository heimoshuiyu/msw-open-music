import { useContext, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQuery } from "./Common";
import FilesTable from "./FilesTable";
import { Tr, tr, langCodeContext } from "../translate";

function GetRandomFiles(props) {
  const [files, setFiles] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [tags, setTags] = useState([]);
  const navigator = useNavigate();
  const query = useQuery();
  const selectedTag = query.get("t") || "";
  const { langCode } = useContext(langCodeContext);

  const fetchRandomFiles = async () => {
    const resp = await fetch("/api/v1/get_random_files");
    const json = await resp.json();
    return json.files;
  };

  async function getRandomFiles() {
    setIsLoading(true);
    fetchRandomFiles()
      .then((data) => {
        setFiles(data);
      })
      .catch((error) => {
        alert("get_random_files error: " + error);
      })
      .finally(() => {
        setIsLoading(false);
      });
  }

  const fetchRandomFilesWithTag = async (selectedTag) => {
    const resp = await fetch("/api/v1/get_random_files_with_tag", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: parseInt(selectedTag),
      }),
    });
    const json = await resp.json();
    return json.files;
  };

  function getRandomFilesWithTag() {
    setIsLoading(true);
    fetchRandomFilesWithTag(selectedTag)
      .then((files) => {
        setFiles(files);
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
          {isLoading ? Tr("Loading...") : Tr("Refresh")}
        </button>
        <select
          className="tag_select"
          onChange={(event) => {
            navigator(`/?t=${event.target.value}`);
          }}
          value={selectedTag}
        >
          <option value="">{tr("All", langCode)}</option>
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
